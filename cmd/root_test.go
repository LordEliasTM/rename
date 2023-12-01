package cmd

import (
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
