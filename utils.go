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
	"strconv"
)

type UtilAPI interface {
	ConvertStrPortsToUint16(portsStr []string) ([]uint16, error)
	GetEth0IP() (string, error)
}

type UtilTools struct {
	sysUtil SystemUtilAPI
}

const (
	Eth0IPCommand = "echo $(ip -o -4 a s eth0 | awk '{ print $4 }' | cut -d/ -f1)"
)

func NewUtilTools(systemUtils SystemUtilAPI) UtilAPI {
	return &UtilTools{sysUtil: systemUtils}
}

func (t *UtilTools) ConvertStrPortsToUint16(portsStr []string) ([]uint16, error) {
	uint16Ports := []uint16{}

	for _, strPort := range portsStr {
		uintPort, err := strconv.ParseUint(strPort, 10, 16)
		if err != nil {
			return uint16Ports, err
		}

		uint16Ports = append(uint16Ports, uint16(uintPort))
	}
	return uint16Ports, nil
}

func (t *UtilTools) GetEth0IP() (string, error) {
	logger.Debug("fetching eth0 IP address")
	ip, err := t.sysUtil.ExecuteCommand(Eth0IPCommand)
	logger.Debugf("fetching ended with IP: %v", ip)
	return ip, err
}
