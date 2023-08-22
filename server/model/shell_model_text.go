package model

// 获取安卓版本
const GetAndroidVersion = "getprop ro.build.version.release"

// 获取全部包名
const allPackages = "pm list packages"

// 默认放行的应用
const defaultAllowPackages = `com.android.phone`

// 包名查uid
const findUid string = `
info=$(grep -oE "^packageReplace ([0-9])+" /data/system/packages.list)
echo $info | cut -d " " -f2
`

// 通过uid放行
const allowUid = `
iptables -t nat -I nat_OUT -m owner --uid uidReplace -j ACCEPT
iptables -t mangle -I man_OUT -m owner --uid uidReplace -j ACCEPT
`

// 启动xray
const startXray = `
ulimit -n 65535
./busybox setuidgid 0:2222 ./xray -config ./config.json &>/dev/null &
`

// 关闭xray
const stopXray = "./busybox killall xray"

// 获取xray状态
const xrayStatus = "./busybox pidof xray"

// 开关网络
const data_control = `
wifiip=$(ip addr show wlan0 2>&- | grep 'inet')
if [ "$wifiip" = "" ]; then
svc data disable && sleep 0.3
svc data enable && sleep 0.3
fi
`

// 开始路由
// $wifiProxy   $denyIPV6   $TunDev  $fxqt
const startRules = `
#WiFi代理
wifiProxy="wifiProxyReplace"
#ipv6联网，0禁止，1联网
#IPV6Net="IPV6NetReplace"
#放行除tcp,udp外的流量(1放行)
allowOutsideTcpUdp="allowOutsideTcpUdpReplace"
#全局放行的app
allowAppsUid=(allowAppsUidReplace)
# 防止usb共享清一些规则造成不免
iptables -t mangle -P FORWARD DROP
iptables -t mangle -A FORWARD -p udp -j ACCEPT
iptables -t mangle -A FORWARD -p icmp -j ACCEPT
iptables -t mangle -A PREROUTING ! -p udp -j ACCEPT
# mangle OUTPUT
iptables -t mangle -N man_OUT
iptables -t mangle -N man_PRE
iptables -t mangle -A OUTPUT -j man_OUT
iptables -t mangle -A PREROUTING -j man_PRE
iptables -t mangle -A man_OUT -m owner --gid-owner 2222 -j ACCEPT
[ "$wifiProxy" = "1" ] || iptables -t mangle -A man_OUT -o wlan+ -j ACCEPT
iptables -t mangle -A man_OUT -o tun+ -j ACCEPT
# mangle PREROUTING
allow_ip="10/8,100/8,127/8,169.254/16,172.16/12,192/24,192.168/16,224/4,240/4"
iptables -t mangle -A man_PRE -d $allow_ip -j ACCEPT
# tun2socks/TPROXY 选择
	if grep -q TPROXY /proc/net/ip_tables_targets; then
		ip route add local default dev lo table 1234
		ip rule add fwmark 0x1234 lookup 1234
		iptables -t mangle -A man_PRE ! -i tun+ -p udp -j TPROXY --on-port 1231 --tproxy-mark 0x1234
		iptables -t mangle -A man_OUT ! -d 192.168/16 ! -o lo -p udp -j MARK --set-mark 0x1234
	else
		# 启动 tun v2tun_start
 		[ ! -e "/dev/net/tun" ] && mkdir -p /dev/net && ln -s /dev/tun /dev/net/tun && echo 1 > /proc/sys/net/ipv4/ip_forward
  		ip tuntap add mode tun TunDev >/dev/null 2>&1
  		ip addr add 10.0.0.10/24 dev TunDev >/dev/null 2>&1
  		ip link set TunDev up >/dev/null 2>&1
  		nohup tun2socks --tundev TunDev --netif-ipaddr 10.0.0.9 --netif-netmask 255.255.255.0 --socks-server-addr 127.0.0.1:1231 --enable-udprelay --loglevel 1 >/dev/null 2>&1 &
		ip route add default dev TunDev table 1234
		ip rule add fwmark 0x1234 lookup 1234
		iptables -t mangle -A man_PRE ! -i tun+ -p udp -j MARK --set-mark 0x1234
		iptables -t mangle -A man_OUT ! -d 192.168/16 ! -o lo -p udp -j MARK --set-mark 0x1234
		iptables -I FORWARD -i TunDev -j ACCEPT
		iptables -I FORWARD -o TunDev -j ACCEPT
	fi
# nat OUTPUT
iptables -t nat -N nat_OUT
iptables -t nat -N nat_PRE
iptables -t nat -A OUTPUT -j nat_OUT
iptables -t nat -A PREROUTING -j nat_PRE
iptables -t nat -A nat_OUT -m owner --gid-owner 2222 -j ACCEPT
[ "$wifiProxy" = "1" ] || iptables -t nat -A nat_OUT -o wlan+ -j ACCEPT
iptables -t nat -A nat_OUT -o tun+ -j ACCEPT
iptables -t nat -A nat_OUT -o lo -j ACCEPT
# 防止WiFi共享获取不到ip
iptables -t nat -A nat_OUT -d 192.168/16 -j ACCEPT
iptables -t nat -A nat_OUT -p tcp -j REDIRECT --to-ports 1230
#放行除tcp,udp外的流量(1放行)
[ "$allowOutsideTcpUdp" != '1' ] && iptables -t nat -A nat_OUT ! -p udp -j REDIRECT --to-ports 1250
# nat PREROUTING
iptables -t nat -A nat_PRE -s 192.168/16 ! -d 192.168/16 -p tcp -j REDIRECT --to-ports 1230
[ "$allowOutsideTcpUdp" != '1' ] && iptables -t nat -A nat_PRE ! -p udp -j REDIRECT --to-ports 1250

# 放行airgo核心
iptables -t nat -I nat_OUT -m owner --gid-owner 3333 -j ACCEPT
iptables -t mangle -I man_OUT -m owner --gid-owner 3333 -j ACCEPT
# 本地全局放行 nat OUTPUT
  for uid in $allowAppsUid; do
    iptables -t nat -I nat_OUT -m owner --uid $uid -j ACCEPT
  done
# 本地全局放行 mangle OUTPUT
  for uid in $allowAppsUid; do
    iptables -t mangle -I man_OUT -m owner --uid $uid -j ACCEPT
  done

# IPV6联网处理
#ip6tables -t mangle -P OUTPUT DROP
#ip6tables -t mangle -A OUTPUT -m owner --gid-owner 3333 -j ACCEPT
#ip6tables -t mangle -A OUTPUT -m owner --gid-owner 2222 -j ACCEPT
#ip6tables -t mangle -A OUTPUT -j REJECT --reject-with tcp-reset

ip6tables -t mangle -A OUTPUT -p icmpv6 -m owner --uid 0 -j ACCEPT
ip6tables -t mangle -A OUTPUT -j MARK --set-mark 0x1122

ip6tables -t mangle -A OUTPUT -m owner --gid-owner 3333 -j MARK --set-mark 0x1133
ip6tables -t mangle -A OUTPUT -m owner --gid-owner 2222 -j MARK --set-mark 0x1133
ip6tables -t mangle -P OUTPUT DROP
ip -6 rule add fwmark 0x1122 unreachable

`

