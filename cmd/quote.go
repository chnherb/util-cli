package cmd

import (
	"github.com/gookit/slog"
	"github.com/spf13/cobra"
	"util-cli/usecase"
)

func newQuoteCommand() *cobra.Command {
	args := usecase.QuoteArgs{}
	comd := &cobra.Command{
		Use:   "quote",
		Short: "exec quote command",
		RunE: func(c *cobra.Command, _ []string) error {
			slog.Infof("Try \"%s [command] --help\" for more options", "util-cli quote")
			return usecase.BatchQuote(&args)
		},
	}
	comd.Flags().StringVar(&args.Src, "src", "", "source directory, default is `${pwd}")
	return comd
}
