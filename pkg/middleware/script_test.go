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
	"reflect"
	"testing"
)

type ScriptContextMock struct {
	vm        *goja.Runtime
	StoreFunc func(name string, value interface{})
}

func (this *ScriptContextMock) GetEnvironment() map[string]interface{} {
	return map[string]interface{}{"io": this}
}

func (this *ScriptContextMock) RegisterRuntime(runtime *goja.Runtime) {
	this.vm = runtime
}

func (this *ScriptContextMock) Store(name string, value interface{}) {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.vm.ToValue(caught))
		}
	}()
	this.StoreFunc(name, value)
}

func TestRunScript(t *testing.T) {
	t.Run("io.store", func(t *testing.T) {
		err := runScript(`io.store("foo", 42)`, &ScriptContextMock{StoreFunc: func(name string, value interface{}) {
			if name != "foo" {
				t.Errorf("name: %#v, value: %#v", name, value)
			}
			if !reflect.DeepEqual(value, int64(42)) {
				t.Errorf("name: %#v, value: %#v", name, value)
			}
		}})
		if err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("io.unknown", func(t *testing.T) {
		err := runScript(`io.foo("foo", 42)`, &ScriptContextMock{StoreFunc: func(name string, value interface{}) {}})
		if err == nil {
			t.Error("expected error like 'TypeError: Object has no member 'foo' at <eval>:1:7(4)'")
			return
		}
	})

	t.Run("go code error", func(t *testing.T) {
		err := runScript(`io.store("foo", 42)`, &ScriptContextMock{StoreFunc: func(name string, value interface{}) {
			panic("my error")
		}})
		if err == nil {
			t.Error("expected error like 'my error at reflect.methodValueCall (native)'")
			return
		}
	})

	t.Run("cached go code error", func(t *testing.T) {
		err := runScript(`try {
    io.store("foo", 42)
} catch(e) {
}`, &ScriptContextMock{StoreFunc: func(name string, value interface{}) {
			panic("my error")
		}})
		if err != nil {
			t.Error(err)
			return
		}
	})
}