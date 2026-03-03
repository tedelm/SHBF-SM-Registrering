# SHBF-SM-Registrering

Registrera tävlande öl till SHBF SM via BeerXML-export (t.ex. från BeerSmith). Scriptet loggar in på event.shbf.se och skickar in en ölregistrering med data från din BeerXML-fil.

## Krav

- Windows PowerShell
- Ett BeerXML-recept (export från BeerSmith eller annat verktyg som stöder BeerXML)
- SHBF-konto (användarnamn och lösenord för event.shbf.se)
- Event-ID och FV Event-ID för aktuell SM (se nedan)

## Användning

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

### Parametrar

| Parameter        | Obligatorisk | Beskrivning |
|-----------------|-------------|-------------|
| `username`      | Ja          | SHBF-inloggning (event.shbf.se) |
| `password`      | Ja          | Lösenord för SHBF |
| `BeerXmlPath`   | Ja          | Sökväg till BeerXML-fil (t.ex. export från BeerSmith) |
| `brewersName`   | Ja          | Bryggarens namn (används om inte i BeerXML) |
| `brewersEmail`  | Ja          | Bryggarens e-post |
| `shbfEventId`   | Ja          | Event-ID för SM (t.ex. 61 för 2026) |
| `shbfFvEventId` | Ja          | FV Event-ID för Folkets val (t.ex. 62 för 2026) |
| `CompDt`        | Nej         | `1` = inkludera Domartävlingen (standard), `0` = utelämna |
| `CompFv`        | Nej         | `1` = inkludera Folkets val, `0` = utelämna (standard) |

### Event-ID för 2026 SM

- **Event ID:** 61  
- **FV Event ID:** 62  

(Uppdatera dessa om scriptet används för andra år.)

### Efter körning

- Scriptet skriver ut HTTP-status och eventuella valideringsfel från servern.
- `body.txt` och `response.txt` sparas i aktuell mapp (för felsökning).
- Vid valideringsfel visas felmeddelanden i rött och scriptet avslutas med felkod 1.
