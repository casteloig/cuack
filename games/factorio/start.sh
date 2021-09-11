#!/bin/bash
sed -i "s/^; max-players=.*/; max-players=$PLAYERS/" /opt/factorio/config/config.ini

/opt/factorio/bin/x64/factorio --create /factorio/saves/"$WORLD_NAME".zip
/opt/factorio/bin/x64/factorio --start-server savegame.zip >> /mnt/cuack/logs
