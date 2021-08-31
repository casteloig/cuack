#!/bin/bash
apt update
apt install -y python3 python3-pip

mkdir -p /root/logs
chmod 777 /root/logs

pip3 install docker