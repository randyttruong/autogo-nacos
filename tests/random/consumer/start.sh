#!/bin/sh
set -e 

cd /home

logDir="logDir: /tmp/log"
cacheDir="cacheDir: /tmp/cache"

for i in $(seq 1 100)
do
    sed -i "/logDir:/c\"  $logDir$i" ./configs/consumer/config.yaml
    sed -i "/cacheDir:/c\  $cacheDir$i" ./configs/consumer/config.yaml
    ./consumer &
done

while true; do
    sleep 60
done

