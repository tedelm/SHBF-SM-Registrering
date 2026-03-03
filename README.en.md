# SHBF SM Registration

Register competition beers for SHBF SM using a BeerXML export (e.g. from BeerSmith).  
*[Svenska](README.md)* The script logs in to event.shbf.se and submits a beer registration with data from your BeerXML file.

## Requirements

- Windows PowerShell
- A BeerXML recipe file (export from BeerSmith or any tool that supports BeerXML)
- SHBF account (username and password for event.shbf.se)
- Event ID and FV Event ID for the current SM (see below)

## Limitations

The registration form on event.shbf.se has a limited number of fields. The script therefore imports **at most the first 4 rows** of each type from the BeerXML; any additional rows are ignored.

| Type   | Max count | Notes |
|--------|-----------|--------|
| Malt   | 4         | Only the first 4 fermentables/malts are imported. |
| Hops   | 4         | Only the first 4 hop varieties are imported. |
| Other  | 4         | Only the first 4 misc additions are imported. |

If your recipe has more than 4 malts (or more hops/other), you need to add or edit those entries manually on event.shbf.se after importing, or simplify the recipe in BeerSmith before export.

## Usage

Run the script from the `powershell` folder (or pass the path to the script) with all required parameters:

```powershell
cd powershell
.\shbf_sm_register.ps1 `
  -username "your_shbf_username" `
  -password "your_password" `
  -BeerXmlPath "C:\path\to\recipe.xml" `
  -brewersName "Your name" `
  -brewersEmail "your@email.com" `
  -shbfEventId "61" `
  -shbfFvEventId "62"
```

### Parameters

| Parameter        | Required | Description |
|-----------------|----------|-------------|
| `username`      | Yes      | SHBF login (event.shbf.se) |
| `password`      | Yes      | SHBF password |
| `BeerXmlPath`   | Yes      | Path to BeerXML file (e.g. export from BeerSmith) |
| `brewersName`   | Yes      | Brewer's name (used if not present in BeerXML) |
| `brewersEmail`  | Yes      | Brewer's email |
| `shbfEventId`   | Yes      | Event ID for SM (e.g. 61 for 2026) |
| `shbfFvEventId` | Yes      | FV Event ID for People's choice (e.g. 62 for 2026) |
| `CompDt`        | No       | `1` = include Judge competition (default), `0` = omit |
| `CompFv`        | No       | `1` = include People's choice, `0` = omit (default) |

### Event IDs for 2026 SM

- **Event ID:** 61  
- **FV Event ID:** 62  

(Update these when using the script for other years.)

### After running

- The script prints HTTP status and any validation errors from the server.
- `body.txt` and `response.txt` are saved in the current directory (for debugging).
- On validation errors, messages are shown in red and the script exits with code 1.
