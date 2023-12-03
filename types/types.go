package types

import "io/fs"

type MatchedRename struct {
	OldName string
	NewName string
}

type Path struct {
	Full     string
	Parts    []string
	Depth    int
	DirEntry fs.DirEntry
}
