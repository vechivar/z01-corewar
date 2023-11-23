package corewar

import (
	cw "corewar"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadInput() {
	i := 1
	var err error

	var playerFiles []string

	// lecture des arguments du fichier
	for i < len(os.Args) {
		if os.Args[i] == "-d" {
			i++
			if i == len(os.Args) {
				WrongInput()
			}
			MaxCycle, err = strconv.Atoi(os.Args[i])
			if err != nil {
				WrongInput()
			}
			if MaxCycle < 0 {
				Exit("Use a positive number for -d argument.")
			}
		} else if os.Args[i] == "-v" {
			VisualMode = true
		} else {
			playerFiles = append(playerFiles, os.Args[i])
			if strings.Split(os.Args[i], ".")[len(strings.Split(os.Args[i], "."))-1] != "cor" {
				Exit("Invalid file format. Only .cor files are allowed")
			}
		}
		i++
	}

	if len(playerFiles) > cw.MAX_PLAYERS {
		Exit("Too much players. Max number of players is " + strconv.Itoa(cw.MAX_PLAYERS) + ".")
	}

	if len(playerFiles) == 0 {
		WrongInput()
	}

	nbPlayers := len(playerFiles)

	for _, x := range playerFiles {
		LoadPlayer(x, nbPlayers)
	}
}

// Initialise un joueur à partir du nom de son programme
func LoadPlayer(fileName string, nbPlayers int) {
	var playerData PlayerData

	data, err := os.ReadFile(fileName)
	if err != nil {
		Exit("Can't open file \"" + fileName + "\".")
	}

	if len(data) < 2192 {
		Exit("Corrupted file \"" + fileName + "\": file is too small.")
	}

	if data[0] != 0x00 || data[1] != 0xea || data[2] != 0x83 || data[3] != 0xf3 {
		Exit("Corrupted file \"" + fileName + "\": wrong signature.")
	}

	playerData.ProgramSize = BytesToInt(data[136:140])
	if len(data)-2192 != playerData.ProgramSize {
		Exit("Corrupted file \"" + fileName + "\": declared size doesn't match program size.")
	}
	if playerData.ProgramSize > cw.PLAYER_MAX_SIZE {
		Exit("Corrupted file \"" + fileName + "\": file is too big, maximum size allowed is " + strconv.Itoa(cw.PLAYER_MAX_SIZE) + " bytes.")
	}

	playerData.PlayerName = strings.TrimRight(string(data[4:132]), string(byte(0)))
	playerData.Description = strings.TrimRight(string(data[140:2188]), string(byte(0)))
	playerData.Id = len(PlayerDatas) + 1
	PlayerDatas = append(PlayerDatas, playerData)

	baseAdress := cw.MEM_SIZE * (playerData.Id - 1) / nbPlayers

	for i, b := range data[2192:] {
		SetArenaValue(baseAdress+i, b, playerData.Id)
	}

	var proc Process

	proc.Pc = baseAdress
	proc.Player = playerData.Id
	proc.Registers[1] = -proc.Player

	Processes = append(Processes, proc)
}

func WrongInput() {
	fmt.Println("------------------ VIRTUAL MACHINE ------------------")
	fmt.Println("Execute assembly code of multiple programs.")
	fmt.Println("./vm [-d nb_cycles] [-v] [player_exec.cor]...")
	fmt.Println("-d: dump memory at cycle [nb_cycles] and quit the program.")
	fmt.Println("-v: enable visual mode.")
	os.Exit(0)
}

func Exit(err string) {
	fmt.Println(err)
	os.Exit(0)
}
