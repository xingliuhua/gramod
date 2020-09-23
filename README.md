[中文版](https://github.com/xingliuhua/gramod/blob/master/README.cn.md)
# gramod


This is a graphical tool for go mod
## Background

The readability of the go mod graph command line is too poor
## Feature
* Generate all dependent graphics
* Generates the graphics for the specified dependencies
## Install
go get github.com/xingliuhua/gramod
## Usage
```text
gramod
// Generate all dependent graphics
```
![](https://github.com/xingliuhua/gramod/blob/master/gramod_eg1.png)

```text
gramod -s github.com/xingliuhua/gramod@v1.0.0
// Generates the graphics for the specified dependencies
```
![](https://github.com/xingliuhua/gramod/blob/master/gramod_eg2.png)

## Maintainers

[@xingliuhua](https://github.com/xingliuhua).

## Contributing

Feel free to dive in! [Open an issue] or submit PRs.
