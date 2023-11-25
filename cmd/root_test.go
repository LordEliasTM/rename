package cmd

import (
	"os"
	"strconv"
	"testing"
)

func TestParseArgs(t *testing.T) {
	args1 := []string{"/asd/\r\n", "kekw\n"}
	regex1, result1, err1 := ParseArgs(args1, false)

	if regex1.String() != "asd" || result1 != "kekw" || err1 != nil {
		t.Errorf("ParseArgs1: %q, %q, %q", regex1.String(), result1, err1)
	}

	args2 := []string{"/asd(/", "kekw"}
	_, _, err2 := ParseArgs(args2, false)

	if err2 == nil {
		t.Errorf("Didn't get error for invalid regex: %q", args2[0])
	}
}

func TestRenameInCurrentDir(t *testing.T) {
	os.Chdir(t.TempDir())

	for i := 0; i < 20; i++ {
		os.Create("asd" + strconv.Itoa(i) + ".txt")
	}

	args1 := []string{"/asd/", "kekw"}
	regex, replace, err := ParseArgs(args1, false)

	if regex.String() != "asd" || replace != "kekw" || err != nil {
		t.Errorf("ParseArgs failed: %q, %q, %q", regex.String(), replace, err)
	}

	content := []byte("y")
	tmpfile, _ := os.CreateTemp(t.TempDir(), "stdin_tmp")
	defer os.Remove(t.TempDir() + "/" + tmpfile.Name())

	tmpfile.Write(content)
	tmpfile.Seek(0, 0)

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	os.Stdin = tmpfile

	RenameInCurrentDir(regex, replace)

	elems, _ := os.ReadDir(".")
	t.Errorf("%q", elems)

	tmpfile.Close()
}
