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

package camunda

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/configuration"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/model"
	"io"
	"net/http"
)

func SendEventTrigger(config configuration.Config, eventId string, variables map[string]model.CamundaVariable) (err error) {
	if variables == nil {
		variables = map[string]model.CamundaVariable{}
	}
	request, err := json.Marshal(map[string]interface{}{
		"messageName":           eventId,
		"all":                   true,
		"resultEnabled":         false,
		"processVariablesLocal": variables,
	})
	if err != nil {
		return err
	}
	resp, err := http.Post(config.CamundaUrl+"/engine-rest/message", "application/json", bytes.NewBuffer(request))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		err = errors.New(string(response))
	}
	return err
}
