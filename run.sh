#!/bin/sh

# パラメータのデフォルト値
DBPORT=${1:-9998}
DTAGO=${2:-15m}
DTFROMNOW=${3:-15m}

# cd /home/chouette/MyProject/Showroom/AddEventuser

export SR_ADD_EVENTUSER_DTAGO=$DTAGO
export SR_ADD_EVENTUSER_DTFROMNOW=$DTFROMNOW
export SR_ADD_EVENTUSER_NOROOMS=40

xport SOPS_AGE_KEY_FILE=/home/chouette/.config/age/key2.txt

./AddEventUser
