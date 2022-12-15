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
	"time"
)

type ScriptContext interface {
	RegisterRuntime(runtime *goja.Runtime)
	GetEnvironment() map[string]interface{}
}

func runScript(script string, ctx ScriptContext) (err error) {
	vm := goja.New()
	time.AfterFunc(2*time.Second, func() {
		vm.Interrupt("script execution timeout")
	})
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	for key, value := range ctx.GetEnvironment() {
		err = vm.Set(key, value)
		if err != nil {
			return err
		}
	}

	_, err = vm.RunString(script)
	return err
}
