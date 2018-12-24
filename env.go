// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logging

import (
	"os"
	"strconv"
	"time"
)

// Variables copied from worker.js.
var (
	supervisorHostname     = os.Getenv("SUPERVISOR_HOSTNAME")
	supervisorInternalPort = os.Getenv("SUPERVISOR_INTERNAL_PORT")
	functionTimeoutSec, _  = strconv.ParseInt(os.Getenv("FUNCTION_TIMEOUT_SEC"), 10, 64)
)

// Constants copied from worker.js.
const (
	maxLogBatchEntries    = 1500
	maxLogBatchLength     = 150000
	supervisorKillTimeout = 5 * time.Second
)

var (
	supervisorLogTimeout = time.Duration(max(60, functionTimeoutSec)) * time.Second
)

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
