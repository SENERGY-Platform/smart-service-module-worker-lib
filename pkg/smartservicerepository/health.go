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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// TracerName identifies this package as the source of the spans it creates.
const TracerName = "github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/smartservicerepository"

func (this *SmartServiceRepository) StartHealthCheck(ctx context.Context, interval time.Duration, query model.ModulQuery, check func(ctx context.Context, module model.SmartServiceModule) (health error, err error)) {
	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				this.RunHealthCheck(ctx, query, check)
			}
		}
	}()
}

func (this *SmartServiceRepository) RunHealthCheck(ctx context.Context, query model.ModulQuery, check func(ctx context.Context, module model.SmartServiceModule) (health error, err error)) {
	//one span per run, not one per checked module: the run is what is triggered and what can be
	//slow, while a span per module would scale with the module count without answering a question
	//anyone asks. the ctx that carries this span is what the module listing, the check callback and
	//the error updates get, so a run is one connected trace instead of unrelated single requests.
	//the ctx of StartHealthCheck is the worker lifetime ctx and carries no span at all; handing
	//that one to the callback is what the workers already do and is exactly what looks connected
	//without being connected.
	ctx, span := otel.Tracer(TracerName).Start(ctx, "health-check")
	defer span.End()
	this.config.GetLogger().InfoContext(ctx, "run health check")
	moduleCount := 0
	checked := 0
	skipped := 0
	healthy := 0
	ill := 0
	updatedAsHealthy := 0
	updatedAsIll := 0
	defer func() {
		span.SetAttributes(
			attribute.Int("health_check_modules", moduleCount),
			attribute.Int("health_check_checked", checked),
			attribute.Int("health_check_skipped", skipped),
			attribute.Int("health_check_healthy", healthy),
			attribute.Int("health_check_ill", ill),
		)
		this.config.GetLogger().InfoContext(ctx, "finished health check", "modules", moduleCount, "checked", checked, "skipped", skipped, "healthy", healthy, "ill", ill, "updatedAsHealthy", updatedAsHealthy, "updatedAsIll", updatedAsIll)
	}()
	for module := range util.IterBatch(100, func(limit int64, offset int64) ([]model.SmartServiceModule, error) {
		query.Limit = limit
		query.Offset = offset
		return this.ListModules(ctx, query)
	}) {
		moduleCount++
		if module.LastUpdate > 0 && time.Since(time.Unix(module.LastUpdate, 0)) < time.Hour {
			//ignore modules that were updated in the last hour
			skipped++
			continue
		}
		checked++
		health, err := check(ctx, module)
		if err != nil {
			this.config.GetLogger().ErrorContext(ctx, "error in health check", "error", err, "module", module)
			continue
		}
		if health == nil {
			healthy++
		} else {
			ill++
		}
		if health != nil {
			updatedAsIll++
			err = this.SetSmartServiceModuleError(ctx, module.Id, health)
			if err != nil {
				this.config.GetLogger().ErrorContext(ctx, "error in health check", "error", err, "module", module)
				continue
			}
		}
		if health == nil && module.Error != "" {
			updatedAsHealthy++
			err = this.RemoveSmartServiceModuleError(ctx, module.Id)
			if err != nil {
				this.config.GetLogger().ErrorContext(ctx, "error in health check", "error", err, "module", module)
			}
		}
	}
}
