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
	"context"
)

// GetInstanceUser returns the user of the smart-service-instance that belongs to the given
// camunda process instance.
// it delegates to GetCachedSmartServiceInstance instead of reading the user-id on its own:
// the middleware has already read the whole instance for this task under the same cache key,
// so the user-id comes out of the cache instead of causing a second request per task.
func (this *SmartServiceRepository) GetInstanceUser(ctx context.Context, processInstanceId string) (userId string, err error) {
	instance, err := this.GetCachedSmartServiceInstance(ctx, processInstanceId)
	if err != nil {
		return "", err
	}
	return instance.UserId, nil
}
