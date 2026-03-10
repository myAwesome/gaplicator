package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var outputPath string

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the project",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Building... output: %s\n", outputPath)
		// TODO: add actual build logic here
		return nil
	},
}

func init() {
	buildCmd.Flags().StringVarP(&outputPath, "output", "o", "dist", "Output path for the build")
}
