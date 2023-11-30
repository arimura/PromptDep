package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the file path: ")
	scanner.Scan()
	filePath := scanner.Text()

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Dummy implementation: Replace this with actual logic to parse Java file and find dependencies
	dependencies := findDependencies(string(content), filePath)

	// Read and combine dependencies
	var combinedContent string
	for _, dep := range dependencies {
		// println(dep)
		depContent, err := ioutil.ReadFile(dep)
		if err != nil {
			fmt.Println("Error reading dependency file:", err)
			return
		}
		combinedContent += string(depContent) + "\n"
	}

	combinedContent += string(content)

	// Output the combined content
	// fmt.Println(combinedContent)
}

func findDependencies(content string, filePath string) []string {
	var dependencies []string

	lines := strings.Split(content, "\n")
	packageName := ""
	for _, line := range lines {
		if strings.HasPrefix(line, "package ") {
			packageName = strings.TrimPrefix(line, "package ")
			packageName = strings.TrimSuffix(packageName, ";")
			break
		}
	}
	if packageName == "" {
		panic("package name not found: " + filePath)
	}

	packagePath := strings.ReplaceAll(packageName, ".", "/")
	index := strings.LastIndex(filePath, packagePath)
	if index == -1 {
		panic("failed to find package root path: " + filePath)
	}
	packageRootDirPath := filePath[:index]

	for _, line := range lines {
		if strings.HasPrefix(line, "import ") {
			importLine := strings.TrimPrefix(line, "import ")
			importLine = strings.TrimSuffix(importLine, ";")
			importLine = strings.ReplaceAll(importLine, ".", "/")
			fp := packageRootDirPath + importLine + ".java"

			if _, err := os.Stat(fp); err == nil {
				dependencies = append(dependencies, fp)
			}
		}
	}

	return dependencies
}
