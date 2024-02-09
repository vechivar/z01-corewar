#!/bin/bash

# Automated tests script for asm
# All results in tests folder tests/res.txt file
#
# tests/asm/invalid folder contains invalid tests
# tests/asm/valid folder for valid tests, tests/asm/expected folder for expected results
# all failed tests should be in tests/testRes/asm folder

echo "Compiling programs"
./compile.sh

if [ -f "asm" ]; then
    echo "Couldn't find asm folder."
    exit
fi

cd tests
TEST_RES="testRes/res.txt"

if [ -d "testRes" ]; then 
    rm -r "testRes"
fi

mkdir testRes
mkdir testRes/asm

echo "Testing asm"

echo "ASM TESTS" > $TEST_RES
echo "---------------------" >> $TEST_RES
echo "Testing invalid files" >> $TEST_RES
echo "---------------------" >> $TEST_RES

echo "Testing invalid files"

# asm invalid files
for f in $(ls asm/invalid); do
    name="$(echo $f | cut -d'.' -f1)"
    echo -e "Testing $name" >> $TEST_RES
    ./../asm_prgrm "asm/invalid/$f" > /dev/null
    if [ -f "$name.cor" ]; then
        echo -e "--> FAILED : file was created" >> $TEST_RES
        echo -e "Test failed for invalid : $name"
        mv "$name.cor" "testRes/asm"
    fi
done

echo "" >> $TEST_RES
echo "---------------------" >> $TEST_RES
echo "Testing valid files" >> $TEST_RES
echo "---------------------" >> $TEST_RES

echo "Testing valid files"

# produce expected results from asm_ref
if [ ! -d "asm/expected" ]; then
    mkdir asm/expected
    echo "Producing expected results from asm_ref. This might take some time..."
    for f in $(ls asm/valid); do
        name="$(echo $f | cut -d'.' -f1)"
        ./asm/asm_ref asm/valid/$f 2>/dev/null
        mv asm/valid/${name}.cor asm/expected
    done
fi

#asm valid files
for f in $(ls asm/valid); do
    name="$(echo $f | cut -d'.' -f1)"
    echo -e "Testing $name" >> $TEST_RES
    ./../asm_prgrm "asm/valid/$f" > /dev/null
    if [ ! -f "$name.cor" ]; then
    #no res file found
        echo -e "--> FAILED : file was not created" >> $TEST_RES
        echo -e "Test failed for valid : $name" 
    else
        if [ ! -f "asm/expected/$name.cor" ]; then
        #expected file not found
            echo -e "--> WARNING : expected result wasn't found, test is invalid" >> $TEST_RES
            echo -e "Test failed for valid : $name (expected file not found)"
            rm "$name.cor"
        else 
        #created and expected file are different
            diff=$(diff "asm/expected/$name.cor" "$name.cor")
            if [ ! -z "$diff" ]; then
                echo -e "--> FAILED" >> $TEST_RES
                echo -e "Test failed for valid : $name" 
                mv "$name.cor" "testRes/asm"
                echo -e "$diff" > testRes/asm/${name}.diff
            else
            #test succeed. remove result file
                rm "$name.cor"
            fi
        fi
    fi
done

echo "Test ended"
