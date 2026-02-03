package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/xingliuhua/gramod/model"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/urfave/cli/v2"
)

// Rainbow-like color palette for coloring nodes/edges.
var colorPalette = []string{
	"#E74C3C", // red
	"#E67E22", // orange
	"#F1C40F", // yellow
	"#2ECC71", // green
	"#1ABC9C", // cyan
	"#3498DB", // blue
	"#9B59B6", // purple
}

const NODE_LABEL_MAX_LEN = 20

// Entry point with urfave/cli/v2
func main() {
	app := &cli.App{
		Name:  "gomodviz",
		Usage: "Generate a visual graph of Go module dependencies",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "target",
				Aliases:  []string{"t"},
				Usage:    "Generate subgraph starting from this module path (optional).",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "focus",
				Aliases:  []string{"f"},
				Usage:    "View both dependents (who depends on it) and dependencies of this module.",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "Output PNG file path (must end with .png). Defaults based on mode.",
				Required: false,
				Value:    "",
			},
			&cli.BoolFlag{
				Name:    "color",
				Aliases: []string{"c"},
				Usage:   "Enable automatic coloring for edges.",
			},
		},
		Action: func(c *cli.Context) error {
			target := c.String("target")
			focus := c.String("focus")
			outFile := c.String("output")
			autoColor := c.Bool("color")

			// validate / prepare output path
			if outFile == "" {
				switch {
				case focus != "":
					outFile = "deps_focus.png"
				case target != "":
					outFile = "deps_sub.png"
				default:
					outFile = "deps_all.png"
				}
			}
			if !strings.HasSuffix(outFile, ".png") {
				return fmt.Errorf("output file must end with .png (got %s)", outFile)
			}
			if strings.HasPrefix(outFile, "~") {
				if home, _ := os.UserHomeDir(); home != "" {
					outFile = filepath.Join(home, strings.TrimPrefix(outFile, "~"))
				}
			}
			if err := os.MkdirAll(filepath.Dir(outFile), 0755); err != nil {
				return err
			}

			ctx := context.Background()

			// load dependency graph
			cmd := exec.Command("go", "mod", "graph")
			out, err := cmd.Output()
			if err != nil {
				return fmt.Errorf("failed to execute go mod graph: %w", err)
			}
			graph := parseGraph(out)

			switch {
			case focus != "":
				fmt.Printf("🔍 Focusing on [%s] (showing dependents & dependencies)...\n", focus)
				sub := collectFocusGraph(graph, focus)
				if len(sub) == 0 {
					return fmt.Errorf("module not found in dependency graph: %s", focus)
				}
				if err := renderGraph(ctx, sub, model.Module{Path: focus}, outFile, autoColor, &model.Module{Path: focus}); err != nil {
					return err
				}
				fmt.Printf("✅ Focus graph saved to: %s\n", outFile)
				return nil

			case target != "":
				fmt.Printf("🎯 Generating dependency subgraph for [%s] ...\n", target)
				sub := collectDependenciesByPath(graph, target)
				if len(sub) == 0 {
					return fmt.Errorf("module not found: %s", target)
				}
				if err := renderGraph(ctx, sub, model.Module{Path: target}, outFile, autoColor, nil); err != nil {
					return err
				}
				fmt.Printf("✅ Subgraph saved to: %s\n", outFile)
				return nil

			default:
				rootOut, _ := exec.Command("go", "list", "-m").Output()
				root := parseModule(strings.TrimSpace(string(rootOut)))
				fmt.Println("🌳 Generating full dependency graph ...")
				if err := renderGraph(ctx, graph, root, outFile, autoColor, nil); err != nil {
					return err
				}
				fmt.Printf("✅ Full graph saved to: %s\n", outFile)
				return nil
			}
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("❌", err)
		os.Exit(1)
	}
}

// Parse `go mod graph` output into adjacency list
func parseGraph(output []byte) model.DependencyMap {
	graph := make(model.DependencyMap)
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}
		parent := parseModule(parts[0])
		child := parseModule(parts[1])
		graph[parent] = append(graph[parent], child)
	}
	return graph
}

// Parse "path@version" string into model.Module struct
func parseModule(s string) model.Module {
	at := strings.LastIndex(s, "@")
	if at == -1 {
		return model.Module{s, ""}
	}
	return model.Module{s[:at], s[at+1:]}
}

