cd ${0%/*}/bin
./busybox setuidgid 0:3333 ./airgo -start >/dev/null 2>&1 &
echo -e "已启动"