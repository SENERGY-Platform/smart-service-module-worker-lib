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

/*
 * generated in github.com/SENERGY-Platform/smart-service-module-worker-lib with a command like:
 * go generate ./...
*/

 
 /**
  * @namespace deviceRepo
  */
 
 /**
  * @namespace inputs
  */
 
 /**
  * @namespace outputs
  */
 
 /**
  * @namespace util
  */
 
 /**
  * @namespace variables
  */
 


/** 
 * @function deviceRepo#getAspect
 * @param { string } id
 * @returns { Aspect }
 * @example
 * deviceRepo.getAspect(id_as_string)
 */

/** 
 * @function deviceRepo#getAspectNode
 * @param { string } id
 * @returns { AspectNode }
 * @example
 * deviceRepo.getAspectNode(id_as_string)
 */

/** 
 * @function deviceRepo#getAspectNodes
 * @returns { AspectNode[] }
 * @example
 * deviceRepo.getAspectNodes()
 */

/** 
 * @function deviceRepo#getAspectNodesByIdList
 * @param { string[] } ids
 * @returns { AspectNode[] }
 * @example
 * deviceRepo.getAspectNodesByIdList(ids_as_string_list)
 */

/** 
 * @function deviceRepo#getAspectNodesMeasuringFunctions
 * @param { string } id
 * @param { boolean } ancestors
 * @param { boolean } descendants
 * @returns { FunctionType[] }
 * @example
 * deviceRepo.getAspectNodesMeasuringFunctions(id_as_string, ancestors_as_boolean, descendants_as_boolean)
 */

/** 
 * @function deviceRepo#getAspectNodesWithMeasuringFunction
 * @param { boolean } ancestors
 * @param { boolean } descendants
 * @returns { AspectNode[] }
 * @example
 * deviceRepo.getAspectNodesWithMeasuringFunction(ancestors_as_boolean, descendants_as_boolean)
 */

/** 
 * @function deviceRepo#getAspects
 * @returns { Aspect[] }
 * @example
 * deviceRepo.getAspects()
 */

/** 
 * @function deviceRepo#getAspectsWithMeasuringFunction
 * @param { boolean } ancestors
 * @param { boolean } descendants
 * @returns { Aspect[] }
 * @example
 * deviceRepo.getAspectsWithMeasuringFunction(ancestors_as_boolean, descendants_as_boolean)
 */

/** 
 * @function deviceRepo#getCharacteristic
 * @param { string } id
 * @returns { Characteristic }
 * @example
 * deviceRepo.getCharacteristic(id_as_string)
 */

/** 
 * @function deviceRepo#getConceptWithCharacteristics
 * @param { string } id
 * @returns { ConceptWithCharacteristics }
 * @example
 * deviceRepo.getConceptWithCharacteristics(id_as_string)
 */

/** 
 * @function deviceRepo#getConceptWithoutCharacteristics
 * @param { string } id
 * @returns { Concept }
 * @example
 * deviceRepo.getConceptWithoutCharacteristics(id_as_string)
 */

/** 
 * @function deviceRepo#getDeviceClass
 * @param { string } id
 * @returns { DeviceClass }
 * @example
 * deviceRepo.getDeviceClass(id_as_string)
 */

/** 
 * @function deviceRepo#getDeviceClasses
 * @returns { DeviceClass[] }
 * @example
 * deviceRepo.getDeviceClasses()
 */

/** 
 * @function deviceRepo#getDeviceClassesControllingFunctions
 * @param { string } id
 * @returns { FunctionType[] }
 * @example
 * deviceRepo.getDeviceClassesControllingFunctions(id_as_string)
 */

/** 
 * @function deviceRepo#getDeviceClassesFunctions
 * @param { string } id
 * @returns { FunctionType[] }
 * @example
 * deviceRepo.getDeviceClassesFunctions(id_as_string)
 */

/** 
 * @function deviceRepo#getDeviceClassesWithControllingFunctions
 * @returns { DeviceClass[] }
 * @example
 * deviceRepo.getDeviceClassesWithControllingFunctions()
 */

