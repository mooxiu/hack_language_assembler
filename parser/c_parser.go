package parser

import (
	"strings"
)

type CParser struct {
}

type CParsed struct {
	Dest string
	Comp string
	Jump string
}

func (c *CParser) Parse(cInstruction string) CParsed {
	if !contains(cInstruction, "=") && !contains(cInstruction, ";") {
		return CParsed{
			Dest: "",
			Comp: "",
			Jump: "",
		}
	}
	if contains(cInstruction, "=") && contains(cInstruction, ";") {
		temp := strings.Split(cInstruction, "=")
		dest := temp[0]
		temp2 := strings.Split(temp[1], ";")
		comp, jump := temp2[0], temp2[1]
		return CParsed{
			Dest: dest,
			Comp: comp,
			Jump: jump,
		}
	}
	if contains(cInstruction, "=") {
		// only dest and comp
		temp := strings.Split(cInstruction, "=")
		return CParsed{
			Dest: temp[0],
			Comp: temp[1],
			Jump: "",
		}
	}
	if contains(cInstruction, ";") {
		// only comp and jump
		temp := strings.Split(cInstruction, ";")
		return CParsed{
			Dest: "",
			Comp: temp[0],
			Jump: temp[1],
		}
	}
	return CParsed{}
}

func contains(s string, sub string) bool {
	return strings.Index(s, sub) >= 0
}
