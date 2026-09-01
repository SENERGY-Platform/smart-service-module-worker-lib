/*
 * Copyright (c) 2026 InfAI (CC SES)
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
	"context"
	"errors"
	"testing"

	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/configuration"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/model"
)

type handlerMock struct {
	DoCtx   context.Context
	UndoCtx context.Context
}

func (this *handlerMock) Do(ctx context.Context, task model.CamundaExternalTask) (modules []model.Module, outputs map[string]interface{}, err error) {
	this.DoCtx = ctx
	return []model.Module{{Id: "module-id", ProcesInstanceId: task.ProcessInstanceId}}, nil, nil
}

func (this *handlerMock) Undo(ctx context.Context, modules []model.Module, reason error) {
	this.UndoCtx = ctx
}

type repoMock struct {
	SendModulesErr error
}

func (this *repoMock) SendWorkerError(ctx context.Context, task model.CamundaExternalTask, err error) error {
	return nil
}

func (this *repoMock) SendWorkerModules(ctx context.Context, modules []model.Module) (result []model.SmartServiceModule, err error) {
	return nil, this.SendModulesErr
}

// TestExecuteTaskUndoGetsTaskContext: the rollback of a failed task must stay in the trace of that
// task. before Undo had a ctx parameter the workers had to use context.Background() in it, which
// made exactly the error case its own unconnected trace.
func TestExecuteTaskUndoGetsTaskContext(t *testing.T) {
	type ctxKeyType string
	const ctxKey = ctxKeyType("camunda-test")
	lifetimeCtx := context.WithValue(context.Background(), ctxKey, "lifetime-value")

	handler := &handlerMock{}
	repo := &repoMock{SendModulesErr: errors.New("send modules failed")}

	New(configuration.Config{CamundaWorkerTopic: "test-topic"}, repo, handler).
		executeTask(lifetimeCtx, model.CamundaExternalTask{Id: "task-id", ProcessInstanceId: "process-instance-id"})

	if handler.UndoCtx == nil {
		t.Error("Undo got no context")
		return
	}
	if handler.UndoCtx != handler.DoCtx {
		t.Error("Undo did not get the context of Do")
	}
	if value, _ := handler.UndoCtx.Value(ctxKey).(string); value != "lifetime-value" {
		t.Errorf("Undo context: %#v", handler.UndoCtx.Value(ctxKey))
	}
}
