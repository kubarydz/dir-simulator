package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrWrongNumberOfArguments = errors.New("command has wrong number of arguments")
)

func handleCommand(input string, fs *filesystem) []string {

	switch getCommand(input) {
	case "dir":
		return handleDir(fs)
	case "mkdir":
		arg, err := getArg(input)
		if err != nil {
			panic(err)
		}
		return handleMkdir(fs, arg)
	case "up":
		return handleUp(fs)
	case "cd":
		arg, err := getArg(input)
		if err != nil {
			panic(err)
		}
		return handleCd(fs, arg)
	case "tree":
		return handleTree(fs)
	case "mv":
		arg1, arg2, err := getArgs(input)
		if err != nil {
			panic(err)
		}
		return handleMv(fs, arg1, arg2)
	case "":
		return nil
	}
	panic("command not known")
}

func getCommandEcho(input string) string {
	if len(input) == 0 {
		return ""
	}
	spaces := regexp.MustCompile(`\s+`)
	s := spaces.ReplaceAllString(input, " ")
	chunks := strings.Split(s, " ")
	echo := fmt.Sprintf("Command: %s", chunks[0])
	if len(chunks) > 1 {
		// first argument in column 18
		echo = fmt.Sprintf("%-17s%s", echo, chunks[1])
	}
	if len(chunks) > 2 {
		// second argument in column 26
		echo = fmt.Sprintf("%-25s%s", echo, chunks[2])
	}

	return echo
}

func getCommand(input string) string {
	spaces := regexp.MustCompile(`\s+`)
	s := spaces.ReplaceAllString(input, " ")
	chunks := strings.Split(s, " ")
	return chunks[0]
}

// if the command has one argument it starts in column 9
func getArg(input string) (string, error) {
	if len(input) < 9 {
		return "", ErrWrongNumberOfArguments
	}
	arg := strings.TrimSpace(input[8:])
	if len(arg) == 0 {
		return "", ErrWrongNumberOfArguments
	}
	return arg, nil
}

// for two arguments there are no column guarantees
func getArgs(input string) (string, string, error) {
	spaces := regexp.MustCompile(`\s+`)
	s := spaces.ReplaceAllString(input, " ")
	chunks := strings.Split(s, " ")

	if len(chunks) != 3 {
		return "", "", ErrWrongNumberOfArguments
	}
	return chunks[1], chunks[2], nil
}

func handleDir(fs *filesystem) []string {
	current := fs.current
	current_path := current.name + ":"
	for current != fs.root {
		current = current.parent
		current_path = current.name + "\\" + current_path
	}

	current_path = "Directory of " + current_path
	if len(current.subs) == 0 {
		return []string{current_path, "No subdirectories"}
	}

	subdirs := []string{fs.current.subs[0].name}
	lineCounter := 0
	for _, subdir := range fs.current.subs[1:] {
		// wrap lines after 10 columns of length 8
		paddingLength := 8 - len(subdirs[lineCounter])%8
		if len(subdirs[lineCounter])+len(subdir.name)+paddingLength > 80 {
			lineCounter++
			subdirs = append(subdirs, subdir.name)
			continue
		}
		padding := paddingLength + len(subdirs[lineCounter])
		subdirs[lineCounter] = fmt.Sprintf("%-*s", padding, subdirs[lineCounter])
		subdirs[lineCounter] += subdir.name
	}

	return append([]string{current_path}, subdirs...)
}

func handleMkdir(fs *filesystem, arg string) []string {
	err := fs.AddSubdir(arg)
	if err != nil {
		return []string{err.Error()}
	}
	return nil
}

func handleUp(fs *filesystem) []string {
	err := fs.Up()
	if err != nil {
		return []string{err.Error()}
	}
	return nil
}

func handleCd(fs *filesystem, arg string) []string {
	err := fs.Cd(arg)
	if err != nil {
		return []string{err.Error()}
	}
	return nil
}

func handleTree(fs *filesystem) []string {
	current := fs.current
	current_path := current.name + ":"
	for current != fs.root {
		current = current.parent
		current_path = current.name + "\\" + current_path
	}

	current_path = "Tree of " + current_path

	tree := []string{current_path, "."}
	branches := getTreeBranches(fs.current, 0)

	return append(tree, branches...)
}

func getTreeBranches(current *dir, level int) []string {
	branches := []string{}
	for i, subdir := range current.subs {
		line := ""
		if i == len(current.subs)-1 {
			line += "└── "
		} else {
			line += "├── "
		}
		line += subdir.name
		branches = append(branches, line)
		subBranches := getTreeBranches(subdir, level+1)
		for subi, subBranch := range subBranches {
			subBranches[subi] = "    " + subBranch
			if i == len(current.subs)-1 {
				continue
			}
			subBranches[subi] = strings.Replace(subBranches[subi], " ", "│", 1)
		}
		branches = append(branches, subBranches...)
	}
	return branches
}

func handleMv(fs *filesystem, from, to string) []string {
	err := fs.Mv(from, to)
	if err != nil {
		return []string{err.Error()}
	}
	return nil

}
