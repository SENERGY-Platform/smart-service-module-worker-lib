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

package middleware

import (
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/camunda"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/model"
	"log"
	"runtime/debug"
	"sort"
	"strings"
)

func New(handler camunda.Handler, repo VariablesRepo) *Middleware {
	return &Middleware{
		handler: handler,
		repo:    repo,
	}
}

type Middleware struct {
	handler camunda.Handler
	repo    VariablesRepo
}

type VariablesRepo interface {
	GetVariables(processId string) (result map[string]interface{}, err error)
	SetVariables(processId string, changes map[string]interface{}) error
}

func (this *Middleware) Do(task model.CamundaExternalTask) (modules []model.Module, outputs map[string]interface{}, err error) {
	variables, err := this.repo.GetVariables(task.ProcessInstanceId)
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return modules, outputs, err
	}
	inputs := map[string]interface{}{}
	for key, value := range task.Variables {
		inputs[key] = value.Value
	}
	variableChanges, outputs, err := this.RunPreScripts(inputs, variables)
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return modules, outputs, err
	}
	for key, value := range variableChanges {
		variables[key] = value
	}
	task.Variables, err = this.handleReferences(task.Variables, variables)
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return modules, outputs, err
	}
	modules, handlerOutputs, err := this.handler.Do(task)
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return modules, handlerOutputs, err
	}
	for key, value := range handlerOutputs {
		outputs[key] = value
	}
	postVarChanges, postOutputs, err := this.RunPostScripts(inputs, outputs, variables)
	if err != nil {
		log.Println("ERROR:", err)
		debug.PrintStack()
		return modules, handlerOutputs, err
	}
	for key, value := range postVarChanges {
		variableChanges[key] = value
	}
	for key, value := range postOutputs {
		outputs[key] = value
	}
	if len(variableChanges) > 0 {
		err = this.repo.SetVariables(task.ProcessInstanceId, variableChanges)
		if err != nil {
			log.Println("ERROR:", err)
			debug.PrintStack()
			return modules, outputs, err
		}
	}
	return modules, outputs, nil
}

func (this *Middleware) Undo(modules []model.Module, reason error) {
	this.handler.Undo(modules, reason)
}

const PreScriptPrefix = "prescript"

func (this *Middleware) RunPreScripts(inputs map[string]interface{}, variables map[string]interface{}) (variableChanges map[string]interface{}, outputs map[string]interface{}, err error) {
	return this.RunScripts(PreScriptPrefix, inputs, nil, variables)
}

const PostScriptPrefix = "postscript"

func (this *Middleware) RunPostScripts(inputs map[string]interface{}, existingOutputs map[string]interface{}, variables map[string]interface{}) (variableChanges map[string]interface{}, outputs map[string]interface{}, err error) {
	return this.RunScripts(PostScriptPrefix, inputs, existingOutputs, variables)
}

type KeyValue struct {
	Key   string
	Value string
}

func (this *Middleware) RunScripts(prefix string, inputs map[string]interface{}, existingOutputs map[string]interface{}, variables map[string]interface{}) (variableChanges map[string]interface{}, outputs map[string]interface{}, err error) {
	scriptsKv := []KeyValue{}
	for name, value := range inputs {
		if str, ok := value.(string); ok && strings.HasPrefix(name, prefix) {
			scriptsKv = append(scriptsKv, KeyValue{
				Key:   name,
				Value: str,
			})
		}
	}
	sort.Slice(scriptsKv, func(i, j int) bool {
		return scriptsKv[i].Key < scriptsKv[j].Key
	})

	scripts := []string{}
	for _, script := range scriptsKv {
		scripts = append(scripts, script.Value)
	}
	script := strings.Join(scripts, "")
	scriptEnv := NewScriptEnvWithOutputs(variables, inputs, existingOutputs)
	err = runScript(script, scriptEnv)
	if err != nil {
		return variableChanges, outputs, err
	}
	return scriptEnv.GetUpdatedVariables(), scriptEnv.GetOutputs(), nil
}
