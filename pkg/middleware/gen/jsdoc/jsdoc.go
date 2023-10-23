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
	_ "embed"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/middleware/gen/util"
	"slices"
	"sort"
	"strings"
	"text/template"
)

//go:embed methods.tmpl
var methodsTmplStr string

//go:embed namespaces.tmpl
var namespacesTmplStr string

//go:embed typedefs.tmpl
var typedefTmplStr string

//go:embed jsdoc.tmpl
var jsdocTemplStr string

func GenerateJsDoc(pathToScriptenv string) (string, error) {
	methods := util.GetScriptEnvMethodTemplateInfos(pathToScriptenv)

	namespaces := []string{}
	for _, method := range methods {
		if !slices.Contains(namespaces, method.Prefix) {
			namespaces = append(namespaces, method.Prefix)
		}
	}

	sort.Strings(namespaces)

	namespacesTmpl, err := template.New("namespacesTmpl").Option("missingkey=zero").Parse(namespacesTmplStr)
	if err != nil {
		return "", err
	}

	namespacesBuilder := strings.Builder{}
	err = namespacesTmpl.Execute(&namespacesBuilder, namespaces)
	if err != nil {
		return "", err
	}

	methodsTmpl, err := template.New("methodsTmpl").Option("missingkey=zero").Parse(methodsTmplStr)
	if err != nil {
		return "", err
	}

	methodsBuilder := strings.Builder{}
	err = methodsTmpl.Execute(&methodsBuilder, methods)
	if err != nil {
		return "", err
	}

	typedefTmpl, err := template.New("typedefTmpl").Option("missingkey=zero").Parse(typedefTmplStr)
	if err != nil {
		return "", err
	}

	typedefBuilder := strings.Builder{}
	err = typedefTmpl.Execute(&methodsBuilder, GetTypeDefs())
	if err != nil {
		return "", err
	}

	jsdocTemplInputs := map[string]string{
		"namespaces": namespacesBuilder.String(),
		"methods":    methodsBuilder.String(),
		"typedefs":   typedefBuilder.String(),
	}

	jsdocTempl, err := template.New("jsdocTempl").Option("missingkey=zero").Parse(jsdocTemplStr)
	if err != nil {
		return "", err
	}

	jsdocBuilder := strings.Builder{}
	err = jsdocTempl.Execute(&jsdocBuilder, jsdocTemplInputs)

	return jsdocBuilder.String(), err
}
