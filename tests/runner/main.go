package main

import (
	"github.com/spf13/cobra"
	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/tests/runner/command"
)



var Runner = &cobra.Command{
	Use:   "qscamel-runner run all",
	Short: constants.ShortDescription,
	Long:  constants.LongDescription,
}


func main()  {
	Runner.AddCommand(command.RunCmd)
	Runner.Execute()
}
