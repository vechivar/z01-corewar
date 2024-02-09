package corewar

import (
	"fmt"

	cw "corewar"
)

// voir la doc pour les buts des fonctions

func Live(proc *Process) int {
	proc.IsAlive = true
	id := BytesToInt(GetArenaValues(proc.Pc+1, cw.REG_SIZE))
	AliveCount++
	if -id >= 1 && id <= len(PlayerDatas) {
		PlayerDatas[-id-1].LiveSinceLastCheck++
		PlayerDatas[-id-1].LastLive = CurrentCycle
		LastAlive = -id
		// fmt.Fprintf(os.Stderr, "Cycle %d : player %d (%s) is alive\n", CurrentCycle, -id, PlayerDatas[-id-1].PlayerName)
		fmt.Printf("cycle %d: Player %d (%s) is alive\n", CurrentCycle, -id, PlayerDatas[-id-1].PlayerName)
	}
	return cw.REG_SIZE
}

func Ld(proc *Process) int {
	argsTypes := GetArgsTypes(GetArenaValue(proc.Pc + 1))

	if CheckExpectedArgsTypes(proc.LoadedCmd, argsTypes) {
		argsValues := GetArgumentsValues(argsTypes, false, proc.Pc+2)
		cmdType := proc.LoadedCmd == 2 // true quand ld, false quand lld
		_, val := CalculateArgValue(argsTypes[0], argsValues[0], *proc, cmdType)
		if CheckRegId(argsValues[1]) {
			proc.Registers[argsValues[1]] = val
			if cmdType {
				proc.Carry = val == 0
			}
		}
	}

	return 1 + ArgSize(argsTypes, false)
}

func St(proc *Process) int {
	argsTypes := GetArgsTypes(GetArenaValue(proc.Pc + 1))

	if CheckExpectedArgsTypes(proc.LoadedCmd, argsTypes) {
		argsValues := GetArgumentsValues(argsTypes, false, proc.Pc+2)

		var val int
		if CheckRegId(argsValues[0]) {
			val = proc.Registers[argsValues[0]]
		} else {
			return 1 + ArgSize(argsTypes, false)
		}

		if argsTypes[1] == 1 {
			if CheckRegId(argsValues[1]) {
				proc.Registers[argsValues[1]] = val
			}
		} else {
			for i, x := range IntToBytes(val, 4) {
				SetArenaValue(proc.Pc+argsValues[1]%cw.IDX_MOD+i, x, proc.Player)
			}
		}
	}

	return 1 + ArgSize(argsTypes, false)
}

func Zjmp(proc *Process) int {
	if proc.Carry {
		offset := BytesToInt(GetArenaValues(proc.Pc+1, 2))
		proc.Pc += offset % cw.IDX_MOD
		return -1
	} else {
		return 2
	}
}

func Ldi(proc *Process) int {
	argsTypes := GetArgsTypes(GetArenaValue(proc.Pc + 1))

	if CheckExpectedArgsTypes(proc.LoadedCmd, argsTypes) {
		argsValues := GetArgumentsValues(argsTypes, true, proc.Pc+2)

		invalid, val1 := CalculateArgValue(argsTypes[0], argsValues[0], *proc, true)
		if invalid {
			return 1 + ArgSize(argsTypes, true)
		}
		invalid, val2 := CalculateArgValue(argsTypes[1], argsValues[1], *proc, false)
		if invalid {
			return 1 + ArgSize(argsTypes, true)
		}

		if CheckRegId(argsValues[2]) {
			if proc.LoadedCmd == 10 {
				// ldi
				// proc.Registers[argsValues[2]] = BytesToInt(GetArenaValues(proc.Pc+(val1+val2)%cw.IDX_MOD, cw.REG_SIZE))
				proc.Registers[argsValues[2]] = BytesToInt(GetArenaValues(proc.Pc+(val1+val2)%cw.IDX_MOD, cw.REG_SIZE))
			} else {
				// lldi
				proc.Registers[argsValues[2]] = BytesToInt(GetArenaValues(proc.Pc+val1+val2, cw.REG_SIZE))
			}
		}
	}
	return 1 + ArgSize(argsTypes, true)
}

