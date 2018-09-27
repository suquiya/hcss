package main

import (
	"fmt"

	"github.com/suquiya/hcss/hcss"
)

func main() {

	str := "$maincolor = {{ .Site.Params.MainColor }};\r\n$text-color: #2b2b2b\r\ndiv{\r\n\tcolor: $maincolor\r\nbackground-color: $text-color\r\n}\r\n"

	css := hcss.Parse(str)

	fmt.Println(css)
}
