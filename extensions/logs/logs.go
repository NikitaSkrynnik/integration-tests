// Copyright (c) 2021-2022 Doc.ai and/or its affiliates.
//
// Copyright (c) 2023 Cisco and/or its affiliates.
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
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	"github.com/networkservicemesh/gotestmd/pkg/bash"
)

var (
	once                       sync.Once
	config                     Config
	ctx                        context.Context
	kubeConfigs                []string
	testNamespacesMatchRegex   *regexp.Regexp
	nsmSystemMatchRegex        *regexp.Regexp
	runner                     *bash.Bash
	clusterDumpSingleOperation *singleOperation
	nsmSystemNumber            int
)

// Config is env config to setup log collecting.
type Config struct {
	ArtifactsDir         string        `default:"logs" desc:"Directory for storing container logs" envconfig:"ARTIFACTS_DIR"`
	Timeout              time.Duration `default:"10s" desc:"Context timeout for kubernetes queries" split_words:"true"`
	WorkerCount          int           `default:"8" desc:"Number of log collector workers" split_words:"true"`
	MaxKubeConfigs       int           `default:"3" desc:"Number of used kubeconfigs" split_words:"true"`
	TestNamespaces       string        `default:"ns-.*" desc:"Regex of allowed namespaces" split_words:"true"`
	NsmSystemNamespace   string        `default:"(nsm-system)|(spire)|(observability)" desc:"Regex of allowed namespaces" split_words:"true"`
	LogCollectionEnabled bool          `default:"true" desc:"Boolean variable which enables log collection" split_words:"true"`
}

// nolint: gocyclo
func initialize() {
	if err := envconfig.Usage("logs", &config); err != nil {
		logrus.Fatal(err.Error())
	}

	if err := envconfig.Process("logs", &config); err != nil {
		logrus.Fatal(err.Error())
	}

	testNamespacesMatchRegex = regexp.MustCompile(config.TestNamespaces)
	nsmSystemMatchRegex = regexp.MustCompile(config.NsmSystemNamespace)

	var singleClusterKubeConfig = os.Getenv("KUBECONFIG")

	if singleClusterKubeConfig == "" {
		singleClusterKubeConfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	kubeConfigs = []string{}

	for i := 1; i <= config.MaxKubeConfigs; i++ {
		kubeConfig := os.Getenv("KUBECONFIG" + fmt.Sprint(i))
		if kubeConfig != "" {
			kubeConfigs = append(kubeConfigs, kubeConfig)
		}
	}

	if len(kubeConfigs) == 0 {
		kubeConfigs = append(kubeConfigs, singleClusterKubeConfig)
	}

	runner, _ = bash.New()

	ctx, _ = signal.NotifyContext(context.Background(),
		os.Interrupt,
		os.Kill,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	clusterDumpSingleOperation = newSingleOperation(func() {})
	nsmSystemNumber = 0
}

// ClusterDump saves logs from all pods in specified namespaces
func ClusterDump(folder string) {
	once.Do(initialize)

	clusterDumpSingleOperation.Body = func() {
		if ctx.Err() != nil {
			return
		}
		nsmSystemNumber++
		for i := range kubeConfigs {
			nsString, _, _, _ := runner.Run(fmt.Sprintf(`kubectl --kubeconfig %v get ns -o go-template='{{range .items}}{{ .metadata.name }} {{end}}'`, kubeConfigs[i]))
			nsList := strings.Split(nsString, " ")

			fmt.Printf("nsList: %v\n", nsList)

			dir := filepath.Join(config.ArtifactsDir, fmt.Sprintf("cluster%v", i), folder)
			runClusterDump(kubeConfigs[i], dir, strings.Join(filterNamespaces(nsList, testNamespacesMatchRegex), ","))

			dir = filepath.Join(dir, fmt.Sprintf("nsm-system%d", nsmSystemNumber))
			runClusterDump(kubeConfigs[i], dir, strings.Join(filterNamespaces(nsList, nsmSystemMatchRegex), ","))
		}
	}
	clusterDumpSingleOperation.Run()
}

func runClusterDump(kubeConfig, suitedir, nsList string) (int, error) {
	_, _, exitCode, err := runner.Run(
		fmt.Sprintf("kubectl --kubeconfig %v cluster-info dump --output-directory=%s --namespaces %s",
			kubeConfig,
			suitedir,
			nsList))

	return exitCode, err
}

func filterNamespaces(nsList []string, matchRegex *regexp.Regexp) []string {
	result := make([]string, 0)

	for i := range nsList {
		if matchRegex.MatchString(nsList[i]) {
			fmt.Printf("Added to nsList: %s\n", nsList[i])
			result = append(result, nsList[i])
		}
	}

	fmt.Println()
	fmt.Println()
	fmt.Println()
	return result
}
