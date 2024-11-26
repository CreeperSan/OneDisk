package input

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func ReadInt() int {
	return ReadIntFallback(0)
}

func ReadIntFallback(fallback int) int {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return fallback
	}
	input = strings.TrimSpace(input)
	number, err := strconv.Atoi(input)
	if err != nil {
		return fallback
	}
	return number
}

func ReadString() string {
	return ReadStringWithFallback("")
}

func ReadStringWithFallback(fallback string) string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return fallback
	}
	return strings.TrimSpace(input)
}