/** 
 * @function deviceRepo#getDeviceTypeSelectables
 * @param { FilterCriteria[] } query
 * @param { string } pathPrefix
 * @param { boolean } includeModified
 * @param { boolean } servicesMustMatchAllCriteria
 * @returns { DeviceTypeSelectable[] }
 * @example
 * deviceRepo.getDeviceTypeSelectables(query_as_FilterCriteria_list, pathPrefix_as_string, includeModified_as_boolean, servicesMustMatchAllCriteria_as_boolean)
 */

/** 
 * @function deviceRepo#getFunction
 * @param { string } id
 * @returns { FunctionType }
 * @example
 * deviceRepo.getFunction(id_as_string)
 */

/** 
 * @function deviceRepo#getFunctionsByType
 * @param { string } rdfType
 * @returns { FunctionType[] }
 * @example
 * deviceRepo.getFunctionsByType(rdfType_as_string)
 */

/** 
 * @function deviceRepo#getLeafCharacteristics
 * @returns { Characteristic[] }
 * @example
 * deviceRepo.getLeafCharacteristics()
 */

/** 
 * @function deviceRepo#getLocation
 * @param { string } id
 * @returns { Location }
 * @example
 * deviceRepo.getLocation(id_as_string)
 */

/** 
 * @function deviceRepo#getService
 * @param { string } id
 * @returns { Service }
 * @example
 * deviceRepo.getService(id_as_string)
 */

/** 
 * @function deviceRepo#listDeviceTypes
 * @param { number } limit
 * @param { number } offset
 * @param { string } sort
 * @param { FilterCriteria[] } filter
 * @param { boolean } includeModified
 * @param { boolean } includeUnmodified
 * @returns { DeviceType[] }
 * @example
 * deviceRepo.listDeviceTypes(limit_as_number, offset_as_number, sort_as_string, filter_as_FilterCriteria_list, includeModified_as_boolean, includeUnmodified_as_boolean)
 */

/** 
 * @function deviceRepo#listHubDeviceIds
 * @param { string } id
 * @param { boolean } asLocalId
 * @returns { string[] }
 * @example
 * deviceRepo.listHubDeviceIds(id_as_string, asLocalId_as_boolean)
 */

/** 
 * @function deviceRepo#readDevice
 * @param { string } id
 * @returns { Device }
 * @example
 * deviceRepo.readDevice(id_as_string)
 */

/** 
 * @function deviceRepo#readDeviceByLocalId
 * @param { string } localId
 * @returns { Device }
 * @example
 * deviceRepo.readDeviceByLocalId(localId_as_string)
 */

/** 
 * @function deviceRepo#readDeviceGroup
 * @param { string } id
 * @returns { DeviceGroup }
 * @example
 * deviceRepo.readDeviceGroup(id_as_string)
 */

/** 
 * @function deviceRepo#readDeviceType
 * @param { string } id
 * @returns { DeviceType }
 * @example
 * deviceRepo.readDeviceType(id_as_string)
 */

/** 
 * @function deviceRepo#readHub
 * @param { string } id
 * @returns { Hub }
 * @example
 * deviceRepo.readHub(id_as_string)
 */

/** 
 * Exists checks if a process worker input exists
 * @function inputs#exists
 * @param { string } name
 * @returns { boolean }
 * @example
 * inputs.exists(name_as_string)
 */

/** 
 * Get value of a process worker input
 * @function inputs#get
 * @param { string } name
 * @returns { Object }
 * @example
 * inputs.get(name_as_string)
 */

/** 
 * List input values sorted by their names
 * @function inputs#list
 * @returns { Object[] }
 * @example
 * inputs.list()
 */

/** 
 * ListNames lists sorted input names
 * @function inputs#listNames
 * @returns { string[] }
 * @example
 * inputs.listNames()
 */

/** 
 * Get a process worker output
 * @function outputs#get
 * @param { string } name
 * @returns { Object }
 * @example
 * outputs.get(name_as_string)
 */

/** 
 * Set a process worker output
 * @function outputs#set
 * @param { string } name
 * @param { Object } value
 * @example
 * outputs.set(name_as_string, value_as_any)
 */

