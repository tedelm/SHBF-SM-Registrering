# SHBF-SM-Registrering

Registrera tävlande öl till SHBF SM via BeerXML-export (t.ex. från BeerSmith).  
*[English version](README.en.md)* Du kan köra antingen **PowerShell-scriptet** (kräver endast Windows PowerShell) eller **Go-programmet** (kräver att Go är installerat). Båda loggar in på event.shbf.se och skickar in en ölregistrering med data från din BeerXML-fil.

**Kör Go-binärerna (Se releaser för färdiga .exe-filer):**  
- **Windows:** `.\shbfsmreg.exe -username pelle -password mittlösenord -beerxmlpath .\mitt_recept.xml -brewername pelle -breweremail pelle@epost.se` 
- **Linux:** `./shbfsmreg_linux -username pelle -password mittlösenord -beerxmlpath .\mitt_recept.xml -brewername pelle -breweremail pelle@epost.se` 
- **macOS:** `./shbfsmreg_darwin -username pelle -password mittlösenord -beerxmlpath .\mitt_recept.xml -brewername pelle -breweremail pelle@epost.se`

Se [Parametrar](#parametrar) för alla flaggor.

## Krav

- **PowerShell:** Windows PowerShell (inget extra behov för scriptet)
- **Go-program:** [Go (Golang)](https://go.dev/) 1.21+ (se [Installera Go](#installera-go) nedan)
- Ett BeerXML-recept (export från BeerSmith eller annat verktyg som stöder BeerXML)
- SHBF-konto (användarnamn och lösenord för event.shbf.se)
- Event-ID och FV Event-ID för aktuell SM (se nedan)

## Begränsningar

Registreringsformuläret på event.shbf.se har ett begränsat antal fält. Scriptet importerar **högst ett konfigurerbart antal rader** (standard 10) av varje typ från BeerXML; övriga rader används inte.

| Typ    | Standard max | Kommentar |
|--------|----------------|-----------|
| Malt   | 10             | Antal rader styr med parametern `IngredientLimit` (PowerShell) eller `-ingredientlimit` (Go). |
| Humle  | 10             | Samma som ovan. |
| Övrigt | 10             | Samma som ovan. |

Vid fler ingredienser än gränsen kan du öka `IngredientLimit` / `-ingredientlimit`, eller redigera registreringen manuellt på event.shbf.se efter import.

---

## Installera Go (för Go-programmet)

Om du vill bygga och köra Go-versionen behöver du Go 1.21 eller senare.

- **Windows:** Ladda ner [installer från go.dev](https://go.dev/dl/) (t.ex. `go1.21.x.windows-amd64.msi`) och kör den. Alternativt: [Chocolatey](https://chocolatey.org/) – `choco install golang`.
- **Linux:** T.ex. `sudo apt install golang-go` (Debian/Ubuntu) eller `sudo dnf install golang` (Fedora). Kontrollera version med `go version`.
- **macOS:** `brew install go`.

Kontrollera installationen:

```powershell
go version
```

---

## Bygga Go-programmet

Källkoden ligger under `go/src`. Bygg från den mappen.

**Alternativ 1 – Bygg-script (rekommenderat)**  
Scriptet `go/build/build.ps1` bygger för Windows, Linux och macOS och lägger binärer i `go/build/bin/`:

```powershell
cd go\build
.\build.ps1
```

Standardutdata: `go\build\bin\shbfsmreg.exe` (Windows), samt `shbfsmreg_linux.exe` och `shbfsmreg_darwin.exe` för andra plattformar. Du kan ändra version och sökvägar via parametrar till scriptet.

**Alternativ 2 – Manuell build (endast aktuell plattform)**  
Bygg från `go/src`:

```powershell
cd go\src
go build -o ..\build\bin\shbfsmreg.exe .\cmd
```

Kör det byggda programmet:

```powershell
.\go\build\bin\shbfsmreg.exe -username "..." -password "..." -beerxmlpath "C:\sökväg\till\recept.xml" -brewername "Ditt namn" -breweremail "din@epost.se" -eventid 61 -fveventid 62
```

Se [Parametrar](#parametrar) för alla flaggor.

---

## Användning

### PowerShell

Kör scriptet från `powershell`-mappen (eller ange sökväg till scriptet) med alla obligatoriska parametrar:

```powershell
cd powershell
.\shbf_sm_register.ps1 `
  -username "ditt_shbf_användarnamn" `
  -password "ditt_lösenord" `
  -BeerXmlPath "C:\sökväg\till\recept.xml" `
  -brewersName "Ditt namn" `
  -brewersEmail "din@epost.se" `
  -shbfEventId "61" `
  -shbfFvEventId "62"
```

### Go (efter build)

```powershell
.\go\build\bin\shbfsmreg.exe `
  -username "ditt_shbf_användarnamn" `
  -password "ditt_lösenord" `
  -beerxmlpath "C:\sökväg\till\recept.xml" `
  -brewername "Ditt namn" `
  -breweremail "din@epost.se" `
  -eventid 61 `
  -fveventid 62
```

### Parametrar

Gemensam betydelse; namn skiljer mellan PowerShell och Go.

| Beskrivning | PowerShell | Go (flagga) | Obligatorisk |
|-------------|------------|-------------|--------------|
| SHBF-inloggning | `username` | `-username` | Ja |
| Lösenord | `password` | `-password` | Ja |
| Sökväg till BeerXML | `BeerXmlPath` | `-beerxmlpath` | Ja |
| Bryggarens namn | `brewersName` | `-brewername` | Ja |
| Bryggarens e-post | `brewersEmail` | `-breweremail` | Ja |
| Event-ID för SM | `shbfEventId` | `-eventid` (standard 61) | Ja |
| FV Event-ID | `shbfFvEventId` | `-fveventid` (standard 62) | Ja |
| Domartävlingen | `CompDt` (1/0) | — | Nej (endast PowerShell) |
| Folkets val | `CompFv` (1/0) | — | Nej (endast PowerShell) |
| Max rader malt/humle/övrigt | `IngredientLimit` | `-ingredientlimit` (standard 10) | Nej |

### Event-ID för 2026 SM

- **Event ID:** 61  
- **FV Event ID:** 62  

(Uppdatera dessa om scriptet används för andra år.)

### Efter körning

- Både PowerShell-scriptet och Go-programmet skriver ut HTTP-status och eventuella valideringsfel från servern.
- PowerShell: `body.txt` och `response.txt` sparas i aktuell mapp (för felsökning).
- Vid valideringsfel visas felmeddelanden i rött och körningen avslutas med felkod 1.

### Projektstruktur

- `powershell/` – PowerShell-script `shbf_sm_register.ps1`
- `go/src/` – Go-källkod (cmd, internal/beerxml, internal/shbf, …)
- `go/build/` – Bygg-script `build.ps1` och utdata i `go/build/bin/`
- `examples/` – Exempel på BeerXML-filer
