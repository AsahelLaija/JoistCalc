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
func bChordtobPanel(fepl, sepl, ipl, epbc float64) float64 {
    return (fepl + sepl + ipl) - epbc
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
    yStr,_ := strconv.ParseFloat(y, 32)
    ixFlt := 2*ixStr
    ixDAngle := strconv.FormatFloat(ixFlt, 'g', -1, 64)
    result := (sbca/2) + yStr
    
    // res is egual to ((sbca/2)+Ytc)^2
    res := math.Pow(result, 2)
    iyFlt := 2*(ixStr + areaStr*res)
    iyDAngle := strconv.FormatFloat(iyFlt, 'g', -1, 64)
    return areaDAngle, ixDAngle, iyDAngle
}
func RoundTo(n float64, decimals uint32) float64 {
  return math.Round(n*math.Pow(10, float64(decimals))) / math.Pow(10, float64(decimals))
}
func anglePropsDOWSlend( mark, i int, crimped string, page *Page ) {
    DOWSlendRats := page.TableDOWSlend.DOWSlendRats
    ElemName := page.WebMember.MemberInputs[i].ElemName

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
    if crimped == "yes" {
	// if is crimped the rx is rz
	RX,_ := strconv.ParseFloat(record[mark][8], 64)
	DOWSlendRats[i].Rx = RoundTo(RX, 4)

	// if is crimped the ry is rw
	RY,_ := strconv.ParseFloat(record[mark][9], 64)
	DOWSlendRats[i].Ry = RoundTo(RY, 4)
    }
    de := page.ResGeom.ResProps[0].Ed
    beta := page.TableDOWForce.DOWForces[0].Beta
    gamma := page.TableDOWForce.DOWForces[0].Gamma
    delta := page.TableDOWForce.DOWForces[0].Delta
    alpha := page.TableDOWForce.DOWForces[0].Alpha
    var Ix float64
    if ElemName == "sv" {
	gammaRadian := gamma / (180/math.Pi)
	Ix = de / math.Sin(gammaRadian)
	DOWSlendRats[i].IX = RoundTo(Ix, 4)
	DOWSlendRats[i].IY = RoundTo(Ix, 4)
	DOWSlendRats[i].IZ = RoundTo(Ix, 4)
	DOWSlendRats[i].Limit = 200
    } else if ElemName == "w2" {
	betaRadian := beta / (180/math.Pi)
	Ix = de / math.Sin(betaRadian)
	DOWSlendRats[i].IX = RoundTo(Ix, 4)
	DOWSlendRats[i].IY = RoundTo(Ix, 4)
	DOWSlendRats[i].IZ = RoundTo(Ix, 4)
	DOWSlendRats[i].Limit = 240
    } else if ElemName == "w3" {
	deltaRadian := delta / (180/math.Pi)
	Ix = de / math.Sin(deltaRadian)
	DOWSlendRats[i].IX = RoundTo(Ix, 4)
	DOWSlendRats[i].IY = RoundTo(Ix, 4)
	DOWSlendRats[i].IZ = RoundTo(Ix, 4)
	DOWSlendRats[i].Limit = 200 
    } else if ElemName[0] == 'v' { 
	DOWSlendRats[i].IX = RoundTo(de, 4) 
	DOWSlendRats[i].IY = RoundTo(de, 4) 
	DOWSlendRats[i].IZ = RoundTo(de, 4)
	DOWSlendRats[i].Limit = 200
    } else{
	alphaRadian := alpha / (180/math.Pi)
	Ix = de / math.Sin(alphaRadian)
	DOWSlendRats[i].IX = RoundTo(Ix, 4)
	DOWSlendRats[i].IY = RoundTo(Ix, 4)
	DOWSlendRats[i].IZ = RoundTo(Ix, 4)
	if 1 == i % 2 {
	    DOWSlendRats[i].Limit = 240
	} else {
	    DOWSlendRats[i].Limit = 200
	}
    }

    DOWSlendRats[i].Ixrx = RoundTo(DOWSlendRats[i].IX / DOWSlendRats[i].Rx, 4)
    DOWSlendRats[i].Iyry = RoundTo(DOWSlendRats[i].IY / DOWSlendRats[i].Ry, 4)
    DOWSlendRats[i].Lrgov = larger(DOWSlendRats[i].Ixrx, DOWSlendRats[i].Iyry)
    if DOWSlendRats[i].Lrgov < DOWSlendRats[i].Limit {
	DOWSlendRats[i].Check = "OK"
    } else {
	DOWSlendRats[i].Check = "NOT OK"
    }
}

func designWeb(mark, i int, page *Page) {
    DOWWeb := page.TableDOWWeb.DOWWebs
    MemberInputs := page.WebMember.MemberInputs
    lbe := page.Geometry.Properties[0].epbc 
    lbe2 := page.ResGeom.ResProps[0].Lbe2

    lep1 := page.Geometry.Properties[0].fepl 
    lep2 := page.Geometry.Properties[0].sepl 
    lip := page.Geometry.Properties[0].ipl

    W := page.TableResForce.ResForces[0].KipUdlWu 

    de := page.ResGeom.ResProps[0].Ed

    L := page.ResGeom.ResProps[0].DLength
    R := page.TableDOWForce.DOWForces[0].Rmax
    betaR := page.TableDOWForce.DOWForces[0].Beta * (math.Pi / 180)
    alphaR := page.TableDOWForce.DOWForces[0].Alpha * (math.Pi / 180)
    gammaR := page.TableDOWForce.DOWForces[0].Gamma * (math.Pi / 180)
    deltaR := page.TableDOWForce.DOWForces[0].Delta * (math.Pi / 180)

    if MemberInputs[i].ElemName == "sv" {
	part1 := W * ((lep1 + lep2) / 2)
	part2 := (R - (W * (lep1 /2))) / math.Tan(betaR)
	vmin := part1  + .005 * part2
	if vmin < page.TableDOWForce.DOWForces[0].Vmin {
	    DOWWeb[i].Vmin = RoundTo(page.TableDOWForce.DOWForces[0].Vmin, 4)
	} else {
	    DOWWeb[i].Vmin = vmin
	}
	DOWWeb[i].FcMin = RoundTo(DOWWeb[i].Vmin / math.Sin(gammaR), 4)

    } else if MemberInputs[i].ElemName == "w2" {
	DOWWeb[i].XPC = lbe
	DOWWeb[i].XPE = 0
	DOWWeb[i].EQV1 = RoundTo(W * ((L / 2) - DOWWeb[i].XPE), 4)
	vmin := DOWWeb[i].EQV1
	if vmin < page.TableDOWForce.DOWForces[0].Vmin {
	    DOWWeb[i].Vmin = RoundTo(page.TableDOWForce.DOWForces[0].Vmin, 4)
	} else {
	    DOWWeb[i].Vmin = vmin
	}
	DOWWeb[i].FtMin = RoundTo(DOWWeb[i].Vmin / math.Sin(betaR), 4)
    } else if MemberInputs[i].ElemName == "w3" {
	DOWWeb[i].XPC = RoundTo(lep1 + lep2, 2)
	DOWWeb[i].XPE = DOWWeb[i - 1].XPC
	DOWWeb[i].EQV1 = RoundTo(W * ((L / 2) - DOWWeb[i].XPE), 4)
	vmin := DOWWeb[i].EQV1
	if vmin < page.TableDOWForce.DOWForces[0].Vmin {
	    DOWWeb[i].Vmin = RoundTo(page.TableDOWForce.DOWForces[0].Vmin, 4)
	} else {
	    DOWWeb[i].Vmin = vmin
	}
	DOWWeb[i].FcMin = RoundTo(DOWWeb[i].Vmin / math.Sin(deltaR), 4)
    } else if MemberInputs[i].ElemName == "w4" {
	DOWWeb[i].XPC = lbe + lbe2
	DOWWeb[i].XPE = DOWWeb[i - 1].XPC
	DOWWeb[i].EQV1 = RoundTo(W * ((L / 2) - DOWWeb[i].XPE), 4)
	vmin := DOWWeb[i].EQV1
	if vmin < page.TableDOWForce.DOWForces[0].Vmin {
	    DOWWeb[i].Vmin = RoundTo(page.TableDOWForce.DOWForces[0].Vmin, 4)
	} else {
	    DOWWeb[i].Vmin = vmin
	}
	DOWWeb[i].FtMin = RoundTo(DOWWeb[i].Vmin / math.Sin(alphaR), 4)
    } else if (MemberInputs[i].ElemName != "w4" && 
    isDiagonal(MemberInputs[i].ElemName)) {

	DOWWeb[i].XPC = DOWWeb[i - 1].XPC + lip

	DOWWeb[i].XPE = DOWWeb[i - 1].XPC

	DOWWeb[i].EQV1 = RoundTo(W * ((L / 2) - DOWWeb[i].XPE), 4)

	vmin := DOWWeb[i].EQV1
	if vmin < page.TableDOWForce.DOWForces[0].Vmin {
	    DOWWeb[i].Vmin = RoundTo(page.TableDOWForce.DOWForces[0].Vmin, 4)
	} else {
	    DOWWeb[i].Vmin = vmin
	}
	wIndex,_ := strconv.Atoi(MemberInputs[i].ElemName[1:])

	if wIndex >= 5 && wIndex % 2 == 0 {
	    DOWWeb[i].FtMin = RoundTo(DOWWeb[i].Vmin / math.Sin(alphaR), 4)
	} else if wIndex % 2 == 1 {
	    DOWWeb[i].FcMin = RoundTo(DOWWeb[i].Vmin / math.Sin(alphaR), 4)
	}
    } else if MemberInputs[i].ElemName[0] == 'v' {
	diagIndex := MemberInputs[i].ElemName[1] - '0'
	countLip := float64((diagIndex * 2) - 1)
	x := lep1 + lep2 + (countLip * lip)
	compChord := ((W * x) / (2 * de)) * (L - x)
	p := (lip * W) + (.005 * compChord)
	fmt.Println(lip, W, countLip, x, compChord)
	DOWWeb[i].FcMin = RoundTo(p, 2)
    }
    DOWWeb[i].DesignFTens = larger(DOWWeb[i].FtMin, DOWWeb[i].LiftFTens)
    DOWWeb[i].DesignFComp = larger(DOWWeb[i].FcMin, DOWWeb[i].LiftFComp)
}
func isDiagonal(elem string) bool {
    iElem,_ := strconv.Atoi(elem[1:])
    if (elem[0:1] == "w" && iElem > 3) {
	return true
    } else {
	return false
    }
    return false
}

