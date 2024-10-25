package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

type Port struct {
	Name  string
	InOut string
	Type  string
	MSB   int // most
	LSB   int // least
}

type Ports []*Port

type VHDL struct {
	Entity    string
	Ports     Ports
	ClkPort   Port
	ResetPort Port
	lines     []string
}

func LoadVHDL(fp *os.File) (*VHDL, error) {
	scanner := bufio.NewScanner(fp)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &VHDL{lines: lines}, nil
}

func (vhdl *VHDL) Parse() error {
	entityStarted := false
	entityName := ""
	for _, line := range vhdl.lines {
		formattedLine := FormatLine(line)
		name, _ := ParseEntityStart(formattedLine)
		if name != "" {
			entityStarted = true
			entityName = name
		}
		if entityStarted {
			if ok, _ := ParseEntityEnd(formattedLine, entityName); ok {
				entityStarted = false
			}

			port, _ := ParsePort(formattedLine)
			if port != nil {
				vhdl.Ports = append(vhdl.Ports, port)
			}
		}
	}

	vhdl.Entity = entityName
	for _, port := range vhdl.Ports {
		if strings.Contains(port.Name, "CLK") && port.Type == "std_logic" {
			vhdl.ClkPort = *port
		} else if (port.Name == "RST" || port.Name == "RESET") && port.Type == "std_logic" {
			vhdl.ResetPort = *port
		}
	}
	return nil
}

func ParsePort(line string) (*Port, error) {
	re, err := regexp.Compile(`(?P<name>\w+)\s*:\s*(?P<inout>in|out|inout)\s*(?P<type>\w+)\s*(?P<range>\(\s*(?P<msb>\d+)\s*(?:downto|to)\s*(?P<lsb>\d+)\s*\))?`)
	if err != nil {
		return nil, err
	}
	matches := re.FindStringSubmatch(line)
	if len(matches) == 0 {
		return nil, fmt.Errorf("no matches")
	}
	port := &Port{
		Name:  matches[1],
		InOut: matches[2],
		Type:  matches[3],
	}
	if len(matches) == 7 {
		port.MSB, _ = strconv.Atoi(matches[5]) // msb
		port.LSB, _ = strconv.Atoi(matches[6]) // lsb
	}
	return port, nil
}

func ParseEntityStart(line string) (string, error) {
	re, err := regexp.Compile(`entity\s+(?P<name>\w+)\s+is`)
	if err != nil {
		return "", err
	}
	matches := re.FindStringSubmatch(line)
	if len(matches) == 0 {
		return "", fmt.Errorf("no matches")
	}
	return matches[1], nil
}

func ParseEntityEnd(line string, name string) (bool, error) {
	re, err := regexp.Compile(`end\s+(?P<name>\w+)\s*;`)
	if err != nil {
		return false, err
	}
	matches := re.FindStringSubmatch(line)
	if len(matches) == 0 {
		return false, fmt.Errorf("no matches")
	}
	return matches[1] == name, nil
}

func FormatLine(line string) string {
	return strings.Join(strings.Fields(line), " ")
}

func main() {
	var (
		inputPath  = flag.String("i", "", "input file path")
		outputPath = flag.String("o", "", "output file path")
	)
	flag.Parse()

	inputFp, err := os.Open(*inputPath)
	if err != nil {
		panic(err)
	}
	defer inputFp.Close()

	vhdl, err := LoadVHDL(inputFp)
	if err != nil {
		panic(err)
	}

	vhdl.Parse()

	tpl, err := template.New("tb").Funcs(template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
	}).ParseFiles("vhd.tpl")
	if err != nil {
		panic(err)
	}

	if *outputPath == "" {
		dir := filepath.Dir(*inputPath)
		*outputPath = filepath.Join(dir, fmt.Sprintf("tb_%s.vhd", vhdl.Entity))
	}

	outputFp, err := os.Create(*outputPath)
	if err != nil {
		panic(err)
	}
	defer outputFp.Close()

	err = tpl.ExecuteTemplate(outputFp, "tb", vhdl)
	if err != nil {
		panic(err)
	}
}
