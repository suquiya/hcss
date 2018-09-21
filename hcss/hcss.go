package hcss

import (
	"strings"
)

//Compile parse hcss string
const (
	cssExt = ".css"
)

//Compile compile hcss string to css string
func Compile(hcss string) string {

	p := NewProcess()

	return p.buf.String()
}

//Process is data of process in compiling hcss
type Process struct {
	Functions map[string]*VariableData
	Variables map[string]*FunctionData
	buf       *strings.Builder
}

//NewProcess create new Process
func NewProcess() *Process {
	var sb strings.Builder
	return &Process{make(map[string]*VariableData, 0), make(map[string]*FunctionData, 0), &sb}
}

//VariableData represents data about a variable
type VariableData struct {
	name    string
	content string
}

//FunctionData represents data about function
type FunctionData struct {
	paramString []string
	content     string
}
