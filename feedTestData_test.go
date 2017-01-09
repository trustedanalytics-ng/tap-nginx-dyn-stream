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
	"strings"
)

const (
	decimalBase      = 10
	uInt16BitSize    = 16
	simpleErrMessage = "bad streamPortsEnvValue"
)

type ErrParseBadValue struct {
	message       string
	listenPort    string
	proxyPassPort string
}

func (e ErrParseBadValue) Error() string {
	return e.message
}

type ServerPortsPair struct {
	listenPort        string
	proxyPassPort     string
	listenPortUint    uint16
	proxyPassPortUint uint16
}

type FeedTestData struct {
	streamPortsEnvValue string
	servers             []ServerPortsPair
}

func NewFeedTestData(envValue string) (FeedTestData, error) {
	feedData := FeedTestData{
		streamPortsEnvValue: envValue,
		servers:             []ServerPortsPair{},
	}

	portPairs := strings.Split(envValue, ",")

	for _, portPair := range portPairs {
		portsSplitted := strings.Split(portPair, ":")
		if len(portsSplitted) != 2 {
			return feedData, ErrParseBadValue{message: simpleErrMessage}
		}

		listenPortUint, err := strconv.ParseUint(portsSplitted[0], decimalBase, uInt16BitSize)
		if err != nil {
			return feedData, ErrParseBadValue{message: simpleErrMessage, listenPort: portsSplitted[0], proxyPassPort: portsSplitted[1]}
		}

		proxyPassPortUint, err := strconv.ParseUint(portsSplitted[1], decimalBase, uInt16BitSize)
		if err != nil {
			return feedData, ErrParseBadValue{message: simpleErrMessage, listenPort: portsSplitted[0], proxyPassPort: portsSplitted[1]}
		}

		parsedServer := ServerPortsPair{
			listenPort:        portsSplitted[0],
			proxyPassPort:     portsSplitted[1],
			listenPortUint:    uint16(listenPortUint),
			proxyPassPortUint: uint16(proxyPassPortUint),
		}

		feedData.servers = append(feedData.servers, parsedServer)
	}

	return feedData, nil
}
