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
	"context"
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

	result, err := New(config, AuthMock).GetSmartServiceInstance(context.Background(), "my-process-instance-id")

	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Error(result)
	}
}

// TestGetCachedSmartServiceInstance checks that the cache that GetInstanceUser used to hold is
// effective on the call that replaced it in the middleware: one request per process-instance,
// not one per task, and the same result out of the cache as from the repository.
func TestGetCachedSmartServiceInstance(t *testing.T) {
	expectedResult := model.SmartServiceInstance{
		Id:        "test-id",
		UserId:    "user",
		DesignId:  "design",
		ReleaseId: "release",
		Ready:     true,
	}
	response, _ := json.Marshal(expectedResult)

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestCount++
		if request.URL.Path != "/instances-by-process-id/my-process-instance-id" {
			t.Error(request.URL.Path)
		}
		writer.Write(response)
	}))
	defer server.Close()

	repo := New(configuration.Config{SmartServiceRepositoryUrl: server.URL}, AuthMock)

	for i := 0; i < 3; i++ {
		result, err := repo.GetCachedSmartServiceInstance(context.Background(), "my-process-instance-id")
		if err != nil {
			t.Error(err)
			return
		}
		if !reflect.DeepEqual(result, expectedResult) {
			t.Errorf("run %v: %#v", i, result)
			return
		}
	}
	if requestCount != 1 {
		t.Error(requestCount)
	}

	//the uncached call is not affected by the cache
	if _, err := repo.GetSmartServiceInstance(context.Background(), "my-process-instance-id"); err != nil {
		t.Error(err)
		return
	}
	if requestCount != 2 {
		t.Error(requestCount)
	}
}

func TestGetCachedSmartServiceInstanceReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, "nope", http.StatusInternalServerError)
	}))
	defer server.Close()

	repo := New(configuration.Config{SmartServiceRepositoryUrl: server.URL}, AuthMock)
	_, err := repo.GetCachedSmartServiceInstance(context.Background(), "my-process-instance-id")
	if err == nil {
		t.Error("expected error")
	}
}

// TestGetInstanceUserUsesInstanceCache is the reason GetInstanceUser delegates to
// GetCachedSmartServiceInstance: the middleware reads the instance for every task, so the
// user-id must come out of that cache entry instead of causing a second request per task.
func TestGetInstanceUserUsesInstanceCache(t *testing.T) {
	expectedInstance := model.SmartServiceInstance{
		Id:        "test-id",
		UserId:    "user",
		DesignId:  "design",
		ReleaseId: "release",
		Ready:     true,
	}
	response, _ := json.Marshal(expectedInstance)

	requestCount := 0
	requestedPaths := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestCount++
		requestedPaths = append(requestedPaths, request.URL.Path)
		writer.Write(response)
	}))
	defer server.Close()

	repo := New(configuration.Config{SmartServiceRepositoryUrl: server.URL}, AuthMock)

	//what the middleware does at the start of a task
	instance, err := repo.GetCachedSmartServiceInstance(context.Background(), "my-process-instance-id")
	if err != nil {
		t.Error(err)
		return
	}
	if instance.UserId != expectedInstance.UserId {
		t.Error(instance.UserId)
		return
	}

	//what the worker handler does afterwards, for the same process instance
	userId, err := repo.GetInstanceUser(context.Background(), "my-process-instance-id")
	if err != nil {
		t.Error(err)
		return
	}
	if userId != expectedInstance.UserId {
		//no early return: the request count below says why
		t.Error(userId)
	}

	if requestCount != 1 {
		t.Error(requestCount, requestedPaths)
	}
	if !reflect.DeepEqual(requestedPaths, []string{"/instances-by-process-id/my-process-instance-id"}) {
		t.Error(requestedPaths)
	}
}

// TestGetInstanceUserReturnsError: GetInstanceUser used to drop the error of the cache getter
// and to answer with an empty user-id, which every caller then used as if it were a user.
func TestGetInstanceUserReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, "nope", http.StatusInternalServerError)
	}))
	defer server.Close()

	repo := New(configuration.Config{SmartServiceRepositoryUrl: server.URL}, AuthMock)
	userId, err := repo.GetInstanceUser(context.Background(), "my-process-instance-id")
	if err == nil {
		t.Error("expected error")
	}
	if userId != "" {
		t.Error(userId)
	}
}

// TestRunHealthCheckPassesContextToCheck: the check callback must run under the context of the
// health check run. before it had a ctx parameter the workers used their lifetime context, which
// carries neither span nor baggage, so the requests of a health check looked connected without
// being connected.
func TestRunHealthCheckPassesContextToCheck(t *testing.T) {
	firstPage, _ := json.Marshal([]model.SmartServiceModule{
		{SmartServiceModuleBase: model.SmartServiceModuleBase{Id: "module-id", UserId: "user"}},
	})
	emptyPage, _ := json.Marshal([]model.SmartServiceModule{})

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/modules" {
			t.Error(request.URL.Path)
		}
		if request.URL.Query().Get("offset") == "" {
			writer.Write(firstPage)
		} else {
			writer.Write(emptyPage)
		}
	}))
	defer server.Close()

	type ctxKeyType string
	const ctxKey = ctxKeyType("health-check-test")
	runCtx := context.WithValue(context.Background(), ctxKey, "caller-value")

	checkedModules := []string{}
	var checkCtx context.Context
	repo := New(configuration.Config{SmartServiceRepositoryUrl: server.URL}, AuthMock)
	repo.RunHealthCheck(runCtx, model.ModulQuery{}, func(ctx context.Context, module model.SmartServiceModule) (health error, err error) {
		checkCtx = ctx
		checkedModules = append(checkedModules, module.Id)
		return nil, nil
	})

	if !reflect.DeepEqual(checkedModules, []string{"module-id"}) {
		t.Error(checkedModules)
		return
	}
	if checkCtx == nil {
		t.Error("check callback got no context")
		return
	}
	//derived from the context of the run, not a fresh background context
	if value, _ := checkCtx.Value(ctxKey).(string); value != "caller-value" {
		t.Errorf("check context: %#v", checkCtx.Value(ctxKey))
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

	result, err := New(config, AuthMock).ListModules(context.Background(), model.ModulQuery{})

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

	result, err := New(config, AuthMock).ListExistingModules(context.Background(), "my-process-instance-id", model.ModulQuery{})

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
	result, err := New(config, AuthMock).ListExistingModules(context.Background(), "my-process-instance-id", model.ModulQuery{TypeFilter: &filter})

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
	result, err := New(config, AuthMock).ListExistingModules(context.Background(), "my-process-instance-id", model.ModulQuery{KeyFilter: &filter})

	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(result, expectedResult) {
		t.Error(result)
	}
}
