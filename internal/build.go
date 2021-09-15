/*
 * Copyright (c) 2021, salesforce.com, inc.
 * All rights reserved.
 * SPDX-License-Identifier: BSD-3-Clause
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/BSD-3-Clause
 */
package internal

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"os"
	"strconv"
)

var (
	successKey      = "SUCCESS"
	fixedKey        = "FIXED"
	unstableKey     = "UNSTABLE"
	unknownKey      = "UNKNOWN"
	failureKey      = "FAILURE"
	stillFailingKey = "STILL FAILING"

	successStatus      = Status{text: "Success", color: "good"}
	fixedStatus        = Status{text: "Fixed", color: "good"}
	unstableStatus     = Status{text: "Unstable", color: "warning"}
	unknownStatus      = Status{text: "Unknown", color: "warning"}
	failedStatus       = Status{text: "Failed", color: "danger"}
	stillFailingStatus = Status{text: "Still Failing", color: "danger"}

	defaultStatus = unknownStatus

	statusMap = map[string]Status{
		successKey:      successStatus,
		fixedKey:        fixedStatus,
		unstableKey:     unstableStatus,
		unknownKey:      unknownStatus,
		failureKey:      failedStatus,
		stillFailingKey: stillFailingStatus,
	}

	branchFieldTitle      = "Branch"
	commitFieldTitle      = "Commit"
	buildTimeFieldTitle   = "Time"
	triggeredByFieldTitle = "Triggered By"

	PickRunModeErrorMessage = "please specify either HOOK_URL or both OAUTH_TOKEN and DEST_CHANNEL_ID"
)

/*
Status represents generic build status with text and an associated color
*/
type Status struct {
	text, color string
}

/*
BuildInfo represents the build information passed in from the caller
*/
type BuildInfo struct {
	JobName         string `required:"true" split_words:"true" desc:"Name of the build's job"`
	BuildURL        string `required:"true" split_words:"true" desc:"Direct URL to the build"`
	BuildStatus     string `required:"true" split_words:"true" desc:"Status of build (e.g. currentBuild.currentResult in Jenkins)"`
	HookURL         string `split_words:"true" desc:"Slack Webhook URL set via Incoming Webhooks"`
	DestChannelId   string `split_words:"true" desc:"Destination Channel ID (not the name of the channel)"`
	OauthToken      string `split_words:"true" desc:"OAuth Token used to send message via app"`
	LastBuildStatus string `split_words:"true" default:"UNKNOWN" desc:"Status of last build used to provide contextual build Status"`
	BranchName      string `split_words:"true" desc:"Name of git branch"`
	GitCommit       string `split_words:"true" desc:"Git commit hash"`
	BuildTime       string `split_words:"true" desc:"Build time (e.g. durationString in Jenkins)"`
	TriggeredBy     string `split_words:"true" desc:"The action which triggered the build"`
	SkipIfSuccess   bool   `split_words:"true" desc:"Skip posting if contextual Status is success"`
}

func (buildInfo *BuildInfo) GetContextualStatus() Status {
	status, present := statusMap[buildInfo.BuildStatus]
	if !present {
		return defaultStatus
	}
	lastBuildStatus := statusMap[buildInfo.LastBuildStatus]
	if lastBuildStatus == failedStatus && status == successStatus {
		return fixedStatus
	} else if lastBuildStatus == failedStatus && status == failedStatus {
		return stillFailingStatus
	}
	return status
}

func (buildInfo *BuildInfo) ShouldSkipPosting() bool {
	return buildInfo.SkipIfSuccess && buildInfo.GetContextualStatus() == successStatus
}

func GetBuildInfoFromEnv() (BuildInfo, error) {
	envConfigPrefix := ""
	var buildInfo BuildInfo
	err := envconfig.Process(envConfigPrefix, &buildInfo)
	if err != nil {
		// Primarily exists for testing use cases
		suppressUsage, _ := strconv.ParseBool(os.Getenv("SUPPRESS_USAGE"))
		if !suppressUsage {
			_ = envconfig.Usage(envConfigPrefix, &buildInfo)
		}
		err = fmt.Errorf("environment variable error: %s", err)
	}
	return buildInfo, err
}
