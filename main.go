package main

import (
	"os"

	"github.com/labstack/echo/v4"
)

func main() {

	app := echo.New()

	app.Logger.Fatal(
		app.Start(os.Getenv("PORT")),
	)
}
