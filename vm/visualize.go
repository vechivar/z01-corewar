package corewar

import (
	cw "corewar"
	"encoding/hex"
	"fmt"
	"strings"
)

// Affiche l'état actuel du jeu
func Visualize() {
	fmt.Println("Processes :")
	PrintProcesses()
	fmt.Println("Players:")
	PrintPlayers()
	fmt.Println("Arena :")
	PrintMemory()
	fmt.Println()
}

// Affiche la mémoire
func PrintMemory() {
	var pcs []int

	// pcs des processus, à surligner
	for _, p := range Processes {
		pcs = append(pcs, p.Pc)
	}

	bytesPerLine := 32
	printingEmptyLines := false
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

// affichage des processus
func PrintProcesses() {
	fmt.Println("Id |Player Id |Pc   |Carry |Instr  |Wait |Registers")
	for n, p := range Processes {
		fmt.Print(Colors[p.Player])
		cmd, _ := GetCmdDatas(p.LoadedCmd)
		fmt.Printf(" %d |%9d |%4d |%5t |%6s |%4d |", n, p.Player, p.Pc, p.Carry, cmd, p.RemainingCycles)
		for i := 1; i <= cw.REG_NUMBER; i++ {
			x := p.Registers[i]
			if x == 0 {
				fmt.Printf(" %d:0", i)
			} else {
				fmt.Printf(" %d:%s", i, strings.TrimLeft(hex.EncodeToString(IntToBytes(x, 4)), "0"))
			}
		}
		fmt.Print("\n")
	}
	fmt.Print(Colors[0])
}

func PrintPlayers() {
	fmt.Println("Id |Last Live |Nb Live since last check")
	for _, p := range PlayerDatas {
		fmt.Printf(Colors[p.Id]+" %d | %9d | %d\n", p.Id, p.LastLive, p.LiveSinceLastCheck)
	}
	fmt.Print(Colors[0])
}
