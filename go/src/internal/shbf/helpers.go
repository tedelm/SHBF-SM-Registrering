package shbf

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"shbfsmreg/internal/beerxml"
	"shbfsmreg/internal/clientsession"
	"strconv"
	"strings"
)

func Login(cs *clientsession.ClientSession, username, password string) error {
	// Form data (application/x-www-form-urlencoded)
	form := url.Values{}
	form.Set("user_name", username)
	form.Set("passwd", password)
	form.Set("submit", "Logga in")
	form.Set("email", "")

	req, err := http.NewRequest("POST", "https://event.shbf.se/login.php", strings.NewReader(form.Encode()))
	if err != nil {
		return err
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
		return err
	}
	// Caller must close resp.Body

	if resp.StatusCode != 200 {
		resp.Body.Close()
		return fmt.Errorf("failed to login: %d %s", resp.StatusCode, resp.Status)
	}

	isLoggedIn, err := CheckIfLoggedIn(cs)
	if err != nil || !isLoggedIn {
		return fmt.Errorf("Error: %w", err)
	}

	return nil
}

func CheckIfLoggedIn(cs *clientsession.ClientSession) (bool, error) {
	// Form data (application/x-www-form-urlencoded)
	req, err := http.NewRequest("GET", "https://event.shbf.se/message.php?message=checkcheck.", nil)
	if err != nil {
		return false, err
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
		return false, err
	}
	// Caller must close resp.Body

	if resp.StatusCode != 200 {
		resp.Body.Close()
		return false, fmt.Errorf("failed to login: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	resp.Body.Close()
	bodyStr := string(body)
	if strings.Contains(bodyStr, "<title>Eventregistrering - (") {
		return true, nil
	}

	return false, fmt.Errorf("Not logged in, check credentials")
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

// UpdateBeerPre fetches the event registration list (event_reg.php), finds the row matching beerNameSubstring
// (e.g. recipe name like "Soleil Sauvage"), and returns the beer_id from that row's "Ändra" link.
func UpdateBeerPre(cs *clientsession.ClientSession, beerNameSubstring string) (string, error) {
	req, err := http.NewRequest("GET", "https://event.shbf.se/event_reg.php", nil)
	if err != nil {
		return "", err
	}
	// Headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://event.shbf.se")
	req.Header.Set("Referer", "https://event.shbf.se/beer_reg.php")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := cs.Client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return "", fmt.Errorf("failed to get event list: %d %s", resp.StatusCode, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return "", err
	}
	resp.Body.Close()
	bodyStr := string(body)
	beerID, err := parseBeerIDFromEventList(bodyStr, beerNameSubstring)
	if err != nil {
		return "", err
	}

	// set session beer_id
	req, err = http.NewRequest("GET", "https://event.shbf.se/beer_reg_pre.php?beer_id="+beerID, nil)
	if err != nil {
		return "", err
	}
	// Headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://event.shbf.se")
	req.Header.Set("Referer", "https://event.shbf.se/event_reg.php")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err = cs.Client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return "", fmt.Errorf("failed to get event list: %d %s", resp.StatusCode, resp.Status)
	}

	return beerID, nil
}

// parseBeerIDFromEventList finds a <tr> row that contains nameSubstring (e.g. beer name like "Soleil Sauvage")
// and extracts beer_id=XXXX from the "Ändra" link (beer_reg_pre.php?beer_id=XXXX) in that row.
func parseBeerIDFromEventList(body, nameSubstring string) (string, error) {
	if nameSubstring == "" {
		return "", fmt.Errorf("beer name substring required to find beer_id")
	}
	idx := strings.Index(body, nameSubstring)
	if idx == -1 {
		return "", fmt.Errorf("no row containing %q found", nameSubstring)
	}
	rowStart := strings.LastIndex(body[:idx], "<tr")
	if rowStart == -1 {
		rowStart = 0
	}
	rowEnd := strings.Index(body[idx:], "</tr>")
	if rowEnd == -1 {
		rowEnd = len(body) - idx
	}
	row := body[rowStart : idx+rowEnd]
	// In this row find beer_id= followed by digits (from link beer_reg_pre.php?beer_id=4761)
	bid := strings.Index(row, "beer_id=")
	if bid == -1 {
		return "", fmt.Errorf("row containing %q has no beer_id link", nameSubstring)
	}
	bid += len("beer_id=")
	end := bid
	for end < len(row) && row[end] >= '0' && row[end] <= '9' {
		end++
	}
	if end == bid {
		return "", fmt.Errorf("beer_id= not followed by digits in row containing %q", nameSubstring)
	}
	return row[bid:end], nil
}

// hopFormID maps BeerXML FORM to SHBF dropdown: 1=Kottar, 2=Pellets
func hopFormID(form string) string {
	f := strings.ToLower(strings.TrimSpace(form))
	switch f {
	case "pellet", "pellets":
		return "2"
	case "plug", "plugs", "leaf", "whole", "cone", "kottar":
		return "1"
	default:
		return "2"
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
func RegisterBeer(cs *clientsession.ClientSession, recipe *beerxml.Recipe, brewerName, brewerEmail string, ingredientLimit int, updateBeer bool) (*http.Response, error) {

	var resp *http.Response
	var err error
	if updateBeer {
		beerID, err := UpdateBeerPre(cs, recipe.Name)
		if err != nil {
			return nil, err
		}
		if beerID == "" {
			return nil, fmt.Errorf("failed to update beer pre, beer not found")
		}
	} else {
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
	}
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
			comment = h.Use

			if h.Use == "Dry Hop" {
				boilTime = "0"
				boilTimeFloat, _ := strconv.ParseFloat(h.Time, 64)
				comment = comment + " (" + strconv.Itoa(int(math.Round(boilTimeFloat/60/24))) + " dagar)"
			}

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

// parseStyleFloat extracts the first number from a string (e.g. "1.063", "27.7 IBUs", "5.9 %", "79.3 EBC").
// Returns (value, true) if a number was parsed, (0, false) otherwise.
func parseStyleFloat(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' || r == '.' && !strings.Contains(b.String(), ".") {
			b.WriteRune(r)
		} else if b.Len() > 0 {
			break
		}
	}
	if b.Len() == 0 {
		return 0, false
	}
	v, err := strconv.ParseFloat(b.String(), 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

// VerifyBeerStyle checks that the recipe's OG, FG, IBU, EBC (color), and ABV fall within the
// style's MIN/MAX bounds from the BeerXML STYLE element. Missing or zero style bounds are skipped.
func VerifyBeerStyle(recipe *beerxml.Recipe) (bool, error) {
	// OG
	if recipe.Style.OGMin != "" || recipe.Style.OGMax != "" {
		og, ok := parseStyleFloat(recipe.OG)
		if !ok {
			og, _ = parseStyleFloat(recipe.EstOG)
		}
		if min, okMin := parseStyleFloat(recipe.Style.OGMin); okMin && ok && og < min {
			return false, fmt.Errorf("OG %g is below style minimum %g", og, min)
		}
		if max, okMax := parseStyleFloat(recipe.Style.OGMax); okMax && ok && og > max {
			return false, fmt.Errorf("OG %g is above style maximum %g", og, max)
		}
	}

	// FG
	if recipe.Style.FGMin != "" || recipe.Style.FGMax != "" {
		fg, ok := parseStyleFloat(recipe.FG)
		if !ok {
			fg, _ = parseStyleFloat(recipe.EstFG)
		}
		if min, okMin := parseStyleFloat(recipe.Style.FGMin); okMin && ok && fg < min {
			return false, fmt.Errorf("FG %g is below style minimum %g", fg, min)
		}
		if max, okMax := parseStyleFloat(recipe.Style.FGMax); okMax && ok && fg > max {
			return false, fmt.Errorf("FG %g is above style maximum %g", fg, max)
		}
	}

	// IBU
	if recipe.Style.IBUMin != "" || recipe.Style.IBUMax != "" {
		ibu, ok := parseStyleFloat(recipe.IBU)
		if !ok {
			ibu, _ = parseStyleFloat(recipe.EstIBU)
		}
		if min, okMin := parseStyleFloat(recipe.Style.IBUMin); okMin && ok && ibu < min {
			return false, fmt.Errorf("IBU %g is below style minimum %g", ibu, min)
		}
		if max, okMax := parseStyleFloat(recipe.Style.IBUMax); okMax && ok && ibu > max {
			return false, fmt.Errorf("IBU %g is above style maximum %g", ibu, max)
		}
	}

	// EBC (color)
	if (recipe.Style.ColorMin != "" || recipe.Style.ColorMax != "") && strings.TrimSpace(recipe.EstColor) != "" {
		ebc, ok := parseStyleFloat(recipe.EstColor)
		if ok {
			if min, okMin := parseStyleFloat(recipe.Style.ColorMin); okMin && ebc < min {
				return false, fmt.Errorf("EBC (color) %g is below style minimum %g", ebc, min)
			}
			if max, okMax := parseStyleFloat(recipe.Style.ColorMax); okMax && ebc > max {
				return false, fmt.Errorf("EBC (color) %g is above style maximum %g", ebc, max)
			}
		}
	}

	// ABV
	if recipe.Style.ABVMin != "" || recipe.Style.ABVMax != "" {
		abv, ok := parseStyleFloat(recipe.ABV)
		if !ok {
			abv, _ = parseStyleFloat(recipe.EstABV)
		}
		if min, okMin := parseStyleFloat(recipe.Style.ABVMin); okMin && ok && abv < min {
			return false, fmt.Errorf("ABV %g%% is below style minimum %g%%", abv, min)
		}
		if max, okMax := parseStyleFloat(recipe.Style.ABVMax); okMax && ok && abv > max {
			return false, fmt.Errorf("ABV %g%% is above style maximum %g%%", abv, max)
		}
	}

	return true, nil
}
