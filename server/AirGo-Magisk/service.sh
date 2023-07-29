#!/system/bin/sh
# 不要假设您的模块将位于何处。
# 如果您需要知道此脚本和模块的放置位置，请使用$MODDIR
# 这将确保您的模块仍能正常工作
# 即使Magisk将来更改其挂载点
MODDIR=${0%/*}
chmod 777 $MODDIR/AirGo/bin/*
while [ "$(getprop sys.boot_completed)" != "1" ]; do
        sleep 1
    done

echo "PowerManagerService.noSuspend" > /sys/power/wake_lock

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

sleep 6

${MODDIR}/AirGo/bin/busybox setuidgid 0:3333 ${MODDIR}/AirGo/bin/airgo -start >/dev/null 2>&1 &
