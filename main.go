package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var dryRun bool
var gitRootOnly bool
var red = "\033[31m"
var reset = "\033[0m"

func init() {
	flag.BoolVar(&dryRun, "dry-run", false, "Perform a dry run without deleting files")
	flag.BoolVar(&gitRootOnly, "git-root-only", false, "Only clean directories that are in the root of a git repository")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "[flags] path_to_workspace_directory")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "%sError: You must provide exactly one path to a workspace directory%s\n", red, reset)
		flag.Usage()
		os.Exit(1)
	}
	startDir := args[0]

	if !dryRun && !confirmAction() {
		fmt.Fprintf(os.Stderr, "Operation cancelled.\n")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Deleting files in %s\nDry run: %t\nGit root only: %t\n", startDir, dryRun, gitRootOnly)

	err := filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				fmt.Fprintf(os.Stderr, "%sError: %s\n%s", red, err, reset)
				return nil
			}
			return err
		}

		// Skip .git folder to speed up scan
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		// Node.js Projects: Delete node_modules
		if info.IsDir() && info.Name() == "node_modules" {
			parentDir := filepath.Dir(path)
			if gitRootOnly && !isGitRoot(parentDir) {
				return nil
			}
			if hasFile(parentDir, "package.json") {
				fmt.Printf("Found node_modules to delete: \t%s\n", path)
				if !dryRun {
					err := os.RemoveAll(path)
					if err != nil {
						return err
					}
				}
				return filepath.SkipDir
			}
		}

		// PHP Projects: Delete vendor
		if info.IsDir() && info.Name() == "vendor" {
			parentDir := filepath.Dir(path)
			if gitRootOnly && !isGitRoot(parentDir) {
				return nil
			}
			if hasFile(parentDir, "composer.json") {
				fmt.Printf("Found vendor to delete: \t%s\n", path)
				if !dryRun {
					err := os.RemoveAll(path)
					if err != nil {
						return err
					}
				}
				return filepath.SkipDir
			}
		}

		// Symfony Projects: Delete var/log
		if info.IsDir() && (filepath.Base(path) == "log" || filepath.Base(path) == "cache") && filepath.Base(filepath.Dir(path)) == "var" {
			parentDir := filepath.Dir(filepath.Dir(path)) // Two levels up from var/log or var/cache
			if gitRootOnly && !isGitRoot(parentDir) {
				return nil
			}
			if hasFile(parentDir, "symfony.lock") {
				fmt.Printf("Found var/log to delete: \t%s\n", path)
				if !dryRun {
					err := os.RemoveAll(path)
					if err != nil {
						return err
					}
				}
				return filepath.SkipDir
			}
		}

		// Flutter Projects: Delete build directory and .dart_tool
		if info.IsDir() && (info.Name() == "build" || info.Name() == ".dart_tool") {
			parentDir := filepath.Dir(path)
			if gitRootOnly && !isGitRoot(parentDir) {
				return nil
			}
			if hasFile(parentDir, "pubspec.yaml") {
				fmt.Printf("Found Flutter %s to delete: \t%s\n", info.Name(), path)
				if !dryRun {
					err := os.RemoveAll(path)
					if err != nil {
						return err
					}
				}
				return filepath.SkipDir
			}
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError: %s\n%s", red, err, reset)
	}
}

// hasFile checks if the specified directory contains the file
func hasFile(dir string, filename string) bool {
	expectedPath := filepath.Join(dir, filename)
	_, err := os.Stat(expectedPath)
	return !os.IsNotExist(err)
}

// User confirmation function
func confirmAction() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintf(os.Stderr, "Are you sure you want to proceed with deletions? (y/n): ")
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError reading response: %s\n%s", red, err, reset)
		return false
	}
	response = strings.TrimSpace(response)
	return response == "y" || response == "Y"
}

// isGitRoot checks if the given path contains a .git folder
func isGitRoot(path string) bool {
	gitPath := filepath.Join(path, ".git")
	_, err := os.Stat(gitPath)
	return !os.IsNotExist(err)
}
