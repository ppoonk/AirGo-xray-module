#!/system/bin/sh

MODDIR=${0%/*}
chmod 777 -R $MODDIR/

while true; do
    # 判断是否已经获取到网卡
    if [[ ! ${net_card} ]]; then
        # 获取流量网卡
        net_card=`dumpsys connectivity | grep "NetworkAgentInfo{" | grep "type: MOBILE" | grep -iv "extra: ims" | grep -o "InterfaceName: [^ ]*" | grep -o "[^: ]*$"`
        net_card=${net_card:-`dumpsys netstats | grep -i "iface=" | grep "metered=true" | grep -Eio -m 1 "rmnet[^ ]*|ccmni[^ ]*"`}
        net_card=${net_card:-`ip route | grep -Eio -m 1 "rmnet[^ ]*|ccmni[^ ]*"`}
        net_card=${net_card:-`ifconfig | grep -Ei -A 1 "rmnet[^ ]*|ccmni[^ ]*" | grep -w -B 1 "inet" | grep -Eo "rmnet[^ ]*|ccmni[^ ]*"`}
    fi

    # 判断是否能获取到内网IP
    [[ ${net_card} && `ifconfig | grep -w -A 1 "${net_card}" | grep -o "inet addr:[^ ]*" | grep -o "[^: ]*$"` ]] && break || sleep 3
done

sleep 10
/system/bin/sh $MODDIR/AirGo/start.sh
