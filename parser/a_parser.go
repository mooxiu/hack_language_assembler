package parser

type AParser struct {
}

type AParsed struct {
	AtMark string
	Addr   string
}

func (a *AParser) Parse(aInstruction string) AParsed {
	return AParsed{
		AtMark: string(aInstruction[0]),
		Addr:   aInstruction[1:],
	}
}
