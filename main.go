package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Port struct {
	Name  string
	InOut string
	Type  string
	MSB   int // most
	LSB   int // least
}

func loadVHDLFile(path string) ([]string, error) {
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

func main() {
	lines, err := loadVHDLFile("test/IM_Gabor.vhd")
	if err != nil {
		panic(err)
	}
	for _, line := range lines {
		formattedLine := strings.Join(strings.Fields(line), " ")
		port, _ := ParsePort(formattedLine)
		if port != nil {
			fmt.Println(port)
		}
	}
}
