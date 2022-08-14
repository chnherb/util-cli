package cmd

import (
	"util-cli/usecase"

	"github.com/gookit/slog"
	"github.com/spf13/cobra"
)

func newImgBase64Command() *cobra.Command {
	comd := &cobra.Command{
		Use:   "imgbase64",
		Short: "exec image base64 relevant command",
		RunE: func(c *cobra.Command, _ []string) error {
			slog.Infof("Try \"%s [command] --help\" for more options", "util-cli imgbase64")
			return nil
		},
	}
	comd.AddCommand(newImgBase64ListCommand())
	return comd
}

func newImgBase64ListCommand() *cobra.Command {
	args := usecase.ImgBase64Args{}
	comd := &cobra.Command{
		Use:   "replace",
		Short: "exec img base64 relevant command",
		RunE: func(c *cobra.Command, _ []string) error {
			slog.Debugf("imgbase64 list origin args: %+v", args)
			// usecase.ParseImgBase64File("/Users/bo/hb_blog/00-util/git/git常用操作.md", true)
			return usecase.Replace(&args)
		},
	}
	comd.Flags().StringVar(&args.Src, "src", "", "source directory, default is `${pwd}")
	comd.Flags().StringVar(&args.Chapter, "chapter", "", "chapter as img's prefix")
	comd.Flags().BoolVar(&args.Rewrite, "rewrite", true, "rewrite origin xx.md")
	return comd
}
