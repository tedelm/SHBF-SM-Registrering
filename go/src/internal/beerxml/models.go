package beerxml

import "encoding/xml"

// Recipe is the root structure parsed from BeerXML (first RECIPE in the document).
type Recipe struct {
	XMLName      xml.Name     `xml:"RECIPE"`
	Name         string       `xml:"NAME"`
	Brewer       string       `xml:"BREWER"`
	Style        Style        `xml:"STYLE"`
	OG           string       `xml:"OG"`
	FG           string       `xml:"FG"`
	EstOG        string       `xml:"EST_OG"`
	EstFG        string       `xml:"EST_FG"`
	BatchSize    string       `xml:"BATCH_SIZE"`
	DisplayBatch string       `xml:"DISPLAY_BATCH_SIZE"`
	IBU          string       `xml:"IBU"`
	EstIBU       string       `xml:"EST_IBU"`
	ABV          string       `xml:"ABV"`
	EstABV       string       `xml:"EST_ABV"`
	EstColor     string       `xml:"EST_COLOR"` // e.g. "79.3 EBC"
	Notes        string       `xml:"NOTES"`
	Hops         Hops         `xml:"HOPS"`
	Fermentables Fermentables `xml:"FERMENTABLES"`
	Mash         Mash         `xml:"MASH"`
	Yeasts       Yeasts       `xml:"YEASTS"`
	Waters       Waters       `xml:"WATERS"`
	Miscs        Miscs        `xml:"MISCS"`
}

// Style holds category and letter for SHBF style dropdown (e.g. "9:J") and optional min/max for validation.
type Style struct {
	CategoryNumber string `xml:"CATEGORY_NUMBER"`
	StyleLetter    string `xml:"STYLE_LETTER"`
	Notes          string `xml:"NOTES"`
	OGMin          string `xml:"OG_MIN"`
	OGMax          string `xml:"OG_MAX"`
	FGMin          string `xml:"FG_MIN"`
	FGMax          string `xml:"FG_MAX"`
	IBUMin         string `xml:"IBU_MIN"`
	IBUMax         string `xml:"IBU_MAX"`
	ColorMin       string `xml:"COLOR_MIN"` // EBC
	ColorMax       string `xml:"COLOR_MAX"`
	ABVMin         string `xml:"ABV_MIN"`
	ABVMax         string `xml:"ABV_MAX"`
}

// Hops wraps hop entries.
type Hops struct {
	Hop []Hop `xml:"HOP"`
}

// Hop is one hop addition.
type Hop struct {
	Name          string `xml:"NAME"`
	Alpha         string `xml:"ALPHA"`
	Amount        string `xml:"AMOUNT"`         // kg
	DisplayAmount string `xml:"DISPLAY_AMOUNT"` // e.g. "50 g"
	Time          string `xml:"TIME"`           // boil minutes
	Form          string `xml:"FORM"`           // Pellet, Leaf, etc.
	Use           string `xml:"USE"`            // Mash, Boil, Primary, Secondary, Dry Hop
	Notes         string `xml:"NOTES"`
}

// Fermentables wraps fermentable (malt) entries.
type Fermentables struct {
	Fermentable []Fermentable `xml:"FERMENTABLE"`
}

// Fermentable is one grain/malt.
type Fermentable struct {
	Name          string `xml:"NAME"`
	Amount        string `xml:"AMOUNT"` // kg
	DisplayAmount string `xml:"DISPLAY_AMOUNT"`
	Notes         string `xml:"NOTES"`
}

// Mash holds mash profile; steps are under MASH_STEPS in BeerXML.
type Mash struct {
	Steps struct {
		MashSteps []MashStep `xml:"MASH_STEP"`
	} `xml:"MASH_STEPS"`
}

// MashStep is one temperature step.
type MashStep struct {
	Name     string `xml:"NAME"`
	StepTemp string `xml:"STEP_TEMP"`
	StepTime string `xml:"STEP_TIME"`
}

// Yeasts wraps yeast entries.
type Yeasts struct {
	Yeast []Yeast `xml:"YEAST"`
}

// Yeast is one yeast.
type Yeast struct {
	Name           string `xml:"NAME"`
	ProductID      string `xml:"PRODUCT_ID"`
	MinTemperature string `xml:"MIN_TEMPERATURE"`
	MaxTemperature string `xml:"MAX_TEMPERATURE"`
}

// Waters wraps water entries.
type Waters struct {
	Water []Water `xml:"WATER"`
}

// Water is one water profile.
type Water struct {
	Name string `xml:"NAME"`
}

// Miscs wraps misc (other) additions.
type Miscs struct {
	Misc []Misc `xml:"MISC"`
}

// Misc is one misc addition (e.g. sugar, spice).
type Misc struct {
	Name           string `xml:"NAME"`
	Amount         string `xml:"AMOUNT"`
	AmountIsWeight string `xml:"AMOUNT_IS_WEIGHT"` // TRUE/false
	Use            string `xml:"USE"`              // Mash, Boil, Primary, Secondary
	Notes          string `xml:"NOTES"`
}

// Recipes is the optional root element (RECIPES with one or more RECIPE).
type Recipes struct {
	Recipe []Recipe `xml:"RECIPE"`
}
