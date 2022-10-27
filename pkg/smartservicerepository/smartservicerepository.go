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
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/auth"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/configuration"
)

type SmartServiceRepository struct {
	config configuration.Config
	auth   Auth
}

type Auth interface {
	Ensure() (token auth.Token, err error)
	ExchangeUserToken(userid string) (token auth.Token, err error)
}

func New(config configuration.Config, auth Auth) *SmartServiceRepository {
	return &SmartServiceRepository{config: config, auth: auth}
}
