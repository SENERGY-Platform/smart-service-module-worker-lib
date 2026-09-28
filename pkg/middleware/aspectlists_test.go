/*
 * Copyright (c) 2026 InfAI (CC SES)
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
	"net/http"
	"reflect"
	"testing"

	"github.com/SENERGY-Platform/device-repository/v2/lib/client"
	devicemodel "github.com/SENERGY-Platform/device-repository/v2/lib/model"
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/configuration"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/model"
)

const (
	testFunctionId = models.URN_PREFIX + "measuring-function:getTemperature"
	testAirAspect  = models.URN_PREFIX + "aspect:air"
	testRoomAspect = models.URN_PREFIX + "aspect:room"
)

// A script names the aspects of a filter criteria in aspect_ids. The runtime shows a struct
// field under its json tag name, so a field the platform model added is only reachable for
// scripts once it is spelled that way - a criteria whose aspects silently never arrive at
// device-repository selects everything instead of failing.
func TestScriptWritesAspectIds(t *testing.T) {
	tests := []struct {
		name     string
		criteria string
		expected []devicemodel.FilterCriteria
	}{
		{
			name:     "an aspect list",
			criteria: `{function_id: fid, aspect_ids: [air, room]}`,
			expected: []devicemodel.FilterCriteria{{FunctionId: testFunctionId, AspectIds: []string{testAirAspect, testRoomAspect}}},
		},
		{
			name:     "an aspect list with one element",
			criteria: `{function_id: fid, aspect_ids: [air]}`,
			expected: []devicemodel.FilterCriteria{{FunctionId: testFunctionId, AspectIds: []string{testAirAspect}}},
		},
		{
			//a design written before the lists keeps working, unfolded: device-repository
			//resolves the alias, this library must not decide for it
			name:     "the deprecated single aspect id",
			criteria: `{function_id: fid, aspect_id: air}`,
			expected: []devicemodel.FilterCriteria{{FunctionId: testFunctionId, AspectId: testAirAspect}},
		},
		{
			name:     "both spellings side by side",
			criteria: `{function_id: fid, aspect_id: air, aspect_ids: [room]}`,
			expected: []devicemodel.FilterCriteria{{FunctionId: testFunctionId, AspectId: testAirAspect, AspectIds: []string{testRoomAspect}}},
		},
		{
			name:     "no aspect at all",
			criteria: `{function_id: fid}`,
			expected: []devicemodel.FilterCriteria{{FunctionId: testFunctionId}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			iotClient := &selectablesStub{}
			_, _, err := runPrescript(t, iotClient, aspectIdConstants+`
				deviceRepo.getDeviceTypeSelectables(`+"["+test.criteria+"]"+`, "", true, true);`)
			if err != nil {
				t.Error(err)
				return
			}
			if !reflect.DeepEqual(iotClient.criteria, test.expected) {
				t.Errorf("device-repository was asked for %#v, not %#v", iotClient.criteria, test.expected)
			}
		})
	}
}

// The answering side of the same move: a path option names every aspect it matched in
// aspect_nodes and keeps the first of them in the deprecated aspect_node. A script has to
// be able to read both.
func TestScriptReadsAspectNodes(t *testing.T) {
	iotClient := &selectablesStub{selectables: []devicemodel.DeviceTypeSelectable{{
		DeviceTypeId: "dt1",
		Services:     []models.Service{{Id: "service1"}},
		ServicePathOptions: map[string][]models.ServicePathOption{
			"service1": {{
				ServiceId:   "service1",
				Path:        "value.temperature",
				AspectNode:  models.AspectNode{Id: testAirAspect},
				AspectNodes: []models.AspectNode{{Id: testAirAspect}, {Id: testRoomAspect}},
			}},
		},
	}}}

	_, outputs, err := runPrescript(t, iotClient, aspectIdConstants+`
		var selectables = deviceRepo.getDeviceTypeSelectables([{function_id: fid, aspect_ids: [air, room]}], "", true, true);
		var option = selectables[0].service_path_options["service1"][0];
		outputs.set("aspect_node_count", option.aspect_nodes.length);
		outputs.set("first_aspect_node", option.aspect_nodes[0].id);
		outputs.set("second_aspect_node", option.aspect_nodes[1].id);
		outputs.set("deprecated_aspect_node", option.aspect_node.id);`)
	if err != nil {
		t.Error(err)
		return
	}

	expected := map[string]interface{}{
		"aspect_node_count":      int64(2),
		"first_aspect_node":      testAirAspect,
		"second_aspect_node":     testRoomAspect,
		"deprecated_aspect_node": testAirAspect,
	}
	if !reflect.DeepEqual(outputs, expected) {
		t.Errorf("%#v", outputs)
	}
}

// aspectIdConstants gives the scripts of these tests the ids as variables, so that the
// criteria under test stays readable next to the urns.
const aspectIdConstants = `
	var fid = "` + testFunctionId + `";
	var air = "` + testAirAspect + `";
	var room = "` + testRoomAspect + `";`

func runPrescript(t *testing.T, iotClient client.Interface, script string) (modules []model.Module, outputs map[string]interface{}, err error) {
	t.Helper()
	handler := &HandlerMock{DoFunc: func(ctx context.Context, task model.CamundaExternalTask) ([]model.Module, map[string]interface{}, error) {
		return nil, map[string]interface{}{}, nil
	}}
	repo := &VariablesRepoMock{GetVariablesFunc: func(processId string) (map[string]interface{}, error) {
		return map[string]interface{}{}, nil
	}}
	middleware := New(configuration.Config{}, handler, repo, AuthMock, iotClient)
	return middleware.Do(context.Background(), model.CamundaExternalTask{
		Variables: map[string]model.CamundaVariable{PreScriptPrefix: {Value: script}},
	})
}

// selectablesStub records the criteria the script env asked device-repository for. Every
// other method of the client interface is inherited unimplemented, so an unexpected call
// fails instead of passing silently.
type selectablesStub struct {
	client.Interface
	selectables []devicemodel.DeviceTypeSelectable

	criteria []devicemodel.FilterCriteria
}

func (this *selectablesStub) GetDeviceTypeSelectablesV2(query []devicemodel.FilterCriteria, pathPrefix string, includeModified bool, servicesMustMatchAllCriteria bool) (result []devicemodel.DeviceTypeSelectable, err error, code int) {
	this.criteria = query
	return this.selectables, nil, http.StatusOK
}
