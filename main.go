package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Please provide a file path\n")
		return
	}
	pathToFile := os.Args[1]

	isValidTsFilePath := strings.HasSuffix(pathToFile, ".ts")

	if !isValidTsFilePath {
		fmt.Printf("File %s is not a valid TypeScript file\n", pathToFile)
		return
	}
	absPath, err := filepath.Abs(pathToFile)
	if err != nil {
		fmt.Printf("Error resolving absolute path for %s: %v\n", absPath, err)
		return
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", absPath, err)
	}

	result := api.Transform(string(content), api.TransformOptions{
		Format: api.FormatCommonJS,
		Loader: api.LoaderTS,
	})
	cmd := exec.Command("node", "-e", string(result.Code))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("%d errors and %d warnings\n",
		len(result.Errors), len(result.Warnings))

	runErr := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to run Node.js code: %v\n", runErr)
	}
}
