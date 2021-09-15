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
	"log"
)

const skippedPostingMessage = "Skipped posting to Slack"
const messageSentTemplate = "Message successfully sent to channel for %s"

func handleRequest(slackClient internal.SlackClient) (string, error) {
	buildInfo, err := internal.GetBuildInfoFromEnv()
	if err != nil {
		return "", err
	}
	if buildInfo.ShouldSkipPosting() {
		return skippedPostingMessage, nil
	}
	err = slackClient.PostToSlack(buildInfo)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(messageSentTemplate, buildInfo.JobName), nil
}

/**
If HTTP_PROXY / HTTPS_PROXY is present then the framework will use the proxy
*/
func main() {
	client := internal.NewSlackClient()
	message, err := handleRequest(client)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(message)
}
