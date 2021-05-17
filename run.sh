#!/bin/sh
# Start or stop receve
# chkconfig:345 89 17
# description: "receve.php"
# Source function library
if [ -f /etc/rc.d/init.d/functions ]; then
. /etc/rc.d/init.d/functions
elif [ -f /etc/init.d/functions ]; then
. /etc/init.d/functions
elif [ -f /etc/rc.d/functions ]; then
. /etc/rc.d/functions
fi
#定义启动关闭时,会显示是否Ok,并有颜色提示.
BOOTUP=color
MOVE_TO_COL="echo -en \\033[60G"
SETCOLOR_SUCCESS="echo -en \\033[1;32m"
SETCOLOR_FAILURE="echo -en \\033[1;31m"
SETCOLOR_WARNING="echo -en \\033[1;33m"
SETCOLOR_NORMAL="echo -en \\033[0;39m"
#
#ok 为绿色
function echo_success(){
   [ "$BOOTUP" = "color" ] && $MOVE_TO_COL
   echo -n "["
   [ "$BOOTUP" = "color" ] && $SETCOLOR_SUCCESS
   echo -n $"  OK  "
   [ "$BOOTUP" = "color" ] && $SETCOLOR_NORMAL
   echo -n "]"
   echo -ne "\r"
   echo -e "\n"
   return 0
 }
#failed为暗红色
function echo_failure() {
   [ "$BOOTUP" = "color" ] && $MOVE_TO_COL
   echo -n "["
   [ "$BOOTUP" = "color" ] && $SETCOLOR_FAILURE
   echo -n $"FAILED"
   [ "$BOOTUP" = "color" ] && $SETCOLOR_NORMAL
   echo -n "]"
   echo -ne "\r"
   echo -e "\n"
   return 1
 }
#判断进程的运行状态,并记录PID,我运行着两个进程,所以要记录多次.
function isrun
{
	   numpro=`/bin/ps -ef | grep "goblog" | grep -v 'grep'|uniq|wc  -l`
       if [ $numpro -ge 1 ];then
               export numpro
               /bin/ps -ef | grep "goblog" | grep -v 'grep'|grep -v  grep |awk '{print $2}' >./app.pid
               export run=1
       else
               export run=0
       fi
}

#需要运行的脚本
RETVAL=0
ShellBin="./goblog"
LockFile=./lock
# See how we were called.
#启动
start() {
    isrun
    if [ $run -eq 0 ];then
       # Start
       if [ ! -f $ShellBin ];then
          echo "FATAL: No such programme";exit 4;
       fi
       echo -n "Starting: "
       $ShellBin &
       RETVAL=$?
       if [ $RETVAL -eq 0 ] ;then
               touch $LockFile
               echo_success
       else
               echo_failure
       fi
       return $RETVAL
   else
       echo $ShellBin " always on!"
   fi
}
#关闭
stop() {
   isrun
   if [ $run -eq 1 ];then
       # Stop
       echo -n $"Shutting down : "
       for i in `cat ./app.pid`
       do
            kill  -9  $i
       done
       RETVAL=$?
       if [ $RETVAL -eq 0 ] ;then
			rm -f $LockFile
			echo_success
       else
			echo_failure
       fi
       return $RETVAL
   else
       echo $ShellBin " always off!"
   fi
}
# call the function we defined
case "$1" in
 start)
       start
       ;;
 stop)
       stop
       ;;
 restart|reload)
       stop
       start
       RETVAL=$?
       ;;
 status)
       status  receve
       RETVAL=$?
       ;;
 *)
       echo $"Usage: $0 {start|stop|restart|reload|status}"
       exit 2
esac
exit $RETVAL