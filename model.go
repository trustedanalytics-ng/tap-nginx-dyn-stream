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

const (
	nginxConfTemplate = `
daemon off;
user nginx;
error_log /dev/stdout {{ .LogLevel }};
worker_processes auto;
worker_rlimit_nofile 8192;
events {
	worker_connections 8000;
}

stream {
    {{ range $streamserver := .StreamServers }}
    server {
        listen {{ $streamserver.ListenIP }}:{{ $streamserver.ListenPort }} ssl;
        proxy_pass 127.0.0.1:{{ $streamserver.ProxyPassPort }};

        include ssl.conf;
    }

    {{ end }}
}

`
)

type StreamServer struct {
	ListenPort    uint16
	ProxyPassPort uint16
	ListenIP      string
}

type TemplateFeed struct {
	StreamServers []StreamServer
	LogLevel      string
}
