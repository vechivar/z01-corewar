package corewar

import (
	cw "corewar"
	"fmt"
)

// Boucle sur les commandes
func ProcessCommands() {
	for _, x := range commands {
		//PrintCommand(x)
		ProcessCommand(x)
	}
}

// Ecrit le code binaire de la commande dans le fichier de sortie
func ProcessCommand(cmd asmCmd) {
	res, idx, pcode := GetCmdValues(cmd.cmd)

	// Gestion du pcode
	if pcode {
		var pcodeVal byte = 0
		var index byte = 64
		for _, arg := range cmd.args {
			switch arg.argType {
			case 'r':
				pcodeVal += index * byte(cw.REG_CODE)
			case 'd':
				pcodeVal += index * byte(cw.DIR_CODE)
			case 'i':
				pcodeVal += index * byte(cw.IND_CODE)
			}
			index = index / 4
		}
		res = append(res, pcodeVal)
	}

	// Gestion des arguments
	for _, arg := range cmd.args {
		if arg.argType == 'r' {
			res = append(res, byte(arg.value))
		} else {
			if arg.label != "" {
				// On récupère la valeur correspondant au décalage du label
				arg.value = GetLabelValue(arg.label) - (cmd.byteCount)
			}
			if idx || arg.argType == 'i' {
				// Argument écrit sur deux bytes
				res = append(res, IntToBytes(arg.value, 2)...)
			} else {
				// Argument écrit sur quatre bytes
				res = append(res, IntToBytes(arg.value, 4)...)
			}
		}
	}
	outputFile.Write(res)
}

// Renvoie les informations (opcode, idx, pcode) liées à une cmd
func GetCmdValues(cmd string) ([]byte, bool, bool) {
	var val byte
	idx := false
	pcode := true

	switch cmd {
	case "live":
		val = 1
		pcode = false
	case "ld":
		val = 2
	case "st":
		val = 3
	case "add":
		val = 4
	case "sub":
		val = 5
	case "and":
		val = 6
	case "or":
		val = 7
	case "xor":
		val = 8
	case "zjmp":
		val = 9
		idx = true
		pcode = false
	case "ldi":
		val = 10
		idx = true
	case "sti":
		val = 11
		idx = true
	case "fork":
		val = 12
		idx = true
		pcode = false
	case "lld":
		val = 13
	case "lldi":
		val = 14
		idx = true
	case "lfork":
		val = 15
		idx = true
		pcode = false
	case "nop":
		val = 16
	}

	return []byte{val}, idx, pcode
}

// Convertit une valeur positive non signée en tableau de bytes de taille size
func IntToBytes(val int, size int) []byte {
	value := val
	if val < 0 {
		// Grosse magouille.
		// A changer.
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

// Renvoie le compteur de bytes associé à un label
func GetLabelValue(label string) int {
	for _, x := range labels {
		if x.label == label {
			return x.value
		}
	}
	Exit("can't find label value. should never happen")
	return 0
}

// Pour débuguer
func PrintCommand(cmd asmCmd) {
	fmt.Print(cmd.cmd + " : ")
	for _, x := range cmd.args {
		PrintArg(x)
	}
	fmt.Print("\n")
}

// Pour débuguer
func PrintArg(arg asmArg) {
	fmt.Printf(" || type  %c ; %d %s (size : %d)", arg.argType, arg.value, arg.label, arg.byteSize)
}
