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
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/middleware/gen/util"
	"sort"
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
	methods := util.GetScriptEnvMethodTemplateInfos(pathToScriptenv)

	statementTemplInputs := []util.Info{}
	assignmentTmplInputs := []util.Info{}

	for _, info := range methods {
		if info.Result != nil {
			assignmentTmplInputs = append(assignmentTmplInputs, info)
		} else {
			statementTemplInputs = append(statementTemplInputs, info)
		}
	}

	sort.Slice(statementTemplInputs, func(i, j int) bool {
		return statementTemplInputs[i].Prefix+"."+statementTemplInputs[i].Method < statementTemplInputs[j].Prefix+"."+statementTemplInputs[j].Method
	})
	sort.Slice(assignmentTmplInputs, func(i, j int) bool {
		return assignmentTmplInputs[i].Prefix+"."+assignmentTmplInputs[i].Method < assignmentTmplInputs[j].Prefix+"."+assignmentTmplInputs[j].Method
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
