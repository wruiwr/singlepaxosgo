#! /bin/bash

set -e

xmlDirName="./xml/testcases"
resultDirName="./results/qfunittests"
coverageDirName="./results/qfcoverage"

testsi1folders=("si-2-1" "si-3-1" "si-5-1")
testsi2folders=("si-2-2" "si-3-2" "si-5-2")
testsi5folders=("si-2-5" "si-3-5" "si-5-5")
testsi10folders=("si-2-10" "si-3-10" "si-5-10")

testssfolders=("ss-2")

cov="cover"
resCov="totalcover"

if [ -d $coverageDirName ]; then
    pwd
    rm -rf $coverageDirName
fi


# funcCreateOneDir is to create a directory for the results
funcCreateOneDir(){
    dir=$1
    mkdir -p $coverageDirName/$dir
}

# funcCreateTwoDirs is to create two directories for the results
funcCreateTwoDirs(){
    dir=$1
    mkdir -p $resultDirName/$dir
    mkdir -p $coverageDirName/$dir
}

#funcCreateDirs is to create all directories for the results
funcCreateDirs(){
    name=("${!1}")

    if [ ! -d $resultDirName ]; then
        for n in "${name[@]}"
        do
             funcCreateTwoDirs $n
        done
    elif [ ! -d $coverageDirName ]; then
        for n in "${name[@]}"
        do
             funcCreateOneDir $n
        done
    else
        for n in "${name[@]}"
        do
             funcCreateOneDir $n
        done
    fi
}

funcCreateDirs testsi1folders[@]
funcCreateDirs testsi2folders[@]
funcCreateDirs testsi5folders[@]
funcCreateDirs testsi10folders[@]
funcCreateDirs testssfolders[@]


# funcRunTests run all the coverages for quorum functions
funcRunTests(){
    name=("${!1}")
    fileNamePrepare=$2
    fileNameAccept=$3
    cover=$4
    for n in "${name[@]}"
    do
        # prepareQF
        for filename in $xmlDirName/$n/$fileNamePrepare
        do
            cns=$(basename ${filename%.*})
            echo $n/$cns
            go test -run TestPrepareQFUnitTest -prepareQFTCsDir="$filename" -covermode=count -coverprofile=$coverageDirName/$n/$cov$cns.out > $coverageDirName/$n/$resCov$cns.out
        done

        # test acceptQF
        for filename in $xmlDirName/$n/$fileNameAccept
        do
            cns=$(basename ${filename%.*})
            echo $n/$cns
            go test -run TestAcceptQFUnitTest -acceptQFTCsDir="$filename" -covermode=count -coverprofile=$coverageDirName/$n/$cov$cns.out > $coverageDirName/$n/$resCov$cns.out
        done

        # obtain coverage files
        cd $coverageDirName/$n
        for filename in $cover
        do
            go tool cover -func=$filename > func$filename
        done
        cd ../../..
        cd $coverageDirName/$n
        for filename in $cover
        do
            go tool cover -html=$filename -o $filename.html
        done
        cd ../../..
    done
}

funcRunTests testsi1folders[@] tcs-paxos-si-1-qfprepare* tcs-paxos-si-1-qfaccept* cov*
funcRunTests testsi2folders[@] tcs-paxos-si-2-qfprepare* tcs-paxos-si-2-qfaccept* cov*
funcRunTests testsi5folders[@] tcs-paxos-si-5-qfprepare* tcs-paxos-si-5-qfaccept* cov*
funcRunTests testsi10folders[@] tcs-paxos-si-10-qfprepare* tcs-paxos-si-10-qfaccept* cov*
funcRunTests testssfolders[@] tcs-paxos-ss-qfprepare* tcs-paxos-ss-qfaccept* cov*

