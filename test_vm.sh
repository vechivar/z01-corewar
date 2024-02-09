#!/bin/bash

# Automated tests script for vm
# All results in tests folder tests/res.txt file
#
# vm/basic-commands contains simple programs to test individual commands
# vm/basic-commands-expected contains expected results generated from reference program
# vm/players contains more complex programs to fight against each others
# vm/players-expected contains expected results, also generated from reference program 


# arguments : (filepath) (command arguments) (players count)
function test_basic_playercount() {
    path=$1
    args=$2
    players=$(($3))

    name="$(echo $path | cut -d'.' -f1 | cut -d'/' -f3)"
    cmd="./../vm_prgrm "

    echo -e "Testing ${name}_${players}" >> $TEST_RES
    for (( i=1; i<=$players; i++ )); do 
        cmd+="${path} "
    done

    cmd+="${args}"
    $cmd 2>/dev/null > "testRes/vm/${name}_${players}"
    expectedPath="vm/"
    expectedPath+=$(echo $path | cut -d'/' -f2)
    expectedPath+="-expected/${name}_${players}"

    diff=$(diff "${expectedPath}" "testRes/vm/${name}_${players}")
    if [ ! -z "$diff" ]; then
        echo -e "--> FAILED" >> $TEST_RES
        echo -e "Test failed for ${name}_${players}"
        echo "$diff" > "testRes/vm/${name}_${players}.diff"
    else
        rm testRes/vm/${name}_${players}
    fi
}

# arguments : (filepath) (command arguments)
function test_basic () {
    for (( i=2; i<=4; i++ )); do 
        test_basic_playercount "$1" "$2" "$i"
    done
}

#arguments: (player1 path) (player2 path) (command arguments) (output extention)
function test_players() {
    p1Path=$1
    p2Path=$2
    args=$3
    outputExt=$4

    name1="$(echo $p1Path | cut -d'.' -f1 | cut -d'/' -f3)"
    name2="$(echo $p2Path | cut -d'.' -f1 | cut -d'/' -f3)"

    cmd="./../vm_prgrm ${p1Path} ${p2Path} ${args}"
    expectedPath="vm/players-expected/${name1}_${name2}_${outputExt}"

    echo -e "Testing ${name1}_${name2}_${outputExt}" >> $TEST_RES

    $cmd 2>/dev/null > "testRes/vm/${name1}_${name2}_${outputExt}"

    diff=$(diff "${expectedPath}" "testRes/vm/${name1}_${name2}_${outputExt}")
    if [ ! -z "$diff" ]; then
        echo -e "--> FAILED" >> $TEST_RES
        echo -e "Test failed for ${name1}_${name2}_${outputExt}"
        echo -e "$diff" > testRes/vm/${name1}_${name2}_${outputExt}.diff
    else
        rm testRes/vm/${name1}_${name2}_${outputExt}
    fi

    #     echo -e "Testing ${name1}_${name2}" >> $TEST_RES
    # ./../vm_prgrm vm/players/${f1} vm/players/${f2} -x 2>/dev/null > testRes/vm/${name1}_${name2}_full
    # diff=$(diff "vm/players-expected/${name1}_${name2}_full" "testRes/vm/${name1}_${name2}_full")
    # if [ ! -z "$diff" ]; then
    #     echo -e "--> FAILED" >> $TEST_RES
    #     echo -e "Test failed for ${name1}_${name2}_full"
    #     echo -e "$diff" > testRes/vm/${name1}_${name2}_full.diff
    # else
    #     rm testRes/vm/${name1}_${name2}_full
    # fi


}

echo "Compiling programs"
./compile.sh

if [ -f "vm" ]; then
    echo "Couldn't find vm folder."
    exit
fi

cd tests
TEST_RES="testRes/res.txt"

if [ -d "testRes" ]; then 
    rm -r "testRes"
fi

mkdir testRes
mkdir testRes/vm

echo "Testing vm"

echo "VM TESTS" > $TEST_RES
echo "---------------------" >> $TEST_RES
echo "Testing basic programs" >> $TEST_RES
echo "---------------------" >> $TEST_RES

