#!/bin/bash
while true
do
    java -Xmx1024M -Xms1024M -jar server.jar nogui >> mine.log
    echo "restarting in 5"
    sleep 5
done