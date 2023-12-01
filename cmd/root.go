/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/LordEliasTM/rename/renamer"

	"github.com/dlclark/regexp2"
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
		yes, _ := cmd.Flags().GetBool("yes")

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

		if recursive {
			renamer.RenameDeep(regex, replace, all, onlyDirs, alsoFiles, yes)
		} else {
			renamer.RenameInCurrentDir(regex, replace, all, onlyDirs, alsoFiles, yes)
		}
	},
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

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("insensitive", "i", false, "case insensitive")
	rootCmd.Flags().BoolP("recursive", "r", false, "recursively rename in sub-directories")
	rootCmd.Flags().BoolP("all", "a", false, "do not ignore entries starting with .")
	rootCmd.Flags().BoolP("directories", "d", false, "rename directories")
	rootCmd.Flags().BoolP("files", "f", false, "rename files, enabled by default unless -d is present")
	rootCmd.Flags().BoolP("yes", "y", false, "skip prompt to accept changes")
}
