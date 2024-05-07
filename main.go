package main
import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "strconv"
    "net/http"
    "fmt"
    "html/template"
    "io"
    "encoding/csv"
    "os"
    "math"
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
func totalInteriorPanel(dLength, fepl, sepl, ipl float64)float64 {
    return (dLength -(2*fepl) - (2*sepl))/ipl
}
func totalDiagonal (tip float64)float64 {
    return tip + 4
}
func totalStruts(tip float64) float64 {
    return tip/2
}
func spanDepth(span float64)float64 {
    return span/24
}
func doubleAngle(area, ix, y string, sbca float64) (string, string, string) {
    areaStr,_  := strconv.ParseFloat(area, 64)
    areaFlt := 2*areaStr

    areaDAngle := strconv.FormatFloat(areaFlt, 'g', -1, 64)

    ixStr,_ := strconv.ParseFloat(ix, 64)
    yStr,_ := strconv.ParseFloat(y, 64)
    
    ixFlt := 2*ixStr
    ixDAngle := strconv.FormatFloat(ixFlt, 'g', -1, 64)
    
    // res is egual to ((sbca/2)+Ytc)^2
    res := math.Pow(((sbca/2)+yStr), 2)
    iyFlt := 2*(ixStr + areaStr*res)

    iyDAngle := strconv.FormatFloat(iyFlt, 'g', -1, 64)
    return areaDAngle, ixDAngle, iyDAngle
}
func getAnglProps(noAngleTop, noAngleBot, sbca float64) AnglProp{
    /* TOOD: 
	    make function to get angle info
    */
    var rT int = int(noAngleTop)
    var rB int = int(noAngleBot)
    // Read CSV File
    file, err := os.Open("Propiedades.csv")
    if err != nil {
	fmt.Println("Error: ", err)
    }
    defer file.Close()
    reader := csv.NewReader(file)
    record, err := reader.ReadAll()
    if err != nil {
	fmt.Println("Error: ", err)
    }
    areaTop, IxTop, IyTop := doubleAngle(
	record[rT][3], 
	record[rT][11],
	record[rT][10],
	sbca,
    )
    areaBot, IxBot, IyBot := doubleAngle(
	record[rB][3], 
	record[rB][11],
	record[rB][10],
	sbca,
    )

    a := newAnglProp(
	record[rT][2],	//secctionTop
	areaTop,	//AreaTop
	record[rT][7],	//RxTop
	record[rT][8],	//RzTop
	record[rT][9],	//RyTop
	record[rT][10],	//YTop
	IxTop,		//IxTop
	IyTop,		//IyTop
	record[rT][5],	//BTop
	record[rT][4],	//TTop
	record[rT][12],	//QTop

	record[rB][2],	//secctionBot
	areaBot,	//AreaBot
	record[rB][7],	//RxBot
	record[rB][8],	//RzBot
	record[rB][9],	//RyBot
	record[rB][10],	//YBot
	IxBot,		//IxBot
	IyBot,	        //IyBot
	record[rB][5],	//BBot
	record[rB][4],	//TBot
	record[rB][12],	//QBot
)
    return a
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
    Lbe2, DLength, Tip, Tod, Ts, Ed, Dmin, Lbrdng1, Lbrdng2 string
}
type Force struct {
    YieldStress, ModElas, SpaceChord, Weight, BSeat, TopChord, BottomChord float64
}
type AnglProp struct {
    SecTop,
    SecBot string
    
    AreaTop,
    RxTop,
    RzTop,
    RyTop,
    YTop,
    IxTop,
    IyTop,
    BTop,
    TTop,
    QTop,
    AreaBot,
    RxBot,
    RzBot,
    RyBot,
    YBot,
    IxBot,
    IyBot,
    BBot,
    TBot,
    QBot float64
}
type BridgingProp struct {
    Brdgng1,
    Brdgng2 float64
}
type MemberInput struct {
    InputName string
    ElemName string
    // request
    Part string
    Mark string
    Crimped string
    // response
    Secction string
    MidPanel string
}
type Contacts = []Contact
type Properties = []Propertie
type ResProps = []ResProp
type Forces = []Force
type AnglProps = []AnglProp
type BridgingProps = []BridgingProp
type MemberInputs = []MemberInput
type Data struct {
    Contacts Contacts
}
type Geometry struct{
    Properties Properties
}
type ResGeom struct {
    ResProps ResProps
}
type Material struct {
    Forces Forces
}
type ResMater struct {
    AnglProps AnglProps
}
type ResBrid struct {
    BridgingProps BridgingProps
}
type WebMember struct {
    MemberInputs MemberInputs
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
func newResProp(lb, dl, ti, to, ts, ed, Dmin string) ResProp {
    return ResProp {
	Lbe2:		lb,
	DLength:	dl,
	Tip:		ti,
	Tod:		to,
	Ts:		ts,
	Ed:		ed,
	Dmin:		Dmin,
    }
}
func newForce(yi, mo, sp, we, bs, to, bo string) Force {
    yieldForce,_ := strconv.ParseFloat(yi, 64)
    modElas,_ := strconv.ParseFloat(mo, 64)
    spaceChord,_ := strconv.ParseFloat(sp, 64)
    weight,_ := strconv.ParseFloat(we, 64)
    bSeat,_ := strconv.ParseFloat(bs, 64)
    topChord,_ := strconv.ParseFloat(to, 64)
    bottomChord,_ := strconv.ParseFloat(bo, 64)
    return Force {
	YieldStress:	yieldForce,
	ModElas:	modElas,
	SpaceChord:	spaceChord,
	Weight:		weight,
	BSeat:		bSeat,
	TopChord:	topChord,
	BottomChord:	bottomChord,
    }
}
func newAnglProp (st, at, rxt, rzt, ryt, yt, ixt, iyt, bt, tt, qt, sb, ab, rxb, rzb, ryb, yb, ixb, iyb, bb, tb, qb string) AnglProp{
    areaTop,_ := strconv.ParseFloat(at, 32)
    rxTop,_ := strconv.ParseFloat(rxt, 32)
    rzTop,_ := strconv.ParseFloat(rzt, 32)
    ryTop,_ := strconv.ParseFloat(ryt, 32)
    yTop,_ := strconv.ParseFloat(yt, 32)
    ixTop,_ := strconv.ParseFloat(ixt, 32)
    iyTop,_ := strconv.ParseFloat(iyt, 32)
    bTop,_ := strconv.ParseFloat(bt, 32)
    tTop,_ := strconv.ParseFloat(tt, 32)
    qTop,_ := strconv.ParseFloat(qt, 32)
    areaBot,_ := strconv.ParseFloat(ab, 32)
    rxBot,_ := strconv.ParseFloat(rxb, 32)
    rzBot,_ := strconv.ParseFloat(rzb, 32)
    ryBot,_ := strconv.ParseFloat(ryb, 32)
    yBot,_ := strconv.ParseFloat(yb, 32)
    ixBot,_ := strconv.ParseFloat(ixb, 32)
    iyBot,_ := strconv.ParseFloat(iyb, 32)
    bBot,_ := strconv.ParseFloat(bb, 32)
    tBot,_ := strconv.ParseFloat(tb, 32)
    qBot,_ := strconv.ParseFloat(qb, 32)
    return AnglProp{
	SecTop:		st,
	AreaTop:	areaTop,
	RxTop:		rxTop,
	RzTop:		rzTop,
	RyTop:		ryTop,
	YTop:		yTop,
	IxTop:		ixTop,
	IyTop:		iyTop,
	BTop:		bTop,
	TTop:		tTop,
	QTop:		qTop,
	SecBot:		sb,
	AreaBot:	areaBot,
	RxBot:		rxBot,
	RzBot:		rzBot,
	RyBot:		ryBot,
	YBot:		yBot,
	IxBot:		ixBot,
	IyBot:		iyBot,
	BBot:		bBot,
	TBot:		tBot,
	QBot:		qBot,
    }
}
func newBridgingProp (brid1, brid2 string) BridgingProp {
    bridging1,_ := strconv.ParseFloat(brid1, 64)
    bridging2,_ := strconv.ParseFloat(brid2, 64)
    return BridgingProp{
	Brdgng1: bridging1,
	Brdgng2: bridging2,
    }
}
func newMemberInput(numInput, elemName, part, mark, crimped, secction, midPanel string) MemberInput {
    return MemberInput{
	InputName: numInput,
	ElemName: elemName,
	Part: part,
	Mark: mark,
	Crimped: crimped,
	Secction: secction,
	MidPanel: midPanel,
    }
}
type Page struct {
    Data Data
    Geometry Geometry
    ResGeom ResGeom
    Material Material
    ResMater ResMater 
    ResBrid ResBrid
    WebMember WebMember
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
	    newResProp("", "", "", "", "", "", ""),
	},
    }
}
func newMaterial() Material{
    return Material{
	Forces: []Force{
	    newForce("50000", "29000", "1", "", "", "", ""),
	},
    }
}
func newResMater() ResMater{
    return ResMater{
	AnglProps: []AnglProp{
	    newAnglProp("", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""),
	},
    }
}
func newResBrid() ResBrid{
    return ResBrid{
	BridgingProps: []BridgingProp{
	    newBridgingProp("", ""),
	},
    }
}
func newWebMember() WebMember{
    return WebMember{
	MemberInputs: []MemberInput{
	    newMemberInput("", "", "", "", "", "", ""),
	},
    }
}

