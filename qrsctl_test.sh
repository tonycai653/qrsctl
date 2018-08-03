#!/bin/bash

COMMAND="./qrsctl"
TESTBUCKET=bucketttmmpp30347
TESTKEY=buckettestfile
TESTFILE=qrsctl_test.sh
NORMAILFILETYPE=0
LINEFILETYPE=1
RULENAME=qrsctlruletest
PREFIX=qrsctlprefix
DELETEAFTERDAYS=0
TOLINEAFTERDAYS=1
UPDATETOLINEAFTERDAYS=2

Exec() {
	echo -e "==> test:" "$1"
	if ${COMMAND} "$@"; then
		echo " PASS"
	else
		echo " FAIL";
	fi
}

Usage() {
    echo -e "Usage:\n ./qrsctl_test.sh <username> <password>"
    echo -e "or:\n ./qrsctl_test.sh <access_key> <secret_key>"
}

if [[ $# -eq 2 ]]; then
    Exec login "$1" "$2"
    Exec mkbucket2 $TESTBUCKET 1
    Exec put $TESTBUCKET $TESTKEY $TESTFILE
    Exec get $TESTBUCKET $TESTKEY /tmp/$TESTFILE
    Exec stat $TESTBUCKET $TESTKEY
    Exec chtype $TESTBUCKET $TESTKEY $LINEFILETYPE
    Exec chtype $TESTBUCKET $TESTKEY $NORMAILFILETYPE
    Exec rule/add $TESTBUCKET $RULENAME $PREFIX $DELETEAFTERDAYS $TOLINEAFTERDAYS
    Exec rule/update $TESTBUCKET $RULENAME $PREFIX $DELETEAFTERDAYS $UPDATETOLINEAFTERDAYS
    Exec rule/get $TESTBUCKET $RULENAME
    Exec rule/del $TESTBUCKET $RULENAME
    Exec cat $TESTBUCKET $TESTKEY 1>/dev/zero
    Exec mv $TESTBUCKET:$TESTKEY $TESTBUCKET:$TESTKEY"1"
    Exec cp $TESTBUCKET:$TESTKEY"1" $TESTBUCKET:$TESTKEY
    Exec del $TESTBUCKET $TESTKEY
    Exec del $TESTBUCKET $TESTKEY"1"
    Exec buckets
    Exec bucketinfo $TESTBUCKET
    Exec protected $TESTBUCKET 0
    Exec separator $TESTBUCKET "-"
    Exec img $TESTBUCKET http://www.google.com
    Exec unimg $TESTBUCKET
    Exec drop -f $TESTBUCKET
    Exec info
    Exec appinfo
    Exec pfop $TESTBUCKET $TESTKEY avinfo
else
    Usage
    exit -1
fi
