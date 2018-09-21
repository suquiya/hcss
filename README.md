# hcss(WIP)
hcss is pure Go css pre-processor. This is created for using process of site-creation with Hugo.
Generated CSS can be used in ExecuteAsTemplate.

Inspired by scss and gcss syntax, but this is simler than it.


## Syntax
### Variables

Before Compile

```
$main-color: #d3381c
$text-color: {{ .Site.Params.TextColor }}

h1{
    color:$main-color
}
div{
    color: $text-color
}
```
After Compile
```
h1{
    color: #d3381c
}

div{
    color: {{ .Site.Params.TextColor }}
}
```

## Links
+ [Hugo][hugo]
+ [GCSS](https://github.com/yosssi/gcss)

[hugo]: https://github.com/gohugoio/hugo