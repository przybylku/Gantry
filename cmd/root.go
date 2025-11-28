package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gantry",
	Short: "Gantry is a Fast and Flexible Static Site Generator",
	Long:  `Gantry is a Fast and Flexible Static Site Generator`,
}

func init(){

	rootCmd.AddCommand(deployCmd)
}


func Execute() {
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

