// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/jvikstedt/alarmy/model"
	"github.com/jvikstedt/alarmy/service"
	"github.com/jvikstedt/alarmy/transform"
	"github.com/spf13/cobra"
)

var objects = map[string]interface{}{
	"project": &model.Project{},
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]

		v, ok := objects[key]

		if !ok {
			fmt.Printf("object %s not found", key)
			return
		}

		err := transform.Modify(v)
		if err != nil {
			panic(err)
		}

		pjson, err := json.MarshalIndent(v, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(pjson))

		fmt.Println("Saving...")

		err = service.ProjectNew(v)
		if err != nil {
			panic(err)
		}

		pjson, err = json.MarshalIndent(v, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(pjson))
	},
}

func init() {
	RootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
