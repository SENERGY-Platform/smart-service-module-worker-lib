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

package references

import (
	"encoding/json"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/model"
	"strings"
	"text/template"
)

func Handle(inputs map[string]model.CamundaVariable, variables map[string]interface{}) (result map[string]model.CamundaVariable, err error) {
	result = map[string]model.CamundaVariable{}
	placeholder, err := GetPlaceholder(variables)
	if err != nil {
		return result, err
	}
	for key, input := range inputs {
		if str, ok := input.Value.(string); ok {
			input.Value, err = Replace(str, placeholder)
			if err != nil {
				return result, err
			}
		}
		result[key] = input
	}
	return result, nil
}

func GetPlaceholder(variables map[string]interface{}) (result map[string]string, err error) {
	result = map[string]string{
		"brl": "{{",
		"brr": "}}",
	}
	for key, value := range variables {
		switch v := value.(type) {
		case string:
			//normal version
			result[key] = v
			//json version
			/*
				temp, err := json.Marshal(v)
				if err != nil {
					return result, err
				}
				result[key+"_json"] = string(temp)
			*/
		default:
			temp, err := json.Marshal(v)
			if err != nil {
				return result, err
			}
			result[key] = string(temp)
		}

	}
	return result, nil
}

func Replace(str string, variables map[string]string) (string, error) {
	tmpl, err := template.New("").Option("missingkey=zero").Parse(str)
	if err != nil {
		return str, err
	}
	builder := strings.Builder{}
	err = tmpl.Execute(&builder, variables)
	return builder.String(), err
}