/** 
 * SetJson marshals the given value to json and sets it as a process worker output
 * @function outputs#setJson
 * @param { string } name
 * @param { Object } value
 * @example
 * outputs.setJson(name_as_string, value_as_any)
 */

/** 
 * GetDevicesWithServiceFromEntityString finds a list of iot-options where the entity is the same the input, but the Service field is set with those that match the input criteria
 * @function util#getDevicesWithServiceFromEntityString
 * @param { string } entityStr
 * @param { FilterCriteria[] } criteria
 * @returns { IotOption[] }
 * @example
 * util.getDevicesWithServiceFromEntityString(entityStr_as_string, criteria_as_FilterCriteria_list)
 */

/** 
 * GetDevicesWithServiceFromIotOption finds a list of iot-options where the entity is the same the input, but the Service field is set with those that match the input criteria
 * @function util#getDevicesWithServiceFromIotOption
 * @param { IotOption } entity
 * @param { FilterCriteria[] } criteria
 * @returns { IotOption[] }
 * @example
 * util.getDevicesWithServiceFromIotOption(entity_as_IotOption, criteria_as_FilterCriteria_list)
 */

/** 
 * @function util#getDevicesWithServiceFromIotOption
 * @param { IotOption } entity
 * @param { FilterCriteria[] } criteria
 * @returns { IotOption[] }
 * @example
 * util.getDevicesWithServiceFromIotOption(entity_as_IotOption, criteria_as_FilterCriteria_list)
 */

/** 
 * GetUserId returns the user-id of the executing user
 * @function util#getUserId
 * @returns { string }
 * @example
 * util.getUserId()
 */

/** 
 * GetUserToken returns a jwt-token for the executing user
 * @function util#getUserToken
 * @returns { string }
 * @example
 * util.getUserToken()
 */

/** 
 * GroupIotOptionsByDevice groups a list of model.IotOption by their device id; options that are not devices will be grouped under ""
 * @function util#groupIotOptionsByDevice
 * @param { IotOption[] } entities
 * @returns { Map<string,IotOption[]> }
 * @example
 * util.groupIotOptionsByDevice(entities_as_IotOption_list)
 */

/** 
 * GroupIotOptionsByService groups a list of IotOption by their service id; options that are not devices or dont hav a service-id will be grouped under ""
 * @function util#groupIotOptionsByService
 * @param { IotOption[] } entities
 * @returns { Map<string,IotOption[]> }
 * @example
 * util.groupIotOptionsByService(entities_as_IotOption_list)
 */

/** 
 * IsDeviceGroupIotOption checks if the input is a device-group
 * @function util#isDeviceGroupIotOption
 * @param { IotOption } entity
 * @returns { boolean }
 * @example
 * util.isDeviceGroupIotOption(entity_as_IotOption)
 */

/** 
 * IsDeviceGroupIotOptionStr checks if the input is a device-group
 * @function util#isDeviceGroupIotOptionStr
 * @param { string } entityStr
 * @returns { boolean }
 * @example
 * util.isDeviceGroupIotOptionStr(entityStr_as_string)
 */

/** 
 * IsDeviceIotOption checks if the input is a device
 * @function util#isDeviceIotOption
 * @param { IotOption } entity
 * @returns { boolean }
 * @example
 * util.isDeviceIotOption(entity_as_IotOption)
 */

/** 
 * IsDeviceIotOptionStr checks if the input is a device
 * @function util#isDeviceIotOptionStr
 * @param { string } entityStr
 * @returns { boolean }
 * @example
 * util.isDeviceIotOptionStr(entityStr_as_string)
 */

/** 
 * IsImportIotOption checks if the input is a import
 * @function util#isImportIotOption
 * @param { IotOption } entity
 * @returns { boolean }
 * @example
 * util.isImportIotOption(entity_as_IotOption)
 */

/** 
 * IsImportIotOptionStr checks if the input is a import
 * @function util#isImportIotOptionStr
 * @param { string } entityStr
 * @returns { boolean }
 * @example
 * util.isImportIotOptionStr(entityStr_as_string)
 */