func Sti(proc *Process) int {
	argsTypes := GetArgsTypes(GetArenaValue(proc.Pc + 1))

	if CheckExpectedArgsTypes(proc.LoadedCmd, argsTypes) {
		argsValues := GetArgumentsValues(argsTypes, true, proc.Pc+2)

		if !CheckRegId(argsValues[0]) {
			return 1 + ArgSize(argsTypes, true)
		}

		invalid, val1 := CalculateArgValue(argsTypes[1], argsValues[1], *proc, false)
		if invalid {
			return 1 + ArgSize(argsTypes, true)
		}
		invalid, val2 := CalculateArgValue(argsTypes[2], argsValues[2], *proc, false)
		if invalid {
			return 1 + ArgSize(argsTypes, true)
		}

		baseAdr := proc.Pc + (val1+val2)%cw.IDX_MOD
		values := IntToBytes(proc.Registers[argsValues[0]], cw.REG_SIZE)
		for i := 0; i < cw.REG_SIZE; i++ {
			SetArenaValue(baseAdr+i, values[i], proc.Player)
		}
	}
	return 1 + ArgSize(argsTypes, true)
}

func Fork(proc *Process) int {
	arg := BytesToInt(GetArenaValues(proc.Pc+1, cw.IND_SIZE))
	var newPc int

	if proc.LoadedCmd == 12 {
		// fork
		newPc = (proc.Pc + arg%cw.IDX_MOD) % cw.MEM_SIZE
	} else {
		// lfork
		newPc = (proc.Pc + arg) % cw.MEM_SIZE
	}

	newProc := *proc
	newProc.Pc = newPc
	newProc.IsAlive = false
	newProc.LoadedCmd = GetArenaValue(newPc)                    // Charge la commande suivante
	_, newProc.RemainingCycles = GetCmdDatas(newProc.LoadedCmd) // Charge

	Processes = append(Processes, newProc)

	return 2
}

func Nop(proc *Process) int {
	argsTypes := GetArgsTypes(GetArenaValue(proc.Pc + 1))
	if CheckExpectedArgsTypes(proc.LoadedCmd, argsTypes) {
		argsValues := GetArgumentsValues(argsTypes, false, proc.Pc+2)
		CheckRegId(argsValues[0])
	}

	return 1 + ArgSize(argsTypes, false)
}

// Utilisé pour add et sub
func Arithmetical(proc *Process) int {
	argsTypes := GetArgsTypes(GetArenaValue(proc.Pc + 1))

	if CheckExpectedArgsTypes(proc.LoadedCmd, argsTypes) {
		argsValues := GetArgumentsValues(argsTypes, false, proc.Pc+2)
		for _, x := range argsValues {
			if !CheckRegId(x) {
				return 4
			}
		}
		switch proc.LoadedCmd {
		case 4: // add
			proc.Registers[argsValues[2]] = proc.Registers[argsValues[0]] + proc.Registers[argsValues[1]]
		case 5: // sub
			proc.Registers[argsValues[2]] = proc.Registers[argsValues[0]] - proc.Registers[argsValues[1]]
		}
		proc.Carry = proc.Registers[argsValues[2]] == 0
		return 4
	}

	return 1 + ArgSize(argsTypes, false)
}

// Utilisé pour and, or et xor
func Logical(proc *Process) int {
	argsTypes := GetArgsTypes(GetArenaValue(proc.Pc + 1))

	if CheckExpectedArgsTypes(proc.LoadedCmd, argsTypes) {
		argsValues := GetArgumentsValues(argsTypes, false, proc.Pc+2)

		invalid, val1 := CalculateArgValue(argsTypes[0], argsValues[0], *proc, false)
		if invalid {
			return 1 + ArgSize(argsTypes, false)
		}
		invalid, val2 := CalculateArgValue(argsTypes[1], argsValues[1], *proc, false)
		if invalid {
			return 1 + ArgSize(argsTypes, false)
		}

		if CheckRegId(argsValues[2]) {
			switch proc.LoadedCmd {
			case 6: // and
				proc.Registers[argsValues[2]] = val1 & val2
			case 7: // or
				proc.Registers[argsValues[2]] = val1 | val2
			case 8: // xor
				proc.Registers[argsValues[2]] = val1 ^ val2
			}
			proc.Carry = proc.Registers[argsValues[2]] == 0
		}
	}
	return 1 + ArgSize(argsTypes, false)
}
