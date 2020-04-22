#! /bin/bash

set -e

xmlDirName="./xml/testcases"
resultDirName="./results/systemtests"

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

#funcCreateDirs is to create all directories for the results
funcCreateDirs(){
    name=("${!1}") # get an array

    for n in "${name[@]}"
    do
        funcCreateDir $n
    done
}

funcCreateDirs testsi1folders[@]
funcCreateDirs testsi2folders[@]
funcCreateDirs testsi5folders[@]
funcCreateDirs testsi10folders[@]
funcCreateDirs testssfolders[@]


# funcRunTests run all the coverages for quorum functions
funcRunTests(){
    name=("${!1}")
    fileName=$2

    for n in "${name[@]}"
    do
        for filename in $xmlDirName/$n/$fileName
        do
            cns=$(basename ${filename%.*})
            echo $n/$cns
            go test -v system_test.go main_test.go adapterconnector.go singlepaxos.pb.go qspecs.go paxosreplica.go paxos_gorums_helper.go proposer.go acceptor.go ld.go fd.go fd_ping_crash.go -tags crash -systemTCsDir="$filename" > $resultDirName/$n/$cns.out
        done
    done
}

funcRunTests testsi1folders[@] tcs-paxos-si-1-system*
funcRunTests testsi2folders[@] tcs-paxos-si-2-system*
funcRunTests testsi5folders[@] tcs-paxos-si-5-system*
funcRunTests testsi10folders[@] tcs-paxos-si-10-system*
funcRunTests testssfolders[@] tcs-paxos-ss-system*

