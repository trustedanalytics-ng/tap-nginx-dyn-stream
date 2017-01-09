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
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
)

func prepareMocksForProcessorUtils(t *testing.T) (*gomock.Controller, *MockUtilAPI) {
	mockCtrl := gomock.NewController(t)
	toolsMock := NewMockUtilAPI(mockCtrl)
	tools = toolsMock
	return mockCtrl, toolsMock
}

func prepareMakeFeedTestCase(t *testing.T, toolsMock *MockUtilAPI, fakeIP string, streamPortsVarValue string) TemplateFeed {
	feedData, err := NewFeedTestData(streamPortsVarValue)

	os.Setenv(streamPortsEnvVarName, streamPortsVarValue)

	toolsMock.EXPECT().GetEth0IP().Return(fakeIP, nil)
	for _, server := range feedData.servers {
		toolsMock.EXPECT().ConvertStrPortsToUint16([]string{server.listenPort, server.proxyPassPort}).Return([]uint16{server.listenPortUint, server.proxyPassPortUint}, nil)
	}

	if err != nil {
		negativeExpectErr, ok := err.(ErrParseBadValue)
		if !ok {
			t.Fatal("bad error type")
		}

		toolsMock.EXPECT().ConvertStrPortsToUint16([]string{negativeExpectErr.listenPort, negativeExpectErr.proxyPassPort}).Return([]uint16{}, negativeExpectErr)
	}

	expectedFeed := TemplateFeed{
		LogLevel:      defaultLogLevel,
		StreamServers: []StreamServer{},
	}

	for _, server := range feedData.servers {
		streamServer := StreamServer{
			ListenPort:    server.listenPortUint,
			ProxyPassPort: server.proxyPassPortUint,
			ListenIP:      fakeIP,
		}
		expectedFeed.StreamServers = append(expectedFeed.StreamServers, streamServer)
	}

	return expectedFeed
}

type MakeFeedTestCaseData struct {
	dataDesc            string
	streamPortsVarValue string
	negativeCase        bool
}

func TestMakeFeed(t *testing.T) {
	fakeIP := "10.0.0.12"

	cases := []MakeFeedTestCaseData{
		{
			dataDesc:            "no env vars",
			streamPortsVarValue: defaultStreamPorts,
		},
		{
			dataDesc:            "simple env var",
			streamPortsVarValue: "22:33",
		},
		{
			dataDesc:            "complicated env var",
			streamPortsVarValue: "55:65,88:143,1654:60000",
		},
		{
			dataDesc:            "bad env var",
			streamPortsVarValue: "55:65,aa:1vcv,1654:60000",
			negativeCase:        true,
		},
	}

	convey.Convey("Test makeFeed", t, func() {
		for _, test := range cases {
			ctrl, toolsMock := prepareMocksForProcessorUtils(t)

			convey.Convey(fmt.Sprintf("given %s", test.dataDesc), func() {
				expectedFeed := prepareMakeFeedTestCase(t, toolsMock, fakeIP, test.streamPortsVarValue)

				feed, err := makeFeed()

				if test.negativeCase {
					convey.Convey("should return error", func() {
						convey.So(err, convey.ShouldNotBeNil)
					})
				} else {
					convey.Convey("should not return error and should return expected default feed", func() {
						convey.So(err, convey.ShouldBeNil)
						convey.So(feed, convey.ShouldResemble, expectedFeed)
					})
				}
			})
			ctrl.Finish()
		}
	})
}
