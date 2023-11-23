# Corewar

This folder contains my implementation of the corewar project from Zone01. All specs can be found in subject.md.

## What is the project "Corewar" ?

The subject consists of detailed rules to implement a version of Corewar, a game created in 1984 by D.G. Jones. In this game, players create an assembler program that will be executed in a shared memory against other players' programs, according to rules that reproduce a simplified version of a real computer's processor executing processes. To win the game, a program has to be the last "alive" program. A program should thus be designed to try to compromise other programs' execution by corrupting their memory.

Students have to write three different parts :
- an asm program, able to compile a player's program into a binary file ready to be utilized (implemented in the asm folder)
- a vm program, able to execute a game between given players (implemented in the vm folder)
- an assembler program, able to beat a very basic opponent (implemented in the vechivar.s file)

## How to use ?

Execute compile.sh to compile both programs.