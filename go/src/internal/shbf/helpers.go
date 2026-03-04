package shbf

import (
	"fmt"
	"net/http"
	"net/url"
	"shbfsmreg/internal/beerxml"
	"shbfsmreg/internal/clientsession"
	"strconv"
	"strings"
)

func Login(cs *clientsession.ClientSession, username, password string) (*http.Response, error) {
	// Form data (application/x-www-form-urlencoded)
	form := url.Values{}
	form.Set("user_name", username)
	form.Set("passwd", password)
	form.Set("submit", "Logga in")
	form.Set("email", "")

	req, err := http.NewRequest("POST", "https://event.shbf.se/login.php", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	// Headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://event.shbf.se")
	req.Header.Set("Referer", "https://event.shbf.se/login.php")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Send request
	resp, err := cs.Client.Do(req)
	if err != nil {
		return nil, err
	}
	// Caller must close resp.Body

	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to login: %d %s", resp.StatusCode, resp.Status)
	}

	return resp, nil
}

func Logout(cs *clientsession.ClientSession) (*http.Response, error) {

	req, err := http.NewRequest("GET", "https://event.shbf.se/logout.php", nil)
	if err != nil {
		return nil, err
	}
	// Headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://event.shbf.se")
	req.Header.Set("Referer", "https://event.shbf.se/login.php")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Send request
	resp, err := cs.Client.Do(req)
	if err != nil {
		return nil, err
	}
	// Caller must close resp.Body

	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to logout: %d %s", resp.StatusCode, resp.Status)
	}

	return resp, nil
}

func SelectEvent(cs *clientsession.ClientSession, eventId, fvEventId int) (*http.Response, error) {

	req, err := http.NewRequest("GET", "https://event.shbf.se/set_sessions.php?dt_event_id="+strconv.Itoa(eventId)+"&fv_event_id="+strconv.Itoa(fvEventId), nil)
	if err != nil {
		return nil, err
	}
	// Headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://event.shbf.se")
	req.Header.Set("Referer", "https://event.shbf.se/login.php")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Send request
	resp, err := cs.Client.Do(req)
	if err != nil {
		return nil, err
	}
	// Caller must close resp.Body

	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to select event: %d %s", resp.StatusCode, resp.Status)
	}

	return resp, nil
}

func RegisterBeerPre(cs *clientsession.ClientSession) (*http.Response, error) {

	req, err := http.NewRequest("GET", "https://event.shbf.se/beer_reg_pre.php", nil)
	if err != nil {
		return nil, err
	}
	// Headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://event.shbf.se")
	req.Header.Set("Referer", "https://event.shbf.se/beer_reg_pre.php")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Send request
	resp, err := cs.Client.Do(req)
	if err != nil {
		return nil, err
	}
	// Caller must close resp.Body

	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to register beer pre: %d %s", resp.StatusCode, resp.Status)
	}

	return resp, nil
}

// hopFormID maps BeerXML FORM to SHBF dropdown: 1=Kottar, 2=Pellets, 3=Färsk, 4=Hel, 5=Krossad, 6=Mald
func hopFormID(form string) string {
	f := strings.ToLower(strings.TrimSpace(form))
	switch f {
	case "pellet", "pellets":
		return "2"
	case "plug", "plugs", "leaf", "whole":
		return "4"
	case "cone", "kottar":
		return "1"
	case "fresh":
		return "3"
	case "extract", "mald":
		return "6"
	default:
		return "0"
	}
}

// miscStageID maps BeerXML USE to form stage: 1=Mäskning, 2=Kokning, 3=Jäsning, 4=Lagring
func miscStageID(use string) string {
	u := strings.ToLower(strings.TrimSpace(use))
	switch {
	case u == "mash":
		return "1"
	case u == "boil":
		return "2"
	case strings.Contains(u, "primary") || strings.Contains(u, "ferment"):
		return "3"
	case strings.Contains(u, "secondary") || strings.Contains(u, "lagering") || strings.Contains(u, "lagring"):
		return "4"
	default:
		return "0"
	}
}

// buildMashing returns mashing text from recipe mash steps.
func buildMashing(r *beerxml.Recipe) string {
	if len(r.Mash.Steps.MashSteps) == 0 {
		return "Proteinrast xx °C, xx min\r\nFörsockringsrast xx °C, xx min"
	}
	var parts []string
	for _, s := range r.Mash.Steps.MashSteps {
		parts = append(parts, s.Name+" "+s.StepTemp+" °C, "+s.StepTime+" min")
	}
	return strings.Join(parts, "\r\n")
}

