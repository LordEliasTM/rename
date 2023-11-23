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
		//insensitive, _ := cmd.Flags().GetBool("insensitive")

		var regex *regexp2.Regexp
		var replace string
		var err error

		if len(args) == 2 {
			// parse regex and replace from args
			regex, replace, err = ParseArgs(args)
		} else if len(args) == 0 {
			// read regex and replace from stdin
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Regex: ")
			regexStr, _ := reader.ReadString('\n')
			fmt.Print("Replace: ")
			replace, _ = reader.ReadString('\n')

			regex, replace, err = ParseArgs([]string{regexStr, replace})
		}

		if err != nil {
			fmt.Println(err)
			return
		}

		//fmt.Println(regex)
		//fmt.Println(replace)

		if recursive {
			// TODO
		} else {
			entries, err := os.ReadDir(".")

			if err != nil {
				fmt.Println(err)
				return
			}

			var matchedRenames [][]string

			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}

				name := entry.Name()
				result, _ := regex.Replace(name, replace, -1, 1)

				if name == result {
					continue
				}

				fmt.Println(name, " -> ", result)

				matchedRenames = append(matchedRenames, []string{name, result})
			}

			if len(matchedRenames) == 0 {
				fmt.Println("No regex match!")
				return
			}

			fmt.Println("\nAccept these changes? (y/N)")
			keyboard.Open()
			char, _, _ := keyboard.GetSingleKey()
			keyboard.Close()
			if char != 'y' {
				return
			}

			for _, match := range matchedRenames {
				os.Rename(match[0], match[1])
			}
		}
	},
}

func ParseArgs(args []string) (*regexp2.Regexp, string, error) {
	regexStr := args[0]
	replace := args[1]

	// reading stdin also reads the newline, so need to remove that
	regexStr = strings.TrimRight(regexStr, "\r\n")
	regexStr = strings.Trim(regexStr, "/")
	replace = strings.TrimRight(replace, "\r\n")

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
}
