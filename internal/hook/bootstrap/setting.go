package bootstrap

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func UpdateAppSettings(app *pocketbase.PocketBase) {
	app.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		if err := e.Next(); err != nil {
			return err
		}

		settings := e.App.Settings()

		settings.Meta.AppName = "Xpense"

		return e.App.Save(settings)
	})
}
