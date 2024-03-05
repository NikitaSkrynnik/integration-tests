// Copyright (c) 2021-2022 Doc.ai and/or its affiliates.
//
// Copyright (c) 2023-2024 Cisco and/or its affiliates.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package logs exports helper functions for storing logs from containers.
package logs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	jaegerAPI      = "api"
	jaegerServices = "services"
	jaegerTraces   = "traces"
)

type jaegerAPIClient struct {
	client        http.Client
	apiServerPort int
}

type Log struct {
	Fields []map[string]string
}

type Span struct {
	TraceID       string
	SpanID        string
	OperationName string
	Logs          []Log
}

type Trace struct {
	TraceID string
	Spans   []Span
}
type Traces struct {
	Data []Trace
}

func TestGetJaegerTraces(t *testing.T) {
	initialize()

	kubeClients := make([]*kubernetes.Clientset, 0)
	for _, cfg := range kubeConfigs {
		kubeconfig, err := clientcmd.BuildConfigFromFlags("", cfg)
		if err != nil {
			logrus.Fatal(err.Error())
		}
		kubeconfig.QPS = float32(config.WorkerCount) * 1000
		kubeconfig.Burst = int(kubeconfig.QPS) * 2
		kubeClient, err := kubernetes.NewForConfig(kubeconfig)
		if err != nil {
			logrus.Fatal(err.Error())
		}
		kubeClients = append(kubeClients, kubeClient)
	}

	//client := kubeClients[0]

	//path := fmt.Sprintf("/api/v1/namespaces/%s/service/%s/portforward", "observability", "jaeger")

	j := &jaegerAPIClient{
		apiServerPort: 16686,
	}

	result := map[string]string{}
	services := j.getServices()
	for _, s := range services {
		result[s] = j.getTracesByService(s)
	}
	service := result["nsmgr-sh9g7"]
	fmt.Printf("service: %v\n", service)

	traces := Traces{}
	err := json.Unmarshal([]byte(service), &traces)
	fmt.Printf("err: %v\n", err)

	//fmt.Printf("traces: %v\n", traces)

	filteredTraces := Traces{}
	filteredTraces.Data = make([]Trace, 0)
	for _, trace := range traces.Data {
		if trace.Spans[0].OperationName != "grpc.health.v1.Health/Check" {
			filteredTraces.Data = append(filteredTraces.Data, trace)
		}
	}

	fmt.Printf("filteredTraces: %v\n", filteredTraces)

	for _, trace := range filteredTraces.Data {
		for _, span := range trace.Spans {
			fmt.Printf("operation: %s\n", span.OperationName)
		}

		fmt.Println()
		fmt.Println()
	}
}

func (j *jaegerAPIClient) getTracesByService(service string) string {
	url := fmt.Sprintf("%v?service=%v&limit=1500&lookback=2d", urlToLocalHost(j.apiServerPort, jaegerAPI, jaegerTraces), service)
	resp, err := j.client.Get(url)
	if err != nil {
		logrus.Errorf("An error during get jaeger traces from API: %v", err)
		return ""
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("An error during read jaeger traces response: %v", err)
		return ""
	}
	return string(bytes)
}

func (j *jaegerAPIClient) getServices() []string {
	resp, err := j.client.Get(urlToLocalHost(j.apiServerPort, jaegerAPI, jaegerServices))
	if err != nil {
		logrus.Errorf("An error during get jaeger services from API: %v", err)
		return nil
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	jsonObject := map[string]interface{}{}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("An error during read jaeger services response: %v", err)
		return nil
	}
	err = json.Unmarshal([]byte(strings.ReplaceAll(string(bytes), "\\\"", "\"")), &jsonObject)
	if err != nil {
		logrus.Errorf("An error during unmarshal jaeger services response: %v", err)
		return nil
	}

	if v, ok := jsonObject["data"].([]interface{}); ok {
		result := []string{}
		for _, item := range v {
			result = append(result, fmt.Sprint(item))
		}
		logrus.Info(v)
		return result
	}

	return nil
}

func urlToLocalHost(port int, parts ...string) string {
	u, _ := url.Parse(fmt.Sprintf("http://0.0.0.0:%v", port))
	fullPath := append([]string{}, u.Path)
	fullPath = append(fullPath, parts...)
	u.Path = path.Join(fullPath...)
	return u.String()
}
