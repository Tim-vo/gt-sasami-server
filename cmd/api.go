package cmd

import (
	"net/http"

	cli "github.com/spf13/cobra"
	"go.uber.org/zap"

	conf "github.com/Tim-vo/gt-sasami-server/config"
	"github.com/Tim-vo/gt-sasami-server/embed"
	"github.com/Tim-vo/gt-sasami-server/gt-sasami-server/rpc"
	"github.com/Tim-vo/gt-sasami-server/server"
	"github.com/Tim-vo/gt-sasami-server/store/postgres"
)

func init() {
	rootCmd.AddCommand(apiCmd)
}

var (
	apiCmd = &cli.Command{
		Use:   "api",
		Short: "Start API",
		Long:  `Start API`,
		Run: func(cmd *cli.Command, args []string) { // Initialize the databse

			migrationSource, err := embed.MigrationSource()
			if err != nil {
				logger.Fatalw("could not get database migrations", "error", err)
			}

			// Database
			pg, err := postgres.New(conf.C, migrationSource)
			if err != nil {
				logger.Fatalw("database error", "error", err)
			}

			// Create the server
			s, err := server.New(conf.C)
			if err != nil {
				logger.Fatalw("could not create server", "error", err)
			}

			s.Router().Get("/version", conf.GetVersion())

			// ThingRPC
			if err = rpc.Setup(s.Router(), pg); err != nil {
				logger.Fatalw("could not setup thingrpc", "error", err)
			}

			// Serve apidocs and swagger-ui
			docsFileServer := http.FileServer(http.FS(embed.PublicHTMLFS()))
			s.Router().Mount("/apidocs", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("vary", "Accept-Encoding")
				w.Header().Set("cache-Control", "no-cache")
				docsFileServer.ServeHTTP(w, r)
			}))

			if err = s.ListenAndServe(conf.C); err != nil {
				logger.Fatalw("could not start server", "error", err)
			}

			conf.Stop.InitInterrupt()
			<-conf.Stop.Chan() // Wait until Stop
			conf.Stop.Wait()   // Wait until everyone cleans up
			_ = zap.L().Sync() // Flush the logger

		},
	}
)
