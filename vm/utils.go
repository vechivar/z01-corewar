package corewar

import (
	cw "corewar"
	"fmt"
	"os"
	"strings"
)

// Vérifie que x est un indice de registre valide.
func CheckRegId(x int) bool {
	if x > 0 && x <= cw.REG_NUMBER {
		return true
	} else {
		fmt.Fprintf(os.Stderr, "Cycle %d : %d is not a valid register id\n", CurrentCycle, x)
		return false
	}
}

// Transforme un pcode en liste de type d'argument
func GetArgsTypes(opcode byte) [3]byte {
	return [3]byte{opcode >> 6, (opcode >> 4) % 4, (opcode >> 2) % 4}
}

// Vérifie que les arguments sont valides pour l'opcode donné.
func CheckExpectedArgsTypes(opcode byte, argsTypes [3]byte) bool {
	var expected []string
	switch opcode {
	case 1:
		expected = []string{"2", "0", "0"}
	case 2:
		expected = []string{"23", "1", "0"}
	case 3:
		expected = []string{"1", "13", "0"}
	case 4:
		expected = []string{"1", "1", "1"}
	case 5:
		expected = []string{"1", "1", "1"}
	case 6:
		expected = []string{"123", "123", "1"}
	case 7:
		expected = []string{"123", "123", "1"}
	case 8:
		expected = []string{"123", "123", "1"}
	case 9:
		expected = []string{"2", "0", "0"}
	case 10:
		expected = []string{"123", "12", "1"}
	case 11:
		expected = []string{"1", "123", "12"}
	case 12:
		expected = []string{"2", "0", "0"}
	case 13:
		expected = []string{"23", "1", "0"}
	case 14:
		expected = []string{"123", "12", "1"}
	case 15:
		expected = []string{"2", "0", "0"}
	case 16:
		expected = []string{"1", "0", "0"}
	}

	for i := 0; i < 3; i++ {
		if !strings.Contains(expected[i], string(argsTypes[i]+'0')) {
			cmd, _ := GetCmdDatas(opcode)
			PrintArgumentError(cmd, argsTypes[i], i+1)
			return false
		}
	}
	return true
}

// Renvoie les valeurs stockées dans la mémoire en fonction de leurs types.
func GetArgumentsValues(argsTypes [3]byte, idx bool, baseAdr int) []int {
	var res []int
	var size int
	offset := 0
	for _, arg := range argsTypes {
		switch arg {
		case 1:
			size = 1

		case 2:
			if idx {
				size = 2
			} else {
				size = 4
			}
		case 3:
			size = 2
		}
		res = append(res, BytesToInt(GetArenaValues(baseAdr+offset, size)))
		offset += size
	}

	return res
}

// nom, cycle correspondant à opcode
func GetCmdDatas(opcode byte) (string, int) {
	var cycles int
	var name string
	switch opcode {
	case 1:
		cycles = 10
		name = "live"
	case 2:
		cycles = 5
		name = "ld"
	case 3:
		cycles = 5
		name = "st"
	case 4:
		cycles = 10
		name = "add"
	case 5:
		cycles = 10
		name = "sub"
	case 6:
		cycles = 6
		name = "and"
	case 7:
		cycles = 6
		name = "or"
	case 8:
		cycles = 6
		name = "xor"
	case 9:
		cycles = 20
		name = "zjmp"
	case 10:
		cycles = 25
		name = "ldi"
	case 11:
		cycles = 25
		name = "sti"
	case 12:
		cycles = 800
		name = "fork"
	case 13:
		cycles = 10
		name = "lld"
	case 14:
		cycles = 50
		name = "lldi"
	case 15:
		cycles = 1000
		name = "lfork"
	case 16:
		cycles = 2
		name = "nop"
	default:
		cycles = 1
		name = "-"
	}

	return name, cycles
}

// Affiche une erreur d'arguments.
func PrintArgumentError(cmd string, argtype byte, argind int) {
	var x string

	switch argtype {
	case 0:
		x = "0"
	case 1:
		x = "register"
	case 2:
		x = "direct"
	case 3:
		x = "indirect"
	}

	fmt.Fprintf(os.Stderr, "Command %s doesn't accept %s as argument %d \n", cmd, x, argind)
}

// Renvoie la taille prise dans la mémoire pour déclarer les arguments des types donnés
func ArgSize(argsType [3]byte, idx bool) int {
	res := 0

	for _, x := range argsType {
		switch x {
		case 1:
			res += 1
		case 2:
			if idx {
				res += 2
			} else {
				res += 4
			}
		case 3:
			res += 2
		}
	}

	return res
}

// Convertit une valeur positive non signée en tableau de bytes de taille size
func IntToBytes(val int, size int) []byte {
	value := val
	if val < 0 {
		// Grosse magouille.
		if size == 2 {
			value += 65536
		} else {
			value += 4294967296
		}
	}
	return RecIntToBytes(value, size)
}

// Convertit une valeur positive non signée en tableau de bytes de taille size
func RecIntToBytes(val int, size int) []byte {
	if size == 1 {
		return []byte{byte(val % 256)}
	} else {
		res := RecIntToBytes(val/256, size-1)
		res = append(res, byte(val%256))
		return res
	}
}

// Renvoie une valeur en fonction du type de l'argument et de sa valeur. Renvoie true si la valeur n'est pas valide.
// 1 (registre), 12 renvoie la valeur contenue dans le registre 12
// 3 (indirect), 12 renvoie la valeur contenue dans la mémoire 12 bytes après le pc du process
func CalculateArgValue(argType byte, argValue int, proc Process, applyIdx bool) (bool, int) {
	switch argType {
	case 1:
		if !CheckRegId(argValue) {
			return true, 0
		}
		return false, proc.Registers[argValue]
	case 2:
		return false, argValue
	case 3:
		if applyIdx {
			return false, BytesToInt(GetArenaValues(proc.Pc+argValue%cw.IDX_MOD, cw.IND_SIZE))
		} else {
			return false, BytesToInt(GetArenaValues(proc.Pc+argValue, cw.IND_SIZE))
		}
	}
	return true, 0
}

// Renvoie la valeur de l'arène située à adr (prend en compte %MEM_SIZE)
func GetArenaValue(adr int) byte {
	x := adr % cw.MEM_SIZE
	if x < 0 {
		x += cw.MEM_SIZE
	}
	return Arena.Memory[x]
}

// Renvoie les size valeur de l'arène situées à partir de adr (prend en compte %MEM_SIZE)
func GetArenaValues(adr int, size int) []byte {
	var res []byte
	for i := 0; i < size; i++ {
		res = append(res, GetArenaValue(adr+i))
	}
	return res
}

// Remplace la valeur de l'arène située à adr par byte (prend en compte %MEM_SIZE)
// Met à jour le player id correspondant.
func SetArenaValue(adr int, value byte, playerId int) {
	x := adr % cw.MEM_SIZE
	if x < 0 {
		x += cw.MEM_SIZE
	}
	Arena.Memory[x] = value
	Arena.LastModified[x] = playerId
}

// Renvoie la valeur correspondant au tableau de bytes (signé avec le bit le plus à gauche)
func BytesToInt(x []byte) int {
	res := 0
	fact := 1
	for k := len(x) - 1; k > 0; k-- {
		res += int(x[k]) * fact
		fact *= 256
	}

	res += int(x[0]%128) * fact
	if x[0]/128 != 0 {
		res -= fact * 128
	}
	return res
}
