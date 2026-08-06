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

package jsdoc

import (
	"fmt"
	"github.com/SENERGY-Platform/device-repository/lib/model"
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/middleware/gen/util"
	model2 "github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/model"
	"github.com/dop251/goja"
	"reflect"
	"slices"
	"strings"
)

// fieldNameMapper must match the mapper runScript() installs in pkg/middleware/scripts.go.
// scripts see struct fields under their json tag name, not under their go field name,
// and fields without a tag usable as js identifier are not visible to scripts at all.
var fieldNameMapper = goja.TagFieldNameMapper("json", true)

type TypeDefField struct {
	Name string
	Type string
}

type TypeDef struct {
	Name   string
	Fields []TypeDefField
}

func typeDefSources() []interface{} {
	return []interface{}{
		models.AspectNode{},
		models.Aspect{},
		models.Characteristic{},
		models.ConceptWithCharacteristics{},
		models.Concept{},
		models.DeviceClass{},
		models.Function{},
		model.FilterCriteria{},
		model.DeviceTypeSelectable{},
		models.Service{},
		models.DeviceType{},
		models.Device{},
		models.DeviceGroup{},
		models.Hub{},
		models.Location{},
		model2.IotOption{},
	}
}

func GetTypeDefs() (result []TypeDef) {
	for _, t := range typeDefSources() {
		result = append(result, GetTypeDef(t)...)
	}

	result = filterDuplicateTypeDef(result)

	slices.SortFunc(result, func(a, b TypeDef) int {
		return strings.Compare(a.Name, b.Name)
	})

	return
}

func GetTypeDef(obj interface{}) (result []TypeDef) {
	return getTypeDef(reflect.TypeOf(obj), nil)
}

func getDeepStruct(t reflect.Type) reflect.Type {
	switch t.Kind() {
	case reflect.Struct:
		return t
	case reflect.Map:
		return getDeepStruct(t.Elem())
	case reflect.Slice:
		return getDeepStruct(t.Elem())
	case reflect.Pointer:
		return getDeepStruct(t.Elem())
	}
	return nil
}

func getTypeDef(t reflect.Type, done []string) (result []TypeDef) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("ERROR: t.Name=%v, Err=%v", t.Name(), err)
			panic(err)
		}
	}()
	this := TypeDef{
		Name:   util.ToJsDocName(t.Name()),
		Fields: nil,
	}
	if slices.Contains(done, this.Name) {
		return nil
	}
	done = append(done, this.Name)
	sub := []TypeDef{}
	for _, field := range reflect.VisibleFields(t) {
		name := getFieldName(t, field)
		if name == "" {
			continue //not visible to scripts
		}
		if slices.ContainsFunc(this.Fields, func(f TypeDefField) bool { return f.Name == name }) {
			continue //shadowed by a field of the same name closer to the surface
		}
		this.Fields = append(this.Fields, TypeDefField{
			Name: name,
			Type: getTypeName(field.Type),
		})
		subStruct := getDeepStruct(field.Type)
		if subStruct != nil {
			sub = append(sub, getTypeDef(subStruct, done)...)
		}
	}
	result = []TypeDef{this}
	result = append(result, sub...)
	return result
}

// getFieldName returns the name scripts use to access the field, or "" if the field is
// not visible to scripts.
func getFieldName(t reflect.Type, field reflect.StructField) string {
	if !field.IsExported() {
		//an unexported field is never visible, even if it carries a json tag
		//(a case go vet rejects, which is why no test covers it)
		return ""
	}
	return fieldNameMapper.FieldName(t, field)
}

func getTypeName(t reflect.Type) string {
	if t.Kind() == reflect.Struct && t.Name() != "" {
		return util.ToJsDocName(t.Name())
	}
	//a named non struct type like models.Interaction has no typedef of its own,
	//scripts see its underlying value
	if primitive := getPrimitiveTypeName(t.Kind()); primitive != "" {
		return primitive
	}
	switch t.Kind() {
	case reflect.Slice:
		return getTypeName(t.Elem()) + "[]"
	case reflect.Map:
		return "Map<string," + getTypeName(t.Elem()) + ">"
	case reflect.Pointer:
		return getTypeName(t.Elem()) + "|null"
	case reflect.Interface:
		return "Object|null"
	}
	return ""
}

func getPrimitiveTypeName(kind reflect.Kind) string {
	switch kind {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	}
	return ""
}

func filterDuplicateTypeDef(list []TypeDef) (result []TypeDef) {
	for _, e := range list {
		if !slices.ContainsFunc(result, func(def TypeDef) bool {
			return def.Name == e.Name
		}) {
			result = append(result, e)
		}
	}
	return result
}
