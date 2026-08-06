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
	_ "embed"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/middleware/gen/jsdoc"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/middleware/gen/util"
	"slices"
	"strings"
	"text/template"
)

//go:embed declarations.tmpl
var declarationsTmplStr string

type Interface struct {
	Name       string
	Properties []Property
}

type Property struct {
	Name string
	Type string
}

type Namespace struct {
	Name    string
	Methods []Method
}

type Method struct {
	Name    string
	Params  string
	Result  string
	Comment []string
}

// GenerateTypescriptDeclarations creates the typescript declarations of the script
// environment, so that editors can offer typed completion for pre/postscripts.
func GenerateTypescriptDeclarations(pathToScriptenv string) (string, error) {
	tmpl, err := template.New("declarationsTmpl").Option("missingkey=zero").Parse(declarationsTmplStr)
	if err != nil {
		return "", err
	}
	builder := strings.Builder{}
	err = tmpl.Execute(&builder, map[string]interface{}{
		"interfaces": GetInterfaces(),
		"namespaces": GetNamespaces(pathToScriptenv),
	})
	return builder.String(), err
}

func GetInterfaces() (result []Interface) {
	for _, def := range jsdoc.GetTypeDefs() {
		element := Interface{Name: def.Name}
		for _, field := range def.Fields {
			element.Properties = append(element.Properties, Property{
				Name: field.Name,
				Type: ToTsType(field.Type),
			})
		}
		result = append(result, element)
	}
	return result
}

func GetNamespaces(pathToScriptenv string) (result []Namespace) {
	//util.GetScriptEnvMethodTemplateInfos() sorts by prefix, so each namespace is
	//encountered as one uninterrupted block
	for _, info := range util.GetScriptEnvMethodTemplateInfos(pathToScriptenv) {
		i := slices.IndexFunc(result, func(n Namespace) bool { return n.Name == info.Prefix })
		if i < 0 {
			result = append(result, Namespace{Name: info.Prefix})
			i = len(result) - 1
		}
		result[i].Methods = append(result[i].Methods, getMethod(info))
	}
	return result
}

func getMethod(info util.Info) (result Method) {
	result.Name = info.Method
	result.Result = "void"
	if info.Result != nil {
		result.Result = ToTsType(info.Result.Type)
	}
	params := []string{}
	for _, param := range info.Inputs {
		params = append(params, param.Name+": "+ToTsType(param.Type))
	}
	result.Params = strings.Join(params, ", ")
	if info.Comment != "" {
		result.Comment = strings.Split(info.Comment, "\n")
	}
	return result
}

// ToTsType translates a jsdoc type as produced by the jsdoc generator into its
// typescript equivalent. Both describe the same value, but jsdoc spells some of them
// differently: a go map reaches scripts as a plain object rather than as a js Map,
// and jsdoc has no "any".
func ToTsType(jsDocType string) string {
	switch {
	case jsDocType == "":
		return "void"
	case strings.HasSuffix(jsDocType, "[]"):
		return ToTsType(strings.TrimSuffix(jsDocType, "[]")) + "[]"
	case strings.HasSuffix(jsDocType, "|null"):
		element := ToTsType(strings.TrimSuffix(jsDocType, "|null"))
		if element == "any" {
			return element //"any" already includes null
		}
		return element + " | null"
	case strings.HasPrefix(jsDocType, "Map<string,") && strings.HasSuffix(jsDocType, ">"):
		element := strings.TrimSuffix(strings.TrimPrefix(jsDocType, "Map<string,"), ">")
		return "Record<string, " + ToTsType(element) + ">"
	case jsDocType == "Object":
		return "any"
	}
	return jsDocType
}
