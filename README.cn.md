[English](https://github.com/xingliuhua/gramod/blob/master/README.md)
# gramod

这是一个go mod的图形化工具

## 背景
go mod graph 生成的依赖报告可读性太差，图形化更方便。
市面上有类似的开源库，但是一旦依赖比较多，生成的图片密密麻麻，可读性极差，而且不能只查看具体某一子依赖的依赖。
## 功能特点
* 支持生成项目所有依赖的图形
* 支持生成指定子依赖的分析图形
* 线条区分度更大
* 版本名称适当折行，可读性更佳

## 安装
go get github.com/xingliuhua/gramod
## 使用
命令行中使用
gramod
// 生成项目所有依赖图
![](https://github.com/xingliuhua/gramod/blob/master/gramod_eg1.png)

gramod -s github.com/xingliuhua/gramod@v1.0.0
// 只生成github.com/xingliuhua/gramod@v1.0.0的依赖
![](https://github.com/xingliuhua/gramod/blob/master/gramod_eg2.png)
## 维护

[@xingliuhua](https://github.com/xingliuhua).

## 贡献

Feel free to dive in! [Open an issue] or submit PRs.
