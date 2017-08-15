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
	"fmt"
	"os"

	edit "github.com/jvikstedt/alarmy/editor"
	"github.com/jvikstedt/alarmy/internal/model"
	"github.com/jvikstedt/alarmy/internal/service"
	"github.com/spf13/cobra"
)

var resources = map[string]edit.Resource{
	"project": edit.Resource{
		Fields: []edit.Field{
			edit.Field{Name: "Name"},
		},
		New: func() interface{} { return &model.Project{} },
	},
	"job": edit.Resource{
		Fields: []edit.Field{
			edit.Field{Name: "Name"},
			edit.Field{Name: "ProjectID"},
			edit.Field{Name: "Spec"},
			edit.Field{Name: "Cmd"},
			edit.Field{Name: "Active"},
			edit.Field{Name: "Triggers", Association: "trigger"},
		},
		New: func() interface{} { return &model.Job{} },
	},
	"trigger": edit.Resource{
		Fields: []edit.Field{
			edit.Field{Name: "JobID"},
			edit.Field{Name: "Target"},
			edit.Field{Name: "Val"},
			edit.Field{Name: "Type"},
		},
		New: func() interface{} { return &model.Trigger{} },
	},
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
		if len(args) == 0 {
			fmt.Println("Give resource name as argument")
			os.Exit(1)
		}
		key := args[0]

		if err := runNewCmd(key); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func runNewCmd(resourceKey string) error {
	editor := edit.NewEditor(os.Stdin, os.Stdout, resources)

	object, err := editor.NewObject(resourceKey)
	if err != nil {
		return err
	}

	return service.PostAsJSON(resourceKey+"s", object)
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
