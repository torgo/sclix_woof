package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const EXTENSION_NAME = "sclix-woof"

func main() {
	inputString, err := ReadInput()
	if err != nil {
		fmt.Println("failed reading input")
		panic(err)
	}

	input, err := ParseInput[WoofInput](inputString)
	if err != nil {
		fmt.Println("failed parsing input JSON to struct")
		panic(err)
	}

	debugLogger := GetDebugLogger(input.Debug)

	debugLogger.Println("input.Debug:", input.Debug)
	debugLogger.Println("input.PoxyPort:", input.ProxyPort)
	debugLogger.Println("input.Args.Lang:", input.Args.Lang)

	// Question: if lang is not set (in InputData.Args), then we want to default to "en".
	// Should the extension do that here in code?
	// Or should it be part of the extension lib (`ParseInput`) in concert with reading the extension.json file?
	if input.Args.Lang == "" {
		input.Args.Lang = "en"
	}

	lang := input.Args.Lang
	fmt.Println("Woof in", lang)
}

type WoofInput struct {
	Lang string `json:"lang"`
}

///////////////////////////////////////////////////////////////////////////////
// The stuff below should go in the extensions library
///////////////////////////////////////////////////////////////////////////////

type InputData[T any] struct {
	// Standard stuff do we want to passed to all extensions
	Debug     bool `json:"debug"`
	ProxyPort int  `json:"proxyPort"`

	// Extension-specific args
	Args T `json:"args"`
}

func GetInputArgs() *WoofInput {
	var args WoofInput
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

// Question: should this fail (return an error) if certain mandatory fields are missing?
// For example, the proxy port?
func ParseInput[T any](inputString string) (*InputData[T], error) {
	rawBytes := []byte(inputString)

	var input InputData[T]
	err := json.Unmarshal(rawBytes, &input)
	if err != nil {
		return nil, err
	}

	return &input, nil
}

// Question: This uses the EXTENSION_NAME constant. Nothing currently forces that constant to be in sync with extension.json.
// How should we enforce this? We could load it from the extension.json file? Or we could do it with some kind of lint check
// when we build it.
func GetDebugLogger(debug bool) *log.Logger {
	logPrefix := fmt.Sprintf("[%s] ", EXTENSION_NAME)
	debugLogger := log.New(os.Stderr, logPrefix, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	if !debug {
		debugLogger.SetOutput(ioutil.Discard)
	}
	return debugLogger
}
