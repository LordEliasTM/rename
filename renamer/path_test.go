package renamer_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/LordEliasTM/rename/renamer"
)

func TestPathStringToStruct(t *testing.T) {
	wd, _ := os.Getwd()
	os.Chdir(t.TempDir())

	path := filepath.Join("path", "test", "asd.txt")

	file, _ := os.Create(path)
	file.Close()
	info, _ := fs.Stat(os.DirFS("."), path)
	entry := fs.FileInfoToDirEntry(info)

	sPath := renamer.PathStringToStruct(path, entry)

	if sPath.Depth != 3 {
		t.Errorf("Expected sPath.Depth to be 3, but got %d", sPath.Depth)
	}

	os.Chdir(wd)
}
