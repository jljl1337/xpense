package bootstrap

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"

	"github.com/jljl1337/xpense/internal/env"
)

func UpsertSuperuser(app *pocketbase.PocketBase) {
	app.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		if err := e.Next(); err != nil {
			return err
		}

		email := env.MustGetString("SUPERUSER_EMAIL", "admin@example.com")
		password := env.MustGetString("SUPERUSER_PASSWORD", "admin12345")

		totalSuperusers, err := app.CountRecords(core.CollectionNameSuperusers)
		if err != nil {
			return err
		}

		// skip if there are multiple superusers
		if totalSuperusers > 1 {
			app.Logger().Warn("Not updating superuser credentials since there are multiple superusers")
			return nil
		}

		// create a default superuser if none exists
		if totalSuperusers == 0 {
			app.Logger().Info("Creating a default superuser")

			superusers, err := e.App.FindCollectionByNameOrId(core.CollectionNameSuperusers)
			if err != nil {
				return err
			}

			record := core.NewRecord(superusers)
			record.Set("email", email)
			record.Set("password", password)

			return e.App.Save(record)
		}

		// update the existing superuser if necessary
		record, err := e.App.FindFirstRecordByFilter(core.CollectionNameSuperusers, "")
		if err != nil {
			return err
		}

		if record.Get("email") != email || !record.ValidatePassword(password) {
			app.Logger().Info("Updating the existing superuser")
			record.Set("email", email)
			record.Set("password", password)
		} else {
			app.Logger().Info("No changes detected for the existing superuser")
		}

		return e.App.Save(record)
	})
}