func designStress(mark, i int, page *Page) {
    DOWDesigns := page.TableDOWDesign.DOWDesigns
    MemberInputs := page.WebMember.MemberInputs
    DOWEfectives := page.TableDOWEfective.DOWEfectives
    Fy := page.Material.Forces[0].YieldStress
    E := page.Material.Forces[0].ModElas
    rs := page.TableResistance.ResistanceFactors[0].TensionValue

    Fy_E := E / Fy

    b := angleProp(mark, 5)
    bF, _ := strconv.ParseFloat(b, 64)
    t := angleProp(mark, 4)
    tF, _ := strconv.ParseFloat(t, 64)
    btF := bF/tF
    DOWDesigns[i].Compbt = RoundTo(btF, 4)
    
    area := angleProp(mark, 3)
    areaF, _ := strconv.ParseFloat(area, 64)
    DOWDesigns[i].A = RoundTo(areaF, 4)

    if (MemberInputs[i].Crimped == "yes" &&
	MemberInputs[i].ElemName == "w2" ||
	MemberInputs[i].ElemName == "w3" ||
	MemberInputs[i].ElemName == "sv" ){
	    Q := (5.25 / btF) + tF
	    DOWDesigns[i].CompQ = RoundTo(Q, 4)

    } else if (MemberInputs[i].Crimped == "yes" &&
	MemberInputs[i].ElemName[0] == 'v' || 
	isDiagonal(MemberInputs[i].ElemName)){
	    if btF <= 0.45 * (math.Sqrt(Fy_E)) {
		DOWDesigns[i].CompQ = 1
	    } else if (btF >= 0.91 * (math.Sqrt(Fy_E))) {
		fmt.Println("2 ", MemberInputs[i].ElemName)
		Q := 0.53 * (E / (Fy * math.Pow(btF, 2)))
		DOWDesigns[i].CompQ = RoundTo(Q, 4)
	    } else if (btF <= 0.91 * (math.Sqrt(Fy_E))) {

		fmt.Println("3 ", MemberInputs[i].ElemName, btF, Fy, E)
		Q := 1.34 - 0.76 * (btF * math.Sqrt(Fy / E))
		DOWDesigns[i].CompQ = RoundTo(Q, 4)

	    } else {
		// fmt.Println("tabla ", MemberInputs[i].ElemName)
		Q := angleProp(mark, 12)
		QF,_ := strconv.ParseFloat(Q, 64)
		DOWDesigns[i].CompQ = RoundTo(QF, 4)
	    }
    } else if (MemberInputs[i].Crimped == "no" &&
	MemberInputs[i].ElemName == "w2" ||
	MemberInputs[i].ElemName == "w3" ||
	MemberInputs[i].ElemName == "sv" ){
	    fmt.Println("tabla ", MemberInputs[i].ElemName)
	    Q := angleProp(mark, 12)
	    QF,_ := strconv.ParseFloat(Q, 64)
	    DOWDesigns[i].CompQ = RoundTo(QF, 4)
    }

    a := math.Sqrt(E / DOWDesigns[i].CompQ * Fy)
    SLRgov := DOWEfectives[i].SLRgov

    Fe := RoundTo(math.Pow(math.Pi, 2) * E / math.Pow(SLRgov, 2), 4)
    DOWDesigns[i].CompFe = Fe
    Qr := DOWDesigns[i].CompQ

    if SLRgov <= a {
	DOWDesigns[i].CompFcr = RoundTo(Qr * (math.Pow(0.658, (Qr * Fy)/Fe)) * Fy, 4)
    } else {
	DOWDesigns[i].CompFcr = RoundTo(0.877 * Fe, 4)
    }

    DOWDesigns[i].CompFc = RoundTo(DOWDesigns[i].CompFcr * .9, 4)
    DOWDesigns[i].CompPuc = RoundTo(DOWDesigns[i].CompFcr * areaF, 4)

    DOWDesigns[i].TenFt = Fy * rs
    DOWDesigns[i].TenPut = RoundTo(DOWDesigns[i].TenFt * areaF, 4)
}

