/*
 * Copyright (c) 2022 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package camunda

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime/debug"
	"sync"
	"time"

	"github.com/SENERGY-Platform/gin-middleware/otelx"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/configuration"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TracerName identifies this library as the source of the spans it creates.
const TracerName = "github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/camunda"

func New(config configuration.Config, smartServiceRepo SmartServiceRepository, handler Handler) *Camunda {
	return &Camunda{
		config:           config,
		handler:          handler,
		smartServiceRepo: smartServiceRepo,
	}
}

func Start(ctx context.Context, wg *sync.WaitGroup, config configuration.Config, smartServiceRepo SmartServiceRepository, handler Handler) {
	New(config, smartServiceRepo, handler).Start(ctx, wg)
}

type Camunda struct {
	config           configuration.Config
	handler          Handler
	smartServiceRepo SmartServiceRepository
}

type SmartServiceRepository interface {
	SendWorkerError(ctx context.Context, task model.CamundaExternalTask, err error) error
	SendWorkerModules(ctx context.Context, modules []model.Module) (result []model.SmartServiceModule, err error)
}

type Handler interface {
	Do(ctx context.Context, task model.CamundaExternalTask) (modules []model.Module, outputs map[string]interface{}, err error)
	Undo(ctx context.Context, modules []model.Module, reason error)
}

func (this *Camunda) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				wg.Done()
				return
			default:
				wait := this.executeNextTasks(ctx)
				if wait {
					duration := time.Duration(this.config.CamundaWorkerWaitDurationInMs) * time.Millisecond
					time.Sleep(duration)
				}
			}
		}
	}()
}

func (this *Camunda) executeNextTasks(ctx context.Context) (wait bool) {
	//no span for the poll itself: it runs continuously and mostly returns nothing,
	//a span per poll would be noise without a task to attach it to.
	tasks, err := retry(this.config.FetchRetries, 0, this.getTasks)
	if err != nil {
		this.config.GetLogger().ErrorContext(ctx, "error on ExecuteNextTasks getTask", "error", err)
		return true
	}
	if len(tasks) == 0 {
		return true
	}
	for _, task := range tasks {
		this.executeTask(ctx, task)
	}
	return false
}

// executeTask handles one fetched task in its own trace.
func (this *Camunda) executeTask(lifetimeCtx context.Context, task model.CamundaExternalTask) {
	//context.WithoutCancel: the task is locked in camunda and already running; a shutdown must not
	//cut it off in the middle. that is the behaviour of this worker today and this change is about
	//telemetry, not about the cancellation semantics of a running task.
	ctx, span := otel.Tracer(TracerName).Start(
		context.WithoutCancel(lifetimeCtx),
		this.config.CamundaWorkerTopic,
		trace.WithAttributes(
			attribute.String("camunda_task_id", task.Id),
			attribute.String("camunda_process_instance_id", task.ProcessInstanceId),
			attribute.String("camunda_activity_id", task.ActivityId),
		),
	)
	defer span.End()

	modules, outputs, err := this.handler.Do(ctx, task)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		repoErr := this.smartServiceRepo.SendWorkerError(ctx, task, err)
		if repoErr == nil {
			_ = this.stopProcessInstance(ctx, task.ProcessInstanceId) //error is sent --> no more retries
		}
		//retry task after lock duration, if stop fails or repoErr != nil
	} else {
		_, err = this.smartServiceRepo.SendWorkerModules(ctx, modules)
		if err != nil {
			//undo module and retry after lock duration
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			this.handler.Undo(ctx, modules, err)
			this.config.GetLogger().ErrorContext(ctx, "error on executeNextTasks getTask", "error", err)
			debug.PrintStack()
		} else {
			err = this.completeTask(ctx, task.Id, outputs)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				this.config.GetLogger().ErrorContext(ctx, "error on executeNextTasks getTask", "error", err, "stack", string(debug.Stack()))
				this.handler.Undo(ctx, modules, err)
				repoErr := this.smartServiceRepo.SendWorkerError(ctx, task, err)
				if repoErr == nil {
					//error is sent --> no more retries
					//if it is a problem with the process we don't want any retries
					//if it is a problem with the process-engine, the stop won't be successful and a future try may succeed
					_ = this.stopProcessInstance(ctx, task.ProcessInstanceId)
				}
			}
		}
	}
}

func (this *Camunda) getTasks() (tasks []model.CamundaExternalTask, err error) {
	fetchRequest := model.CamundaFetchRequest{
		WorkerId: this.config.CamundaWorkerId,
		MaxTasks: this.config.CamundaFetchMaxTasks,
		Topics:   []model.CamundaTopic{{LockDuration: this.config.CamundaLockDurationInMs, Name: this.config.CamundaWorkerTopic}},
	}
	client := http.Client{Timeout: 5 * time.Second}
	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(fetchRequest)
	if err != nil {
		return
	}
	endpoint := this.config.CamundaUrl + "/engine-rest/external-task/fetchAndLock"
	resp, err := client.Post(endpoint, "application/json", b)
	if err != nil {
		return tasks, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		temp, err := io.ReadAll(resp.Body)
		err = errors.New(fmt.Sprintln(endpoint, resp.Status, resp.StatusCode, string(temp), err))
		return tasks, err
	}
	err = json.NewDecoder(resp.Body).Decode(&tasks)
	return
}

func (this *Camunda) completeTask(ctx context.Context, taskId string, outputs map[string]interface{}) (err error) {
	this.config.GetLogger().DebugContext(ctx, "complete task", "taskId", taskId, "outputs", outputs)
	client := http.Client{Timeout: 5 * time.Second}

	variables := map[string]model.CamundaVariable{}
	for key, value := range outputs {
		variables[key] = model.CamundaVariable{Value: value}
	}

	var completeRequest = model.CamundaCompleteRequest{WorkerId: this.config.CamundaWorkerId, Variables: variables}
	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(completeRequest)
	if err != nil {
		return
	}
	request, err := http.NewRequest("POST", this.config.CamundaUrl+"/engine-rest/external-task/"+url.PathEscape(taskId)+"/complete", b)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	err = otelx.InjectContextToRequest(ctx, request)
	if err != nil {
		return err
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	pl, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		this.config.GetLogger().ErrorContext(ctx, "unable to complete task", "statuscode", resp.StatusCode, "response", string(pl))
		return fmt.Errorf("unable to complete task: %v, %v", resp.StatusCode, string(pl))
	} else {
		this.config.GetLogger().DebugContext(ctx, "complete camunda task", "request", completeRequest, "response", string(pl))
	}
	return nil
}

func (this *Camunda) stopProcessInstance(ctx context.Context, id string) (err error) {
	client := &http.Client{Timeout: 5 * time.Second}
	request, err := http.NewRequest("DELETE", this.config.CamundaUrl+"/engine-rest/process-instance/"+url.PathEscape(id)+"?skipIoMappings=true", nil)
	if err != nil {
		return err
	}
	err = otelx.InjectContextToRequest(ctx, request)
	if err != nil {
		return err
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if resp.StatusCode == 200 || resp.StatusCode == 204 {
		return nil
	}
	msg, _ := io.ReadAll(resp.Body)
	err = errors.New("error on delete in engine for /engine-rest/process-instance/" + url.PathEscape(id) + ": " + resp.Status + " " + string(msg))
	return err
}

func retry[T any](attempts int, sleep time.Duration, f func() (T, error)) (result T, err error) {
	if attempts == 0 {
		attempts = 1
	}
	for i := 0; i < attempts; i++ {
		result, err = f()
		if err == nil {
			return
		}
		if sleep > 0 {
			time.Sleep(sleep)
		}
	}
	return result, fmt.Errorf("after %d attempts, last error: %w", attempts, err)
}
