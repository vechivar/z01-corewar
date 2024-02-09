.name "vechivar"
.description "bulldozer"

# The program binary code is stored in the registers from 10 to 3. Other registers contain values used by the program.
# The program calls live, then duplicates itself a little further in the memory, from end to beginning.
# It's a simple loop that changes the argument of the command copying the content of registers.
# On the last loop, the last zjmp command is overwritten, allowing the program to keep progressing in memory.

ld %1879769140, r13
ld %65540, r14

st r1, r2
ld %1, r1

ld %55577871, r3
ld %89394958, r4
ld %251883535, r5
ld %394096, r6
ld %16777218, r7
ld %-1879048192, r8
ld %1051135, r9
ld %-369098752, r10

st r2, 6            # writes players id for live command

prgrm:
live %0                 # 5 bytes       live
st r13, r15             # 4 bytes       initializes r15 (arguments on line 25)           
sub r15, r14, r15       # 5 bytes       decreases the arguments : register -1, address -4
st r15, 6               # 5 bytesx      writes arguments of next line
st r1, 0                # 5 bytes       duplicates part of the program a little further in memory
ld %0, r16              # 7 bytes       sets carry to true
zjmp %-22               # 3 bytes       loops on beginning of the program. overritten at the end