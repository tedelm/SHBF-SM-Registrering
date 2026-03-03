<#
SHBF SM Registeration

This script is used to register a beer for the SHBF SM.
It uses the BeerXML format to import the beer information.
It uses the BeerSmith export format to import the beer information.
It uses the SHBF event ID to set the event and fv event.
Eventids for the 2026 SM are:
- Event ID: 61
- FV Event ID: 62
#>

param(
    [int] $CompDt = 1,  # 1 = include comp_dt=1 (Domartavlingen), 0 = omit
    [int] $CompFv = 0,   # 1 = include comp_fv=1 (Folkets val), 0 = omit
    [parameter(mandatory = $true)] [string] $username,
    [parameter(mandatory = $true)] [string] $password,
    [parameter(mandatory = $true)] [string] $BeerXmlPath,
    [parameter(mandatory = $true)] [string] $brewersName,
    [parameter(mandatory = $true)] [string] $brewersEmail,
    [parameter(mandatory = $true)] [string] $shbfEventId,
    [parameter(mandatory = $true)] [string] $shbfFvEventId
)

function New-SHBFSession {
    param(
        [parameter(mandatory = $true)] [string] $username,
        [parameter(mandatory = $true)] [string] $password,
        [parameter(mandatory = $true)] [string] $eventId,
        [parameter(mandatory = $true)] [string] $fvEventId
    )

    $eventUrl = "https://event.shbf.se/set_sessions.php?dt_event_id=$eventId&fv_event_id=$fvEventId"
    $regPreUrl = "https://event.shbf.se/beer_reg_pre.php"
    $regUrl = "https://event.shbf.se/beer_reg.php"

    $session = New-Object Microsoft.PowerShell.Commands.WebRequestSession
    $session.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36"
    #$session.Cookies.Add((New-Object System.Net.Cookie("regrate", "rt41i1sfenhbbs132m07uko6g0", "/", ".shbf.se")))
    $response = Invoke-WebRequest -UseBasicParsing -Uri "https://event.shbf.se/login.php" `
    -Method "POST" `
    -WebSession $session `
    -Headers @{
    "authority"="event.shbf.se"
    "method"="POST"
    "path"="/login.php"
    "scheme"="https"
    "accept"="text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"
    "accept-encoding"="gzip, deflate, br, zstd"
    "accept-language"="en-US,en;q=0.9"
    "cache-control"="max-age=0"
    "origin"="https://event.shbf.se"
    "priority"="u=0, i"
    "referer"="https://event.shbf.se/login.php"
    "sec-ch-ua"="`"Not:A-Brand`";v=`"99`", `"Google Chrome`";v=`"145`", `"Chromium`";v=`"145`""
    "sec-ch-ua-mobile"="?0"
    "sec-ch-ua-platform"="`"Windows`""
    "sec-fetch-dest"="document"
    "sec-fetch-mode"="navigate"
    "sec-fetch-site"="same-origin"
    "sec-fetch-user"="?1"
    "upgrade-insecure-requests"="1"
    } `
    -ContentType "application/x-www-form-urlencoded" `
    -Body "user_name=$username&passwd=$password&submit=Logga+in&email=";

    if ($response.StatusCode -ne 200) {
        Write-Error "Failed to login: $($response.StatusCode) $($response.StatusDescription)"
        return $null
    }

    $response = Invoke-WebRequest -UseBasicParsing -Uri $eventUrl -WebSession $session
    if ($response.StatusCode -ne 200) {
        Write-Error "Failed to get event page ($eventUrl): $($response.StatusCode) $($response.StatusDescription)"
        return $null
    }

    $response = Invoke-WebRequest -UseBasicParsing -Uri $regPreUrl -WebSession $session
    if ($response.StatusCode -ne 200) {
        Write-Error "Failed to get registration pre page ($regPreUrl): $($response.StatusCode) $($response.StatusDescription)"
        return $null
    }

    $response = Invoke-WebRequest -UseBasicParsing -Uri $regUrl -WebSession $session
    if ($response.StatusCode -ne 200) {
        Write-Error "Failed to get registration page ($regUrl): $($response.StatusCode) $($response.StatusDescription)"
        return $null
    }

    return [Microsoft.PowerShell.Commands.WebRequestSession]$session
}

$session = New-SHBFSession -username $username -password $password -eventId $shbfEventId -fvEventId $shbfFvEventId

