package renamer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/LordEliasTM/rename/cmd"
	"github.com/LordEliasTM/rename/renamer"
)

func TestRenameDeep(t *testing.T) {
	wd, _ := os.Getwd()
	os.Chdir(t.TempDir())

	// create testing files

	paths := []string{
		filepath.Join("uhm", ".asduhm", "0"),
		filepath.Join("uhm", "kekuhm", "0"),
		filepath.Join("uhm", "uhm", ".huhm"),
		filepath.Join("uhm", "uhm", "fuhm"),
		filepath.Join("uhm", "uhm", "uhm", "0"),
	}

	for _, path := range paths {
		file, _ := os.Create(path)
		file.Close()
	}

	// do test

	args := []string{"/uhm/", "xxx"}
	regex, replace, _ := cmd.ParseArgs(args, false, false)

	renamer.RenameDeep(regex, replace, false, false, false, true, 1)

	//TODO finish writing test

	// cleanup

	os.Chdir(wd)
}
