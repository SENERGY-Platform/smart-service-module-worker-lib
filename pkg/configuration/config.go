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

package configuration

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	struct_logger "github.com/SENERGY-Platform/go-service-base/struct-logger"
	"github.com/SENERGY-Platform/go-service-base/struct-logger/handlers"
)

// DefaultServiceName is used as open-telemetry service name if the name of the worker
// can not be read from the build info.
const DefaultServiceName = "smart-service-module-worker"

type Config struct {
	DeviceRepositoryUrl                  string `json:"device_repository_url"`
	SmartServiceRepositoryUrl            string `json:"smart_service_repository_url"`
	CamundaUrl                           string `json:"camunda_url" config:"secret"`
	CamundaWorkerId                      string `json:"camunda_worker_id"`
	CamundaWorkerTopic                   string `json:"camunda_worker_topic"`
	CamundaLockDurationInMs              int64  `json:"camunda_lock_duration_in_ms"`
	CamundaWorkerWaitDurationInMs        int64  `json:"camunda_worker_wait_duration_in_ms"`
	CamundaFetchMaxTasks                 int64  `json:"camunda_fetch_max_tasks"`
	AuthEndpoint                         string `json:"auth_endpoint"`
	AuthClientId                         string `json:"auth_client_id" config:"secret"`
	AuthClientSecret                     string `json:"auth_client_secret" config:"secret"`
	TokenCacheDefaultExpirationInSeconds int    `json:"token_cache_default_expiration_in_seconds"`

	FetchRetries int `json:"fetch_retries"`

	LogLevel     string       `json:"log_level"`
	OtelEndpoint string       `json:"otel_endpoint"`
	logger       *slog.Logger `json:"-"`
}

func LoadLibConfig(location string) (config Config, err error) {
	return Load[Config](location)
}

// loads config from json in location and used environment variables (e.g KafkaUrl --> KAFKA_URL)
func Load[T any](location string) (config T, err error) {
	file, err := os.Open(location)
	if err != nil {
		return config, err
	}
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return config, err
	}
	handleEnvironmentVars(&config)
	return config, nil
}

var camel = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")

func fieldNameToEnvName(s string) string {
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}
	return strings.ToUpper(strings.Join(a, "_"))
}

// preparations for docker
func handleEnvironmentVars[T any](config *T) {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	configType := configValue.Type()
	for index := 0; index < configType.NumField(); index++ {
		fieldName := configType.Field(index).Name
		fieldConfig := configType.Field(index).Tag.Get("config")
		envName := fieldNameToEnvName(fieldName)
		envValue := os.Getenv(envName)
		if envValue != "" {
			if !strings.Contains(fieldConfig, "secret") {
				fmt.Println("use environment variable: ", envName, " = ", envValue)
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.Int64 || configValue.FieldByName(fieldName).Kind() == reflect.Int {
				i, _ := strconv.ParseInt(envValue, 10, 64)
				configValue.FieldByName(fieldName).SetInt(i)
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.String {
				configValue.FieldByName(fieldName).SetString(envValue)
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.Bool {
				b, _ := strconv.ParseBool(envValue)
				configValue.FieldByName(fieldName).SetBool(b)
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.Float64 {
				f, _ := strconv.ParseFloat(envValue, 64)
				configValue.FieldByName(fieldName).SetFloat(f)
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.Slice {
				val := []string{}
				for _, element := range strings.Split(envValue, ",") {
					val = append(val, strings.TrimSpace(element))
				}
				configValue.FieldByName(fieldName).Set(reflect.ValueOf(val))
			}
			if configValue.FieldByName(fieldName).Kind() == reflect.Map {
				value := map[string]string{}
				for _, element := range strings.Split(envValue, ",") {
					keyVal := strings.Split(element, ":")
					key := strings.TrimSpace(keyVal[0])
					val := strings.TrimSpace(keyVal[1])
					value[key] = val
				}
				configValue.FieldByName(fieldName).Set(reflect.ValueOf(value))
			}
		}
	}
}

func (this *Config) GetLogger() *slog.Logger {
	if this.logger == nil {
		info, ok := debug.ReadBuildInfo()
		project := ""
		org := ""
		if ok {
			if parts := strings.Split(info.Main.Path, "/"); len(parts) > 2 {
				project = strings.Join(parts[2:], "/")
				org = strings.Join(parts[:2], "/")
			}
		}
		base := struct_logger.New(
			struct_logger.Config{
				Handler:    struct_logger.JsonHandlerSelector,
				Level:      this.LogLevel,
				TimeFormat: time.RFC3339Nano,
				TimeUtc:    true,
				AddMeta:    true,
			},
			os.Stdout,
			org,
			project,
		)
		//the open-telemetry handler adds the baggage of the context to the log record and the record to the current span.
		//it only takes effect if the context based log methods are used (ErrorContext(), InfoContext(), ...)
		this.logger = slog.New(handlers.NewOpenTelemetryHandler(base.Handler())).With("project-group", "smart-service")
	}
	return this.logger
}

// ServiceName returns the name this process is reported under in open-telemetry.
// This library is shared by multiple worker services; every one of them has to show up in
// jaeger under its own name, so the name is taken from the main module of the binary
// (e.g. github.com/SENERGY-Platform/smart-service-module-worker-analytics -->
// smart-service-module-worker-analytics) and not from a constant of this library.
func ServiceName() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		if name := serviceNameFromModulePath(info.Main.Path); name != "" {
			return name
		}
	}
	return DefaultServiceName
}

// serviceNameFromModulePath returns the last segment of a go module path.
// a major-version suffix (/v2, /v3, ...) is not part of the name of the service.
// returns "" if no usable name can be found, the caller falls back to DefaultServiceName.
func serviceNameFromModulePath(modulePath string) string {
	parts := strings.Split(strings.Trim(modulePath, "/"), "/")
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		if part == "" {
			continue
		}
		if len(part) > 1 && part[0] == 'v' {
			if _, err := strconv.Atoi(part[1:]); err == nil {
				continue //major-version suffix, the name is in front of it
			}
		}
		return part
	}
	return ""
}
