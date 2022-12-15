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
	"errors"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/model"
	"reflect"
	"testing"
)

func TestMiddleware(t *testing.T) {
	handler := &HandlerMock{DoFunc: func(task model.CamundaExternalTask) (modules []model.Module, outputs map[string]interface{}, err error) {
		return []model.Module{{
				SmartServiceModuleInit: model.SmartServiceModuleInit{
					ModuleData: map[string]interface{}{"result": task.Variables},
				},
			}}, map[string]interface{}{
				"result": task.Variables,
			}, nil
	}}
	repo := &VariablesRepoMock{GetVariablesFunc: func(processId string) (result map[string]interface{}, err error) {
		return map[string]interface{}{
			"v1": "str",
			"v2": float64(42),
			"v3": true,
			"v4": nil,
			"v5": map[string]interface{}{
				"foo":  "bar",
				"batz": float64(42),
			},
		}, nil
	}}
	middleware := New(handler, repo)

	_, outputs, err := middleware.Do(model.CamundaExternalTask{
		Variables: map[string]model.CamundaVariable{
			"templ": {Value: "{{.brl}}placeholder{{.brr}}"},
			"str":   {Value: "{{.v1}}"},
			//"strJson":   {Value: "{{.v1_json}}"},
			"number":    {Value: "{{.v2}}"},
			"bool":      {Value: "{{.v3}}"},
			"null":      {Value: "{{.v4}}"},
			"obj":       {Value: "{{.v5}}"},
			"unknown":   {Value: "{{.unknown}}"},
			"rawStr":    {Value: "raw"},
			"rawNumber": {Value: 13},
			"rawBool":   {Value: true},
			"rawNull":   {Value: nil},
		},
	})

	if err != nil {
		t.Error(err)
		return
	}
	expected := map[string]interface{}{
		"result": map[string]model.CamundaVariable{
			"templ": {Value: "{{placeholder}}"},
			"str":   {Value: "str"},
			//"strJson":   {Value: `"str"`},
			"number":    {Value: "42"},
			"bool":      {Value: "true"},
			"null":      {Value: "null"},
			"obj":       {Value: `{"batz":42,"foo":"bar"}`},
			"unknown":   {Value: ""},
			"rawStr":    {Value: "raw"},
			"rawNumber": {Value: 13},
			"rawBool":   {Value: true},
			"rawNull":   {Value: nil},
		},
	}

	if !reflect.DeepEqual(outputs, expected) {
		t.Errorf("\n%#v\n%#v", outputs, expected)
	}
}

func TestMiddlewareScripts(t *testing.T) {
	handler := &HandlerMock{DoFunc: func(task model.CamundaExternalTask) (modules []model.Module, outputs map[string]interface{}, err error) {
		return []model.Module{}, map[string]interface{}{
			"bar":         "batz",
			"overwriting": 2,
			"overwritten": 2,
		}, nil
	}}
	repo := &VariablesRepoMock{
		GetVariablesFunc: func(processId string) (result map[string]interface{}, err error) {
			return map[string]interface{}{
				"v1": "str",
				"v2": float64(42),
				"v3": true,
				"v4": nil,
				"v5": map[string]interface{}{
					"foo":  "bar",
					"batz": float64(42),
				},
				"toBeUpdatesInPre":        1,
				"toBeUpdatesInPost":       1,
				"toBeUpdatesInPreAndPost": 1,
			}, nil
		},
		SetVariablesFunc: func(processId string, changes map[string]interface{}) (err error) {
			if !reflect.DeepEqual(changes, map[string]interface{}{
				"toBeUpdatesInPre":        int64(2),
				"toBeUpdatesInPost":       int64(2),
				"toBeUpdatesInPreAndPost": int64(3),
				"added":                   "foo",
				"addedJson":               `"foo"`,
				"long_result":             "a long text",
			}) {
				t.Errorf("%#v", changes)
			}
			return nil
		},
	}
	middleware := New(handler, repo)

	_, outputs, err := middleware.Do(model.CamundaExternalTask{
		Variables: map[string]model.CamundaVariable{
			"inp1": {Value: "42"},
			"inp2": {Value: 43},
			"prescript_1": {Value: `
					io.store("toBeUpdatesInPre", io.read("toBeUpdatesInPre") + 1);
					io.store('added', 'foo'); //use single quote
					outputs.set("overwriting", 1);
					outputs.set("input_v1", inputs.get("inp1"));
			`},
			"prescript_2": {Value: `
					io.store("toBeUpdatesInPreAndPost", io.read("toBeUpdatesInPreAndPost") + 1);
					io.store("addedJson", JSON.stringify("foo"));
					io.store("long_result", "a long text");
					outputs.set("long_result_output", io.ref("long_result"));
			`},
			"postscript": {Value: `
					io.store("toBeUpdatesInPreAndPost", io.read("toBeUpdatesInPreAndPost") + 1);
					io.store("toBeUpdatesInPost", io.read("toBeUpdatesInPost") + 1);
					outputs.setJson("input_list", inputs.list());
					outputs.setJson("input_listNames", inputs.listNames());
					outputs.set("overwritten", 3);
			`},
		},
	})

	if err != nil {
		t.Error(err)
		return
	}
	expectedOutput := map[string]interface{}{
		"bar":                "batz",
		"overwriting":        2,
		"overwritten":        int64(3),
		"input_v1":           "42",
		"input_list":         `["42",43]`,
		"input_listNames":    `["inp1","inp2"]`,
		"long_result_output": "{{.long_result}}",
	}

	if !reflect.DeepEqual(outputs, expectedOutput) {
		t.Errorf("\n%#v\n%#v", outputs, expectedOutput)
	}
}

type HandlerMock struct {
	DoFunc func(task model.CamundaExternalTask) (modules []model.Module, outputs map[string]interface{}, err error)
}

func (this *HandlerMock) Do(task model.CamundaExternalTask) (modules []model.Module, outputs map[string]interface{}, err error) {
	if this.DoFunc == nil {
		return modules, outputs, errors.New("missing mock DoFunc")
	}
	return this.DoFunc(task)
}

func (this *HandlerMock) Undo(modules []model.Module, reason error) {}

type VariablesRepoMock struct {
	GetVariablesFunc func(processId string) (result map[string]interface{}, err error)
	SetVariablesFunc func(processId string, changes map[string]interface{}) (err error)
}

func (this *VariablesRepoMock) SetVariables(processId string, changes map[string]interface{}) error {
	if this.SetVariablesFunc != nil {
		return this.SetVariablesFunc(processId, changes)
	}
	return nil
}

func (this *VariablesRepoMock) GetVariables(processId string) (result map[string]interface{}, err error) {
	if this.GetVariablesFunc == nil {
		return result, errors.New("missing mock GetVariablesFunc")
	}
	return this.GetVariablesFunc(processId)
}
