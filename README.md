# hcss(WIP)
hcss is pure Go css pre-processor. This is created for using process of site-creation with Golang HTML Template.
Generated CSS can be used as Golang Template(Because of Part of Go HTML Template remain).

Inspired by scss and gcss syntax, but this may be simpler than sass(scss).

suquiya is a beginner of Go programming, so pull requests and issues are appropriated.
To be honest, please help suquiya...

## Syntax
### Variables(Now implimating....)

Before Compile

```
$main-color= #d3381c
$text-color: {{ .Site.Params.TextColor }}

h1{
    color:$main-color;
}
div{
    color: $text-color
}
```
After Compile
```
h1{
    color: #d3381c;
}

div{
    color: {{ .Site.Params.TextColor }};
}
```

Either of "=" or ":" can be used in Var definition.
newline and ";" is interpreted as end of Var definition.

## TO DO

+ Add features
  + mixin
  + nest
  + Comment
  + Extend(I can not decide whether implement it or do not now.)
  + Import

## Links
+ [Hugo][hugo]
+ [GCSS](https://github.com/yosssi/gcss)

[hugo]: https://github.com/gohugoio/hugo