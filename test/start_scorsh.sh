#!/bin/sh

. ./scorsh_functions

## start the scorsh daemon
cd ${SCORSH_APP}
./scorshd -c scorsh_example.cfg &
echo "$!" > scorsh.pid
is_running=$(ps $(cat scorsh.pid) | grep -c scorshd )
cd - >/dev/null
check_fatal "[ \"${is_running}\" = \"1\" ]" $0 "start_scorshd"

##check_fatal "[ 1 -eq 0 ]" $0 "exit_on_purpose"

## check if workers were started
cd ${SCORSH_APP}
ret=$(grep -c "Workers started correctly" scorsh.log)
cd - > /dev/null
check_fatal "[ \"$ret\" = \"1\" ]" $0 "workers_started"

## check if spooler was started
cd ${SCORSH_APP}
ret=$(grep -c "Spooler started correctly" scorsh.log)
cd - > /dev/null
check_fatal "[ \"$ret\" = \"1\" ]" $0 "spooler_started"

return_results
