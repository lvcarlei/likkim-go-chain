#!/bin/bash  
  
# 设置容器名称  
CONTAINER_NAME="wallet_api"  
# 设置容器启动命令（根据需要修改）  
DOCKER_RUN_COMMAND="docker run -d --network host --name $CONTAINER_NAME goapp:$CONTAINER_NAME"  
docker load -i /home/java/$CONTAINER_NAME.tar 
# 检查容器是否存在  
if docker ps -a --filter "name=$CONTAINER_NAME" -q | grep -q .; then  
    # 容器存在，先停止并删除容器  
    echo "Container $CONTAINER_NAME exists, stopping and removing..."  
    docker stop $CONTAINER_NAME  
    docker rm $CONTAINER_NAME  
    # 删除后重新启动容器  
    echo "Starting new container $CONTAINER_NAME..."  
    $DOCKER_RUN_COMMAND  
else  
    # 容器不存在，直接启动新容器  
    echo "Container $CONTAINER_NAME does not exist, starting new container..."  
    $DOCKER_RUN_COMMAND  
fi
