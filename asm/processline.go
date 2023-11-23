package corewar

import (
	cw "corewar"
	"strconv"
	"strings"
)

// Traitement d'un ligne de code du fichier .s
// La commande correspondante est ajoutée à la liste des commandes si elle est valide.
func ProcessLine(line string) {
	var splitLine []string
	for _, x := range strings.Split(line, " ") {
		if x != "" {
			splitLine = append(splitLine, x)
		}
	}
	cmd := splitLine[0]
	if rune(cmd[len(cmd)-1]) == cw.LABEL_CHAR {
		// Déclaration d'un label
		lab := cmd[:len(cmd)-1]
		CheckLabel(lab)
		labels = append(labels, label{byteCount, lab})
		if len(splitLine) == 1 {
			// Label seul sur sa ligne
			return
		}
		if rune(splitLine[1][len(splitLine[1])-1]) == cw.LABEL_CHAR {
			ExitAtLine("two consecutive labels", lineCount)
		}
		// On traite le reste de la ligne
		ProcessLine(strings.Join(splitLine[1:], " "))
		return
	}

	args := ProcessArgs(strings.Join(splitLine[1:], " "))
	commands = append(commands, CheckArgsType(cmd, args))
}

// Traitement des arguments
func ProcessArgs(x string) []asmArg {
	var res []asmArg

	// On enlève le commentaire éventuel, on enlève les espaces et on sépare selon le caractère de séparation
	// Rq : donne une certaine permissivité sur l'écriture des commandes (exemple : 'r 1' équivalent à 'r1')
	args := strings.Split(strings.ReplaceAll(strings.Split(x, "#")[0], " ", ""), string(cw.SEPARATOR_CHAR))

	if len(args) > 3 {
		ExitAtLine("wrong number of arguments", lineCount)
	}

	for _, arg := range args {
		res = append(res, BuildAsmArg(arg))
	}

	return res
}

// Récupération du type et de la valeur d'un argument.
func BuildAsmArg(arg string) asmArg {
	var res asmArg
	res.value = 0
	res.label = ""

	if rune(arg[0]) == cw.DIRECT_CHAR {
		// Argument direct
		res.argType = 'd'
		if rune(arg[1]) == cw.LABEL_CHAR {
			res.label = arg[2:]
			usedLabels = append(usedLabels, label{lineCount, res.label})
		} else {
			if rune(arg[1]) == cw.LABEL_CHAR {
				res.label = arg[2:]
			} else {
				n, err := strconv.Atoi(arg[1:])
				if err != nil {
					ExitAtLine("invalid argument", lineCount)
				}
				if !CheckIntRange(n) {
					ExitAtLine("int out of range", lineCount)
				}
				res.value = n
			}
		}
	} else if rune(arg[0]) == 'r' {
		// Argument registre
		res.argType = 'r'
		n, err := strconv.Atoi(arg[1:])
		if err != nil {
			ExitAtLine("invalid argument", lineCount)
		}
		if n > cw.REG_NUMBER || n < 1 {
			ExitAtLine("invalid register id", lineCount)
		}
		res.value = n
	} else {
		// Argument indirect
		res.argType = 'i'
		if rune(arg[0]) == cw.LABEL_CHAR {
			res.label = arg[1:]
			usedLabels = append(usedLabels, label{lineCount, res.label})
		} else {
			n, err := strconv.Atoi(arg)
			if err != nil {
				ExitAtLine("invalid argument", lineCount)
			}
			if !CheckIntRange(n) {
				ExitAtLine("int out of range", lineCount)
			}
			res.value = n
		}
	}

	return res
}

// Vérifie que le label déclaré utilise les caractères autorisés et n'est pas déjà utilisé
func CheckLabel(x string) {
	for _, c := range x {
		if !strings.Contains(cw.LABEL_CHARS, string(c)) {
			ExitAtLine("invalid character '"+string(c)+"' in label", lineCount)
		}
	}

	for _, l := range labels {
		if l.label == x {
			ExitAtLine("label is already used", lineCount)
		}
	}
}

// Vérifie que les arguments sont conformes à la commande utilisée.
// Calcule le nombre de bytes utilisés par la commande (indispensable pour la gestion des labels)
func CheckArgsType(cmd string, args []asmArg) asmCmd {
	var argsType []string // type des arguments attendus en fonction de la commande
	var pcode, idx bool
	var res asmCmd

	res.byteCount = byteCount

	switch cmd {
	case "live":
		pcode = false
		idx = false
		argsType = []string{"d"}
	case "ld":
		pcode = true
		idx = false
		argsType = []string{"id", "r"}
	case "st":
		pcode = true
		idx = false
		argsType = []string{"r", "ri"}
	case "add":
		pcode = true
		idx = false
		argsType = []string{"r", "r", "r"}
	case "sub":
		pcode = true
		idx = false
		argsType = []string{"r", "r", "r"}
	case "and":
		pcode = true
		idx = false
		argsType = []string{"rid", "rid", "r"}
	case "or":
		pcode = true
		idx = false
		argsType = []string{"rid", "rid", "r"}
	case "xor":
		pcode = true
		idx = false
		argsType = []string{"rid", "rid", "r"}
	case "zjmp":
		pcode = false
		idx = true
		argsType = []string{"d"}
	case "ldi":
		pcode = true
		idx = true
		argsType = []string{"rid", "rd", "r"}
	case "sti":
		pcode = true
		idx = true
		argsType = []string{"r", "rid", "rd"}
	case "fork":
		pcode = false
		idx = true
		argsType = []string{"d"}
	case "lld":
		pcode = true
		idx = false
		argsType = []string{"id", "r"}
	case "lldi":
		pcode = true
		idx = true
		argsType = []string{"rid", "rd", "r"}
	case "lfork":
		pcode = false
		idx = true
		argsType = []string{"d"}
	case "nop":
		pcode = true
		idx = false
		argsType = []string{"r"}
	default:
		ExitAtLine("invalid instruction", lineCount)
	}

	if len(args) != len(argsType) {
		ExitAtLine("wrong number of arguments", lineCount)
	}

	if pcode {
		byteCount = 2
	} else {
		byteCount = 1
	}

	for i, arg := range args {
		var byteSize int = 0
		if !strings.Contains(argsType[i], string(arg.argType)) {
			ExitAtLine("wrong type of argument", lineCount)
		}
		if arg.argType == 'r' {
			byteSize = 1
		} else {
			if idx || arg.argType == 'i' {
				byteSize = 2
			} else {
				byteSize = 4
			}
		}
		args[i].byteSize = byteSize
		byteCount += byteSize
	}

	res.line = lineCount
	res.cmd = cmd
	res.args = args
	byteCount += res.byteCount

	return res
}

// Vérifie que la valeur ne dépasse pas les valeurs possibles d'un int 32 bits
func CheckIntRange(x int) bool {
	return x >= -2147483648 && x <= 2147483647
}

// Vérifie que tous les labels utilisés correspondent à un label déclaré.
func CheckUsedLabels() {
	for _, ul := range usedLabels {
		flag := false
		for _, l := range labels {
			if ul.label == l.label {
				flag = true
				break
			}
		}
		if !flag {
			ExitAtLine("label does not exist", ul.value)
		}
	}
}
