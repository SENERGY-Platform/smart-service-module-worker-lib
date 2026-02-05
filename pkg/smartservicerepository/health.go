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

package smartservicerepository

import (
	"context"
	"time"

	"github.com/SENERGY-Platform/service-commons/pkg/util"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/model"
)

func (this *SmartServiceRepository) StartHealthCheck(ctx context.Context, interval time.Duration, query model.ModulQuery, check func(module model.SmartServiceModule) (health error, err error)) {
	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				this.RunHealthCheck(query, check)
			}
		}
	}()
}

func (this *SmartServiceRepository) RunHealthCheck(query model.ModulQuery, check func(module model.SmartServiceModule) (health error, err error)) {
	for module := range util.IterBatch(100, func(limit int64, offset int64) ([]model.SmartServiceModule, error) {
		query.Limit = limit
		query.Offset = offset
		return this.ListModules(query)
	}) {
		health, err := check(module)
		if err != nil {
			this.config.GetLogger().Error("error in health check", "error", err, "module", module)
			continue
		}
		if health != nil {
			err = this.SetSmartServiceError(module.InstanceId, health)
			if err != nil {
				this.config.GetLogger().Error("error in health check", "error", err, "module", module)
				continue
			}
		}
	}
}
