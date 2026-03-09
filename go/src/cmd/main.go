package main

import (
	"flag"
	"fmt"
	"os"
	"shbfsmreg/internal/beerxml"
	"shbfsmreg/internal/clientsession"
	"shbfsmreg/internal/shbf"
)

func main() {

	username := flag.String("username", "", "SHBF username")
	password := flag.String("password", "", "SHBF password")
	eventId := flag.Int("eventid", 61, "SHBF judge's event ID")
	fvEventId := flag.Int("fveventid", 62, "SHBF people's choice event ID")
	beerXmlPath := flag.String("beerxmlpath", "", "Path to BeerXML file")
	brewerName := flag.String("brewername", "", "Brewer name")
	brewerEmail := flag.String("breweremail", "", "Brewer email")
	ingredientLimit := flag.Int("ingredientlimit", 10, "Max number of rows for hops, malts, and others (default 10)")
	updateBeer := flag.Bool("updatebeer", false, "Update beer instead of creating a new one")
	verifyBeerStyle := flag.Bool("verifybeerstyle", true, "Verify beer style")
	flag.Parse()

	fmt.Println("### Starting SHBF SM registration process ###")

	if *username == "" || *password == "" || *eventId == 0 || *fvEventId == 0 || *beerXmlPath == "" || *brewerName == "" || *brewerEmail == "" {
		fmt.Println("Missing required parameters")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Read BeerXML file
	beerXml, err := beerxml.ParseFile(*beerXmlPath)
	if err != nil {
		fmt.Println("[ParseFile] Failed to parse BeerXML file: ", err)
		os.Exit(1)
	}

	if *verifyBeerStyle {
		valid, err := shbf.VerifyBeerStyle(beerXml)
		if err != nil {
			fmt.Println("[VerifyBeerStyle] Failed to verify beer style: ", err)
			os.Exit(1)
		}
		if !valid {
			fmt.Println("[VerifyBeerStyle] Beer style is not valid")
			os.Exit(1)
		}
		fmt.Println("[VerifyBeerStyle] Beer style is valid")
	}

	beerXml, err = beerxml.MakeSHBFCompatible(beerXml)
	if err != nil {
		fmt.Println("[MakeSHBFCompatible] Failed to make BeerXML compatible with SHBF: ", err)
		os.Exit(1)
	}

	fmt.Println("[MakeSHBFCompatible] BeerXML file parsed successfully")

	_, err = shbf.CreateBeerRegistrationForm(beerXml, *brewerName, *brewerEmail, *ingredientLimit)
	if err != nil {
		fmt.Println("[CreateBeerRegistrationForm] Failed to create beer registration form: ", err)
		os.Exit(1)
	}
	fmt.Println("[CreateBeerRegistrationForm] Beer registration form created successfully")

	// Create cookie jar (equivalent to WebRequestSession)
	cs, err := clientsession.NewClientSession()
	if err != nil {
		fmt.Println("[NewClientSession] Failed to create client session: ", err)
		os.Exit(1)
	}

	// Login
	err = shbf.Login(cs, *username, *password)
	if err != nil {
		fmt.Println("[Login] Failed to login: ", err)
		os.Exit(1)
	}
	// Select event
	resp, err := shbf.SelectEvent(cs, *eventId, *fvEventId)
	if err != nil {
		fmt.Println("[SelectEvent] Failed to select event: ", err)
		fmt.Println("[SelectEvent] Status:", resp.StatusCode, resp.Status)
		os.Exit(1)
	}

	fmt.Println("[SelectEvent] Select event successful")

	resp.Body.Close()

	// Register beer
	resp, err = shbf.RegisterBeer(cs, beerXml, *brewerName, *brewerEmail, *ingredientLimit, *updateBeer)
	if err != nil {
		fmt.Println("[RegisterBeer] Failed to register beer: ", err)
		fmt.Println("[RegisterBeer] Status:", resp.StatusCode, resp.Status)
		os.Exit(1)
	}
	fmt.Println("[RegisterBeer] Beer registered successfully")
	resp.Body.Close()

	// Logout
	resp, err = shbf.Logout(cs)
	if err != nil {
		fmt.Println("[Logout] Failed to logout: ", err)
		fmt.Println("[Logout] Status:", resp.StatusCode, resp.Status)
		os.Exit(1)
	}

	fmt.Println("[Logout] Logout successful")
	resp.Body.Close()

	fmt.Println("### SHBF SM registration process completed ###")
}
