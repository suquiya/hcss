package hcss

import (
	"fmt"
	"strings"
)

//Compile parse hcss string
const (
	CSSExt       = ".css"
	VMPrefix     = "$"
	HugoTmpBegin = "{{"
	HugoTmpEnd   = "}}"

	CBBegin = "{"
	CBEnd   = "}"

	RBBegin          = "("
	RBEnd            = ")"
	CommentLine      = "//"
	CommentAreaBegin = "/*"
	CommentAreaEnd   = "*/"
	AtSign           = "@"

	EQ    = "="
	COLON = ":"
	COMMA = ","

	NLC = "\r\n"

	Normal             = 0
	StyleStatementFlag = 1

	Type1 = 0
	Type2 = 1

	Var     = 2 //CallVar is 3
	DefVar  = 2
	CallVar = 3
	Mix     = 4 //CallMix is 5
	DefMix  = 4
	CallMix = 5
	HugoTmp = 8  //10000
	NewLine = 16 //100000
	UNKNOWN = 0
	COMMENT = 32  //LCOMMENT is 33
	RSet    = 64  //include call is 65
	AtR     = 128 //include call is 129
	ERROR   = 1024

	VMS        = "=:({;\r\n"
	CallVarEnd = ";{\r\n"
	DefVarSep  = ":="
)

//Parse parse hcss string
func Parse(src string) *ParsedDataStorage {

	LineReplacer := strings.NewReplacer("\r\n", NLC, "\r", NLC, "\n", NLC)
	src = LineReplacer.Replace(src)
	ds := NewDataStorage()

	processing := true

	//SelectStyles := make(map[string]string)

	if processing {
		ds, src = PartParse(src, ds, Normal)
		if len(src) < 1 {
			processing = false
		}
	}

	return ds
}

//PartParse parse part of hcss string
func PartParse(src string, pds *ParsedDataStorage, Cond int) (*ParsedDataStorage, string) {
	if strings.HasPrefix(src, NLC) {
		pds.Statements = append(pds.Statements, NewLineString(NLC))
	}
	src = strings.TrimSpace(src)

	if strings.HasPrefix(src, HugoTmpBegin) {

		endIndex := strings.Index(src[len(HugoTmpBegin):], HugoTmpEnd) + len(HugoTmpEnd)
		pds.Statements = append(pds.Statements, HugoTemplate(src[:endIndex]))
		src = src[endIndex:]
	} else if strings.HasPrefix(src, VMPrefix) {
		s, err := VMParse(src, -1, pds)
		if err != nil {
			fmt.Println(err)
		}
		pds.Statements = append(pds.Statements, s)
	} else {

	}

	return pds, src
}

//VMParse evaluate and parse
func VMParse(src string, DCType int, pds *ParsedDataStorage) (ContentTyper, error) {
	sepIndex := strings.IndexAny(src, VMS)

	if sepIndex < 0 {
		name := src
		c, v := pds.Variables[name]
		if v {
			return c, nil
		}

		err := fmt.Errorf("Variable %s is not defined", name)
		/*
			fmt.Println(err)
		*/
		return InvalidStatement(src), err
	}

	sep := string(src[sepIndex])

	name := strings.TrimSpace(src[:sepIndex])
	src = strings.TrimSpace(src[sepIndex+1:])

	dcType := DCType
	if dcType < 0 {
		//Not clear def or call
		_, existVariableName := pds.Variables[name]
		_, existMixInName := pds.MixIns[name]
		if existVariableName {
			if existMixInName {
				if sep == RBBegin {
					//call Mixin
					dcType = CallMix
				} else if strings.Contains(CallVarEnd, sep) {
					dcType = CallVar
				} else {
					if strings.Contains(DefVarSep, sep) {
						dcType = DefVar
					} else {
						err := fmt.Errorf("invalid sep: %s", sep)
						fmt.Println(err)

						panic(err)

					}
				}
			} else {
				if strings.Contains(DefVarSep, sep) {
					dcType = DefVar
				} else {
					dcType = CallVar
				}
			}
		}

	}

	return nil, fmt.Errorf("Cannot parse as Variable or MixIn! ")
}

//ParsedDataStorage is data of process in compiling hcss
type ParsedDataStorage struct {
	Variables  map[string]*Variable
	MixIns     map[string]*MixIn
	Statements []ContentTyper
}

//GetStorageStrings get strings which is stored in dsp
/*
func (dsp *ParsedDataStorage) GetStorageStrings() string {

	var sb strings.Builder
	for _, sd := range dsp.StringStorage {
		sb.WriteString(sd.Content)
	}
	return sb.String()
}
*/

//NewDataStorage create new Data Storage
func NewDataStorage() *ParsedDataStorage {
	return &ParsedDataStorage{make(map[string]*Variable, 0), make(map[string]*MixIn, 0), make([]ContentTyper, 0)}
}

//ContentTyper is interface for distinguish ContentType
type ContentTyper interface {
	ContentType() int
}

//StyleStatement strange information for style
type StyleStatement struct {
	Content     string
	ContentType int
}

//AtRule represents AtRule
type AtRule struct {
	Identifier string
	Str1       string
	InBracket  string
}

//ContentType return *AtRule a's ContentType
func (a *AtRule) ContentType() int {
	return AtR
}

//ContentString is simple string.
type ContentString string

//ContentType return ContentString cs's ContentType
func (cs ContentString) ContentType() int {
	return Normal
}

//NewLineString is new Line
type NewLineString string

//ContentType return typemark NewLine
func (nls NewLineString) ContentType() int {
	return NewLine
}

//HugoTemplate is another name of string for storage hugo template string
type HugoTemplate string

//ContentType retrun HugoTmplate ht's ContentType
func (ht HugoTemplate) ContentType() int {
	return HugoTmp
}

//InvalidStatement represents invalid statement
type InvalidStatement string

//ContentType of InvalidStatement return ERROR Type.
func (is InvalidStatement) ContentType() int {
	return ERROR
}
