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

package configuration

import (
	"testing"
)

func TestServiceNameFromModulePath(t *testing.T) {
	tests := map[string]string{
		"github.com/SENERGY-Platform/smart-service-module-worker-analytics":  "smart-service-module-worker-analytics",
		"github.com/SENERGY-Platform/smart-service-module-worker-lib":        "smart-service-module-worker-lib",
		"github.com/SENERGY-Platform/smart-service-module-worker-process/v2": "smart-service-module-worker-process",
		"worker": "worker",
		"":       "",
		"/":      "",
	}
	for modulePath, expected := range tests {
		if actual := serviceNameFromModulePath(modulePath); actual != expected {
			t.Errorf("%q: got %q, want %q", modulePath, actual, expected)
		}
	}
}

// TestServiceName only checks that a usable name comes out; the concrete value depends on the
// binary the library is linked into and is the module path of this repository under go test.
func TestServiceName(t *testing.T) {
	name := ServiceName()
	if name == "" {
		t.Error("empty service name")
	}
	t.Log(name)
}

func TestOtelEndpointEnvName(t *testing.T) {
	if actual := fieldNameToEnvName("OtelEndpoint"); actual != "OTEL_ENDPOINT" {
		t.Error(actual)
	}
}
