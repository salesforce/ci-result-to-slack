/*
 * Copyright (c) 2021, salesforce.com, inc.
 * All rights reserved.
 * SPDX-License-Identifier: BSD-3-Clause
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/BSD-3-Clause
 */
package main

import (
	"fmt"
	"github.com/salesforce/ci-result-to-slack/internal"
	"strconv"
	"testing"
)

func Test_handleRequest(t *testing.T) {
	tests := []struct {
		name             string
		buildInfo        internal.BuildInfo
		slackClient      internal.SlackClient
		wantReturnString string
		wantErr          bool
		wantErrString    string
	}{
		{
			"success",
			internal.BuildInfo{JobName: "job", BuildURL: "https://sometest", BuildStatus: "SUCCESS", HookURL: "https://slack.com/hook"},
			internal.NewTestClient(false, false),
			fmt.Sprintf(messageSentTemplate, "job"),
			false,
			"",
		},
		{
			"skip posting if success and skipIfSuccess",
			internal.BuildInfo{JobName: "job", BuildURL: "https://sometest", BuildStatus: "SUCCESS", HookURL: "https://slack.com/hook", SkipIfSuccess: true},
			internal.NewTestClient(false, false),
			skippedPostingMessage,
			false,
			"",
		},
		{
			"error processing env vars should throw error",
			internal.BuildInfo{},
			internal.NewTestClient(false, false),
			"",
			true,
			"environment variable error: required key JOB_NAME missing value",
		},
		{
			"error posting to slack should throw error for webhook",
			internal.BuildInfo{JobName: "job", BuildURL: "https://sometest", BuildStatus: "SUCCESS", HookURL: "https://slack.com/hook"},
			internal.NewTestClient(false, true),
			"",
			true,
			internal.WebhookMessageTestErr,
		},
		{
			"error posting to slack should throw error for webhook",
			internal.BuildInfo{JobName: "job", BuildURL: "https://sometest", BuildStatus: "SUCCESS", DestChannelId: "8675309", OauthToken: "token"},
			internal.NewTestClient(true, false),
			"",
			true,
			internal.ChannelMessageTestErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.buildInfo.JobName != "" {
				t.Setenv("JOB_NAME", tt.buildInfo.JobName)
			}
			if tt.buildInfo.BuildURL != "" {
				t.Setenv("BUILD_URL", tt.buildInfo.BuildURL)
			}
			if tt.buildInfo.BuildStatus != "" {
				t.Setenv("BUILD_STATUS", tt.buildInfo.BuildStatus)
			}
			if tt.buildInfo.HookURL != "" {
				t.Setenv("HOOK_URL", tt.buildInfo.HookURL)
			}
			if tt.buildInfo.OauthToken != "" {
				t.Setenv("OAUTH_TOKEN", tt.buildInfo.OauthToken)
			}
			if tt.buildInfo.DestChannelId != "" {
				t.Setenv("DEST_CHANNEL_ID", tt.buildInfo.DestChannelId)
			}
			t.Setenv("SKIP_IF_SUCCESS", strconv.FormatBool(tt.buildInfo.SkipIfSuccess))
			t.Setenv("SUPPRESS_USAGE", "T")
			got, err := handleRequest(tt.slackClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantReturnString {
				t.Errorf("handleRequest() got = %v, wantReturnString %v", got, tt.wantReturnString)
			}
			if err != nil && err.Error() != tt.wantErrString {
				t.Errorf("handleRequest() err = %s, wantErrString %v", err.Error(), tt.wantErrString)
			}
		})
	}
}
