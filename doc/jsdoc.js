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
 */

/** 
 * @function deviceRepo#getAspectNode
 * @param { string } id
 * @returns { AspectNode }
 */

/** 
 * @function deviceRepo#getAspectNodes
 * @returns { AspectNode[] }
 */

/** 
 * @function deviceRepo#getAspectNodesByIdList
 * @param { string[] } ids
 * @returns { AspectNode[] }
 */

/** 
 * @function deviceRepo#getAspectNodesMeasuringFunctions
 * @param { string } id
 * @param { boolean } ancestors
 * @param { boolean } descendants
 * @returns { Function[] }
 */

/** 
 * @function deviceRepo#getAspectNodesWithMeasuringFunction
 * @param { boolean } ancestors
 * @param { boolean } descendants
 * @returns { AspectNode[] }
 */

/** 
 * @function deviceRepo#getAspects
 * @returns { Aspect[] }
 */

/** 
 * @function deviceRepo#getAspectsWithMeasuringFunction
 * @param { boolean } ancestors
 * @param { boolean } descendants
 * @returns { Aspect[] }
 */

/** 
 * @function deviceRepo#getCharacteristic
 * @param { string } id
 * @returns { Characteristic }
 */

/** 
 * @function deviceRepo#getConceptWithCharacteristics
 * @param { string } id
 * @returns { ConceptWithCharacteristics }
 */

/** 
 * @function deviceRepo#getConceptWithoutCharacteristics
 * @param { string } id
 * @returns { Concept }
 */

/** 
 * @function deviceRepo#getDeviceClass
 * @param { string } id
 * @returns { DeviceClass }
 */

/** 
 * @function deviceRepo#getDeviceClasses
 * @returns { DeviceClass[] }
 */

/** 
 * @function deviceRepo#getDeviceClassesControllingFunctions
 * @param { string } id
 * @returns { Function[] }
 */

/** 
 * @function deviceRepo#getDeviceClassesFunctions
 * @param { string } id
 * @returns { Function[] }
 */

/** 
 * @function deviceRepo#getDeviceClassesWithControllingFunctions
 * @returns { DeviceClass[] }
 */

/** 
 * @function deviceRepo#getDeviceTypeSelectables
 * @param { FilterCriteria[] } query
 * @param { string } pathPrefix
 * @param { boolean } includeModified
 * @param { boolean } servicesMustMatchAllCriteria
 * @returns { DeviceTypeSelectable[] }
 */

/** 
 * @function deviceRepo#getFunction
 * @param { string } id
 * @returns { Function }
 */

/** 
 * @function deviceRepo#getFunctionsByType
 * @param { string } rdfType
 * @returns { Function[] }
 */

/** 
 * @function deviceRepo#getLeafCharacteristics
 * @returns { Characteristic[] }
 */

/** 
 * @function deviceRepo#getLocation
 * @param { string } id
 * @returns { Location }
 */

/** 
 * @function deviceRepo#getService
 * @param { string } id
 * @returns { Service }
 */

/** 
 * @function deviceRepo#listDeviceTypes
 * @param { int64 } limit
 * @param { int64 } offset
 * @param { string } sort
 * @param { FilterCriteria[] } filter
 * @param { boolean } includeModified
 * @param { boolean } includeUnmodified
 * @returns { DeviceType[] }
 */

/** 
 * @function deviceRepo#listHubDeviceIds
 * @param { string } id
 * @param { boolean } asLocalId
 * @returns { string[] }
 */

/** 
 * @function deviceRepo#readDevice
 * @param { string } id
 * @returns { Device }
 */

/** 
 * @function deviceRepo#readDeviceByLocalId
 * @param { string } localId
 * @returns { Device }
 */

/** 
 * @function deviceRepo#readDeviceGroup
 * @param { string } id
 * @returns { DeviceGroup }
 */

/** 
 * @function deviceRepo#readDeviceType
 * @param { string } id
 * @returns { DeviceType }
 */

/** 
 * @function deviceRepo#readHub
 * @param { string } id
 * @returns { Hub }
 */

/** 
 * Exists checks if a process worker input exists
 * @function inputs#exists
 * @param { string } name
 * @returns { boolean }
 */

/** 
 * Get value of a process worker input
 * @function inputs#get
 * @param { string } name
 * @returns { Object }
 */

/** 
 * List input values sorted by their names
 * @function inputs#list
 * @returns { Object[] }
 */

/** 
 * ListNames lists sorted input names
 * @function inputs#listNames
 * @returns { string[] }
 */

/** 
 * Get a process worker output
 * @function outputs#get
 * @param { string } name
 * @returns { Object }
 */

/** 
 * Set a process worker output
 * @function outputs#set
 * @param { string } name
 * @param { Object } value
 */

/** 
 * SetJson marshals the given value to json and sets it as a process worker output
 * @function outputs#setJson
 * @param { string } name
 * @param { Object } value
 */

/** 
 * @function util#getDevicesWithServiceFromEntityString
 * @param { string } entityStr
 * @param { FilterCriteria[] } criteria
 * @returns { IotOption[] }
 */

/** 
 * @function util#getDevicesWithServiceFromIotOption
 * @param { IotOption } entity
 * @param { FilterCriteria[] } criteria
 * @returns { IotOption[] }
 */

/** 
 * @function util#getDevicesWithServiceFromIotOption
 * @param { IotOption } entity
 * @param { FilterCriteria[] } criteria
 * @returns { IotOption[] }
 */

/** 
 * @function util#groupIotOptionsByDevice
 * @param { IotOption[] } entities
 * @returns {  }
 */

/** 
 * @function util#groupIotOptionsByService
 * @param { IotOption[] } entities
 * @returns {  }
 */

/** 
 * @function util#isDeviceGroupIotOption
 * @param { IotOption } entity
 * @returns { boolean }
 */

/** 
 * @function util#isDeviceGroupIotOptionStr
 * @param { string } entityStr
 * @returns { boolean }
 */

/** 
 * @function util#isDeviceIotOption
 * @param { IotOption } entity
 * @returns { boolean }
 */

/** 
 * @function util#isDeviceIotOptionStr
 * @param { string } entityStr
 * @returns { boolean }
 */

/** 
 * @function util#isImportIotOption
 * @param { IotOption } entity
 * @returns { boolean }
 */

/** 
 * @function util#isImportIotOptionStr
 * @param { string } entityStr
 * @returns { boolean }
 */

/** 
 * DerefName returns the name of a smart-service instance variable referenced in parameter ref
 * @function variables#derefName
 * @param { string } ref
 * @returns { string }
 */

/** 
 * DerefTemplate replaces variable references in the input string with the corresponding variable values
 * @function variables#derefTemplate
 * @param { string } templ
 * @returns { string }
 */

/** 
 * DerefValue returns the value of a smart-service instance variable referenced in parameter ref
 * @function variables#derefValue
 * @param { string } ref
 * @returns { Object }
 */

/** 
 * Exists checks if a smart-service instance variable exists
 * @function variables#exists
 * @param { string } name
 * @returns { boolean }
 */

/** 
 * Read value of a smart-service instance variable
 * @function variables#read
 * @param { string } name
 * @returns { Object }
 */

/** 
 * Ref creates a reference to a variable (e.g. "my_var_name" --> "{{.my_var_name}}")
throws exception if variable is unknown
 * @function variables#ref
 * @param { string } name
 * @returns { string }
 */

/** 
 * Write value as smart-service instance variable
 * @function variables#write
 * @param { string } name
 * @param { Object } value
 */
