/*
 * Copyright (c) 2021, salesforce.com, inc.
 * All rights reserved.
 * SPDX-License-Identifier: BSD-3-Clause
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/BSD-3-Clause
 */
package internal

import (
	"fmt"
	"github.com/slack-go/slack"
	"reflect"
	"testing"
)

var (
	jobName    = "my test job"
	buildURL   = "https://www.com"
	branchName = "test branch name"

	branchNameField  = getAttachmentField(branchFieldTitle, branchName)
	commit           = "8675309"
	commitField      = getAttachmentField(commitFieldTitle, commit)
	buildTime        = "0m 3s"
	buildTimeField   = getAttachmentField(buildTimeFieldTitle, buildTime)
	triggeredBy      = "Pull request"
	triggeredByField = getAttachmentField(triggeredByFieldTitle, triggeredBy)

	emptyAttachmentFields []slack.AttachmentField
)

func Test_PostToSlack(t *testing.T) {
	happySlackClient := NewTestClient(false, false)

	type args struct {
		buildInfo   BuildInfo
		slackClient SlackClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		errTest string
	}{
		{
			"run mode error thrown if HOOK_URL, DEST_CHANNEL_ID, and OAUTH_TOKEN not specified",
			args{
				buildInfo: BuildInfo{
					JobName:     "somejob",
					BuildURL:    "someurl",
					BuildStatus: "somestatus",
				},
				slackClient: happySlackClient,
			},
			true,
			PickRunModeErrorMessage,
		},
		{
			"run mode error thrown if OAUTH_TOKEN not specified with DEST_CHANNEL_ID",
			args{
				buildInfo: BuildInfo{
					JobName:       "somejob",
					BuildURL:      "someurl",
					BuildStatus:   "somestatus",
					DestChannelId: "somechannelid",
				},
				slackClient: happySlackClient,
			},
			true,
			PickRunModeErrorMessage,
		},
		{
			"run mode error thrown if DEST_CHANNEL_ID not specified with OAUTH_TOKEN",
			args{
				buildInfo: BuildInfo{
					JobName:     "somejob",
					BuildURL:    "someurl",
					BuildStatus: "somestatus",
					OauthToken:  "sometoken",
				},
				slackClient: happySlackClient,
			},
			true,
			PickRunModeErrorMessage,
		},
		{
			"successful post if client is happy and we've filled in all required values and a hookURL",
			args{
				buildInfo: BuildInfo{
					JobName:     "somejob",
					BuildURL:    "someurl",
					BuildStatus: "somestatus",
					HookURL:     "somehookurl",
				},
				slackClient: happySlackClient,
			},
			false,
			"",
		},
		{
			"successful post if client is happy and we've filled in all required values and a destChannelID / Oauth",
			args{
				buildInfo: BuildInfo{
					JobName:       "somejob",
					BuildURL:      "someurl",
					BuildStatus:   "somestatus",
					DestChannelId: "somechannel",
					OauthToken:    "sometoken",
				},
				slackClient: happySlackClient,
			},
			false,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.slackClient.PostToSlack(tt.args.buildInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostToSlack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_getWebhookMessage(t *testing.T) {
	type args struct {
		buildInfo   BuildInfo
		buildStatus Status
	}
	tests := []struct {
		name string
		args args
		want slack.WebhookMessage
	}{
		{"success - required fields only",
			args{
				buildInfo: BuildInfo{
					JobName:     jobName,
					BuildURL:    buildURL,
					BuildStatus: successKey,
				},
				buildStatus: successStatus,
			},
			slack.WebhookMessage{Attachments: []slack.Attachment{
				{
					Title:     fmt.Sprintf("%s: %s", successStatus.text, jobName),
					TitleLink: buildURL,
					Color:     "good",
				}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getWebhookMessage(tt.args.buildInfo, tt.args.buildStatus); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getWebhookMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSpecifiedAttachmentFields(t *testing.T) {
	type args struct {
		buildInfo BuildInfo
	}
	tests := []struct {
		name string
		args args
		want []slack.AttachmentField
	}{
		{"all the fields",
			args{buildInfo: BuildInfo{
				BranchName:  branchName,
				GitCommit:   commit,
				BuildTime:   buildTime,
				TriggeredBy: triggeredBy,
			}},
			[]slack.AttachmentField{
				branchNameField,
				commitField,
				buildTimeField,
				triggeredByField,
			},
		},
		{"branch name",
			args{buildInfo: BuildInfo{BranchName: branchName}},
			[]slack.AttachmentField{branchNameField}},
		{"commit",
			args{buildInfo: BuildInfo{GitCommit: commit}},
			[]slack.AttachmentField{commitField}},
		{"build time",
			args{buildInfo: BuildInfo{BuildTime: buildTime}},
			[]slack.AttachmentField{buildTimeField}},
		{"triggered by",
			args{buildInfo: BuildInfo{TriggeredBy: triggeredBy}},
			[]slack.AttachmentField{triggeredByField}},
		{"no fields in build info",
			args{buildInfo: BuildInfo{}},
			emptyAttachmentFields},
		{"empty value",
			args{buildInfo: BuildInfo{BuildTime: ""}},
			emptyAttachmentFields},
		{"whitespace value",
			args{buildInfo: BuildInfo{BuildTime: "  "}},
			emptyAttachmentFields},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSpecifiedAttachmentFields(tt.args.buildInfo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSpecifiedAttachmentFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAttachment(t *testing.T) {
	type args struct {
		buildInfo   BuildInfo
		buildStatus Status
	}
	tests := []struct {
		name string
		args args
		want slack.Attachment
	}{
		{"success - required fields only",
			args{
				buildInfo: BuildInfo{
					JobName:     jobName,
					BuildURL:    buildURL,
					BuildStatus: successKey,
				},
				buildStatus: successStatus,
			},
			slack.Attachment{
				Title:     fmt.Sprintf("%s: %s", successStatus.text, jobName),
				TitleLink: buildURL,
				Color:     "good",
			}},
		{"success - with attachments",
			args{
				buildInfo: BuildInfo{
					JobName:     jobName,
					BuildURL:    buildURL,
					BuildStatus: successKey,
					BranchName:  branchName,
					GitCommit:   commit,
					BuildTime:   buildTime,
					TriggeredBy: triggeredBy,
				},
				buildStatus: successStatus,
			},
			slack.Attachment{
				Title:     fmt.Sprintf("%s: %s", successStatus.text, jobName),
				TitleLink: buildURL,
				Color:     "good",
				Fields: []slack.AttachmentField{
					branchNameField,
					commitField,
					buildTimeField,
					triggeredByField,
				},
			}},
		{"failure - required fields only",
			args{
				buildInfo: BuildInfo{
					JobName:     jobName,
					BuildURL:    buildURL,
					BuildStatus: "FAILURE",
				},
				buildStatus: failedStatus,
			},
			slack.Attachment{
				Title:     fmt.Sprintf("%s: %s", failedStatus.text, jobName),
				TitleLink: buildURL,
				Color:     "danger",
			}},
		{"unstable - required fields only",
			args{
				buildInfo: BuildInfo{
					JobName:     jobName,
					BuildURL:    buildURL,
					BuildStatus: "UNSTABLE",
				},
				buildStatus: unstableStatus,
			},
			slack.Attachment{
				Title:     fmt.Sprintf("%s: %s", unstableStatus.text, jobName),
				TitleLink: buildURL,
				Color:     "warning",
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAttachment(tt.args.buildInfo, tt.args.buildStatus); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAttachment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPostMessage(t *testing.T) {
	buildInfo := BuildInfo{
		JobName:     "test-job",
		BuildURL:    "https://example.com/build/1",
		BuildStatus: successKey,
		BranchName:  "main",
		GitCommit:   "abc123",
	}

	msgOptions := getPostMessage(buildInfo, successStatus)

	if len(msgOptions) == 0 {
		t.Error("Expected message options to be returned")
	}

	// Verify the structure contains message options (we can't easily inspect MsgOption internals)
	// This ensures the function executes and returns the expected slice structure
}

// fakeSlackAPI is a test double for the slackAPI interface.
type fakeSlackAPI struct {
	capturedChannelID string
	capturedOptions   []slack.MsgOption
	err               error
}

func (f *fakeSlackAPI) PostMessage(channelID string, options ...slack.MsgOption) (string, string, error) {
	f.capturedChannelID = channelID
	f.capturedOptions = options
	return "", "", f.err
}

func Test_productionSlackClientWorker_postChannelMessage(t *testing.T) {
	t.Run("calls PostMessage with correct channelID and non-empty options", func(t *testing.T) {
		fakeAPI := &fakeSlackAPI{}
		worker := &productionSlackClientWorker{
			apiFactory: func(token string) slackAPI { return fakeAPI },
		}
		buildInfo := BuildInfo{
			JobName:       "test-job",
			BuildURL:      "https://example.com",
			BuildStatus:   successKey,
			OauthToken:    "token",
			DestChannelId: "C12345",
		}
		err := worker.postChannelMessage(buildInfo)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if fakeAPI.capturedChannelID != "C12345" {
			t.Errorf("expected channelID %q, got %q", "C12345", fakeAPI.capturedChannelID)
		}
		if len(fakeAPI.capturedOptions) == 0 {
			t.Error("expected non-empty message options")
		}
	})

	t.Run("propagates error from PostMessage", func(t *testing.T) {
		fakeAPI := &fakeSlackAPI{err: fmt.Errorf("api error")}
		worker := &productionSlackClientWorker{
			apiFactory: func(token string) slackAPI { return fakeAPI },
		}
		buildInfo := BuildInfo{
			OauthToken:    "token",
			DestChannelId: "C12345",
			BuildStatus:   successKey,
		}
		err := worker.postChannelMessage(buildInfo)
		if err == nil || err.Error() != "api error" {
			t.Errorf("expected 'api error', got %v", err)
		}
	})
}

func Test_productionSlackClientWorker_postWebhookMessage(t *testing.T) {
	t.Run("calls webhookPoster with correct URL and non-nil message", func(t *testing.T) {
		var capturedURL string
		var capturedMsg *slack.WebhookMessage
		worker := &productionSlackClientWorker{
			webhookPoster: func(url string, msg *slack.WebhookMessage) error {
				capturedURL = url
				capturedMsg = msg
				return nil
			},
		}
		buildInfo := BuildInfo{
			JobName:     "test-job",
			BuildURL:    "https://example.com",
			BuildStatus: successKey,
			HookURL:     "https://hooks.slack.com/test",
		}
		err := worker.postWebhookMessage(buildInfo)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if capturedURL != "https://hooks.slack.com/test" {
			t.Errorf("expected URL %q, got %q", "https://hooks.slack.com/test", capturedURL)
		}
		if capturedMsg == nil {
			t.Error("expected non-nil webhook message")
		}
	})

	t.Run("propagates error from webhookPoster", func(t *testing.T) {
		worker := &productionSlackClientWorker{
			webhookPoster: func(url string, msg *slack.WebhookMessage) error {
				return fmt.Errorf("webhook error")
			},
		}
		buildInfo := BuildInfo{
			HookURL:     "https://hooks.slack.com/test",
			BuildStatus: successKey,
		}
		err := worker.postWebhookMessage(buildInfo)
		if err == nil || err.Error() != "webhook error" {
			t.Errorf("expected 'webhook error', got %v", err)
		}
	})
}

func Test_NewSlackClient(t *testing.T) {
	client := NewSlackClient()
	if client.slackClient == nil {
		t.Error("Expected non-nil slack client")
	}
}

func Test_PostToSlack_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		buildInfo   BuildInfo
		slackClient SlackClient
		wantErr     bool
		errContains string
	}{
		{
			"both webhook and oauth set - should prefer oauth and succeed",
			BuildInfo{
				JobName:       "job",
				BuildURL:      "url",
				BuildStatus:   "SUCCESS",
				HookURL:       "https://hook",
				OauthToken:    "token",
				DestChannelId: "channel",
			},
			NewTestClient(false, false),
			false,
			"",
		},
		{
			"empty strings should be treated as not set",
			BuildInfo{
				JobName:       "job",
				BuildURL:      "url",
				BuildStatus:   "SUCCESS",
				HookURL:       "",
				OauthToken:    "",
				DestChannelId: "",
			},
			NewTestClient(false, false),
			true,
			PickRunModeErrorMessage,
		},
		{
			"only oauth token without channel id should error",
			BuildInfo{
				JobName:     "job",
				BuildURL:    "url",
				BuildStatus: "SUCCESS",
				OauthToken:  "token",
			},
			NewTestClient(false, false),
			true,
			PickRunModeErrorMessage,
		},
		{
			"only channel id without oauth token should error",
			BuildInfo{
				JobName:       "job",
				BuildURL:      "url",
				BuildStatus:   "SUCCESS",
				DestChannelId: "channel",
			},
			NewTestClient(false, false),
			true,
			PickRunModeErrorMessage,
		},
		{
			"channel message client returns error - PostToSlack propagates it",
			BuildInfo{
				JobName:       "job",
				BuildURL:      "url",
				BuildStatus:   "SUCCESS",
				OauthToken:    "token",
				DestChannelId: "channel",
			},
			NewTestClient(true, false),
			true,
			ChannelMessageTestErr,
		},
		{
			"webhook client returns error - PostToSlack propagates it",
			BuildInfo{
				JobName:     "job",
				BuildURL:    "url",
				BuildStatus: "SUCCESS",
				HookURL:     "https://hooks.slack.com/test",
			},
			NewTestClient(false, true),
			true,
			WebhookMessageTestErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.slackClient.PostToSlack(tt.buildInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostToSlack() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.errContains != "" && err != nil && err.Error() != tt.errContains {
				t.Errorf("PostToSlack() error = %v, expected to contain %v", err.Error(), tt.errContains)
			}
		})
	}
}
