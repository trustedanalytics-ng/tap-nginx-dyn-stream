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
	"io/ioutil"
	"os"
	"text/template"

	"github.com/op/go-logging"

	commonLogger "github.com/trustedanalytics/tap-go-common/logger"
)

var logger = initLogger()

const (
	streamPortsEnvVarName = "STREAM_PORTS"
	logLevelEnvVarName    = "LOG_LEVEL"

	nginxConfPath = "/etc/nginx/nginx.conf"

	defaultStreamPorts = "443:80"
	defaultLogLevel    = "notice"
)

func initLogger() *logging.Logger {
	logger, err := commonLogger.InitLogger("main")
	if err != nil {
		panic("cannot initialize logger! err: " + err.Error())
	}
	return logger
}

func main() {
	logger.Info("nginx config initialization started")

	err := initConfigFromEnv()
	if err != nil {
		logger.Fatal("config initialization failed! Error: ", err)
	}

	if logger.IsEnabledFor(logging.DEBUG) {
		generatedConfBytes, err := ioutil.ReadFile(nginxConfPath)
		if err != nil {
			logger.Errorf("failed to read %v for debugging, err: %v", nginxConfPath, err)
		} else {
			logger.Debug("generated nginx configuration: ", string(generatedConfBytes))
		}
	}

	logger.Info("config initialization application ended successfully")
}

func initConfigFromEnv() error {
	logger.Debug("parsing nginx configuration template")
	tmpl, err := template.New("nginx").Parse(nginxConfTemplate)
	if err != nil {
		logger.Error("parsing nginx configuration template failed: ", err)
		return err
	}

	logger.Debug("making feed for template filling")
	feed, err := makeFeed()
	if err != nil {
		logger.Error("making feed failed: ", err)
		return err
	}

	logger.Debug("creating nginx.conf file")
	w, err := os.Create(nginxConfPath)
	if err != nil {
		logger.Error("creation of nginx.conf file failed: ", err)
		return err
	}

	logger.Debug("executing template fulfillment")
	err = tmpl.Execute(w, feed)
	if err != nil {
		logger.Error("template filling error: ", err)
		return err
	}
	return nil
}
