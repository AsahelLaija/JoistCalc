package main
import (
    "fmt"
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func getUser (c echo.Context) error {
    // User ID from path `users/:id`
    id := c.Param("id")
    return c.String(http.StatusOK, id)
}

func saveUser (c echo.Context) error {
    // Get name and email
    name := c.FormValue("name")
    email := c.FormValue("email")
    return c.String(http.StatusOK, "name: " + name + "\nmail: " + email)
}

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
    
    e.POST("/geometry", func(c echo.Context) error {
	trussType := c.FormValue("trussType")
	joistType := c.FormValue("joistType")
	deflexion := c.FormValue("deflexion")
	span := c.FormValue("span")
	fepl := c.FormValue("fepl")
	sepl := c.FormValue("sepl")
	depth := c.FormValue("depth")

	fmt.Print("\tTruss Type:\t\t", trussType, "\n")
	fmt.Print("\tjoistType:\t\t", joistType, "\n")
	fmt.Print("\tdeflexion:\t\t", deflexion, "\n")
	fmt.Print("\tspan:\t\t", span, "\n")
	fmt.Print("\tfepl:\t\t", fepl, "\n")
	fmt.Print("\tsepl:\t\t", sepl, "\n")
	fmt.Print("\tdepth:\t\t", depth, "\n")
	return c.String(http.StatusOK, "whre good")
    })

    e.POST("/contacts", func(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	// return c.String(http.StatusOK, "name: " + name + "\nemail: " + email )
	return c.String(http.StatusOK, "<tr> <td>" + name + "</td>" + "<td>" + email + "</td> </tr>" )
    })

    // test of routing and htmx
    e.POST("/save", saveUser)
    // e.GET("/users/:id", getUser)
    // e.PUT("/users/:id", updateUser)
    // e.DELETE("/users/:id", deleteUser)


    e.Logger.Fatal(e.Start(":3000"))
}
