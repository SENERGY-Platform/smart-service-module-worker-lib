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
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/SENERGY-Platform/device-repository/v2/lib/model"
)

type DeprecationTestStruct struct {
	DeprecationTestEmbedded `json:"emb"`

	BesideTheField string `json:"beside_the_field"` //deprecated: use something else

	//deprecated: above the field
	AboveTheField string `json:"above_the_field"`

	//Deprecated: the go convention capitalises it
	Capitalised string `json:"capitalised"`

	//deprecated:
	WithoutANote string `json:"without_a_note"`

	Bare string `json:"bare"` //deprecated

	//Deprecated
	BareAndCapitalised string `json:"bare_and_capitalised"`

	//deprecated - a dash introduces the note just as well
	WithADash string `json:"with_a_dash"`

	//deprecated use the other field
	WithoutPunctuation string `json:"without_punctuation"`

	//deprecatedName was the old spelling of this field
	NamesADeprecation string `json:"names_a_deprecation"`

	//a comment that is not a deprecation
	Current string `json:"current"`

	Untouched string `json:"untouched"`
}

type DeprecationTestEmbedded struct {
	Promoted string `json:"promoted"` //deprecated: declared by the embedded struct
}

// The models mark a deprecated field in its go comment, and that is where the generator has
// to read it: the marking then lives in the one place that decides it, and a field without a
// replacement to be recognised by is marked just as well as an alias that has one.
func TestDeprecationIsReadFromTheGoComment(t *testing.T) {
	defs := GetTypeDef(DeprecationTestStruct{})
	i := slices.IndexFunc(defs, func(d TypeDef) bool { return d.Name == "DeprecationTestStruct" })
	if i < 0 {
		t.Fatalf("no typedef generated, got %#v", defs)
	}

	byName := map[string]TypeDefField{}
	for _, field := range defs[i].Fields {
		byName[field.Name] = field
	}
	expected := map[string]TypeDefField{
		"beside_the_field": {Name: "beside_the_field", Type: "string", Deprecated: true, DeprecationNote: "use something else"},
		"above_the_field":  {Name: "above_the_field", Type: "string", Deprecated: true, DeprecationNote: "above the field"},
		"capitalised":      {Name: "capitalised", Type: "string", Deprecated: true, DeprecationNote: "the go convention capitalises it"},
		//a marked field stays marked even when there is nothing to say about it, which is
		//the usual shape when there is no replacement worth naming
		"without_a_note":       {Name: "without_a_note", Type: "string", Deprecated: true},
		"bare":                 {Name: "bare", Type: "string", Deprecated: true},
		"bare_and_capitalised": {Name: "bare_and_capitalised", Type: "string", Deprecated: true},
		"with_a_dash":          {Name: "with_a_dash", Type: "string", Deprecated: true, DeprecationNote: "a dash introduces the note just as well"},
		"without_punctuation":  {Name: "without_punctuation", Type: "string", Deprecated: true, DeprecationNote: "use the other field"},
		//the comment names something that is deprecated, it does not mark this field
		"names_a_deprecation": {Name: "names_a_deprecation", Type: "string"},
		"current":             {Name: "current", Type: "string"},
		"untouched":           {Name: "untouched", Type: "string"},
		"emb":                 {Name: "emb", Type: "DeprecationTestEmbedded"},
		//the comment sits on the embedded struct, not on the type the script reads it from
		"promoted": {Name: "promoted", Type: "string", Deprecated: true, DeprecationNote: "declared by the embedded struct"},
	}
	if !reflect.DeepEqual(byName, expected) {
		t.Errorf("expected %#v, got %#v", expected, byName)
	}
}

// The same on the type a script actually writes, so that the rule is tied to the model it
// exists for and not only to a fixture of this test.
func TestFilterCriteriaAspectIdIsDeprecated(t *testing.T) {
	defs := GetTypeDef(model.FilterCriteria{})
	if len(defs) != 1 {
		t.Fatalf("expected exactly one typedef, got %#v", defs)
	}
	expected := []TypeDefField{
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
	if !reflect.DeepEqual(defs[0].Fields, expected) {
		t.Errorf("expected %#v, got %#v", expected, defs[0].Fields)
	}
}

// The note has to reach the generated jsdoc, that is the only place a script author reads it.
func TestGenerateJsDocMarksDeprecatedFields(t *testing.T) {
	result, err := GenerateJsDoc("../../scriptenv")
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		" * @property { string } aspect_id deprecated: alias for a single element AspectIds; normalized into AspectIds at the controller boundary",
		" * @property { string[] } aspect_ids",
		" * @property { AspectNode } aspect_node deprecated: alias for a single element AspectNodes; holds the node with the alphabetically first id",
		" * @property { AspectNode[] } aspect_nodes",
	} {
		if !strings.Contains(result, expected) {
			t.Errorf("missing %#v in:\n%v", expected, result)
		}
	}
	//a field the models do not mark must not pick up a note
	if strings.Contains(result, "function_id deprecated") {
		t.Error("function_id was marked as deprecated")
	}
}
