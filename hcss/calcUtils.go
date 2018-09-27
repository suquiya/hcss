package hcss

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//StrSorter is sort Tool with string length
type StrSorter []string

func (ss StrSorter) Len() int           { return len(ss) }
func (ss StrSorter) Less(i, j int) bool { return len(ss[i]) > len(ss[j]) }
func (ss StrSorter) Swap(i, j int)      { ss[i], ss[j] = ss[j], ss[i] }
