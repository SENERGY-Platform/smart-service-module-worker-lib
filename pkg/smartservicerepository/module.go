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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime/debug"

	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/model"
)

func (this *SmartServiceRepository) SendWorkerModules(modules []model.Module) (result []model.SmartServiceModule, err error) {
	for _, module := range modules {
		temp, err := this.SendWorkerModule(module)
		if err != nil {
			return result, err
		}
		result = append(result, temp)
	}
	return result, nil
}

func (this *SmartServiceRepository) SendWorkerModule(module model.Module) (result model.SmartServiceModule, err error) {
	body := new(bytes.Buffer)
	err = json.NewEncoder(body).Encode(module.SmartServiceModuleInit)
	if err != nil {
		this.config.GetLogger().Error("error in SmartServiceRepository.SendWorkerModule", "error", err, "stack", string(debug.Stack()))
		return result, err
	}
	req, err := http.NewRequest("PUT", this.config.SmartServiceRepositoryUrl+"/instances-by-process-id/"+url.PathEscape(module.ProcesInstanceId)+"/modules/"+url.PathEscape(module.Id), body)
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

func (this *SmartServiceRepository) UseModuleDeleteInfo(info model.ModuleDeleteInfo) error {
	req, err := http.NewRequest("DELETE", info.Url, nil)
	if err != nil {
		return err
	}
	if info.UserId != "" {
		token, err := this.auth.ExchangeUserToken(info.UserId)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", token.Jwt())
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 && resp.StatusCode != http.StatusNotFound {
		temp, _ := io.ReadAll(resp.Body)
		err = fmt.Errorf("unexpected response: %v, %v", resp.StatusCode, string(temp))
		this.config.GetLogger().Error("error in SmartServiceRepository.UseModuleDeleteInfo", "error", err, "stack", string(debug.Stack()))
		return err
	}
	_, _ = io.ReadAll(resp.Body)
	return nil
}

func (this *SmartServiceRepository) ListExistingModules(processInstanceId string, query model.ModulQuery) (result []model.SmartServiceModule, err error) {
	queryValues := url.Values{}
	if query.KeyFilter != nil {
		queryValues.Set("key", *query.KeyFilter)
	}
	if query.TypeFilter != nil {
		queryValues.Set("module_type", *query.TypeFilter)
	}
	queryStr := ""
	if len(queryValues) > 0 {
		queryStr = "?" + queryValues.Encode()
	}

	req, err := http.NewRequest("GET", this.config.SmartServiceRepositoryUrl+"/instances-by-process-id/"+url.PathEscape(processInstanceId)+"/modules"+queryStr, nil)
	if err != nil {
		return result, err
	}
	token, err := this.auth.Ensure()
	if err != nil {
		return result, err
	}
	req.Header.Set("Authorization", token.Jwt())
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

func (this *SmartServiceRepository) GetModule(userId string, moduleId string) (result model.SmartServiceModule, err error, code int) {
	req, err := http.NewRequest("GET", this.config.SmartServiceRepositoryUrl+"/modules/"+url.PathEscape(moduleId), nil)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	token, err := this.auth.ExchangeUserToken(userId)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", token.Jwt())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, err, http.StatusInternalServerError
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		temp, _ := io.ReadAll(resp.Body)
		err = errors.New(string(temp))
		return result, err, resp.StatusCode
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, err, resp.StatusCode
	}
	return result, nil, resp.StatusCode
}
