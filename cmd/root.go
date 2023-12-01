/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/LordEliasTM/rename/utils"

	"github.com/dlclark/regexp2"
	"github.com/eiannone/keyboard"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rename <regex> <replace>",
	Short: "Rename all files in the current folder with regex",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 && len(args) != 2 {
			argNumStr := strconv.Itoa(len(args))
			return errors.New("accetps 0 or 2 args, recieved " + argNumStr)
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		recursive, _ := cmd.Flags().GetBool("recursive")
		insensitive, _ := cmd.Flags().GetBool("insensitive")
		all, _ := cmd.Flags().GetBool("all")
		onlyDirs, _ := cmd.Flags().GetBool("directories")
		alsoFiles, _ := cmd.Flags().GetBool("files")

		var regex *regexp2.Regexp
		var replace string
		var err error

		if len(args) == 2 {
			// parse regex and replace from args
			regex, replace, err = ParseArgs(args, insensitive)
		} else if len(args) == 0 {
			// read regex and replace from stdin
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Regex: ")
			regexStr, _ := reader.ReadString('\n')
			fmt.Print("Replace: ")
			replace, _ = reader.ReadString('\n')

			regex, replace, err = ParseArgs([]string{regexStr, replace}, insensitive)
		}

		if err != nil {
			fmt.Println(err)
			return
		}

		//fmt.Println(regex)
		//fmt.Println(replace)
		if recursive && onlyDirs {
			HateMyLifeRenameDirsTreeWalkerThingyProbablyBuggyAsHell2(regex, replace, all, onlyDirs, alsoFiles)
		} else if recursive {
			RenameRecursively(regex, replace, all, onlyDirs, alsoFiles)
		} else {
			RenameInCurrentDir(regex, replace, all, onlyDirs, alsoFiles)
		}
	},
}

type NestedMap map[string]NestedMap
type Path struct {
	full     string
	parts    []string
	depth    int
	dirEntry fs.DirEntry
}

func PathStringToStruct(path string, dirEntry fs.DirEntry) *Path {
	parts := strings.Split(path, string(os.PathSeparator))
	return &Path{
		full:     path,
		parts:    parts,
		depth:    len(parts),
		dirEntry: dirEntry,
	}
}

// Walks shitty file tree from bottom to top
// Memory intensive, so better get dad's credit card
// to buy some more of that juicy gigabytes
func HateMyLifeRenameDirsTreeWalkerThingyProbablyBuggyAsHell2(regex *regexp2.Regexp, replace string, all bool, onlyDirs bool, alsoFiles bool) {
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

	var matchedRenames [][]string

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

		matchedRenames = append(matchedRenames, []string{path, result})

		fmt.Println(p.depth, path, "->", result)
	}

	RenameMatched(matchedRenames)
}

func HateMyLifeRenameDirsTreeWalkerThingyProbablyBuggyAsHell() {
	m := make(NestedMap)

	filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		sPath := PathStringToStruct(path, d)
		current := m
		for _, v := range sPath.parts {
			if current[v] == nil {
				current[v] = make(NestedMap)
			}
			current = current[v]
		}
		fmt.Println(len(m))
		return nil
	})

	regex, replace, _ := ParseArgs([]string{"a", "b"}, false)

	DoSomeRecursiveAssDogballGarbage(m, regex, replace)
}

func DoSomeRecursiveAssDogballGarbage(ballz NestedMap, regex *regexp2.Regexp, replace string) {
	for k, v := range ballz {
		DoSomeRecursiveAssDogballGarbage(v, regex, replace)
		fmt.Println(k)
		time.Sleep(10 * time.Millisecond)
	}
}

func RenameRecursively(regex *regexp2.Regexp, replace string, all bool, onlyDirs bool, alsoFiles bool) {
	// list with the files that matched the regex
	var matchedRenames [][]string

	// Walk the directory tree and files
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		// skip . and ..
		if info.Name() == "." || info.Name() == ".." {
			return nil
		}
		// exclude hidden files/folders by default, "all" allows hidden files/folders
		if utils.IsHidden(path) && !all {
			// skip whole tree
			return filepath.SkipDir
		}
		// ignore dirs by default, except when -d flag
		if info.IsDir() && !onlyDirs {
			fmt.Println(info.Name())
			return nil
		}
		// if file && only rename dirs && but not also rename files
		if !info.IsDir() && onlyDirs && !alsoFiles {
			return nil
		}

		//TODO Fix this folder renaming bug:
		//   uhm -> khm
		//   uhm/kek/uha -> uhm/kek/kha
		// second one will not work because parent folder is already renamed

		name := info.Name()
		result, _ := regex.Replace(name, replace, -1, 1)

		// skip if filename didnt change
		if name == result {
			return nil
		}

		// My garbage go skills
		// Removes file name from end of path "asd/kek/file" -> "asd/kek/"
		// and then appends the result to it "asd/kek/" -> "asd/kek/fefle"
		filenameRegex := regexp2.MustCompile(name+"$", regexp2.None)
		pathWithoutFilename, _ := filenameRegex.Replace(path, "", -1, 1)
		result = pathWithoutFilename + result

		fmt.Println(path, " -> ", result)

		matchedRenames = append(matchedRenames, []string{path, result})

		return nil
	})

	RenameMatched(matchedRenames)
}

func RenameInCurrentDir(regex *regexp2.Regexp, replace string, all bool, onlyDirs bool, alsoFiles bool) {
	// get all entries in this current directory
	entries, err := os.ReadDir(".")

	if err != nil {
		fmt.Println(err)
		return
	}

	// list with the files that matched the regex
	var matchedRenames [][]string

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

		matchedRenames = append(matchedRenames, []string{name, result})
	}

	RenameMatched(matchedRenames)
}

func RenameMatched(matchedRenames [][]string) {
	if len(matchedRenames) == 0 {
		fmt.Println("No regex match!")
		return
	}

	if !AcceptsChanges() {
		return
	}

	for _, match := range matchedRenames {
		// match is in format [oldName, newName]
		os.Rename(match[0], match[1])
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

func ParseArgs(args []string, insensitive bool) (*regexp2.Regexp, string, error) {
	regexStr := args[0]
	replace := args[1]

	// reading stdin also reads the newline, so need to remove that
	regexStr = strings.TrimRight(regexStr, "\r\n")
	regexStr = strings.Trim(regexStr, "/")
	replace = strings.TrimRight(replace, "\r\n")

	// case insensitive
	if insensitive {
		regexStr = "(?i)" + regexStr
	}

	regex, err := regexp2.Compile(regexStr, regexp2.None)

	return regex, replace, err
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rename.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("insensitive", "i", false, "case insensitive")
	rootCmd.Flags().BoolP("recursive", "r", false, "recursively rename in sub-directories")
	rootCmd.Flags().BoolP("all", "a", false, "do not ignore entries starting with .")
	rootCmd.Flags().BoolP("directories", "d", false, "rename directories")
	rootCmd.Flags().BoolP("files", "f", false, "rename files, enabled by default unless -d is present")
}
