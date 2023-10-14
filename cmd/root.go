/*
Copyright © 2023 Chihiro Hasegawa <encry1024@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/owlinux1000/bincookie/internal"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:        "bincookie",
	Short:      "A simple tool to parse \"Cookies.binarycookies\"",
	Long:       ``,
	Args:       cobra.MinimumNArgs(1),
	ArgAliases: []string{"path"},
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		format, err := cmd.Flags().GetString("format")
		if err != nil {
			log.Fatal(err)
		}
		c, err := internal.NewBinaryCookie(path)
		if err != nil {
			log.Fatal(err)
		}
		if !c.IsBinaryCookie() {
			log.Printf("`%s` is not binarycookie format\n", path)
			os.Exit(0)
		}
		numberOfPages := c.ReadNumberOfPages()
		eachPageSize := c.ReadEachPageSize(numberOfPages)
		fmt.Printf(c.String(eachPageSize, format))
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("format", "f", "curl", "The output format [json|csv|curl]")
}