// 清除路由
const clearRules = `
  while iptables -t nat -D OUTPUT -j nat_OUT; do :; done
  while iptables -t nat -D PREROUTING -j nat_PRE; do :; done
  while iptables -t mangle -D OUTPUT -j man_OUT; do :; done
  while iptables -t mangle -D PREROUTING -j man_PRE; do :; done
  iptables -t nat -F nat_OUT
  iptables -t nat -X nat_OUT
  iptables -t nat -F nat_PRE
  iptables -t nat -X nat_PRE
  iptables -t mangle -F man_OUT
  iptables -t mangle -X man_OUT
  iptables -t mangle -F man_PRE
  iptables -t mangle -X man_PRE
  while iptables -D FORWARD -i TunDev -j ACCEPT; do :; done
  while iptables -D FORWARD -o TunDev -j ACCEPT; do :; done
  iptables -t mangle -P FORWARD ACCEPT
  while iptables -t mangle -D FORWARD -p udp -j ACCEPT; do :; done
  while iptables -t mangle -D FORWARD -p icmp -j ACCEPT; do :; done
  while iptables -t mangle -D PREROUTING ! -p udp -j ACCEPT; do :; done
  # 清除ip规则
  while ip route del local default dev lo table 1234; do :; done
  while ip rule del fwmark 0x1234 lookup 1234; do :; done
  while ip route del default dev TunDev table 1234; do :; done
  ip tuntap del mode tun TunDev

  #停止v6
  #ip6tables -t mangle -P OUTPUT ACCEPT
  while ip6tables -t mangle -D OUTPUT -p icmpv6 -m owner --uid 0 -j ACCEPT; do :; done
  while ip6tables -t mangle -D OUTPUT -j MARK --set-mark 0x1122; do :; done
  while ip6tables -t mangle -D OUTPUT -m owner --gid-owner 3333 -j MARK --set-mark 0x1133; do :; done
  while ip6tables -t mangle -D OUTPUT -m owner --gid-owner 2222 -j MARK --set-mark 0x1133; do :; done
  ip6tables -t mangle -P OUTPUT ACCEPT
  while ip -6 rule del fwmark 0x1122 unreachable; do :; done
`
const OpenFirewall = `
iptables -P INPUT ACCEPT
iptables -P FORWARD ACCEPT
iptables -P OUTPUT ACCEPT
`
