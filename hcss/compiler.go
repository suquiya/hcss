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

	RBBegin = "("
	RBEnd   = ")"

	EQ    = "="
	COLON = ":"
	COMMA = ","

	NLC = "\r\n"

	Normal  = 0
	Var     = 1
	Mix     = 2
	HugoTmp = 4
	UNKNOWN = 16
	VMS     = "=:("
	VEnd    = ";\r\n"
)

//Compile compile hcss string to css string
func Compile(src string) string {

	LineReplacer := strings.NewReplacer("\r\n", NLC, "\r", NLC, "\n", NLC)
	src = LineReplacer.Replace(src)
	dsp := NewDataStorageForProcess()

	processing := true

	pos := 0
	//SelectStyles := make(map[string]string)
	src = strings.TrimSpace(src[pos:])
	if processing {
		if strings.HasPrefix(src, HugoTmpBegin) {

			endIndex := strings.Index(src[len(HugoTmpBegin):], HugoTmpEnd) + len(HugoTmpEnd)
			dsp.StringStorage = append(dsp.StringStorage, &StringData{src[:endIndex], HugoTmp})
			src = src[endIndex:]
		} else if strings.HasPrefix(src, VMPrefix) {
			defname, contentIndex, which, err := VMNameStrip(src)
			if err != nil {
				fmt.Println(err)
				return dsp.GetStorageStrings()
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
					dsp.Variables[defname] = NewVariable(src[:endIndex], Normal)
					src = src[endIndex+1:]
				}
			} else {
				//if mixin
				argAreaEndIndex := strings.Index(src, RBEnd)
				if argAreaEndIndex < 0 {
					serr := fmt.Errorf("cannot purse argumant of %s", defname)
					fmt.Println(serr)
					return dsp.GetStorageStrings()
				}

				args := strings.Split(src[:argAreaEndIndex], COMMA)

			}
		} else {

		}
	}

	return dsp.GetStorageStrings()
}

//VMNameStrip strip name of var or mixin
func VMNameStrip(src string) (string, int, int, error) {

	sepIndex := strings.IndexAny(src, VMS)

	if sepIndex < 0 {
		err := fmt.Errorf("cannot strip name - not found : or = or {. If You define variables, you should use : or =. ")
		return src, sepIndex, UNKNOWN, err
	}
	name := src[:sepIndex]
	if strings.HasPrefix(src[sepIndex:], RBBegin) {
		return name, sepIndex, Mix, nil
	}
	return name, sepIndex, Var, nil
}

//DataStorageForProcess is data of process in compiling hcss
type DataStorageForProcess struct {
	Variables     map[string]*Variable
	MixIns        map[string]*MixIn
	StringStorage []*StringData
}

//GetStorageStrings get strings which is stored in dsp
func (dsp *DataStorageForProcess) GetStorageStrings() string {

	var sb strings.Builder
	for _, sd := range dsp.StringStorage {
		sb.WriteString(sd.Content)
	}
	return sb.String()
}

//NewDataStorageForProcess create new Process
func NewDataStorageForProcess() *DataStorageForProcess {
	return &DataStorageForProcess{make(map[string]*Variable, 0), make(map[string]*MixIn, 0), make([]*StringData, 0)}
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
	Content        string
	CompiledString string
	ContentType    int
}

//NewVariable get new variable
func NewVariable(vcontent string, vtype int) *Variable {
	return &Variable{vcontent, "", vtype}
}

//MixIn represents data about MixIn
type MixIn struct {
	paramString []string
	content     string
}
