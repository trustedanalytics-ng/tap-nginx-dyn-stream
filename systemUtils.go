/**
 * Copyright (c) 2016 Intel Corporation
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
	"os/exec"
	"strings"
)

type SystemUtilAPI interface {
	ExecuteCommand(cmd string) (string, error)
}

type SystemUtils struct{}

func (sU *SystemUtils) ExecuteCommand(cmd string) (string, error) {
	logger.Debugf("executing command: %v with string output", cmd)
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	outputWithoutNewline := strings.TrimRight(string(out), "\n")
	logger.Debugf("executing: %v ended, output: %v", cmd, outputWithoutNewline)
	return outputWithoutNewline, err
}
