package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHandleDir(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		fs             func() *filesystem
		expectedOutput []string
	}{
		{
			name: "root with no subdirs",
			fs:   CreateFilesystem,
			expectedOutput: []string{
				"Directory of root:",
				"No subdirectories",
			},
		},
		{
			name: "root with subdirs",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				fs.AddSubdir("sub2")
				return fs
			},
			expectedOutput: []string{
				"Directory of root:",
				"sub1    sub2",
			},
		},
		{
			name: "inside dir with subdirs",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				fs.AddSubdir("sub2")
				fs.current = fs.current.subs[0]
				fs.AddSubdir("sub4")
				fs.AddSubdir("sub3")
				return fs
			},
			expectedOutput: []string{
				"Directory of root\\sub1:",
				"sub3    sub4",
			},
		},
		{
			name: "exactly 10 subdirs should not wrap output",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				fs.AddSubdir("sub2")
				fs.AddSubdir("sub4")
				fs.AddSubdir("sub3")
				fs.AddSubdir("sub5")
				fs.AddSubdir("sub6")
				fs.AddSubdir("sub7")
				fs.AddSubdir("sub8")
				fs.AddSubdir("sub9")
				fs.AddSubdir("sub10")
				return fs
			},
			expectedOutput: []string{
				"Directory of root:",
				"sub1    sub10   sub2    sub3    sub4    sub5    sub6    sub7    sub8    sub9",
			},
		},
		{
			name: "more than 10 subdirs",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				fs.AddSubdir("sub2")
				fs.AddSubdir("sub4")
				fs.AddSubdir("sub3")
				fs.AddSubdir("sub5")
				fs.AddSubdir("sub6")
				fs.AddSubdir("sub7")
				fs.AddSubdir("sub8")
				fs.AddSubdir("sub9")
				fs.AddSubdir("sub10")
				fs.AddSubdir("sub11")
				fs.AddSubdir("sub12")
				return fs
			},
			expectedOutput: []string{
				"Directory of root:",
				"sub1    sub10   sub11   sub12   sub2    sub3    sub4    sub5    sub6    sub7",
				"sub8    sub9",
			},
		},
		{
			name: "long subdir names",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("longsubname1")
				fs.AddSubdir("longsubname2")
				fs.AddSubdir("longsubname3")
				fs.AddSubdir("longsubname4")
				fs.AddSubdir("longsubname5")
				fs.AddSubdir("longsubname6")
				return fs
			},
			expectedOutput: []string{
				"Directory of root:",
				"longsubname1    longsubname2    longsubname3    longsubname4    longsubname5",
				"longsubname6",
			},
		},
		{
			name: "subdir names with length of column width",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("longsub1")
				fs.AddSubdir("longsub2")
				fs.AddSubdir("longsub3")
				fs.AddSubdir("longsub4")
				fs.AddSubdir("longsub5")
				fs.AddSubdir("longsub6")
				return fs
			},
			expectedOutput: []string{
				"Directory of root:",
				"longsub1        longsub2        longsub3        longsub4        longsub5",
				"longsub6",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := handleCommand("dir", tt.fs())

			if diff := cmp.Diff(tt.expectedOutput, output); diff != "" {
				t.Fatalf("output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleMkdir(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		cmdArg           string
		fs               func() *filesystem
		expectedOutput   []string
		expectedSubNames []string
	}{
		{
			name:             "make new dir when no subdir exists",
			cmdArg:           "sub1",
			fs:               CreateFilesystem,
			expectedOutput:   nil,
			expectedSubNames: []string{"sub1"},
		},
		{
			name:   "make new dir when subdirs exist",
			cmdArg: "sub1",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub2")
				fs.AddSubdir("sub3")
				return fs
			},
			expectedOutput:   nil,
			expectedSubNames: []string{"sub1", "sub2", "sub3"},
		},
		{
			name:   "cannot make subdir with same name as existing",
			cmdArg: "sub1",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				return fs
			},
			expectedOutput:   []string{ErrSubdirAlreadyExists.Error()},
			expectedSubNames: []string{"sub1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.fs()
			output := handleCommand("mkdir   "+tt.cmdArg, fs)

			if diff := cmp.Diff(tt.expectedOutput, output); diff != "" {
				t.Fatalf("output mismatch (-want +got):\n%s", diff)
			}
			subDirNames := []string{}
			for _, subdir := range fs.current.subs {
				subDirNames = append(subDirNames, subdir.name)
			}
			if diff := cmp.Diff(tt.expectedSubNames, subDirNames); diff != "" {
				t.Fatalf("output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleUp(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name               string
		fs                 func() *filesystem
		expectedOutput     []string
		expectedCurrentDir func(fs *filesystem) *dir
	}{
		{
			name: "move up a dir",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				fs.current = fs.current.subs[0]
				return fs
			},
			expectedOutput: nil,
			expectedCurrentDir: func(fs *filesystem) *dir {
				return fs.current.parent
			},
		},
		{
			name:           "cannot move up from root",
			fs:             CreateFilesystem,
			expectedOutput: []string{ErrCannotMoveUpFromRoot.Error()},
			expectedCurrentDir: func(fs *filesystem) *dir {
				return fs.current
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.fs()
			expectedCurrent := tt.expectedCurrentDir(fs)
			output := handleCommand("up", fs)
			if diff := cmp.Diff(tt.expectedOutput, output); diff != "" {
				t.Fatalf("output mismatch (-want +got):\n%s", diff)
			}
			if expectedCurrent != fs.current {
				t.Fatalf("current dir mismatch:\nwant: %s\ngot: %s", expectedCurrent.name, fs.current.name)
			}

		})
	}
}

func TestHandleCd(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name               string
		cmdArg             string
		fs                 func() *filesystem
		expectedOutput     []string
		expectedCurrentDir func(fs *filesystem) *dir
	}{
		{
			name:   "move to subdir",
			cmdArg: "sub1",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				fs.AddSubdir("sub2")
				return fs
			},
			expectedOutput: nil,
			expectedCurrentDir: func(fs *filesystem) *dir {
				return fs.current.subs[0]
			},
		},
		{
			name:   "cannot move to non existing subdir",
			cmdArg: "sub3",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				fs.AddSubdir("sub2")
				return fs
			},
			expectedOutput: []string{ErrSubdirDoesNotExist.Error()},
			expectedCurrentDir: func(fs *filesystem) *dir {
				return fs.current
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.fs()
			expectedCurrent := tt.expectedCurrentDir(fs)
			output := handleCommand("cd      "+tt.cmdArg, fs)
			if diff := cmp.Diff(tt.expectedOutput, output); diff != "" {
				t.Fatalf("output mismatch (-want +got):\n%s", diff)
			}
			if expectedCurrent != fs.current {
				t.Fatalf("current dir mismatch:\nwant: %s\ngot: %s", expectedCurrent.name, fs.current.name)
			}
		})
	}
}

