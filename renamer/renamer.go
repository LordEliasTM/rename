package renamer

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/LordEliasTM/rename/strdiff"
	"github.com/LordEliasTM/rename/utils"
	"github.com/dlclark/regexp2"
	"github.com/eiannone/keyboard"

	. "github.com/LordEliasTM/rename/types"
)

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
		return asd[i].Depth > asd[j].Depth
	})

	var matchedRenames []MatchedRename

	for _, p := range asd {
		name := p.DirEntry.Name()
		path := p.Full

		result, _ := regex.Replace(name, replace, -1, 1)

		if name == result {
			continue
		}

		filenameRegex := regexp2.MustCompile(name+"$", regexp2.None)
		pathWithoutFilename, _ := filenameRegex.Replace(path, "", -1, 1)
		result = pathWithoutFilename + result

		matchedRenames = append(matchedRenames, MatchedRename{path, result})

		fmt.Println(strdiff.SheeshDiff(path, result))
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

		fmt.Println(strdiff.SheeshDiff(name, result))

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
		// if file exists
		if _, err := os.Stat(match.NewName); err == nil {
			fmt.Printf("File %s already exists! (o = overwrite, s = skip, r = choose new name, c = cancel)\n", match.NewName)
			keyboard.Open()
			char, _, _ := keyboard.GetSingleKey()
			keyboard.Close()
			switch char {
			case 'o':
				break
			case 's':
				continue
			case 'r':
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("New Name (write full path): ")
				newerName, _ := reader.ReadString('\n')
				newerName = strings.TrimRight(newerName, "\r\n")
				match.NewName = newerName
				break
			case 'c':
				return
			default:
				continue
			}
		}
		os.Rename(match.OldName, match.NewName)
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
