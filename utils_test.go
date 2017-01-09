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
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
)

type StringToUintTestCase struct {
	dataDesc       string
	testData       []string
	expectedResult []uint16
}

func TestConvertStrPortsToUint16(t *testing.T) {
	tools = NewUtilTools(&SystemUtils{})

	tests := []StringToUintTestCase{
		{
			dataDesc:       "proper string ports",
			testData:       []string{"443", "80", "12", "24"},
			expectedResult: []uint16{443, 80, 12, 24},
		},
		{
			dataDesc: "string ports out of uint16 range",
			testData: []string{"443", "80", "12", "65537"},
		},
		{
			dataDesc: "string ports and not numeric value",
			testData: []string{"443", "80", "12", "asdasd"},
		},
	}

	for _, test := range tests {
		convey.Convey(fmt.Sprintf("Given array with %s", test.dataDesc), t, func() {
			uPorts, err := tools.ConvertStrPortsToUint16(test.testData)

			if len(test.expectedResult) > 0 {
				convey.Convey("error should be nil, actual result should resemble expected result", func() {
					convey.So(err, convey.ShouldBeNil)
					convey.So(uPorts, convey.ShouldResemble, test.expectedResult)
				})
			} else {
				convey.Convey("error should not be nil", func() {
					convey.So(err, convey.ShouldNotBeNil)
				})
			}
		})
	}
}

func prepareMocksForUtils(t *testing.T) (*gomock.Controller, UtilAPI, *MockSystemUtilAPI) {
	mockCtrl := gomock.NewController(t)
	systemUtilMock := NewMockSystemUtilAPI(mockCtrl)

	tools := NewUtilTools(systemUtilMock)

	return mockCtrl, tools, systemUtilMock
}

func TestGetEth0IP(t *testing.T) {
	convey.Convey("Test GetEth0IP", t, func() {
		mockCtrl, tools, systemUtilsMock := prepareMocksForUtils(t)

		convey.Convey("Given that system command will execute properly", func() {
			fakeIP := "10.10.12.13"

			systemUtilsMock.EXPECT().ExecuteCommand(Eth0IPCommand).Return(fakeIP, nil)

			ip, err := tools.GetEth0IP()

			convey.Convey("GetEth0IP should not return any error, should return proper IP address", func() {
				convey.So(err, convey.ShouldBeNil)
				convey.So(ip, convey.ShouldEqual, fakeIP)
			})
		})

		convey.Convey("Given that system command will return some error", func() {
			systemUtilsMock.EXPECT().ExecuteCommand(Eth0IPCommand).Return("", errors.New("cannot find command!"))

			_, err := tools.GetEth0IP()

			convey.Convey("GetEth0IP should return also error", func() {
				convey.So(err, convey.ShouldNotBeNil)
			})
		})

		convey.Reset(func() {
			mockCtrl.Finish()
		})
	})
}
