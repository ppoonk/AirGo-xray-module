
${0%/*}/bin/busybox setuidgid 0:3333 ${0%/*}/bin/airgo -start >/dev/null 2>&1 &

echo -e "已启动"