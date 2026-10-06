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

package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"reflect"
	"strings"
	"testing"

	"github.com/SENERGY-Platform/go-service-base/struct-logger/handlers"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/middleware/scriptenv"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/tracing"
)

// runConsoleScript runs script in a script-env whose logger is built like configuration.Config.GetLogger(),
// with the open-telemetry handler in front, and returns the written log records
func runConsoleScript(t *testing.T, script string, variables map[string]interface{}, inputs map[string]interface{}) []map[string]interface{} {
	t.Helper()
	ctx, err := tracing.AddToBaggage(context.Background(), tracing.BaggageKeyInstanceId, "test-instance-id")
	if err != nil {
		t.Fatal(err)
	}
	buf := &bytes.Buffer{}
	logger := slog.New(handlers.NewOpenTelemetryHandler(slog.NewJSONHandler(buf, nil))).With("script", PreScriptPrefix)
	env := scriptenv.NewScriptEnv(ctx, logger, AuthMock, nil, "user", variables, inputs, nil)
	err = runScript(script, env)
	if err != nil {
		t.Fatal(err)
	}
	result := []map[string]interface{}{}
	for _, line := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		if line == "" {
			continue
		}
		record := map[string]interface{}{}
		err = json.Unmarshal([]byte(line), &record)
		if err != nil {
			t.Fatal(err, line)
		}
		result = append(result, record)
	}
	return result
}

func TestConsoleLog(t *testing.T) {
	records := runConsoleScript(t, `console.log("foo", 42, true, null, {"bar": [1, "batz"]});`, map[string]interface{}{}, map[string]interface{}{})
	if len(records) != 1 {
		t.Fatalf("expected 1 log record, got %#v", records)
	}
	record := records[0]
	if record["msg"] != `foo 42 true null {"bar":[1,"batz"]}` {
		t.Errorf("unexpected msg %#v", record["msg"])
	}
	if record["level"] != "INFO" {
		t.Errorf("unexpected level %#v", record["level"])
	}
	if record["script"] != PreScriptPrefix {
		t.Errorf("unexpected script %#v", record["script"])
	}
	if record[tracing.BaggageKeyInstanceId] != "test-instance-id" {
		t.Errorf("missing baggage in log record %#v", record)
	}
}

func TestConsoleDump(t *testing.T) {
	variables := map[string]interface{}{"v1": "str", "v2": float64(42)}
	inputs := map[string]interface{}{
		"in1":                   "foo",
		PreScriptPrefix + "_1":  "console.dump();",
		PostScriptPrefix + "_1": "console.log('post');",
	}
	records := runConsoleScript(t, `variables.write("v3", {"foo": "bar"}); outputs.set("out1", 13); console.dump();`, variables, inputs)
	if len(records) != 1 {
		t.Fatalf("expected 1 log record, got %#v", records)
	}
	record := records[0]
	if record[tracing.BaggageKeyInstanceId] != "test-instance-id" {
		t.Errorf("missing baggage in log record %#v", record)
	}
	//changes made earlier in the same script are part of the dump
	expectedVariables := map[string]interface{}{"v1": "str", "v2": float64(42), "v3": map[string]interface{}{"foo": "bar"}}
	if !reflect.DeepEqual(record["variables"], expectedVariables) {
		t.Errorf("unexpected variables %#v", record["variables"])
	}
	//only the changes of this script, although they are part of the variables as well
	if !reflect.DeepEqual(record["variablesUpdates"], map[string]interface{}{"v3": map[string]interface{}{"foo": "bar"}}) {
		t.Errorf("unexpected variablesUpdates %#v", record["variablesUpdates"])
	}
	//script sources are no inputs of the script-env
	if !reflect.DeepEqual(record["inputs"], map[string]interface{}{"in1": "foo"}) {
		t.Errorf("unexpected inputs %#v", record["inputs"])
	}
	if !reflect.DeepEqual(record["outputs"], map[string]interface{}{"out1": float64(13)}) {
		t.Errorf("unexpected outputs %#v", record["outputs"])
	}
}
