package extension_lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/snyk/cli/cliv2/exported/extension_metadata"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

type ExtensionInput[T any] struct {
	// Standard stuff do we want to passed to all extensions
	Debug     bool `json:"debug"`
	ProxyPort int  `json:"proxyPort"`

	// Extension-specific args
	Args T `json:"args"`
}

func ReadInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	consecutiveNewlinesCount := 0
	chars := []rune{}

	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			// this may just be EOF so don't propagate the error
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

func ParseInput[T any](inputString string) (*ExtensionInput[T], error) {
	rawBytes := []byte(inputString)

	var input ExtensionInput[T]
	err := json.Unmarshal(rawBytes, &input)
	if err != nil {
		return nil, err
	}

	return &input, nil
}

func GetDebugLogger(debug bool, extensionName string) *log.Logger {
	logPrefix := fmt.Sprintf("[%s] ", extensionName)
	debugLogger := log.New(os.Stderr, logPrefix, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	if !debug {
		debugLogger.SetOutput(ioutil.Discard)
	}
	return debugLogger
}

func GetExtensionRoot() (string, error) {
	exPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	dirPath := filepath.Dir(exPath)
	return dirPath, nil
}

func DeserExtensionMetadata(extensionRoot string) (*extension_metadata.ExtensionMetadata, error) {
	extensionMetadataFile := path.Join(extensionRoot, "extension.json")

	bytes, err := os.ReadFile(extensionMetadataFile)
	if err != nil {
		return nil, err
	}

	var extMeta extension_metadata.ExtensionMetadata
	err = json.Unmarshal(bytes, &extMeta)
	if err != nil {
		return nil, err
	}

	return &extMeta, nil
}

func LoadExtensionMetadata() (*extension_metadata.ExtensionMetadata, error) {
	extensionRoot, err := GetExtensionRoot()
	if err != nil {
		return nil, err
	}
	return DeserExtensionMetadata(extensionRoot)
}
