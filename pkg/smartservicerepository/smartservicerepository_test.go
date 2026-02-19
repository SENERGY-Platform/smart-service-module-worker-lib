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

package smartservicerepository

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/auth"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/configuration"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/model"
)

type AuthMockType string

func (this AuthMockType) Ensure() (token auth.Token, err error) {
	return auth.Token{
		Token:       string(this),
		Sub:         "",
		RealmAccess: nil,
	}, nil
}

func (this AuthMockType) ExchangeUserToken(userid string) (token auth.Token, err error) {
	return auth.Token{
		Token:       string(this),
		Sub:         "",
		RealmAccess: nil,
	}, nil
}

var AuthMock AuthMockType = "token"

func TestGetSmartServiceInstance(t *testing.T) {
	expectedMethod := "GET"
	expectedEndpoint := "/instances-by-process-id/my-process-instance-id"
	expectedResult := model.SmartServiceInstance{
		Id:        "test-id",
		UserId:    "user",
		DesignId:  "design",
		ReleaseId: "release",
		Ready:     false,
	}
	response, _ := json.Marshal(expectedResult)

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != expectedMethod {
			t.Error(request.Method)
		}
		path := request.URL.Path
		if request.URL.RawQuery != "" {
			path = path + "?" + request.URL.RawQuery
		}
		if path != expectedEndpoint {
			t.Error(path)
		}
		writer.Write(response)
	}))

	defer server.Close()

	config := configuration.Config{SmartServiceRepositoryUrl: server.URL}

	result, err := New(config, AuthMock).GetSmartServiceInstance("my-process-instance-id")

	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Error(result)
	}
}

func TestListModules(t *testing.T) {
	expectedMethod := "GET"
	expectedEndpoint := "/modules"
	expectedResult := []model.SmartServiceModule{
		{
			SmartServiceModuleBase: model.SmartServiceModuleBase{
				Id:         "module-id",
				UserId:     "user",
				InstanceId: "instance",
				DesignId:   "design",
				ReleaseId:  "release",
			},
			SmartServiceModuleInit: model.SmartServiceModuleInit{
				ModuleType: "module-type",
				Keys:       []string{"key1"},
				ModuleData: map[string]interface{}{
					"foo": "bar",
				},
			},
		},
	}
	response, _ := json.Marshal(expectedResult)

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != expectedMethod {
			t.Error(request.Method)
		}
		path := request.URL.Path
		if request.URL.RawQuery != "" {
			path = path + "?" + request.URL.RawQuery
		}
		if path != expectedEndpoint {
			t.Error(path)
		}
		writer.Write(response)
	}))

	defer server.Close()

	config := configuration.Config{SmartServiceRepositoryUrl: server.URL}

	result, err := New(config, AuthMock).ListModules(model.ModulQuery{})

	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Error(result)
	}
}

func TestListExistingModules(t *testing.T) {
	expectedMethod := "GET"
	expectedEndpoint := "/instances-by-process-id/my-process-instance-id/modules"
	expectedResult := []model.SmartServiceModule{
		{
			SmartServiceModuleBase: model.SmartServiceModuleBase{
				Id:         "module-id",
				UserId:     "user",
				InstanceId: "instance",
				DesignId:   "design",
				ReleaseId:  "release",
			},
			SmartServiceModuleInit: model.SmartServiceModuleInit{
				ModuleType: "module-type",
				Keys:       []string{"key1"},
				ModuleData: map[string]interface{}{
					"foo": "bar",
				},
			},
		},
	}
	response, _ := json.Marshal(expectedResult)

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != expectedMethod {
			t.Error(request.Method)
		}
		path := request.URL.Path
		if request.URL.RawQuery != "" {
			path = path + "?" + request.URL.RawQuery
		}
		if path != expectedEndpoint {
			t.Error(path)
		}
		writer.Write(response)
	}))

	defer server.Close()

	config := configuration.Config{SmartServiceRepositoryUrl: server.URL}

	result, err := New(config, AuthMock).ListExistingModules("my-process-instance-id", model.ModulQuery{})

	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Error(result)
	}
}

func TestListExistingModulesFilterByType(t *testing.T) {
	expectedMethod := "GET"
	expectedEndpoint := "/instances-by-process-id/my-process-instance-id/modules?module_type=filter"
	expectedResult := []model.SmartServiceModule{
		{
			SmartServiceModuleBase: model.SmartServiceModuleBase{
				Id:         "module-id",
				UserId:     "user",
				InstanceId: "instance",
				DesignId:   "design",
				ReleaseId:  "release",
			},
			SmartServiceModuleInit: model.SmartServiceModuleInit{
				ModuleType: "mt",
				Keys:       []string{"key1"},
				ModuleData: map[string]interface{}{
					"foo": "bar",
				},
			},
		},
	}
	response, _ := json.Marshal(expectedResult)

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != expectedMethod {
			t.Error(request.Method)
		}
		path := request.URL.Path
		if request.URL.RawQuery != "" {
			path = path + "?" + request.URL.RawQuery
		}
		if path != expectedEndpoint {
			t.Error(path)
		}
		writer.Write(response)
	}))

	defer server.Close()

	config := configuration.Config{SmartServiceRepositoryUrl: server.URL}

	filter := "filter"
	result, err := New(config, AuthMock).ListExistingModules("my-process-instance-id", model.ModulQuery{TypeFilter: &filter})

	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Error(result)
	}
}

func TestListExistingModulesFilterByKey(t *testing.T) {
	expectedMethod := "GET"
	expectedEndpoint := "/instances-by-process-id/my-process-instance-id/modules?key=filter"
	expectedResult := []model.SmartServiceModule{
		{
			SmartServiceModuleBase: model.SmartServiceModuleBase{
				Id:         "module-id",
				UserId:     "user",
				InstanceId: "instance",
				DesignId:   "design",
				ReleaseId:  "release",
			},
			SmartServiceModuleInit: model.SmartServiceModuleInit{
				ModuleType: "mt",
				Keys:       []string{"key1"},
				ModuleData: map[string]interface{}{
					"foo": "bar",
				},
			},
		},
	}
	response, _ := json.Marshal(expectedResult)

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != expectedMethod {
			t.Error(request.Method)
		}
		path := request.URL.Path
		if request.URL.RawQuery != "" {
			path = path + "?" + request.URL.RawQuery
		}
		if path != expectedEndpoint {
			t.Error(path)
		}
		writer.Write(response)
	}))

	defer server.Close()

	config := configuration.Config{SmartServiceRepositoryUrl: server.URL}

	filter := "filter"
	result, err := New(config, AuthMock).ListExistingModules("my-process-instance-id", model.ModulQuery{KeyFilter: &filter})

	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Error(result)
	}
}
