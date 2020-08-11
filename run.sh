#!/bin/bash

if [[ ! -d "/root/.pocd/config" ]]; then
  echo "Initialize /root/.pocdd"
  cp -r /go/initial-node/* /root/.pocd/
fi

cd /go; ./pocd start $@
