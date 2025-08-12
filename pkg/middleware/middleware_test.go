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
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/SENERGY-Platform/device-repository/lib/client"
	deviceRepoModel "github.com/SENERGY-Platform/device-repository/lib/model"
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/auth"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/configuration"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/model"
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

	testIotClient, testDb, err := client.NewTestClient()
	if err != nil {
		t.Error(err)
		return
	}
	sec := testDb

	middleware := New(configuration.Config{}, handler, repo, AuthMock, testIotClient)

	t.Run("check placeholder substitution", func(t *testing.T) {
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
	})

	t.Run("check script error handling", func(t *testing.T) {
		err = testDb.SetDevice(context.Background(), deviceRepoModel.DeviceWithConnectionState{
			Device: models.Device{
				Id:           "device1",
				LocalId:      "device1lid",
				Name:         "device1name",
				Attributes:   nil,
				DeviceTypeId: "dtid",
			},
		})
		if err != nil {
			t.Error(err)
			return
		}
		err = sec.SetRights("devices", "device1", deviceRepoModel.ResourceRights{
			UserRights: map[string]deviceRepoModel.Right{},
			GroupRights: map[string]deviceRepoModel.Right{
				"user":  {Read: true, Write: true, Execute: true, Administrate: true},
				"admin": {Read: true, Write: true, Execute: true, Administrate: true},
			},
		})
		if err != nil {
			t.Error(err)
			return
		}

		t.Run("ok, no error check", func(t *testing.T) {
			_, outputs, err := middleware.Do(model.CamundaExternalTask{
				ProcessInstanceId: "test-instance",
				Variables: map[string]model.CamundaVariable{
					PreScriptPrefix + "_1": {
						Value: `
					var result_as_Device = deviceRepo.readDevice("device1");
					outputs.set("device_name", result_as_Device.name);`,
					},
				},
			})

			if err != nil {
				t.Error(err)
				return
			}

			if outputs["device_name"] != "device1name" {
				t.Errorf("%#v\n", outputs)
				return
			}
		})

		t.Run("unknown device, no error check", func(t *testing.T) {
			_, _, err = middleware.Do(model.CamundaExternalTask{
				Variables: map[string]model.CamundaVariable{
					PreScriptPrefix + "_1": {
						Value: `var result_as_Device = deviceRepo.readDevice("unknown");`,
					},
				},
			})

			if err == nil {
				t.Error("expected error")
				return
			}
		})

		t.Run("null device-id, no error check", func(t *testing.T) {
			_, _, err = middleware.Do(model.CamundaExternalTask{
				Variables: map[string]model.CamundaVariable{
					PreScriptPrefix + "_1": {
						Value: `var result_as_Device = deviceRepo.readDevice(null);`,
					},
				},
			})

			if err == nil {
				t.Error("expected error")
				return
			}
		})

		t.Run("null ref, no error check", func(t *testing.T) {
			_, _, err = middleware.Do(model.CamundaExternalTask{
				Variables: map[string]model.CamundaVariable{
					PreScriptPrefix + "_1": {
						Value: `
							var result_as_Device = null;
							outputs.set("device_name", result_as_Device.name)`,
					},
				},
			})

			if err == nil {
				t.Error("expected error")
				return
			}
		})

		t.Run("unknown device, try-catch", func(t *testing.T) {
			_, outputs, err := middleware.Do(model.CamundaExternalTask{
				Variables: map[string]model.CamundaVariable{
					PreScriptPrefix + "_1": {
						Value: `
					try{
						var result_as_Device = deviceRepo.readDevice("unknown");
						outputs.set("device_name", result_as_Device.name);
					} catch (error) {
						outputs.set("error", error);
					}`,
					},
				},
			})

			if err != nil {
				t.Error(err)
				return
			}

			if _, ok := outputs["error"]; !ok {
				t.Errorf("%#v\n", outputs)
				return
			}
		})
	})
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
	testIotClient, _, err := client.NewTestClient()
	if err != nil {
		t.Error(err)
		return
	}

	middleware := New(configuration.Config{}, handler, repo, auth.New(configuration.Config{}), testIotClient)

	_, outputs, err := middleware.Do(model.CamundaExternalTask{
		Variables: map[string]model.CamundaVariable{
			"inp1": {Value: "42"},
			"inp2": {Value: 43},
			"prescript_1": {Value: `
					variables.write("toBeUpdatesInPre", variables.read("toBeUpdatesInPre") + 1);
					variables.write('added', 'foo'); //use single quote
					outputs.set("overwriting", 1);
					outputs.set("inp`},
			"prescript_2": {Value: `ut_v1", inputs.get("inp1"));
					variables.write("toBeUpdatesInPreAndPost", variables.read("toBeUpdatesInPreAndPost") + 1);
					variables.write("addedJson", JSON.stringify("foo"));
					variables.write("long_result", "a long text");
					outputs.set("long_result_output", variables.ref("long_result"));
			`},
			"postscript": {Value: `
					variables.write("toBeUpdatesInPreAndPost", variables.read("toBeUpdatesInPreAndPost") + 1);
					variables.write("toBeUpdatesInPost", variables.read("toBeUpdatesInPost") + 1);
					outputs.setJson("input_list", inputs.list());
					outputs.setJson("input_listNames", inputs.listNames());
					outputs.set("overwritten", 3);
					outputs.set("outputsGet", outputs.get("bar")+"_read");
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
		"outputsGet":         "batz_read",
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

func (this *VariablesRepoMock) GetInstanceUser(instanceId string) (userId string, err error) {
	return "user-id", nil
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

type AuthMockType string

const AuthMock AuthMockType = `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiIwOGM0N2E4OC0yYzc5LTQyMGYtODEwNC02NWJkOWViYmU0MWUiLCJleHAiOjE1NDY1MDcyMzMsIm5iZiI6MCwiaWF0IjoxNTQ2NTA3MTczLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjgwMDEvYXV0aC9yZWFsbXMvbWFzdGVyIiwiYXVkIjoiZnJvbnRlbmQiLCJzdWIiOiJ0ZXN0T3duZXIiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJmcm9udGVuZCIsIm5vbmNlIjoiOTJjNDNjOTUtNzViMC00NmNmLTgwYWUtNDVkZDk3M2I0YjdmIiwiYXV0aF90aW1lIjoxNTQ2NTA3MDA5LCJzZXNzaW9uX3N0YXRlIjoiNWRmOTI4ZjQtMDhmMC00ZWI5LTliNjAtM2EwYWUyMmVmYzczIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJ1c2VyIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsibWFzdGVyLXJlYWxtIjp7InJvbGVzIjpbInZpZXctcmVhbG0iLCJ2aWV3LWlkZW50aXR5LXByb3ZpZGVycyIsIm1hbmFnZS1pZGVudGl0eS1wcm92aWRlcnMiLCJpbXBlcnNvbmF0aW9uIiwiY3JlYXRlLWNsaWVudCIsIm1hbmFnZS11c2VycyIsInF1ZXJ5LXJlYWxtcyIsInZpZXctYXV0aG9yaXphdGlvbiIsInF1ZXJ5LWNsaWVudHMiLCJxdWVyeS11c2VycyIsIm1hbmFnZS1ldmVudHMiLCJtYW5hZ2UtcmVhbG0iLCJ2aWV3LWV2ZW50cyIsInZpZXctdXNlcnMiLCJ2aWV3LWNsaWVudHMiLCJtYW5hZ2UtYXV0aG9yaXphdGlvbiIsIm1hbmFnZS1jbGllbnRzIiwicXVlcnktZ3JvdXBzIl19LCJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJyb2xlcyI6WyJ1c2VyIl19.ykpuOmlpzj75ecSI6cHbCATIeY4qpyut2hMc1a67Ycg`

func (this AuthMockType) ExchangeUserToken(userid string) (token auth.Token, err error) {
	return auth.Parse(string(this))
}
