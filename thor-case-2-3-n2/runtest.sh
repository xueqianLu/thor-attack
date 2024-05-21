#!/bin/bash
testcase=("base" "case1" "case2" "case3" "case4" "case5" "case6" "case7")

#for var in ${testcase[@]};
#do
#	echo "build image for test-$var"
#
#	git checkout test-$var && TAG="$var" make docker
#done

git checkout test-base

for var in ${testcase[@]};
do
	echo "test $var start"
	cd deploy && ./test.sh $var && cd -
	echo "test $var end"
done
