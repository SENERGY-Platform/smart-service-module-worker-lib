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
	this.config.GetLogger().Info("run health check")
	moduleCount := 0
	checked := 0
	skipped := 0
	healthy := 0
	ill := 0
	updatedAsHealthy := 0
	updatedAsIll := 0
	defer func() {
		this.config.GetLogger().Info("finished health check", "modules", moduleCount, "checked", checked, "skipped", skipped, "healthy", healthy, "ill", ill, "updatedAsHealthy", updatedAsHealthy, "updatedAsIll", updatedAsIll)
	}()
	for module := range util.IterBatch(100, func(limit int64, offset int64) ([]model.SmartServiceModule, error) {
		query.Limit = limit
		query.Offset = offset
		return this.ListModules(query)
	}) {
		moduleCount++
		if module.LastUpdate > 0 && time.Since(time.Unix(module.LastUpdate, 0)) < time.Hour {
			//ignore modules that were updated in the last hour
			skipped++
			continue
		}
		checked++
		health, err := check(module)
		if err != nil {
			this.config.GetLogger().Error("error in health check", "error", err, "module", module)
			continue
		}
		if health == nil {
			healthy++
		} else {
			ill++
		}
		if health != nil {
			updatedAsIll++
			err = this.SetSmartServiceModuleError(module.Id, health)
			if err != nil {
				this.config.GetLogger().Error("error in health check", "error", err, "module", module)
				continue
			}
		}
		if health == nil && module.Error != "" {
			updatedAsHealthy++
			err = this.RemoveSmartServiceModuleError(module.Id)
			if err != nil {
				this.config.GetLogger().Error("error in health check", "error", err, "module", module)
			}
		}
	}
}
