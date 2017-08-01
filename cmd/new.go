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
	"strconv"

	"github.com/jvikstedt/alarmy/edit"
	"github.com/jvikstedt/alarmy/model"
	"github.com/jvikstedt/alarmy/service"
	"github.com/spf13/cobra"
)

type Resource struct {
	Path   string
	Fields []edit.Field
	New    func() interface{}
}

var resources = map[string]Resource{
	"project": Resource{
		Path: "projects",
		Fields: []edit.Field{
			edit.Field{Name: "Name", Kind: edit.String},
		},
		New: func() interface{} { return &model.Project{} },
	},
	"job": Resource{
		Path: "jobs",
		Fields: []edit.Field{
			edit.Field{Name: "Name", Kind: edit.String},
			edit.Field{Name: "ProjectID", Kind: edit.Int},
			edit.Field{Name: "Spec", Kind: edit.String},
			edit.Field{Name: "Cmd", Kind: edit.String},
			edit.Field{Name: "Active", Kind: edit.Bool},
			edit.Field{Name: "Triggers", FieldEditor: triggersEditor},
		},
		New: func() interface{} { return &model.Job{} },
	},
}

var editor edit.Editor

func triggersEditor(object interface{}, field edit.Field) error {
	job, ok := object.(*model.Job)
	if !ok {
		return fmt.Errorf("Not a *model.Job object")
	}

	triggers := []model.Trigger{}

	for {
		fmt.Printf("Create a new trigger? (y/n): ")
		line, err := editor.GetLine()
		if err != nil {
			return err
		}

		if line == "y" {
		} else if line == "n" {
			break
		} else {
			continue
		}

		trigger := model.Trigger{}

		fmt.Printf("FieldName: ")
		line, err = editor.GetLine()
		if err != nil {
			return err
		}
		trigger.FieldName = line

		fmt.Printf("Target: ")
		line, err = editor.GetLine()
		if err != nil {
			return err
		}
		trigger.Target = line

	TriggerType:
		for i, v := range model.TriggerTypes {
			fmt.Printf("%d: %s\n", i, v)
		}

		fmt.Printf("TriggerType: ")
		line, err = editor.GetLine()
		if err != nil {
			return err
		}

		selected, err := strconv.Atoi(line)
		if err != nil {
			fmt.Println(err)
			goto TriggerType
		}

		if selected < 0 || selected > len(model.TriggerTypes)-1 {
			fmt.Printf("Select between %d and %d\n", 0, len(model.TriggerTypes)-1)
			goto TriggerType
		}

		trigger.TriggerType = model.TriggerType(selected)

		triggers = append(triggers, trigger)
	}
	job.Triggers = triggers

	return nil
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

		if err := runNewCmd(key); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func runNewCmd(resourceKey string) error {
	editor = edit.NewEditor(os.Stdin)

	resource, err := resourceByKey(resourceKey)
	if err != nil {
		return err
	}

	object := resource.New()

	fmt.Printf("New resource %s\n", resourceKey)

	if err := editor.Edit(object, resource.Fields); err != nil {
		return err
	}

	if err := service.PostAsJSON(resource.Path, object); err != nil {
		return err
	}

	pretty, err := edit.ObjectPrettyFormat(object)
	if err != nil {
		return err
	}
	fmt.Println(pretty)

	return nil
}

func resourceByKey(key string) (Resource, error) {
	resource, ok := resources[key]

	if !ok {
		return Resource{}, fmt.Errorf("object %s not found", key)
	}

	return resource, nil
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
