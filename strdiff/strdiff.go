package strdiff

import (
	"bytes"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func CoolFunc(diffs []diffmatchpatch.Diff) string {
	var buff1 bytes.Buffer
	var buff2 bytes.Buffer

	for _, diff := range diffs {
		text := diff.Text

		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			_, _ = buff2.WriteString("\x1b[42m")
			_, _ = buff2.WriteString(text)
			_, _ = buff2.WriteString("\x1b[0m")
		case diffmatchpatch.DiffDelete:
			_, _ = buff1.WriteString("\x1b[41m")
			_, _ = buff1.WriteString(text)
			_, _ = buff1.WriteString("\x1b[0m")
		case diffmatchpatch.DiffEqual:
			_, _ = buff1.WriteString(text)
			_, _ = buff2.WriteString(text)
		}
	}

	buff1.WriteString(" -> ")
	buff1.WriteString(buff2.String())
	return buff1.String()
}

func SheeshDiff(old string, new string) string {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(old, new, false)

	return CoolFunc(diffs)
}
