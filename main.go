package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	inputString, err := ReadInput()
	if err != nil {
		fmt.Println("failed reading input")
		panic(err)
	}

	input, err := ParseInput[Input](inputString)
	if err != nil {
		fmt.Println("failed parsing input JSON to struct")
		panic(err)
	}

	lang := input.Args.Lang
	fmt.Println("Woof in", lang)
}

type Input struct {
	Lang string `json:"lang"`
}

///////////////////////////////////////////////////////////////////////////////
// The stuff below should go in the extensions library
///////////////////////////////////////////////////////////////////////////////

type InputData[T any] struct {
	// TODO: what standard stuff needs to go here?
	ProxyPort int `json:"proxyPort"`

	// Extension-specific args
	Args T `json:"args"`
}

func GetInputArgs() *Input {
	var args Input
	args.Lang = "en"
	return &args
}

func ReadInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	consecutiveNewlinesCount := 0
	chars := []rune{}

	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			// assuming this is EOF, but should probably verify at some point!
			// TODO: check if this is an EOF
			// if so, break
			// if not return "", err
			break
		}
		if char == '\n' {
			consecutiveNewlinesCount++
		} else {
			consecutiveNewlinesCount = 0
			chars = append(chars, char)
		}

		if consecutiveNewlinesCount == 2 {
			break
		}
	}

	inputString := string(chars)
	return inputString, nil
}

func ParseInput[T any](inputString string) (*InputData[T], error) {
	rawBytes := []byte(inputString)

	var input InputData[T]
	err := json.Unmarshal(rawBytes, &input)
	if err != nil {
		return nil, err
	}

	return &input, nil
}
