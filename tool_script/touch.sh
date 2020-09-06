#!/bin/bash

# 触发构建
# $ ./touch.sh backend BRANCH
# $ ./touch.sh fronend BRANCH

echo '$(GITHUB_WORKFLOW_TOKEN)'
exit 1

GITHUB_WORKFLOW_TOKEN=$GITHUB_WORKFLOW_TOKEN

BACKEND_WORKFLOW_ID=772220
FRONEND_WORKFLOW_ID=772221

TARGET_WORKFLOW_ID=$BACKEND_WORKFLOW_ID

ACTION_TARGET_BACKEND='backend'
ACTION_TARGET_FRONEND='fronend'

ACTION_TARGET_BRANCH=$2

if [ $1 = $ACTION_TARGET_BACKEND ]; then
    TARGET_WORKFLOW_ID=$BACKEND_WORKFLOW_ID
else
    TARGET_WORKFLOW_ID=$FRONEND_WORKFLOW_ID
fi

if [ -z $GITHUB_WORKFLOW_TOKEN ]; then
    echo "err!"
    echo "github token is required"
    exit 1
fi

if [[ -z $ACTION_TARGET_BACKEND || -z $ACTION_TARGET_BRANCH ]]; then
    echo "err!"
    echo "eg: ./touch.sh backend master"
    exit 1
fi

echo "touch $ACTION_TARGET_BACKEND $ACTION_TARGET_BRANCH ..."

curl --location --request POST "https://api.github.com/repos/growerlab/growerlab/actions/workflows/$TARGET_WORKFLOW_ID/dispatches" \
    --header "Accept: application/vnd.github.v3+json" \
    --header "Authorization: token $GITHUB_WORKFLOW_TOKEN" \
    --header "Content-Type: application/json" \
    --header "Cookie: logged_in=no" \
    --data-raw "{\"ref\":\"${ACTION_TARGET_BRANCH}\"}"

if test $? -eq 0; then
    echo "done"
fi
