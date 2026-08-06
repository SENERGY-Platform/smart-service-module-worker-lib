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

package typescript

import (
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"
)

const pathToScriptenv = "../../scriptenv"

func TestToTsType(t *testing.T) {
	//the left side is the complete set of type expressions the jsdoc generator produces
	cases := map[string]string{
		"":                        "void",
		"string":                  "string",
		"number":                  "number",
		"boolean":                 "boolean",
		"Aspect":                  "Aspect",
		"Aspect[]":                "Aspect[]",
		"string[]":                "string[]",
		"Object":                  "any",
		"Object|null":             "any",
		"Object|null[]":           "any[]",
		"string|null":             "string | null",
		"DeviceSelection|null":    "DeviceSelection | null",
		"Map<string,IotOption[]>": "Record<string, IotOption[]>",
		"Map<string,Object>":      "Record<string, any>",
		"Map<string,string>[]":    "Record<string, string>[]",
	}
	for jsDocType, expected := range cases {
		t.Run(jsDocType, func(t *testing.T) {
			if actual := ToTsType(jsDocType); actual != expected {
				t.Errorf("expected %#v, got %#v", expected, actual)
			}
		})
	}
}

func TestInterfacePropertiesUseRuntimeNames(t *testing.T) {
	interfaces := GetInterfaces()
	i := slices.IndexFunc(interfaces, func(e Interface) bool { return e.Name == "Aspect" })
	if i < 0 {
		t.Fatalf("no interface generated for Aspect")
	}
	expected := []Property{
		{Name: "id", Type: "string"},
		{Name: "name", Type: "string"},
		{Name: "sub_aspects", Type: "Aspect[]"},
	}
	if !reflect.DeepEqual(interfaces[i].Properties, expected) {
		t.Errorf("expected %#v, got %#v", expected, interfaces[i].Properties)
	}
}

// TestEveryReferencedTypeIsDeclared keeps the declarations compilable: typescript
// rejects a reference to a type that is never declared, where jsdoc would tolerate it.
func TestEveryReferencedTypeIsDeclared(t *testing.T) {
	interfaces := GetInterfaces()
	declared := []string{}
	for _, element := range interfaces {
		declared = append(declared, element.Name)
	}

	assertDeclared := func(t *testing.T, tsType string) {
		t.Helper()
		for _, name := range typeNamesOf(tsType) {
			if !slices.Contains(declared, name) {
				t.Errorf("%v references undeclared type %v", tsType, name)
			}
		}
	}

	for _, element := range interfaces {
		for _, property := range element.Properties {
			t.Run(element.Name+"."+property.Name, func(t *testing.T) {
				assertDeclared(t, property.Type)
			})
		}
	}
	for _, namespace := range GetNamespaces(pathToScriptenv) {
		for _, method := range namespace.Methods {
			t.Run(namespace.Name+"."+method.Name, func(t *testing.T) {
				assertDeclared(t, method.Result)
				assertDeclared(t, method.Params)
			})
		}
	}
}

// TestNoDuplicateMembers guards against declaring the same member twice, which
// typescript rejects but jsdoc does not.
func TestNoDuplicateMembers(t *testing.T) {
	for _, element := range GetInterfaces() {
		names := []string{}
		for _, property := range element.Properties {
			names = append(names, property.Name)
		}
		assertUnique(t, "interface "+element.Name, names)
	}
	for _, namespace := range GetNamespaces(pathToScriptenv) {
		names := []string{}
		for _, method := range namespace.Methods {
			names = append(names, method.Name)
		}
		assertUnique(t, "namespace "+namespace.Name, names)
	}
}

func TestGenerateTypescriptDeclarations(t *testing.T) {
	result, err := GenerateTypescriptDeclarations(pathToScriptenv)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"interface Aspect {\n    id: string;\n    name: string;\n    sub_aspects: Aspect[];\n}",
		"declare const deviceRepo: {",
		"    getAspect(id: string): Aspect;",
		"    setJson(name: string, value: any): void;",
	} {
		if !strings.Contains(result, expected) {
			t.Errorf("missing %#v in:\n%v", expected, result)
		}
	}
}

func assertUnique(t *testing.T, context string, names []string) {
	t.Helper()
	seen := []string{}
	for _, name := range names {
		if slices.Contains(seen, name) {
			t.Errorf("%v declares %v more than once", context, name)
		}
		seen = append(seen, name)
	}
}

var typeNameExpr = regexp.MustCompile(`[A-Za-z_][A-Za-z_0-9]*`)

// typeNamesOf returns the named types a typescript type expression or parameter list
// refers to, ignoring parameter names and anything typescript declares itself.
func typeNamesOf(tsType string) (result []string) {
	builtin := []string{"string", "number", "boolean", "any", "void", "null", "Record"}
	for _, part := range strings.Split(tsType, ",") {
		//drop the parameter name of a "name: type" pair
		if index := strings.Index(part, ":"); index >= 0 {
			part = part[index+1:]
		}
		for _, name := range typeNameExpr.FindAllString(part, -1) {
			if !slices.Contains(builtin, name) {
				result = append(result, name)
			}
		}
	}
	return result
}
