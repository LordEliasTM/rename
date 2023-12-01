package renamer

import (
	"io/fs"
	"os"
	"strings"
)

type Path struct {
	full     string
	parts    []string
	depth    int
	dirEntry fs.DirEntry
}

func PathStringToStruct(path string, dirEntry fs.DirEntry) *Path {
	parts := strings.Split(path, string(os.PathSeparator))
	return &Path{
		full:     path,
		parts:    parts,
		depth:    len(parts),
		dirEntry: dirEntry,
	}
}
