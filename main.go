package main
import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "strconv"
    "net/http"
    "fmt"
    "html/template"
    "io"
    // "log"
)

// templates Code
type Template struct {
    templates *template.Template
}
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

// logic Code
type Geometry struct {
    trussType	string
    joistType	string
    deflexion	float64
    span	float64
    fepl	float64
    sepl	float64
    ipl		float64
    lbe		float64
    depth	float64
}

func (g *Geometry) dataGeometry(tt, jt, df, sp, fe, se, ip, lb, de string) {
    g.trussType		= tt
    g.joistType		= jt
    g.deflexion, _	= strconv.ParseFloat(df, 64)
    g.span, _		= strconv.ParseFloat(sp, 64)
    g.fepl, _		= strconv.ParseFloat(fe, 64)
    g.sepl, _		= strconv.ParseFloat(se, 64)
    g.ipl, _		= strconv.ParseFloat(ip, 64)
    g.lbe, _ 		= strconv.ParseFloat(lb, 64)
    g.depth, _		= strconv.ParseFloat(de, 64)
}

// Films Struct

type Film struct {
    Title	string
    Director	string
}
func bChordtobPanel(fepl, sepl, ipl, lbe float64) float64 {
    return (fepl + sepl + ipl) - lbe
}

func designLength(span float64) float64 {
    return span*12 - 4
}
/*
func Hello(c echo.Context) error {
    films := map[string][]Film {
	"Films" : {
	    {Title: "The Godfather", Director: "Francis Ford Copola"},
	    {Title: "Blade Runner", Director: "Ridley Scott"},
	    {Title: "The Thing", Director: "John Carpenter"},
	},
    }
    return c.Render(http.StatusOK, "hello", films)
}
*/
func main() {
    t := &Template{
	templates: template.Must(template.ParseGlob("views/*.html")),
    }
    e := echo.New()
    e.Renderer = t

    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    /*		Home		*/
    e.GET("/", func(c echo.Context) error {
	return c.Render(http.StatusOK, "Home", "Joist Calculator")
	// return c.File("views/index.html")
    })

    e.GET("/hello", func (c echo.Context) error {
	films := map[string][]Film {
	    "Films" : {
		{Title: "The Godfather", Director: "Francis Ford Copola"},
		{Title: "Blade Runner", Director: "Ridley Scott"},
		{Title: "The Thing", Director: "John Carpenter"},
	    },
	}
	return c.Render(http.StatusOK, "hello", films)
    })

    /*		Styles		*/
    e.GET("/styles", func(c echo.Context) error {
	return c.File("static/styles.css")
    })

    var newGeometry Geometry
    
    e.POST("/geometry", func(c echo.Context) error {

	trussType := c.FormValue("trussType")
	joistType := c.FormValue("joistType")
	deflexion := c.FormValue("deflexion")
	span := c.FormValue("span")
	fepl := c.FormValue("fepl")
	sepl := c.FormValue("sepl")
	ipl := c.FormValue("ipl")
	lbe := c.FormValue("lbe")
	depth := c.FormValue("depth")

	newGeometry.dataGeometry(
	    trussType,
	    joistType,
	    deflexion,
	    span,
	    fepl,
	    sepl,
	    ipl,
	    lbe,
	    depth,
	)

	lbe2 := bChordtobPanel(
	    newGeometry.fepl,
	    newGeometry.sepl,
	    newGeometry.ipl,
	    newGeometry.lbe,
	)

	dLength := designLength(newGeometry.span)

	s := strconv.FormatFloat(lbe2, 'g', -1, 64)
	d := strconv.FormatFloat(dLength, 'g', -1, 64)

	var result string
	result = s + "\n" + d

	fmt.Println(newGeometry)

	return c.String(http.StatusOK, result)
    })
    e.POST("/add-film", func (c echo.Context) error {
	title := c.FormValue("title")
	director := c.FormValue("director")
	htmlStr := fmt.Sprintf("<li class='list-group-item bg-primary text-white'>%s - %s</li>", title, director)
	fmt.Print("\n\n", title, "\n", director, "\n")
	return c.String(http.StatusOK, htmlStr)
    })

    // test of routing and htmx
    // e.POST("/save", saveUser)
    // e.GET("/users/:id", getUser)
    // e.PUT("/users/:id", updateUser)
    // e.DELETE("/users/:id", deleteUser)


    e.Logger.Fatal(e.Start(":3000"))
}
