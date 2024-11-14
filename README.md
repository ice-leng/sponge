<p align="center">
<img width="500px" src="https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/logo.png">
</p>

<div align=center>

[![Go Report](https://goreportcard.com/badge/github.com/zhufuyi/sponge)](https://goreportcard.com/report/github.com/zhufuyi/sponge)
[![codecov](https://codecov.io/gh/zhufuyi/sponge/branch/main/graph/badge.svg)](https://codecov.io/gh/zhufuyi/sponge)
[![Go Reference](https://pkg.go.dev/badge/github.com/zhufuyi/sponge.svg)](https://pkg.go.dev/github.com/zhufuyi/sponge)
[![Go](https://github.com/zhufuyi/sponge/workflows/Go/badge.svg?branch=main)](https://github.com/zhufuyi/sponge/actions)
[![Awesome Go](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go)
[![License: MIT](https://img.shields.io/github/license/zhufuyi/sponge)](https://img.shields.io/github/license/zhufuyi/sponge)

</div>

## 当前版本为魔改 restful 版本，具体用法请看 [官方文档](https://github.com/zhufuyi/sponge)

// 注意 暂时实现 http web 版本，  pb 版本没测试
>>>>>>> main

### 源码安装 使用说明
```git
    git clone https://github.com/ice-leng/sponge.git
    cd sponge/cmd/sponge
    go run ./main.go upgrade
```

### 主要魔改功能有
- 基于数据库dsn 添加表前缀
```html
    数据库dsn: root:@(127.0.0.1:3306)/hyperf;prefix=t_
```

- 由于本人是mac 下载替换文件，那是真的全替换... 
- 下载代码功能兼容，支持成mac，命令行 在那个目录，代码就在这个目录下生成

```shell
    mkdir xxx
    cd xxx
    sponge run 
    ... // web 操作 代码下载 
	ls -al
```