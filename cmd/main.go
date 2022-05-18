package main

import (
	"fmt"
	"github.com/snyk/sclix-woof/internal/extension_lib"
)

const EXTENSION_NAME = "sclix-woof"

func main() {
	inputString, err := extension_lib.ReadInput()
	if err != nil {
		fmt.Println("failed reading input")
		panic(err)
	}

	input, err := extension_lib.ParseInput[WoofInput](inputString)
	if err != nil {
		fmt.Println("failed parsing input JSON to struct")
		panic(err)
	}

	debugLogger := extension_lib.GetDebugLogger(input.Debug, EXTENSION_NAME)

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
