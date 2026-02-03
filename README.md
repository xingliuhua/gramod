[中文版](https://github.com/xingliuhua/gramod/blob/master/README.cn.md)

# gramod

A lightweight **Go module dependency visualizer**.

`gramod` parses `go mod graph` output and generates beautiful, accurate PNG graphs  
to help you understand how your Go modules depend on each other.



## ✨ Features

- **Full graph** – visualize all module dependencies in your project.
- **Target mode** – display only one module’s dependency tree (`--target`).
- **Focus mode** – show who depends on a specific module and what it depends on (`--focus`).
- **Color mode** – optionally colorize nodes and edges for clarity.
- **Deterministic output** – stable order, the same graph each time.
- **Pure Go renderer** – no external Graphviz installation required.



## image
all:
![deps_all.png](image/deps_all.png)
sub:
![deps_sub.png](image/deps_sub.png)
focus:
![deps_focus.png](image/deps_focus.png)



## 🚀 Installation

### 🧰 Option 1 – Go developers

```bash
go install github.com/xingliuhua/gramod@latest
```

> Requires Go 1.18 or newer  
The binary will be placed in `$GOPATH/bin` or `$HOME/go/bin`


### 💡 Option 2 – Pre‑built binaries

Download the binary for your system from [Releases](https://github.com/xingliuhua/gramod/releases/latest):

Unpack and place **`gramod`** in your `PATH`, then run:

```bash
gramod --help
```

## 🧰 Usage

```bash
gramod [flags]
```

| Flag     | Alias | Description                                                                 |
|----------|-------|-----------------------------------------------------------------------------|
| --target | -t    | Show the dependency tree starting from a specific module                    |
| --focus  | -f    | Show both direct dependents (who depends on it) and its direct dependencies |
| --output | -o    | Output file path (must end with .png)                                       |
| --color  | -c    | Enable automatic coloring of edges                                          |
| --help   | -h    | Show help and usage information                                             |

## 📘 Examples

1️⃣ Full dependency graph

```Bash
gramod
```

Generates deps_all.png.

2️⃣ A specific module’s dependency tree

```Bash
gramod -t github.com/gin-gonic/gin -o gin_tree.png
```

3️⃣ Focus on a module’s “neighbors”

```Bash
gramod -f github.com/bwmarrin/snowflake -o snowflake_focus.png
```

Shows who depends on snowflake and what snowflake depends on.

4️⃣ Enable colorized output

```Bash
gramod --color -t github.com/gin-gonic/gin
```

## Maintainers

[@xingliuhua](https://github.com/xingliuhua).

## Contributing

Feel free to dive in! [Open an issue] or submit PRs.
