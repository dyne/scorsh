#!/bin/sh

. ./scorsh_functions

LINE_FILE=.line_file
cd ${SCORSH_APP}
LAST_LINE=$(wc -l scorsh.log | cut -d " " -f 1)
echo ${LAST_LINE} > ${LINE_FILE}
cd - > /dev/null

### create an empty git commit without a scorsh-command the first
### commit will be ignored by scorsh for the moment, due to an error
### on the 0000...0000 oid
cd ${LOCAL_REPO}
git commit --allow-empty -m "this is an empty commit"
check "[ $? -eq 0 ]" $0 "create_first_commit"

LAST_LINE=$(cat ${SCORSH_APP}/${LINE_FILE})
git push

cd - > /dev/null

cd ${SCORSH_APP}
ret=$(tail -n +${LAST_LINE} scorsh.log | grep -c "Invalid commit ID")
sleep 1
LAST_LINE=$(wc -l scorsh.log | cut -d " " -f 1)
echo ${LAST_LINE} > ${LINE_FILE}
cd - > /dev/null

check "[ \"$ret\" = \"2\" ] " $0 "check_first_commit"

### create two more commits without scorsh-commands
cd ${LOCAL_REPO}
git commit --allow-empty -m "second commit"
check "[ $? -eq 0 ]" $0 "create_second_commit"
git commit --allow-empty -m "third commit"
check "[ $? -eq 0 ]" $0 "create_third_commit"

LAST_LINE=$(cat ${SCORSH_APP}/${LINE_FILE})
commits=$(git log | grep "^commit " | cut -d " " -f 2 | head -2)
message_id=$(git push | grep "remote: id:" | cut -d " " -f 3)
cd - > /dev/null

cd ${SCORSH_APP}
sleep 1
for c in ${commits}; do
    ret=$(tail -n +${LAST_LINE} scorsh.log | grep -c "error parsing commit ${c}: no SCORSH message found")
    check "[ \"$ret\" = \"2\" ]" $0 "process_empty_commit"
done

LAST_LINE=$(wc -l scorsh.log | cut -d " " -f 1)
echo ${LAST_LINE} > ${LINE_FILE}
cd - > /dev/null


##check_fatal "[ 1 -eq 0 ]" $0 "abort_on_purpose"

rm ${LINE_FILE}

return_results
