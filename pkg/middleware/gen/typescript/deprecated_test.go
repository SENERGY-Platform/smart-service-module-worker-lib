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
	"slices"
	"strings"
	"testing"
)

// A criteria a script writes may name several aspects in aspect_ids; aspect_id is the
// deprecated alias for a list with one element. An editor reads that from the @deprecated
// tag, which is the reason the declarations carry a comment at all. The note behind the tag
// is what the go source of the model says, so it does not have to be kept in step here.
func TestDeclarationsMarkDeprecatedFields(t *testing.T) {
	result, err := GenerateTypescriptDeclarations(pathToScriptenv)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"interface FilterCriteria {\n" +
			"    interaction: string;\n" +
			"    function_id: string;\n" +
			"    device_class_id: string;\n" +
			"    /** @deprecated alias for a single element AspectIds; normalized into AspectIds at the controller boundary */\n" +
			"    aspect_id: string;\n" +
			"    aspect_ids: string[];\n}",
		"    /** @deprecated alias for a single element AspectNodes; holds the node with the alphabetically first id */\n" +
			"    aspect_node: AspectNode;\n" +
			"    aspect_nodes: AspectNode[];",
	} {
		if !strings.Contains(result, expected) {
			t.Errorf("missing %#v in:\n%v", expected, result)
		}
	}
}

// The comment is only emitted for the deprecated property, so the one carrying it has to be
// the one the jsdoc generator flagged rather than every property of the interface.
func TestInterfacePropertiesCarryDeprecation(t *testing.T) {
	interfaces := GetInterfaces()
	i := slices.IndexFunc(interfaces, func(e Interface) bool { return e.Name == "FilterCriteria" })
	if i < 0 {
		t.Fatalf("no interface generated for FilterCriteria")
	}
	expected := []Property{
		{Name: "interaction", Type: "string"},
		{Name: "function_id", Type: "string"},
		{Name: "device_class_id", Type: "string"},
		{
			Name: "aspect_id", Type: "string",
			Deprecated:      true,
			DeprecationNote: "alias for a single element AspectIds; normalized into AspectIds at the controller boundary",
		},
		{Name: "aspect_ids", Type: "string[]"},
	}
	if !reflect.DeepEqual(interfaces[i].Properties, expected) {
		t.Errorf("expected %#v, got %#v", expected, interfaces[i].Properties)
	}
}
