package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
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

func LoadVHDLFile(path string) ([]string, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

func ParsePort(line string) (*Port, error) {
	re, err := regexp.Compile(`(?P<name>\w+)\s*:\s*(?P<inout>in|out)\s*(?P<type>\w+)\s*(?P<range>\(\s*(?P<msb>\d+)\s*(?:downto|to)\s*(?P<lsb>\d+)\s*\))?`)
	if err != nil {
		return nil, err
	}
	matches := re.FindStringSubmatch(line)
	if len(matches) == 0 {
		return nil, fmt.Errorf("No matches")
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
		return "", fmt.Errorf("No matches")
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
		return false, fmt.Errorf("No matches")
	}
	return matches[1] == name, nil
}

func FormatLine(line string) string {
	return strings.Join(strings.Fields(line), " ")
}

func (ports Ports) ConvPorts() string {
	format := ""
	for _, port := range ports {
		if port.Type == "std_logic" {
			format += fmt.Sprintf("\t%s : %s %s;\n", port.Name, port.InOut, port.Type)
		} else {
			format += fmt.Sprintf("\t%s : %s %s(%d downto %d);\n", port.Name, port.InOut, port.Type, port.MSB, port.LSB)
		}
	}
	return format[:len(format)-2]
}

func (ports Ports) ConvInputs() string {
	format := ""
	for _, port := range ports {
		if port.InOut == "in" {
			if port.Type == "std_logic" {
				format += fmt.Sprintf("\tsignal %s : %s := '0';\n", port.Name, port.Type)
			} else {
				format += fmt.Sprintf("\tsignal %s : %s(%d downto %d) := (others => '0');\n", port.Name, port.Type, port.MSB, port.LSB)
			}
		}
	}
	return format[:len(format)-1]
}

func (ports Ports) ConvOutputs() string {
	format := ""
	for _, port := range ports {
		if port.InOut == "out" {
			if port.Type == "std_logic" {
				format += fmt.Sprintf("\tsignal %s : %s;\n", port.Name, port.Type)
			} else {
				format += fmt.Sprintf("\tsignal %s : %s(%d downto %d);\n", port.Name, port.Type, port.MSB, port.LSB)
			}
		}
	}
	return format[:len(format)-1]
}

func (ports Ports) ConvPortMap() string {
	format := ""
	for _, port := range ports {
		format += fmt.Sprintf("\t%s => %s,\n", port.Name, port.Name)
	}
	return format[:len(format)-2]
}

func main() {
	var (
		inputPath  = flag.String("i", "", "input file path")
		outputPath = flag.String("o", "", "output file path")
	)
	flag.Parse()

	lines, err := LoadVHDLFile(*inputPath)
	if err != nil {
		panic(err)
	}
	ports := make(Ports, 0)
	entityStarted := false
	entityName := ""
	for _, line := range lines {
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
				ports = append(ports, port)
			}
		}
	}
	tpl, err := template.ParseFiles("templates/test.vhd")
	if err != nil {
		panic(err)
	}
	m := map[string]interface{}{
		"entityName":    fmt.Sprintf("Test_%s", entityName),
		"componentName": entityName,
		"ports":         ports.ConvPorts(),
		"inputs":        ports.ConvInputs(),
		"outputs":       ports.ConvOutputs(),
		"portMap":       ports.ConvPortMap(),
	}
	fp, err := os.Create(*outputPath)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	err = tpl.Execute(fp, m)
	if err != nil {
		panic(err)
	}
}
