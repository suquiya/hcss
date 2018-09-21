package main

func main() {

	string := "$maincolor = {{ .Site.Params.MainColor }}\r\ndiv{color: $maincolor}"

	css := hcss.Parse()
}
