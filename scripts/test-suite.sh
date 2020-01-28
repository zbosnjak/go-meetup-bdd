#!/bin/bash




# Note:
#   The $feature-$filename.json naming convention allows us to support 
#   more than one feature file in a directory.

basedir=`pwd`
specsdir=$basedir/tests/specs
resultsdir=$basedir/tests/results

mkdir -p $resultsdir

# exit on error
set -e

# generate a report for a feature
#
# ex: make tests-run REPORT=version
if [[ $REPORT ]] ; then
    cd $specsdir/$REPORT
    echo -e "Running scenarios for $REPORT..."
    for file in *.feature ; do
#        godog -f cucumber $file > $resultsdir/$REPORT-$file.json ; \
        godog -f cucumber $file > $resultsdir/$REPORT.json ; \
    done
    exit
fi

# generate a full report
#
# ex: make tests-run FULL_REPORT=1
if [[ $FULL_REPORT ]] ; then
    cd $specsdir
    for feature in * ; do
        cd $specsdir/$feature
        echo -e "Running scenarios for $feature..."
        for file in *.feature ; do
#            godog -f cucumber $file > $resultsdir/$feature-$file.json
            godog -f cucumber $file > $resultsdir/$feature.json
        done
    done
    exit
fi

# if a feature is specified, run that feature
# uses --strict to fail when there is a pending or undefined step
# in case this is used in a CI script
#
# ex: make tests-run  FEATURE=version
if [[ $FEATURE ]] ; then
    cd $specsdir/$FEATURE
    godog -t "~@manual" --strict *.feature
    exit
fi

# default: run the entire suite
#
# ex: make tests-run
sleep 5s
count=1
for dir in $specsdir/*/ ; do
    cd $dir
    echo "Running test #$count: $dir"
    godog -t "~@manual" --strict *.feature
    count=$((count+1))
done