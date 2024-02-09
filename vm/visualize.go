package corewar

import (
	"encoding/hex"
	"fmt"
	"strings"

	cw "corewar"
)

// Affiche l'état actuel du jeu
func Visualize() {
	fmt.Println("Processes:")
	PrintProcesses()
	fmt.Println("Players:")
	PrintPlayers()
	fmt.Println("Arena:")
	PrintMemory()
	fmt.Println()
}

func PrintMemory() {
	if PerfectOutputMode {
		PrintMemoryPerfectMatch()
	} else {
		PrintMemoryCustom()
	}
}

// Affiche la mémoire
func PrintMemoryCustom() {
	var pcs []int

	// pcs des processus, à surligner
	for _, p := range Processes {
		pcs = append(pcs, p.Pc)
	}

	bytesPerLine := 32
	printingEmptyLines := true
	currentPlayer := 0
	for i := 0; i < cw.MEM_SIZE/bytesPerLine; i++ {
		// Affichage d'une ligne
		notEmpty := false
		line := ""

		// affichage en blanc de l'adresse de base de la ligne
		if currentPlayer != 0 {
			line += Colors[0]
			currentPlayer = 0
		}
		line += fmt.Sprintf("%08x  ", i*bytesPerLine)

		// récupération des valeurs à afficher
		for j := 0; j < bytesPerLine; j++ {
			ind := bytesPerLine*i + j // indice du byte à afficher
			player := Arena.LastModified[ind]
			val := hex.EncodeToString(GetArenaValues(ind, 1))

			if player != currentPlayer {
				// changement de couleur
				line += Colors[player]
				currentPlayer = player
			}

			if len(pcs) > 0 && pcs[0] == ind {
				// pc d'un processus : surlignage
				pcs = pcs[1:]
				line += "\033[47m" + val + "\033[0m" + Colors[currentPlayer] + " "
				notEmpty = true
			} else {
				line += val + " "
			}

			if GetArenaValue(ind) != 0 {
				// valeur non nulle, on affiche la ligne
				notEmpty = true
			}
		}
		if notEmpty {
			fmt.Println(line)
			printingEmptyLines = false
		} else {
			if !printingEmptyLines {
				currentPlayer = 0
				fmt.Println(Colors[0] + "...")
				printingEmptyLines = true
			}
		}
	}
	fmt.Print(Colors[0])
}

func PrintMemoryPerfectMatch() {
	lastPrintedAdr := -1
	lastBytePrinted := false
	bytesPerLine := 32
	for i := 0; i < cw.MEM_SIZE/bytesPerLine; i++ {
		// Affichage d'une ligne
		notEmpty := false
		line := ""

		// affichage de l'adresse de base de la ligne
		line += fmt.Sprintf("%08x ", i*bytesPerLine)
		for j := 0; j < bytesPerLine; j++ {
			ind := bytesPerLine*i + j // indice du byte à afficher
			val := hex.EncodeToString(GetArenaValues(ind, 1))
			line += " " + val

			if GetArenaValue(ind) != 0 {
				// valeur non nulle, on affiche la ligne
				notEmpty = true
				if j == bytesPerLine-1 {
					lastBytePrinted = true
				} else {
					lastBytePrinted = false
				}
			}
		}

		if notEmpty {
			if lastPrintedAdr != i-1 {
				for j := lastPrintedAdr + 1; j < i; j++ {
					PrintEmptyLine(j, bytesPerLine)
				}
			}
			fmt.Println(line)
			lastPrintedAdr = i
		}
	}
	if lastPrintedAdr == -1 {
		PrintEmptyLine(0, bytesPerLine)
	}

	if lastPrintedAdr != cw.MEM_SIZE/bytesPerLine-1 {
		if lastBytePrinted {
			PrintEmptyLine(lastPrintedAdr+1, bytesPerLine)
		}
		fmt.Println("...")
	}
}

func PrintEmptyLine(adr int, bytesPerLine int) {
	line := ""
	line += fmt.Sprintf("%08x ", adr*bytesPerLine)
	for k := 0; k < bytesPerLine; k++ {
		line += " 00"
	}
	fmt.Println(line)
}

// affichage des processus
func PrintProcesses() {
	fmt.Println("Id |Player Id |Pc   |Carry |Instr  |Wait |Registers           ")
	for n, p := range Processes {
		if !p.IsAlive && CurrentCycle-lastCheckCycle == CycleToDie+1 {
			continue
		}
		PrintColor(p.Player)
		// magouille pour affichage demandé
		// cmd : ___ et remainingCycles : 0 quand une commande est exécutée
		cmd, cmdCycles := GetCmdDatas(p.LoadedCmd)
		if cmdCycles == p.RemainingCycles || p.RemainingCycles == 0 {
			cmd = "___"
			cmdCycles = 0
		} else {
			cmdCycles = p.RemainingCycles
		}
		fmt.Printf(" %d |%9d |%4d |%-5t |%-6s |%4d |", n, p.Player, p.Pc, p.Carry, cmd, cmdCycles)

		for i := 1; i <= cw.REG_NUMBER; i++ {
			x := p.Registers[i]
			if x == 0 {
				fmt.Printf("%2d:0 ", i)
			} else {
				fmt.Printf("%2d:%s ", i, strings.TrimLeft(hex.EncodeToString(IntToBytes(x, 4)), "0"))
			}
		}
		fmt.Print("\n")
	}
	PrintColor(0)
}

func PrintPlayers() {
	fmt.Println("Id |Last Live |Nb Live since last check")
	for _, p := range PlayerDatas {
		PrintColor(p.Id)
		// if !PerfectOutputMode {
		// 	fmt.Printf(Colors[p.Id]+" %d | %9d | %d\n", p.Id, p.LastLive, p.LiveSinceLastCheck)
		// } else {
		fmt.Printf(" %d | %8d | %d\n", p.Id, p.LastLive, p.LiveSinceLastCheck)
	}
	PrintColor(0)
}

func PrintColor(id int) {
	if !PerfectOutputMode {
		fmt.Print(Colors[id])
	}
}
