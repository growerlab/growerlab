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
BRANCH=$BUILD_BRANCH
SERVICES_PATH="/data/$BRANCH"
DATABASE_NAME="growerlab_$BRANCH"
DEPLOY_KEY=$SERVER_SSH_KEY
POSTGRES_PASSWORD=$SERVER_POSTGRES_PASSWORD
DOMAIN="$BRANCH.dev.growerlab.net"

# PWD
ROOT_DIR=$GITHUB_WORKSPACE

cloneAndBuildProject() {
  repoName=$1
  echo "------ clone and build $repoName -------"

  git clone "https://github.com/growerlab/$repoName.git" --depth=1 --branch=$BRANCH
  if test $? -ne 0; then
    git clone "https://github.com/growerlab/$repoName.git" --depth=1 --branch=master
  fi

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

  mkdir -p "$servicesDir"
  cp -R "$BACKEND/db" "$servicesDir"
  cp -R "$BACKEND/conf" "$servicesDir"
  cp "$BACKEND/$BACKEND" "$servicesDir"

  sed -i "s/namespace: master/namespace: $BRANCH/g" "$servicesDir/conf/config.yaml"
  sed -i "s/postgresql:.*/postgresql:\/\/growerlab:$POSTGRES_PASSWORD@postgres:5432\/$DATABASE_NAME?sslmode=disable/g" "$servicesDir/conf/config.yaml"
  sed -i "s/host: 127.0.0.1/host: keydb/g" "$servicesDir/conf/config.yaml"
  sed -i "s/localhost/$DOMAIN/g" "$servicesDir/conf/config.yaml"

  echo "------ done backend -------"
}

bindSVC() {
  echo "------ bind svc -------"
  cd "$ROOT_DIR" || exit 1
  servicesDir="$ROOT_DIR/data/services/$SVC"

  mkdir -p "$servicesDir"
  cp -R "$SVC/.env.example" "$servicesDir"/.env
  cp -R "$SVC/template" "$servicesDir"/
  cp "$SVC/$SVC" "$servicesDir"

  # ESCAPE_SERVICES_PATH=$(echo "${SERVICES_PATH}" | sed 's/\//\\\//g')
  sed -i "s/logs\//\/data\/logs\/svc-/g" $servicesDir/.env
  sed -i "s/REPOS_DIR=repos\//REPOS_DIR=\/data\/repositories\//g" $servicesDir/.env

  echo "------ done svc -------"
}

bindMensa() {
  echo "------ bind mensa -------"
  cd "$ROOT_DIR" || exit 1
  servicesDir="$ROOT_DIR/data/services/$MENSA"

  mkdir -p "$servicesDir"
  cp -R "$MENSA/conf" "$servicesDir"
  cp "$MENSA/$MENSA" "$servicesDir"

  sed -i "s/growerlab:growerlab@127.0.0.1:5432\/growerlab/growerlab:$POSTGRES_PASSWORD@postgres:5432\/$DATABASE_NAME/g" $servicesDir/conf/config.yaml
  sed -i "s/host: 127.0.0.1/host: keydb/g" $servicesDir/conf/config.yaml
  sed -i "s/master/${BRANCH}/g" $servicesDir/conf/config.yaml
  sed -i "s/test\/repos/\/data\/repositories/g" $servicesDir/conf/config.yaml

  echo "------ done mensa -------"
}

