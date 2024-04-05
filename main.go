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
	return c.File("views/index.html")
    })

    e.GET("/styles", func(c echo.Context) error {
	return c.File("static/styles.css")
    })

    e.POST("/dataEnter" func(c echo.Context) error {
	name := c.FormValue("name")
	return c.String(http.StatusOK, "name: " + name)
    })

    e.Logger.Fatal(e.Start(":3000"))
} 