func TestHandleTree(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		fs             func() *filesystem
		expectedOutput []string
	}{
		{

			name: "root with no subdirs",
			fs:   CreateFilesystem,
			expectedOutput: []string{
				"Tree of root:",
				".",
			},
		},
		{
			name: "root with 3 subdir levels",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				fs.AddSubdir("sub2")
				fs.current = fs.current.subs[0]
				fs.AddSubdir("sub3")
				fs.current = fs.current.subs[0]
				fs.AddSubdir("sub4")
				fs.current = fs.root
				return fs
			},
			expectedOutput: []string{
				"Tree of root:",
				".",
				"├── sub1",
				"│   └── sub3",
				"│       └── sub4",
				"└── sub2",
			},
		},
		{
			name: "root with 3 subdir levels with more subdirs",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				fs.current = fs.current.subs[0]
				fs.AddSubdir("sub3")
				fs.current = fs.current.subs[0]
				fs.AddSubdir("sub4")
				fs.AddSubdir("sub5")
				fs.current = fs.root
				return fs
			},
			expectedOutput: []string{
				"Tree of root:",
				".",
				"└── sub1",
				"    └── sub3",
				"        ├── sub4",
				"        └── sub5",
			},
		},
		{

			name: "subdir with no subdirs",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				fs.AddSubdir("sub2")
				fs.current = fs.current.subs[0]
				return fs
			},
			expectedOutput: []string{
				"Tree of root\\sub1:",
				".",
			},
		},
		{
			name: "root with 2 subdir levels one subdir per level",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				fs.current = fs.current.subs[0]
				fs.AddSubdir("sub3")
				fs.current = fs.current.subs[0]
				fs.AddSubdir("sub4")
				fs.current = fs.root
				return fs
			},
			expectedOutput: []string{
				"Tree of root:",
				".",
				"└── sub1",
				"    └── sub3",
				"        └── sub4",
			},
		},
		{
			name: "root with complicated subdir levels",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				fs.AddSubdir("sub2")
				fs.current = fs.current.subs[0]
				fs.AddSubdir("sub3")
				fs.AddSubdir("sub4")
				fs.current = fs.current.subs[0]
				fs.AddSubdir("sub4")
				fs.current = fs.root
				return fs
			},
			expectedOutput: []string{
				"Tree of root:",
				".",
				"├── sub1",
				"│   ├── sub3",
				"│   │   └── sub4",
				"│   └── sub4",
				"└── sub2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.fs()
			output := handleCommand("tree", fs)
			if diff := cmp.Diff(tt.expectedOutput, output); diff != "" {
				for _, o := range output {
					t.Log(o)
				}
				t.Fatalf("output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleMv(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		argFrom        string
		argTo          string
		fs             func() *filesystem
		expectedOutput []string
	}{
		{
			name:           "mv non existing subdir",
			argFrom:        "nosub",
			argTo:          "sub1",
			fs:             CreateFilesystem,
			expectedOutput: []string{ErrSubdirDoesNotExist.Error()},
		},
		{
			name:    "rename subdir",
			argFrom: "sub1",
			argTo:   "sub2",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				return fs
			},
			expectedOutput: nil,
		},
		{
			name:    "rename subdir using relative path",
			argFrom: "sub1",
			argTo:   ".\\sub2",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				return fs
			},
			expectedOutput: nil,
		},
		{
			name:    "move and rename subdir",
			argFrom: "sub1",
			argTo:   "sub2\\sub111",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				fs.AddSubdir("sub2")
				return fs
			},
			expectedOutput: nil,
		},
		{
			name:    "subdirectory already exists",
			argFrom: "sub1",
			argTo:   "sub2",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				fs.AddSubdir("sub2")
				fs.Cd("sub2")
				fs.AddSubdir("sub1")
				fs.current = fs.root
				return fs
			},
			expectedOutput: []string{ErrSubdirAlreadyExists.Error()},
		},
		{
			name:    "cannot move to illegal path",
			argFrom: "sub1",
			argTo:   "..",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				fs.current = fs.root
				return fs
			},
			expectedOutput: []string{ErrSubdirDoesNotExist.Error()},
		},
		{
			name:    "cannot move to illegal intermediate path",
			argFrom: "sub1",
			argTo:   "sub2\\..\\..\\root",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				fs.AddSubdir("sub2")
				fs.Cd("sub2")
				fs.AddSubdir("sub1")
				fs.current = fs.root
				return fs
			},
			expectedOutput: []string{ErrSubdirDoesNotExist.Error()},
		},
		{
			name:    "parent of destination does not exist",
			argFrom: "sub1",
			argTo:   "sub2\\sub3\\sub4",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				fs.AddSubdir("sub2")
				return fs
			},
			expectedOutput: []string{ErrSubdirDoesNotExist.Error()},
		},
		{
			name:    "complicated relative path",
			argFrom: "sub1",
			argTo:   "sub2\\.\\..\\sub2\\.\\.\\sub1",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				fs.AddSubdir("sub2")
				fs.Cd("sub2")
				fs.AddSubdir("sub1")
				fs.current = fs.root
				return fs
			},
			expectedOutput: nil,
		},
		{
			name:    "move to current dir",
			argFrom: "sub1",
			argTo:   ".",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				return fs
			},
			expectedOutput: nil,
		},
		{
			name:    "move to current dir with same name specified",
			argFrom: "sub1",
			argTo:   "sub1",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				return fs
			},
			expectedOutput: nil,
		},
		{
			name:    "move to current complicated relative path",
			argFrom: "sub1",
			argTo:   "sub2\\..",
			fs: func() *filesystem {
				fs := CreateFilesystem()
				fs.AddSubdir("sub1")
				fs.AddSubdir("sub2")
				fs.Cd("sub2")
				fs.AddSubdir("sub1")
				fs.current = fs.root
				return fs
			},
			expectedOutput: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.fs()
			output := handleCommand("mv "+tt.argFrom+" "+tt.argTo, fs)
			if diff := cmp.Diff(tt.expectedOutput, output); diff != "" {
				t.Fatalf("output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHandleCommand(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		command        string
		fs             func() *filesystem
		shouldPanic    bool
		expectedOutput []string
	}{
		{
			name:           "empty command",
			fs:             CreateFilesystem,
			shouldPanic:    false,
			expectedOutput: nil,
		},
		{
			name:           "not known command",
			command:        "notarealcommand",
			fs:             CreateFilesystem,
			shouldPanic:    true,
			expectedOutput: nil,
		},
		{
			name:           "missing arguments",
			command:        "mkdir    ",
			fs:             CreateFilesystem,
			shouldPanic:    true,
			expectedOutput: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := tt.fs()
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected to panic")
					}
				}()
			}
			output := handleCommand(tt.command, fs)
			if diff := cmp.Diff(tt.expectedOutput, output); diff != "" {
				t.Fatalf("output mismatch (-want +got):\n%s", diff)
			}
		})
	}

}