echo "Testing basic programs"

# produce basic programs expected files
if [ ! -d "vm/basic-commands-expected" ]; then
    mkdir vm/basic-commands-expected
    echo "Producing expected results from vm_ref for basic programs. This might take some time..."
    for f in $(ls vm/basic-commands); do
        name="$(echo $f | cut -d'.' -f1)"
        path="vm/basic-commands/$f"
        ./vm/vm_ref $path $path -d 55 -v 2>/dev/null > vm/basic-commands-expected/${name}_2
        ./vm/vm_ref $path $path $path -d 55 -v 2>/dev/null  > vm/basic-commands-expected/${name}_3
        ./vm/vm_ref $path $path $path $path -d 55 -v 2>/dev/null  > vm/basic-commands-expected/${name}_4
    done
fi

# testing basic-commands
for f in $(ls vm/basic-commands); do
    # name="$(echo $f | cut -d'.' -f1)"
    path="vm/basic-commands/$f"
    test_basic "$path" "-d 55 -v -x"
done

# testing special programs
echo "Testing special programs"
echo "---------------------" >> $TEST_RES
echo "Testing special programs" >> $TEST_RES
echo "---------------------" >> $TEST_RES


if [ ! -d "vm/special-expected" ]; then
    mkdir vm/special-expected
    echo "Producing expected results from vm_ref for special programs. This might take some time..."

    path="vm/special/pierino_fork.cor"
    resPath="vm/special-expected/pierino_fork"
    ./vm/vm_ref $path $path -v -d 2000  > "${resPath}_2"
    ./vm/vm_ref $path $path $path -v -d 2000 2>/dev/null > "${resPath}_3"
    ./vm/vm_ref $path $path $path $path -v -d 2000 2>/dev/null > "${resPath}_4"

    path="vm/special/pierino_lfork.cor"
    resPath="vm/special-expected/pierino_lfork"
    ./vm/vm_ref $path $path -v -d 10000  > "${resPath}_2"
    ./vm/vm_ref $path $path $path -v -d 10000 2>/dev/null > "${resPath}_3"
    ./vm/vm_ref $path $path $path $path -v -d 10000 2>/dev/null > "${resPath}_4"
fi

# pierino_fork
path="vm/special/pierino_fork.cor"
test_basic "$path" "-d 2000 -v -x"

#pierino_lfork
path="vm/special/pierino_lfork.cor"
test_basic "$path" "-d 10000 -v -x"

echo "---------------------" >> $TEST_RES
echo "Testing players" >> $TEST_RES
echo "---------------------" >> $TEST_RES

echo "Testing complex programs"

if [ ! -d "vm/players-expected" ]; then
    mkdir vm/players-expected
    echo "Producing expected results from vm_ref for players programs. This might take some time..."
    for f1 in $(ls vm/players); do
        for f2 in $(ls vm/players); do
            name1="$(echo $f1 | cut -d'.' -f1)"
            name2="$(echo $f2 | cut -d'.' -f1)"
            ./vm/vm_ref vm/players/${f1} vm/players/${f2} 2>/dev/null > vm/players-expected/${name1}_${name2}_full
            ./vm/vm_ref vm/players/${f1} vm/players/${f2} -v -d 2000  2>/dev/null > vm/players-expected/${name1}_${name2}_2000_v
            ./vm/vm_ref vm/players/${f1} vm/players/${f2} -d 10000 2>/dev/null > vm/players-expected/${name1}_${name2}_10000
            ./vm/vm_ref vm/players/${f1} vm/players/${f2} -d 50000 2>/dev/null > vm/players-expected/${name1}_${name2}_50000
        done
    done
fi

for f1 in $(ls vm/players); do
    for f2 in $(ls vm/players); do
        p1=vm/players/$f1
        p2=vm/players/$f2

        test_players $p1 $p2 "-x" "full"
        test_players $p1 $p2 "-v -d 2000 -x" "2000_v"
        test_players $p1 $p2 "-d 10000 -x" "10000"
        test_players $p1 $p2 "-d 50000 -x" "50000"
    done
done

echo "Test ended"