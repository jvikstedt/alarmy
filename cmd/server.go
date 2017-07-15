package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jvikstedt/alarmy/api"
	"github.com/jvikstedt/alarmy/schedule"
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
		port := os.Getenv("ALARMY_PORT")
		if port == "" {
			port = "8080"
		}

		err := setupServer(":" + port)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func setupServer(addr string) error {
	// Store / Database Setup
	boltStore, err := store.NewBoltStore("alarmy_dev.db")
	if err != nil {
		return err
	}
	defer boltStore.Close()

	// Logger setup
	f, err := os.OpenFile("dev.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	logger := log.New(f, "", log.LstdFlags)

	// Scheduler setup
	scheduler := schedule.NewCronScheduler(logger)
	go scheduler.Start()
	defer scheduler.Stop()

	// Server & http.Handler setup
	api := api.NewApi(boltStore.Store(), logger, scheduler)
	handler, err := api.Handler()
	if err != nil {
		return err
	}
	s := http.Server{Addr: addr, Handler: handler}

	// Signalling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		<-stop

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Shutdown server gracefully
		s.Shutdown(ctx)
	}()

	// Server startup
	printf(logger, "Listening on %s\n", addr)
	if err := s.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			logger.Println(err)
		} else {
			return err
		}
	}

	return nil
}

func printf(logger *log.Logger, format string, v ...interface{}) {
	if logger != nil {
		logger.Printf(format, v...)
	}
	fmt.Printf(format, v...)
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
