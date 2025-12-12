/*
 * Copyright (c) 2021, salesforce.com, inc.
 * All rights reserved.
 * SPDX-License-Identifier: BSD-3-Clause
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/BSD-3-Clause
 */
package internal

import (
	"reflect"
	"testing"
)

func Test_GetBuildInfoFromEnvReturnsNoErrorWhenEnvVarsSet(t *testing.T) {
	t.Setenv("HOOK_URL", "test")
	t.Setenv("JOB_NAME", "test")
	t.Setenv("BUILD_STATUS", "test")
	t.Setenv("BUILD_URL", "test")
	_, err := GetBuildInfoFromEnv()
	if err != nil {
		t.Error("Expected no error when environment variables set")
	}
}

func Test_GetBuildInfoFromEnvReturnsErrorWhenEnvVarsNotSet(t *testing.T) {
	t.Setenv("SUPPRESS_USAGE", "T")
	_, err := GetBuildInfoFromEnv()
	if err == nil {
		t.Error("Expected error when environment variables set")
	}
}

func Test_GetContextualStatus(t *testing.T) {
	type args struct {
		buildInfo BuildInfo
	}
	tests := []struct {
		name string
		args args
		want Status
	}{
		{
			"success",
			args{BuildInfo{BuildStatus: successKey, LastBuildStatus: ""}},
			successStatus,
		},
		{
			"fixed - direct access",
			args{BuildInfo{BuildStatus: fixedKey, LastBuildStatus: ""}},
			fixedStatus,
		},
		{
			"fixed - contextual",
			args{BuildInfo{BuildStatus: successKey, LastBuildStatus: failureKey}},
			fixedStatus,
		},
		{
			"unstable",
			args{BuildInfo{BuildStatus: unstableKey, LastBuildStatus: ""}},
			unstableStatus,
		},
		{
			"unknown",
			args{BuildInfo{BuildStatus: "blah", LastBuildStatus: ""}},
			unknownStatus,
		},
		{
			"failure",
			args{BuildInfo{BuildStatus: failureKey, LastBuildStatus: ""}},
			failedStatus,
		},
		{
			"still failing - direct access",
			args{BuildInfo{BuildStatus: stillFailingKey, LastBuildStatus: ""}},
			stillFailingStatus,
		},
		{
			"still failing - contextual",
			args{BuildInfo{BuildStatus: failureKey, LastBuildStatus: failureKey}},
			stillFailingStatus,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.buildInfo.GetContextualStatus(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetContextualStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ShouldSkipPosting(t *testing.T) {
	tests := []struct {
		name           string
		buildInfo      BuildInfo
		expectedResult bool
	}{
		{
			"should skip when SkipIfSuccess=true and status is success",
			BuildInfo{BuildStatus: successKey, SkipIfSuccess: true},
			true,
		},
		{
			"should NOT skip when status is fixed (success after failure) - fixed is newsworthy",
			BuildInfo{BuildStatus: successKey, LastBuildStatus: failureKey, SkipIfSuccess: true},
			false,
		},
		{
			"should NOT skip when SkipIfSuccess=false",
			BuildInfo{BuildStatus: successKey, SkipIfSuccess: false},
			false,
		},
		{
			"should NOT skip when SkipIfSuccess=true but status is failure",
			BuildInfo{BuildStatus: failureKey, SkipIfSuccess: true},
			false,
		},
		{
			"should NOT skip when SkipIfSuccess=true but status is unstable",
			BuildInfo{BuildStatus: unstableKey, SkipIfSuccess: true},
			false,
		},
		{
			"should NOT skip when SkipIfSuccess=true but status is still failing",
			BuildInfo{BuildStatus: failureKey, LastBuildStatus: failureKey, SkipIfSuccess: true},
			false,
		},
		{
			"should NOT skip when SkipIfSuccess=true but status is unknown",
			BuildInfo{BuildStatus: "UNKNOWN", SkipIfSuccess: true},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.buildInfo.ShouldSkipPosting(); got != tt.expectedResult {
				t.Errorf("ShouldSkipPosting() = %v, want %v", got, tt.expectedResult)
			}
		})
	}
}

func Test_GetBuildInfoFromEnv_UsageOutput(t *testing.T) {
	// Don't set SUPPRESS_USAGE, so usage will be printed to stderr
	// This tests the error path where usage is displayed
	_, err := GetBuildInfoFromEnv()
	if err == nil {
		t.Error("Expected error when required environment variables missing")
	}
}
