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
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/networkservicemesh/gotestmd/pkg/bash"
)

const (
	defaultQPS = 5 // this is default value for QPS of kubeconfig. See at documentation.
)

var (
	once                       sync.Once
	config                     Config
	ctx                        context.Context
	kubeClients                []kubernetes.Interface
	kubeConfigs                []string
	matchRegex                 *regexp.Regexp
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

	matchRegex = regexp.MustCompile(config.AllowedNamespaces)

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

	var apiVersions = []string{"client.authentication.k8s.io/v1", "client.authentication.k8s.io/v1beta1", "client.authentication.k8s.io/v1alpha1"}

	for _, cfg := range kubeConfigs {
		for _, apiVersion := range apiVersions {
			kubeconfig, err := clientcmd.BuildConfigFromFlags("", cfg)
			if err != nil {
				logrus.Warn(err.Error())
				continue
			}
			kubeconfig.QPS = float32(config.WorkerCount) * defaultQPS
			kubeconfig.Burst = int(kubeconfig.QPS) * 2

			if kubeconfig.ExecProvider != nil {
				kubeconfig.ExecProvider.APIVersion = apiVersion
			}

			kubeClient, err := kubernetes.NewForConfig(kubeconfig)
			if err != nil {
				logrus.Warn(err.Error())
				continue
			}
			kubeClients = append(kubeClients, kubeClient)
			break
		}
	}
	if len(kubeClients) == 0 {
		logrus.Warn("k8s clients weren't initialized properly. loggig is disabled")
		return
	}

	runner, _ = bash.New()

	ctx, _ = signal.NotifyContext(context.Background(),
		os.Interrupt,
		os.Kill,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	clusterDumpSingleOperation = newSingleOperation(func() {
		if ctx.Err() != nil {
			return
		}
		for i, client := range kubeClients {
			suitedir := filepath.Join(config.ArtifactsDir, fmt.Sprintf("cluster%v", i))
			nsList, _ := client.CoreV1().Namespaces().List(ctx, v1.ListOptions{})

			os.Setenv("KUBECONFIG", kubeConfigs[i])
			_, _, exitCode, err := runner.Run(
				fmt.Sprintf("kubectl cluster-info dump --output-directory=%s --namespaces %s",
					suitedir,
					strings.Join(filterNamespaces(nsList), ",")))

			if exitCode != 0 {
				logrus.Errorf("An error while getting cluster dump. Exit Code: %v", exitCode)
			}
			if err != nil {
				logrus.Errorf("An error while getting cluster dump. Error: %s", err.Error())
			}
		}
	})
}

// ClusterDump saves logs from all pods in specified namespaces
func ClusterDump() {
	once.Do(initialize)
	clusterDumpSingleOperation.Run()
}

func filterNamespaces(nsList *corev1.NamespaceList) []string {
	result := make([]string, 0)

	for i := range nsList.Items {
		if matchRegex.MatchString(nsList.Items[i].Name) {
			result = append(result, nsList.Items[i].Name)
		}
	}

	return result
}
