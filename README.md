# CI Results to Slack

This project has been built as a somewhat CI agnostic approach to posting 
build results to Slack. Some of the plugins for CI servers work great in an environment
where you don't have proxy requirements but not so great when some of the tradeoffs are that
you have to open the proxies up for the whole CI system (Jenkins). Granted, there are knobs
you can pull, but why not also run this on an agent (again in the Jenkins case).

# Approach
`main.go` accepts inputs via environment variables. The idea is to potentially
run this on a host or within a Docker container.

# Building
`make local-docker-build` will do all you need but there are other make targets available.

# Usage
The following environment variables can be used. You *MUST* specify either `HOOK_URL` for incoming webhook integration 
or both `OAUTH_TOKEN` and `DEST_CHANNEL_ID` for app integration which calls the Slack APIs (more flexible):
```
KEY                  TYPE             DEFAULT    REQUIRED    DESCRIPTION
JOB_NAME             String                      true        Name of the build's job
BUILD_URL            String                      true        Direct URL to the build
BUILD_STATUS         String                      true        Status of build (e.g. currentBuild.currentResult in Jenkins)
HOOK_URL             String                                  Slack Webhook URL set via Incoming Webhooks
DEST_CHANNEL_ID      String                                  Destination Channel ID (not the name of the channel)
OAUTH_TOKEN          String                                  OAuth Token used to send message via app
LAST_BUILD_STATUS    String           UNKNOWN                Status of last build used to provide contextual build Status
BRANCH_NAME          String                                  Name of git branch
GIT_COMMIT           String                                  Git commit hash
BUILD_TIME           String                                  Build time (e.g. durationString in Jenkins)
TRIGGERED_BY         String                                  The action which triggered the build
SKIP_IF_SUCCESS      True or False                           Skip posting if contextual Status is success
```

## Example
`docker run --rm=true -e OAUTH_TOKEN -e JOB_NAME -e BUILD_URL -e BUILD_STATUS -e DEST_CHANNEL_ID -e TRIGGERED_BY -e SKIP_IF_SUCCESS -e BUILD_TIME -e LAST_BUILD_STATUS -e BRANCH_NAME ci-result-to-slack`

# Setup

## Slack Bot

### OAUTH_TOKEN
To use a Slack Bot, you'll get the `OAUTH_TOKEN` by following 
[Slack's Bot instructions](https://api.slack.com/bot-users#getting-started). You'll want to give the bot access to 
[chat.postMessage](https://api.slack.com/methods/chat.postMessage) by giving it the 
[chat:write:bot scope](https://api.slack.com/scopes/chat:write:bot).

### DEST_CHANNEL_ID
**NOTE: This is not the name of the channel**
To get the `DEST_CHANNEL_ID` (destination channel ID), you can do the following:

#### In a Browser
The Channel ID is in the URL at the end: `https://${YOUR_WORKSPACE}.slack.com/messages/${DEST_CHANNEL_ID}/`

#### In the Slack App
If you right-click on the channel / messages, you can hit `Copy Link` and also see the `DEST_CHANNEL_ID` at the end
(e.g. `https://${YOUR_WORKSPACE}.slack.com/messages/${DEST_CHANNEL_ID}`)

## Slack Incoming Webhooks
To get the `HOOK_URL` for Incoming Webhooks, you'll follow 
[Slack's Incoming Webhooks instructions](https://api.slack.com/incoming-webhooks).
