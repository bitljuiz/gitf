package executor

import (
	"fmt"
	"gitf/configs"
	"gitf/internal/cli"
	"gitf/internal/flags"
	"gitf/internal/git"
	"gitf/internal/globs"
	"gitf/internal/statistics"
	"github.com/spf13/cobra"
	"os"
)

func Run() {
	var gitfCmd = cobra.Command{
		Use:   "gitf",
		Short: "Gitf is an utility for collecting detailed information about commits",
		Long: "Fast utility that can collect statistics about\n1. Lines count\n2. Commits count\n3. Files count\n" +
			"in git repository.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 && args[0] == "help" {
				_ = cmd.Help()
				os.Exit(1)
			}
		},
	}

	flagsMap := flags.Create()

	for _, flag := range flagsMap {
		switch flagValue := flag.Value.(type) {
		case *flags.BoolValue:
			gitfCmd.Flags().BoolVarP(
				flagValue.Pointer(),
				flag.Name,
				flag.ShorthandName,
				flagValue.DefaultValue(),
				flag.Use,
			)
		case *flags.StringValue:
			gitfCmd.Flags().StringVarP(
				flagValue.Pointer(),
				flag.Name,
				flag.ShorthandName,
				flagValue.DefaultValue(),
				flag.Use,
			)
		case *flags.StringSliceValue:
			gitfCmd.Flags().StringSliceVarP(
				flagValue.Pointer(),
				flag.Name,
				flag.ShorthandName,
				flagValue.DefaultValue(),
				flag.Use,
			)
		default:
			panic("unhandled default case")
		}
	}

	if err := gitfCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if helpCmd, _ := gitfCmd.Flags().GetBool("help"); helpCmd {
		os.Exit(1)
	}

	repository, _ := flagsMap[flags.Repository].GetString()
	revision, _ := flagsMap[flags.Revision].GetString()

	files, err := git.GetRepositoryFiles(repository, revision)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	exclude, _ := flagsMap[flags.Exclude].GetStringSlice()

	files2 := globs.Exclude(exclude, files)

	restrictTo, _ := flagsMap[flags.RestrictTo].GetStringSlice()

	files3 := globs.RestrictTo(restrictTo, files2)

	languages, _ := flagsMap[flags.Languages].GetStringSlice()
	extensions, _ := flagsMap[flags.Extensions].GetStringSlice()
	onlyProgramming, _ := flagsMap[flags.OnlyProgramming].GetBool()

	requiredFiles, err := configs.LimitExtensions(files3, languages, extensions, onlyProgramming)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	useCommitter, _ := flagsMap[flags.UseCommitter].GetBool()

	var stats map[string]*statistics.Author
	var errList []error

	if stats, errList = git.GetStatistics(repository, revision, useCommitter, requiredFiles); len(errList) != 0 {
		for _, err = range errList {
			fmt.Println(err)
		}
	}

	if len(stats) != 0 {
		order, _ := flagsMap[flags.OrderBy].GetString()
		format, _ := flagsMap[flags.Format].GetString()

		sortedStats := statistics.Sort(order, stats)
		err = cli.FormatData(sortedStats, format)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
