#!/bin/bash

# Resolve absolute path of the metricsd root directory
ABSPATH="$(cd "${0%/*}" 2>/dev/null; echo "$PWD"/"${0##*/}")"
ROOT=$(dirname "$ABSPATH")
# Some useful variables
SERVER="${ROOT}/metricsd"
SERVER_ARGS=""
SERVER_PID="${ROOT}/log/metricsd.pid"
SERVER_LOG="${ROOT}/log/metricsd.log"

mkdir -p "${ROOT}/log"
# Make sure $ROOT is a current working directory
cd $ROOT

case "$1" in
    run)
        exec ${SERVER} ${SERVER_ARGS} >> ${SERVER_LOG} 2>&1
    ;;

    start)
        if [ -f "${SERVER_PID}" ]; then
          if ps -o pid= -p $(cat ${SERVER_PID}) 2&>1 > /dev/null; then
            echo "ERROR: server is already started"
            exit 1
          else
            echo "WARNING: stale PID file found, removing"
          fi
        fi
        echo -n "Starting server: "
        exec nohup ${SERVER} ${SERVER_ARGS} >> ${SERVER_LOG} 2>&1 &
        echo $! > ${SERVER_PID}
        echo "Ok"
    ;;

    stop)
        echo -n "Stopping server: "
        kill `cat ${SERVER_PID}` &> /dev/null
        rm -f ${SERVER_PID}
        echo "Ok"
    ;;

    stopkill)
        echo -n "Killing server: "
        kill $(cat ${SERVER_PID}) &> /dev/null
        for i in {1..7};
        do
          if ps -o pid= -p $(cat ${SERVER_PID}) 2&>1 > /dev/null; then
            sleep 10
          else
            killed=1
            break
          fi
        done
        if [ "$killed" -eq 0 ]; then
          kill -9 $(cat ${SERVER_PID}) &> /dev/null
        fi
        rm -f ${SERVER_PID}
        echo "Ok"
    ;;

    restart)
        $0 stop
        sleep 1
        $0 start
        ;;

    chkconfig)
        exec ${SERVER} ${SERVER_ARGS} -test
        ;;

    *)
        echo "Usage: $0 {start|stop|restart|stopkill|chkconfig}"
        exit 1
    ;;
esac

exit 0
