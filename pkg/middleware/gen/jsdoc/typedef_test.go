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

package jsdoc

import (
	"github.com/SENERGY-Platform/models/go/models"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/middleware/gen/util"
	"github.com/dop251/goja"
	"reflect"
	"slices"
	"sort"
	"testing"
)

func TestGetTypeDefUsesJsonTagNames(t *testing.T) {
	defs := GetTypeDef(models.Aspect{})
	if len(defs) != 1 {
		t.Fatalf("expected exactly one typedef, got %#v", defs)
	}
	expected := []TypeDefField{
		{Name: "id", Type: "string"},
		{Name: "name", Type: "string"},
		{Name: "sub_aspects", Type: "Aspect[]"},
	}
	if !reflect.DeepEqual(defs[0].Fields, expected) {
		t.Errorf("expected %#v, got %#v", expected, defs[0].Fields)
	}
}

type TypeDefTestStruct struct {
	Plain               string `json:"plain"`
	WithOptions         string `json:"with_options,omitempty"`
	Untagged            string
	Excluded            string `json:"-"`
	NoIdentifier        string `json:"no-identifier"`
	EmptyTag            string `json:""`
	unexported          string
	Hidden              TypeDefTestHidden
	TypeDefTestEmbedded `json:"emb"`
}

type TypeDefTestHidden struct {
	Foo string `json:"foo"`
}

type TypeDefTestEmbedded struct {
	Promoted string `json:"promoted"`
}

func TestGetTypeDefFieldVisibility(t *testing.T) {
	defs := GetTypeDef(TypeDefTestStruct{})
	def, found := findTypeDef(defs, "TypeDefTestStruct")
	if !found {
		t.Fatalf("missing typedef, got %#v", defs)
	}

	//an untagged, excluded, unexported or not-a-js-identifier field is unreachable
	//for scripts, so documenting it would send script authors down a dead end
	expected := []string{"emb", "plain", "promoted", "with_options"}
	actual := fieldNames(def)
	sort.Strings(actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %#v, got %#v", expected, actual)
	}

	//the field visibility rules above are goja's, so check them against goja itself
	fromRuntime := scriptVisibleProperties(t, TypeDefTestStruct{})
	sort.Strings(fromRuntime)
	if !reflect.DeepEqual(fromRuntime, expected) {
		t.Errorf("script runtime disagrees with the expectation of this test: expected %#v, got %#v", expected, fromRuntime)
	}

	//the type of a documented field needs a typedef of its own
	if _, found = findTypeDef(defs, "TypeDefTestEmbedded"); !found {
		t.Errorf("missing typedef of embedded struct, got %#v", defs)
	}

	//the type of a hidden field does not
	if _, found = findTypeDef(defs, "TypeDefTestHidden"); found {
		t.Errorf("unexpected typedef of hidden field type, got %#v", defs)
	}
}

// TestTypeDefsMatchScriptRuntime guards against the generated documentation drifting
// from what pkg/middleware/scripts.go actually exposes to scripts.
func TestTypeDefsMatchScriptRuntime(t *testing.T) {
	defs := GetTypeDefs()
	for _, source := range typeDefSources() {
		name := util.ToJsDocName(reflect.TypeOf(source).Name())
		t.Run(name, func(t *testing.T) {
			def, found := findTypeDef(defs, name)
			if !found {
				t.Fatalf("no typedef generated for %v", name)
			}
			expected := scriptVisibleProperties(t, source)
			actual := fieldNames(def)
			sort.Strings(expected)
			sort.Strings(actual)
			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("documented properties differ from the properties scripts see\nscript runtime: %#v\ndocumented:     %#v", expected, actual)
			}
		})
	}
}

// scriptVisibleProperties reports the property names a script can read off value,
// using the same runtime setup as runScript() in pkg/middleware/scripts.go.
func scriptVisibleProperties(t *testing.T, value interface{}) []string {
	t.Helper()
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	err := vm.Set("value", value)
	if err != nil {
		t.Fatal(err)
	}
	keys, err := vm.RunString(`Object.keys(value)`)
	if err != nil {
		t.Fatal(err)
	}
	result := []string{}
	methods := scriptVisibleMethods(value)
	for _, key := range keys.Export().([]interface{}) {
		name := key.(string)
		if slices.Contains(methods, name) {
			continue //methods are documented as @function, not as @property
		}
		result = append(result, name)
	}
	return result
}

func scriptVisibleMethods(value interface{}) (result []string) {
	mapper := goja.TagFieldNameMapper("json", true)
	//goja makes the value addressable, so pointer receiver methods are exposed too
	for _, t := range []reflect.Type{reflect.TypeOf(value), reflect.PointerTo(reflect.TypeOf(value))} {
		for i := 0; i < t.NumMethod(); i++ {
			result = append(result, mapper.MethodName(t, t.Method(i)))
		}
	}
	return result
}

func findTypeDef(list []TypeDef, name string) (TypeDef, bool) {
	i := slices.IndexFunc(list, func(def TypeDef) bool { return def.Name == name })
	if i < 0 {
		return TypeDef{}, false
	}
	return list[i], true
}

func fieldNames(def TypeDef) (result []string) {
	result = []string{}
	for _, field := range def.Fields {
		result = append(result, field.Name)
	}
	return result
}
