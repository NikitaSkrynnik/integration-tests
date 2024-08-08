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
	LogsConfig                 Config
	once                       sync.Once
	ctx                        context.Context
	kubeConfigs                []string
	matchRegex                 *regexp.Regexp
	dockerRegex                *regexp.Regexp
	runner                     *bash.Bash
	clusterDumpSingleOperation *singleOperation
)

// Config is env config to setup log collecting.
type Config struct {
	ArtifactsDir         string        `default:"logs" desc:"Directory for storing container logs" envconfig:"ARTIFACTS_DIR"`
	Timeout              time.Duration `default:"10s" desc:"Context timeout for kubernetes queries" split_words:"true"`
	WorkerCount          int           `default:"8" desc:"Number of log collector workers" split_words:"true"`
	MaxKubeConfigs       int           `default:"3" desc:"Number of used kubeconfigs" split_words:"true"`
	AllowedNamespaces    string        `default:"(ns-.*)|(nsm-system)|(spire)|(observability)" desc:"Regex of allowed namespaces" split_words:"true"`
	AllowedContainers    string        `default:"(nsc-.*)|(nse-.*)" desc:"Regexp of allowed docker containers" split_words:"true"`
	LogCollectionEnabled bool          `default:"true" desc:"Boolean variable which enables log collection" split_words:"true"`
}

// nolint: gocyclo
func initialize() {
	if err := envconfig.Usage("logs", &LogsConfig); err != nil {
		logrus.Fatal(err.Error())
	}

	if err := envconfig.Process("logs", &LogsConfig); err != nil {
		logrus.Fatal(err.Error())
	}

	matchRegex = regexp.MustCompile(LogsConfig.AllowedNamespaces)
	dockerRegex = regexp.MustCompile(LogsConfig.AllowedContainers)

	var singleClusterKubeConfig = os.Getenv("KUBECONFIG")

	if singleClusterKubeConfig == "" {
		singleClusterKubeConfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	kubeConfigs = []string{}

	for i := 1; i <= LogsConfig.MaxKubeConfigs; i++ {
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

	clusterDumpSingleOperation = newSingleOperation()
}

func Initialize() {
	once.Do(initialize)
}

// ClusterDump saves logs from all pods in specified namespaces
func ClusterDump(suiteName, testName string) {
	clusterDumpSingleOperation.Run(func() {
		if ctx.Err() != nil {
			return
		}
		for i := range kubeConfigs {
			suitedir := filepath.Join(LogsConfig.ArtifactsDir, fmt.Sprintf("cluster%v", i), suiteName, testName)

			nsString, _, _, _ := runner.Run(fmt.Sprintf(`kubectl --kubeconfig %v get ns -o go-template='{{range .items}}{{ .metadata.name }} {{end}}'`, kubeConfigs[i]))
			nsList := strings.Split(nsString, " ")

			_, _, exitCode, err := runner.Run(fmt.Sprintf("kubectl --kubeconfig %v cluster-info dump --output-directory=%s --namespaces %s",
				kubeConfigs[i],
				suitedir,
				strings.Join(filterNamespaces(nsList), ",")))

			if exitCode != 0 {
				logrus.Errorf("An error while getting cluster dump. Exit Code: %v", exitCode)
			}
			if err != nil {
				logrus.Errorf("An error while getting cluster dump. Error: %s", err.Error())
			}

			allContainers, _, _, err := runner.Run(`docker ps --format '{{.Names}}'`)
			if err != nil {
				logrus.Errorf("An error while getting docker containers. Error: %s", err.Error())
				continue
			}
			containerList := strings.Split(strings.TrimSpace(allContainers), "\n")
			containers := filterContainers(containerList)

			for _, container := range containers {
				_, _, dExitCode, dErr := runner.Run(fmt.Sprintf("docker logs %s >& %s/%s.log", container, suitedir, container))

				if dExitCode != 0 {
					logrus.Errorf("An error while getting docker logs. Exit Code: %v", dExitCode)
				}
				if dErr != nil {
					logrus.Errorf("An error while getting docker logs. Error: %s", dErr.Error())
				}
			}
		}
	})
}

func filterNamespaces(nsList []string) []string {
	result := make([]string, 0)

	for i := range nsList {
		if matchRegex.MatchString(nsList[i]) {
			result = append(result, nsList[i])
		}
	}

	return result
}

func filterContainers(containerList []string) []string {
	result := make([]string, 0)

	for i := range containerList {
		if dockerRegex.MatchString(containerList[i]) {
			result = append(result, containerList[i])
		}
	}

	return result
}

// func MetricsDump(suiteName, testName string) {
// 	testDir := filepath.Join(config.ArtifactsDir, suiteName, testName)
// 	runner.Run(fmt.Sprintf(`kubectl top pod -A > %s/metrics`, testDir))
// }
