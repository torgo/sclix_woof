package main

import (
	"bufio"
	"fmt"
	"github.com/snyk/cli-extension-lib-go"
	"os"
	"path"
)

func main() {
	extensionRoot, err := cli_extension_lib_go.GetExtensionRoot()
	if err != nil {
		fmt.Println("failed to get extension root:", err)
		os.Exit(1)
	}
	extensionMetadataPath := path.Join(extensionRoot, "extension.json")
	extMeta, err := cli_extension_lib_go.DeserExtensionMetadata(extensionMetadataPath)
	if err != nil {
		fmt.Println("failed deserializing extension metadata file:", err)
		os.Exit(1)
	}
	extensionName := extMeta.Name

	stdinReader := bufio.NewReader(os.Stdin)
	inputString, err := cli_extension_lib_go.ReadInput(stdinReader)
	if err != nil {
		fmt.Println("failed reading input")
		panic(err)
	}

	input, err := cli_extension_lib_go.ParseInput[WoofInput](inputString)
	if err != nil {
		fmt.Println("failed parsing input JSON to struct")
		panic(err)
	}

	debugLogger := cli_extension_lib_go.GetDebugLogger(input.Debug, extensionName)
	debugLogger.Println("extension name:", extensionName)
	debugLogger.Println("input.Debug:", input.Debug)
	debugLogger.Println("input.PoxyPort:", input.ProxyPort)
	debugLogger.Println("input.Args.Lang:", input.Args.Lang)

	lang := input.Args.Lang
	fmt.Println("Woof in", lang)
}

type WoofInput struct {
	Lang string `json:"lang"`
}