if (-not $session) {
    Write-Error "Failed to create session"
    exit 1
}


[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

$url = "https://event.shbf.se/beer_reg.php"

$Headers = @{
  "authority" = "event.shbf.se"
  "method" = "POST"
  "path" = "/beer_reg.php"
  "scheme" = "https"
  "accept" = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"
  "accept-encoding" = "gzip, deflate, br, zstd"
  "accept-language" = "en-US,en;q=0.9,de;q=0.8,sv;q=0.7"
  "cache-control" = "max-age=0"
  "origin" = "https://event.shbf.se"
  "priority" = "u=0, i"
  "referer" = "https://event.shbf.se/beer_reg.php"
  "sec-ch-ua" = "`"Not:A-Brand`";v=`"99`", `"Google Chrome`";v=`"145`", `"Chromium`";v=`"145`""
  "sec-ch-ua-mobile" = "?0"
  "sec-ch-ua-platform" = "`"Windows`""
  "sec-fetch-dest" = "document"
  "sec-fetch-mode" = "navigate"
  "sec-fetch-site" = "same-origin"
  "sec-fetch-user" = "?1"
  "upgrade-insecure-requests" = "1"
}

function Get-BeerRegFormBody {
    [CmdletBinding()]
    param(
        [int] $CompDt = 1,
        [int] $CompFv = 0,
        [parameter(mandatory = $true)] [string] $BrewersName,
        [parameter(mandatory = $true)] [string] $BrewersEmail,
        [Parameter(ParameterSetName = "Manual")][string] $BeerName = "",
        [Parameter(ParameterSetName = "Manual")][string] $BeerType = "",
        [Parameter(ParameterSetName = "Manual")][string] $Og = "",
        [Parameter(ParameterSetName = "Manual")][string] $Fg = "",
        [Parameter(ParameterSetName = "Manual")][string] $Bu = "",
        [Parameter(ParameterSetName = "Manual")][string] $Alc = "",
        [Parameter(ParameterSetName = "Manual")][string] $Volume = "",
        [Parameter(ParameterSetName = "Manual")][string] $Mashing = "",
        [Parameter(ParameterSetName = "Manual")][string] $Ferment = "",
        [Parameter(ParameterSetName = "Manual")][string] $Water = "",
        [Parameter(ParameterSetName = "Manual")][string] $Comment = "",
        # BeerXML: path to .xml file or raw XML string
        [Parameter(ParameterSetName = "BeerXml", Mandatory = $true)]
        [string] $BeerXmlPath,
        # If set, POST the form body to the registration URL (use with -WebSession from script)
        [switch] $Post,
        # Web session (cookies) to use when -Post is set; default uses script $session if available
        [Microsoft.PowerShell.Commands.WebRequestSession] $WebSession
    )
    Add-Type -AssemblyName System.Web -ErrorAction SilentlyContinue

    # If BeerXML path provided, parse and override from recipe
    if ($PSCmdlet.ParameterSetName -eq "BeerXml") {
        $xmlContent = if (Test-Path -LiteralPath $BeerXmlPath -PathType Leaf) {
            Get-Content -Path $BeerXmlPath -Raw -Encoding UTF8
        } else {
            $BeerXmlPath
        }
        $doc = [xml]$xmlContent
        $ns = New-Object System.Xml.XmlNamespaceManager $doc.NameTable
        $recipe = $doc.SelectSingleNode("//*[local-name()='RECIPE']")
        if (-not $recipe) { throw "No RECIPE found in BeerXML." }

        $get = { param($n) $node = $recipe.SelectSingleNode(".//*[local-name()='$n']"); if ($node) { $node.InnerText.Trim() } else { "" } }

        $BeerName = & $get "NAME"
        # beer_type = dropdown value "category:letter" from STYLE (e.g. 9:J for Supersaison)
        $styleNode = $recipe.SelectSingleNode(".//*[local-name()='STYLE']")
        $catNum = ""; $styleLetter = ""
        if ($styleNode) {
            $cn = $styleNode.SelectSingleNode("*[local-name()='CATEGORY_NUMBER']")
            $sl = $styleNode.SelectSingleNode("*[local-name()='STYLE_LETTER']")
            if ($cn) { $catNum = $cn.InnerText.Trim() }
            if ($sl) { $styleLetter = $sl.InnerText.Trim() }
        }
        if ($catNum -and $styleLetter) { $BeerType = "${catNum}:${styleLetter}" } else { $BeerType = "" }
        $OgRaw = & $get "OG"
        $FgRaw = & $get "FG"
        # Form expects OG/FG in g/l (SG * 1000), integer, no decimals
        $ogNum = 0; if ([double]::TryParse($OgRaw, [System.Globalization.NumberStyles]::Any, [System.Globalization.CultureInfo]::InvariantCulture, [ref]$ogNum) -and $ogNum -gt 0) { $Og = [string][int][Math]::Round($ogNum * 1000) } else { $Og = "" }
        $fgNum = 0; if ([double]::TryParse($FgRaw, [System.Globalization.NumberStyles]::Any, [System.Globalization.CultureInfo]::InvariantCulture, [ref]$fgNum) -and $fgNum -gt 0) { $Fg = [string][int][Math]::Round($fgNum * 1000) } else { $Fg = "" }
        $Volume = & $get "BATCH_SIZE"
        if (-not $Volume) { $Volume = & $get "DISPLAY_BATCH_SIZE" }
        # Volume as integer (liters)
        $volNum = 0; if ([double]::TryParse($Volume, [System.Globalization.NumberStyles]::Any, [System.Globalization.CultureInfo]::InvariantCulture, [ref]$volNum)) { $Volume = [string][int][Math]::Round($volNum) } else { $Volume = "" }
        # BU: BeerXML often has "27.9 IBUs" – form expects integer
        $buRaw = & $get "IBU"; if (-not $buRaw) { $buRaw = & $get "EST_IBU" }
        $buNum = 0; if ([double]::TryParse(($buRaw -replace '[^\d\.].*', '').Trim(), [System.Globalization.NumberStyles]::Any, [System.Globalization.CultureInfo]::InvariantCulture, [ref]$buNum)) { $Bu = [string][int][Math]::Round($buNum) } else { $Bu = "" }
        # Alc: BeerXML often has "8.6 %" – form expects number, max 10.0
        $alcRaw = & $get "ABV"; if (-not $alcRaw) { $alcRaw = & $get "EST_ABV" }
        $alcNum = 0; if ([double]::TryParse(($alcRaw -replace '[^\d\.].*', '').Trim(), [System.Globalization.NumberStyles]::Any, [System.Globalization.CultureInfo]::InvariantCulture, [ref]$alcNum)) { $Alc = [string][Math]::Min(10.0, [Math]::Round($alcNum, 1)) } else { $Alc = "" }
        $BrewersName = & $get "BREWER"
        if (-not $BrewersName) { $BrewersName = "Ted Elmenheim" }
        # Recipe-level NOTES (direct child); fallback to STYLE/NOTES for style description
        $notesNode = $recipe.SelectSingleNode("*[local-name()='NOTES']")
        $Comment = if ($notesNode) { $notesNode.InnerText.Trim() }
        if (-not $Comment) {
            $styleNotes = $recipe.SelectSingleNode(".//*[local-name()='STYLE']/*[local-name()='NOTES']")
            if ($styleNotes) { $Comment = $styleNotes.InnerText.Trim() }
        }

        function Get-NodeText($node, $tag) {
            $child = $node.SelectSingleNode("*[local-name()='$tag']"); if ($child) { $child.InnerText.Trim() } else { "" }
        }
        # Map BeerXML hop FORM to form dropdown: 1=Kottar, 2=Pellets, 3=Färsk, 4=Hel, 5=Krossad, 6=Mald
        function Get-HopFormId($formText) {
            $f = ($formText -replace "\s+", " ").Trim().ToLowerInvariant()
            if ($f -eq "pellet" -or $f -eq "pellets") { return "2" }
            if ($f -eq "plug" -or $f -eq "plugs") { return "4" }
            if ($f -eq "leaf" -or $f -eq "whole") { return "4" }
            if ($f -eq "cone" -or $f -eq "kottar") { return "1" }
            if ($f -eq "fresh") { return "3" }
            if ($f -eq "extract" -or $f -eq "mald") { return "6" }
            return "0"
        }
        function Get-HopWeightGrams($hopNode) {
            $disp = Get-NodeText $hopNode "DISPLAY_AMOUNT"
            if ($disp -match "(\d+(?:[.,]\d+)?)\s*g") { return [string][int][Math]::Round(([double]($Matches[1] -replace ",", "."))) }
            $amt = Get-NodeText $hopNode "AMOUNT"
            $a = 0; if ([double]::TryParse($amt, [System.Globalization.NumberStyles]::Any, [System.Globalization.CultureInfo]::InvariantCulture, [ref]$a)) { return [string][int][Math]::Round($a * 1000) }
            return ""
        }

        # Hops (form has 4 slots)
        $hopNodes = @($recipe.SelectNodes(".//*[local-name()='HOPS']/*[local-name()='HOP']"))
        $hops = @()
        foreach ($h in $hopNodes) {
            if ($hops.Count -ge 4) { break }
            $hops += [PSCustomObject]@{
                Name = (Get-NodeText $h "NAME"); Alpha = (Get-NodeText $h "ALPHA"); Weight = (Get-HopWeightGrams $h)
                BoilTime = (Get-NodeText $h "TIME"); FormId = (Get-HopFormId (Get-NodeText $h "FORM")); Comment = (Get-NodeText $h "NOTES")
            }
        }

        # Fermentables -> malts (form has 4 slots); weight in grams (BeerXML AMOUNT is kg)
        function Get-MaltWeightGrams($fermNode) {
            $disp = Get-NodeText $fermNode "DISPLAY_AMOUNT"
            if ($disp -match "(\d+(?:[.,]\d+)?)\s*(?:kg|g)") { $val = [double]($Matches[1] -replace ",", "."); if ($disp -match "kg") { $val = $val * 1000 }; return [string][int][Math]::Round($val) }
            $amt = Get-NodeText $fermNode "AMOUNT"
            $a = 0; if ([double]::TryParse($amt, [System.Globalization.NumberStyles]::Any, [System.Globalization.CultureInfo]::InvariantCulture, [ref]$a)) { return [string][int][Math]::Round($a * 1000) }
            return ""
        }
        $fermNodes = @($recipe.SelectNodes(".//*[local-name()='FERMENTABLES']/*[local-name()='FERMENTABLE']"))
        $malts = @()
        foreach ($f in $fermNodes) {
            if ($malts.Count -ge 4) { break }
            $malts += [PSCustomObject]@{
                Name = (Get-NodeText $f "NAME"); Weight = (Get-MaltWeightGrams $f); Comment = (Get-NodeText $f "NOTES")
            }
        }

        # Mash steps -> mashing text
        $mash = $recipe.SelectSingleNode(".//*[local-name()='MASH']")
        if ($mash) {
            $steps = $mash.SelectNodes(".//*[local-name()='MASH_STEP']")
            $Mashing = ($steps | ForEach-Object {
                $sn = Get-NodeText $_ "NAME"; $st = Get-NodeText $_ "STEP_TEMP"; $stm = Get-NodeText $_ "STEP_TIME"
                "${sn} ${st} °C, ${stm} min"
            }) -join "`r`n"
        }
        if (-not $Mashing) { $Mashing = "Proteinrast xx °C, xx min`r`nFörsockringsrast xx °C, xx min" }

        # Yeast -> ferment text; prefer PRODUCT_ID (e.g. M29) for "Jästsort" when present
        $yeastNodes = $recipe.SelectNodes(".//*[local-name()='YEASTS']/*[local-name()='YEAST']")
        if ($yeastNodes -and $yeastNodes.Count -gt 0) {
            $y = $yeastNodes[0]
            $yn = Get-NodeText $y "PRODUCT_ID"; if (-not $yn) { $yn = Get-NodeText $y "NAME" }; if (-not $yn) { $yn = "xxx" }
            $minT = Get-NodeText $y "MIN_TEMPERATURE"; $maxT = Get-NodeText $y "MAX_TEMPERATURE"
            $tempStr = if ($minT -and $maxT) { "${minT}-${maxT}" } else { "xx" }
            $Ferment = "Jästsort: $yn`r`nJäsning ${tempStr} °C, xx dagar"
        }
        if (-not $Ferment) { $Ferment = "Jästsort: xxx`r`nJäsning xx °C, xx dagar" }

        # Water
        $waterNodes = $recipe.SelectNodes(".//*[local-name()='WATERS']/*[local-name()='WATER']")
        if ($waterNodes -and $waterNodes.Count -gt 0) {
            $w = $waterNodes[0]; $wn = Get-NodeText $w "NAME"
            $Water = "Vatten från $wn`r`nTillsatser:`r`nxx g någonting"
        }
        if (-not $Water) { $Water = "Vatten från xx-stad`r`nTillsatser:`r`nxx g någonting" }

        if (-not $Comment) { $Comment = "Här ska t.ex. öltyp anges för underklasserna Övriga klassiska." }

        # Others (MISC): form stage 1=Mäskning, 2=Kokning, 3=Jäsning, 4=Lagring; BeerXML USE: Mash, Boil, Primary, Secondary
        function Get-MiscStageId($useText) {
            $u = ($useText -replace "\s+", " ").Trim().ToLowerInvariant()
            if ($u -eq "mash") { return "1" }
            if ($u -eq "boil") { return "2" }
            if ($u -match "primary|ferment") { return "3" }
            if ($u -match "secondary|lagering|lagring") { return "4" }
            return "0"
        }
        function Get-MiscWeight($miscNode) {
            $isWeight = (Get-NodeText $miscNode "AMOUNT_IS_WEIGHT")
            $amt = Get-NodeText $miscNode "AMOUNT"
            $a = 0; if (-not [double]::TryParse($amt, [System.Globalization.NumberStyles]::Any, [System.Globalization.CultureInfo]::InvariantCulture, [ref]$a)) { return "" }
            if ($isWeight -eq "TRUE" -or $isWeight -eq "true" -or $isWeight -eq "1") { return [string][int][Math]::Round($a * 1000) }
            return [string][int][Math]::Round($a)
        }
        $miscNodes = @($recipe.SelectNodes(".//*[local-name()='MISCS']/*[local-name()='MISC']"))
        $others = @()
        foreach ($m in $miscNodes) {
            if ($others.Count -ge 4) { break }
            $others += [PSCustomObject]@{
                Name = (Get-NodeText $m "NAME"); StageId = (Get-MiscStageId (Get-NodeText $m "USE")); Weight = (Get-MiscWeight $m); Comment = (Get-NodeText $m "NOTES")
            }
        }
    }

    # Defaults when not from BeerXML
    if (-not $Mashing) { $Mashing = "Proteinrast xx °C, xx min`r`nFörsockringsrast xx °C, xx min" }
    if (-not $Ferment) { $Ferment = "Jästsort: xxx`r`nJäsning xx °C, xx dagar" }
    if (-not $Water) { $Water = "Vatten från xx-stad`r`nTillsatser:`r`nxx g någonting" }
    if (-not $Comment) { $Comment = "Här ska t.ex. öltyp anges för underklasserna Övriga klassiska." }

    # Ensure we have 4 hop slots and 4 malt slots (and 4 others)
    if (-not $hops) {
        $hops = 1..4 | ForEach-Object { [PSCustomObject]@{ Name = ""; Alpha = ""; Weight = ""; BoilTime = ""; FormId = "0"; Comment = "" } }
    }
    while ($hops.Count -lt 4) {
        $hops += [PSCustomObject]@{ Name = ""; Alpha = ""; Weight = ""; BoilTime = ""; FormId = "0"; Comment = "" }
    }
    if (-not $malts) {
        $malts = 1..4 | ForEach-Object { [PSCustomObject]@{ Name = ""; Weight = ""; Comment = "" } }
    }
    while ($malts.Count -lt 4) {
        $malts += [PSCustomObject]@{ Name = ""; Weight = ""; Comment = "" }
    }
    if (-not $others) {
        $others = 1..4 | ForEach-Object { [PSCustomObject]@{ Name = ""; StageId = "0"; Weight = ""; Comment = "" } }
    }
    while ($others.Count -lt 4) {
        $others += [PSCustomObject]@{ Name = ""; StageId = "0"; Weight = ""; Comment = "" }
    }

    $parts = @()
    if ($CompDt -eq 1) { $parts += "comp_dt=1" }
    if ($CompFv -eq 1) { $parts += "comp_fv=1" }
    $parts += "beer_type=" + [System.Web.HttpUtility]::UrlEncode($BeerType), "beer_name=" + [System.Web.HttpUtility]::UrlEncode($BeerName), "og=" + [System.Web.HttpUtility]::UrlEncode($Og), "fg=" + [System.Web.HttpUtility]::UrlEncode($Fg), "bu=" + [System.Web.HttpUtility]::UrlEncode($Bu), "alc=" + [System.Web.HttpUtility]::UrlEncode($Alc), "volume=" + [System.Web.HttpUtility]::UrlEncode($Volume), "brewers_name%5B%5D=" + [System.Web.HttpUtility]::UrlEncode($BrewersName), "brewers_email%5B%5D=" + [System.Web.HttpUtility]::UrlEncode($BrewersEmail)
    foreach ($h in $hops[0..3]) {
        $fid = if ($h.FormId) { $h.FormId } else { "0" }
        $parts += "hops_name%5B%5D=" + [System.Web.HttpUtility]::UrlEncode($h.Name), "hops_form_id_sel%5B%5D=" + $fid, "hops_alpha%5B%5D=" + [System.Web.HttpUtility]::UrlEncode($h.Alpha), "hops_weight%5B%5D=" + [System.Web.HttpUtility]::UrlEncode($h.Weight), "hops_boil_time%5B%5D=" + [System.Web.HttpUtility]::UrlEncode($h.BoilTime), "hops_comment%5B%5D=" + [System.Web.HttpUtility]::UrlEncode($h.Comment)
    }
    foreach ($m in $malts[0..3]) {
        $parts += "malts_name%5B%5D=" + [System.Web.HttpUtility]::UrlEncode($m.Name), "malts_weight%5B%5D=" + [System.Web.HttpUtility]::UrlEncode($m.Weight), "malts_comment%5B%5D=" + [System.Web.HttpUtility]::UrlEncode($m.Comment)
    }
    $parts += "mashing=" + [System.Web.HttpUtility]::UrlEncode($Mashing)
    foreach ($o in $others[0..3]) {
        $sid = if ($o.StageId) { $o.StageId } else { "0" }
        $parts += "others_name%5B%5D=" + [System.Web.HttpUtility]::UrlEncode($o.Name), "others_stage_id_sel%5B%5D=" + $sid, "others_weight%5B%5D=" + [System.Web.HttpUtility]::UrlEncode($o.Weight), "others_comment%5B%5D=" + [System.Web.HttpUtility]::UrlEncode($o.Comment)
    }
    $parts += "ferment=" + [System.Web.HttpUtility]::UrlEncode($Ferment), "water=" + [System.Web.HttpUtility]::UrlEncode($Water), "comment=" + [System.Web.HttpUtility]::UrlEncode($Comment), "add_beer=Spara"

    $body = $parts -join "&"
    $body = $body.Replace("%20", "+")
    if ($Post) {
        $sessionToUse = if ($WebSession) { $WebSession } else { $session }
        if (-not $sessionToUse) { throw "Use -WebSession or set script variable `$session before calling with -Post." }
        $postUrl = if ($url) { $url } else { "https://event.shbf.se/beer_reg.php" }
        Invoke-WebRequest -UseBasicParsing -Uri $postUrl -WebSession $sessionToUse -Method POST -ContentType "application/x-www-form-urlencoded" -Body $body -Headers $headers
    } else {
        $body
    }
}

