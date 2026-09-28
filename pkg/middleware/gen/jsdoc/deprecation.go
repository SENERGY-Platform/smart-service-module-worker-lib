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
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

// deprecationMarker is the word a go comment opens with to mark the field below or beside it
// as deprecated, the convention the platform models are written with:
//
//	AspectId  string   `json:"aspect_id"` //deprecated: alias for a single element AspectIds
//	AspectIds []string `json:"aspect_ids,omitempty"`
//
// What follows the word is the note, and there is often nothing to follow: a field that has
// no replacement worth naming is usually marked with a bare "deprecated". The marking is
// therefore the word, not the colon.
const deprecationMarker = "deprecated"

// deprecationOf answers whether the go source marks a struct field as deprecated, and with
// what note. The declaring type is needed rather than the type a script reads the field on,
// because a promoted field is declared by the embedded struct.
//
// The source is the one place that knows: the model that declares the field also decides
// that it is deprecated, and says so in the same commit. A rule over the field names would
// only recognise the aliases that happen to have a replacement named after them, and a list
// kept in this repository would be a second place to maintain, silently stale from the next
// deprecation in the models on.
func deprecationOf(declaring reflect.Type, fieldName string) (note string, isDeprecated bool) {
	if declaring == nil || declaring.PkgPath() == "" || declaring.Name() == "" {
		return "", false //an anonymous struct has no declaration to look up
	}
	note, isDeprecated = notesOfPackage(declaring.PkgPath())[declaring.Name()+"."+fieldName]
	return note, isDeprecated
}

// notes maps "<type name>.<field name>" to the deprecation note of that field, per package.
// A package is parsed once, however many of its types are asked about.
var notes = struct {
	sync.Mutex
	byPackage map[string]map[string]string
}{byPackage: map[string]map[string]string{}}

func notesOfPackage(pkgPath string) map[string]string {
	notes.Lock()
	defer notes.Unlock()
	if known, ok := notes.byPackage[pkgPath]; ok {
		return known
	}
	result := map[string]string{}
	//also remembers an empty result: a package whose source is not readable stays that way
	notes.byPackage[pkgPath] = result

	pkg, err := build.Import(pkgPath, ".", build.FindOnly)
	if err != nil {
		return result
	}
	parsed, err := parser.ParseDir(token.NewFileSet(), pkg.Dir, nil, parser.ParseComments)
	if err != nil {
		return result
	}
	for _, astPackage := range parsed {
		for _, file := range astPackage.Files {
			collectDeprecations(file, result)
		}
	}
	return result
}

func collectDeprecations(file *ast.File, result map[string]string) {
	ast.Inspect(file, func(node ast.Node) bool {
		typeSpec, isTypeSpec := node.(*ast.TypeSpec)
		if !isTypeSpec {
			return true
		}
		structType, isStruct := typeSpec.Type.(*ast.StructType)
		if !isStruct {
			return true
		}
		for _, field := range structType.Fields.List {
			note, isDeprecated := deprecationComment(field)
			if !isDeprecated {
				continue
			}
			for _, name := range field.Names {
				result[typeSpec.Name.Name+"."+name.Name] = note
			}
		}
		return true
	})
}

// deprecationComment reads the deprecation out of the comment above the field and out of the
// one beside it, the two places a field can carry one.
func deprecationComment(field *ast.Field) (note string, isDeprecated bool) {
	for _, comment := range []string{field.Doc.Text(), field.Comment.Text()} {
		for _, line := range strings.Split(comment, "\n") {
			if note, isDeprecated := deprecationOfLine(line); isDeprecated {
				return note, true
			}
		}
	}
	return "", false
}

// deprecationOfLine reads one comment line. The line marks the field when it opens with the
// word "deprecated" - alone, or followed by a note that the usual ":" or "-" may introduce.
// The word has to end there: a line about a "deprecatedName" says nothing about this field.
func deprecationOfLine(line string) (note string, isDeprecated bool) {
	line = strings.TrimSpace(line)
	if len(line) < len(deprecationMarker) {
		return "", false
	}
	if !strings.EqualFold(line[:len(deprecationMarker)], deprecationMarker) {
		return "", false
	}
	rest := line[len(deprecationMarker):]
	if rest != "" {
		if first, _ := utf8.DecodeRuneInString(rest); unicode.IsLetter(first) || unicode.IsDigit(first) || first == '_' {
			return "", false
		}
	}
	return strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(rest), ":,-")), true
}

// declaringStruct returns the struct type that declares field, which is t itself unless the
// field is promoted from an embedded struct.
func declaringStruct(t reflect.Type, field reflect.StructField) reflect.Type {
	if len(field.Index) == 0 {
		return t
	}
	for _, i := range field.Index[:len(field.Index)-1] {
		t = t.Field(i).Type
		if t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
	}
	return t
}
