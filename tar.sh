#!/bin/sh
tar czvf AddEventUser_$(date +%Y%m%d-%H%M%S).tar.gz \
AddEventUser \
*.enc.yml \
Env.yml \
run.sh \
tar.sh
