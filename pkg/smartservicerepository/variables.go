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

package smartservicerepository

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"runtime/debug"
)

func (this *SmartServiceRepository) GetVariables(processId string) (result map[string]interface{}, err error) {
	req, err := http.NewRequest("GET", this.config.SmartServiceRepositoryUrl+"/instances-by-process-id/"+url.PathEscape(processId)+"/variables-map", nil)
	if err != nil {
		return result, err
	}
	token, err := this.auth.Ensure()
	if err != nil {
		return result, err
	}
	req.Header.Set("Authorization", token.Jwt())
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		temp, _ := io.ReadAll(resp.Body)
		err = errors.New(string(temp))
		return result, err
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (this *SmartServiceRepository) SetVariables(processId string, variableChanges map[string]interface{}) (err error) {
	body := new(bytes.Buffer)
	err = json.NewEncoder(body).Encode(variableChanges)
	if err != nil {
		this.config.GetLogger().Error("error in SmartServiceRepository.SetVariables", "error", err, "stack", string(debug.Stack()))
		return err
	}
	req, err := http.NewRequest("PUT", this.config.SmartServiceRepositoryUrl+"/instances-by-process-id/"+url.PathEscape(processId)+"/variables-map", body)
	if err != nil {
		return err
	}
	token, err := this.auth.Ensure()
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", token.Jwt())
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		temp, _ := io.ReadAll(resp.Body)
		err = errors.New(string(temp))
		return err
	}
	return nil
}
