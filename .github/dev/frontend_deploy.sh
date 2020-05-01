#!/bin/bash
set -eu

# server
SERVER_HOST=$SERVER
SERVER_PORT=$SERVER_PORT
SERVER_USER=$SERVER_USER
BRANCH="$(cat "$BRANCH_FILE")"
SERVICES_PATH="/data/$BRANCH"
DEPLOY_KEY=$SERVER_SSH_KEY

# PWD
ROOT_DIR=$GITHUB_WORKSPACE

FRONTEND="frontend"

cloneFrontend() {
    echo "------ clone and build $FRONTEND -------"
    cd $ROOT_DIR || exit 1

    git clone "https://github.com/growerlab/$FRONTEND.git" --depth=1 --branch=$BRANCH
    if test $? -ne 0; then
        git clone "https://github.com/growerlab/$FRONTEND.git" --depth=1 --branch=master
    fi

    cd -
}

buildFrontend() {
    echo "------ buildFrontend -------"
    cd $FRONTEND || exit 1

    npm install
    npm run build --if-present

    cd -
    echo "------ buildFrontend done -------"
}

syncDist() {
    echo "------ syncDist -------"
    cd $FRONTEND || exit 1

    SSHPATH="$HOME/.ssh"
    mkdir -p "$SSHPATH"
    echo "$DEPLOY_KEY" >"$SSHPATH/key"
    chmod 600 "$SSHPATH/key"
    SERVER_DEPLOY_STRING="$SERVER_USER@$SERVER_HOST:$SERVICES_PATH"

    rsync -e -c -r --delete -e "ssh -i $SSHPATH/key -o StrictHostKeyChecking=no -p $SERVER_PORT" ./dist/ "$SERVER_DEPLOY_STRING"/website

    cd -
    echo "------ syncDist done -------"
}

main() {
    cloneFrontend
    buildFrontend
    syncDist
}

main
