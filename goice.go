package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/rest"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// add new "GET /api/hello" route to the app router (echo)
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/quote",
			Handler: func(c echo.Context) error {
				result := struct {
					Quote  string `db:"quote" json:"quote"`
					Author string `db:"author" json:"author"`
				}{}
				queryErr := app.Dao().DB().
					NewQuery("SELECT quote, author FROM quotes ORDER BY RANDOM() LIMIT 1").
					One(&result)
				if queryErr != nil {
					return rest.NewBadRequestError("Failed to fetch.", queryErr)
				}
				return c.JSON(200, result)
			},
			Name: "",
		})

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
