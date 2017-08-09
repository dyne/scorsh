#!/bin/sh
##
## This script is part of the scorsh testbed. It must be executed by
## the scorsh_testsuite script
##

. ./scorsh_functions

### remove the directories if they exist already
[ -d ${REMOTE_REPO} ] && rm -rf ${REMOTE_REPO}
[ -d ${LOCAL_REPO} ] && rm -rf ${LOCAL_REPO}

### create the repository 
git init --bare ${REMOTE_REPO}
check "[ $? -eq 0 ]" $0 "create_remote_repo"

### clone it
git clone ${REMOTE_REPO} ${LOCAL_REPO}
check "[ $? -eq 0 ]" $0 "clone_remote_repo"





### create the directory where scorsh will be cloned
mkdir ${SCORSH_REPO}
check "[ $? -eq 0 ]" $0 "create_scorsh_repo_folder"

### clone the scorsh repo
olddir=$(pwd)
cd ${SCORSH_REPO}
git clone ${SCORSH_URL} ./
check "[ $? -eq 0 ]" $0 "clone_scorsh_repo"
cd ${olddir}

### make the scorshd executable
olddir=$(pwd)
cd ${SCORSH_REPO}
make
check "[ $? -eq 0 ]" $0 "make_scorshd"
cd ${olddir}


### create the directory where the scorsh app will run
mkdir ${SCORSH_APP}
check "[ $? -eq 0 ]" $0 "create_scorsh_app_folder"


### create spool directory
mkdir "${SCORSH_APP}/spool"
check "[ -$? -eq 0 ]" $0 "create_spool_folder"


### configure the remote to be used with scorsh
cd ${REMOTE_REPO}
git config -f scorsh scorsh.spooldir $(realpath "${SCORSH_APP}/spool")
check "[ $? -eq 0 ]" $0 "config_remote_repo"
cd -


return_results
