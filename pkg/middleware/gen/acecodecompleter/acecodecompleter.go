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

package acecodecompleter

import (
	_ "embed"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/middleware/scriptenv"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

//go:embed statement.tmpl
var statementTmplStr string

//go:embed assignment.tmpl
var assignmentTmplStr string

//go:embed completer.tmpl
var completerTmplStr string

func GenerateTsAceCodeCompleter(pathToScriptenv string) (string, error) {
	scriptEnvMapping := getScriptEnvMapping()

	f, err := parser.ParseDir(token.NewFileSet(), pathToScriptenv, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	statementTemplInputs := []map[string]interface{}{}
	assignmentTmplInputs := []map[string]interface{}{}

	for key, typeName := range scriptEnvMapping {
		for _, method := range filterMethods(f, typeName) {
			templInputs := map[string]interface{}{
				"prefix":     key,
				"method":     method.Name.Name,
				"inputs":     "",
				"resultName": nil,
			}
			inputList := []string{}
			for i, param := range method.Type.Params.List {
				inputList = append(inputList, getInputAsNameAndType(param, "param"+strconv.Itoa(i)))
			}
			if method.Type.Results != nil && len(method.Type.Results.List) > 0 {
				templInputs["resultName"] = getInputAsNameAndType(method.Type.Results.List[0], "result")
			}
			templInputs["inputs"] = strings.Join(inputList, ", ")
			if templInputs["resultName"] != nil {
				assignmentTmplInputs = append(assignmentTmplInputs, templInputs)
			} else {
				statementTemplInputs = append(statementTemplInputs, templInputs)
			}
		}
	}

	sort.Slice(statementTemplInputs, func(i, j int) bool {
		return statementTemplInputs[i]["prefix"].(string)+"."+statementTemplInputs[i]["method"].(string) < statementTemplInputs[j]["prefix"].(string)+"."+statementTemplInputs[j]["method"].(string)
	})
	sort.Slice(assignmentTmplInputs, func(i, j int) bool {
		return assignmentTmplInputs[i]["prefix"].(string)+"."+assignmentTmplInputs[i]["method"].(string) < assignmentTmplInputs[j]["prefix"].(string)+"."+assignmentTmplInputs[j]["method"].(string)
	})

	statementTempl, err := template.New("").Option("missingkey=zero").Parse(statementTmplStr)
	if err != nil {
		return "", err
	}
	assignmentTempl, err := template.New("").Option("missingkey=zero").Parse(assignmentTmplStr)
	if err != nil {
		return "", err
	}
	completerTempl, err := template.New("").Option("missingkey=zero").Parse(completerTmplStr)
	if err != nil {
		return "", err
	}

	newLineStatements := []string{}
	statementInputs := []string{}

	for _, input := range assignmentTmplInputs {
		builder := strings.Builder{}
		err = assignmentTempl.Execute(&builder, input)
		if err != nil {
			return "", err
		}
		newLineStatements = append(newLineStatements, builder.String())
	}

	for _, input := range append(statementTemplInputs, assignmentTmplInputs...) {
		builder := strings.Builder{}
		err = statementTempl.Execute(&builder, input)
		if err != nil {
			return "", err
		}
		statementInputs = append(statementInputs, builder.String())
	}

	completerTemplInputs := map[string]string{
		"newLineStatements": strings.Join(newLineStatements, ",\n"),
		"statements":        strings.Join(statementInputs, ",\n"),
	}

	builder := strings.Builder{}
	err = completerTempl.Execute(&builder, completerTemplInputs)
	return builder.String(), err
}

func getInputAsNameAndType(param *ast.Field, defaultName string) string {
	name := defaultName
	if len(param.Names) > 0 {
		name = param.Names[0].Name
	}
	switch t := param.Type.(type) {
	case *ast.Ident:
		name = name + "_as_" + t.Name
	case *ast.SelectorExpr:
		name = name + "_as_" + t.Sel.Name
	case *ast.InterfaceType:
		name = name + "_as_any"
	case *ast.ArrayType:
		switch element := t.Elt.(type) {
		case *ast.Ident:
			name = name + "_as_" + element.Name + "_list"
		case *ast.SelectorExpr:
			name = name + "_as_" + element.Sel.Name + "_list"
		case *ast.InterfaceType:
			name = name + "_as_list"
		}
	}
	return name
}

func filterMethods(dirAst map[string]*ast.Package, typeName string) (result []*ast.FuncDecl) {
	for _, packageAst := range dirAst {
		for _, fileAst := range packageAst.Files {
			for _, decl := range fileAst.Decls {
				fdecl, ok := decl.(*ast.FuncDecl)
				if ok && fdecl.Recv != nil && len(fdecl.Recv.List) > 0 {
					ptr, isStarExp := fdecl.Recv.List[0].Type.(*ast.StarExpr)
					if isStarExp {
						receiverType, isIdent := ptr.X.(*ast.Ident)
						if isIdent && receiverType.Name == typeName {
							result = append(result, fdecl)
						}
					}
				}
			}
		}
	}
	return result
}

func getScriptEnvMapping() map[string]string {
	result := map[string]string{}
	for key, obj := range (&scriptenv.ScriptEnv{}).GetEnvironment() {
		name := reflect.TypeOf(obj).Name()
		if name == "" && reflect.TypeOf(obj).Kind() == reflect.Pointer {
			name = reflect.TypeOf(obj).Elem().Name()
		}
		result[key] = name
	}
	return result
}
