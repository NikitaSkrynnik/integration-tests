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

package logs

import (
	"sync/atomic"
	"time"
)

const (
	notScheduled        = int32(0)
	running             = int32(1)
	scheduledAndRunning = int32(3)
)

type singleOperation struct {
	Body  func()
	state int32
}

func (o *singleOperation) Wait() {
	for atomic.AddInt32(&o.state, 0) != notScheduled {
		<-time.After(time.Millisecond * 25)
	}
}

// newSingleOperation creates an operation which should be invoked once by run period. Can be used in cases where required the last run.
func newSingleOperation(body func()) *singleOperation {
	if body == nil {
		panic("body can not be nil")
	}
	return &singleOperation{Body: body, state: notScheduled}
}

func (o *singleOperation) Run() {
	if !atomic.CompareAndSwapInt32(&o.state, notScheduled, running) {
		if !atomic.CompareAndSwapInt32(&o.state, running, scheduledAndRunning) {
			if !atomic.CompareAndSwapInt32(&o.state, notScheduled, running) {
				return
			}
		} else {
			return
		}
	}
	o.Body()
	if !atomic.CompareAndSwapInt32(&o.state, running, notScheduled) {
		o.Body()
		atomic.StoreInt32(&o.state, notScheduled)
	}
}
