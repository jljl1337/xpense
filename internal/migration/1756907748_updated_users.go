package migration

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"indexes": [
				"CREATE UNIQUE INDEX ` + "`" + `idx_tokenKey__pb_users_auth_` + "`" + ` ON ` + "`" + `user` + "`" + ` (` + "`" + `tokenKey` + "`" + `)",
				"CREATE UNIQUE INDEX ` + "`" + `idx_email__pb_users_auth_` + "`" + ` ON ` + "`" + `user` + "`" + ` (` + "`" + `email` + "`" + `) WHERE ` + "`" + `email` + "`" + ` != ''"
			],
			"name": "user",
			"oauth2": {
				"mappedFields": {
					"avatarURL": "",
					"name": ""
				}
			}
		}`), &collection); err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("text1579384326")

		// remove field
		collection.Fields.RemoveById("file376926767")

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(0, []byte(`{
			"autogeneratePattern": "[A-Za-z0-9_-]{10}",
			"hidden": false,
			"id": "text3208210256",
			"max": 10,
			"min": 10,
			"name": "id",
			"pattern": "[A-Za-z0-9_-]{10}",
			"presentable": false,
			"primaryKey": true,
			"required": true,
			"system": true,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(6, []byte(`{
			"hidden": false,
			"id": "autodate2990389176",
			"name": "createdAt",
			"onCreate": true,
			"onUpdate": false,
			"presentable": false,
			"system": false,
			"type": "autodate"
		}`)); err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(7, []byte(`{
			"hidden": false,
			"id": "autodate3332085495",
			"name": "updatedAt",
			"onCreate": true,
			"onUpdate": true,
			"presentable": false,
			"system": false,
			"type": "autodate"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"indexes": [
				"CREATE UNIQUE INDEX ` + "`" + `idx_tokenKey__pb_users_auth_` + "`" + ` ON ` + "`" + `users` + "`" + ` (` + "`" + `tokenKey` + "`" + `)",
				"CREATE UNIQUE INDEX ` + "`" + `idx_email__pb_users_auth_` + "`" + ` ON ` + "`" + `users` + "`" + ` (` + "`" + `email` + "`" + `) WHERE ` + "`" + `email` + "`" + ` != ''"
			],
			"name": "users",
			"oauth2": {
				"mappedFields": {
					"avatarURL": "avatar",
					"name": "name"
				}
			}
		}`), &collection); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(6, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text1579384326",
			"max": 255,
			"min": 0,
			"name": "name",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(7, []byte(`{
			"hidden": false,
			"id": "file376926767",
			"maxSelect": 1,
			"maxSize": 0,
			"mimeTypes": [
				"image/jpeg",
				"image/png",
				"image/svg+xml",
				"image/gif",
				"image/webp"
			],
			"name": "avatar",
			"presentable": false,
			"protected": false,
			"required": false,
			"system": false,
			"thumbs": null,
			"type": "file"
		}`)); err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(0, []byte(`{
			"autogeneratePattern": "[a-z0-9]{15}",
			"hidden": false,
			"id": "text3208210256",
			"max": 15,
			"min": 15,
			"name": "id",
			"pattern": "^[a-z0-9]+$",
			"presentable": false,
			"primaryKey": true,
			"required": true,
			"system": true,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(8, []byte(`{
			"hidden": false,
			"id": "autodate2990389176",
			"name": "created",
			"onCreate": true,
			"onUpdate": false,
			"presentable": false,
			"system": false,
			"type": "autodate"
		}`)); err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(9, []byte(`{
			"hidden": false,
			"id": "autodate3332085495",
			"name": "updated",
			"onCreate": true,
			"onUpdate": true,
			"presentable": false,
			"system": false,
			"type": "autodate"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
