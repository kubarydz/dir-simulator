package main

import (
	"errors"
	"sort"
	"strings"
)

var (
	ErrSubdirAlreadyExists  = errors.New("Subdirectory already exists")
	ErrCannotMoveUpFromRoot = errors.New("Cannot move up from root directory")
	ErrSubdirDoesNotExist   = errors.New("Subdirectory does not exist")
)

type dir struct {
	name   string
	parent *dir
	subs   []*dir
}

type filesystem struct {
	current *dir
	root    *dir
}

// creates represenation of the filesystem as a file tree
// only creates root directory
func CreateFilesystem() *filesystem {
	root := dir{
		name: "root",
	}
	return &filesystem{
		current: &root,
		root:    &root,
	}
}

// adds subdirectory with given name to the current directory
// returns error if subdirectory already exists
func (fs *filesystem) AddSubdir(subName string) error {
	for _, subdir := range fs.current.subs {
		if subdir.name == subName {
			return ErrSubdirAlreadyExists
		}
	}

	fs.current.subs = append(fs.current.subs, &dir{
		name:   subName,
		parent: fs.current})

	sort.Slice(fs.current.subs, func(i, j int) bool {
		return fs.current.subs[i].name < fs.current.subs[j].name
	})
	return nil
}

// moves one directory upword
// returns error if moving up is impossible
func (fs *filesystem) Up() error {
	if fs.current == fs.root {
		return ErrCannotMoveUpFromRoot
	}
	fs.current = fs.current.parent
	return nil
}

// changes directory to given subdirectory
// does not support relative path, only direct subdirectory
// returns error if subdirectory doesn't exist
func (fs *filesystem) Cd(dirName string) error {
	for _, subdir := range fs.current.subs {
		if subdir.name == dirName {
			fs.current = subdir
			return nil
		}
	}
	return ErrSubdirDoesNotExist
}

// moves given subdirectory to destination
// creates destination if it doesn't exist
// returns error if moving is impossible
func (fs *filesystem) Mv(from, to string) error {
	var dirToMove *dir
	for _, subdir := range fs.current.subs {
		if subdir.name == from {
			dirToMove = subdir
			break
		}
	}
	if dirToMove == nil {
		return ErrSubdirDoesNotExist
	}

	destinationSteps := strings.Split(to, "\\")
	destination := fs.current

StepsLoop:
	for i, step := range destinationSteps {
		if step == "." {
			continue
		}
		if step == ".." {
			destination = destination.parent
			if destination == nil {
				return ErrSubdirDoesNotExist
			}
			continue
		}
		for _, subdir := range destination.subs {
			if subdir.name == step {
				destination = subdir
				continue StepsLoop
			}
		}

		if i == len(destinationSteps)-1 {
			dirToMove.name = step
			moveDirectory(dirToMove, destination)
			return nil
		}
		return ErrSubdirDoesNotExist
	}

	// mv to the same dir without changing name == do nothing
	if destination == dirToMove || destination == fs.current {
		return nil
	}

	// check if directory name already exists in subdirs of destination
	for _, subdir := range destination.subs {
		if subdir.name == dirToMove.name {
			return ErrSubdirAlreadyExists
		}
	}

	moveDirectory(dirToMove, destination)
	return nil
}

func moveDirectory(dirToMove *dir, destination *dir) {
	for i, subdir := range dirToMove.parent.subs {
		if subdir == dirToMove {
			newSubs := dirToMove.parent.subs[:i]
			newSubs = append(newSubs, dirToMove.parent.subs[i+1:]...)
			dirToMove.parent.subs = newSubs
		}
	}
	dirToMove.parent = destination
	destination.subs = append(destination.subs, dirToMove)
	sort.Slice(destination.subs, func(i, j int) bool {
		return destination.subs[i].name < destination.subs[j].name
	})
}
