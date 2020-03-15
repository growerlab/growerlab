#!/bin/bash

# 编译golang相关的项目，并部署到dev服务器

SVC="svc"
BACKEND="backend"
MENSA="mensa"

# server
SERVER=$GOWERLAB_SERVER
SERVER_PORT=$GOWERLAB_SERVER_PORT
SERVICES_PATH="/data/$DEPLOY_DIR/"

# PWD
ROOT_DIR=pwd

cloneAndBuildProject() {
  repoName=$1
  echo "------ clone and build $repoName -------"

  git clone "https://github.com/growerlab/$repoName.git" --depth=1
  cd "$repoName" || exit 1
  go get -v -t -d ./...
  go build -o "$repoName" -v main.go
  cd ..
}

# 将二进制文件移动到本仓库的data/services目录下
bindBackend() {
  echo "------ bind backend -------"
  cd $ROOT_DIR || exit 1
  servicesDir="$ROOT_DIR/data/services/$BACKEND"
  sourceDir="$ROOT_DIR/$BACKEND"

  mkdir $sourceDir
  cp -R "$BACKEND/conf" $servicesDir
  cp "$BACKEND/$BACKEND" $servicesDir
}

bindSVC() {
  echo "------ bind svc -------"
  cd $ROOT_DIR || exit 1
  servicesDir="$ROOT_DIR/data/services/$SVC"
  sourceDir="$ROOT_DIR/$SVC"

  mkdir $sourceDir
  cp -R "$SVC/.env.example" $servicesDir/.env
  cp "$SVC/$SVC" $servicesDir
}

bindMensa() {
  echo "------ bind mensa -------"
  cd $ROOT_DIR || exit 1
  servicesDir="$ROOT_DIR/data/services/$MENSA"
  sourceDir="$ROOT_DIR/$MENSA"

  mkdir $sourceDir
  cp -R "$MENSA/conf" $servicesDir
  cp "$MENSA/$MENSA" $servicesDir
}

# 将本仓库rsync到服务器
deploy() {
  echo "------ deploy -------"

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

  deploy
  restartService
}

main
