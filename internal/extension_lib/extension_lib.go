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

type InputData[T any] struct {
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

func GetDebugLogger(debug bool, extensionName string) *log.Logger {
	logPrefix := fmt.Sprintf("[%s] ", extensionName)
	debugLogger := log.New(os.Stderr, logPrefix, log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	if !debug {
		debugLogger.SetOutput(ioutil.Discard)
	}
	return debugLogger
}

func DeserExtensionMetadata() (*extension_metadata.ExtensionMetadata, error) {
	exPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	dirPath := filepath.Dir(exPath)
	extensionMetadataFile := path.Join(dirPath, "extension.json")

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
