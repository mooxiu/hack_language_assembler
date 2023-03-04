package main

import (
	"assembler/parser"
	"assembler/symbol_table"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var sTable symbol_table.SymbolTable
var emptyRAM int

func main() {
	// read a file
	filePath := os.Args[1]
	fileName := strings.Split(filePath, ".")[0]
	dat, err := os.ReadFile(filePath)
	if err != nil {
		os.Exit(1)
	}
	fileStr := string(dat)
	programSlice := strings.Split(fileStr, "\n")

	// pass one, deal with all the labels, white space and comment
	if len(sTable) == 0 {
		sTable = symbol_table.GetSymbolTable()
	}
	emptyRAM = 16
	unlabeledProgram := labelHandler(programSlice)
	fmt.Println(unlabeledProgram)
	fmt.Println(sTable)

	// pass two
	machineCodes := make([]string, 0)
	// TODO: suppose we have no symbols now
	for lineNum, line := range unlabeledProgram {
		line = strings.TrimSpace(line)
		if line[0] == '@' {
			// is A instruction
			machineCodes = append(machineCodes, handleAInstruction(lineNum, line))
			continue
		}
		// is C instruction
		machineCodes = append(machineCodes, handleCInstruction(lineNum, line))
	}

	// flatten machine codes into output
	sb := strings.Builder{}
	for _, machineCode := range machineCodes {
		sb.WriteString(machineCode)
		sb.WriteByte('\n')
	}

	// create a new file and write into it
	err = os.WriteFile(fileName+".hack", []byte(sb.String()), os.ModePerm)
	if err != nil {
		os.Exit(1)
	}
}

// TODO: deal with all of the labels
func labelHandler(program []string) []string {
	processed := make([]string, 0)
	count := 0
	for _, line := range program {
		trimmed := strings.TrimSpace(line)
		// skip blank
		if trimmed == "" {
			continue
		}
		if len(trimmed) < 2 {
			processed = append(processed, trimmed)
			count++
			continue
		}
		if trimmed[0:2] == "//" {
			continue
		}
		// skip comment
		idx := strings.Index(trimmed, "//")
		var cleanTrimmed string
		if idx >= 0 {
			cleanTrimmed = trimmed[:idx]
		} else {
			cleanTrimmed = trimmed
		}
		if Contains(cleanTrimmed, "(") {
			label := strings.TrimSuffix(strings.TrimPrefix(cleanTrimmed, "("), ")")
			sTable[label] = count
		} else {
			processed = append(processed, cleanTrimmed)
			count++
		}
	}
	return processed
}

func handleAInstruction(lineNum int, line string) string {
	var sb strings.Builder
	sb.WriteByte('0')
	var p *parser.AParser
	parsed := p.Parse(strings.TrimSpace(line))
	// if addr is not number
	eval := ""
	if ContainsNonNum(parsed.Addr) {
		v, ok := sTable[parsed.Addr]
		if ok {
			eval = strconv.Itoa(v)
		} else {
			eval = strconv.Itoa(emptyRAM)
			sTable[parsed.Addr] = emptyRAM
			emptyRAM += 1
		}
	} else {
		eval = parsed.Addr
	}

	// if addr is number
	binaryAddr, err := strconv.ParseInt(eval, 10, 64)
	if err != nil {
		panic(err)
	}
	binaryAddrStr := fmt.Sprintf("%b", binaryAddr)
	for i := 0; i < 15-len(binaryAddrStr); i++ {
		sb.WriteByte('0')
	}
	sb.WriteString(binaryAddrStr)
	return sb.String()
}

func handleCInstruction(lineNum int, line string) string {
	var p *parser.CParser
	line = strings.TrimSpace(line)
	parsed := p.Parse(line)
	var sb strings.Builder
	// start with 1
	sb.WriteString("111")
	// write comp part
	sb.WriteString(handleCInstructionComp(parsed.Comp))
	// write dest part
	sb.WriteString(handleCInstructionDest(parsed.Dest))
	// write jump
	sb.WriteString(handleCInstructionJump(parsed.Jump))
	return sb.String()
}

func handleCInstructionComp(comp string) string {
	var s string
	switch comp {
	case "0":
		s = "0101010"
	case "1":
		s = "0111111"
	case "-1":
		s = "0111010"
	case "D":
		s = "0001100"
	case "A":
		s = "0110000"
	case "!D":
		s = "0001101"
	case "!A":
		s = "0110001"
	case "-D":
		s = "0001111"
	case "-A":
		s = "0110011"
	case "D+1":
		s = "0011111"
	case "A+1":
		s = "0110111"
	case "D-1":
		s = "0001110"
	case "A-1":
		s = "0110010"
	case "D+A":
		s = "0000010"
	case "D-A":
		s = "0010011"
	case "A-D":
		s = "0000111"
	case "D&A":
		s = "0000000"
	case "D|A":
		s = "0010101"
	case "M":
		s = "1110000"
	case "!M":
		s = "1110001"
	case "-M":
		s = "1110011"
	case "M+1":
		s = "1110111"
	case "M-1":
		s = "1110010"
	case "D+M":
		s = "1000010"
	case "D-M":
		s = "1010011"
	case "M-D":
		s = "1000111"
	case "D&M":
		s = "1000000"
	case "D|M":
		s = "1010101"
	default:
		s = "0000000"
	}
	return s
}

func handleCInstructionDest(dest string) string {
	var s string
	switch dest {
	case "":
		s = "000"
	case "M":
		s = "001"
	case "D":
		s = "010"
	case "DM", "MD":
		s = "011"
	case "A":
		s = "100"
	case "AM", "MA":
		s = "101"
	case "AD", "DA":
		s = "110"
	case "ADM", "AMD", "DMA", "DAM", "MAD", "MDA":
		s = "111"
	default:
		s = "000"
	}
	return s
}

func handleCInstructionJump(jump string) string {
	var s string
	switch jump {
	case "":
		s = "000"
	case "JGT":
		s = "001"
	case "JEQ":
		s = "010"
	case "JGE":
		s = "011"
	case "JLT":
		s = "100"
	case "JNE":
		s = "101"
	case "JLE":
		s = "110"
	case "JMP":
		s = "111"
	default:
		s = "000"
	}
	return s
}

func Contains(s string, sub string) bool {
	return strings.Index(s, sub) >= 0
}

func ContainsNonNum(s string) bool {
	for i := range s {
		if s[i] < '0' || s[i] > '9' {
			return true
		}
	}
	return false
}