func efectiveSlendernes(i int, page *Page) {
    MemberInputs := page.WebMember.MemberInputs
    DOWEfectives := page.TableDOWEfective.DOWEfectives
    DOWSlendRats := page.TableDOWSlend.DOWSlendRats

    if MemberInputs[i].Crimped == "yes" {
	DOWEfectives[i].Klrx = 0.75
	DOWEfectives[i].Klry = 0.90
	DOWEfectives[i].Klrz = 0
	DOWEfectives[i].Klsrz = 0
    } else if MemberInputs[i].Fill == "yes" {
	DOWEfectives[i].Klrx = 0.75
	DOWEfectives[i].Klry = 0.94
	DOWEfectives[i].Klrz = 0
	DOWEfectives[i].Klsrz = 1.0
    } else if MemberInputs[i].Fill == "no" {
	DOWEfectives[i].Klrx = 0
	DOWEfectives[i].Klry = 0
	DOWEfectives[i].Klrz = 0.90
	DOWEfectives[i].Klsrz = 0
    }
    DOWEfectives[i].SlendKlrx = DOWSlendRats[i].Ixrx * DOWEfectives[i].Klrx
    DOWEfectives[i].SlendKlry = DOWSlendRats[i].Iyry * DOWEfectives[i].Klry
    DOWEfectives[i].SlendKlrz = DOWSlendRats[i].Izrz * DOWEfectives[i].Klrz

    DOWEfectives[i].SLRgov = larger(
	DOWEfectives[i].SlendKlrx, 
	DOWEfectives[i].SlendKlry,
	DOWEfectives[i].SlendKlrz)

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
	// record[rT][9],	//RyTop
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
	// record[rB][9],	//RyBot
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
type Propertie struct {
    TrussType,
    JoistType	string
    deflexion,
    span,
    fepl,
    sepl,
    ipl,
    epbc,
    depth	float64
}
type Properties = []Propertie
type Geometry struct{
    Properties Properties
}

func newPropertie(tt, jt, df, sp, fe, se, ip, lb, de string) Propertie {
    deflexion,_ := strconv.ParseFloat(df, 64)
    span,_ := strconv.ParseFloat(sp, 64)
    fepl,_ := strconv.ParseFloat(fe, 64)
    sepl,_ := strconv.ParseFloat(se, 64)
    ipl,_ := strconv.ParseFloat(ip, 64)
    epbc,_ := strconv.ParseFloat(lb, 64)
    depth,_ := strconv.ParseFloat(de, 64)

    return Propertie {
	TrussType:	tt,
	JoistType:	jt,
	deflexion:	deflexion,
	span:		span,
	fepl:		fepl,
	sepl:		sepl,
	ipl:		ipl,
	epbc:		epbc,
	depth:		depth,
    }
} 
func newGeometry() Geometry {
    return Geometry{
	Properties: []Propertie{
	     newPropertie("warren", " roof", "240", "49.21", "27.28", "26", "24", "41.1", "28"),
	},
    }
}

/*	Geometry Response	*/
type ResProp struct {
    Lbe2, DLength, Tip, Tod, Ts, Ed, Dmin, Lbrdng1, Lbrdng2 float64
}
type ResProps = []ResProp
type ResGeom struct {
    ResProps ResProps
}
func newResProp() ResProp {
    return ResProp {}
}
func newResGeom() ResGeom {
    return ResGeom{
	ResProps: []ResProp{
	    newResProp(),
	},
    }
}

/*	Stress		*/
type Force struct {
    YieldStress, ModElas, SpaceChord, Weight, BSeat, TopChord, BottomChord, Udlw, LlWLL, NsLRFD float64
    FillTopChord, FillBotChord, TopChordEP1, TopChordEP2 bool
}
type Forces = []Force
type Material struct {
    Forces Forces
}
func newForce(yi, mo, sp, we, bs, to, bo, tf, bf, udlwA, llWLLA, nsLRFDA, tcep1, tcep2 string) Force {
    // Note the capital A means "function argument"
    // Stress
    yieldForce,_ := strconv.ParseFloat(yi, 64)
    modElas,_ := strconv.ParseFloat(mo, 64)
    spaceChord,_ := strconv.ParseFloat(sp, 64)
    weight,_ := strconv.ParseFloat(we, 64)
    bSeat,_ := strconv.ParseFloat(bs, 64)

    // Loads
    udlw,_ := strconv.ParseFloat(udlwA, 64)
    llWLL,_ := strconv.ParseFloat(llWLLA, 64)
    nsLRFD,_ := strconv.ParseFloat(nsLRFDA, 64)

    // Chords
    topChord,_ := strconv.ParseFloat(to, 64)
    bottomChord,_ := strconv.ParseFloat(bo, 64)
    var fillTopChord bool
    var fillBotChord bool
    var topChordEP1 bool
    var topChordEP2 bool
    if tf == "yes" {
	fillTopChord = true
    } else {
	fillTopChord = false
    }
    if bf == "yes" {
	fillBotChord = true
    } else {
	fillBotChord = false
    }
    if tcep1 == "yes" {
	topChordEP1 = true
    } else {
	topChordEP1 = false
    }
    if tcep2 == "yes" {
	topChordEP2 = true
    } else {
	topChordEP2 = false
    }
    return Force {
	YieldStress:	yieldForce,
	ModElas:	modElas,
	SpaceChord:	spaceChord,
	Weight:		weight,
	BSeat:		bSeat,
	TopChord:	topChord,
	BottomChord:	bottomChord,
	FillTopChord:	fillTopChord,
	FillBotChord:	fillBotChord,
	TopChordEP1:	topChordEP1,
	TopChordEP2:	topChordEP2,
	Udlw:	udlw,
	LlWLL:	llWLL,
	NsLRFD: nsLRFD,
    }
}
func newMaterial() Material{
    return Material{
	Forces: []Force{
	    newForce("50000", "29000", "1", "", "", "", "", "", "", "", "", "", "", ""),
	},
    }
}

// 3. Design Loads
type ResForce struct {
    KipUdlWu float64

    KipLlWLL float64
    KipNsLRFD float64
    MaxDsgnMoment float64
    ChordForce float64
}
type ResForces = []ResForce
type TableResForce struct {
    ResForces ResForces
}

func newResForce() ResForce{
    return ResForce{}
}
func newTableResForce() TableResForce {
    return TableResForce {
	ResForces: []ResForce{
	    newResForce(),
 	},
    }
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
type AnglProps = []AnglProp
type ResMater struct {
    AnglProps AnglProps
}

func newAnglProp (st, at, rxt, rzt, yt, ixt, iyt, bt, tt, qt, sb, ab, rxb, rzb, yb, ixb, iyb, bb, tb, qb string) AnglProp{
    areaTop,_ := strconv.ParseFloat(at, 32)
    areaTop = roundFloat(areaTop, 4)
    rxTop,_ := strconv.ParseFloat(rxt, 32)
    rxTop = roundFloat(rxTop, 4)
    rzTop,_ := strconv.ParseFloat(rzt, 32)
    rzTop = roundFloat(rzTop, 4)
    // ryTop,_ := strconv.ParseFloat(ryt, 32)
    yTop,_ := strconv.ParseFloat(yt, 32)
    yTop = roundFloat(yTop, 4)
    ixTop,_ := strconv.ParseFloat(ixt, 32)
    ixTop = roundFloat(ixTop, 4)
    iyTop,_ := strconv.ParseFloat(iyt, 32)
    iyTop = roundFloat(iyTop, 4)
    bTop,_ := strconv.ParseFloat(bt, 32)
    bTop = roundFloat(bTop, 4)
    tTop,_ := strconv.ParseFloat(tt, 32)
    tTop = roundFloat(tTop, 4)
    qTop,_ := strconv.ParseFloat(qt, 32)
    qTop = roundFloat(qTop, 4)
    areaBot,_ := strconv.ParseFloat(ab, 32)
    areaBot = roundFloat(areaBot, 4)
    rxBot,_ := strconv.ParseFloat(rxb, 32)
    rxBot = roundFloat(rxBot, 4)
    rzBot,_ := strconv.ParseFloat(rzb, 32)
    rzBot = roundFloat(rzBot, 4)
    // ryBot,_ := strconv.ParseFloat(ryb, 32)
    yBot,_ := strconv.ParseFloat(yb, 32)
    yBot = roundFloat(yBot, 4)
    ixBot,_ := strconv.ParseFloat(ixb, 32)
    ixBot = roundFloat(ixBot, 4)
    iyBot,_ := strconv.ParseFloat(iyb, 32)
    iyBot = roundFloat(iyBot, 4)
    bBot,_ := strconv.ParseFloat(bb, 32)
    bBot = roundFloat(bBot, 4)
    tBot,_ := strconv.ParseFloat(tb, 32)
    tBot = roundFloat(tBot, 4)
    qBot,_ := strconv.ParseFloat(qb, 32)
    qBot = roundFloat(qBot, 4)

    ryTop := math.Sqrt(iyTop / areaTop)
    ryTop = roundFloat(ryTop, 4)
    ryBot := math.Sqrt(iyBot / areaBot)
    ryBot = roundFloat(ryBot, 4)

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
func newResMater() ResMater{
    return ResMater{
	AnglProps: []AnglProp{
	    newAnglProp("", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""),
	},
    }
}

type BridgingProp struct {
    Brdgng1,
    Brdgng2 float64
}
type BridgingProps = []BridgingProp
type ResBrid struct {
    BridgingProps BridgingProps
}

func newBridgingProp () BridgingProp {
    return BridgingProp{}
}
func newResBrid() ResBrid{
    return ResBrid{
	BridgingProps: []BridgingProp{
	    newBridgingProp(),
	},
    }
}

/*	 2.- Material and Section "Table Web member" :)		*/
type MemberInput struct {
    InputName string 
    ElemName string
    // request
    Part string
    Mark string
    Crimped string
    Fill string
    // response
    Secction string
    MidPanel string
    rxDouble float64
    ryDouble float64
    rzDouble float64
    rxCrimped float64
    ryCrimped float64
    rxUncrimped float64
    ryUncrimped float64
    A float64
    Q float64
    S float64
    T float64
}
type MemberInputs = []MemberInput
type WebMember struct {
    MemberInputs MemberInputs
}

func newMemberInput() MemberInput {
    return MemberInput{}
}
func newWebMember() WebMember{
    return WebMember{
	MemberInputs: []MemberInput{
	    newMemberInput(),
	},
    }
}

type AnalysisResult struct {
    BotChrdTnsForce float64
    BotChrdTnsFt float64
    BotChrdCmpForce float64
    BotChrdCmpFc float64
    BotChrdMomDwnMu float64
    BotChrdMomDwnFb float64
    BotChrdMomUpMu float64
    BotChrdMomUpfb float64

    TopChrdTnsForceMid float64
    TopChrdTnsFtMid float64
    TopChrdCmpForceMid float64
    TopChrdCmpFcMid float64
    TopChrdMomDwnMuMid float64
    TopChrdMomDwnFbMid float64
    TopChrdMomUpMuMid float64
    TopChrdMomUpfbMid float64

    TopChrdTnsForcePoint float64
    TopChrdTnsFtPoint float64
    TopChrdCmpForcePoint float64
    TopChrdCmpFcPoint float64
    TopChrdMomDwnMuPoint float64
    TopChrdMomDwnFbPoint float64
    TopChrdMomUpMuPoint float64
    TopChrdMomUpfbPoint float64

    TopChrdEP1TnsForceMid float64 
    TopChrdEP1TnsFtMid float64    
    TopChrdEP1CmpForceMid float64 
    TopChrdEP1CmpFcMid float64  
    TopChrdEP1MomDwnMuMid float64 
    TopChrdEP1MomDwnFbMid float64 
    TopChrdEP1MomUpMuMid float64  
    TopChrdEP1MomUpfbMid float64  

    TopChrdEP1TnsForcePoint float64 
    TopChrdEP1TnsFtPoint float64    
    TopChrdEP1CmpForcePoint float64 
    TopChrdEP1CmpFcPoint float64    
    TopChrdEP1MomDwnMuPoint float64 
    TopChrdEP1MomDwnFbPoint float64 
    TopChrdEP1MomUpMuPoint float64  
    TopChrdEP1MomUpfbPoint float64  

    TopChrdEP2TnsForceMid float64 
    TopChrdEP2TnsFtMid float64    
    TopChrdEP2CmpForceMid float64 
    TopChrdEP2CmpFcMid float64    
    TopChrdEP2MomDwnMuMid float64 
    TopChrdEP2MomDwnFbMid float64 
    TopChrdEP2MomUpMuMid float64  
    TopChrdEP2MomUpfbMid float64  

    TopChrdEP2TnsForcePoint float64 
    TopChrdEP2TnsFtPoint float64    
    TopChrdEP2CmpForcePoint float64 
    TopChrdEP2CmpFcPoint float64    
    TopChrdEP2MomDwnMuPoint float64 
    TopChrdEP2MomDwnFbPoint float64 
    TopChrdEP2MomUpMuPoint float64  
    TopChrdEP2MomUpfbPoint float64  
}
type AnalysisResults = []AnalysisResult
type TableAnalysis struct {
    AnalysisResults AnalysisResults
}
func newAnalysisResult () AnalysisResult {
    return AnalysisResult {}
}
func newTableAnalysis() TableAnalysis {
    return TableAnalysis {
	AnalysisResults: []AnalysisResult{
	    newAnalysisResult(),
 	},
    }
}
// 5.- Design of Chords Begin
// Resistance Factor
type ResistanceFactor struct {
    TensionFactor string // value as default is "Phi t"
    TensionValue float64 // val
	CompressionDStress: compressionDStress,
	BendingFactor: bendingFactor,
	BendingValue: bendingValue,
	BendingDStress: bendingDStress,
    }
}
func newTableResistance() TableResistance {
    return TableResistance {
	ResistanceFactors: []ResistanceFactor{
	    newResistanceFactor("t", "c", "b", "0.9 Fy", "0.9Fcr", "0.9 Fy", 0.9, 0.9, 0.9),
 	},
    }
}
// Check Slenderness ratio
type CheckSlendernes struct {
    BCIx float64
    BCIy float64
    BCIz float64
    BCrx float64
    BCry float64
    BCrz float64
    BCIxrx float64
    BCIyry float64
    BCIzrz float64
    BCSLRym float64
    BCrgov float64
    BClimit float64
    BCcheck string

    TCIpMidIx float64
    TCIpMidIy float64
    TCIpMidIz float64
    TCIpMidrx float64
    TCIpMidry float64
    TCIpMidrz float64
    TCIpMidIxrx float64
    TCIpMidIyry float64
    TCIpMidIzrz float64
    TCIpMidSLRym float64
    TCIpMidrgov float64
    TCIpMidlimit float64
    TCIpMidcheck string
    
    TCIpPointIx float64
    TCIpPointIy float64
    TCIpPointIz float64
    TCIpPointrx float64
    TCIpPointry float64
    TCIpPointrz float64
    TCIpPointIxrx float64
    TCIpPointIyry float64
    TCIpPointIzrz float64
    TCIpPointSLRym float64
    TCIpPointrgov float64
    TCIpPointlimit float64
    TCIpPointcheck string

    TCEp1MidIx float64
    TCEp1MidIy float64
    TCEp1MidIz float64
    TCEp1Midrx float64
    TCEp1Midry float64
    TCEp1Midrz float64
    TCEp1MidIxrx float64
    TCEp1MidIyry float64
    TCEp1MidIzrz float64
    TCEp1MidSLRym float64
    TCEp1Midrgov float64
    TCEp1Midlimit float64
    TCEp1Midcheck string

    TCEp1PointIx float64
    TCEp1PointIy float64
    TCEp1PointIz float64
    TCEp1Pointrx float64
    TCEp1Pointry float64
    TCEp1Pointrz float64
    TCEp1PointIxrx float64
    TCEp1PointIyry float64
    TCEp1PointIzrz float64
    TCEp1PointSLRym float64
    TCEp1Pointrgov float64
    TCEp1Pointlimit float64
    TCEp1Pointcheck string

    TCEp2MidIx float64
    TCEp2MidIy float64
    TCEp2MidIz float64
    TCEp2Midrx float64
    TCEp2Midry float64
    TCEp2Midrz float64
    TCEp2MidIxrx float64
    TCEp2MidIyry float64
    TCEp2MidIzrz float64
    TCEp2MidSLRym float64
    TCEp2Midrgov float64
    TCEp2Midlimit float64
    TCEp2Midcheck string

    TCEp2PointIx float64
    TCEp2PointIy float64
    TCEp2PointIz float64
    TCEp2Pointrx float64
    TCEp2Pointry float64
    TCEp2Pointrz float64
    TCEp2PointIxrx float64
    TCEp2PointIyry float64
    TCEp2PointIzrz float64
    TCEp2PointSLRym float64
    TCEp2Pointrgov float64
    TCEp2Pointlimit float64
    TCEp2Pointcheck string
}
type CheckSlenderness = []CheckSlendernes
type TableCheck struct {
    CheckSlenderness CheckSlenderness
}
func newCheckSlendernes()CheckSlendernes{
    return CheckSlendernes{}
}
func newTableCheck() TableCheck{
    return TableCheck {
	CheckSlenderness: []CheckSlendernes{
	    newCheckSlendernes(),
	},
    }
}

// Efective Slenderness ratio
type EfectiveSlendernes struct {
    BCklrx float64
    BCklry float64
    BCklrz float64
    BCklsrz float64
    BCKLrx float64
    BCSlendklrx float64
    BCSlendklry float64 
    BCSlendklrz float64 
    BCSlendklsrz float64
    BCSLRgov float64
    BCFeKlrx float64

    TCIpMidklrx float64
    TCIpMidklry float64
    TCIpMidklrz float64
    TCIpMidklsrz float64
    TCIpMidKLrx float64
    TCIpMidSlendklrx float64
    TCIpMidSlendklry float64
    TCIpMidSlendklrz float64
    TCIpMidSlendklsrz float64
    TCIpMidSLRgov float64
    TCIpMidFeKlrx float64

    TCIpPointklrx float64
    TCIpPointklry float64
    TCIpPointklrz float64
    TCIpPointklsrz float64
    TCIpPointKLrx float64
    TCIpPointSlendklrx float64
    TCIpPointSlendklry float64
    TCIpPointSlendklrz float64
    TCIpPointSlendklsrz float64
    TCIpPointSLRgov float64
    TCIpPointFeKlrx float64

    TCEp1Midklrx float64
    TCEp1Midklry float64
    TCEp1Midklrz float64
    TCEp1Midklsrz float64
    TCEp1MidKLrx float64
    TCEp1MidSlendklrx float64
    TCEp1MidSlendklry float64
    TCEp1MidSlendklrz float64
    TCEp1MidSlendklsrz float64
    TCEp1MidSLRgov float64
    TCEp1MidFeKlrx float64

    TCEp1Pointklrx float64
    TCEp1Pointklry float64
    TCEp1Pointklrz float64
    TCEp1Pointklsrz float64
    TCEp1PointKLrx float64
    TCEp1PointSlendklrx float64
    TCEp1PointSlendklry float64
    TCEp1PointSlendklrz float64
    TCEp1PointSlendklsrz float64
    TCEp1PointSLRgov float64
    TCEp1PointFeKlrx float64

    TCEp2Midklrx float64
    TCEp2Midklry float64
    TCEp2Midklrz float64
    TCEp2Midklsrz float64
    TCEp2MidKLrx float64
    TCEp2MidSlendklrx float64
    TCEp2MidSlendklry float64
    TCEp2MidSlendklrz float64
    TCEp2MidSlendklsrz float64
    TCEp2MidSLRgov float64
    TCEp2MidFeKlrx float64

    TCEp2Pointklrx float64
    TCEp2Pointklry float64
    TCEp2Pointklrz float64
    TCEp2Pointklsrz float64
    TCEp2PointKLrx float64
    TCEp2PointSlendklrx float64
    TCEp2PointSlendklry float64
    TCEp2PointSlendklrz float64
    TCEp2PointSlendklsrz float64
    TCEp2PointSLRgov float64
    TCEp2PointFeKlrx float64
}
type EfectiveSlenderness = []EfectiveSlendernes
type TableEfective struct {
    EfectiveSlenderness EfectiveSlenderness
}
func newEfectiveSlendernes() EfectiveSlendernes{
	return EfectiveSlendernes{}
}
func newTableEfective() TableEfective{
    return TableEfective {
	EfectiveSlenderness: []EfectiveSlendernes{
	    newEfectiveSlendernes(),
	},
    }
}

// Design Stress and Check of relation capabilites vs solicitations
type Design struct{
    BCFt float64
    BCPut float64
    BCbt float64
    BCQ float64
    BCFetc float64
    BCFex float64
    BCFcr float64
    BCFc float64
    BCPuc float64
    BCcm float64
    BCFb float64

    TCIpMidFt float64
    TCIpMidPut float64
    TCIpMidbt float64
    TCIpMidQ float64
    TCIpMidFetc float64
    TCIpMidFex float64
    TCIpMidFcr float64
    TCIpMidFc float64
    TCIpMidPuc float64
    TCIpMidcm float64
    TCIpMidFb float64

    TCIpPointFt float64
    TCIpPointPut float64
    TCIpPointbt float64
    TCIpPointQ float64
    TCIpPointFetc float64
    TCIpPointFex float64
    TCIpPointFcr float64
    TCIpPointFc float64
    TCIpPointPuc float64
    TCIpPointcm float64
    TCIpPointFb float64

    TCEp1MidFt float64
    TCEp1MidPut float64
    TCEp1Midbt float64
    TCEp1MidQ float64
    TCEp1MidFetc float64
    TCEp1MidFex float64
    TCEp1MidFcr float64
    TCEp1MidFc float64
    TCEp1MidPuc float64
    TCEp1Midcm float64
    TCEp1MidFb float64

    TCEp1PointFt float64
    TCEp1PointPut float64
    TCEp1Pointbt float64
    TCEp1PointQ float64
    TCEp1PointFetc float64
    TCEp1PointFex float64
    TCEp1PointFcr float64
    TCEp1PointFc float64
    TCEp1PointPuc float64
    TCEp1Pointcm float64
    TCEp1PointFb float64

    TCEp2MidFt float64
    TCEp2MidPut float64
    TCEp2Midbt float64
    TCEp2MidQ float64
    TCEp2MidFetc float64
    TCEp2MidFex float64
    TCEp2MidFcr float64
    TCEp2MidFc float64
    TCEp2MidPuc float64
    TCEp2Midcm float64
    TCEp2MidFb float64

    TCEp2PointFt float64
    TCEp2PointPut float64
    TCEp2Pointbt float64
    TCEp2PointQ float64
    TCEp2PointFetc float64
    TCEp2PointFex float64
    TCEp2PointFcr float64
    TCEp2PointFc float64
    TCEp2PointPuc float64
    TCEp2Pointcm float64
    TCEp2PointFb float64
}
type Designs = []Design
type TableDesign struct {
    Designs Designs
}

func newDesign() Design{
	return Design{}
}

func newTableDesign() TableDesign{
    return TableDesign {
	Designs: []Design{
	    newDesign(),
	},
    }
}

// Check of relatin capabilities vs solicitations
type CapSol struct {
    BCPut float64
    BCPusol float64
    BCTenRat float64
    BCTenAllow float64
    BCFau float64
    BCFbu float64
    BCFauFc float64
    BCAxRatio float64
    BCAxAllow float64
    BCcheck float64
    TCIpMidPut float64
    TCIpMidPusol float64
    TCIpMidTenRat float64
    TCIpMidTenAllow float64
    TCIpMidFau float64
    TCIpMidFbu float64
    TCIpMidFauFc float64
    TCIpMidAxRatio float64
    TCIpMidAxAllow float64
    TCIpMidcheck float64
    TCIpPointPut float64
    TCIpPointPusol float64
    TCIpPointTenRat float64
    TCIpPointTenAllow float64
    TCIpPointFau float64
    TCIpPointFbu float64
    TCIpPointFauFc float64
    TCIpPointAxRatio float64
    TCIpPointAxAllow float64
    TCIpPointcheck float64
    TCEp1MidPut float64
    TCEp1MidPusol float64
    TCEp1MidTenRat float64
    TCEp1MidTenAllow float64
    TCEp1MidFau float64
    TCEp1MidFbu float64
    TCEp1MidFauFc float64
    TCEp1MidAxRatio float64
    TCEp1MidAxAllow float64
    TCEp1Midcheck float64
    TCEp1PointPut float64
    TCEp1PointPusol float64
    TCEp1PointTenRat float64
    TCEp1PointTenAllow float64
    TCEp1PointFau float64
    TCEp1PointFbu float64
    TCEp1PointFauFc float64
    TCEp1PointAxRatio float64
    TCEp1PointAxAllow float64
    TCEp1Pointcheck float64
    TCEp2MidPut float64
    TCEp2MidPusol float64
    TCEp2MidTenRat float64
    TCEp2MidTenAllow float64
    TCEp2MidFau float64
    TCEp2MidFbu float64
    TCEp2MidFauFc float64
    TCEp2MidAxRatio float64
    TCEp2MidAxAllow float64
    TCEp2Midcheck float64
    TCEp2PointPut float64
    TCEp2PointPusol float64
    TCEp2PointTenRat float64
    TCEp2PointTenAllow float64
    TCEp2PointFau float64
    TCEp2PointFbu float64
    TCEp2PointFauFc float64
    TCEp2PointAxRatio float64
    TCEp2PointAxAllow float64
    TCEp2Pointcheck float64
}

type CapSols = []CapSol
type TableCap struct {
    CapSols CapSols
}

func newCapSol() CapSol{
    return CapSol{}
}

func newTableCap() TableCap{
    return TableCap {
	CapSols: []CapSol{
	    newCapSol(),
	},
    }
}
type ShearCap struct {
    Vep float64
    Pu_Ep float64
    Fn float64
    Fv float64

    Botfv float64
    Botfa float64
    Botfvmod float64
    Botcheck string

    Topfv float64
    Topfa float64
    Topfvmod float64
    Topcheck string
}
type ShearCaps = []ShearCap
type TableShear struct {
    ShearCaps ShearCaps
}

func newShearCap()ShearCap{
	return ShearCap{}
}
func newTableShear() TableShear{
    return TableShear {
	ShearCaps: []ShearCap{
	    newShearCap(),
	},
    }
}

// Design of Web
type DOWSlendRat struct {
    InputName string
    ElemName string
    IX float64
    IY float64
    IZ float64
    Rx float64
    Ry float64
    Rz float64
    Ixrx float64
    Iyry float64
    Izrz float64
    SLRym float64
    Lrgov float64
    Limit float64
    Check string
}
type DOWSlendRats = []DOWSlendRat
type TableDOWSlend struct {
    DOWSlendRats DOWSlendRats
}

func newDOWSlendRat() DOWSlendRat{
	return DOWSlendRat{}
}
func newTableDOWSlend() TableDOWSlend{
    return TableDOWSlend{
	DOWSlendRats: []DOWSlendRat{
	    newDOWSlendRat(), 
	},
    }
}

type DOWDesign struct {
    InputName string
    ElemName string
    A float64
    TenFt float64
    TenPut float64
    TenPuSol float64
    TenRat float64
    TenAllow string
    Compbt float64
    CompQ float64
    CompFe float64
    CompFcr float64
    CompFc float64
    CompPuc float64
    CompPuSol float64
    CompRat float64
    CompAllow string
}
type DOWDesigns = []DOWDesign
type TableDOWDesign struct {
    DOWDesigns DOWDesigns
}

func newDOWDesign()DOWDesign{
    return DOWDesign{}
}
func newTableDOWDesign() TableDOWDesign{
    return TableDOWDesign{
	DOWDesigns: []DOWDesign{
	    newDOWDesign(), 
	},
    }
}

type DOWWeb struct {
    InputName string
    ElemName string
    XPE float64
    XPC float64
    EQV1 float64
    EQV2 float64
    Vmin float64
    FtMin float64
    FcMin float64
    LiftFTens float64
    LiftFComp float64
    DesignFTens float64
    DesignFComp float64
}
type DOWWebs = []DOWWeb
type TableDOWWeb struct {
    DOWWebs DOWWebs
}

func newDOWWeb()DOWWeb{
    return DOWWeb {}
}
func newTableDOWWeb() TableDOWWeb{
    return TableDOWWeb{
	DOWWebs: []DOWWeb{
	    newDOWWeb(), 
	},
    }
}

type DOWEfective struct {
    InputName string
    ElemName string
    Klrx float64
    Klry float64
    Klrz float64
    Klsrz float64
    SlendKlrx float64
    SlendKlry float64
    SlendKlrz float64
    SlendKlsrz float64
    SLRgov float64
}
type DOWEfectives =[]DOWEfective
type TableDOWEfective struct {
    DOWEfectives DOWEfectives
}

func newDOWEfective()DOWEfective{
    return DOWEfective{}
}
func newTableDOWEfective() TableDOWEfective{
    return TableDOWEfective{
	DOWEfectives: []DOWEfective{
	    newDOWEfective(), 
	},
    }
}

// forces and stresses in Web member
type DOWForce struct {
    Rmax float64
    Vmin float64
    Beta float64
    Gamma float64
    Delta float64
    Alpha float64
}
type DOWForces = []DOWForce
type TableDOWForce struct {
    DOWForces DOWForces
}

func newDOWForce()DOWForce{
    return DOWForce{}
}
func newTableDOWForce() TableDOWForce{
    return TableDOWForce{
	DOWForces: []DOWForce{
	    newDOWForce(), 
	},
    }
}

// Moment of Inertia
type M struct {
    v float64
}
type Ms = []M
type TMs struct {
    Ms Ms
}
func newM() M {
    return M{}
}
func newTMs() TMs{
    return TMs{
	Ms: []M{
	    newM(),
	},
    }
}

type Moment struct {
    Ijoist float64
    Ireq360 float64
    Ireq240 float64
    Check string
}
type Moments = []Moment
type TableMoment struct {
    Moments Moments
}

func newMoment() Moment{
    return Moment{}
}
func newTableMoment() TableMoment{
    return TableMoment{
	Moments: []Moment{
	    newMoment(), 
	},
    }
}

// Lateral Stability of Joist during erection
type Lateral struct {
    P float64
    K float64
    G float64
    J float64
    Y float64
    Iytc float64
    Iybc float64
    Iy float64
    Ix float64
    Y0 float64
    Cw float64
    Betax float64
    Ae float64
    A float64
    B float64
    C float64
    W1 float64
    W2 float64
    WuMin float64
    Wactual float64
    Wmax float64
}
type Laterals = []Lateral
type TableLateral struct {
    Laterals Laterals
}

func newLateral()Lateral{
    return Lateral{
    }
}
func newTableLateral() TableLateral{
    return TableLateral{
	Laterals: []Lateral{
	    newLateral(), 
	},
    }
}

// Design of Weld
type DesignWeld struct {
    Resistance float64
    Tensile float64
    Nominal float64
    Angle string
    Fillet float64
    WeldL float64
    Force float64
    Lw float64
}
type DesignWelds = []DesignWeld
type TableDesignWeld struct {
    DesignWelds DesignWelds
}

func newDesignWeld() DesignWeld{
    return DesignWeld{
    }
}
func newTableDesignWeld() TableDesignWeld{
    return TableDesignWeld{
	DesignWelds: []DesignWeld{
	    newDesignWeld(), 
	},
    }
}

type DesignWeldCon struct {
    InputName string
    ElemName string
    Conclusion string
    Ftens float64
    FComp float64
    F50 float64
    FGov float64
    Tw float64
    Lwmin float64
    Lw float64
}
type DesignWeldCons = []DesignWeldCon
type TableDesignWeldCon struct {
    DesignWeldCons DesignWeldCons
}

func newDesignWeldCon()DesignWeldCon{
    return DesignWeldCon{
    }
}
func newTableDesignWeldCon() TableDesignWeldCon{
    return TableDesignWeldCon{
	DesignWeldCons: []DesignWeldCon{
	    newDesignWeldCon(), 
	},
    }
}

func roundFloat(val float64, precision uint) float64 {
    ratio := math.Pow(10, float64(precision))
    return math.Round(val * ratio) / ratio
}

type Page struct {
    Geometry Geometry
    ResGeom ResGeom
    Material Material
    ResMater ResMater 
    ResBrid ResBrid
    WebMember WebMember
    TableResistance TableResistance

    TableCheck TableCheck

    TableEfective TableEfective
    TableDesign TableDesign
    TableResForce TableResForce

    TableCap TableCap
    TableShear TableShear
    TableDOWSlend TableDOWSlend
    TableDOWDesign TableDOWDesign


    TableDOWForce TableDOWForce
    TableDOWWeb TableDOWWeb
    TableDOWEfective TableDOWEfective
    TableLateral TableLateral
    TableDesignWeld TableDesignWeld
    TableDesignWeldCon TableDesignWeldCon

    TableMoment TableMoment
    TMs TMs
    TableAnalysis TableAnalysis
}

func newPage() Page {
    return Page {
	Geometry: newGeometry(),
	ResGeom: newResGeom(),
	Material: newMaterial(),
	ResMater: newResMater(),
	ResBrid: newResBrid(),
	WebMember: newWebMember(),
	TableResistance: newTableResistance(),
	TableCheck: newTableCheck(),
	TableEfective: newTableEfective(),
	TableDesign: newTableDesign(),
	TableResForce: newTableResForce(),
	TableCap: newTableCap(),
	TableShear: newTableShear(),
	TableDOWSlend: newTableDOWSlend(),
	TableDOWDesign: newTableDOWDesign(),
	TableDOWForce: newTableDOWForce(),
	TableDOWWeb: newTableDOWWeb(),
	TableDOWEfective: newTableDOWEfective(),
	TableLateral: newTableLateral(),
	TableDesignWeld: newTableDesignWeld(),
	TableDesignWeldCon: newTableDesignWeldCon(),
	TableMoment: newTableMoment(),
	TMs: newTMs(),
	TableAnalysis: newTableAnalysis(),
    }
}

func larger (array ...float64) float64{
    var holder float64
    for _, num := range array{
	if num > holder {
	    holder = num
	}
    }
    return holder
}

func largeThree(num1, num2, num3 float64) float64 {
    var larger float64
    if num1 >= num2 && num1 >= num3 {
	larger = num1
    } else if num2 >= num1 && num2 >= num3 {
	larger = num2
    } else {
	larger = num3
    }
    return larger
}
func checkerOK(rgov, limit *float64) string{
    if *rgov < *limit {
	return "OK"
    } else {
	return "NOT OK"
    }
}
func mmm(page *Page) {
     ms := &page.TMs.Ms[0]
     ms.v = math.Pi
}

func analysis(page *Page) {

    W := page.TableResForce.ResForces[0].KipUdlWu 
    lep1 := page.Geometry.Properties[0].fepl 
    lep2 := page.Geometry.Properties[0].sepl 

    rmax := page.TableDOWForce.DOWForces[0].Rmax
    betaR := page.TableDOWForce.DOWForces[0].Beta / (180/math.Pi)
    gammaR := page.TableDOWForce.DOWForces[0].Gamma / (180/math.Pi)

     
    forces := &page.TableAnalysis.AnalysisResults[0]    
    pChord := page.TableResForce.ResForces[0].ChordForce
    Atc := page.ResMater.AnglProps[0].AreaTop

    FcMin := W * ((lep1 / 2) + (lep2 / 2))
    minusPanel := math.Cos(gammaR) * FcMin

    forces.TopChrdCmpForceMid = RoundTo(pChord, 4)
    forces.TopChrdCmpFcMid = RoundTo(forces.TopChrdCmpForceMid / Atc, 4)

    EP1Force := RoundTo(rmax / math.Tan(betaR), 4)
    forces.TopChrdEP1CmpForceMid = EP1Force

    EP2Force := RoundTo(EP1Force - minusPanel, 4)
    forces.TopChrdEP2CmpForceMid = EP2Force

    forces.TopChrdEP1CmpFcMid  = RoundTo(EP1Force / Atc, 4)
    forces.TopChrdEP2CmpFcMid  = RoundTo(EP2Force / Atc, 4)
} 

func checkShear(page *Page) {
    shearCap := &page.TableShear.ShearCaps[0]
    rmax := page.TableDOWForce.DOWForces[0].Rmax
    W := page.TableResForce.ResForces[0].KipUdlWu 
    lep1 := page.Geometry.Properties[0].fepl 
    betaR := page.TableDOWForce.DOWForces[0].Beta / (180/math.Pi)
    ABot := page.ResMater.AnglProps[0].AreaBot
    bBot := page.ResMater.AnglProps[0].BBot
    tBot := page.ResMater.AnglProps[0].TBot

    ATop := page.ResMater.AnglProps[0].AreaTop
    bTop := page.ResMater.AnglProps[0].BTop
    tTop := page.ResMater.AnglProps[0].TTop


    vep := rmax - ((lep1 * W) / 2)
    puep := vep / math.Tan(betaR)
    shearCap.Vep = RoundTo(vep, 4)
    shearCap.Pu_Ep = RoundTo(puep, 4)
    Fy := page.Material.Forces[0].YieldStress

    shearCap.Fn = 0.6 * Fy
    shearCap.Fv = shearCap.Fn

    shearCap.Botfv = RoundTo(vep / (2 * bBot * tBot), 4)
    shearCap.Botfa = RoundTo(puep / ABot, 4)
    shearCap.Botfvmod = RoundTo(.5 * math.Sqrt(math.Pow(shearCap.Botfa, 2) + 4 * math.Pow(shearCap.Botfv, 2)), 4)
    if shearCap.Botfvmod < shearCap.Fv {
	shearCap.Botcheck = "Ok"
    } else {
	shearCap.Botcheck = "Not Ok"
    }

    shearCap.Topfv = RoundTo(vep / (2 * bTop * tTop), 4)
    shearCap.Topfa = RoundTo(puep / ATop, 4)
    shearCap.Topfvmod = RoundTo(.5 * math.Sqrt(math.Pow(shearCap.Topfa, 2) + 4 * math.Pow(shearCap.Topfv, 2)), 4)
    if shearCap.Topfvmod < shearCap.Fv {
	shearCap.Topcheck = "Ok"
    } else {
	shearCap.Topcheck = "Not Ok"
    }
}

func slendernesRadio(page *Page) {
    geom := &page.Geometry.Properties[0]
    resProp := &page.ResGeom.ResProps[0]
    check := &page.TableCheck.CheckSlenderness[0]
    bot := &page.ResMater.AnglProps[0]

    check.BCIx = geom.ipl*2
    check.BCIy = (geom.span*12)/4
    check.BCIz = check.BCIx
    check.BCrx = bot.RxBot
    check.BCry = bot.RyBot
    check.BCrz = bot.RzBot
    check.BCIxrx = roundFloat(check.BCIx / check.BCrx, 4)
    check.BCIyry = roundFloat(check.BCIy / check.BCry, 4)
    check.BCIzrz = roundFloat(check.BCIz / check.BCrz, 4)
    check.BCrgov = largeThree(check.BCIxrx, check.BCIyry, check.BCIzrz)
    check.BClimit = 240
    check.BCcheck = checkerOK(&check.BCrgov, &check.BClimit)
    
    check.TCIpMidIx = geom.ipl
    check.TCIpMidIy = math.Floor(resProp.Lbe2)
    check.TCIpMidIz = geom.ipl/2
    check.TCIpMidrx = bot.RxTop
    check.TCIpMidry = bot.RyTop
    check.TCIpMidrz = bot.RzTop
    check.TCIpMidIxrx = roundFloat(check.TCIpMidIx / check.TCIpMidrx, 4)
    check.TCIpMidIyry = roundFloat(check.TCIpMidIy / check.TCIpMidry, 4)
    check.TCIpMidIzrz = roundFloat(check.TCIpMidIz / check.TCIpMidrz, 4)
    check.TCIpMidrgov = larger(check.TCIpMidIxrx, check.TCIpMidIyry, check.TCIpMidIzrz)
    check.TCIpMidlimit = 90
    check.TCIpMidcheck = checkerOK(&check.TCIpMidrgov, &check.TCIpMidlimit)

    check.TCIpPointIx = geom.ipl
    check.TCIpPointIy = math.Floor(resProp.Lbe2)
    check.TCIpPointIz = geom.ipl/2
    check.TCIpPointrx = bot.RxTop
    check.TCIpPointry = bot.RyTop
    check.TCIpPointrz = bot.RzTop
    check.TCIpPointIxrx = roundFloat(check.TCIpPointIx / check.TCIpPointrx, 4)
    check.TCIpPointIyry = roundFloat(check.TCIpPointIy / check.TCIpPointry, 4)
    check.TCIpPointIzrz = roundFloat(check.TCIpPointIz / check.TCIpPointrz, 4)
    check.TCIpPointrgov = larger(check.TCIpPointIxrx, check.TCIpPointIyry, check.TCIpPointIzrz)
    check.TCIpPointlimit = 90
    check.TCIpPointcheck = checkerOK(&check.TCIpPointrgov, &check.TCIpPointlimit)


    check.TCEp1MidIx = geom.fepl
    check.TCEp1MidIy = math.Floor(resProp.Lbe2)
    check.TCEp1MidIz = geom.fepl
    check.TCEp1Midrx = bot.RxTop
    check.TCEp1Midry = bot.RyTop
    check.TCEp1Midrz = bot.RzTop
    check.TCEp1MidIxrx = roundFloat(check.TCEp1MidIx / check.TCEp1Midrx, 4)
    check.TCEp1MidIyry = roundFloat(check.TCEp1MidIy / check.TCEp1Midry, 4)
    check.TCEp1MidIzrz = roundFloat(check.TCEp1MidIz / check.TCEp1Midrz, 4)
    check.TCEp1Midrgov = larger(check.TCEp1MidIxrx, check.TCEp1MidIyry, check.TCEp1MidIzrz)
    check.TCEp1Midlimit = 120
    check.TCEp1Midcheck = checkerOK(&check.TCEp1Midrgov, &check.TCEp1Midlimit)

    check.TCEp1PointIx = geom.fepl
    check.TCEp1PointIy = math.Floor(resProp.Lbe2)
    check.TCEp1PointIz = geom.fepl
    check.TCEp1Pointrx = bot.RxTop
    check.TCEp1Pointry = bot.RyTop
    check.TCEp1Pointrz = bot.RzTop
    check.TCEp1PointIxrx = roundFloat(check.TCEp1PointIx / check.TCEp1Pointrx, 4)
    check.TCEp1PointIyry = roundFloat(check.TCEp1PointIy / check.TCEp1Pointry, 4)
    check.TCEp1PointIzrz = roundFloat(check.TCEp1PointIz / check.TCEp1Pointrz, 4)
    check.TCEp1Pointrgov = larger(check.TCEp1PointIxrx, check.TCEp1PointIyry, check.TCEp1PointIzrz)
    check.TCEp1Pointlimit = 120
    check.TCEp1Pointcheck = checkerOK(&check.TCEp1Pointrgov, &check.TCEp1Pointlimit)

    check.TCEp2MidIx = geom.sepl
    check.TCEp2MidIy = math.Floor(resProp.Lbe2)
    check.TCEp2MidIz = geom.sepl
    check.TCEp2Midrx = bot.RxTop
    check.TCEp2Midry = bot.RyTop
    check.TCEp2Midrz = bot.RzTop
    check.TCEp2MidIxrx = roundFloat(check.TCEp2MidIx / check.TCEp2Midrx, 4)
    check.TCEp2MidIyry = roundFloat(check.TCEp2MidIy / check.TCEp2Midry, 4)
    check.TCEp2MidIzrz = roundFloat(check.TCEp2MidIz / check.TCEp2Midrz, 4)
    check.TCEp2Midrgov = larger(check.TCEp2MidIxrx, check.TCEp2MidIyry, check.TCEp2MidIzrz)
    check.TCEp2Midlimit = 120
    check.TCEp2Midcheck = checkerOK(&check.TCEp2Midrgov, &check.TCEp2Midlimit)

    check.TCEp2PointIx = geom.sepl
    check.TCEp2PointIy = math.Floor(resProp.Lbe2)
    check.TCEp2PointIz = geom.sepl
    check.TCEp2Pointrx = bot.RxTop
    check.TCEp2Pointry = bot.RyTop
    check.TCEp2Pointrz = bot.RzTop
    check.TCEp2PointIxrx = roundFloat(check.TCEp2PointIx / check.TCEp2Pointrx, 4)
    check.TCEp2PointIyry = roundFloat(check.TCEp2PointIy / check.TCEp2Pointry, 4)
    check.TCEp2PointIzrz = roundFloat(check.TCEp2PointIz / check.TCEp2Pointrz, 4)
    check.TCEp2Pointrgov = larger(check.TCEp2PointIxrx, check.TCEp2PointIyry, check.TCEp2PointIzrz)
    check.TCEp2Pointlimit = 120
    check.TCEp2Pointcheck = checkerOK(&check.TCEp2Pointrgov, &check.TCEp2Pointlimit)
}
func deCalculation(page *Page) {
    d := &page.Geometry.Properties[0].depth
    ytc := &page.ResMater.AnglProps[0].YTop
    ybc := &page.ResMater.AnglProps[0].YBot
    de := (*d - *ytc) - *ybc
    page.ResGeom.ResProps[0].Ed = de
    // fmt.Println(de)
}
func momentOfInertia(page *Page) {
    Moment := &page.TableMoment.Moments[0]
    Ixtc := page.ResMater.AnglProps[0].IxTop
    Ixbc := page.ResMater.AnglProps[0].IxBot
    Atc := page.ResMater.AnglProps[0].AreaTop
    Abc := page.ResMater.AnglProps[0].AreaBot
    de := page.ResGeom.ResProps[0].Ed
    Wll := page.TableResForce.ResForces[0].KipLlWLL
    E := page.Material.Forces[0].ModElas
    // check := &page.TableMoment.Moments[0].check
    L := page.ResGeom.ResProps[0].DLength

    ijoist := Ixtc + Ixbc + ((Atc * Abc * math.Pow(de, 2)) / (Atc + Abc))
    ireq360 := 1.15 * 5 * 360 * Wll * math.Pow(L, 3) / (384 * E)
    ireq240 := ireq360 * 2 / 3

    Moment.Ijoist = RoundTo(ijoist, 2)
    Moment.Ireq360 = RoundTo(ireq360, 2)
    Moment.Ireq240 = RoundTo(ireq240, 2)

    if ireq240 < ijoist {
	Moment.Check = "OK"
    } else {
	Moment.Check ="NOT OK"
    }
}
func efectiveSlend(page *Page) {
    // if Top chord (Lip max) mid panel fill yes check criteria B in the table
    EfectiveSlend := &page.TableEfective.EfectiveSlenderness[0]

    fillTopChord := &page.Material.Forces[0].FillTopChord
    fillBotChord := &page.Material.Forces[0].FillBotChord
    TCEP1MP := &page.Material.Forces[0].TopChordEP1
    TCEP2MP := &page.Material.Forces[0].TopChordEP2
    CheckSlend := &page.TableCheck.CheckSlenderness[0]

	//EfectiveSlend.BCklrz = 1.0

    // T.C. (Ip) mid-panel
    if *fillTopChord {
	EfectiveSlend.TCIpMidklrx = 0.75
	EfectiveSlend.TCIpMidklry = 0.94
	EfectiveSlend.TCIpMidklrz = 0
	EfectiveSlend.TCIpMidklsrz = 1
    } else {
	EfectiveSlend.TCIpMidklrx = 0
	EfectiveSlend.TCIpMidklry = 0
	EfectiveSlend.TCIpMidklrz = 0.75
	EfectiveSlend.TCIpMidklsrz = 0
    }
    if *fillBotChord {
	EfectiveSlend.BCklrx = 0.9
	EfectiveSlend.BCklry = 0.94
	EfectiveSlend.BCklrz = 0
	EfectiveSlend.BCklsrz = 1
    } else {
	EfectiveSlend.BCklrx = 0
	EfectiveSlend.BCklry = 0
	EfectiveSlend.BCklrz = 1
	EfectiveSlend.BCklsrz = 0
    }
    // T.C. (Ep1) mid-panel
    if *TCEP1MP{
	EfectiveSlend.TCEp1Midklrx = 1.0
	EfectiveSlend.TCEp1Midklry = 0.94
	EfectiveSlend.TCEp1Midklrz = 0
	EfectiveSlend.TCEp1Midklsrz = 1
    } else {
	EfectiveSlend.TCEp1Midklrx = 0
	EfectiveSlend.TCEp1Midklry = 0
	EfectiveSlend.TCEp1Midklrz = 1
	EfectiveSlend.TCEp1Midklsrz = 0
    }
    if *TCEP2MP{
	EfectiveSlend.TCEp2Midklrx = 1.0
	EfectiveSlend.TCEp2Midklry = 0.94
	EfectiveSlend.TCEp2Midklrz = 0
	EfectiveSlend.TCEp2Midklsrz = 1
    } else {
	EfectiveSlend.TCEp2Midklrx = 0
	EfectiveSlend.TCEp2Midklry = 0
	EfectiveSlend.TCEp2Midklrz = 1
	EfectiveSlend.TCEp2Midklsrz = 0
    }

    EfectiveSlend.BCSlendklrx = EfectiveSlend.BCklrx * CheckSlend.BCIxrx
    EfectiveSlend.BCSlendklry = EfectiveSlend.BCklry * CheckSlend.BCIyry
    EfectiveSlend.BCSlendklrz = EfectiveSlend.BCklrz * CheckSlend.BCIzrz
    EfectiveSlend.BCSlendklsrz = EfectiveSlend.BCklsrz * CheckSlend.BCIzrz

    EfectiveSlend.TCIpMidSlendklrx = EfectiveSlend.TCIpMidklrx * CheckSlend.TCIpMidIxrx
    EfectiveSlend.TCIpMidSlendklry = EfectiveSlend.TCIpMidklry * CheckSlend.TCIpMidIyry
    EfectiveSlend.TCIpMidSlendklrz = EfectiveSlend.TCIpMidklrz * CheckSlend.TCIpMidIzrz
    EfectiveSlend.TCIpMidSlendklrz = EfectiveSlend.TCIpMidklsrz * CheckSlend.TCIpMidIzrz

    // SLR Ggov
    EfectiveSlend.BCSLRgov = larger(EfectiveSlend.BCSlendklrx, EfectiveSlend.BCSlendklry, EfectiveSlend.BCSlendklrz)

    EfectiveSlend.TCIpMidSLRgov = larger(EfectiveSlend.TCIpMidSlendklrx, EfectiveSlend.TCIpMidSlendklry, EfectiveSlend.TCIpMidSlendklrz) 
    EfectiveSlend.TCEp1MidSlendklrx = EfectiveSlend.TCEp1Midklrx * CheckSlend.TCEp1MidIxrx 
    EfectiveSlend.TCEp1MidSlendklry = EfectiveSlend.TCEp1Midklry * CheckSlend.TCEp1MidIyry
    EfectiveSlend.TCEp1MidSlendklrz = EfectiveSlend.TCEp1Midklrz * CheckSlend.TCEp1MidIzrz

    EfectiveSlend.TCEp1MidSLRgov = larger(EfectiveSlend.TCEp1MidSlendklrx, EfectiveSlend.TCEp1MidSlendklry, EfectiveSlend.TCEp1MidSlendklrz)
    
    EfectiveSlend.TCEp2MidSlendklrx = EfectiveSlend.TCEp2Midklrx * CheckSlend.TCEp2MidIxrx
    EfectiveSlend.TCEp2MidSlendklry = EfectiveSlend.TCEp2Midklry * CheckSlend.TCEp2MidIyry
    EfectiveSlend.TCEp2MidSlendklrz = EfectiveSlend.TCEp2Midklrz * CheckSlend.TCEp2MidIzrz

    EfectiveSlend.TCEp2MidSLRgov = larger(EfectiveSlend.TCEp2MidSlendklrx, EfectiveSlend.TCEp2MidSlendklry, EfectiveSlend.TCEp2MidSlendklrz)


    // change this values. Probably not constants
    EfectiveSlend.BCKLrx = 0.75
    EfectiveSlend.TCIpMidKLrx = 0.75
    EfectiveSlend.TCIpPointKLrx = 0.75
    EfectiveSlend.TCEp1MidKLrx = 1
    EfectiveSlend.TCEp1PointKLrx = 1
    EfectiveSlend.TCEp2MidKLrx = 1
    EfectiveSlend.TCEp2PointKLrx = 1

    EfectiveSlend.BCFeKlrx = EfectiveSlend.BCKLrx * CheckSlend.BCIxrx
    EfectiveSlend.TCIpMidFeKlrx = EfectiveSlend.TCIpMidKLrx * CheckSlend.TCIpMidIxrx
    EfectiveSlend.TCIpPointFeKlrx = EfectiveSlend.TCIpPointKLrx * CheckSlend.TCIpPointIxrx
    EfectiveSlend.TCEp1MidFeKlrx = EfectiveSlend.TCEp1MidKLrx * CheckSlend.TCEp1MidIxrx
    EfectiveSlend.TCEp1PointFeKlrx = EfectiveSlend.TCEp1PointKLrx * CheckSlend.TCEp1PointIxrx
    EfectiveSlend.TCEp2MidFeKlrx = EfectiveSlend.TCEp2MidKLrx * CheckSlend.TCEp2MidIxrx
    EfectiveSlend.TCEp2PointFeKlrx = EfectiveSlend.TCEp2PointKLrx * CheckSlend.TCEp2PointIxrx
 
}
func design (page *Page) {
    C := page.TableResistance.ResistanceFactors[0].CompressionValue
    B := page.TableResistance.ResistanceFactors[0].BendingValue
    Fy := page.Material.Forces[0].YieldStress
    dStress := &page.TableDesign.Designs[0]
    Ft := C * Fy
    Fb := B * Fy
    E := page.Material.Forces[0].ModElas
    dStress.BCFt = Ft
    dStress.BCFb = Fb
    Abc := page.ResMater.AnglProps[0].AreaBot
    Atc := page.ResMater.AnglProps[0].AreaTop
    Btc := page.ResMater.AnglProps[0].BTop
    Ttc := page.ResMater.AnglProps[0].TTop
    Qtc := page.ResMater.AnglProps[0].QTop
    Bbc := page.ResMater.AnglProps[0].BBot
    Tbc := page.ResMater.AnglProps[0].TBot
    Qbc := page.ResMater.AnglProps[0].QBot
    Klrxbc := page.TableEfective.EfectiveSlenderness[0].BCFeKlrx

    Klrxtc := page.TableEfective.EfectiveSlenderness[0].TCIpMidFeKlrx
    KlrxEp1 := page.TableEfective.EfectiveSlenderness[0].TCEp1MidFeKlrx
    KlrxEp2 := page.TableEfective.EfectiveSlenderness[0].TCEp2MidFeKlrx

    SLRbc := page.TableEfective.EfectiveSlenderness[0].BCSLRgov
    SLRtc := page.TableEfective.EfectiveSlenderness[0].TCIpMidSLRgov
    SLREp1 := page.TableEfective.EfectiveSlenderness[0].TCEp1MidSLRgov
    SLREp2 := page.TableEfective.EfectiveSlenderness[0].TCEp2MidSLRgov


    atc := math.Sqrt(E / Qtc * Fy)
    abc := math.Sqrt(E / Qbc * Fy)

    dStress.BCPut = Ft * page.ResMater.AnglProps[0].AreaBot
    dStress.BCbt = RoundTo(Bbc / Tbc, 4)
    dStress.BCQ = Qbc
    dStress.BCFex = RoundTo((math.Pow(math.Pi, 2) * E) / math.Pow(Klrxbc, 2), 4)
    fmt.Println(SLRbc)

    dStress.TCIpMidPut = Ft * page.ResMater.AnglProps[0].AreaTop
    dStress.TCIpMidbt = RoundTo(Btc / Ttc , 4)
    dStress.TCIpMidQ = Qtc

    dStress.TCIpMidFex = RoundTo((math.Pow(math.Pi, 2) * E) / math.Pow(Klrxtc, 2), 4)
    dStress.TCEp1MidFex = RoundTo((math.Pow(math.Pi, 2) * E) / math.Pow(KlrxEp1, 2), 4)
    dStress.TCEp2MidFex = RoundTo((math.Pow(math.Pi, 2) * E) / math.Pow(KlrxEp2, 2), 4)

    bcFetc := RoundTo((math.Pow(math.Pi, 2) * E) / math.Pow(SLRbc, 2), 4)
    dStress.BCFetc = bcFetc
    tcIpFetc := RoundTo((math.Pow(math.Pi, 2) * E) / math.Pow(SLRtc, 2), 4)
    dStress.TCIpMidFetc = tcIpFetc
    tcEp1Fetc := RoundTo((math.Pow(math.Pi, 2) * E) / math.Pow(SLREp1, 2), 4)
    dStress.TCEp1MidFetc = tcEp1Fetc
    tcEp2Fetc := RoundTo((math.Pow(math.Pi, 2) * E) / math.Pow(SLREp2, 2), 4)
    dStress.TCEp2MidFetc = tcEp2Fetc

    expbc := (Qbc * Fy) / bcFetc
    exptc := (Qtc * Fy) / tcIpFetc
    expEp1 := (Qtc * Fy) / tcEp1Fetc
    expEp2 := (Qtc * Fy) / tcEp2Fetc


    if bcFetc <= abc {
	dStress.BCFcr = RoundTo(.877 * bcFetc, 4)
    } else {
	dStress.BCFcr = RoundTo(Qbc * math.Pow(0.658, expbc) * Fy, 4)
    }
    if tcIpFetc <= atc {
	dStress.TCIpMidFcr = RoundTo(Qtc * (math.Pow(0.658, exptc)) * Fy, 4)
    } else {
	dStress.TCIpMidFcr = .877 * tcIpFetc
    }
    if tcEp1Fetc <= atc {
	dStress.TCEp1MidFcr = RoundTo(Qtc * (math.Pow(0.658, expEp1)) * Fy, 4)
    } else {
	dStress.TCEp1MidFcr = .877 * tcEp1Fetc
    }
    if tcEp2Fetc <= atc {
	dStress.TCEp2MidFcr = RoundTo(Qtc * (math.Pow(0.658, expEp2)) * Fy, 4)
    } else {
	dStress.TCEp2MidFcr = .877 * tcEp2Fetc
    }

    dStress.BCFc = RoundTo(C * dStress.BCFcr, 4)
    dStress.TCIpMidFc = RoundTo(C * dStress.TCIpMidFcr, 4)
    dStress.TCEp1MidFc = RoundTo(C * dStress.TCEp1MidFcr, 4)
    dStress.TCEp2MidFc = RoundTo(C * dStress.TCEp2MidFcr, 4)

    dStress.BCPuc = RoundTo(dStress.BCFc * Abc, 4)
    dStress.TCIpMidPuc = RoundTo(dStress.TCIpMidFc * Atc, 4)
    dStress.TCEp1MidPuc = RoundTo(dStress.TCEp1MidFc * Atc, 4)
    dStress.TCEp2MidPuc = RoundTo(dStress.TCEp2MidFc * Atc, 4)

    forces := &page.TableAnalysis.AnalysisResults[0]    
    if forces.BotChrdCmpFc == 0 {
	dStress.BCcm = 1
    }
    Fc := forces.TopChrdCmpFcMid

    Ep1Fc := forces.TopChrdEP1CmpFcMid
    Ep2Fc := forces.TopChrdEP2CmpFcMid
    Ipcm := RoundTo(1 - ((0.4 * Fc) / (C * dStress.TCIpMidFex)), 4)
    dStress.TCIpMidcm = Ipcm

    Ep1cm := RoundTo(1 - (0.3 * Ep1Fc) / (C * dStress.TCEp1MidFex), 4)
    dStress.TCEp1Midcm = Ep1cm
    Ep2cm := RoundTo(1 - (0.3 * Ep2Fc) / (C * dStress.TCEp2MidFex), 4)
    dStress.TCEp2Midcm = Ep2cm
}

func designLoads(page *Page) {
    udlw := &page.Material.Forces[0].Udlw
    llWLL := &page.Material.Forces[0].LlWLL
    L := &page.ResGeom.ResProps[0].DLength
    de := &page.ResGeom.ResProps[0].Ed 
    // divide by 12000 is the same than divide by 1000 and then by 12
    Wu := *udlw/12000
    Wll := *llWLL/12000

    tst := math.Pow(*L, 2) 
    Msji := (Wu * tst)/8
    Pchord := Msji / *de

    page.TableResForce.ResForces[0].KipUdlWu = Wu
    page.TableResForce.ResForces[0].KipLlWLL = Wll
    page.TableResForce.ResForces[0].MaxDsgnMoment = Msji
    page.TableResForce.ResForces[0].ChordForce = Pchord
}
func webMemberAngles(page *Page) {
    de := &page.ResGeom.ResProps[0].Ed
    lbe := &page.Geometry.Properties[0].epbc
    lpe1 := &page.Geometry.Properties[0].fepl
    lpe2 := &page.Geometry.Properties[0].sepl
    lip := &page.Geometry.Properties[0].ipl
    W := page.TableResForce.ResForces[0].KipUdlWu 
    L := page.ResGeom.ResProps[0].DLength
    
    beta := math.Atan(*de / *lbe) * (180 / math.Pi)
    gamma := math.Atan(*de / (*lbe - *lpe1)) * (180 / math.Pi)
    delta := math.Atan(*de / (*lpe1 + *lpe2 - *lbe)) * (180 / math.Pi)
    alpha := math.Atan(*de / *lip) * (180 / math.Pi)
    rmax := (W * L) / 2
    vmin := rmax * .25

    page.TableDOWForce.DOWForces[0].Beta = beta
    page.TableDOWForce.DOWForces[0].Gamma = gamma
    page.TableDOWForce.DOWForces[0].Delta = delta
    page.TableDOWForce.DOWForces[0].Alpha = alpha
    page.TableDOWForce.DOWForces[0].Rmax = rmax
    page.TableDOWForce.DOWForces[0].Vmin = vmin

} 

func angleProp(col, row int) string {
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
    return record[col][row]
}

func main() {
    t := &Template{
	templates: template.Must(template.ParseGlob("views/*.html")),
    }
    e := echo.New()
    page := newPage()

    // fmt.Printf("%+v", page)

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
    e.GET("/icon", func(c echo.Context) error {
	return c.File("static/favicon.ico")
    })

    e.POST("/sendform", func(c echo.Context) error {
	page.WebMember.MemberInputs = page.WebMember.MemberInputs[:0]
	page.TableDOWSlend.DOWSlendRats = page.TableDOWSlend.DOWSlendRats[:0]
	page.TableDOWDesign.DOWDesigns = page.TableDOWDesign.DOWDesigns[:0]
	page.TableDOWWeb.DOWWebs = page.TableDOWWeb.DOWWebs[:0]
	page.TableDOWEfective.DOWEfectives = page.TableDOWEfective.DOWEfectives[:0]
	page.TableDesignWeldCon.DesignWeldCons = page.TableDesignWeldCon.DesignWeldCons[:0]
	// Geometry input table
	trussType := c.FormValue("trussType")
	joistType := c.FormValue("joistType")
	deflexion := c.FormValue("deflexion")
	span := c.FormValue("span")
	depth := c.FormValue("depth")
	fepl := c.FormValue("fepl")
	sepl := c.FormValue("sepl")
	ipl := c.FormValue("ipl")
	epbc := c.FormValue("epbc")

	// Stress input table
	yieldStress := c.FormValue("yieldStress")
	weight := c.FormValue("weight")
	bSeat := c.FormValue("bSeat")
	modElas := c.FormValue("modElas")
	spaceChord := c.FormValue("spaceChord")

	fmt.Println("Stress ", yieldStress, modElas)

	// Chords input
	topChord := c.FormValue("topChord")
	bottomChord := c.FormValue("bottomChord")
	fillTopChord := c.FormValue("fillTopChord")
	fillBotChord := c.FormValue("fillBotChord")
	topChordEP1 := c.FormValue("topChordEP1")
	topChordEP2 := c.FormValue("topChordEP2")

	// Loads
	udlw := c.FormValue("udlw")
	llWLL := c.FormValue("llWLL")
	nsLRFD := c.FormValue("nsLRFD")
	// nsASD := c.FormValue("nsASD")

	/*
	// Chords input table
	topChord := c.FormValue("topChord")
	bottomChord := c.FormValue("bottomChord")
	topChordEP1 := c.FormValue("topChordEP1")
	topChordEP2 := c.FormValue("topChordEP2")
	fillBottomChord := c.FormValue("fillBottomChord")
	*/
	propertie := newPropertie(
	    trussType,
	    joistType,
	    deflexion,
	    span,
	    fepl,
	    sepl,
	    ipl,
	    epbc,
	    depth,
	)
	page.Geometry.Properties = append(page.Geometry.Properties, propertie)
	page.Geometry.Properties = page.Geometry.Properties[1:]
	// change name of lbe2
	lbe2 := bChordtobPanel(
	    page.Geometry.Properties[0].fepl,
	    page.Geometry.Properties[0].sepl,
	    page.Geometry.Properties[0].ipl,
	    page.Geometry.Properties[0].epbc,
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
	page.ResGeom.ResProps[0].Lbe2 = lbe2
	page.ResGeom.ResProps[0].DLength = dLength
	page.ResGeom.ResProps[0].Tip = tip
	page.ResGeom.ResProps[0].Tod = tod
	page.ResGeom.ResProps[0].Ts = ts
	page.ResGeom.ResProps[0].Dmin = dmin


	force := newForce(
	    yieldStress,
	    modElas,
	    spaceChord,
	    weight,
	    bSeat,
	    topChord,
	    bottomChord,
	    fillTopChord,
	    fillBotChord,
	    udlw,
	    llWLL,
	    nsLRFD,
	    topChordEP1,
	    topChordEP2,
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

	halfTod := tod/2
	halfTs := ts/2

	totalAngles := halfTod +1 + halfTs

	elementSec := newMemberInput()
	elementSlend := newDOWSlendRat()
	elementDesign := newDOWDesign()
	elementWeb := newDOWWeb()
	elementEfective := newDOWEfective()
	elementDesignWeldCon := newDesignWeldCon()

	tas := int(totalAngles)
	hts := int(halfTod)

	for i := 1; i < tas; i++{
	    elementSec.InputName = fmt.Sprintf("%v", i)
	    elementSlend.InputName = fmt.Sprintf("%v", i)
	    elementDesign.InputName = fmt.Sprintf("%v", i)
	    elementWeb.InputName = fmt.Sprintf("%v", i)
	    elementEfective.InputName = fmt.Sprintf("%v", i)
	    elementDesignWeldCon.InputName = fmt.Sprintf("%v", i)
	    if i == 1{
		elementSec.ElemName = "sv"
		elementSlend.ElemName = "sv"
		elementDesign.ElemName = "sv"
		elementWeb.ElemName = "sv"
		elementEfective.ElemName = "sv"
		elementDesignWeldCon.ElemName = "sv"
	    } else if i <= hts+1 && i != 1{
		elementSec.ElemName = fmt.Sprintf("w%v", i)
		elementSlend.ElemName = fmt.Sprintf("w%v", i)
		elementDesign.ElemName = fmt.Sprintf("w%v", i)
		elementWeb.ElemName = fmt.Sprintf("w%v", i)
		elementEfective.ElemName = fmt.Sprintf("w%v", i)
		elementDesignWeldCon.ElemName = fmt.Sprintf("w%v", i)
	    } else {
		elementSec.ElemName = fmt.Sprintf("v%v", i - (hts+1))
		elementSlend.ElemName = fmt.Sprintf("v%v", i - (hts+1))
		elementDesign.ElemName = fmt.Sprintf("v%v", i - (hts+1))
		elementWeb.ElemName = fmt.Sprintf("v%v", i - (hts+1))
		elementEfective.ElemName = fmt.Sprintf("v%v", i - (hts+1))
		elementDesignWeldCon.ElemName = fmt.Sprintf("v%v", i - (hts+1))
	    }
	    page.WebMember.MemberInputs = append(page.WebMember.MemberInputs, elementSec)
	    page.TableDOWSlend.DOWSlendRats = append(page.TableDOWSlend.DOWSlendRats, elementSlend)
	    page.TableDOWDesign.DOWDesigns = append(page.TableDOWDesign.DOWDesigns, elementDesign)
	    page.TableDOWWeb.DOWWebs = append(page.TableDOWWeb.DOWWebs, elementWeb)
	    page.TableDOWEfective.DOWEfectives = append(page.TableDOWEfective.DOWEfectives, elementEfective)
	    page.TableDesignWeldCon.DesignWeldCons = append(page.TableDesignWeldCon.DesignWeldCons, elementDesignWeldCon)
	}

	slendernesRadio(&page)
	efectiveSlend(&page)
	deCalculation(&page)
	designLoads(&page)
	webMemberAngles(&page)
	checkShear(&page)

	fmt.Println(page.TableMoment.Moments)
	momentOfInertia(&page)
	fmt.Println(page.TableMoment.Moments)

	analysis(&page)
	design(&page)


	//fmt.Println()
	//fmt.Printf("%+v", page)
	//fmt.Println()
	return c.Render(http.StatusOK, "res", page)
    })

    e.POST("/member", func(c echo.Context) error {
	MemberInputs := page.WebMember.MemberInputs
	//angleProps()
	/* TOOD: 
		make function to get angle props
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

	for i := 0; i < len(page.WebMember.MemberInputs); i++ {
	    si := strconv.Itoa(i + 1)

	    part := "part" + si
	    mark := "mark" + si
	    crimped := "crimped" + si
	    fill := "fill" + si

	    MemberInputs[i].Part = c.FormValue(part)
	    MemberInputs[i].Mark = c.FormValue(mark)
	    MemberInputs[i].Crimped = c.FormValue(crimped)
	    MemberInputs[i].Fill = c.FormValue(fill)

	    mmark,_ := strconv.Atoi(MemberInputs[i].Mark)

	    MemberInputs[i].Secction = record[mmark][2]

	    anglePropsDOWSlend(mmark, i, MemberInputs[i].Crimped, &page)

	    efectiveSlendernes(i, &page)

	    designWeb(mmark, i, &page)

	    designStress(mmark, i, &page)

	}

	return c.Render(http.StatusOK, "webMem", page )
    })
    e.Logger.Fatal(e.Start(":8080"))
}
