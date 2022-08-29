package cmd

import (
	"os"
	"util-cli/utils"

	"github.com/gookit/color"
	"github.com/gookit/slog"
	"github.com/spf13/cobra"
)

func Init(verbose bool) {
	utils.InitSlog(verbose)
}

func newRootCmd() *cobra.Command {
	verbose := false
	cmd := &cobra.Command{
		Use:           "util-cli",
		Short:         "util-cli, more see: xxx",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			Init(verbose)
		},
	}
	cmd.SetHelpFunc(func(c *cobra.Command, _ []string) {
		color.Info.Print(c.UsageString())
	})
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "show more detail message")
	return cmd
}

func newCommand() *cobra.Command {
	rootCmd := newRootCmd()
	rootCmd.AddCommand(newImgBase64Command())
	rootCmd.AddCommand(newCollapseCommand())
	return rootCmd
}

func Execute() {
	rootCmd := newCommand()
	if err := rootCmd.Execute(); err != nil {
		// Manual handle error
		slog.Error(err)
		color.Magenta.Printf("Try \"%s [command] --help\" for more options\n", rootCmd.Use)
		os.Exit(1)
	}
}
