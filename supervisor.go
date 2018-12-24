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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)



func init() {
	if supervisorHostname != "" && supervisorInternalPort != "" {
		loggingCtx.initialize()
	}
}

func newSupervisorRequest(path string, v interface{}) (*http.Request, error) {
	postData, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", (&url.URL{
		Scheme: "http",
		Host:   supervisorHostname + ":" + supervisorInternalPort,
		Path:   path,
	}).String(), bytes.NewBuffer(postData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(postData)))

	return req, nil
}

func postToSupervisor(path string, v interface{}, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := newSupervisorRequest(path, v)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		if err == ctx.Err() {
			return errors.New("timeout when calling supervisor")
		}

		return fmt.Errorf("error when calling supervisor: %s\n", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bad response code from supervisor: %d\n", resp.StatusCode)
	}

	return nil
}

func killInstance() {
	err := postToSupervisor("/_ah/kill", nil, supervisorKillTimeout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	// Exit code 16 is copied over from worker.js.
	os.Exit(16)
}
