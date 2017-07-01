package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jvikstedt/alarmy/api"
	"github.com/jvikstedt/alarmy/store"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		boltStore, err := store.NewBoltStore("alarmy_dev.db")
		if err != nil {
			panic(err)
		}
		defer boltStore.Close()

		store := store.Store{
			ProjectStore: store.NewBoltProjectStore(boltStore),
		}

		api := api.NewApi(store)
		handler, err := api.Handler()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := http.ListenAndServe(":8080", handler); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
