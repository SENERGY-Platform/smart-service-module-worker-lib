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

package main

import (
	"fmt"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/middleware/gen/acecodecompleter"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/middleware/gen/jsdoc"
	"os"
	"path"
	"sync"
)

//go:generate go run gen.go

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	docDir := "../../../doc"

	go func() {
		defer wg.Done()
		aceCodeCompleter, err := acecodecompleter.GenerateTsAceCodeCompleter("../scriptenv")
		if err != nil {
			panic(err)
		}
		aceCompleterFile := path.Join(docDir, "ace-code-completer.ts")
		storeAceCompleterFile(aceCompleterFile, aceCodeCompleter)
	}()

	go func() {
		defer wg.Done()
		jsDocOutput, err := jsdoc.GenerateJsDoc("../scriptenv")
		if err != nil {
			panic(err)
		}
		jsDocFile := path.Join(docDir, "jsdoc.js")
		storeAceCompleterFile(jsDocFile, jsDocOutput)
	}()

	wg.Wait()
}

func storeAceCompleterFile(location string, output string) {
	file, err := os.OpenFile(location, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			panic(fmt.Errorf("unable to close output file %v %w", location, err))
		}
	}()
	_, err = file.WriteString(output)
	if err != nil {
		panic(fmt.Errorf("unable to open write to output file %v %w", location, err))
	}
}
