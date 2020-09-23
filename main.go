package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/xingliuhua/gramod/model"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const NAME_MAX_LINE_LEN = 20

const help =
`go mod graph tool
usage:
	gramod [-s <dependency-name-and-version>]
eg: gramod -s github.com/xingliuhua/gramod@v1.0.0
`

var s = flag.String("s", "", "special dependency name and version,eg: gramod -s github.com/xingliuhua/gramod@v1.0.0")

func main() {
	flag.Usage = func() {
		fmt.Print(help)
	}

	flag.Parse()
	goModGraphCom := exec.Command("bash", "-c", "go mod graph")
	goModGraphComOutput, _ := goModGraphCom.CombinedOutput()
	reader := bufio.NewReader(bytes.NewBufferString(string(goModGraphComOutput)))

	AllDependencyIdMap := make(map[string]string)
	AllDependencyLines := make([]model.DependencyLine, 0)
	dependencyCount := 0
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		splits := strings.Split(string(line), " ")
		if _, b := AllDependencyIdMap[splits[0]]; !b {
			AllDependencyIdMap[splits[0]] = fmt.Sprintf("id%d", dependencyCount)
			dependencyCount++
		}
		if _, b := AllDependencyIdMap[splits[1]]; !b {
			AllDependencyIdMap[splits[1]] = fmt.Sprintf("id%d", dependencyCount)
			dependencyCount++
		}
		AllDependencyLines = append(AllDependencyLines, model.DependencyLine{
			Name:           splits[0],
			DependencyName: splits[1],
		})
	}

	if s == nil || *s == "" {
		writeDependencyGraph(AllDependencyIdMap, AllDependencyLines)
	} else {
		generateSpecialDependency(*s, AllDependencyLines)
	}

}

func generateSpecialDependency(speKey string, dependencySlice []model.DependencyLine) {
	speDependencys := getSpecialDependencys(speKey, dependencySlice)
	keyMap := make(map[string]string)
	i := 0
	for _, v := range speDependencys {
		if _, b := keyMap[v.Name]; !b {
			keyMap[v.Name] = fmt.Sprintf("id%d", i)
			i++
		}
		if _, b := keyMap[v.DependencyName]; !b {
			keyMap[v.DependencyName] = fmt.Sprintf("id%d", i)
			i++
		}
	}
	writeDependencyGraph(keyMap, speDependencys)
}

func getSpecialDependencys(speKey string, dependencySlice []model.DependencyLine) []model.DependencyLine {
	dest := make([]model.DependencyLine, 0)
	que := make([]string, 0)
	que = append(que, speKey)

	for ; len(que) > 0; {
		for _, key := range que {
			for _, depentLine := range dependencySlice {
				if depentLine.Name == key {
					que = append(que, depentLine.DependencyName)
					dest = append(dest, depentLine)
				}
			}
			que = que[1:]
		}
	}

	for i := 0; i < len(dest); i++ {
		for j := i + 1; j < len(dest); j++ {
			if dest[i].Name == dest[j].Name && dest[i].DependencyName == dest[j].DependencyName {
				dest = append(dest[:i], dest[i+1:]...)
			}
		}
	}
	return dest
}

// generate png file
func writeDependencyGraph(keyMap map[string]string, dependencySlice []model.DependencyLine) {
	bufferString := bytes.NewBufferString("digraph {\nnode [shape=box];\n")
	keys := make([]string, 0)
	for k, _ := range keyMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := keyMap[k]
		k = collapseKey(k)
		bufferString.WriteString(fmt.Sprintf("%s [label = \"%s\" color = gainsboro];\n", v, k))
	}
	for _, dependency := range dependencySlice {
		nameId := keyMap[dependency.Name]
		dependencyNameId := keyMap[dependency.DependencyName]
		bufferString.WriteString(fmt.Sprintf("%s -> %s[color=%s];\n", nameId, dependencyNameId, getLineColor(nameId)))
	}
	bufferString.WriteString("}")

	command := exec.Command("bash", "-c", "echo '"+bufferString.String()+"' | dot -T png -o gramod.png")
	_, err := command.CombinedOutput()
	if err != nil {
		fmt.Println("failed:", err)
	}
	fmt.Println("success! generate a gramod.png file")
}
func collapseKey(key string) string {
	if len(key) < 1 {
		return ""
	}
	bufferString := bytes.NewBufferString("")
	spe := ""
	for i := 0; ; {
		if i+NAME_MAX_LINE_LEN > len(key) {
			bufferString.WriteString(spe + key[i:])
			break
		}
		bufferString.WriteString(spe + key[i:i+NAME_MAX_LINE_LEN])
		spe = "\n"
		i += NAME_MAX_LINE_LEN
	}
	return bufferString.String()
}
func getLineColor(nameId string) string {
	id, _ := strconv.Atoi(nameId[2:])
	switch id % 4 {
	case 0:
		return "aquamarine3"
	case 1:
		return "bisque2"
	case 2:
		return "chocolate4"
	case 3:
		return "firebrick3"
	}
	return "black"
}
