#!/bin/sh



##function
write_message(){

    orev=${3:-"a1b2c3d4e5f6"}
    nrev=${4:-"9a8b7c6d5e4f"}
    
    
    cat <<EOF
---
m_id: 123456
m_repo: $1
m_branch: $2
m_oldrev: $orev
m_newrev: $nrev
...
EOF
    
}


if [ $# -le 1 ]; then
    echo "Usage: $0 <repo> <branch> [<oldrev> [<newrev]]"
    exit 1
fi
     
write_message $@
