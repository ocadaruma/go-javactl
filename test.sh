#!/bin/bash

retval=0

for testsuite in "." logger setting util
do
    pushd ${testsuite} > /dev/null
    go test $@
    [ $? -ne 0 ] && retval=1
    popd > /dev/null
done

exit ${retval}
