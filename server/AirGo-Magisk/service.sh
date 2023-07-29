#!/system/bin/sh

MODDIR=${0%/*}
chmod 777 -R $MODDIR/
/system/bin/sh $MODDIR/AirGo/start.sh
