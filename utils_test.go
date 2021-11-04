// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2021 Canonical Ltd
 *
 *  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *
 * SPDX-License-Identifier: Apache-2.0'
 */

package hooks

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: add more tests (see trello: )

func TestEnvVars(t *testing.T) {
	// Arrange
	os.Setenv(snapEnv, "/snap/testsnap/x1")
	os.Setenv(snapCommonEnv, "/snap/testsnap/common")
	os.Setenv(snapDataEnv, "/var/snap/testsnap/x1")
	os.Setenv(snapInstNameEnv, "testsnap")
	os.Setenv(snapRevEnv, "2112")

	// Test
	err := getEnvVars()

	// Assert values
	assert.Nil(t, err)
	assert.Equal(t, Snap, "/snap/testsnap/x1")
	assert.Equal(t, SnapCommon, "/snap/testsnap/common")
	assert.Equal(t, SnapData, "/var/snap/testsnap/x1")
	assert.Equal(t, SnapInst, "testsnap")
	assert.Equal(t, SnapRev, "2112")
	assert.Equal(t, SnapConf, "/snap/testsnap/x1/config")
	assert.Equal(t, SnapDataConf, "/var/snap/testsnap/x1/config")
}

func TestGetConfigEnvVar(t *testing.T) {
	var env string
	var ok bool

	// TODO: make this a data driven test (ie. reduce dup code)
	// Test
	env, ok = getConfigEnvVar("service.port", nil)
	assert.True(t, ok)
	assert.Equal(t, env, "SERVICE_PORT")

	// test invalid key
	env, ok = getConfigEnvVar("service.foo", nil)
	assert.False(t, ok)

	// test extra key
	var extraConf = map[string]string{
		"service.mykey":   "SERVICE_MYKEY",
		"service.mykey-2": "SERVICE_MYKEY2",
	}

	// extra key exists
	env, ok = getConfigEnvVar("service.mykey", extraConf)
	assert.True(t, ok)
	assert.Equal(t, env, "SERVICE_MYKEY")

	// extra key exists w/hyphen
	env, ok = getConfigEnvVar("service.mykey-2", extraConf)
	assert.True(t, ok)
	assert.Equal(t, env, "SERVICE_MYKEY2")

	// extra key doesn't exist
	env, ok = getConfigEnvVar("service.fubar", extraConf)
	assert.False(t, ok)
}

func TestSetConfig(t *testing.T) {
	key, value := "mykey", "myvalue"

	cli := NewSnapCtl()
	err := cli.SetConfig(key, value)
	require.Nilf(t, err, "Error setting config.", err)

	// check using snapctl
	require.Equal(t, value, getConfigValue(t, key))
}

func TestUnsetConfig(t *testing.T) {
	key, value := "mykey2", "myvalue"

	// make sure this isn't already set
	require.Equal(t, "", getConfigValue(t, key))

	// set using snapctl
	setConfigValue(t, key, value)

	// check using snapctl
	require.Equal(t, value, getConfigValue(t, key))

	// set using the library
	cli := NewSnapCtl()
	err := cli.UnsetConfig(key)
	require.Nilf(t, err, "Error un-setting config.", err)

	// make sure it has been unset
	require.Equal(t, "", getConfigValue(t, key))
}

func TestStartMultiple(t *testing.T) {
	t.Skipf("need to run in an active snap")
	cli := NewSnapCtl()
	err := cli.StartMultiple(false, "edgexfoundry.consul")
	require.NoError(t, err, "Error getting services.")
}
// utility testing functions

func setConfigValue(t *testing.T, key, value string) {
	err := exec.Command("snapctl", "set", fmt.Sprintf("%s=%s", key, value)).Run()
	require.Nilf(t, err, "Error setting config value via snapctl.")
}

func getConfigValue(t *testing.T, key string) string {
	out, err := exec.Command("snapctl", "get", key).Output()
	require.Nilf(t, err, "Error getting config value via snapctl.")
	return strings.TrimSpace(string(out))
}
