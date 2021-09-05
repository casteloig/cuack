#!/bin/bash
sed -i "s/^max-players=.*/max-players=$PLAYERS/" server.properties
sed -i "s/^difficulty=.*/difficulty=$DIFFICULTY/" server.properties

while true
do
    java -Xmx1024M -Xms1024M -jar server.jar nogui >> /mnt/cuack/logs
    echo "restarting in 5"
    sleep 5
done