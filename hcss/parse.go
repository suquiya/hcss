package hcss

import (
	"fmt"
	"strings"
	"unicode"
)

//Compile parse hcss string
const (
	CSSExt       = ".css"
	HCSSExt      = ".hcss"
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

	EQ        = "="
	COLON     = ":"
	SEMICOLON = ";"
	COMMA     = ","

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

	VMS            = "=:({;\r\n"
	CallEnd        = ";{\r\n"
	CallMixInBegin = "(:="
	DefVarSep      = ":="
	DefEnd         = ";\r\n"
)

//Parse parse hcss string
func Parse(src string) *ParsedData {

	LineReplacer := strings.NewReplacer("\r\n", NLC, "\r", NLC, "\n", NLC)
	src = LineReplacer.Replace(src)
	ds := NewDataStorage()

	processing := true

	//SelectStyles := make(map[string]string)

	src = strings.TrimSpace(src)

	if processing {
		ds, src = PartParse(src, ds, Normal)
		if len(src) < 1 {
			processing = false
		}
	}

	return ds
}

//PartParse parse part of hcss string
func PartParse(src string, pds *ParsedData, Cond int) (*ParsedData, string) {
	if strings.HasPrefix(src, NLC) {
		pds.Statements = append(pds.Statements, NewLineString(NLC))
	}
	src = strings.TrimLeftFunc(src, unicode.IsSpace)

	if strings.HasPrefix(src, HugoTmpBegin) {

		endIndex := strings.Index(src[len(HugoTmpBegin):], HugoTmpEnd)
		if endIndex < 0 {
			err := fmt.Errorf("Hugo Template not closed")
			fmt.Println(err)
			src = ""
		} else {
			endIndex += len(HugoTmpEnd)

			pds.Statements = append(pds.Statements, HugoTemplate(src[:endIndex]))
			src = src[endIndex:]
		}
	} else if strings.HasPrefix(src, VMPrefix) {
		var s ContentTyper
		var err error
		s, src, err = VMParse(src, -1, pds)
		if err != nil {
			fmt.Println(err)
		}
		pds.Statements = append(pds.Statements, s)
	} else {

	}

	return pds, src
}

//VMParse evaluate and parse
func VMParse(src string, DCType int, pds *ParsedData) (ContentTyper, string, error) {
	sepIndex := strings.IndexAny(src, VMS)

	if sepIndex < 0 {
		name := src
		c, v := pds.Variables[name]
		if v {
			return c, "", nil
		}

		err := fmt.Errorf("Variable %s is not defined", name)
		/*
			fmt.Println(err)
		*/
		return InvalidStatement(src), "", err
	}

	sep := string(src[sepIndex])

	name := strings.TrimSpace(src[:sepIndex])
	src = strings.TrimLeftFunc(src[sepIndex+1:], unicode.IsSpace)

	dcType := DCType
	if dcType < 0 {
		//Not clear def or call
		_, existVariableName := pds.Variables[name]
		mi, existMixInName := pds.MixIns[name]
		if existVariableName {
			if existMixInName {
				if sep == RBBegin {
					//call Mixin
					var err error
					err = nil
					endIndex := strings.Index(src, RBEnd)
					mica := NewMixInCallArgs(name)
					if endIndex < 0 {
						err = fmt.Errorf("() is not closed")
						fmt.Println(err)
						return mica, src, err
					}

					args := src[:endIndex]
					src = strings.TrimLeftFunc(src, IsNotNewLineSpace)
					var content string
					content = ""
					if HasPrefixOfMixInContentBegin(src) {
						endContent := strings.Index(src, CBEnd)
						if endContent < 0 {
							err = fmt.Errorf("MixIn's Content missing")
							fmt.Println(err)
						} else {
							content = src[1:endContent]
							src = src[endContent+1:]
						}
					}
					mica, err = mica.ArgsParse(mi, args, content)

					src = strings.TrimLeftFunc(src, unicode.IsSpace)
					if strings.HasPrefix(src, SEMICOLON) {
						src = src[1:]
					}
					return mica, src, err
				} else if strings.Contains(DefVarSep, sep) {
					dcType = CallMix

					mica := NewMixInCallArgs(name)
					argsEndIndex := strings.IndexAny(src, CallEnd)
					argString := src[:argsEndIndex]

					src = strings.TrimLeftFunc(src, IsNotNewLineSpace)

					content := ""
					var err error
					if HasPrefixOfMixInContentBegin(src) {
						endContent := strings.Index(src, CBEnd)

						if endContent < 0 {
							err = fmt.Errorf("MixIn's Content Missing")
							fmt.Println(err)
						} else {
							content = src[1:endContent]
							src = src[endContent+1:]
						}
					}
					mica, err = mica.ArgsParse(mi, argString, content)
					src = strings.TrimLeftFunc(src, unicode.IsSpace)
					if strings.HasPrefix(src, SEMICOLON) {
						src = src[1:]
					}
					return mica, src, err

				} else if strings.Contains(CallEnd, sep) {
					if sep == CBBegin {
						mica := NewMixInCallArgs(name)
						endContent := strings.Index(src, CBEnd)
						var err error
						var content string
						if endContent < 0 {
							err = fmt.Errorf("MixIn's Content Missing")
							fmt.Println(err)
						} else {
							content = src[1:endContent]
							src = src[endContent+1:]
						}
						mica, err = mica.ArgsParse(mi, "", content)
						src = strings.TrimLeftFunc(src, unicode.IsSpace)
						if strings.HasPrefix(src, SEMICOLON) {
							src = src[1:]
						}
						return mica, src, err
					}
					dcType = CallVar
					return VarCall(name), src, nil
				} else {
					if strings.Contains(DefVarSep, sep) {
						dcType = DefVar
						v, trimedString, err := GetVariableDefContentString(src)

						src = trimedString
						return NewSimpleVariable(name, v), src, err
					}
					err := fmt.Errorf("invalid sep: %s", sep)
					fmt.Println(err)
					panic(err)
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

	return nil, src, fmt.Errorf("Cannot parse as Variable or MixIn! ")
}

//IsNotNewLineSpace is function for space
func IsNotNewLineSpace(r rune) bool {
	if unicode.IsSpace(r) {
		return !strings.ContainsRune(NLC, r)
	}
	return false
}

//ParsedData is data of process in compiling hcss
type ParsedData struct {
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
func NewDataStorage() *ParsedData {
	return &ParsedData{make(map[string]*Variable, 0), make(map[string]*MixIn, 0), make([]ContentTyper, 0)}
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
