#!/bin/bash
set -eu

# 编译golang相关的项目，并部署到dev服务器

SVC="svc"
BACKEND="backend"
MENSA="mensa"

# server
SERVER_HOST=$SERVER
SERVER_PORT=$SERVER_PORT
SERVER_USER=$SERVER_USER
SERVICES_PATH="/data/$(cat "$DEPLOY_DIR")"
DEPLOY_KEY=$SERVER_SSH_KEY

# PWD
ROOT_DIR=$GITHUB_WORKSPACE

cloneAndBuildProject() {
  repoName=$1
  echo "------ clone and build $repoName -------"

  git clone "https://github.com/growerlab/$repoName.git" --depth=1
  cd "$repoName" || exit 1
  go get -v -t -d ./...
  go build -o "$repoName" -v main.go
  echo "build $repoName => $(pwd)/$repoName"
  cd -
}

# 将二进制文件移动到本仓库的data/services目录下
bindBackend() {
  echo "------ bind backend -------"
  cd "$ROOT_DIR" || exit 1
  servicesDir="$ROOT_DIR/data/services/$BACKEND"
  sourceDir="$ROOT_DIR/$BACKEND"

  mkdir -p "$servicesDir"
  cp -R "$BACKEND/conf" "$servicesDir"
  cp "$BACKEND/$BACKEND" "$servicesDir"
  echo "------ done backend -------"
}

bindSVC() {
  echo "------ bind svc -------"
  cd "$ROOT_DIR" || exit 1
  servicesDir="$ROOT_DIR/data/services/$SVC"
  sourceDir="$ROOT_DIR/$SVC"

  mkdir -p "$servicesDir"
  cp -R "$SVC/.env.example" "$servicesDir"/.env
  cp "$SVC/$SVC" "$servicesDir"
  echo "------ done svc -------"
}

bindMensa() {
  echo "------ bind mensa -------"
  cd "$ROOT_DIR" || exit 1
  servicesDir="$ROOT_DIR/data/services/$MENSA"
  sourceDir="$ROOT_DIR/$MENSA"

  mkdir -p "$servicesDir"

  echo "---------test"
  stat "$MENSA/conf"
  stat "$MENSA/$MENSA"
  echo "---------test"

  cp -R "$MENSA/conf" "$servicesDir"
  cp "$MENSA/$MENSA" "$servicesDir"
  echo "------ done mensa -------"
}

# 将本仓库rsync到服务器
syncData() {
  echo "------ deploy -------"
  cd "$ROOT_DIR" || exit 1
  SSHPATH="$HOME/.ssh"
  mkdir -p "$SSHPATH"
  echo "$DEPLOY_KEY" > "$SSHPATH/key"
  chmod 600 "$SSHPATH/key"
  SERVER_DEPLOY_STRING="$SERVER_USER@$SERVER_HOST:$SERVICES_PATH"

  rsync -avzP --delete -e "ssh -i $SSHPATH/key -o StrictHostKeyChecking=no -p $SERVER_PORT" $ROOT_DIR/data $SERVER_DEPLOY_STRING
  rsync -avzP --delete -e "ssh -i $SSHPATH/key -o StrictHostKeyChecking=no -p $SERVER_PORT" $ROOT_DIR/docker-compose.prod.yaml $SERVER_DEPLOY_STRING
  echo "------ done deploy -------"
}

# 重启docker
restartService() {
  echo "------- restart service -------"
}

main() {
  cloneAndBuildProject $SVC
  cloneAndBuildProject $BACKEND
  cloneAndBuildProject $MENSA

  bindBackend
  bindSVC
  bindMensa

  syncData
  restartService
}

main
