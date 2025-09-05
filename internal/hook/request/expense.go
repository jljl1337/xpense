package request

import (
	"errors"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func AddExpenseChecks(app *pocketbase.PocketBase) {
	app.OnRecordCreateRequest("expense").BindFunc(checkSameBook)
	app.OnRecordUpdateRequest("expense").BindFunc(checkSameBook)
	app.OnRecordCreateRequest("expense").BindFunc(checkDateAtMidnight)
	app.OnRecordUpdateRequest("expense").BindFunc(checkDateAtMidnight)
}

func checkSameBook(e *core.RecordRequestEvent) error {
	category := e.Record.GetString("category")
	paymentMethod := e.Record.GetString("paymentMethod")

	record, err := e.App.FindRecordById("category", category)
	if err != nil {
		return err
	}
	categoryBook := record.GetString("book")

	record, err = e.App.FindRecordById("paymentMethod", paymentMethod)
	if err != nil {
		return err
	}
	paymentMethodBook := record.GetString("book")

	if categoryBook != paymentMethodBook {
		return errors.New("category and payment method must belong to the same book")
	}

	return e.Next()
}

func checkDateAtMidnight(e *core.RecordRequestEvent) error {
	date := e.Record.GetDateTime("date").Time()

	if date.Hour() != 0 || date.Minute() != 0 || date.Second() != 0 {
		return errors.New("date must be at midnight")
	}

	return e.Next()
}
