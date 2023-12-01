package renamer

import (
	"io/fs"
	"os"
	"strings"
)

type Path struct {
	Full     string
	Parts    []string
	Depth    int
	DirEntry fs.DirEntry
}

func PathStringToStruct(path string, dirEntry fs.DirEntry) *Path {
	parts := strings.Split(path, string(os.PathSeparator))
	return &Path{
		Full:     path,
		Parts:    parts,
		Depth:    len(parts),
		DirEntry: dirEntry,
	}
}
