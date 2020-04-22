#! /bin/bash

set -e

xmlDirName="./xml/testcases"
resultDirName="./results/qfunittests"

testsi1folders=("si-2-1" "si-3-1" "si-5-1")
testsi2folders=("si-2-2" "si-3-2" "si-5-2")
testsi5folders=("si-2-5" "si-3-5" "si-5-5")
testsi10folders=("si-2-10" "si-3-10" "si-5-10")

testssfolders=("ss-2")


if [ -d $resultDirName ]; then
    pwd
    rm -rf $resultDirName
fi

# funcCreateDir is to create folders for results
funcCreateDir(){
 dir=$1
 mkdir -p $resultDirName/$dir
}

# funcRunTests is to run tests
funcRunTests(){
    name=("${!1}")
    fileNamePrepare=$2
    fileNameAccept=$3
    for n in "${name[@]}"
    do
        echo $n
        funcCreateDir $n

        for filename in $xmlDirName/$n/$fileNamePrepare
        do
            cns=$(basename ${filename%.*})
            echo $n/$cns
            go test -v -run TestPrepareQFUnitTest main_test.go qspecs_unit_test.go qspecs.go singlepaxos.pb.go -prepareQFTCsDir="$filename" > $resultDirName/$n/$cns.out
        done

        for filename in $xmlDirName/$n/$fileNameAccept
        do
            cns=$(basename ${filename%.*})
            echo $n/$cns
            go test -v -run TestAcceptQFUnitTest main_test.go qspecs_unit_test.go qspecs.go singlepaxos.pb.go -acceptQFTCsDir="$filename" > $resultDirName/$n/$cns.out
        done
    done
}

# test case execution
if [ ! -d $resultDirName ]; then
    funcRunTests testsi1folders[@] tcs-paxos-si-1-qfprepare* tcs-paxos-si-1-qfaccept*
    funcRunTests testsi2folders[@] tcs-paxos-si-2-qfprepare* tcs-paxos-si-2-qfaccept*
    funcRunTests testsi5folders[@] tcs-paxos-si-5-qfprepare* tcs-paxos-si-5-qfaccept*
    funcRunTests testsi10folders[@] tcs-paxos-si-10-qfprepare* tcs-paxos-si-10-qfaccept*
    funcRunTests testssfolders[@] tcs-paxos-ss-qfprepare* tcs-paxos-ss-qfaccept*
else
    echo "the directories for results exist"
fi


