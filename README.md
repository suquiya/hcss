# hcss(WIP)
hcss is pure Go css pre-processor. This is created for using process of site-creation with Hugo.
Generated CSS can be used in ExecuteAsTemplate(Part of Hugo Template remain).

Inspired by scss and gcss syntax, but this may be simpler than sass(scss).

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

+ mixin
+ nest
+ Comment

## Links
+ [Hugo][hugo]
+ [GCSS](https://github.com/yosssi/gcss)

[hugo]: https://github.com/gohugoio/hugo