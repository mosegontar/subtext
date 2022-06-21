/*
Copyright © 2022 Alex Gontar mosegontar@gmail.com

*/
package cmd

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mosegontar/stegapix/lib"
	"github.com/spf13/cobra"
)

// encodeCmd represents the encode command
var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		imageFilename := args[0]

		var messageBytes []byte

		if len(args) == 2 {
			messageBytes = parseMessage(args[1])
		} else {
			messageBytes = parseMessage("")
		}

		img := lib.EncodeByteToPixel(imageFilename, messageBytes)
		lib.SaveImage(img, outputFilename(imageFilename))
	},
}

func splitOnExtension(filename string) (string, string) {
	ext := path.Ext(filename)
	basename := strings.TrimSuffix(filename, ext)
	return basename, ext
}

func outputFilename(originalFilename string) string {
	filename, ext := splitOnExtension(originalFilename)

	currentTime := time.Now()
	formattedTime := currentTime.Format("20060102150405")

	parts := []string{filename, "_", formattedTime, ext}

	return strings.Join(parts, "")

}

func parseMessage(message string) []byte {
	if message != "" {
		return []byte(message + "\000")
	}

	byts, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err.Error())
	}
	byts = append(byts, byte('\000'))
	return []byte(string(byts[:]))
}

func init() {
	rootCmd.AddCommand(encodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
