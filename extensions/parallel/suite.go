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

// Package parallel provides functions to run suite tests in parallel
package parallel

import (
	"reflect"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

func recoverAndFailOnPanic(t *testing.T) {
	r := recover()
	failOnPanic(t, r)
}

func failOnPanic(t *testing.T, r interface{}) {
	if r != nil {
		t.Errorf("test panicked: %v\n%s", r, debug.Stack())
		t.FailNow()
	}
}

// Run runs suite tests in parallel
func Run(t *testing.T, s suite.TestingSuite, excludedTests ...string) {
	excludedTestsSet := make(map[string]struct{})
	for _, test := range excludedTests {
		excludedTestsSet[test] = struct{}{}
	}

	defer recoverAndFailOnPanic(t)
	var suiteSetupDone bool

	s.SetT(t)
	tests := []testing.InternalTest{}
	methodFinder := reflect.TypeOf(s)

	t.Cleanup(func() {
		if suiteSetupDone {
			if tearDownAllSuite, ok := s.(suite.TearDownAllSuite); ok {
				tearDownAllSuite.TearDownSuite()
			}
		}
	})

	for i := 0; i < methodFinder.NumMethod(); i++ {
		method := methodFinder.Method(i)
		if ok := strings.HasPrefix(method.Name, "Test"); !ok {
			continue
		}
		parallel := true
		if _, ok := excludedTestsSet[method.Name]; ok {
			parallel = false
		}

		if !suiteSetupDone {
			if setupAllSuite, ok := s.(suite.SetupAllSuite); ok {
				setupAllSuite.SetupSuite()
			}

			suiteSetupDone = true
		}

		test := newTest(t, s, methodFinder, &method, parallel)
		tests = append(tests, test)
	}

	if len(tests) == 0 {
		t.Log("warning: no tests to run")
		return
	}

	// run sub-tests in a group so tearDownSuite is called in the right order
	for _, test := range tests {
		t.Run(test.Name, test.F)
	}
}

func newTest(t *testing.T, s suite.TestingSuite, methodFinder reflect.Type, method *reflect.Method, parallel bool) testing.InternalTest {
	return testing.InternalTest{
		Name: method.Name,
		F: func(testingT *testing.T) {
			defer recoverAndFailOnPanic(t)

			if parallel {
				testingT.Parallel()
			}

			subS := reflect.New(reflect.ValueOf(s).Elem().Type())
			subS.MethodByName("SetT").Call([]reflect.Value{reflect.ValueOf(testingT)})

			defer func() {
				r := recover()

				if afterTestSuite, ok := subS.Interface().(suite.AfterTest); ok {
					afterTestSuite.AfterTest(s.T().Name(), method.Name)
				}

				if tearDownTestSuite, ok := subS.Interface().(suite.TearDownTestSuite); ok {
					tearDownTestSuite.TearDownTest()
				}

				failOnPanic(t, r)
			}()

			if setupTestSuite, ok := subS.Interface().(suite.SetupTestSuite); ok {
				setupTestSuite.SetupTest()
			}
			if beforeTestSuite, ok := subS.Interface().(suite.BeforeTest); ok {
				beforeTestSuite.BeforeTest(methodFinder.Elem().Name(), method.Name)
			}

			method.Func.Call([]reflect.Value{subS})
		},
	}
}