# 将本仓库rsync到服务器
syncData() {
  echo "------ deploy -------"
  cd "$ROOT_DIR" || exit 1
  SSHPATH="$HOME/.ssh"
  mkdir -p "$SSHPATH"
  echo "$DEPLOY_KEY" >"$SSHPATH/key"
  chmod 600 "$SSHPATH/key"
  SERVER_DEPLOY_STRING="$SERVER_USER@$SERVER_HOST:$SERVICES_PATH"

  if [ ! -d "$ROOT_DIR"/data/pgdata ]; then
    mkdir "$ROOT_DIR"/data/pgdata || echo "pgdata exists"
  fi

  echo "rsync /data/keydb..."
  rsync -e -c -r -u --ignore-errors -e "ssh -i $SSHPATH/key -o StrictHostKeyChecking=no -p $SERVER_PORT" "$ROOT_DIR"/data/keydb "$SERVER_DEPLOY_STRING"/data || $(case "$?" in 0 | 3 | 23) exit 0 ;; *) exit $? ;; esac)
  echo "/rsync /data/pgdata"
  rsync -e -c -r -u --ignore-errors -e "ssh -i $SSHPATH/key -o StrictHostKeyChecking=no -p $SERVER_PORT" "$ROOT_DIR"/data/pgdata "$SERVER_DEPLOY_STRING"/data || $(case "$?" in 0 | 3 | 23) exit 0 ;; *) exit $? ;; esac)
  echo "/rsync /data/logs"
  rsync -e -c -r -u -e "ssh -i $SSHPATH/key -o StrictHostKeyChecking=no -p $SERVER_PORT" "$ROOT_DIR"/data/logs "$SERVER_DEPLOY_STRING"/data
  echo "/rsync /data/repositories"
  rsync -e -c -r -u -e "ssh -i $SSHPATH/key -o StrictHostKeyChecking=no -p $SERVER_PORT" "$ROOT_DIR"/data/repositories "$SERVER_DEPLOY_STRING"/data
  echo "/rsync /data/nginx"
  rsync -e -c -r --delete -e "ssh -i $SSHPATH/key -o StrictHostKeyChecking=no -p $SERVER_PORT" "$ROOT_DIR"/data/nginx "$SERVER_DEPLOY_STRING"/data
  echo "/rsync /data/supervisor"
  rsync -e -c -r --delete -e "ssh -i $SSHPATH/key -o StrictHostKeyChecking=no -p $SERVER_PORT" "$ROOT_DIR"/data/supervisor "$SERVER_DEPLOY_STRING"/data
  echo "/rsync /data/website"
  rsync -e -c -r -u --ignore-errors -e "ssh -i $SSHPATH/key -o StrictHostKeyChecking=no -p $SERVER_PORT" "$ROOT_DIR"/data/website "$SERVER_DEPLOY_STRING"/data || $(case "$?" in 0 | 3 | 23) exit 0 ;; *) exit $? ;; esac)
  echo "/rsync /data/services"
  rsync -e -c -r --delete --ignore-errors -e "ssh -i $SSHPATH/key -o StrictHostKeyChecking=no -p $SERVER_PORT" "$ROOT_DIR"/data/services "$SERVER_DEPLOY_STRING"/data
  echo "/rsync /data/router"
  rsync -e -c -r --delete -e "ssh -i $SSHPATH/key -o StrictHostKeyChecking=no -p $SERVER_PORT" "$ROOT_DIR"/router "$SERVER_DEPLOY_STRING"
  echo "/rsync /data/dev.compose.yaml"
  rsync -e -c -r --delete -e "ssh -i $SSHPATH/key -o StrictHostKeyChecking=no -p $SERVER_PORT" "$ROOT_DIR"/dev.compose.yaml "$SERVER_DEPLOY_STRING"
  echo "/rsync /data/Dockerfile"
  rsync -e -c -r --delete -e "ssh -i $SSHPATH/key -o StrictHostKeyChecking=no -p $SERVER_PORT" "$ROOT_DIR"/Dockerfile "$SERVER_DEPLOY_STRING"
  echo "------ done deploy -------"
}

# 重启docker
restartServices() {
  echo "------- restart services -------"
  cd "$ROOT_DIR" || exit 1
  SSHPATH="$HOME/.ssh"

  DB_SEED=$(cat "$ROOT_DIR/data/services/$BACKEND/db/seed.sql")
  DB_STRUCTURE=$(cat "$ROOT_DIR/data/services/$BACKEND/db/growerlab.sql")

  (
    cat <<EOF
cd $SERVICES_PATH || exit 1

# 写环境变量
cat > .env <<-EOENV
POSTGRES_PASSWORD=$POSTGRES_PASSWORD
EOENV

# build router
./router/build.sh

# docker-compose 编排
echo "docker compose..."
runOrRestartContainer() {
    name=\$1
    alias=\$2
    service_name=\$name
    if [[ \$alias != "" ]]; then
      service_name=\$alias
    fi

    if docker ps --format "{{.Names}}" | grep -qw \$name ; then
      echo "\$service_name 已启动，重启中.."
      docker-compose -f ./dev.compose.yaml restart \$service_name
    else
      echo "\$service_name 未启动，启动中.."
      docker-compose -f ./dev.compose.yaml up -d \$service_name
    fi
}

runOrRestartContainer "postgres"
runOrRestartContainer "keydb"
runOrRestartContainer "nginx"
# waiting for services
sleep 2
runOrRestartContainer "router"
runOrRestartContainer "services_$BRANCH" "growerlab"

# init database
echo "init database..."
docker exec -i postgres /bin/bash <<-EODOCKER
  if ! psql --username growerlab -lqt | cut -d \| -f 1 | grep -qw $DATABASE_NAME; then
    psql -v ON_ERROR_STOP=1 --username growerlab --dbname $DATABASE_NAME <<-EOSQL
      create database $DATABASE_NAME;
      grant all privileges on database $DATABASE_NAME to growerlab;
      ${DB_STRUCTURE}
      ${DB_SEED}
EOSQL
  fi
EODOCKER

EOF
  ) >"$HOME"/growerlab.sh

  sh -c "ssh -i $SSHPATH/key -o StrictHostKeyChecking=no -p $SERVER_PORT $SERVER_USER@$SERVER_HOST < $HOME/growerlab.sh"
}

main() {
  cloneAndBuildProject $SVC
  cloneAndBuildProject $BACKEND
  cloneAndBuildProject $MENSA

  bindBackend
  bindSVC
  bindMensa

  syncData
  restartServices
}

main
