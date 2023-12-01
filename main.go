package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprint(os.Stderr, "Enter the file path: ")
	scanner.Scan()
	filePath := scanner.Text()

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	targetFile := string(content)

	dependencies, err := findDependencies(targetFile, filePath)
	if err != nil {
		panic(err)
	}

	template := ""
	if _, err := os.Stat("./template/default.txt"); err == nil {
		b, err := os.ReadFile("./template/default.txt")
		if err == nil {
			template = string(b)
		}
	}
	fmt.Println(template)
	fmt.Println(targetFile)

	for _, depPath := range dependencies {
		fileContent, err := os.ReadFile(depPath)
		if err != nil {
			panic("Error reading file: " + depPath)
		}
		fmt.Println(string(fileContent))
	}
}

// findDependencies parses the Java file content and returns file paths
// of all dependencies that are part of the same package.
func findDependencies(content string, filePath string) ([]string, error) {
	lines := strings.Split(content, "\n")
	packageName, err := extractPackageName(lines)
	if err != nil {
		return nil, fmt.Errorf("extracting package name: %w", err)
	}

	packageRootDirPath, err := findPackageRootDirPath(filePath, packageName)
	if err != nil {
		return nil, fmt.Errorf("finding package root directory path: %w", err)
	}

	return parseImports(lines, packageRootDirPath), nil
}

// extractPackageName extracts the package name from the given lines of a Java file.
func extractPackageName(lines []string) (string, error) {
	for _, line := range lines {
		if strings.HasPrefix(line, "package ") {
			packageName := strings.TrimPrefix(line, "package ")
			return strings.TrimSuffix(packageName, ";"), nil
		}
	}
	return "", fmt.Errorf("package name not found")
}

// findPackageRootDirPath finds the root directory path of the package in the file path.
func findPackageRootDirPath(filePath, packageName string) (string, error) {
	packagePath := strings.ReplaceAll(packageName, ".", "/")
	index := strings.LastIndex(filePath, packagePath)
	if index == -1 {
		return "", fmt.Errorf("package root path not found")
	}
	return filePath[:index], nil
}

// parseImports processes the provided lines from a Java file, identifies import statements,
// and returns a slice of file paths for those imports that are from the same package.
func parseImports(lines []string, packageRootDirPath string) []string {
	var dependencies []string

	for _, line := range lines {
		// Skip non-import lines quickly
		if !strings.HasPrefix(line, "import ") {
			continue
		}

		// Process import line
		importPath := extractImportPath(line)
		fullPath := filepath.Join(packageRootDirPath, importPath) + ".java"

		// Check if file exists and add to dependencies if it does
		if _, err := os.Stat(fullPath); err == nil {
			dependencies = append(dependencies, fullPath)
		}
	}

	return dependencies
}

// extractImportPath takes an import line from a Java file and converts it into a file path.
func extractImportPath(importLine string) string {
	importLine = strings.TrimPrefix(importLine, "import ")
	importLine = strings.TrimSuffix(importLine, ";")
	return strings.ReplaceAll(importLine, ".", "/")
}
