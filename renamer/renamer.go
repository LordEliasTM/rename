package renamer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/LordEliasTM/rename/utils"
	"github.com/dlclark/regexp2"
	"github.com/eiannone/keyboard"
)

type MatchedRename struct {
	oldName string
	newName string
}

// Walks shitty file tree from bottom to top
// Memory intensive, so better get dad's credit card
// to buy some more of that juicy gigabytes
func RenameDeep(regex *regexp2.Regexp, replace string, all bool, onlyDirs bool, alsoFiles bool, yes bool) {
	var asd []*Path

	filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		name := d.Name()
		isDir := d.IsDir()

		if name == "." || name == ".." {
			return nil
		}
		if !all && utils.IsHidden(path) {
			if d.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}
		if !onlyDirs && isDir {
			return nil
		}
		if onlyDirs && !alsoFiles && !isDir {
			return nil
		}

		sPath := PathStringToStruct(path, d)
		asd = append(asd, sPath)
		return nil
	})

	// sort descending by path depth
	sort.Slice(asd, func(i, j int) bool {
		return asd[i].depth > asd[j].depth
	})

	var matchedRenames []MatchedRename

	for _, p := range asd {
		name := p.dirEntry.Name()
		path := p.full

		result, _ := regex.Replace(name, replace, -1, 1)

		if name == result {
			continue
		}

		filenameRegex := regexp2.MustCompile(name+"$", regexp2.None)
		pathWithoutFilename, _ := filenameRegex.Replace(path, "", -1, 1)
		result = pathWithoutFilename + result

		matchedRenames = append(matchedRenames, MatchedRename{path, result})

		fmt.Println(p.depth, path, "->", result)
	}

	RenameMatched(matchedRenames, yes)
}

func RenameInCurrentDir(regex *regexp2.Regexp, replace string, all bool, onlyDirs bool, alsoFiles bool, yes bool) {
	// get all entries in this current directory
	entries, err := os.ReadDir(".")

	if err != nil {
		fmt.Println(err)
		return
	}

	// list with the files that matched the regex
	var matchedRenames []MatchedRename

	for _, entry := range entries {
		// those if statements look quite messy and noone will understand
		// them. I hate myself for this but wharever
		if utils.IsHidden(entry.Name()) && !all {
			continue
		}
		// ignore dirs by default, except when -d flag
		if entry.IsDir() && !onlyDirs {
			fmt.Println(entry.Name())
			continue
		}
		if !entry.IsDir() && onlyDirs && !alsoFiles {
			continue
		}

		name := entry.Name()
		result, _ := regex.Replace(name, replace, -1, 1)

		// skip if filename didnt change
		if name == result {
			continue
		}

		fmt.Println(name, " -> ", result)

		matchedRenames = append(matchedRenames, MatchedRename{name, result})
	}

	RenameMatched(matchedRenames, yes)
}

func RenameMatched(matchedRenames []MatchedRename, yes bool) {
	if len(matchedRenames) == 0 {
		fmt.Println("No regex match!")
		return
	}

	if !yes && !AcceptsChanges() {
		return
	}

	for _, match := range matchedRenames {
		os.Rename(match.oldName, match.newName)
	}
}

func AcceptsChanges() bool {
	fmt.Println("\nAccept these changes? (y/N)")
	// immediately open and close to not interfer with the stdinput when user doesnt specify arguments
	keyboard.Open()
	char, _, _ := keyboard.GetSingleKey()
	keyboard.Close()
	return char == 'y'
}
