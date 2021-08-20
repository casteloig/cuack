#!/bin/bash
while true
do
    java -Xmx1024M -Xms1024M -jar server.jar nogui
    echo "restarting in 10"
    sleep 10
done