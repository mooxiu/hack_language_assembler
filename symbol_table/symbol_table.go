package symbol_table

type SymbolTable map[string]int

func (s SymbolTable) get(key string) int {
	return s[key]
}

func GetSymbolTable() SymbolTable {
	s := make(map[string]int)
	s["R0"] = 0
	s["R1"] = 1
	s["R2"] = 2
	s["R3"] = 3
	s["R4"] = 4
	s["R5"] = 5
	s["R6"] = 6
	s["R7"] = 7
	s["R8"] = 8
	s["R9"] = 9
	s["R10"] = 10
	s["R11"] = 11
	s["R12"] = 12
	s["R13"] = 13
	s["R14"] = 14
	s["R15"] = 15
	s["SCREEN"] = 16384
	s["KBD"] = 24576
	s["SP"] = 0
	s["LCL"] = 1
	s["ARG"] = 2
	s["THIS"] = 3
	s["THAT"] = 4
	return s
}
