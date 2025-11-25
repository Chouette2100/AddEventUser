#!/bin/sh
tar czvf AddEventUser_$(date +%Y%m%d-%H%M%S).tar.gz \
AddEventUser \
DBConfig.yml \
Env.yml \
tar.sh
