package main

import (
	"bufio"
	"os"
	"strings"
)

func getInput() (str string) {
	in := bufio.NewReader(os.Stdin)
	line, err := in.ReadString('\n')
	line = strings.TrimSuffix(line, "\n")
	if err == nil {
		return line
	}
	return ""
}
