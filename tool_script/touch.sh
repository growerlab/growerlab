#!/bin/bash

# 触发构建
# $ ./touch.sh backend BRANCH
# $ ./touch.sh fronend BRANCH

GITHUB_WORKFLOW_TOKEN=$GITHUB_WORKFLOW_TOKEN

BACKEND_WORKFLOW_ID=772220
FRONEND_WORKFLOW_ID=772221

TARGET_WORKFLOW_ID=$BACKEND_WORKFLOW_ID

ACTION_TARGET='backend'

if [ $1 = $ACTION_TARGET ]; then
    TARGET_WORKFLOW_ID=$BACKEND_WORKFLOW_ID
else
    TARGET_WORKFLOW_ID=$FRONEND_WORKFLOW_ID
fi

curl --location --request POST "https://api.github.com/repos/growerlab/growerlab/actions/workflows/$TARGET_WORKFLOW_ID/dispatches" \
    --header "Accept: application/vnd.github.v3+json" \
    --header "Authorization: token $GITHUB_WORKFLOW_TOKEN" \
    --header "Content-Type: application/json" \
    --header "Cookie: logged_in=no" \
    --data-raw '{
    "ref":"master"
}'

if test $? -eq 0; then
    echo "done"
fi
