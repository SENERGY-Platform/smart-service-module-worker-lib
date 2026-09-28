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

package util

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/SENERGY-Platform/device-repository/v2/lib/client"
	devicemodel "github.com/SENERGY-Platform/device-repository/v2/lib/model"
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/model"
)

// A filter criteria may name several aspects in AspectIds; AspectId is the deprecated alias
// for a list with one element. This library only hands the criteria to device-repository,
// which resolves the alias at its own controller boundary, so what has to hold here is that
// both spellings arrive there exactly as the caller wrote them. Folding one into the other
// would change what a device-repository predating the lists sees.
func TestGetDevicesWithServiceForwardsAspects(t *testing.T) {
	const functionId = models.URN_PREFIX + "measuring-function:getTemperature"
	const airAspect = models.URN_PREFIX + "aspect:air"
	const insideAspect = models.URN_PREFIX + "aspect:inside"

	tests := []struct {
		name     string
		criteria []devicemodel.FilterCriteria
	}{
		{
			name:     "an aspect list",
			criteria: []devicemodel.FilterCriteria{{FunctionId: functionId, AspectIds: []string{airAspect, insideAspect}}},
		},
		{
			name:     "the deprecated single aspect id",
			criteria: []devicemodel.FilterCriteria{{FunctionId: functionId, AspectId: airAspect}},
		},
		{
			name:     "both spellings side by side",
			criteria: []devicemodel.FilterCriteria{{FunctionId: functionId, AspectId: airAspect, AspectIds: []string{insideAspect}}},
		},
		{
			name:     "no aspect at all",
			criteria: []devicemodel.FilterCriteria{{FunctionId: functionId}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			iotClient := &iotClientStub{devices: map[string]models.Device{
				"device1": {Id: "device1", DeviceTypeId: "dt1"},
			}}
			_, err := GetDevicesWithService(iotClient, "token", model.IotOption{
				DeviceSelection: &model.DeviceSelection{DeviceId: "device1"},
			}, test.criteria)
			if err != nil {
				t.Error(err)
				return
			}
			if !reflect.DeepEqual(iotClient.criteria, test.criteria) {
				t.Errorf("device-repository was asked for %#v, not %#v", iotClient.criteria, test.criteria)
			}
		})
	}
}

// A path option names every aspect its content variable matched. The list has to survive the
// translation into an iot-option, next to the deprecated single node the answer still carries.
func TestGetDevicesWithServiceReadsMultiAspectPathOptions(t *testing.T) {
	const airAspect = models.URN_PREFIX + "aspect:air"
	const insideAspect = models.URN_PREFIX + "aspect:inside"

	iotClient := &iotClientStub{
		devices: map[string]models.Device{"device1": {Id: "device1", DeviceTypeId: "dt1"}},
		selectables: []devicemodel.DeviceTypeSelectable{{
			DeviceTypeId: "dt1",
			Services:     []models.Service{{Id: "service1"}},
			ServicePathOptions: map[string][]models.ServicePathOption{
				"service1": {{
					ServiceId:        "service1",
					Path:             "value.temperature",
					CharacteristicId: "characteristic1",
					AspectNode:       models.AspectNode{Id: airAspect},
					AspectNodes:      []models.AspectNode{{Id: airAspect}, {Id: insideAspect}},
				}},
			},
		}},
	}

	options, err := GetDevicesWithService(iotClient, "token", model.IotOption{
		DeviceSelection: &model.DeviceSelection{DeviceId: "device1"},
	}, []devicemodel.FilterCriteria{{AspectIds: []string{airAspect, insideAspect}}})
	if err != nil {
		t.Error(err)
		return
	}

	serviceId := "service1"
	path := "value.temperature"
	characteristicId := "characteristic1"
	expected := []model.IotOption{{DeviceSelection: &model.DeviceSelection{
		DeviceId:         "device1",
		ServiceId:        &serviceId,
		Path:             &path,
		CharacteristicId: &characteristicId,
	}}}
	if !reflect.DeepEqual(options, expected) {
		t.Errorf("%#v", options)
	}
}

// iotClientStub answers the two device-repository calls GetDevicesWithService makes and
// records the criteria it was asked with. Every other method of the client interface is
// inherited unimplemented, so a call this test does not expect fails instead of passing
// silently.
type iotClientStub struct {
	client.Interface
	devices     map[string]models.Device
	selectables []devicemodel.DeviceTypeSelectable

	criteria []devicemodel.FilterCriteria
}

func (this *iotClientStub) GetDeviceTypeSelectablesV2(query []devicemodel.FilterCriteria, pathPrefix string, includeModified bool, servicesMustMatchAllCriteria bool) (result []devicemodel.DeviceTypeSelectable, err error, code int) {
	this.criteria = query
	return this.selectables, nil, http.StatusOK
}

func (this *iotClientStub) ReadDevice(id string, token string, action devicemodel.AuthAction) (result models.Device, err error, errCode int) {
	device, ok := this.devices[id]
	if !ok {
		return result, errors.New("unknown device " + id), http.StatusNotFound
	}
	return device, nil, http.StatusOK
}
