package main

import (
	"log"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"github.com/jljl1337/xpense/internal/hook/bootstrap"
	"github.com/jljl1337/xpense/internal/hook/request"
	"github.com/jljl1337/xpense/internal/hook/serve"
	_ "github.com/jljl1337/xpense/internal/migration"
)

func main() {
	app := pocketbase.New()

	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Dir: "./internal/migration",
		// enable auto creation of migration files when making collection changes in the Dashboard
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	bootstrap.UpdateAppSettings(app)
	bootstrap.UpsertSuperuser(app)

	serve.ServeSiteFiles(app)

	request.AddExpenseChecks(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
