package migration

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		jsonData := `{
			"createRule": "book.owner.id = @request.auth.id",
			"deleteRule": "book.owner.id = @request.auth.id",
			"fields": [
				{
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
				},
				{
					"cascadeDelete": true,
					"collectionId": "pbc_3974677633",
					"hidden": false,
					"id": "relation3420824369",
					"maxSelect": 1,
					"minSelect": 0,
					"name": "book",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "relation"
				},
				{
					"cascadeDelete": true,
					"collectionId": "pbc_1174553048",
					"hidden": false,
					"id": "relation105650625",
					"maxSelect": 1,
					"minSelect": 0,
					"name": "category",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "relation"
				},
				{
					"cascadeDelete": true,
					"collectionId": "pbc_1840306185",
					"hidden": false,
					"id": "relation2223302008",
					"maxSelect": 1,
					"minSelect": 0,
					"name": "paymentMethod",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "relation"
				},
				{
					"hidden": false,
					"id": "date2862495610",
					"max": "",
					"min": "",
					"name": "date",
					"presentable": false,
					"required": true,
					"system": false,
					"type": "date"
				},
				{
					"hidden": false,
					"id": "number2392944706",
					"max": null,
					"min": null,
					"name": "amount",
					"onlyInt": false,
					"presentable": false,
					"required": false,
					"system": false,
					"type": "number"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "text3788167225",
					"max": 0,
					"min": 0,
					"name": "remark",
					"pattern": "",
					"presentable": false,
					"primaryKey": false,
					"required": false,
					"system": false,
					"type": "text"
				},
				{
					"hidden": false,
					"id": "autodate2261412156",
					"name": "createdAt",
					"onCreate": true,
					"onUpdate": false,
					"presentable": false,
					"system": false,
					"type": "autodate"
				},
				{
					"hidden": false,
					"id": "autodate3175243278",
					"name": "updatedAt",
					"onCreate": true,
					"onUpdate": true,
					"presentable": false,
					"system": false,
					"type": "autodate"
				}
			],
			"id": "pbc_3161059916",
			"indexes": [],
			"listRule": "book.owner.id = @request.auth.id",
			"name": "expense",
			"system": false,
			"type": "base",
			"updateRule": "book.owner.id = @request.auth.id",
			"viewRule": "book.owner.id = @request.auth.id"
		}`

		collection := &core.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3161059916")
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}