// Collect dependencies (downstream) of a given module path
func collectDependenciesByPath(g model.DependencyMap, path string) model.DependencyMap {
	result := make(model.DependencyMap)
	visited := make(map[model.Module]bool)

	var starts []model.Module
	for m := range g {
		if m.Path == path {
			starts = append(starts, m)
		}
	}
	for _, children := range g {
		for _, c := range children {
			if c.Path == path {
				starts = append(starts, c)
			}
		}
	}
	if len(starts) == 0 {
		return nil
	}

	var dfs func(model.Module)
	dfs = func(m model.Module) {
		if visited[m] {
			return
		}
		visited[m] = true
		for _, child := range g[m] {
			result[m] = append(result[m], child)
			dfs(child)
		}
	}
	for _, s := range starts {
		dfs(s)
	}
	return result
}

// Collect both dependents (who depends on it) and dependencies.
func collectFocusGraph(g model.DependencyMap, path string) model.DependencyMap {
	result := make(model.DependencyMap)
	targetMods := findModulesByPath(g, path)
	if len(targetMods) == 0 {
		return nil
	}
	// who depends on it
	for parent, children := range g {
		for _, c := range children {
			for _, t := range targetMods {
				if c.Path == t.Path && c.Version == t.Version {
					result[parent] = append(result[parent], c)
				}
			}
		}
	}
	// what it depends on
	for _, t := range targetMods {
		if deps, ok := g[t]; ok {
			result[t] = append(result[t], deps...)
		}
	}
	return result
}
func findModulesByPath(g model.DependencyMap, path string) []model.Module {
	var mods []model.Module
	for m := range g {
		if m.Path == path {
			mods = append(mods, m)
		}
	}
	for _, children := range g {
		for _, c := range children {
			if c.Path == path {
				mods = append(mods, c)
			}
		}
	}
	return mods
}

// Wrap label string for better multi-line layout
func wrapLabelSmart(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	var lines []string
	start := 0
	lastSplit := -1
	for i, r := range runes {
		if r == '/' || r == '.' || r == '-' || r == '@' {
			lastSplit = i + 1
		}
		if i-start+1 >= maxLen {
			split := i + 1
			if lastSplit > start && lastSplit < i {
				split = lastSplit
			}
			lines = append(lines, strings.TrimSpace(string(runes[start:split])))
			start = split
			lastSplit = -1
		}
	}
	if start < len(runes) {
		lines = append(lines, strings.TrimSpace(string(runes[start:])))
	}
	if len(lines) > 1 && len([]rune(lines[len(lines)-1])) < 5 {
		lines[len(lines)-2] += lines[len(lines)-1]
		lines = lines[:len(lines)-1]
	}
	return strings.Join(lines, "\n")
}

// Render the graph into PNG using go-graphviz.
func renderGraph(ctx context.Context, deps model.DependencyMap, root model.Module, file string, autoColor bool, focus *model.Module) error {
	gv, err := graphviz.New(ctx)
	if err != nil {
		return err
	}
	graph, err := gv.Graph()
	if err != nil {
		return err
	}
	defer func() {
		graph.Close()
		gv.Close()
	}()

	if autoColor {
		rand.Seed(42)
	}

	nodes := make(map[string]*graphviz.Node)

	getNode := func(m model.Module) *graphviz.Node {
		key := m.Path + "@" + m.Version
		if n, ok := nodes[key]; ok {
			return n
		}
		n, _ := graph.CreateNodeByName(key)
		n.SetLabel(wrapLabelSmart(key, NODE_LABEL_MAX_LEN))
		n.SetShape("rect")
		n.SetStyle("solid")
		n.SetColor("black")

		if focus != nil && strings.HasPrefix(key, focus.Path+"@") {
			n.SetStyle("filled,bold,rounded")
			n.SetFillColor("lightskyblue")
		}
		nodes[key] = n
		return n
	}

	modules := make([]model.Module, 0, len(deps))
	for m := range deps {
		modules = append(modules, m)
	}
	sort.Slice(modules, func(i, j int) bool {
		if modules[i].Path == modules[j].Path {
			return modules[i].Version < modules[j].Version
		}
		return modules[i].Path < modules[j].Path
	})

	for _, parent := range modules {
		pNode := getNode(parent)
		children := deps[parent]
		sort.Slice(children, func(i, j int) bool {
			if children[i].Path == children[j].Path {
				return children[i].Version < children[j].Version
			}
			return children[i].Path < children[j].Path
		})
		color := "gray"
		if autoColor {
			color = colorPalette[rand.Intn(len(colorPalette))]
		}
		for _, child := range children {
			cNode := getNode(child)
			e, _ := graph.CreateEdgeByName("", pNode, cNode)
			e.SetColor(color)
		}
	}

	graph.SetRankDir("TB")
	graph.SetSplines("true")

	buf := new(bytes.Buffer)
	if err := gv.Render(ctx, graph, "png", buf); err != nil {
		return err
	}
	return os.WriteFile(file, buf.Bytes(), 0o644)
}
