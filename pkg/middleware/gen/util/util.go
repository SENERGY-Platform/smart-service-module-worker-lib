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

package util

import (
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/middleware/scriptenv"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

type Info struct {
	Prefix    string
	Method    string
	Inputs    []Parameter
	InputsStr string
	Result    *Parameter
	Comment   string
}

type Parameter struct {
	Name         string
	NameWithType string
	Type         string
}

func GetScriptEnvMethodTemplateInfos(pathToScriptenv string) (result []Info) {
	scriptEnvMapping := GetScriptEnvMapping()

	f, err := parser.ParseDir(token.NewFileSet(), pathToScriptenv, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	for key, typeName := range scriptEnvMapping {
		for _, method := range FilterMethods(f, typeName) {
			info := Info{
				Prefix:    key,
				Method:    Uncapitalize(method.Name.Name),
				Inputs:    nil,
				Result:    nil,
				InputsStr: "",
				Comment:   strings.TrimSpace(method.Doc.Text()),
			}
			inputList := []string{}
			for i, param := range method.Type.Params.List {
				input := GetInputParameter(param, "param"+strconv.Itoa(i))
				info.Inputs = append(info.Inputs, input)
				inputList = append(inputList, input.NameWithType)
			}
			if method.Type.Results != nil && len(method.Type.Results.List) > 0 {
				resultParam := GetInputParameter(method.Type.Results.List[0], "result")
				info.Result = &resultParam
			}
			info.InputsStr = strings.Join(inputList, ", ")
			result = append(result, info)
		}
	}
	slices.SortFunc(result, func(a, b Info) int {
		temp := strings.Compare(a.Prefix, b.Prefix)
		if temp != 0 {
			return temp
		}
		return strings.Compare(a.Method, b.Method)
	})
	return result
}

func GetScriptEnvMapping() map[string]string {
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

func FilterMethods(dirAst map[string]*ast.Package, typeName string) (result []*ast.FuncDecl) {
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

func Uncapitalize(s string) string {
	if s == "" {
		return ""
	}
	if len(s) == 1 {
		return strings.ToLower(s[0:1])
	}
	return strings.ToLower(s[0:1]) + s[1:]
}

func GetInputParameter(param *ast.Field, defaultName string) (result Parameter) {
	result.Name = GetInputAsName(param, defaultName)
	result.Type = GetInputAsJsDocType(param)
	result.NameWithType = result.Name
	if result.Type != "" {
		result.NameWithType = result.Name + "_as_" + GetInputAsTypeName(param)
	}
	return result
}

func GetInputAsTypeName(param *ast.Field) string {
	switch t := param.Type.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return t.Sel.Name
	case *ast.InterfaceType:
		return "any"
	case *ast.ArrayType:
		switch element := t.Elt.(type) {
		case *ast.Ident:
			return element.Name + "_list"
		case *ast.SelectorExpr:
			return element.Sel.Name + "_list"
		case *ast.InterfaceType:
			return "list"
		}
	}
	return ""
}

func GetInputAsJsDocType(param *ast.Field) string {
	switch t := param.Type.(type) {
	case *ast.Ident:
		return toJsDocName(t.Name)
	case *ast.SelectorExpr:
		return toJsDocName(t.Sel.Name)
	case *ast.InterfaceType:
		return "Object"
	case *ast.ArrayType:
		switch element := t.Elt.(type) {
		case *ast.Ident:
			return toJsDocName(element.Name) + "[]"
		case *ast.SelectorExpr:
			return toJsDocName(element.Sel.Name) + "[]"
		case *ast.InterfaceType:
			return "Object[]"
		}
	}
	return ""
}

func toJsDocName(name string) string {
	switch name {
	case "bool":
		return "boolean"
	default:
		return name
	}
}

func GetInputAsName(param *ast.Field, defaultName string) string {
	name := defaultName
	if len(param.Names) > 0 {
		name = param.Names[0].Name
	}
	return name
}
