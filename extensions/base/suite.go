// Copyright (c) 2021-2023 Doc.ai and/or its affiliates.
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

// Package base exports base suite type that will be injected into each generated suite.
package base

import (
	"fmt"
	"strings"

	"github.com/networkservicemesh/gotestmd/pkg/suites/shell"
	"github.com/networkservicemesh/integration-tests/extensions/checkout"
	"github.com/networkservicemesh/integration-tests/extensions/logs"
	"github.com/networkservicemesh/integration-tests/extensions/prefetch"
)

// Suite is a base suite for generating tests. Contains extensions that can be used for assertion and automation goals.
type Suite struct {
	shell.Suite
	// Add other extensions here
	checkout checkout.Suite
	prefetch prefetch.Suite
}

// AfterTest stores logs after each test in the suite.
func (s *Suite) AfterTest(_, _ string) {
	logs.ClusterDump()
}

// TearDownSuite stores logs from containers that spawned during SuiteSetup.
func (s *Suite) TearDownSuite() {
	logs.ClusterDump()
}

const (
	sha = "607a6275cb85275ed96c0afd6981f2c24eee275c"
)

// SetupSuite runs all extensions
func (s *Suite) SetupSuite() {
	repo := "NikitaSkrynnik/deployments-k8s"
	version := "grpc"

	s.checkout.Version = version

	if strings.Contains(sha, "tags") {
		s.checkout.Version = sha
		version = strings.ReplaceAll(sha, "tags/", "")
	}

	s.checkout.Dir = "../" // Note: this should be synced with input parameters in gen.go file
	s.checkout.Repository = repo
	s.checkout.SetT(s.T())
	s.checkout.SetupSuite()

	// prefetch
	s.prefetch.SourcesURLs = []string{
		// Note: use urls for local image files.
		// For example:
		//    "file://my-debug-images-for-prefetch.yaml"
		//    "file://deployments-k8s/apps/"
		fmt.Sprintf("https://raw.githubusercontent.com/%v/%v/external-images.yaml", repo, version),
		fmt.Sprintf("https://api.github.com/repos/%v/contents/apps?ref=%v", repo, version),
	}

	s.prefetch.SetT(s.T())
	s.prefetch.SetupSuite()
}
