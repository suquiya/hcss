package hcss

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
