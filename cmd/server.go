package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jvikstedt/alarmy/internal"
	"github.com/jvikstedt/alarmy/internal/api"
	"github.com/jvikstedt/alarmy/internal/model"
	"github.com/jvikstedt/alarmy/internal/store"
	"github.com/jvikstedt/alarmy/schedule"
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
	// Find out root directory for everything app related
	rootDir, err := getRootDir()
	if err != nil {
		return err
	}

	db, err := gorm.Open("sqlite3", filepath.Join(rootDir, "alarmy.db"))
	if err != nil {
		return err
	}
	defer db.Close()

	db.AutoMigrate(
		&model.Project{},
		&model.Job{},
		&model.Trigger{},
	)

	store := store.NewGormStore(db)

	// Logger setup
	f, err := os.OpenFile(filepath.Join(rootDir, "server.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	logger := log.New(f, "", log.LstdFlags)

	// Scheduler setup
	scheduler := schedule.NewCronScheduler(logger)
	go scheduler.Start()
	defer scheduler.Stop()

	executor := internal.NewExecutor(db, logger)

	//jobs, err := sqlStore.Job().All()
	//if err != nil {
	//	return err
	//}
	//for _, j := range jobs {
	//	if j.Active {
	//		scheduler.AddEntry(schedule.EntryID(j.ID), j.Spec, executor.Execute)
	//	}
	//}

	// Server & http.Handler setup
	api := api.NewApi(store, logger, scheduler, executor)
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

func getRootDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	homeDir := usr.HomeDir

	rootDir := filepath.Join(homeDir, ".alarmy")
	os.MkdirAll(rootDir, os.ModePerm)

	return rootDir, nil
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
