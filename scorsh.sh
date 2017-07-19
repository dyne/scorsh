#!/bin/sh

##
## Rudimentary implementation of a scorsh client, in POSIX sh
##

## func
build_command(){

    TAG=$1
    shift
    ARGS=$@

    ARGLIST=""
    for a in ${ARGS}; do
        ARGLIST="${ARGLIST}\"$a\","
    done
    ARGLIST=$(echo ${ARGLIST}| sed -r -e 's/,$//g')
    
    cmd_str=$(cat <<EOF
---
scorsh:
  [ 
    {
     s_tag: "$TAG",
     s_args: [${ARGLIST}]
    }
  ]
...
EOF
           )
    
}


if [ $# -le 0 ]; then
    echo "Usage: $0 <tag> [<arg>...]" 
    exit 1
fi

echo $@

build_command $@

echo $@

echo "${cmd_str}"

echo "$0"


## Check if we have to create the commit
script_name=$(basename $0)
echo "script_name: ${script_name}"
if [ "${script_name}" = "scorsh-commit" ]; then
    
    git commit --allow-empty -Sscorsh -m "${cmd_str}"
    
fi
