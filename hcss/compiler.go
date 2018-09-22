package hcss

import (
	"strings"
)

//Compile parse hcss string
const (
	CSSExt       = ".css"
	VMStr        = '$'
	HugoTmpBegin = "{{"
	HugoTmpEnd   = "}}"

	CBBegin = '{'
	CBEnd   = '}'

	RBBegin = '('
	RBEnd   = ')'

	EQ    = '='
	COLON = ':'

	NLC = "\r\n"

	Normal     = 1
	InVar      = 2
	InMixIn    = 4
	InHugoTmp  = 8
	VarOrMixIn = 2 & 4
)

//Compile compile hcss string to css string
func Compile(src string) string {

	LineReplacer := strings.NewReplacer("\r\n", NLC, "\r", NLC, "\n", NLC)
	src = LineReplacer.Replace(src)
	p := NewProcess()

	nowStat := Normal
	IsPrevBrackets := false
	processing := true
	srcLetters := []rune(src)

	pos := 0
	if processing {
		srcLetter := srcLetters[pos]
		if nowStat == Normal {
			if srcLetter == CBBegin {
				if IsPrevBrackets {
					nowStat = InHugoTmp
					endIndex := strings.Index(src, HugoTmpEnd) + len(HugoTmpEnd)
					p.buf.WriteString(src[pos:endIndex])
					pos = endIndex
				}
			}
		}
		pos++
	}

	return p.buf.String()
}

//CProcess is data of process in compiling hcss
type CProcess struct {
	Variables map[string]*Variable
	MixIns    map[string]*MixIn
	buf       *strings.Builder
}

//NewProcess create new Process
func NewProcess() *CProcess {

	var sb strings.Builder

	return &CProcess{make(map[string]*Variable, 0), make(map[string]*MixIn, 0), &sb}
}

//Variable represents data about a variable
type Variable struct {
	Content string
	Type    string
}

//MixIn represents data about MixIn
type MixIn struct {
	paramString []string
	content     string
}
