package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	lines := "aaa\nbbb\nccc\n\n"
	br := bufio.NewReader(strings.NewReader(lines))
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			return
		}
		fmt.Println(line)
	}

}
