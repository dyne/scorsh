#!/bin/sh


. ./scorsh_functions

### kill scroshd, if it's running
if [ -f "${SCORSH_APP}/scorshd.pid" ]; then
    kill -9 $(cat "${SCORSH_APP}/scorshd.pid")
fi

### remove all the folders
rm -rf ${SCORSH_REPO}
rm -rf ${SCORSH_APP}
rm -rf ${REMOTE_REPO}
rm -rf ${LOCAL_REPO}

check "[ ! -d \"${SCORSH_REPO}\" ]" $0 "remove_scorsh_repo"
check "[ ! -d \"${SCORSH_APP}\" ]" $0 "remove_scorsh_app"
check "[ ! -d \"${REMOTE_REPO}\" ]" $0 "remove_remote_repo"
check "[ ! -d \"${LOCAL_REPO}\" ]" $0 "remove_local_repo"

return_results
