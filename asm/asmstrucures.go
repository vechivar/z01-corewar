package corewar

import (
	"bufio"
	"os"
)

// Commande assembleur
type asmCmd struct {
	cmd       string   // Nom de la commande
	args      []asmArg // Liste des arguments
	line      int      // Ligne correspondant dans le fichier d'entrée
	byteCount int      // Position de la commande dans le binaire finale
}

// Argument assembleur
type asmArg struct {
	argType  byte   // Type d'argument. 'r' registre, 'd' direct, 'i' indirect
	value    int    // Valeur associée
	label    string // Nom du label. Vide si l'argument n'en utilise pas
	byteSize int    // Taille en bytes de l'argument
}

// Label assembleur
type label struct {
	value int    // Position du label dans le binaire
	label string // Nom du label
}

// Variables du programme

var outputFile *os.File      // Fichier de sortie
var inputScan *bufio.Scanner // Scanner du fichier de sortie

var lineCount int // Indice de la ligne actuelle
var byteCount = 0 // Nombre de bytes utilisés par les commandes précédentes. Contient la taille totale à la fin du programme.

var labels []label     // Tous les labels déclarés dans le programme.
var usedLabels []label // Tous les labels utilisés (value contient la ligne d'apparition)
var commands []asmCmd  // Liste des commandes
