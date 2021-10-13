#!/usr/bin/env bash
set -x
set -e

TAG_NAME=$1
if [[ -z ${TAG_NAME} ]]; then
    echo "第一个参数必须指定tag版本号"
    exit0
fi
echo 使用 ${TAG_NAME} 版本运行

NAME=file-server
DOCKER_USER='robot$biomind_robot'
DOCKER_PASSWORD=zjyNJBrmerOIWgtai9wtjxLRh8djdJVa
DOCKER_HUB=smart.hub.biomind.com.cn

mkdir -p /data/${NAME}/data:

docker login  --username "${DOCKER_USER}" --password "${DOCKER_PASSWORD}" ${DOCKER_HUB}

docker ps -a|grep ${NAME}|awk '{print $1}'|xargs -r docker rm -f

docker run --name ${NAME} --network host -e MODE=test -v /data/${NAME}/data:/home/data -v /data/${NAME}/log:/home/log -e MODE=test -d ${DOCKER_HUB}/intelligent-system/${NAME}:${TAG_NAME}