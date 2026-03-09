package beerxml

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// recipesRoot matches <RECIPES><RECIPE>...</RECIPE></RECIPES>
type recipesRoot struct {
	XMLName xml.Name `xml:"RECIPES"`
	Recipe  []Recipe `xml:"RECIPE"`
}

// charsetReader returns a reader that decodes from the given charset (e.g. "ISO-8859-1") to UTF-8.
// BeerSmith and other tools often export with encoding="ISO-8859-1"; Go's xml package requires a CharsetReader for that.
func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch strings.ToUpper(strings.TrimSpace(charset)) {
	case "ISO-8859-1", "ISO8859-1", "LATIN1":
		return transform.NewReader(input, charmap.ISO8859_1.NewDecoder()), nil
	default:
		return input, nil
	}
}

// Parse reads BeerXML from b and returns the first RECIPE.
// Supports both a single RECIPE root and RECIPES containing one or more RECIPE.
// Handles encoding="ISO-8859-1" (and UTF-8) in the XML declaration.
func Parse(b []byte) (*Recipe, error) {
	dec := xml.NewDecoder(strings.NewReader(string(b)))
	dec.CharsetReader = charsetReader

	// Try RECIPES (plural) root first, common in BeerSmith export
	var recipes recipesRoot
	err := dec.Decode(&recipes)
	if err == nil && len(recipes.Recipe) > 0 {
		return &recipes.Recipe[0], nil
	}
	// Try single RECIPE as root (reuse decoder on fresh input)
	dec = xml.NewDecoder(strings.NewReader(string(b)))
	dec.CharsetReader = charsetReader
	var recipe Recipe
	if err := dec.Decode(&recipe); err != nil {
		return nil, fmt.Errorf("beerxml parse: %w", err)
	}
	if recipe.Name == "" && recipe.Brewer == "" && len(recipe.Hops.Hop) == 0 && len(recipe.Fermentables.Fermentable) == 0 {
		return nil, fmt.Errorf("beerxml parse: no RECIPE found in document")
	}
	return &recipe, nil
}

// ParseFromReader reads BeerXML from r and returns the first RECIPE.
func ParseFromReader(r io.Reader) (*Recipe, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("beerxml read: %w", err)
	}
	return Parse(b)
}

// ParseFile reads BeerXML from the file at path and returns the first RECIPE.
func ParseFile(path string) (*Recipe, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("beerxml open %s: %w", path, err)
	}
	defer f.Close()
	return ParseFromReader(f)
}

func MakeSHBFCompatible(beerXml *Recipe) (*Recipe, error) {
	if beerXml.Name == "" || beerXml.Brewer == "" || beerXml.Style.CategoryNumber == "" || beerXml.Style.StyleLetter == "" || beerXml.OG == "" || beerXml.FG == "" || beerXml.IBU == "" || beerXml.ABV == "" {
		return nil, fmt.Errorf("beerxml make shbf compatible: missing required fields")
	}

	//Batch size to 1 decimal place
	batchSizeFloat, _ := strconv.ParseFloat(beerXml.BatchSize, 64)
	beerXml.BatchSize = strconv.Itoa(int(math.Round(batchSizeFloat)))

	ogFloat, _ := strconv.ParseFloat(beerXml.OG, 64)
	fgFloat, _ := strconv.ParseFloat(beerXml.FG, 64)

	beerXml.IBU = strings.Split(beerXml.IBU, ".")[0]
	beerXml.ABV = strings.Split(beerXml.ABV, " ")[0]

	beerXml.OG = strconv.Itoa(int(math.Round(ogFloat * 1000)))
	beerXml.FG = strconv.Itoa(int(math.Round(fgFloat * 1000)))

	for i := range beerXml.Fermentables.Fermentable {
		fermentableAmount, _ := strconv.ParseFloat(beerXml.Fermentables.Fermentable[i].Amount, 64)
		beerXml.Fermentables.Fermentable[i].Amount = strconv.Itoa(int(math.Round(fermentableAmount * 1000)))
	}
	for i := range beerXml.Hops.Hop {
		hopAmount, _ := strconv.ParseFloat(beerXml.Hops.Hop[i].Amount, 64)
		beerXml.Hops.Hop[i].Amount = strconv.Itoa(int(math.Round(hopAmount * 1000)))

		//alpha to 1 decimal place
		alphaFloat, _ := strconv.ParseFloat(beerXml.Hops.Hop[i].Alpha, 64)
		beerXml.Hops.Hop[i].Alpha = strconv.FormatFloat(alphaFloat, 'f', 1, 64)

		//Boil time to integer
		boilTimeFloat, _ := strconv.ParseFloat(beerXml.Hops.Hop[i].Time, 64)
		beerXml.Hops.Hop[i].Time = strconv.Itoa(int(math.Round(boilTimeFloat)))
	}
	for i := range beerXml.Miscs.Misc {
		miscAmount, _ := strconv.ParseFloat(beerXml.Miscs.Misc[i].Amount, 64)
		beerXml.Miscs.Misc[i].Amount = strconv.Itoa(int(math.Round(miscAmount * 1000)))
	}

	//Mash step time and temp to integer
	for i := range beerXml.Mash.Steps.MashSteps {
		mashStepTimeFloat, _ := strconv.ParseFloat(beerXml.Mash.Steps.MashSteps[i].StepTime, 64)
		beerXml.Mash.Steps.MashSteps[i].StepTime = strconv.Itoa(int(math.Round(mashStepTimeFloat)))

		mashStepTempFloat, _ := strconv.ParseFloat(beerXml.Mash.Steps.MashSteps[i].StepTemp, 64)
		beerXml.Mash.Steps.MashSteps[i].StepTemp = strconv.Itoa(int(math.Round(mashStepTempFloat)))
	}

	//Yeasts min and max temperature to 1 decimal place
	for i := range beerXml.Yeasts.Yeast {
		yeastMinTemperatureFloat, _ := strconv.ParseFloat(beerXml.Yeasts.Yeast[i].MinTemperature, 64)
		beerXml.Yeasts.Yeast[i].MinTemperature = strconv.FormatFloat(yeastMinTemperatureFloat, 'f', 1, 64)

		yeastMaxTemperatureFloat, _ := strconv.ParseFloat(beerXml.Yeasts.Yeast[i].MaxTemperature, 64)
		beerXml.Yeasts.Yeast[i].MaxTemperature = strconv.FormatFloat(yeastMaxTemperatureFloat, 'f', 1, 64)
	}

	return beerXml, nil
}
