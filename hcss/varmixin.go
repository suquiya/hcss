package hcss

import (
	"fmt"
	"strings"
)

//Variable represents data about a variable
type Variable struct {
	Name           string
	Content        []ContentTyper
	CompiledString string
}

//VarCall is another name of string, for record calling variable
type VarCall string

//ContentType return VarCall v's ContentType
func (v VarCall) ContentType() int {
	return CallVar
}

//ContentType return *variable v's ContentType and impliment interface BlockData
func (v *Variable) ContentType() int {
	return DefVar
}

//NewVariable get new variable
func NewVariable(defname string, vcontent []ContentTyper) *Variable {
	return &Variable{defname, vcontent, ""}
}

//MixIn represents data about MixIn
type MixIn struct {
	Name              string
	ParamNames        []string
	SortedParamString StrSorter
	Content           string
}

//ContentType return *MixIn mi's ContentType
func (mi *MixIn) ContentType() int {
	return DefMix
}

//MinInCallArgs is struct for storage args in Call MixIn
type MinInCallArgs struct {
	CalledMixinName string
	Content         string
	Args            map[string]string
}

//NewMixInCallArgs return new MixInCallArgs instance
func NewMixInCallArgs(MixInName string) *MinInCallArgs {
	return &MinInCallArgs{MixInName, "", make(map[string]string)}
}

//ArgsParse parse mixin's arg
func (mica *MinInCallArgs) ArgsParse(mi *MixIn, args string, content string) (*MinInCallArgs, error) {
	mica.Content = content
	argArray := strings.Split(args, COMMA)

	if len(argArray) > len(mi.ParamNames) {
		err := fmt.Errorf("Too many args")
		for i, pn := range mi.ParamNames {
			mica.Args[pn] = argArray[i]
		}
		return mica, err
	}
	if len(argArray) < len(mi.ParamNames) {
		err := fmt.Errorf("insufficient arg number")
		return mica, err
	}

	for i, pn := range mi.ParamNames {
		mica.Args[pn] = argArray[i]
	}

	return mica, nil
}

const (
	cb0 = CBBegin
	cb1 = "\r\n{"
)

//HasPrefixOfMixInContentBegin return src has or does not have mixInContentBeginStrings
func HasPrefixOfMixInContentBegin(src string) bool {
	prefix := strings.HasPrefix(src, cb0) || strings.HasPrefix(src, cb1)
	return prefix && !strings.HasPrefix(src, HugoTmpBegin)
}

//ContentType return *MinInCallArgs ContentType
func (mica *MinInCallArgs) ContentType() int {
	return CallMix
}