# Build form body from BeerXML, then POST once (CompDt determines if comp_dt=1 is sent)
$body = Get-BeerRegFormBody -BeerXmlPath $BeerXmlPath -BrewersName "$brewersName" -BrewersEmail "$brewersEmail" -CompDt $CompDt -CompFv $CompFv -WebSession $session
$body = $body.Replace(" ", "&")
$response = Invoke-WebRequest -UseBasicParsing -Uri $url -WebSession $session -Method POST -ContentType "application/x-www-form-urlencoded" -Body $body -Headers $Headers

$body | Out-File -FilePath '.\body.txt' -Encoding UTF8 -Force
$response.Content | Out-File -FilePath '.\response.txt' -Encoding UTF8 -Force


Write-Host "Status: $($response.StatusCode) $($response.StatusDescription)"
$hasError = $false
if ($response.Content -match '<pre\s+class="error"[^>]*>([\s\S]*?)</pre>') {
    $errorText = $Matches[1].Trim()
    if ($errorText) {
        $hasError = $true
        Write-Host ""
        Write-Host "Validation errors from server:" -ForegroundColor Red
        $errorText -split "`r?`n" | ForEach-Object { $line = $_.Trim(); if ($line) { Write-Host "  $line" -ForegroundColor Red } }
    }
}
if (-not $hasError) {
    if ($response.StatusCode -ge 200 -and $response.StatusCode -lt 300) {
        if ($response.Content -match 'add_beer|name="add_beer"') { Write-Host "OK: Registration form loaded (add beer)." }
        else { Write-Host "OK: Page loaded." }
    } else {
        Write-Host $response.Content.Substring(0, [Math]::Min(500, $response.Content.Length))
    }
}
if ($hasError) { exit 1 }
