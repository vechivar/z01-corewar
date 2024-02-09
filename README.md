# Corewar

This folder contains my implementation of the corewar project from Zone01. All specs can be found in subject.md.

## What is the project "Corewar" ?

The subject consists of detailed rules to implement a version of Corewar, a game created in 1984 by D.G. Jones. In this game, players create an assembler program that will be executed in a shared memory against other players' programs, according to rules that reproduce a simplified version of a real computer's processor executing processes. To win the game, a program has to be the last "alive" program. A program should thus be designed to try to compromise other programs' execution by corrupting their memory.

Students have to write three different parts :
- an asm program, able to compile a player's program into a binary file ready to be utilized (implemented in the asm folder)
- a vm program, able to execute a game between given players (implemented in the vm folder)
- an assembler program, able to beat a very basic opponent (implemented in the vechivar.s file)

## How to use ?

`./compile.sh` to compile both programs.
`./asm_prgrm "filename.s"` to compile binary from a player's program
`./vm_prgrm "player1.cor" "player2.cor" ["player3.cor"] ["player4.cor"] [-d N -v -x]`
    -d N option to stop the game after N cycles
    -v option to visualize game state every cycle
    -x option for a display closer to reference vm output.
    This option should also be used when outputing in a file instead of terminal.

## Tests

./test_asm.sh and ./test_vm.sh contain basic scripts for automated testing of asm and vm.
tests folder contains all files needed for tests, including asm_ref and vm_ref programs provided in the subject. These programs can be downloaded from https://zone01normandie.org/git/root/public/src/branch/master/subjects/corewar/data > playground.zip rather than being trusted and executed from this repository. Total size for exepected files is around 600MB.

### Asm tests

tests/asm contains a valid and invalid folder. Invalid files are expected to not output anything, whereas valid files are expected to produce an output strictly equal to the asm_ref output.

### Vm tests

All vm tests are expected to produce an output (almost...) identical to the vm_ref program. -x option is mandatory.

tests/vm/basic-commands contains very simple programs to test individual commands. All those programs are tested in a 2, 3 and 4 players context for 55 cycles. -v option is used to strictly control memory behaviour.

tests/vm/special contains files that need specific command lines, for example tests for fork and lfork commands that need more than 55 cycles to execute.

tests/vm/players contains players to run against each others. Each match is tested with :
    - no parameters (normal game)
    - visualize with 2000 cycles deadline
    - 50000 cycles deadline
    - 10000 cycles deadline

## Project status

This corewar isn't validated from the zone01 Rouen community yet.
Vm requires additional testing, especially with stronger, more complex players to face each others that should be provided from zone01-edu teams.
So far, this version is able to pass rigourously every test, except for a single empty line added before winning message when the last memory printed before winning message includes the last memory line adress. This can be seen in ameba_vechivar_50000, crab_vechivar_50000 and vechivar_vechivar_50000 tests.
