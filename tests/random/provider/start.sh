#!/bin/sh
set -e 

cd /home

logDir="logDir: /tmp/log"
cacheDir="cacheDir: /tmp/cache"

for i in $(seq 1 100)
do
    sed -i "/logDir:/c\"  $logDir$i" ./configs/provider/config.yaml
    sed -i "/cacheDir:/c\  $cacheDir$i" ./configs/provider/config.yaml
    ./provider &
done

while true; do
    sleep 60
done

