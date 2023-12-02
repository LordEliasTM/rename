package strdiff

import (
	"bytes"

	"github.com/fatih/color"
	"github.com/sergi/go-diff/diffmatchpatch"
)

var red = color.New(color.BgRed).Add(color.FgBlack)
var green = color.New(color.BgGreen).Add(color.FgBlack)

func CoolFunc(diffs []diffmatchpatch.Diff) string {
	var buff1 bytes.Buffer
	var buff2 bytes.Buffer

	for _, diff := range diffs {
		text := diff.Text

		switch diff.Type {
		case diffmatchpatch.DiffDelete:
			_, _ = buff1.WriteString(red.Sprint(text))
		case diffmatchpatch.DiffInsert:
			_, _ = buff2.WriteString(green.Sprint(text))
		case diffmatchpatch.DiffEqual:
			_, _ = buff1.WriteString(text)
			_, _ = buff2.WriteString(text)
		}
	}

	// merge buffers and add "->" in between
	buff1.WriteString(" -> ")
	buff1.WriteString(buff2.String())
	return buff1.String()
}

func SheeshDiff(old string, new string) string {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(old, new, false)

	return CoolFunc(diffs)
}
