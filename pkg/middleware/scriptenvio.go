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

import "strings"

type ScriptEnvIo struct {
	env *ScriptEnv
}

func NewIoScriptEnv(env *ScriptEnv) *ScriptEnvIo {
	return &ScriptEnvIo{env: env}
}

// Store value as smart-service instance variable
func (this *ScriptEnvIo) Store(name string, value interface{}) {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	this.env.Variables[name] = value
	this.env.VariablesUpdates[name] = value
}

// Read value of a smart-service instance variable
func (this *ScriptEnvIo) Read(name string) interface{} {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	return this.env.Variables[name]
}

// Exists checks if a smart-service instance variable exists
func (this *ScriptEnvIo) Exists(name string) bool {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	_, exists := this.env.Variables[name]
	return exists
}

// Ref creates a reference to a variable (e.g. "my_var_name" --> "{{.my_var_name}}")
// throws exception if variable is unknown
func (this *ScriptEnvIo) Ref(name string) string {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	_, ok := this.env.Variables[name]
	if !ok {
		panic("unknown variable name")
	}
	return "{{." + name + "}}"
}

var TrowErrorForUnknownDerefName = false

// DerefName returns the name of a smart-service instance variable referenced in parameter ref
func (this *ScriptEnvIo) DerefName(ref string) string {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	if !strings.HasPrefix(ref, "{{.") {
		panic("derefName input is not valid reference: missing '{{.' prefix")
	}
	if !strings.HasSuffix(ref, "}}") {
		panic("derefName input is not valid reference: missing '}}' suffix")
	}
	name := strings.TrimPrefix(strings.TrimSuffix(ref, "}}"), "{{.")
	if TrowErrorForUnknownDerefName {
		_, ok := this.env.Variables[name]
		if !ok {
			panic("unknown variable name")
		}
	}
	return name
}

// DerefValue returns the value of a smart-service instance variable referenced in parameter ref
func (this *ScriptEnvIo) DerefValue(ref string) interface{} {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	return this.env.Variables[this.DerefName(ref)]
}

func (this *ScriptEnvIo) DerefTemplate(templ string) string {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	placeholder, err := getPlaceholder(this.env.Variables)
	if err != nil {
		panic(err)
	}
	result, err := replaceReferences(templ, placeholder)
	if err != nil {
		panic(err)
	}
	return result
}