/** 
 * DerefName returns the name of a smart-service instance variable referenced in parameter ref
 * @function variables#derefName
 * @param { string } ref
 * @returns { string }
 * @example
 * variables.derefName(ref_as_string)
 */

/** 
 * DerefTemplate replaces variable references in the input string with the corresponding variable values
 * @function variables#derefTemplate
 * @param { string } templ
 * @returns { string }
 * @example
 * variables.derefTemplate(templ_as_string)
 */

/** 
 * DerefValue returns the value of a smart-service instance variable referenced in parameter ref
 * @function variables#derefValue
 * @param { string } ref
 * @returns { Object }
 * @example
 * variables.derefValue(ref_as_string)
 */

/** 
 * Exists checks if a smart-service instance variable exists
 * @function variables#exists
 * @param { string } name
 * @returns { boolean }
 * @example
 * variables.exists(name_as_string)
 */

/** 
 * Read value of a smart-service instance variable
 * @function variables#read
 * @param { string } name
 * @returns { Object }
 * @example
 * variables.read(name_as_string)
 */

/** 
 * Ref creates a reference to a variable (e.g. "my_var_name" --> "{{.my_var_name}}")
throws exception if variable is unknown
 * @function variables#ref
 * @param { string } name
 * @returns { string }
 * @example
 * variables.ref(name_as_string)
 */

/** 
 * Write value as smart-service instance variable
 * @function variables#write
 * @param { string } name
 * @param { Object } value
 * @example
 * variables.write(name_as_string, value_as_any)
 */

/**
 * @typedef {Object} Aspect
 * @property { string } Id
 * @property { string } Name
 * @property { Aspect[] } SubAspects
 */

/**
 * @typedef {Object} AspectNode
 * @property { string } Id
 * @property { string } Name
 * @property { string } RootId
 * @property { string } ParentId
 * @property { string[] } ChildIds
 * @property { string[] } AncestorIds
 * @property { string[] } DescendentIds
 */

/**
 * @typedef {Object} Attribute
 * @property { string } Key
 * @property { string } Value
 * @property { string } Origin
 */

/**
 * @typedef {Object} Characteristic
 * @property { string } Id
 * @property { string } Name
 * @property { string } DisplayUnit
 * @property { Type } Type
 * @property { Object|null } MinValue
 * @property { Object|null } MaxValue
 * @property { Object|null[] } AllowedValues
 * @property { Object|null } Value
 * @property { Characteristic[] } SubCharacteristics
 */

/**
 * @typedef {Object} Concept
 * @property { string } Id
 * @property { string } Name
 * @property { string[] } CharacteristicIds
 * @property { string } BaseCharacteristicId
 * @property { ConverterExtension[] } Conversions
 */

/**
 * @typedef {Object} ConceptWithCharacteristics
 * @property { string } Id
 * @property { string } Name
 * @property { string } BaseCharacteristicId
 * @property { Characteristic[] } Characteristics
 * @property { ConverterExtension[] } Conversions
 */

/**
 * @typedef {Object} Configurable
 * @property { string } Path
 * @property { string } CharacteristicId
 * @property { AspectNode } AspectNode
 * @property { string } FunctionId
 * @property { Object|null } Value
 * @property { Type } Type
 */

/**
 * @typedef {Object} Content
 * @property { string } Id
 * @property { ContentVariable } ContentVariable
 * @property { Serialization } Serialization
 * @property { string } ProtocolSegmentId
 */

/**
 * @typedef {Object} ContentVariable
 * @property { string } Id
 * @property { string } Name
 * @property { boolean } IsVoid
 * @property { boolean } OmitEmpty
 * @property { Type } Type
 * @property { ContentVariable[] } SubContentVariables
 * @property { string } CharacteristicId
 * @property { Object|null } Value
 * @property { string[] } SerializationOptions
 * @property { string } UnitReference
 * @property { string } FunctionId
 * @property { string } AspectId
 */

/**
 * @typedef {Object} ConverterExtension
 * @property { string } From
 * @property { string } To
 * @property { number } Distance
 * @property { string } Formula
 * @property { string } PlaceholderName
 */