func newPage() Page {
    return Page {
	Data: newData(),
	Geometry: newGeometry(),
	ResGeom: newResGeom(),
	Material: newMaterial(),
	ResMater: newResMater(),
	ResBrid: newResBrid(),
	WebMember: newWebMember(),
    }
}
func main() {
    t := &Template{
	templates: template.Must(template.ParseGlob("views/*.html")),
    }
    e := echo.New()
    page := newPage()
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
    /*		Styles		*/
    e.GET("/styles", func(c echo.Context) error {
	return c.File("static/styles.css")
    })
    e.POST("/material", func(c echo.Context) error {
	page.WebMember.MemberInputs = page.WebMember.MemberInputs[:0]
	yieldStress := c.FormValue("yieldStress")
	modElas := c.FormValue("modElas")
	spaceChord := c.FormValue("spaceChord")
	weight := c.FormValue("weight")
	bSeat := c.FormValue("bSeat")
	topChord := c.FormValue("topChord")
	bottomChord := c.FormValue("bottomChord")
	force := newForce(
	    yieldStress,
	    modElas,
	    spaceChord,
	    weight,
	    bSeat,
	    topChord,
	    bottomChord,
	)
	page.Material.Forces = append(page.Material.Forces, force)
	page.Material.Forces = page.Material.Forces[1:]
	anglsProp := getAnglProps(
	    page.Material.Forces[0].TopChord, 
	    page.Material.Forces[0].BottomChord,
	    page.Material.Forces[0].SpaceChord,
	)
	page.ResMater.AnglProps = append(page.ResMater.AnglProps, anglsProp)
	page.ResMater.AnglProps = page.ResMater.AnglProps[1:]

	tod,_  := strconv.Atoi(page.ResGeom.ResProps[0].Tod)
	ts,_  := strconv.Atoi(page.ResGeom.ResProps[0].Ts)

	halfTod := tod/2
	halfTs := ts/2

	totalAngles := halfTod +1 + halfTs
	// totalAngles = totalAngles/2
	fmt.Println()
	fmt.Println(totalAngles)
	fmt.Println()

	obj := newMemberInput("", "", "", "", "", "", "")
	fmt.Printf("obj\t%v\n", obj)

	for i := 1; i < totalAngles; i++{
	    obj.InputName = fmt.Sprintf("%v", i)
	    if i == 1 {
		obj.ElemName = "sv"
	    } else if i <= halfTod+1 && i != 1{
		obj.ElemName = fmt.Sprintf("w%v", i)
	    } else {
		obj.ElemName = fmt.Sprintf("v%v", i - (halfTod+1))
	    }
	    page.WebMember.MemberInputs = append(page.WebMember.MemberInputs, obj)
	    //fmt.Println(page.WebMember.MemberInputs)
	}
	// page.WebMember.MemberInputs = page.WebMember.MemberInputs[1:]
	// fmt.Println(page.WebMember.MemberInputs)
	return c.Render(http.StatusOK, "materialResponse", page )
    })

    e.POST("/member", func(c echo.Context) error {
	/* TOOD: 
		make function to get angle info
	*/
	// read csv file 

	file, err := os.Open("Propiedades.csv")
	if err != nil {
	    fmt.Println("Error: ", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	record, err := reader.ReadAll()
	if err != nil {
	    fmt.Println("Error: ", err)
	}

	membersTable := len(page.WebMember.MemberInputs)
	for i := 0; i < membersTable; i++ {
	    si := strconv.Itoa(i + 1)

	    part := "part" + si
	    mark := "mark" + si
	    crimped := "crimped" + si

	    page.WebMember.MemberInputs[i].Part = c.FormValue(part)
	    page.WebMember.MemberInputs[i].Mark = c.FormValue(mark)
	    page.WebMember.MemberInputs[i].Crimped = c.FormValue(crimped)

	    mmark,_ := strconv.Atoi(page.WebMember.MemberInputs[i].Mark)
	    page.WebMember.MemberInputs[i].Secction = record[mmark][2]

	    fmt.Println(page.WebMember.MemberInputs[i])

	}
	fmt.Println(page.WebMember.MemberInputs)
	return c.Render(http.StatusOK, "webMem", page )
	return nil
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
	lbe2 := bChordtobPanel(
	    page.Geometry.Properties[0].fepl,
	    page.Geometry.Properties[0].sepl,
	    page.Geometry.Properties[0].ipl,
	    page.Geometry.Properties[0].lbe,
	)
	dLength := designLength(page.Geometry.Properties[0].span)
	tip := totalInteriorPanel(
	    dLength,
	    page.Geometry.Properties[0].fepl,
	    page.Geometry.Properties[0].sepl,
	    page.Geometry.Properties[0].ipl,
	)
	tod := totalDiagonal(tip)
	ts := totalStruts(tip)
	tip = math.Round(tip)
	tod = math.Round(tod)
	ts = math.Round(ts)

	dmin := spanDepth(page.Geometry.Properties[0].span)
	q := strconv.FormatFloat(lbe2, 'g', -1, 64)
	w := strconv.FormatFloat(dLength, 'g', -1, 64)
	e := strconv.FormatFloat(tip, 'g', -1, 64)
	r := strconv.FormatFloat(tod, 'g', -1, 64)
	t := strconv.FormatFloat(ts, 'g', -1, 64)
	y := strconv.FormatFloat(dmin, 'g', -1, 64)
	page.ResGeom.ResProps[0].Lbe2 = q
	page.ResGeom.ResProps[0].DLength = w
	page.ResGeom.ResProps[0].Tip = e
	page.ResGeom.ResProps[0].Tod = r
	page.ResGeom.ResProps[0].Ts = t
	page.ResGeom.ResProps[0].Dmin = y
	return c.Render(http.StatusOK, "geometryResponse", page.ResGeom)
    })
    e.Logger.Fatal(e.Start(":8080"))
}
