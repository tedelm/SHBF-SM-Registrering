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
	flag.Parse()

	if *username == "" || *password == "" || *eventId == 0 || *fvEventId == 0 || *beerXmlPath == "" || *brewerName == "" || *brewerEmail == "" {
		fmt.Println("Missing required parameters")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Read BeerXML file
	beerXml, err := beerxml.ParseFile(*beerXmlPath)
	if err != nil {
		fmt.Println("Failed to parse BeerXML file: ", err)
		os.Exit(1)
	}

	beerXml, err = beerxml.MakeSHBFCompatible(beerXml)
	if err != nil {
		fmt.Println("Failed to make BeerXML compatible with SHBF: ", err)
		os.Exit(1)
	}

	fmt.Println("BeerXML file parsed successfully")

	_, err = shbf.CreateBeerRegistrationForm(beerXml, *brewerName, *brewerEmail, *ingredientLimit)
	if err != nil {
		fmt.Println("Failed to create beer registration form: ", err)
		os.Exit(1)
	}
	fmt.Println("Beer registration form created successfully")

	// Create cookie jar (equivalent to WebRequestSession)
	cs, err := clientsession.NewClientSession()
	if err != nil {
		fmt.Println("Failed to create client session: ", err)
		os.Exit(1)
	}

	// Login
	resp, err := shbf.Login(cs, *username, *password)
	if err != nil {
		fmt.Println("Failed to login: ", err)
		os.Exit(1)
	}

	// Check if login was successful
	if resp.StatusCode != 200 {
		fmt.Println("Failed to login: ", resp.StatusCode, resp.Status)
		os.Exit(1)
	}

	fmt.Println("Login successful")
	fmt.Println("Status:", resp.StatusCode, resp.Status)
	resp.Body.Close()

	// Select event
	resp, err = shbf.SelectEvent(cs, *eventId, *fvEventId)
	if err != nil {
		fmt.Println("Failed to select event: ", err)
		os.Exit(1)
	}

	fmt.Println("Select event successful")
	fmt.Println("Status:", resp.StatusCode, resp.Status)
	resp.Body.Close()

	// Register beer
	resp, err = shbf.RegisterBeer(cs, beerXml, *brewerName, *brewerEmail, *ingredientLimit)
	if err != nil {
		fmt.Println("Failed to register beer: ", err)
		os.Exit(1)
	}
	fmt.Println("Beer registered successfully")
	fmt.Println("Status:", resp.StatusCode, resp.Status)
	resp.Body.Close()

	// Logout
	resp, err = shbf.Logout(cs)
	if err != nil {
		fmt.Println("Failed to logout: ", err)
		os.Exit(1)
	}

	fmt.Println("Logout successful")
	fmt.Println("Status:", resp.StatusCode, resp.Status)
	resp.Body.Close()

}
