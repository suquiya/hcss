package hcss

//Variable represents data about a variable
type Variable struct {
	Name           string
	Content        []Block
	CompiledString string
	ContentType    int
}

//WhichContentType return *variable v's ContentType and impliment interface BlockData
func (v *Variable) WhichContentType() int {
	return v.ContentType
}

//NewVariable get new variable
func NewVariable(defname string, vcontent string, vtype int) *Variable {
	return &Variable{defname, vcontent, "", vtype}
}

//MixIn represents data about MixIn
type MixIn struct {
	Name              string
	ParamNames        []string
	SortedParamString StrSorter
	Content           string
	ContentType       int
}

//WhichContentType return *MixIn mi's ContentType
func (mi *MixIn) WhichContentType() int {
	return mi.ContentType
}

//MinInCallArgs is struct for storage args in Call MixIn
type MinInCallArgs struct {
	Content string
	Args    map[string]string
}
