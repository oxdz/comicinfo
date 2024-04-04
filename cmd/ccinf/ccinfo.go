package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := newRootCmd()
	rootCmd.AddCommand(newShowCmd())
	rootCmd.AddCommand(newStartCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ccinf",
		Short: `ccinf is a tool used to obtain comic information`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return nil
		},
	}
	return cmd
}

func newShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "show supported sites",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprint(os.Stdout, "supported sites:\n")
			return nil
		},
	}
	return cmd
}

func newStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start to get comic info",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprint(os.Stdout, "hello\n")
			return nil
		},
	}
	return cmd
}
