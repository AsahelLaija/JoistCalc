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
type ResProp struct {
    Lbe2, DLength, Tip, Tod, Ts, Ed, Dmin, Lbrdng1, Lbrdng2 float64
}
type ResProps = []ResProp
type ResGeom struct {
    ResProps ResProps
}
func newResProp(lb, dl, ti, to, ts, ed, Dmin float64) ResProp {
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
func newResGeom() ResGeom {
    return ResGeom{
	ResProps: []ResProp{
	    newResProp(0, 0, 0, 0, 0, 0, 0),
	},
    }
}

type Force struct {
    YieldStress, ModElas, SpaceChord, Weight, BSeat, TopChord, BottomChord, Udlw, LlWLL, NsLRFD float64
    FillTopChord, TopChordEP1, TopChordEP2 bool
}
type Forces = []Force
type Material struct {
    Forces Forces
}
func newForce(yi, mo, sp, we, bs, to, bo, tf, udlwA, llWLLA, nsLRFDA, tcep1, tcep2 string) Force {
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
    var topChordEP1 bool
    var topChordEP2 bool
    if tf == "yes" {
	fillTopChord = true
    } else {
	fillTopChord = false
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
	    newForce("50000", "29000", "1", "", "", "", "", "", "", "", "", "", ""),
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
func newResForce(kipUdlWu, kipLlWLL, kipNsLRFD, maxDsgnMoment, chordForce float64) ResForce{
    return ResForce{
	KipUdlWu: kipUdlWu,
	KipLlWLL: kipLlWLL,
	KipNsLRFD: kipNsLRFD,
	MaxDsgnMoment: maxDsgnMoment,
	ChordForce: chordForce,
    }
}
func newTableResForce() TableResForce {
    return TableResForce {
	ResForces: []ResForce{
	    newResForce(0, 0, 0, 0, 0),
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
type AnalysisResult struct {
    BotChrdTnsForce float64
    BotChrdTnsFt float64
    BotChrdCmpForce float64
    BotChrdMomDwnMu float64
    BotChrdMomDwnFb float64
    BotChrdMomUpMu float64
    BotChrdMomUpfb float64

    TopChrdTnsForceMid float64
    TopChrdTnsFtMid float64
    TopChrdCmpForceMid float64
    TopChrdTnsFcMid float64
    TopChrdMomDwnMuMid float64
    TopChrdMomDwnFbMid float64
    TopChrdMomUpMuMid float64
    TopChrdMomUpfbMid float64

    TopChrdTnsForcePoint float64
    TopChrdTnsFtPoint float64
    TopChrdCmpForcePoint float64
    TopChrdTnsFcPoint float64
    TopChrdMomDwnMuPoint float64
    TopChrdMomDwnFbPoint float64
    TopChrdMomUpMuPoint float64
    TopChrdMomUpfbPoint float64

    TopChrdEP1TnsForceMid float64 
    TopChrdEP1TnsFtMid float64    
    TopChrdEP1CmpForceMid float64 
    TopChrdEP1TnsFcMid float64    
    TopChrdEP1MomDwnMuMid float64 
    TopChrdEP1MomDwnFbMid float64 
    TopChrdEP1MomUpMuMid float64  
    TopChrdEP1MomUpfbMid float64  

    TopChrdEP1TnsForcePoint float64 
    TopChrdEP1TnsFtPoint float64    
    TopChrdEP1CmpForcePoint float64 
    TopChrdEP1TnsFcPoint float64    
    TopChrdEP1MomDwnMuPoint float64 
    TopChrdEP1MomDwnFbPoint float64 
    TopChrdEP1MomUpMuPoint float64  
    TopChrdEP1MomUpfbPoint float64  

    TopChrdEP2TnsForceMid float64 
    TopChrdEP2TnsFtMid float64    
    TopChrdEP2CmpForceMid float64 
    TopChrdEP2TnsFcMid float64    
    TopChrdEP2MomDwnMuMid float64 
    TopChrdEP2MomDwnFbMid float64 
    TopChrdEP2MomUpMuMid float64  
    TopChrdEP2MomUpfbMid float64  

    TopChrdEP2TnsForcePoint float64 
    TopChrdEP2TnsFtPoint float64    
    TopChrdEP2CmpForcePoint float64 
    TopChrdEP2TnsFcPoint float64    
    TopChrdEP2MomDwnMuPoint float64 
    TopChrdEP2MomDwnFbPoint float64 
    TopChrdEP2MomUpMuPoint float64  
    TopChrdEP2MomUpfbPoint float64  
}
// 5.- Design of Chords Begin
// Resistance Factor
type ResistanceFactor struct {
    TensionFactor string // value as default is "Phi t"
    TensionValue float64 // value as default is 0.9 Fy
    TensionDStress  string// string value as default is 0.9 Fy

    CompressionFactor string // value as default is "Phi c"
    CompressionValue float64 // value as default is 0.9 Fy
    CompressionDStress string // string value as default is 0.9 Fy

    BendingFactor string // value as default is "Phi t"
    BendingValue float64 // value as default is 0.9 Fy
    BendingDStress string // string value as default is 0.9 Fy
}
type ResistanceFactors = []ResistanceFactor
type TableResistance struct {
    ResistanceFactors ResistanceFactors
}
func newResistanceFactor(
    tensionFactor, compressionFactor, bendingFactor,
    tensionDStress, compressionDStress, bendingDStress string,
    tensionValue, compressionValue, bendingValue float64) ResistanceFactor{
    return ResistanceFactor{
	TensionFactor: tensionFactor,
	TensionValue: tensionValue,
	TensionDStress: tensionDStress,
	CompressionFactor: compressionFactor,
	CompressionValue: compressionValue,
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
    // BCcheck bool

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
}
type CheckSlenderness = []CheckSlendernes
type TableCheck struct {
    CheckSlenderness CheckSlenderness
}
func newCheckSlendernes( bCIx,  bCIy,   bCIz,   bCrx,   bCry,   bCrz,   bCIxrx,   bCIyry,   bCIzrz,   bCSLRym,   bCrgov,   bClimit,   tCIpMidIx,   tCIpMidIy,   tCIpMidIz,   tCIpMidrx,   tCIpMidry,   tCIpMidrz,   tCIpMidIxrx,   tCIpMidIyry,   tCIpMidIzrz,   tCIpMidSLRym,   tCIpMidrgov,   tCIpMidlimit,   tCIpPointIx,   tCIpPointIy,   tCIpPointIz,   tCIpPointrx,   tCIpPointry,   tCIpPointrz,   tCIpPointIxrx,   tCIpPointIyry,   tCIpPointIzrz,   tCIpPointSLRym,   tCIpPointrgov,   tCIpPointlimit,   tCEp1MidIx,   tCEp1MidIy,   tCEp1MidIz,   tCEp1Midrx,   tCEp1Midry,   tCEp1Midrz,   tCEp1MidIxrx,   tCEp1MidIyry,   tCEp1MidIzrz,   tCEp1MidSLRym,   tCEp1Midrgov,   tCEp1Midlimit,   tCEp1PointIx,   tCEp1PointIy,   tCEp1PointIz,   tCEp1Pointrx,   tCEp1Pointry,   tCEp1Pointrz,   tCEp1PointIxrx,   tCEp1PointIyry,   tCEp1PointIzrz,   tCEp1PointSLRym,   tCEp1Pointrgov,   tCEp1Pointlimit,   tCEp2MidIx,   tCEp2MidIy,   tCEp2MidIz,   tCEp2Midrx,   tCEp2Midry,   tCEp2Midrz,   tCEp2MidIxrx,   tCEp2MidIyry,   tCEp2MidIzrz,   tCEp2MidSLRym,   tCEp2Midrgov,   tCEp2Midlimit,   tCEp2PointIx,   tCEp2PointIy,   tCEp2PointIz,   tCEp2Pointrx,   tCEp2Pointry,   tCEp2Pointrz,   tCEp2PointIxrx,   tCEp2PointIyry,   tCEp2PointIzrz,   tCEp2PointSLRym,   tCEp2Pointrgov, tCEp2Pointlimit float64)CheckSlendernes{
    return CheckSlendernes{ BCIx:             bCIx, BCIy:             bCIy, BCIz:             bCIz, BCrx:             bCrx, BCry:             bCry,   BCrz:             bCrz,   BCIxrx:           bCIxrx, BCIyry:           bCIyry, BCIzrz:           bCIzrz,  BCSLRym:          bCSLRym,   BCrgov:           bCrgov,   BClimit:          bClimit, TCIpMidIx:        tCIpMidIx, TCIpMidIy:        tCIpMidIy,   TCIpMidIz:        tCIpMidIz,   TCIpMidrx:        tCIpMidrx,   TCIpMidry:        tCIpMidry,   TCIpMidrz:        tCIpMidrz,   TCIpMidIxrx:      tCIpMidIxrx,   TCIpMidIyry:      tCIpMidIyry,   TCIpMidIzrz:      tCIpMidIzrz,   TCIpMidSLRym:     tCIpMidSLRym,   TCIpMidrgov:      tCIpMidrgov,   TCIpMidlimit:     tCIpMidlimit,   TCIpPointIx:      tCIpPointIx,   TCIpPointIy:      tCIpPointIy,   TCIpPointIz:      tCIpPointIz,   TCIpPointrx:      tCIpPointrx,   TCIpPointry:      tCIpPointry,   TCIpPointrz:      tCIpPointrz,   TCIpPointIxrx:    tCIpPointIxrx,  TCIpPointIyry:    tCIpPointIyry,  TCIpPointIzrz:    tCIpPointIzrz,  TCIpPointSLRym:   tCIpPointSLRym, TCIpPointrgov:    tCIpPointrgov,  TCIpPointlimit:   tCIpPointlimit, TCEp1MidIx:       tCEp1MidIx,   TCEp1MidIy:       tCEp1MidIy,   TCEp1MidIz:       tCEp1MidIz,   TCEp1Midrx:       tCEp1Midrx,   TCEp1Midry:       tCEp1Midry,   TCEp1Midrz:       tCEp1Midrz,   TCEp1MidIxrx:     tCEp1MidIxrx,   TCEp1MidIyry:     tCEp1MidIyry,   TCEp1MidIzrz:     tCEp1MidIzrz,   TCEp1MidSLRym:    tCEp1MidSLRym,  TCEp1Midrgov:     tCEp1Midrgov,   TCEp1Midlimit:    tCEp1Midlimit,  TCEp1PointIx:     tCEp1PointIx,   TCEp1PointIy:     tCEp1PointIy,   TCEp1PointIz:     tCEp1PointIz,   TCEp1Pointrx:     tCEp1Pointrx,   TCEp1Pointry:     tCEp1Pointry,   TCEp1Pointrz:     tCEp1Pointrz,   TCEp1PointIxrx:   tCEp1PointIxrx, TCEp1PointIyry:   tCEp1PointIyry, TCEp1PointIzrz:   tCEp1PointIzrz, TCEp1PointSLRym:  tCEp1PointSLRym, TCEp1Pointrgov:   tCEp1Pointrgov, TCEp1Pointlimit:  tCEp1Pointlimit, TCEp2MidIx:       tCEp2MidIx,   TCEp2MidIy:       tCEp2MidIy,   TCEp2MidIz:       tCEp2MidIz,   TCEp2Midrx:       tCEp2Midrx,   TCEp2Midry:       tCEp2Midry,   TCEp2Midrz:       tCEp2Midrz,   TCEp2MidIxrx:     tCEp2MidIxrx,   TCEp2MidIyry:     tCEp2MidIyry,   TCEp2MidIzrz:     tCEp2MidIzrz,   TCEp2MidSLRym:    tCEp2MidSLRym,  TCEp2Midrgov:     tCEp2Midrgov,   TCEp2Midlimit:    tCEp2Midlimit,  TCEp2PointIx:     tCEp2PointIx,   TCEp2PointIy:     tCEp2PointIy,   TCEp2PointIz:     tCEp2PointIz,   TCEp2Pointrx:     tCEp2Pointrx,   TCEp2Pointry:     tCEp2Pointry,   TCEp2Pointrz:     tCEp2Pointrz,   TCEp2PointIxrx:   tCEp2PointIxrx, TCEp2PointIyry:   tCEp2PointIyry, TCEp2PointIzrz:   tCEp2PointIzrz, TCEp2PointSLRym:  tCEp2PointSLRym, TCEp2Pointrgov:   tCEp2Pointrgov, TCEp2Pointlimit:  tCEp2Pointlimit, 
    }
}
func newTableCheck() TableCheck{
    return TableCheck {
	CheckSlenderness: []CheckSlendernes{
	    newCheckSlendernes(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0),
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
func newEfectiveSlendernes( bCklrx, bCklry, bCklrz, bCklsrz, bCKLrx, bCSlendklrx, bCSlendklry, bCSlendklrz, bCSlendklsrz, bCSLRgov, bCFeKlrx, tCIpMidklrx, tCIpMidklry, tCIpMidklrz, tCIpMidklsrz, tCIpMidKLrx, tCIpMidSlendklrx, tCIpMidSlendklry, tCIpMidSlendklrz, tCIpMidSlendklsrz, tCIpMidSLRgov, tCIpMidFeKlrx, tCIpPointklrx, tCIpPointklry, tCIpPointklrz, tCIpPointklsrz, tCIpPointKLrx, tCIpPointSlendklrx, tCIpPointSlendklry, tCIpPointSlendklrz, tCIpPointSlendklsrz, tCIpPointSLRgov, tCIpPointFeKlrx, tCEp1Midklrx, tCEp1Midklry, tCEp1Midklrz, tCEp1Midklsrz, tCEp1MidKLrx, tCEp1MidSlendklrx, tCEp1MidSlendklry, tCEp1MidSlendklrz, tCEp1MidSlendklsrz, tCEp1MidSLRgov, tCEp1MidFeKlrx, tCEp1Pointklrx, tCEp1Pointklry, tCEp1Pointklrz, tCEp1Pointklsrz, tCEp1PointKLrx, tCEp1PointSlendklrx, tCEp1PointSlendklry, tCEp1PointSlendklrz, tCEp1PointSlendklsrz, tCEp1PointSLRgov, tCEp1PointFeKlrx, tCEp2Midklrx, tCEp2Midklry, tCEp2Midklrz, tCEp2Midklsrz, tCEp2MidKLrx, tCEp2MidSlendklrx, tCEp2MidSlendklry, tCEp2MidSlendklrz, tCEp2MidSlendklsrz, tCEp2MidSLRgov, tCEp2MidFeKlrx, tCEp2Pointklrx, tCEp2Pointklry, tCEp2Pointklrz, tCEp2Pointklsrz, tCEp2PointKLrx, tCEp2PointSlendklrx, tCEp2PointSlendklry, tCEp2PointSlendklrz, tCEp2PointSlendklsrz, tCEp2PointSLRgov, tCEp2PointFeKlrx float64) EfectiveSlendernes{
	return EfectiveSlendernes{
	    BCklrx:                bCklrx,
	    BCklry:                bCklry,
	    BCklrz:                bCklrz,
	    BCklsrz:               bCklsrz,
	    BCKLrx:                bCKLrx,
	    BCSlendklrx:           bCSlendklrx,
	    BCSlendklry:           bCSlendklry,
	    BCSlendklrz:           bCSlendklrz,
	    BCSlendklsrz:          bCSlendklsrz,
	    BCSLRgov:              bCSLRgov,
	    BCFeKlrx:              bCFeKlrx,
	    
	    TCIpMidklrx:           tCIpMidklrx,
	    TCIpMidklry:           tCIpMidklry,
	    TCIpMidklrz:           tCIpMidklrz,
	    TCIpMidklsrz:          tCIpMidklsrz,
	    TCIpMidKLrx:           tCIpMidKLrx,

	    TCIpMidSlendklrx:      tCIpMidSlendklrx,

	    TCIpMidSlendklry:      tCIpMidSlendklry,
	    TCIpMidSlendklrz:      tCIpMidSlendklrz,
	    TCIpMidSlendklsrz:     tCIpMidSlendklsrz,
	    TCIpMidSLRgov:         tCIpMidSLRgov,
	    TCIpMidFeKlrx:         tCIpMidFeKlrx,

	    TCIpPointklrx:         tCIpPointklrx,
	    TCIpPointklry:         tCIpPointklry,
	    TCIpPointklrz:         tCIpPointklrz,
	    TCIpPointklsrz:        tCIpPointklsrz,
	    TCIpPointKLrx:         tCIpPointKLrx,
	    TCIpPointSlendklrx:    tCIpPointSlendklrx,
	    TCIpPointSlendklry:    tCIpPointSlendklry,
	    TCIpPointSlendklrz:    tCIpPointSlendklrz,
	    TCIpPointSlendklsrz:   tCIpPointSlendklsrz,
	    TCIpPointSLRgov:       tCIpPointSLRgov,
	    TCIpPointFeKlrx:       tCIpPointFeKlrx,

	    TCEp1Midklrx:          tCEp1Midklrx,
	    TCEp1Midklry:          tCEp1Midklry,
	    TCEp1Midklrz:          tCEp1Midklrz,
	    TCEp1Midklsrz:         tCEp1Midklsrz,
	    TCEp1MidKLrx:          tCEp1MidKLrx,
	    TCEp1MidSlendklrx:     tCEp1MidSlendklrx,
	    TCEp1MidSlendklry:     tCEp1MidSlendklry,
	    TCEp1MidSlendklrz:     tCEp1MidSlendklrz,
	    TCEp1MidSlendklsrz:    tCEp1MidSlendklsrz,
	    TCEp1MidSLRgov:        tCEp1MidSLRgov,
	    TCEp1MidFeKlrx:        tCEp1MidFeKlrx,

	    TCEp1Pointklrx:        tCEp1Pointklrx,
	    TCEp1Pointklry:        tCEp1Pointklry,
	    TCEp1Pointklrz:        tCEp1Pointklrz,
	    TCEp1Pointklsrz:       tCEp1Pointklsrz,
	    TCEp1PointKLrx:        tCEp1PointKLrx,
	    TCEp1PointSlendklrx:   tCEp1PointSlendklrx,
	    TCEp1PointSlendklry:   tCEp1PointSlendklry,
	    TCEp1PointSlendklrz:   tCEp1PointSlendklrz,
	    TCEp1PointSlendklsrz:  tCEp1PointSlendklsrz,
	    TCEp1PointSLRgov:      tCEp1PointSLRgov,
	    TCEp1PointFeKlrx:      tCEp1PointFeKlrx,

	    TCEp2Midklrx:          tCEp2Midklrx,
	    TCEp2Midklry:          tCEp2Midklry,
	    TCEp2Midklrz:          tCEp2Midklrz,
	    TCEp2Midklsrz:         tCEp2Midklsrz,
	    TCEp2MidKLrx:          tCEp2MidKLrx,
	    TCEp2MidSlendklrx:     tCEp2MidSlendklrx,
	    TCEp2MidSlendklry:     tCEp2MidSlendklry,
	    TCEp2MidSlendklrz:     tCEp2MidSlendklrz,
	    TCEp2MidSlendklsrz:    tCEp2MidSlendklsrz,
	    TCEp2MidSLRgov:        tCEp2MidSLRgov,
	    TCEp2MidFeKlrx:        tCEp2MidFeKlrx,

	    TCEp2Pointklrx:        tCEp2Pointklrx,
	    TCEp2Pointklry:        tCEp2Pointklry,
	    TCEp2Pointklrz:        tCEp2Pointklrz,
	    TCEp2Pointklsrz:       tCEp2Pointklsrz,
	    TCEp2PointKLrx:        tCEp2PointKLrx,
	    TCEp2PointSlendklrx:   tCEp2PointSlendklrx,
	    TCEp2PointSlendklry:   tCEp2PointSlendklry,
	    TCEp2PointSlendklrz:   tCEp2PointSlendklrz,
	    TCEp2PointSlendklsrz:  tCEp2PointSlendklsrz,
	    TCEp2PointSLRgov:      tCEp2PointSLRgov,
	    TCEp2PointFeKlrx:      tCEp2PointFeKlrx,      
	}
}
func newTableEfective() TableEfective{
    return TableEfective {
	EfectiveSlenderness: []EfectiveSlendernes{
	    newEfectiveSlendernes(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0),
	},
    }
}
// Design Stress
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

type ShearCap struct {
    Vep float64
    Pu_Ep float64
    Fn float64
    Fv float64

    Botfv float64
    Botfa float64
    Botfvmod float64
    BotFv float64
    Botcheck float64

    Topfv float64
    Topfa float64
    Topfvmod float64
    TopFv float64
    Topcheck float64
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
    Check float64
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

type DOWDesign struct {
    InputName string
    ElemName string
    TenFt float64
    TenPut float64
    TenPuSol float64
    TenRat float64
    TenAllow float64
    Compbt float64
    CompQ float64
    CompFe float64
    CompFcr float64
    CompFc float64
    CompPuc float64
    CompPuSol float64
    CompRat float64
    CompAllow float64
}
type DOWForce struct {
    Rmax float64
    VSheard float64
    Beta float64
    Gamma float64
    Delta float64
    Alpha float64
}
// Moment of Inertia
type Moment struct {
    Ijoist float64
    Ireq360 float64
    Ireq240 float64
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
    ae float64
    a float64
    b float64
    c float64
    W float64
    Wactual float64
    WuMin float64
}
// Design of Weld
type DesignWeld struct {
    Resistance float64
    Tensile float64
    Nominal float64
    Angle float64
    Fillet float64
    WeldL float64
    Force float64
    Lw float64
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

type Properties = []Propertie


type AnglProps = []AnglProp
type BridgingProps = []BridgingProp
type MemberInputs = []MemberInput
type AnalysisResults = []AnalysisResult



type Designs = []Design
type CapSols = []CapSol
type ShearCaps = []ShearCap
type DOWSlendRats = []DOWSlendRat
type DOWWebs = []DOWWeb
type DOWEfectives =[]DOWEfective
type DOWDesigns = []DOWDesign
type DOWForces = []DOWForce
type Moments = []Moment
type Laterals = []Lateral
type DesignWelds = []DesignWeld
type DesignWeldCons = []DesignWeldCon

type Geometry struct{
    Properties Properties
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
type TableAnalysis struct {
    AnalysisResults AnalysisResults
}
type TableDesign struct {
    Designs Designs
}
type TableCap struct {
    CapSols CapSols
}
type TableShear struct {
    ShearCaps ShearCaps
}
type TableDOWSlend struct {
    DOWSlendRats DOWSlendRats
}
type TableDOWWeb struct {
    DOWWebs DOWWebs
}
type TableDOWEfective struct {
    DOWEfectives DOWEfectives
}
type TableDOWDesign struct {
    DOWDesigns DOWDesigns
}
type TableDOWForce struct {
    DOWForces DOWForces
}
type TableMoment struct {
    Moments Moments
}
type TableLateral struct {
    Laterals Laterals
}
type TableDesignWeld struct {
    DesignWelds DesignWelds
}
type TableDesignWeldCon struct {
    DesignWeldCons DesignWeldCons
}
func newDesignWeldCon(
    InputName,
    ElemName,
    Conclusion string,
    Ftens,
    FComp,
    F50,
    FGov,
    Tw,
    Lwmin,
    Lw float64)DesignWeldCon{
    return DesignWeldCon{
	InputName:  InputName,
	ElemName:   ElemName,
	Conclusion: Conclusion,
	Ftens:      Ftens,
	FComp:      FComp,
	F50:        F50,
	FGov:       FGov,
	Tw:         Tw,
	Lwmin:      Lwmin,
	Lw:         Lw,        
    }
}
func newDesignWeld(
    Resistance,
    Tensile,
    Nominal,
    Angle,
    Fillet,
    WeldL,
    Force,
    Lw float64)DesignWeld{
    return DesignWeld{
	Resistance: Resistance, 
	Tensile:    Tensile,
	Nominal:    Nominal,
	Angle:      Angle,
	Fillet:     Fillet,
	WeldL:      WeldL,
	Force:      Force,
	Lw:         Lw,         
    }
}
func newLateral(
    p,
    k,
    g,
    j,
    y,
    iytc,
    iybc,
    iy,
    ix,
    y0,
    cw,
    betax,
    ae,
    a,
    b,
    c,
    w,
    wactual,
    wuMin float64)Lateral{
    return Lateral{
	P:        p,      
	K:        k,
	G:        g,
	J:        j,
	Y:        y,
	Iytc:     iytc,
	Iybc:     iybc,
	Iy:       iy,
	Ix:       ix,
	Y0:       y0,
	Cw:       cw,
	Betax:    betax,
	ae:       ae,
	a:        a,
	b:        b,
	c:        c,
	W:        w,
	Wactual:  wactual,
	WuMin:    wuMin,  
    }
}
func newMoment(
    ijoist,
    ireq360,
    ireq240 float64)Moment{
    return Moment{
	Ijoist:  ijoist,  
	Ireq360: ireq360,
	Ireq240: ireq240,
    }
}
func newDOWForce(
    rmax,
    vSheard,
    beta,
    gamma,
    delta,
    alpha float64)DOWForce{
    return DOWForce{
	Rmax:     rmax,   
	VSheard:  vSheard,
	Beta:     beta,
	Gamma:    gamma,
	Delta:    delta,
	Alpha:    alpha,  
    }
}
func newDOWDesign(
    inputName,
    elemName string,
    tenFt,
    tenPut,
    tenPuSol,
    tenRat,
    tenAllow,
    compbt,
    compQ,
    compFe,
    compFcr,
    compFc,
    compPuc,
    compPuSol,
    compRat,
    compAllow float64)DOWDesign{
    return DOWDesign{
	InputName:   inputName,
	ElemName:    elemName,
	TenFt:       tenFt,    
	TenPut:      tenPut,
	TenPuSol:    tenPuSol,
	TenRat:      tenRat,
	TenAllow:    tenAllow,
	Compbt:      compbt,
	CompQ:       compQ,
	CompFe:      compFe,
	CompFcr:     compFcr,
	CompFc:      compFc,
	CompPuc:     compPuc,
	CompPuSol:   compPuSol,
	CompRat:     compRat,
	CompAllow:   compAllow,
    }
}
func newDOWEfective(
    inputName,
    elemName string,
    klrx,
    klry,
    klrz,
    klsrz,
    slendKlrx,
    slendKlry,
    slendKlrz,
    slendKlsrz,
    sLRgov float64)DOWEfective{
    return DOWEfective{
	InputName:   inputName,
	ElemName:    elemName,
	Klrx:       klrx,      
	Klry:       klry,
	Klrz:       klrz,
	Klsrz:      klsrz,
	SlendKlrx:  slendKlrx,
	SlendKlry:  slendKlry,
	SlendKlrz:  slendKlrz,
	SlendKlsrz: slendKlsrz,
	SLRgov:     sLRgov,    
    }
}
func newDOWWeb(
    inputName,
    elemName string,
    xPE,
    xPC,
    eQV1,
    eQV2,
    vmin,
    ftMin,
    fcMin,
    liftFTens,
    liftFComp,
    designFTens,
    designFComp float64)DOWWeb{
    return DOWWeb {
	InputName:   inputName,
	ElemName:    elemName,
	XPE:         xPE,
	XPC:         xPC,
	EQV1:        eQV1,
	EQV2:        eQV2,
	Vmin:        vmin,
	FtMin:       ftMin,
	FcMin:       fcMin,
	LiftFTens:   liftFTens,
	LiftFComp:   liftFComp,
	DesignFTens: designFTens,
	DesignFComp: designFComp,
    }
}
func newDOWSlendRat(
    inputName,
    elemName string,
    iX,
    iY,
    iZ,
    rx,
    ry,
    rz,
    ixrx,
    iyry,
    izrz,
    sLRym,
    lrgov,
    limit,
    check float64) DOWSlendRat{
	return DOWSlendRat{
	    InputName:  inputName,
	    ElemName:   elemName,
	    IX:         iX,
	    IY:         iY,
	    IZ:         iZ,
	    Rx:         rx,
	    Ry:         ry,
	    Rz:         rz,
	    Ixrx:       ixrx,
	    Iyry:       iyry,
	    Izrz:       izrz,
	    SLRym:      sLRym,
	    Lrgov:      lrgov,
	    Limit:      limit,
	    Check:      check,    
	}
}
func newShearCap(
    vep,
    pu_Ep,
    fn,
    fv,
    botfv,
    botfa,
    botfvmod,
    botFv,
    botcheck,
    topfv,
    topfa,
    topfvmod,
    topFv,
    topcheck float64)ShearCap{
	return ShearCap{
	    Vep:       vep,
	    Pu_Ep:     pu_Ep,
	    Fn:        fn,
	    Fv:        fv,
	    Botfv:     botfv,
	    Botfa:     botfa,
	    Botfvmod:  botfvmod,
	    BotFv:     botFv,
	    Botcheck:  botcheck,
	    Topfv:     topfv,
	    Topfa:     topfa,
	    Topfvmod:  topfvmod,
	    TopFv:     topFv,
	    Topcheck:  topcheck,
	}
}
func newCapSol( bCPut, bCPusol, bCTenRat, bCTenAllow, bCFau, bCFbu, bCFauFc, bCAxRatio, bCAxAllow, bCcheck, tCIpMidPut, tCIpMidPusol, tCIpMidTenRat, tCIpMidTenAllow, tCIpMidFau, tCIpMidFbu, tCIpMidFauFc, tCIpMidAxRatio, tCIpMidAxAllow, tCIpMidcheck, tCIpPointPut, tCIpPointPusol, tCIpPointTenRat, tCIpPointTenAllow, tCIpPointFau, tCIpPointFbu, tCIpPointFauFc, tCIpPointAxRatio, tCIpPointAxAllow, tCIpPointcheck, tCEp1MidPut, tCEp1MidPusol, tCEp1MidTenRat, tCEp1MidTenAllow, tCEp1MidFau, tCEp1MidFbu, tCEp1MidFauFc, tCEp1MidAxRatio, tCEp1MidAxAllow, tCEp1Midcheck, tCEp1PointPut, tCEp1PointPusol, tCEp1PointTenRat, tCEp1PointTenAllow, tCEp1PointFau, tCEp1PointFbu, tCEp1PointFauFc, tCEp1PointAxRatio, tCEp1PointAxAllow, tCEp1Pointcheck, tCEp2MidPut, tCEp2MidPusol, tCEp2MidTenRat, tCEp2MidTenAllow, tCEp2MidFau, tCEp2MidFbu, tCEp2MidFauFc, tCEp2MidAxRatio, tCEp2MidAxAllow, tCEp2Midcheck, tCEp2PointPut, tCEp2PointPusol, tCEp2PointTenRat, tCEp2PointTenAllow, tCEp2PointFau, tCEp2PointFbu, tCEp2PointFauFc, tCEp2PointAxRatio, tCEp2PointAxAllow, tCEp2Pointcheck float64) CapSol{
	return CapSol{
	    BCPut:                bCPut,             
	    BCPusol:              bCPusol,
	    BCTenRat:             bCTenRat,
	    BCTenAllow:           bCTenAllow,
	    BCFau:                bCFau,
	    BCFbu:                bCFbu,
	    BCFauFc:              bCFauFc,
	    BCAxRatio:            bCAxRatio,
	    BCAxAllow:            bCAxAllow,
	    BCcheck:              bCcheck,
	    TCIpMidPut:           tCIpMidPut,
	    TCIpMidPusol:         tCIpMidPusol,
	    TCIpMidTenRat:        tCIpMidTenRat,
	    TCIpMidTenAllow:      tCIpMidTenAllow,
	    TCIpMidFau:           tCIpMidFau,
	    TCIpMidFbu:           tCIpMidFbu,
	    TCIpMidFauFc:         tCIpMidFauFc,
	    TCIpMidAxRatio:       tCIpMidAxRatio,
	    TCIpMidAxAllow:       tCIpMidAxAllow,
	    TCIpMidcheck:         tCIpMidcheck,
	    TCIpPointPut:         tCIpPointPut,
	    TCIpPointPusol:       tCIpPointPusol,
	    TCIpPointTenRat:      tCIpPointTenRat,
	    TCIpPointTenAllow:    tCIpPointTenAllow,
	    TCIpPointFau:         tCIpPointFau,
	    TCIpPointFbu:         tCIpPointFbu,
	    TCIpPointFauFc:       tCIpPointFauFc,
	    TCIpPointAxRatio:     tCIpPointAxRatio,
	    TCIpPointAxAllow:     tCIpPointAxAllow,
	    TCIpPointcheck:       tCIpPointcheck,
	    TCEp1MidPut:          tCEp1MidPut,
	    TCEp1MidPusol:        tCEp1MidPusol,
	    TCEp1MidTenRat:       tCEp1MidTenRat,
	    TCEp1MidTenAllow:     tCEp1MidTenAllow,
	    TCEp1MidFau:          tCEp1MidFau,
	    TCEp1MidFbu:          tCEp1MidFbu,
	    TCEp1MidFauFc:        tCEp1MidFauFc,
	    TCEp1MidAxRatio:      tCEp1MidAxRatio,
	    TCEp1MidAxAllow:      tCEp1MidAxAllow,
	    TCEp1Midcheck:        tCEp1Midcheck,
	    TCEp1PointPut:        tCEp1PointPut,
	    TCEp1PointPusol:      tCEp1PointPusol,
	    TCEp1PointTenRat:     tCEp1PointTenRat,
	    TCEp1PointTenAllow:   tCEp1PointTenAllow,
	    TCEp1PointFau:        tCEp1PointFau,
	    TCEp1PointFbu:        tCEp1PointFbu,
	    TCEp1PointFauFc:      tCEp1PointFauFc,
	    TCEp1PointAxRatio:    tCEp1PointAxRatio,
	    TCEp1PointAxAllow:    tCEp1PointAxAllow,
	    TCEp1Pointcheck:      tCEp1Pointcheck,
	    TCEp2MidPut:          tCEp2MidPut,
	    TCEp2MidPusol:        tCEp2MidPusol,
	    TCEp2MidTenRat:       tCEp2MidTenRat,
	    TCEp2MidTenAllow:     tCEp2MidTenAllow,
	    TCEp2MidFau:          tCEp2MidFau,
	    TCEp2MidFbu:          tCEp2MidFbu,
	    TCEp2MidFauFc:        tCEp2MidFauFc,
	    TCEp2MidAxRatio:      tCEp2MidAxRatio,
	    TCEp2MidAxAllow:      tCEp2MidAxAllow,
	    TCEp2Midcheck:        tCEp2Midcheck,
	    TCEp2PointPut:        tCEp2PointPut,
	    TCEp2PointPusol:      tCEp2PointPusol,
	    TCEp2PointTenRat:     tCEp2PointTenRat,
	    TCEp2PointTenAllow:   tCEp2PointTenAllow,
	    TCEp2PointFau:        tCEp2PointFau,
	    TCEp2PointFbu:        tCEp2PointFbu,
	    TCEp2PointFauFc:      tCEp2PointFauFc,
	    TCEp2PointAxRatio:    tCEp2PointAxRatio,
	    TCEp2PointAxAllow:    tCEp2PointAxAllow,
	    TCEp2Pointcheck:      tCEp2Pointcheck,   
	}
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

func roundFloat(val float64, precision uint) float64 {
    ratio := math.Pow(10, float64(precision))
    return math.Round(val * ratio) / ratio
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
func newBridgingProp (brid1, brid2 string) BridgingProp {
    bridging1,_ := strconv.ParseFloat(brid1, 64)
    bridging2,_ := strconv.ParseFloat(brid2, 64)
    return BridgingProp{
	Brdgng1: bridging1,
	Brdgng2: bridging2,
    }
}
func newMemberInput(
    numInput, 
    elemName, 
    part, 
    mark, 
    crimped, 
    secction, 
    midPanel string, a, q, s, 
    t,
    rxdouble, 
    rydouble,
    rzdouble,
    rxcrimped, 
    rycrimped, 
    rxuncrimped,
    ryuncrimped float64) MemberInput {
    return MemberInput{
	InputName: numInput,
	ElemName: elemName,
	Part: part,
	Mark: mark,
	Crimped: crimped,
	Secction: secction,
	MidPanel: midPanel,
	A: a,
	Q: q,
	S: s,
	T: t,
	rxDouble: rxdouble,
	ryDouble: rydouble,
	rzDouble: rzdouble,

	rxCrimped: rxcrimped,
	ryCrimped: rycrimped,
	rxUncrimped: rxuncrimped,
	ryUncrimped: ryuncrimped,
    }
}

func newAnalysisResult ( botChrdTnsForce, botChrdTnsFt, botChrdCmpForce, botChrdMomDwnMu, botChrdMomDwnFb, botChrdMomUpMu, botChrdMomUpfb, topChrdTnsForceMid, topChrdTnsFtMid, topChrdCmpForceMid, topChrdTnsFcMid, topChrdMomDwnMuMid, topChrdMomDwnFbMid, topChrdMomUpMuMid, topChrdMomUpfbMid, topChrdTnsForcePoint, topChrdTnsFtPoint, topChrdCmpForcePoint, topChrdTnsFcPoint, topChrdMomDwnMuPoint, topChrdMomDwnFbPoint, topChrdMomUpMuPoint, topChrdMomUpfbPoint, topChrdEP1TnsForceMid, topChrdEP1TnsFtMid, topChrdEP1CmpForceMid, topChrdEP1TnsFcMid, topChrdEP1MomDwnMuMid, topChrdEP1MomDwnFbMid, topChrdEP1MomUpMuMid, topChrdEP1MomUpfbMid, topChrdEP1TnsForcePoint, topChrdEP1TnsFtPoint, topChrdEP1CmpForcePoint, topChrdEP1TnsFcPoint, topChrdEP1MomDwnMuPoint, topChrdEP1MomDwnFbPoint, topChrdEP1MomUpMuPoint, topChrdEP1MomUpfbPoint, topChrdEP2TnsForceMid, topChrdEP2TnsFtMid, topChrdEP2CmpForceMid, topChrdEP2TnsFcMid, topChrdEP2MomDwnMuMid, topChrdEP2MomDwnFbMid, topChrdEP2MomUpMuMid, topChrdEP2MomUpfbMid, topChrdEP2TnsForcePoint, topChrdEP2TnsFtPoint, topChrdEP2CmpForcePoint, topChrdEP2TnsFcPoint, topChrdEP2MomDwnMuPoint, topChrdEP2MomDwnFbPoint, topChrdEP2MomUpMuPoint, topChrdEP2MomUpfbPoint float64) AnalysisResult {
    return AnalysisResult { BotChrdTnsForce:         botChrdTnsForce, BotChrdTnsFt:            botChrdTnsFt, BotChrdCmpForce:         botChrdCmpForce, BotChrdMomDwnMu:         botChrdMomDwnMu, BotChrdMomDwnFb:         botChrdMomDwnFb, BotChrdMomUpMu:          botChrdMomUpMu, BotChrdMomUpfb:          botChrdMomUpfb, TopChrdTnsForceMid:      topChrdTnsForceMid, TopChrdTnsFtMid:         topChrdTnsFtMid, TopChrdCmpForceMid:      topChrdCmpForceMid, TopChrdTnsFcMid:         topChrdTnsFcMid, TopChrdMomDwnMuMid:      topChrdMomDwnMuMid, TopChrdMomDwnFbMid:      topChrdMomDwnFbMid, TopChrdMomUpMuMid:       topChrdMomUpMuMid, TopChrdMomUpfbMid:       topChrdMomUpfbMid, TopChrdTnsForcePoint:    topChrdTnsForcePoint, TopChrdTnsFtPoint:       topChrdTnsFtPoint, TopChrdCmpForcePoint:    topChrdCmpForcePoint, TopChrdTnsFcPoint:       topChrdTnsFcPoint, TopChrdMomDwnMuPoint:    topChrdMomDwnMuPoint, TopChrdMomDwnFbPoint:    topChrdMomDwnFbPoint, TopChrdMomUpMuPoint:     topChrdMomUpMuPoint, TopChrdMomUpfbPoint:     topChrdMomUpfbPoint, TopChrdEP1TnsForceMid:   topChrdEP1TnsForceMid, TopChrdEP1TnsFtMid:      topChrdEP1TnsFtMid, TopChrdEP1CmpForceMid:   topChrdEP1CmpForceMid, TopChrdEP1TnsFcMid:      topChrdEP1TnsFcMid, TopChrdEP1MomDwnMuMid:   topChrdEP1MomDwnMuMid, TopChrdEP1MomDwnFbMid:   topChrdEP1MomDwnFbMid, TopChrdEP1MomUpMuMid:    topChrdEP1MomUpMuMid, TopChrdEP1MomUpfbMid:    topChrdEP1MomUpfbMid, TopChrdEP1TnsForcePoint: topChrdEP1TnsForcePoint, TopChrdEP1TnsFtPoint:    topChrdEP1TnsFtPoint, TopChrdEP1CmpForcePoint: topChrdEP1CmpForcePoint, TopChrdEP1TnsFcPoint:    topChrdEP1TnsFcPoint, TopChrdEP1MomDwnMuPoint: topChrdEP1MomDwnMuPoint, TopChrdEP1MomDwnFbPoint: topChrdEP1MomDwnFbPoint, TopChrdEP1MomUpMuPoint:  topChrdEP1MomUpMuPoint, TopChrdEP1MomUpfbPoint:  topChrdEP1MomUpfbPoint, TopChrdEP2TnsForceMid:   topChrdEP2TnsForceMid, TopChrdEP2TnsFtMid:      topChrdEP2TnsFtMid, TopChrdEP2CmpForceMid:   topChrdEP2CmpForceMid, TopChrdEP2TnsFcMid:      topChrdEP2TnsFcMid, TopChrdEP2MomDwnMuMid:   topChrdEP2MomDwnMuMid, TopChrdEP2MomDwnFbMid:   topChrdEP2MomDwnFbMid, TopChrdEP2MomUpMuMid:    topChrdEP2MomUpMuMid, TopChrdEP2MomUpfbMid:    topChrdEP2MomUpfbMid, TopChrdEP2TnsForcePoint: topChrdEP2TnsForcePoint, TopChrdEP2TnsFtPoint:    topChrdEP2TnsFtPoint, TopChrdEP2CmpForcePoint: topChrdEP2CmpForcePoint, TopChrdEP2TnsFcPoint:    topChrdEP2TnsFcPoint, TopChrdEP2MomDwnMuPoint: topChrdEP2MomDwnMuPoint, TopChrdEP2MomDwnFbPoint: topChrdEP2MomDwnFbPoint, TopChrdEP2MomUpMuPoint:  topChrdEP2MomUpMuPoint, TopChrdEP2MomUpfbPoint:  topChrdEP2MomUpfbPoint, }
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
    TableResForce TableResForce
}

func newGeometry() Geometry {
    return Geometry{
	Properties: []Propertie{
	     newPropertie("warren", " roof", "240", "49.21", "27.28", "26", "24", "41.1", "28"),
	},
    }
}
func newResMater() ResMater{
    return ResMater{
	AnglProps: []AnglProp{
	    newAnglProp("", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""),
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
	    newMemberInput("", "", "", "", "", "", "", 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0),
	},
    }
}
func newTableAnalysis() TableAnalysis {
    return TableAnalysis {
	AnalysisResults: []AnalysisResult{
	    newAnalysisResult(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0),
 	},
    }
}

func newTableCap() TableCap{
    return TableCap {
	CapSols: []CapSol{
	    newCapSol(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0),
	},
    }
}
func newTableShear() TableShear{
    return TableShear {
	ShearCaps: []ShearCap{
	    newShearCap(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0),
	},
    }
}
func newTableDOWSlend() TableDOWSlend{
    return TableDOWSlend{
	DOWSlendRats: []DOWSlendRat{
	    newDOWSlendRat("", "", 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0), 
	},
    }
}
func newTableDOWWeb() TableDOWWeb{
    return TableDOWWeb{
	DOWWebs: []DOWWeb{
	    newDOWWeb("", "", 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0), 
	},
    }
}

func newTableDOWEfective() TableDOWEfective{
    return TableDOWEfective{
	DOWEfectives: []DOWEfective{
	    newDOWEfective("", "", 0, 0, 0, 0, 0, 0, 0, 0, 0), 
	},
    }
}
func newTableDOWDesign() TableDOWDesign{
    return TableDOWDesign{
	DOWDesigns: []DOWDesign{
	    newDOWDesign("", "", 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0), 
	},
    }
}
func newTableDOWForce() TableDOWForce{
    return TableDOWForce{
	DOWForces: []DOWForce{
	    newDOWForce(0, 0, 0, 0, 0, 0), 
	},
    }
}
func newTableMoment() TableMoment{
    return TableMoment{
	Moments: []Moment{
	    newMoment(0, 0, 0), 
	},
    }
}
func newTableLateral() TableLateral{
    return TableLateral{
	Laterals: []Lateral{
	    newLateral(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0), 
	},
    }
}
func newTableDesignWeld() TableDesignWeld{
    return TableDesignWeld{
	DesignWelds: []DesignWeld{
	    newDesignWeld( 0, 0, 0, 0, 0, 0, 0, 0), 
	},
    }
}
func newTableDesignWeldCon() TableDesignWeldCon{
    return TableDesignWeldCon{
	DesignWeldCons: []DesignWeldCon{
	    newDesignWeldCon("", "", "", 0, 0, 0, 0, 0, 0, 0), 
	},
    }
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
	TableResForce: newTableResForce(),
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

    /*
    check.TCIpMidIx = geom.ipl*2
    check.TCIpMidIz = check.TCIpMidIx
    */
    
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

    check.TCIpPointIx = 24
    check.TCIpPointIy = 36 
    check.TCIpPointIz = 12 
    check.TCIpPointrx = bot.RxTop
    check.TCIpPointry = bot.RyTop
    check.TCIpPointrz = bot.RzTop
    check.TCIpPointIxrx = roundFloat(check.TCIpPointIx / check.TCIpPointrx, 4)
    check.TCIpPointIyry = roundFloat(check.TCIpPointIy / check.TCIpPointry, 4)
    check.TCIpPointIzrz = roundFloat(check.TCIpPointIz / check.TCIpPointrz, 4)
    check.TCIpPointrgov = larger(check.TCIpPointIxrx, check.TCIpPointIyry, check.TCIpPointIzrz)
    check.TCIpPointlimit = 90


    check.TCEp1MidIx = 27.26
    check.TCEp1MidIy = 36
    check.TCEp1MidIz = 27.26
    check.TCEp1Midrx = bot.RxTop
    check.TCEp1Midry = bot.RyTop
    check.TCEp1Midrz = bot.RzTop
    check.TCEp1MidIxrx = roundFloat(check.TCEp1MidIx / check.TCEp1Midrx, 4)
    check.TCEp1MidIyry = roundFloat(check.TCEp1MidIy / check.TCEp1Midry, 4)
    check.TCEp1MidIzrz = roundFloat(check.TCEp1MidIz / check.TCEp1Midrz, 4)
    check.TCEp1Midrgov = larger(check.TCEp1MidIxrx, check.TCEp1MidIyry, check.TCEp1MidIzrz)
    check.TCEp1Midlimit = 120

    check.TCEp1PointIx = 27.26
    check.TCEp1PointIy = 36
    check.TCEp1PointIz = 27.26
    check.TCEp1Pointrx = bot.RxTop
    check.TCEp1Pointry = bot.RyTop
    check.TCEp1Pointrz = bot.RzTop
    check.TCEp1PointIxrx = roundFloat(check.TCEp1PointIx / check.TCEp1Pointrx, 4)
    check.TCEp1PointIyry = roundFloat(check.TCEp1PointIy / check.TCEp1Pointry, 4)
    check.TCEp1PointIzrz = roundFloat(check.TCEp1PointIz / check.TCEp1Pointrz, 4)
    check.TCEp1Pointrgov = larger(check.TCEp1PointIxrx, check.TCEp1PointIyry, check.TCEp1PointIzrz)
    check.TCEp1Pointlimit = 120

    check.TCEp2MidIx = 26
    check.TCEp2MidIy = 36
    check.TCEp2MidIz = 26
    check.TCEp2Midrx = bot.RxTop
    check.TCEp2Midry = bot.RyTop
    check.TCEp2Midrz = bot.RzTop
    check.TCEp2MidIxrx = roundFloat(check.TCEp2MidIx / check.TCEp2Midrx, 4)
    check.TCEp2MidIyry = roundFloat(check.TCEp2MidIy / check.TCEp2Midry, 4)
    check.TCEp2MidIzrz = roundFloat(check.TCEp2MidIz / check.TCEp2Midrz, 4)
    check.TCEp2Midrgov = larger(check.TCEp2MidIxrx, check.TCEp2MidIyry, check.TCEp2MidIzrz)
    check.TCEp2Midlimit = 120

    check.TCEp2PointIx = 26
    check.TCEp2PointIy = 36
    check.TCEp2PointIz = 26
    check.TCEp2Pointrx = bot.RxTop
    check.TCEp2Pointry = bot.RyTop
    check.TCEp2Pointrz = bot.RzTop
    check.TCEp2PointIxrx = roundFloat(check.TCEp2PointIx / check.TCEp2Pointrx, 4)
    check.TCEp2PointIyry = roundFloat(check.TCEp2PointIy / check.TCEp2Pointry, 4)
    check.TCEp2PointIzrz = roundFloat(check.TCEp2PointIz / check.TCEp2Pointrz, 4)
    check.TCEp2Pointrgov = larger(check.TCEp2PointIxrx, check.TCEp2PointIyry, check.TCEp2PointIzrz)
    check.TCEp2Pointlimit = 120
}
func deCalculation(page *Page) {
    d := &page.Geometry.Properties[0].depth
    ytc := &page.ResMater.AnglProps[0].YTop
    ybc := &page.ResMater.AnglProps[0].YBot
    de := (*d - *ytc) - *ybc
    page.ResGeom.ResProps[0].Ed = de
    // fmt.Println(de)
}
func efectiveSlend(page *Page) {
    // if Top chord (Lip max) mid panel fill yes check criteria B in the table
    EfectiveSlend := &page.TableEfective.EfectiveSlenderness[0]
    fillTopChord := &page.Material.Forces[0].FillTopChord
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
	EfectiveSlend.TCEp2Midklsrz = 1 } else {
	EfectiveSlend.TCEp2Midklrx = 0
	EfectiveSlend.TCEp2Midklry = 0
	EfectiveSlend.TCEp2Midklrz = 1
	EfectiveSlend.TCEp2Midklsrz = 0
    }
    EfectiveSlend.TCIpMidSlendklrx = EfectiveSlend.TCIpMidklrx * CheckSlend.TCIpMidIxrx
    EfectiveSlend.TCIpMidSlendklry = EfectiveSlend.TCIpMidklry * CheckSlend.TCIpMidIyry
    EfectiveSlend.TCIpMidSlendklrz = EfectiveSlend.TCIpMidklrz * CheckSlend.TCIpMidIzrz

    // SLR Ggov
    EfectiveSlend.TCIpMidSLRgov = larger(EfectiveSlend.TCIpMidSlendklrx, EfectiveSlend.TCIpMidSlendklry, EfectiveSlend.TCIpMidSlendklrz)

    EfectiveSlend.TCEp1MidSlendklrx = EfectiveSlend.TCEp1Midklrx * CheckSlend.TCEp1MidIxrx
    EfectiveSlend.TCEp1MidSlendklry = EfectiveSlend.TCEp1Midklry * CheckSlend.TCEp1MidIyry
    EfectiveSlend.TCEp1MidSlendklrz = EfectiveSlend.TCEp1Midklrz * CheckSlend.TCEp1MidIzrz

    EfectiveSlend.TCEp1MidSLRgov = larger(EfectiveSlend.TCEp1MidSlendklrx, EfectiveSlend.TCEp1MidSlendklry, EfectiveSlend.TCEp1MidSlendklrz)
    
    EfectiveSlend.TCEp2MidSlendklrx = EfectiveSlend.TCEp2Midklrx * CheckSlend.TCEp2MidIxrx
    EfectiveSlend.TCEp2MidSlendklry = EfectiveSlend.TCEp2Midklry * CheckSlend.TCEp2MidIyry
    EfectiveSlend.TCEp2MidSlendklrz = EfectiveSlend.TCEp2Midklrz * CheckSlend.TCEp2MidIzrz

    EfectiveSlend.TCEp2MidSLRgov = larger(EfectiveSlend.TCEp2MidSlendklrx, EfectiveSlend.TCEp2MidSlendklry, EfectiveSlend.TCEp2MidSlendklrz)
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
    e.POST("/sendform", func(c echo.Context) error {
	page.WebMember.MemberInputs = page.WebMember.MemberInputs[:0]
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

	// Chords input
	topChord := c.FormValue("topChord")
	bottomChord := c.FormValue("bottomChord")
	fillTopChord := c.FormValue("fillTopChord")
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
	obj := newMemberInput("", "", "", "", "", "", "", 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)

	tas := int(totalAngles)
	hts := int(halfTod)

	for i := 1; i < tas; i++{
	    obj.InputName = fmt.Sprintf("%v", i)
	    if i == 1{
		obj.ElemName = "sv"
	    } else if i <= hts+1 && i != 1{
		obj.ElemName = fmt.Sprintf("w%v", i)
	    } else {
		obj.ElemName = fmt.Sprintf("v%v", i - (hts+1))
	    }
	    page.WebMember.MemberInputs = append(page.WebMember.MemberInputs, obj)
	}
	slendernesRadio(&page)
	efectiveSlend(&page)
	deCalculation(&page)
	designLoads(&page)

	//fmt.Println()
	//fmt.Printf("%+v", page)
	//fmt.Println()
	return c.Render(http.StatusOK, "res", page)
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
	    //fmt.Println(page.WebMember.MemberInputs[i])
	}
	fmt.Println(page.WebMember.MemberInputs)
	return c.Render(http.StatusOK, "webMem", page )
    })
    e.Logger.Fatal(e.Start(":8080"))
}
