package cmd

import (
	"fmt"
	"strings"
	"sync"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
	"github.com/fatih/color"
)

var (
	Values string
	boldColor, successColor, errorColor *color.Color
	mutex *sync.Mutex
	bold = color.New(color.Bold).SprintFunc()
	longDesc = `Example: 
$ with -v "$(cat myurls.txt)" "curl -L this"

With is a CLI to execute a single shell command with
multiple variables concurrently. xargs is a similar utility.
Place the 'this' keyword anywhere in your command and
with will iterate through all the input values, replace
"this" with a value and execute the resulting command.

If you have multiple "this" in your command, each one
of them will be replaced (this can be controlled
with a flag in a future release).

Use cases:

   1) You have a list of URLs in a text file and
      want to curl all of them in parallel, with
      one command.
	   
      $ with -v "$(cat myurls.txt)" "curl -L this"

   2) You want to quickly check AWS instance status
      for all the regions.
	  
      $ with -v "$(cat myregions.txt)" "aws --region=this ec2 describe-instance-status"
	  
   3) You have a directory with a lot of kubeconfig
      files and want to get pods from all the different
      clusters.

      $ with -v "$(ls)" "kubectl --kubeconfig=this get pods"
`

	rootCmd = &cobra.Command{
		Use: "with",
		Short: "with let's you run any shell command with variables, in parallel",
		Long: longDesc,
		Version: "v0.1-beta.0",
		Args: cobra.ExactArgs(1),
		Run: withCmd,
	}
)

func init() {
	rootCmd.Flags().StringVarP(&Values, "values", "v", "", "Iterate with these values")
	rootCmd.MarkFlagRequired("values")	
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor = color.New(color.FgRed, color.Bold)
	mutex = &sync.Mutex{}
}

func withCmd(cmd *cobra.Command, arguments []string) {
	sep := "------------------------------------------------------------"
	failures := 0
	start := time.Now()
	with := strings.Fields(Values)
	keyLoc := strings.Index(arguments[0], "this")

	if (keyLoc == -1) {
		errorColor.Println("this not found in command.")
	}

	var wg sync.WaitGroup

	for _, v := range with {
		wg.Add(1)
		toRun := strings.Replace(arguments[0], "this", v, -1)
		go func(r string) {
			prog, args, res, err := execute(toRun)
			mutex.Lock()
			if err != nil {
				failures++
				errorColor.Println("$", prog, args)
				errorColor.Println(err)
				fmt.Printf("%s\n\n", res)

			} else {
				successColor.Println("$", prog, args)
				fmt.Printf("%s\n\n", res)
			}
			mutex.Unlock()
			wg.Done()

		}(toRun)
	}

	wg.Wait()

	t := time.Now()
	elapsed := t.Sub(start)
	command := strings.Replace(arguments[0], "this", "____", -1)

	fmt.Println(sep)
	fmt.Println(bold(" Command:"), command)
	fmt.Println(bold(" Iterations:"), len(with))
	fmt.Println(bold(" Failures:"), failures)
	fmt.Println(bold(" Duration:"), elapsed)
	fmt.Println(sep)
}

func execute(toRun string) (prog string, args string, stdoutstderr []byte, err error) {
	cmdArr := strings.Fields(toRun)
	prog = cmdArr[0]
	args = strings.Join(cmdArr[1:], " ")
	cmd := exec.Command(prog, cmdArr[1:]...)
	stdoutstderr, err = cmd.CombinedOutput()
	return 
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
