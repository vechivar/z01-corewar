package corewar

import (
	"fmt"
	"os"
	"strings"

	cw "corewar"
)

var (
	debugMode  = false // true pour afficher l'état quand une commande est exécutée
	printCycle = false
)

// Fait tourner la machine virtuelle
func RunVm() {
	InitVm()

	CurrentCycle = 1
	lastCheckCycle = 0

	for MaxCycle == 0 || CurrentCycle < MaxCycle {
		RunCycle()
		if VisualMode || (debugMode && printCycle) {
			fmt.Printf("Cycle %d || Cycles before life check: %d || Cycles between checks: %d", CurrentCycle, CycleToDie-(CurrentCycle-lastCheckCycle-1), CycleToDie)
			if LastAlive != 0 {
				fmt.Print(" || ")
				PrintColor(LastAlive)
				fmt.Printf("Last alive: %d", LastAlive)
				PrintColor(0)
			}
			fmt.Print("\n")
			Visualize()
		}
		if CurrentCycle-lastCheckCycle == CycleToDie+1 {
			LiveCheck()
			lastCheckCycle = CurrentCycle
		}
		CurrentCycle++
	}

	EndGame()
}

// Exécute la commande du processus et renvoie la taille des arguments utilisés
func RunProcess(proc *Process) int {
	switch proc.LoadedCmd {
	case 1: // live
		return Live(proc)
	case 2: // ld
		return Ld(proc)
	case 3: // st
		return St(proc)
	case 4: // add
		return Arithmetical(proc)
	case 5: // sub
		return Arithmetical(proc)
	case 6: // and
		return Logical(proc)
	case 7: // or
		return Logical(proc)
	case 8: // xor
		return Logical(proc)
	case 9: // zjmp
		return Zjmp(proc)
	case 10: // ldi
		return Ldi(proc)
	case 11: // sti
		return Sti(proc)
	case 12:
		return Fork(proc) // fork
	case 13: // lld
		return Ld(proc)
	case 14: // lldi
		return Ldi(proc)
	case 15:
		return Fork(proc) // lfork
	case 16:
		return Nop(proc) // nop
	default:
		fmt.Fprintf(os.Stderr, "cycle %d: Opcode %d is not a valid instruction\n", CurrentCycle, proc.LoadedCmd)
		return 0
	}
}

// Vérifie si des processus sont encore en vie.
// Met fin à la partie selon les règles du jeu
func LiveCheck() {
	var liveProc []Process

	if LastAlive == 0 { // aucun joueur n'a effectué de live
		EndGame()
	} else {
		for i := 0; i < len(PlayerDatas); i++ {
			PlayerDatas[i].LiveSinceLastCheck = 0
		}
	}

	for i := 0; i < len(Processes); i++ {
		if Processes[i].IsAlive {
			Processes[i].IsAlive = false
			liveProc = append(liveProc, Processes[i])
		}
	}

	if len(liveProc) == 0 { // aucun processus n'a exécuté de live
		EndGame()
	}

	Processes = liveProc // mise à jour des processus en vie

	LifeCheckCount++

	if LifeCheckCount > cw.MAX_CHECKS || AliveCount >= cw.NBR_LIVE { // réduction de cycletodie selon les règles du jeu
		prev := CycleToDie
		CycleToDie -= cw.CYCLE_DELTA
		if CycleToDie < 0 {
			CycleToDie = 0
		}
		// fmt.Fprintf(os.Stderr, "cycle %d: Cycles to die decreased: %d -> %d\n", CurrentCycle, prev, CycleToDie)
		fmt.Printf("cycle %d: Cycles to die decreased: %d -> %d\n", CurrentCycle, prev, CycleToDie)

		LifeCheckCount = 0
	}

	AliveCount = 0
}

// Fin de partie
func EndGame() {
	if MaxCycle != 0 && !VisualMode {
		PrintMemory()
		fmt.Println()
	}
	fmt.Printf("cycle %d: ", CurrentCycle)
	if LastAlive == 0 {
		fmt.Print("Nobody wins!\n")
	} else {
		fmt.Printf("The winner is player %d: %s!\n", LastAlive, PlayerDatas[LastAlive-1].PlayerName)
	}
	os.Exit(0)
}

// Exécute un cycle de la machine virtuelle
// Décrémente le nombre de cycle de chaque processus.
// Exécute la commande si ce nombre tombe à 0, et charge la commande suivante
func RunCycle() {
	printCycle = false
	for i := len(Processes) - 1; i >= 0; i-- {
		Processes[i].RemainingCycles--
		if Processes[i].RemainingCycles == 0 {
			Processes[i].Pc += RunProcess(&Processes[i]) + 1 // Exécute la commande et met à jour le pc
			Processes[i].Pc = Processes[i].Pc % cw.MEM_SIZE
			Processes[i].LoadedCmd = GetArenaValue(Processes[i].Pc)               // Charge la commande suivante
			_, Processes[i].RemainingCycles = GetCmdDatas(Processes[i].LoadedCmd) // Charge le nombre de tours correspondant
			printCycle = true
		}
	}
}

// Initialise les valeurs
func InitVm() {
	InitMessage()

	for i := 0; i < len(PlayerDatas); i++ {
		PlayerDatas[i].LastLive = 0
		PlayerDatas[i].LiveSinceLastCheck = 0
	}

	for i := 0; i < len(Processes); i++ {
		_, Processes[i].RemainingCycles = GetCmdDatas(GetArenaValue(Processes[i].Pc))
		Processes[i].LoadedCmd = GetArenaValue(Processes[i].Pc)
		Processes[i].Carry = false
	}

	CycleToDie = cw.CYCLE_TO_DIE
	LifeCheckCount = 0
	AliveCount = 0
	LastAlive = 0
}

func InitMessage() {
	fmt.Println("For this match the players will be:")
	for _, p := range PlayerDatas {
		fmt.Printf("Player %d (%d bytes): %s (%s)\n", p.Id, p.ProgramSize, strings.Trim(p.PlayerName, string(byte(0))), strings.Trim(p.Description, string(byte(0))))
	}
}
