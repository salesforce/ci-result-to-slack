/*
 * Copyright (c) 2021, salesforce.com, inc.
 * All rights reserved.
 * SPDX-License-Identifier: BSD-3-Clause
 * For full license text, see the LICENSE file in the repo root or https://opensource.org/licenses/BSD-3-Clause
 */
package internal

import (
	"errors"
	"fmt"
	"github.com/slack-go/slack"
	"strings"
)

const ChannelMessageTestErr = "error from postChannelMessage"
const WebhookMessageTestErr = "error from postWebhookMessage"

type SlackClient struct {
	slackClient
}

type slackClient interface {
	postChannelMessage(buildInfo BuildInfo) error
	postWebhookMessage(buildInfo BuildInfo) error
}

type productionSlackClientWorker struct {
}

func (client *productionSlackClientWorker) postChannelMessage(buildInfo BuildInfo) error {
	api := slack.New(buildInfo.OauthToken)
	postMessage := getPostMessage(buildInfo, buildInfo.GetContextualStatus())
	_, _, err := api.PostMessage(buildInfo.DestChannelId, postMessage...)
	return err
}

func (client *productionSlackClientWorker) postWebhookMessage(buildInfo BuildInfo) error {
	message := getWebhookMessage(buildInfo, buildInfo.GetContextualStatus())
	err := slack.PostWebhook(buildInfo.HookURL, &message)
	return err
}

type testSlackClientWorker struct {
	postChannelMessageShouldError bool
	postWebhookMessageShouldError bool
}

func (client *testSlackClientWorker) postChannelMessage(buildInfo BuildInfo) error {
	if client.postChannelMessageShouldError {
		return errors.New(ChannelMessageTestErr)
	}
	return nil
}

func (client *testSlackClientWorker) postWebhookMessage(buildInfo BuildInfo) error {
	if client.postWebhookMessageShouldError {
		return errors.New(WebhookMessageTestErr)
	}
	return nil
}

func NewTestClient(postChannelMessageShouldError bool, postWebhookMessageShouldError bool) SlackClient {
	return SlackClient{
		&testSlackClientWorker{
			postWebhookMessageShouldError: postWebhookMessageShouldError,
			postChannelMessageShouldError: postChannelMessageShouldError,
		},
	}
}

func (client *SlackClient) PostToSlack(buildInfo BuildInfo) error {
	if buildInfo.OauthToken != "" && buildInfo.DestChannelId != "" {
		err := client.postChannelMessage(buildInfo)
		if err != nil {
			return err
		}
	} else if buildInfo.HookURL != "" {
		err := client.postWebhookMessage(buildInfo)
		if err != nil {
			return err
		}
	} else {
		return errors.New(PickRunModeErrorMessage)
	}
	return nil
}

func NewSlackClient() SlackClient {
	return SlackClient{&productionSlackClientWorker{}}
}

func getPostMessage(buildInfo BuildInfo, buildStatus Status) []slack.MsgOption {
	msgOptions := []slack.MsgOption{
		slack.MsgOptionAttachments(getAttachment(buildInfo, buildStatus)),
	}
	return msgOptions
}

func getWebhookMessage(buildInfo BuildInfo, buildStatus Status) slack.WebhookMessage {
	attachment := getAttachment(buildInfo, buildStatus)
	message := slack.WebhookMessage{Attachments: []slack.Attachment{attachment}}
	return message
}

func getAttachment(buildInfo BuildInfo, buildStatus Status) slack.Attachment {
	attachment := slack.Attachment{
		Title:     fmt.Sprintf("%s: %s", buildStatus.text, buildInfo.JobName),
		TitleLink: buildInfo.BuildURL,
		Color:     buildStatus.color,
		Fields:    getSpecifiedAttachmentFields(buildInfo),
	}
	return attachment
}

func getSpecifiedAttachmentFields(buildInfo BuildInfo) []slack.AttachmentField {
	var attachmentFields []slack.AttachmentField

	appendAttachmentField(&attachmentFields, branchFieldTitle, buildInfo.BranchName)
	appendAttachmentField(&attachmentFields, commitFieldTitle, buildInfo.GitCommit)
	appendAttachmentField(&attachmentFields, buildTimeFieldTitle, buildInfo.BuildTime)
	appendAttachmentField(&attachmentFields, triggeredByFieldTitle, buildInfo.TriggeredBy)
	return attachmentFields
}

func appendAttachmentField(attachmentFields *[]slack.AttachmentField, fieldTitle string, fieldValue string) {
	if strings.TrimSpace(fieldValue) != "" {
		*attachmentFields = append(*attachmentFields, getAttachmentField(fieldTitle, fieldValue))
	}
}

func getAttachmentField(title string, value string) slack.AttachmentField {
	return slack.AttachmentField{
		Title: title,
		Value: value,
		Short: true,
	}
}
