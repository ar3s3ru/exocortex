// Copyright © 2017 Spencer Dixon <spencercdixon@gmail.com>

package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/apex/httplog"
	"github.com/apex/log"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spencercdixon/exocortex/config"
	"github.com/spencercdixon/exocortex/wiki"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(startCmd)
}

// startCmd boots up the wiki server
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Boot up the wiki for a specific repository",
	Run: func(cmd *cobra.Command, args []string) {
		// Boot up config
		var configPath string
		if len(args) < 1 {
			configPath = "."
		} else {
			configPath = args[0]
		}
		config.New(configPath)

		// Root router
		router := mux.NewRouter()

		// Context to signal sync goroutine to stop
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Mount wiki under '/api'
		wikiHandler := wiki.New(ctx)
		router.PathPrefix("/api").Handler(
			http.StripPrefix("/api", wikiHandler),
		)

		// Handle Static Assets - catch all for SPA + serving bundled JS/CSS
		box := packr.NewBox("../ui/build")
		router.PathPrefix("/static").Handler(http.FileServer(box))
		router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(box.Bytes("index.html"))
		})

		// Add some utility handlers
		var handler http.Handler
		handler = httplog.New(router)
		handler = handlers.CompressHandler(handler)
		handler = handlers.RecoveryHandler()(handler)

		// Boot up wiki
		address := fmt.Sprintf(
			"%s:%d",
			viper.GetString("server.host"),
			viper.GetInt("server.port"),
		)
		log.Infof("Booting up wiki on: %s", address)
		log.Fatalf("Listening: %s", http.ListenAndServe(address, handler))
	},
}
