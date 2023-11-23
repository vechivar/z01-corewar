package corewar

import (
	"bufio"
	cw "corewar"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Fonction principale.
// Lit le fichier d'entrée ligne par ligne.
// Les lignes correspondant à du code sont traduites en binaire et ajoutées au fichier de sortie.
func ReadInput() {
	inputPath := ReadParams()

	// Ouverture du fichier d'entrée
	inputFile, err := os.Open(inputPath)
	if err != nil {
		Exit("Can't open input file")
	}

	OpenOutputPath(inputPath)

	inputScan = bufio.NewScanner(inputFile)
	inputScan.Split(bufio.ScanLines)

	ProcessHeader()

	byteCount = 0

	flag, line := GetNextLine()
	for !flag {
		ProcessLine(line)
		flag, line = GetNextLine()
	}

	CheckUsedLabels()
	ProcessCommands()

	if byteCount > cw.PLAYER_MAX_SIZE {
		Exit("Error : program is too long. Max byte size is " + strconv.Itoa(cw.PLAYER_MAX_SIZE))
	}
	// Ecriture de la taille du programme
	_, err = outputFile.WriteAt([]byte{byte(byteCount / 256), byte(byteCount % 256)}, 138)
	if err != nil {
		fmt.Println(err)
	}

	inputFile.Close()
	outputFile.Close()
}

// Lecture des paramètres du programme. Renvoie le nom du fichier d'entrée.
func ReadParams() string {
	if len(os.Args) != 2 {
		Exit("Wrong number of arguments")
	}

	argExt := strings.Split(os.Args[1], ".")
	if argExt[len(argExt)-1] != "s" {
		Exit("Invalid input extension. Only .s files are allowed.")
	}

	return os.Args[1]
}

// Crée le fichier de sortie.
// Ecrit la signature en tête du fichier.
func OpenOutputPath(inputPath string) {
	splitTemp := strings.Split(inputPath, "/")
	outputPath := splitTemp[len(splitTemp)-1]
	splitTemp = strings.Split(outputPath, ".")
	outputPath = strings.Join(splitTemp[:len(splitTemp)-1], ".") + ".cor"

	var err error
	outputFile, err = os.Create(outputPath)
	if err != nil {
		Exit("Can't create output file : " + outputPath)
	}

	sign := []byte{0x00, 0xea, 0x83, 0xf3}

	outputFile.Write(sign)
}

// Traite le nom et la description du programme.
func ProcessHeader() {
	// Lecture et écriture du nom du programme
	x, nameLine := GetNextLine()
	if x {
		Exit("Invalid file. Use " + cw.NAME_CMD_STRING + " first to specify a name")
	}

	if strings.Split(nameLine, " ")[0] != cw.NAME_CMD_STRING {
		Exit("Invalid file. Use " + cw.NAME_CMD_STRING + " first to specify a name")
	}

	name := strings.Trim(strings.Join(strings.Split(nameLine, " ")[1:], " "), "\"")

	n, _ := outputFile.WriteString(name)

	if n > cw.PROG_NAME_LENGTH {
		Exit("Program name is too long")
	}

	// rq : on réserve 4 bytes supplémentaires pour écrire la taille du programme
	// on réserve aussi 4 bytes supplémentaires parce que le sujet n'est pas cohérent avec les exemples
	for i := 0; i < cw.PROG_NAME_LENGTH-n+8; i++ {
		outputFile.Write([]byte{0x00})
	}

	// Lecture et écriture de la description du programme
	x, descLine := GetNextLine()
	if x {
		Exit("Invalid file. Use " + cw.DESCRIPTION_CMD_STRING + " second to specify a description")
	}

	if strings.Split(descLine, " ")[0] != cw.DESCRIPTION_CMD_STRING {
		Exit("Invalid file. Use " + cw.DESCRIPTION_CMD_STRING + " second to specify a description")
	}

	desc := strings.Trim(strings.Join(strings.Split(descLine, " ")[1:], " "), "\"")

	n, _ = outputFile.WriteString(desc)

	if n > cw.DESCRIPTION_LENGTH {
		Exit("Program description is too long")
	}

	for i := 0; i < cw.DESCRIPTION_LENGTH-n+4; i++ {
		outputFile.Write([]byte{0x00})
	}
}

// Renvoie la prochaine ligne à traiter du programme (ignore les commentaires et les lignes vides)
// Renvoie true si la lecture du fichier est terminée
func GetNextLine() (bool, string) {
	for inputScan.Scan() {
		lineCount++
		line := strings.Trim(strings.ReplaceAll(inputScan.Text(), "\t", " "), " ")
		if len(line) > 0 && rune(line[0]) != cw.COMMENT_CHAR {
			return false, line
		}
	}
	return true, ""
}

// Affiche l'erreur et quitte le programme.
// Supprime le fichier de sortie s'il a été créé.
func Exit(err string) {
	fmt.Println(err)
	if outputFile != nil {
		os.Remove(outputFile.Name())
	}
	os.Exit(0)
}

// Affiche l'erreur à une ligne donnée et quitte le programme.
// Supprime le fichier de sortie s'il a été créé.
func ExitAtLine(err string, line int) {
	fmt.Printf("Error on line %d : ", line)
	fmt.Println(err)
	if outputFile != nil {
		os.Remove(outputFile.Name())
	}
	os.Exit(0)
}
