#!/bin/sh

. ./scorsh_functions


### create spool directory
mkdir "${SCORSH_APP}/spool"
check_fatal "[ -$? -eq 0 ]" $0 "create_spool_folder"


### configure the remote to be used with scorsh
cd ${REMOTE_REPO}
SPOOL_DIR="${SCORSH_APP}/spool"
git config -f scorsh scorsh.spooldir $(realpath "${SPOOL_DIR}")
ret=$(git config -f scorsh scorsh.spooldir)
check_fatal "[ \"${ret}\" = \"${SPOOL_DIR}\" ]" $0 "config_remote_repo"
cd - > /dev/null
###


### copy the post-receive hook in REMOTE_REPO/hooks
cp ${SCORSH_REPO}/hooks/post-receive ${REMOTE_REPO}/hooks/
check_fatal "[ $? -eq 0 ]" $0 "copy_post-receive_hook"


### copy the scorshd program under SCORSH_APP
cp ${SCORSH_REPO}/scorshd ${SCORSH_APP}
check_fatal "[ $? -eq 0 ]" $0 "copy_scorshd"

### copy the files under "examples" into SCORSH_APP
cp -a ${SCORSH_REPO}/examples/* ${SCORSH_APP}
check_fatal "[ $? -eq 0 ]" $0 "copy_scorsh_config"




##check_fatal "[ 1 -eq 0 ]" $0 "aborting_on_purpose"

return_results
