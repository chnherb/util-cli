package cmd

import (
	"github.com/gookit/slog"
	"github.com/spf13/cobra"
	"util-cli/usecase"
)

func newCollapseCommand() *cobra.Command {
	args := usecase.CollapseArgs{}
	comd := &cobra.Command{
		Use:   "collapse",
		Short: "exec collapse code command",
		RunE: func(c *cobra.Command, _ []string) error {
			slog.Infof("Try \"%s [command] --help\" for more options", "util-cli code")
			return usecase.BatchCollapseCode(&args)
		},
	}
	comd.Flags().StringVar(&args.Src, "src", "", "source directory, default is `${pwd}")
	comd.Flags().IntVar(&args.MinRow, "minrow", 20, "collapse code when row(code) >= minrow")
	comd.Flags().BoolVar(&args.NeedTitle, "needtitle", false, "whether need the collapse title, defalt is 'Expand/Collapse Code Block'")
	return comd
}
