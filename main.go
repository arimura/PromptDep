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
	targetFile := string(content)

	dependencies := findDependencies(targetFile, filePath)

	template := ""
	if _, err := os.Stat("./template/default.txt"); err == nil {
		b, err := ioutil.ReadFile("./template/default.txt")
		if err != nil {
			fmt.Println("Error reading default file:", err)
			return
		}
		template = string(b)
	}
	println(template)
	println(targetFile)

	for _, depPath := range dependencies {
		fileContent, err := ioutil.ReadFile(depPath)
		if err != nil {
			panic("Error reading file: " + depPath)
		}
		fmt.Println(string(fileContent))
	}
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
