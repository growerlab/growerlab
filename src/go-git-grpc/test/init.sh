#!/bin/bash

baseDir=$GO_GIT_GRPC_TEST_DIR
repoDir=$baseDir/testrepo_bare
tmpRepoDir=/tmp/testrepo_bare

if [ ! -d $repoDir ]; then
    mkdir -p $repoDir
    cd $repoDir
    git init --bare
fi


if [ ! -d $tmpRepoDir ]; then
    mkdir -p $tmpRepoDir
    cd $tmpRepoDir || true

    git init
    git remote add origin $repoDir
    for i in {1..10} ; do
        touch ${i}.txt
        echo ${i} >> ${i}.txt
        git add .
        git commit -m "commit ${i}"
    done

    for i in {1..5} ; do
        git checkout -b branch_${i}
        touch ${i}.log
        echo ${i} >> ${i}.log
        git add .
        git commit -m "commit ${i} in 'branch_${i}'"
        git tag -a v${i}.0 -m "v${i}.0"
    done

    git push --all
    git push --tags
fi