// buildFerment returns ferment text from first yeast.
func buildFerment(r *beerxml.Recipe) string {
	if len(r.Yeasts.Yeast) == 0 {
		return "Jästsort: xxx\r\nJäsning xx °C, xx dagar"
	}
	y := r.Yeasts.Yeast[0]
	name := y.ProductID
	if name == "" {
		name = y.Name
	}
	if name == "" {
		name = "xxx"
	}
	temp := "xx"
	if y.MinTemperature != "" && y.MaxTemperature != "" {
		minTemperatureFloat, _ := strconv.ParseFloat(y.MinTemperature, 64)
		maxTemperatureFloat, _ := strconv.ParseFloat(y.MaxTemperature, 64)
		y.MinTemperature = strconv.FormatFloat(minTemperatureFloat, 'f', 1, 64)
		y.MaxTemperature = strconv.FormatFloat(maxTemperatureFloat, 'f', 1, 64)
		temp = strconv.FormatFloat(minTemperatureFloat, 'f', 1, 64) + "-" + strconv.FormatFloat(maxTemperatureFloat, 'f', 1, 64)
	}
	return "Jästsort: " + name + "\r\nJäsning " + temp + " °C, 14 dagar"
}

// buildWater returns water text from first water profile.
func buildWater(r *beerxml.Recipe) string {
	if len(r.Waters.Water) == 0 {
		return "Vatten från xx-stad\r\nTillsatser:\r\nxx g någonting"
	}
	return "Vatten från " + r.Waters.Water[0].Name + "\r\nTillsatser:\r\nxx g någonting"
}

// RegisterBeer POSTs the beer registration form to event.shbf.se/beer_reg.php.
// recipe should already be SHBF-compatible (e.g. after beerxml.MakeSHBFCompatible).
// ingredientLimit is the max number of rows for hops, malts, and others (e.g. 10).
func RegisterBeer(cs *clientsession.ClientSession, recipe *beerxml.Recipe, brewerName, brewerEmail string, ingredientLimit int) (*http.Response, error) {

	// Register beer pre, some session magic happens here
	resp, err := RegisterBeerPre(cs)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to get registration pre page: %d %s", resp.StatusCode, resp.Status)
	}
	resp.Body.Close()

	// Build form
	form, err := CreateBeerRegistrationForm(recipe, brewerName, brewerEmail, ingredientLimit)
	if err != nil {
		return nil, err
	}

	// Send request
	req, err := http.NewRequest("POST", "https://event.shbf.se/beer_reg.php", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://event.shbf.se")
	req.Header.Set("Referer", "https://event.shbf.se/beer_reg.php")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err = cs.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to register beer: %d %s", resp.StatusCode, resp.Status)
	}
	resp.Body.Close()

	return resp, nil
}

func CreateBeerRegistrationForm(recipe *beerxml.Recipe, brewerName, brewerEmail string, ingredientLimit int) (url.Values, error) {

	// Build form
	form := url.Values{}
	form.Set("comp_dt", "1")
	form.Set("comp_fv", "0")
	form.Set("beer_type", recipe.Style.CategoryNumber+":"+recipe.Style.StyleLetter)
	form.Set("beer_name", recipe.Name)
	form.Set("og", recipe.OG)
	form.Set("fg", recipe.FG)
	form.Set("bu", recipe.IBU)
	form.Set("alc", recipe.ABV)
	form.Set("volume", recipe.BatchSize)
	form.Add("brewers_name[]", brewerName)
	form.Add("brewers_email[]", brewerEmail)

	// Hops (ingredientLimit slots)
	for i := 0; i < ingredientLimit; i++ {
		var name, alpha, weight, boilTime, comment, formID string
		if i < len(recipe.Hops.Hop) {
			h := &recipe.Hops.Hop[i]
			name = h.Name
			alpha = h.Alpha
			weight = h.Amount
			boilTime = h.Time
			comment = ""
			formID = hopFormID(h.Form)
		} else {
			formID = "0"
		}

		if name != "" {
			form.Add("hops_name[]", name)
			form.Add("hops_form_id_sel[]", formID)
			form.Add("hops_alpha[]", alpha)
			form.Add("hops_weight[]", weight)
			form.Add("hops_boil_time[]", boilTime)
			form.Add("hops_comment[]", comment)
		}
	}
	// Malts (ingredientLimit slots)
	for i := 0; i < ingredientLimit; i++ {
		var name, weight, comment string
		if i < len(recipe.Fermentables.Fermentable) {
			m := &recipe.Fermentables.Fermentable[i]
			name = m.Name
			weight = m.Amount
			comment = ""

			form.Add("malts_name[]", name)
			form.Add("malts_weight[]", weight)
			form.Add("malts_comment[]", comment)
		}

	}
	form.Set("mashing", buildMashing(recipe))
	// Others/Misc (ingredientLimit slots)
	for i := 0; i < ingredientLimit; i++ {
		var name, weight, comment, stageID string
		if i < len(recipe.Miscs.Misc) {
			m := &recipe.Miscs.Misc[i]
			name = m.Name
			weight = m.Amount
			comment = ""
			stageID = miscStageID(m.Use)
		} else {
			stageID = "0"
		}

		if name != "" {
			form.Add("others_name[]", name)
			form.Add("others_stage_id_sel[]", stageID)
			form.Add("others_weight[]", weight)
			form.Add("others_comment[]", comment)
		}
	}
	comment := recipe.Notes
	if comment == "" {
		comment = ""
	}

	form.Set("ferment", buildFerment(recipe))
	form.Set("water", buildWater(recipe))
	form.Set("comment", comment)
	form.Set("add_beer", "Spara")

	return form, nil
}
