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
func bChordtobPanel(fepl, sepl, ipl, lbe float64) float64 {
    return (fepl + sepl + ipl) - lbe
}

func designLength(span float64) float64 {
    return span*12 - 4
}

// Films Struct

type Film struct {
    Title	string
    Director	string
}
/* 		page		*/
type Contact struct {
    Name	string
    Email	string
}
type Propertie struct {
    TrussType,
    JoistType	string
    deflexion,
    span,
    fepl,
    sepl,
    ipl,
    lbe,
    depth	float64
}
type ResProp struct {
    Lbe2, DLength, Tip, Tod, Ts, Ed string
}

type Contacts = []Contact
type Properties = []Propertie
type ResProps = []ResProp

type Data struct {
    Contacts Contacts
}
type Geometry struct{
    Properties Properties
}
type ResGeom struct {
    ResProps ResProps
}

func newContact(name, email string) Contact {
    return Contact {
	Name: name,
	Email: email,
    }
}
func newPropertie(tt, jt, df, sp, fe, se, ip, lb, de string) Propertie {
    deflexion,_ := strconv.ParseFloat(df, 64)
    span,_ := strconv.ParseFloat(sp, 64)
    fepl,_ := strconv.ParseFloat(fe, 64)
    sepl,_ := strconv.ParseFloat(se, 64)
    ipl,_ := strconv.ParseFloat(ip, 64)
    lbe,_ := strconv.ParseFloat(lb, 64)
    depth,_ := strconv.ParseFloat(de, 64)

    return Propertie {
	TrussType:	tt,
	JoistType:	jt,
	deflexion:	deflexion,
	span:		span,
	fepl:		fepl,
	sepl:		sepl,
	ipl:		ipl,
	lbe:		lbe,
	depth:		depth,
    }
} 
func newResProp(lb, dl, ti, to, ts, ed string) ResProp {
    return ResProp {
	Lbe2:		lb,
	DLength:	dl,
	Tip:		ti,
	Tod:		to,
	Ts:		ts,
	Ed:		ed,
    }
}

type Page struct {
    Data Data
    Geometry Geometry
    ResGeom ResGeom
}
func newData() Data {
    return Data{
	Contacts: []Contact{
	    newContact("Joist", " Calculation"),
	},
    }
}
func newGeometry() Geometry {
    return Geometry{
	Properties: []Propertie{
	     newPropertie("warren", " roof", "240", "49.21", "27.28", "26", "24", "41.1", "28"),
	},
    }
}
func newResGeom() ResGeom {
    return ResGeom{
	ResProps: []ResProp{
	    //newResProp("36.16", "586.52", "20", "24", "10", "27.01456693"),
	    newResProp("", "", "", "", "", ""),
	},
    }
}
func newPage() Page {
    return Page {
	Data: newData(),
	Geometry: newGeometry(),
	ResGeom: newResGeom(),
    }
}

/*tmp		send Data back		*/

func main() {
    t := &Template{
	templates: template.Must(template.ParseGlob("views/*.html")),
    }
    e := echo.New()
    page := newPage()
    fmt.Println(page.Geometry)

    e.Renderer = t

    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    /*		Home		*/
    e.GET("/", func(c echo.Context) error {
	return c.Render(http.StatusOK, "Home", page)
	/* TOOD: 
		Refresh response data to zero
	*/
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

	/* cannot use deflexion (variable of type string) as float64 value in assignment
	page.Geometry.Properties[0].TrussType = trussType
	page.Geometry.Properties[0].JoistType = JoistType
	page.Geometry.Properties[0].deflexion = deflexion
	page.Geometry.Properties[0].span = span
	page.Geometry.Properties[0].fepl = fepl
	page.Geometry.Properties[0].sepl = sepl
	page.Geometry.Properties[0].ipl = ipl
	page.Geometry.Properties[0].lbe = lbe
	page.Geometry.Properties[0].depth = depth
	*/

	propertie := newPropertie(
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

	page.Geometry.Properties = append(page.Geometry.Properties, propertie)
	page.Geometry.Properties = page.Geometry.Properties[1:]

	dLength := designLength(page.Geometry.Properties[0].span)
	lbe2 := bChordtobPanel(
	    page.Geometry.Properties[0].fepl,
	    page.Geometry.Properties[0].sepl,
	    page.Geometry.Properties[0].ipl,
	    page.Geometry.Properties[0].lbe,
	)
	/*
	tip := totalInteriorPanel(
	    dLength,
	    page.Geometry.Properties[0].fepl,
	    page.Geometry.Properties[0].sepl,
	)
	*/
 
	s := strconv.FormatFloat(lbe2, 'g', -1, 64)
	d := strconv.FormatFloat(dLength, 'g', -1, 64)

	page.ResGeom.ResProps[0].Lbe2 = s
	page.ResGeom.ResProps[0].DLength = d

	// res := s + " " + d

	fmt.Print("\n\n\n", page.Geometry, "\n\n\n")
	// fmt.Print("\n\n", page.ResGeom, "\n\n")

	return c.Render(http.StatusOK, "geometryResponse", page.ResGeom)
	// return c.String(http.StatusOK, res)
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


    e.Logger.Fatal(e.Start(":8080"))
}