/**
 * @typedef {Object} Device
 * @property { string } Id
 * @property { string } LocalId
 * @property { string } Name
 * @property { Attribute[] } Attributes
 * @property { string } DeviceTypeId
 * @property { string } OwnerId
 */

/**
 * @typedef {Object} DeviceClass
 * @property { string } Id
 * @property { string } Image
 * @property { string } Name
 */

/**
 * @typedef {Object} DeviceGroup
 * @property { string } Id
 * @property { string } Name
 * @property { string } Image
 * @property { DeviceGroupFilterCriteria[] } Criteria
 * @property { string[] } DeviceIds
 * @property { string[] } CriteriaShort
 * @property { Attribute[] } Attributes
 * @property { string } AutoGeneratedByDevice
 */

/**
 * @typedef {Object} DeviceGroupFilterCriteria
 * @property { Interaction } Interaction
 * @property { string } FunctionId
 * @property { string } AspectId
 * @property { string } DeviceClassId
 */

/**
 * @typedef {Object} DeviceGroupSelection
 * @property { string } Id
 */

/**
 * @typedef {Object} DeviceSelection
 * @property { string } DeviceId
 * @property { string|null } ServiceId
 * @property { string|null } Path
 * @property { string|null } CharacteristicId
 */

/**
 * @typedef {Object} DeviceType
 * @property { string } Id
 * @property { string } Name
 * @property { string } Description
 * @property { ServiceGroup[] } ServiceGroups
 * @property { Service[] } Services
 * @property { string } DeviceClassId
 * @property { Attribute[] } Attributes
 */

/**
 * @typedef {Object} DeviceTypeSelectable
 * @property { string } DeviceTypeId
 * @property { Service[] } Services
 * @property { Map<string,ServicePathOption[]> } ServicePathOptions
 */

/**
 * @typedef {Object} FilterCriteria
 * @property { Interaction } Interaction
 * @property { string } FunctionId
 * @property { string } DeviceClassId
 * @property { string } AspectId
 */

/**
 * @typedef {Object} FunctionType
 * @property { string } Id
 * @property { string } Name
 * @property { string } DisplayName
 * @property { string } Description
 * @property { string } ConceptId
 * @property { string } RdfType
 */

/**
 * @typedef {Object} GenericEventSource
 * @property { string } FilterType
 * @property { string } FilterIds
 * @property { string } Topic
 * @property { string } Path
 * @property { string|null } CharacteristicId
 */

/**
 * @typedef {Object} Hub
 * @property { string } Id
 * @property { string } Name
 * @property { string } Hash
 * @property { string[] } DeviceLocalIds
 * @property { string[] } DeviceIds
 * @property { string } OwnerId
 */

/**
 * @typedef {Object} ImportSelection
 * @property { string } Id
 * @property { string|null } Path
 * @property { string|null } CharacteristicId
 */

/**
 * @typedef {Object} IotOption
 * @property { DeviceSelection|null } DeviceSelection
 * @property { DeviceGroupSelection|null } DeviceGroupSelection
 * @property { ImportSelection|null } ImportSelection
 * @property { GenericEventSource|null } GenericEventSource
 */

/**
 * @typedef {Object} Service
 * @property { string } Id
 * @property { string } LocalId
 * @property { string } Name
 * @property { string } Description
 * @property { Interaction } Interaction
 * @property { string } ProtocolId
 * @property { Content[] } Inputs
 * @property { Content[] } Outputs
 * @property { Attribute[] } Attributes
 * @property { string } ServiceGroupKey
 */

/**
 * @typedef {Object} ServiceGroup
 * @property { string } Key
 * @property { string } Name
 * @property { string } Description
 */

/**
 * @typedef {Object} ServicePathOption
 * @property { string } ServiceId
 * @property { string } Path
 * @property { string } CharacteristicId
 * @property { AspectNode } AspectNode
 * @property { string } FunctionId
 * @property { boolean } IsVoid
 * @property { Object|null } Value
 * @property { boolean } IsControllingFunction
 * @property { Configurable[] } Configurables
 * @property { Type } Type
 * @property { Interaction } Interaction
 */
