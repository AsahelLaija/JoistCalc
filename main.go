package main
import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)
func main() {
    e := echo.New()

    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    e.GET("/", func(c echo.Context) error {
	return c.File("static/index.html")
    })
    e.GET("/styles", func(c echo.Context) error {
	return c.File("static/styles.css")
    })


    e.Logger.Fatal(e.Start(":8080"))
} 
