package cmd

import (
	"net/http"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/messagebird/internal"
	"github.com/messagebird/internal/config"
)

// var SMSChan = make(chan *internal.Messages, 10)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts messagebird server",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		configFile, err := filepath.Abs(cmd.Flag("config").Value.String())
		if err != nil {
			log.Fatal("can not get full path to config file:", err)
		}
		cnf, err := config.MakeServerConfigFromFile(configFile)

		// Set logging level
		logLevel, err := log.ParseLevel(cnf.LogLevel)
		if err != nil {
			log.Fatal("error parsing log level: %v", err)
		}
		log.SetLevel(logLevel)
		log.Infof("starting server at %s", cnf.ServerAddr)

		q := make(chan *internal.Messages, 100)
		s := internal.NewServer(cnf,&q)
		// Start SMS worker
		go func(q <-chan *internal.Messages) {
			s.MessagebirdWorker(q)
		}(q)

		http.HandleFunc("/", s.Handler )

		//@TODO- gracefull stop
		// Run server
		log.Println("Listen and service on", cnf.ServerAddr)
		err = http.ListenAndServe(cnf.ServerAddr, nil)
		if err != nil {
			log.Fatal("Could not start server: ", err.Error())
		}
	},
}

func init() {
	rootCmd.PersistentFlags().String("config", "configs/config.default.json", "Config file to load")
	rootCmd.AddCommand(serverCmd)
}