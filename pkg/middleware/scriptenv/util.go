/*
 * Copyright (c) 2023 InfAI (CC SES)
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

package scriptenv

import (
	"encoding/json"
	"errors"
	devicemodel "github.com/SENERGY-Platform/device-repository/lib/model"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/model"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/util"
)

type ScriptEnvUtil struct {
	env *ScriptEnv
}

func NewScriptEnvUtil(env *ScriptEnv) *ScriptEnvUtil {
	return &ScriptEnvUtil{env: env}
}

// GetDevicesWithServiceFromIotOption finds a list of iot-options where the entity is the same the input, but the Service field is set with those that match the input criteria
func (this *ScriptEnvUtil) GetDevicesWithServiceFromIotOption(entity model.IotOption, criteria []devicemodel.FilterCriteria) []model.IotOption {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	return this.getDevicesWithServiceFromIotOption(entity, criteria)
}

// GetDevicesWithServiceFromEntityString finds a list of iot-options where the entity is the same the input, but the Service field is set with those that match the input criteria
func (this *ScriptEnvUtil) GetDevicesWithServiceFromEntityString(entityStr string, criteria []devicemodel.FilterCriteria) []model.IotOption {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	entity := model.IotOption{}
	err := json.Unmarshal([]byte(entityStr), &entity)
	if err != nil {
		panic(err)
	}
	return this.getDevicesWithServiceFromIotOption(entity, criteria)
}

func (this *ScriptEnvUtil) getDevicesWithServiceFromIotOption(entity model.IotOption, criteria []devicemodel.FilterCriteria) (result []model.IotOption) {
	var err error
	result, err = util.GetDevicesWithService(this.env.iotClient, this.env.getToken(), entity, criteria)
	if errors.Is(err, util.ErrNoDeviceOrGroupSelection) {
		return []model.IotOption{}
	}
	if err != nil {
		panic(err)
	}
	return result
}

// IsDeviceIotOption checks if the input is a device
func (this *ScriptEnvUtil) IsDeviceIotOption(entity model.IotOption) bool {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	return entity.DeviceSelection != nil
}

// IsDeviceIotOptionStr checks if the input is a device
func (this *ScriptEnvUtil) IsDeviceIotOptionStr(entityStr string) bool {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	entity := model.IotOption{}
	err := json.Unmarshal([]byte(entityStr), &entity)
	if err != nil {
		panic(err)
	}
	return entity.DeviceSelection != nil
}

// IsDeviceGroupIotOption checks if the input is a device-group
func (this *ScriptEnvUtil) IsDeviceGroupIotOption(entity model.IotOption) bool {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	return entity.DeviceGroupSelection != nil
}

// IsDeviceGroupIotOptionStr checks if the input is a device-group
func (this *ScriptEnvUtil) IsDeviceGroupIotOptionStr(entityStr string) bool {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	entity := model.IotOption{}
	err := json.Unmarshal([]byte(entityStr), &entity)
	if err != nil {
		panic(err)
	}
	return entity.DeviceGroupSelection != nil
}

// IsImportIotOption checks if the input is a import
func (this *ScriptEnvUtil) IsImportIotOption(entity model.IotOption) bool {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	return entity.ImportSelection != nil
}

// IsImportIotOptionStr checks if the input is a import
func (this *ScriptEnvUtil) IsImportIotOptionStr(entityStr string) bool {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	entity := model.IotOption{}
	err := json.Unmarshal([]byte(entityStr), &entity)
	if err != nil {
		panic(err)
	}
	return entity.ImportSelection != nil
}

// GroupIotOptionsByService groups a list of IotOption by their service id; options that are not devices or dont hav a service-id will be grouped under ""
func (this *ScriptEnvUtil) GroupIotOptionsByService(entities []model.IotOption) map[string][]model.IotOption {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	return util.GroupIotOptionsByService(entities)
}

// GroupIotOptionsByDevice groups a list of model.IotOption by their device id; options that are not devices will be grouped under ""
func (this *ScriptEnvUtil) GroupIotOptionsByDevice(entities []model.IotOption) map[string][]model.IotOption {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	return util.GroupIotOptionsByDevice(entities)
}
