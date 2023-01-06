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
	"github.com/SENERGY-Platform/device-repository/lib/model"
	"github.com/SENERGY-Platform/models/go/models"
)

type ScriptEnvIot struct {
	env *ScriptEnv
}

func NewIotScriptEnv(env *ScriptEnv) *ScriptEnvIot {
	return &ScriptEnvIot{env: env}
}

func (this *ScriptEnvIot) ReadDevice(id string) models.Device {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.ReadDevice(id, this.env.getToken(), model.READ)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) ReadDeviceByLocalId(localId string) models.Device {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.ReadDeviceByLocalId(localId, this.env.getToken(), model.READ)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) ReadHub(id string) models.Hub {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.ReadHub(id, this.env.getToken(), model.READ)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) ListHubDeviceIds(id string, asLocalId bool) []string {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.ListHubDeviceIds(id, this.env.getToken(), model.READ, asLocalId)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) ReadDeviceType(id string) models.DeviceType {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.ReadDeviceType(id, this.env.getToken())
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) ListDeviceTypes(limit int64, offset int64, sort string, filter []model.FilterCriteria, includeModified bool, includeUnmodified bool) []models.DeviceType {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.ListDeviceTypesV2(this.env.getToken(), limit, offset, sort, filter, includeModified, includeUnmodified)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetDeviceTypeSelectables(query []model.FilterCriteria, pathPrefix string, includeModified bool) []model.DeviceTypeSelectable {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetDeviceTypeSelectablesV2(query, pathPrefix, includeModified)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) ReadDeviceGroup(id string) models.DeviceGroup {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.ReadDeviceGroup(id, this.env.getToken())
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetService(id string) models.Service {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetService(id)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetAspects() []models.Aspect {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetAspects()
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetAspectsWithMeasuringFunction(ancestors bool, descendants bool) []models.Aspect {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetAspectsWithMeasuringFunction(ancestors, descendants)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetAspect(id string) models.Aspect {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetAspect(id)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetAspectNode(id string) models.AspectNode {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetAspectNode(id)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetAspectNodes() []models.AspectNode {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetAspectNodes()
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetAspectNodesMeasuringFunctions(id string, ancestors bool, descendants bool) []models.Function {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetAspectNodesMeasuringFunctions(id, ancestors, descendants)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetAspectNodesWithMeasuringFunction(ancestors bool, descendants bool) []models.AspectNode {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetAspectNodesWithMeasuringFunction(ancestors, descendants)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetAspectNodesByIdList(ids []string) []models.AspectNode {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetAspectNodesByIdList(ids)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetLeafCharacteristics() []models.Characteristic {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetLeafCharacteristics()
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetCharacteristic(id string) models.Characteristic {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetCharacteristic(id)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetConceptWithCharacteristics(id string) models.ConceptWithCharacteristics {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetConceptWithCharacteristics(id)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetConceptWithoutCharacteristics(id string) models.Concept {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetConceptWithoutCharacteristics(id)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetDeviceClasses() []models.DeviceClass {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetDeviceClasses()
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetDeviceClassesWithControllingFunctions() []models.DeviceClass {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetDeviceClassesWithControllingFunctions()
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetDeviceClassesFunctions(id string) []models.Function {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetDeviceClassesFunctions(id)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetDeviceClassesControllingFunctions(id string) []models.Function {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetDeviceClassesControllingFunctions(id)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetDeviceClass(id string) models.DeviceClass {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetDeviceClass(id)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetFunctionsByType(rdfType string) []models.Function {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetFunctionsByType(rdfType)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetFunction(id string) models.Function {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetFunction(id)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *ScriptEnvIot) GetLocation(id string) models.Location {
	defer func() {
		if caught := recover(); caught != nil {
			panic(this.env.GetVm().ToValue(caught))
		}
	}()
	result, err, _ := this.env.iotClient.GetLocation(id, this.env.getToken())
	if err != nil {
		panic(err)
	}
	return result
}
