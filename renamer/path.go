package renamer

import (
	"io/fs"
	"os"
	"strings"

	. "github.com/LordEliasTM/rename/types"
)

func PathStringToStruct(path string, dirEntry fs.DirEntry) *Path {
	parts := strings.Split(path, string(os.PathSeparator))
	return &Path{
		Full:     path,
		Parts:    parts,
		Depth:    len(parts),
		DirEntry: dirEntry,
	}
}
