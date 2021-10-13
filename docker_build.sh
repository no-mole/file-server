#!/usr/bin/env bash
set -x
set -e

TAG_NAME=$1
if [[ -z ${TAG_NAME} ]]; then
    echo "第一个参数必须指定tag版本号"
    exit
fi
echo 使用 ${TAG_NAME} 版本构建镜像

NAME=file-server
PROJECT=intelligent-system
IMAGE_ORIGIN=smart.hub.biomind.com.cn

if [[ $NAME == "biogo-example" ]]; then
  echo "Default project name is not allowed，please modify it"
  exit
fi

IMAGE_URI=${IMAGE_ORIGIN}/${PROJECT}/${NAME}
TAG_IMAGE_URI=${IMAGE_URI}:${TAG_NAME}
docker build --no-cache --compress -f Dockerfile -t ${TAG_IMAGE_URI} --build-arg TAG_NAME=${TAG_NAME} .

DOCKER_USER='robot$biomind_robot'
DOCKER_PASSWORD=zjyNJBrmerOIWgtai9wtjxLRh8djdJVa
docker login  --username "${DOCKER_USER}" --password "${DOCKER_PASSWORD}" ${IMAGE_ORIGIN}
docker push  ${TAG_IMAGE_URI}

