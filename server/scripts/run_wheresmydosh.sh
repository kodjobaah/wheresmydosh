#!/usr/bin/env bash

nohup /usr/bin/wheresmydosh > /var/log/wheresmydosh/wheresmydosh.log 2>&1 &
echo $! >  /var/run/wheresmydosh.pid
