package main

import (
	"fmt"

	"github.com/suquiya/hcss/hcss"
)

func main() {

	string := "$maincolor = {{ .Site.Params.MainColor }};\r\n$text-color: #2b2b2b\r\ndiv{\r\n\tcolor: $maincolor\r\nbackground-color: $text-color\r\n}\r\n"

	css := hcss.Compile(string)

	fmt.Println(css)
}
