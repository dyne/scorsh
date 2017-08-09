#!/bin/sh

. ./scorsh_functions

## is the daemon running?
cd ${SCORSH_APP}
is_running=$(ps $(cat scorsh.pid) | grep -c scorshd )
cd - >/dev/null
check "[ \"${is_running}\" = \"1\" ]" $0 "is_scorshd_running"


## stop the scorsh daemon
cd ${SCORSH_APP}
kill -15 $(cat scorsh.pid)
ret=$?
rm scorsh.pid
cd - >/dev/null
check "[ $ret -eq 0 ]" $0 "" 

return_results
