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

package util

import (
	"os"
	"os/signal"
	"sync"

	commonLogger "github.com/trustedanalytics-ng/tap-go-common/logger"
)

var logger, _ = commonLogger.InitLogger("util")

func GetTerminationObserverChannel() chan os.Signal {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	return channel
}

func TerminationObserver(waitGroup *sync.WaitGroup, appName string) {
	<-GetTerminationObserverChannel()
	logger.Info(appName, "is going to be stopped now...")
	waitGroup.Wait()
	logger.Info(appName, "stopped!")
	os.Exit(0)
}

func GetEnvValueOrDefault(envName string, defaultValue string) string {
	if value := os.Getenv(envName); value != "" {
		return value
	}
	return defaultValue
}
