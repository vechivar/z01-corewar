package corewar

import (
	cw "corewar"
)

// Processus en train d'être exécuté
type Process struct {
	Pc              int                    // pc actuel
	Registers       [cw.REG_NUMBER + 1]int // valeurs des registres
	Player          int                    // joueur ayant créé le processus
	RemainingCycles int                    // nombre de cycles avant d'excéuter la commande chargée
	IsAlive         bool                   // le processus a exécuté un live depuis le dernier check
	LoadedCmd       byte                   // commande actuellement chargée
	Carry           bool                   // carry du processus
}

// Joueur
type PlayerData struct {
	Id                 int    // id du joueur
	PlayerName         string // nom du joueur
	Description        string // description du joueur
	ProgramSize        int    // taille du joueur
	LastLive           int    // dernier cycle où le joueur a reçu un live
	LiveSinceLastCheck int    // nombre de lives reçus depuis le dernier check
}

// Arène
type Arena_ struct {
	Memory       [cw.MEM_SIZE]byte // valeurs de la mémoire
	LastModified [cw.MEM_SIZE]int  // id du joueur ayant modifié la valeur en dernier (pour l'impression)
}

var PlayerDatas []PlayerData // liste des joueurs
var MaxCycle int             // argument -d du programme
var VisualMode = false       // argument -v du programme

var Arena Arena_        // arène
var Processes []Process // liste des processus en cours d'exécution

var LastAlive int      // id du dernier joueur à avoir reçu un live
var CurrentCycle int   // cycle actuel
var CycleToDie int     // nombre de cycles entre deux check d'élimination des processus
var LifeCheckCount int // nombre de check exécutés sans avoir baissé CycleToDie
var AliveCount int     // nombre de live exécutés depuis le dernier check

var Colors []string = []string{"\033[37m", "\033[31m", "\033[32m", "\033[33m", "\033[34m"} // couleurs à utiliser pour les joueurs
