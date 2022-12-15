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
	"github.com/dop251/goja"
)

func NewScriptEnv(variables map[string]interface{}, inputs map[string]interface{}) *ScriptEnv {
	return &ScriptEnv{
		vm:               nil,
		Variables:        variables,
		VariablesUpdates: map[string]interface{}{},
		Inputs:           inputs,
		Outputs:          map[string]interface{}{},
	}
}

type ScriptEnv struct {
	vm               *goja.Runtime
	Variables        map[string]interface{}
	VariablesUpdates map[string]interface{}
	Inputs           map[string]interface{}
	Outputs          map[string]interface{}
}

func (this *ScriptEnv) RegisterRuntime(runtime *goja.Runtime) {
	this.vm = runtime
}

func (this *ScriptEnv) GetOutputs() map[string]interface{} {
	return this.Outputs
}

func (this *ScriptEnv) GetUpdatedVariables() map[string]interface{} {
	return this.VariablesUpdates
}

func (this *ScriptEnv) GetEnvironment() map[string]interface{} {
	return map[string]interface{}{
		"io":      NewIoScriptEnv(this),
		"inputs":  NewInputsScriptEnv(this),
		"outputs": NewOutputsScriptEnv(this),
	}
}

func (this *ScriptEnv) GetVm() *goja.Runtime {
	return this.vm
}
