#!/bin/bash

ip=`ifconfig -a|grep inet|grep -v 127.0.0.1|grep -v inet6|grep -v 172.\*.\*.\*|awk '{print $2}'|tr -d "addr:"`
echo $ip

cp conf/config.yaml config.yaml
sed -i -e "s/127.0.0.1/$ip/g" config.yaml
docker stop fiber
docker rm fiber
docker rmi fiber

docker build --progress=plain -t fiber .
docker run --name fiber -d -v ./logs/:/app/logs/ -v ./uploads:/app/uploads -p 3100:3100 fiber
