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
}

func (this *Middleware) Do(task model.CamundaExternalTask) (modules []model.Module, outputs map[string]interface{}, err error) {
	variables, err := this.repo.GetVariables(task.ProcessInstanceId)
	if err != nil {
		return modules, outputs, err
	}
	task.Variables, err = this.handleReferences(task.Variables, variables)
	if err != nil {
		return modules, outputs, err
	}
	return this.handler.Do(task)
}

func (this *Middleware) Undo(modules []model.Module, reason error) {
	this.handler.Undo(modules, reason)
}
