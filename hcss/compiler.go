package hcss

import (
	"fmt"
	"sort"
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

	Normal     = 0
	Style      = 1
	Var        = 2
	Mix        = 4
	HugoTmp    = 8
	NL         = 16
	UNKNOWN    = 32
	COMMENT    = 64
	MediaQuery = 128
	LCOMMENT   = 256
	VMS        = "=:("
	VEnd       = ";\r\n"
)

//Parse compile hcss string to css string
func Parse(src string) *ParsedDataStorage {

	LineReplacer := strings.NewReplacer("\r\n", NLC, "\r", NLC, "\n", NLC)
	src = LineReplacer.Replace(src)
	ds := NewDataStorage()

	processing := true

	pos := 0
	//SelectStyles := make(map[string]string)

	if processing {
		if strings.HasPrefix(src, NLC) {
			ds.StringStorage = append(ds.StringStorage, &StringData{NLC, NL})
		}
		src = strings.TrimSpace(src[pos:])

		if strings.HasPrefix(src, HugoTmpBegin) {

			endIndex := strings.Index(src[len(HugoTmpBegin):], HugoTmpEnd) + len(HugoTmpEnd)
			ds.StringStorage = append(ds.StringStorage, &StringData{src[:endIndex], HugoTmp})
			src = src[endIndex:]
		} else if strings.HasPrefix(src, VMPrefix) {
			defname, contentIndex, which, err := VMNameStrip(src)
			if err != nil {
				fmt.Println(err)
				return ds
			}
			contentIndex++
			src = src[contentIndex:]
			if strings.Contains(defname, HugoTmpBegin) || strings.Contains(defname, HugoTmpEnd) {
				serr := fmt.Errorf("You cannot use hugo template in variable name.\r\nName: %s", defname)
				fmt.Println(serr)
			}

			if which&Var > 0 {
				//if Variable
				endIndex := strings.IndexAny(src, VEnd)
				if endIndex < 0 {
					serr := fmt.Errorf("end string of %s: ( ; or newLine) is missing", defname)
					fmt.Println(serr)
					processing = false
				} else {
					ds.Variables[defname] = NewVariable(defname, src[:endIndex], Normal)
					src = src[endIndex+1:]
				}
			} else {
				//if mixin
				argAreaEndIndex := strings.Index(src, RBEnd)
				if argAreaEndIndex < 0 {
					serr := fmt.Errorf("cannot purse argumant of %s", defname)
					fmt.Println(serr)
					return ds
				}

				mi := &MixIn{defname, nil, nil, ""}

				args := strings.Split(src[:argAreaEndIndex], COMMA)
				mi.ParamString = make([]string, 0, len(args))

				argNameMap := make(map[string]bool)

				for _, arg := range args {
					arg = strings.TrimSpace(arg)

					_, ae := argNameMap[arg]

					if ae {
						fmt.Printf("argument name %s  is duplicate", arg)
					}
					mi.ParamString = append(mi.ParamString, arg)
				}

				tmp := make([]string, len(mi.ParamString), len(mi.ParamString))
				copy(mi.ParamString, tmp)
				mi.SortedParamString = StrSorter(tmp)
				sort.Sort(mi.SortedParamString)
				src = src[argAreaEndIndex+1:]
				cb := strings.Index(src, CBBegin)
				ce := strings.Index(src, CBEnd)
				if cb < 0 || ce < 0 {
					serr := fmt.Errorf("mixin definitation Error: missing { or }")
					fmt.Println(serr)
					processing = false
				} else if cb < ce {
					mi.Content = src[cb+1 : ce]
				} else {
					serr := fmt.Errorf("} exist before {")
					fmt.Println(serr)
					processing = false
				}

			}
		} else {

			if len(src) < 1 {
				processing = false
			} else {

			}
		}
	}

	return ds
}

//VMNameStrip strip name of var or mixin
func VMNameStrip(src string) (string, int, int, int, error) {

	sepIndex := strings.IndexAny(src, VMS)

	if sepIndex < 0 {
		err := fmt.Errorf("cannot strip name - not found : or = or {. If You define variables, you should use : or =. ")
		return "", src, sepIndex, UNKNOWN, err
	}
	name := src[:sepIndex]
	if strings.HasPrefix(src[sepIndex:], RBBegin) {
		return "", name, sepIndex, Mix, nil
	}
	return "", name, sepIndex, Var, nil
}

//ParsedDataStorage is data of process in compiling hcss
type ParsedDataStorage struct {
	Variables     map[string]*Variable
	MixIns        map[string]*MixIn
	StringStorage []*StringData
}

//GetStorageStrings get strings which is stored in dsp
func (dsp *ParsedDataStorage) GetStorageStrings() string {

	var sb strings.Builder
	for _, sd := range dsp.StringStorage {
		sb.WriteString(sd.Content)
	}
	return sb.String()
}

//NewDataStorage create new Data Storage
func NewDataStorage() *ParsedDataStorage {
	return &ParsedDataStorage{make(map[string]*Variable, 0), make(map[string]*MixIn, 0), make([]*StringData, 0)}
}

//StringData storage string data for output
type StringData struct {
	Content     string
	ContentType int
}

//StyleInfo strange information for style
type StyleInfo struct {
	Content     string
	ContentType int
}

//Variable represents data about a variable
type Variable struct {
	Name           string
	Content        string
	CompiledString string
	ContentType    int
}

//NewVariable get new variable
func NewVariable(defname string, vcontent string, vtype int) *Variable {
	return &Variable{defname, vcontent, "", vtype}
}

//MixIn represents data about MixIn
type MixIn struct {
	Name              string
	ParamString       []string
	SortedParamString StrSorter
	Content           string
}

//AtRule represents AtRule
type AtRule struct {
	Identifier string
	Str1       string
	InBracket  string
}
