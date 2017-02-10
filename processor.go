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
	"errors"
	"strings"

	"github.com/trustedanalytics-ng/tap-go-common/util"
)

var tools UtilAPI = NewUtilTools(&SystemUtils{})

func makeFeed() (TemplateFeed, error) {
	feed := TemplateFeed{}
	var err error

	streamPortsString := util.GetEnvValueOrDefault(streamPortsEnvVarName, defaultStreamPorts)

	feed.StreamServers, err = makeStreamServerFeeds(streamPortsString)
	if err != nil {
		return feed, err
	}

	logLevelRaw := util.GetEnvValueOrDefault(logLevelEnvVarName, defaultLogLevel)
	feed.LogLevel = strings.ToLower(logLevelRaw)

	return feed, nil
}

func makeStreamServerFeeds(envVarString string) ([]StreamServer, error) {
	servers := []StreamServer{}
	ports := strings.Split(envVarString, ",")
	ip, err := tools.GetEth0IP()
	if err != nil {
		return servers, err
	}

	for _, port := range ports {
		var portPairStr []string
		if strings.Contains(port, ":") {
			portPairStr = strings.Split(port, ":")
			if len(portPairStr) != 2 {
				return servers, errors.New("bad ports format!")
			}
		} else {
			portPairStr = []string{port, port}
		}

		uint16Ports, err := tools.ConvertStrPortsToUint16(portPairStr)
		if err != nil {
			return servers, err
		}

		server := StreamServer{
			ListenPort:    uint16Ports[0],
			ListenIP:      ip,
			ProxyPassPort: uint16Ports[1],
		}

		servers = append(servers, server)
	}
	return servers, nil
}
