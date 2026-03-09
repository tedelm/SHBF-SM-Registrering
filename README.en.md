# SHBF SM Registration

Register competition beers for SHBF SM using a BeerXML export (e.g. from BeerSmith).  
*[Svenska](README.md)* You can run either the pre-built binaries listed on the releases page, **PowerShell script** (requires only Windows PowerShell) or the **Go program** (requires Go to be installed). Both log in to event.shbf.se and submit a beer registration with data from your BeerXML file.

**Running the pre-compiled binaries (See releases for latest .exe, no need to build):**  
- **Windows:** `.\shbfsmreg.exe -username pelle -password mittlösenord -beerxmlpath .\mitt_recept.xml -brewername pelle -breweremail pelle@epost.se` 
- **Linux:** `./shbfsmreg_linux -username pelle -password mittlösenord -beerxmlpath .\mitt_recept.xml -brewername pelle -breweremail pelle@epost.se` 
- **macOS:** `./shbfsmreg_darwin -username pelle -password mittlösenord -beerxmlpath .\mitt_recept.xml -brewername pelle -breweremail pelle@epost.se`

See [Parameters](#parameters) for all flags.

## Requirements (There is no requirements if you choose to run the pre-compiled binaries)

- **PowerShell:** Windows PowerShell (no extra install needed for the script)
- **Go program:** [Go (Golang)](https://go.dev/) 1.21+ (see [Installing Go](#installing-go) below)
- A BeerXML recipe file (export from BeerSmith or any tool that supports BeerXML)
- SHBF account (username and password for event.shbf.se)
- Event ID and FV Event ID for the current SM (see below)

## Limitations

The registration form on event.shbf.se has a limited number of fields. The script imports **at most a configurable number of rows** (default 10) of each type from the BeerXML; any additional rows are ignored.

| Type   | Default max | Notes |
|--------|-------------|--------|
| Malt   | 10          | Control with `IngredientLimit` (PowerShell) or `-ingredientlimit` (Go). |
| Hops   | 10          | Same as above. |
| Other  | 10          | Same as above. |

For more ingredients than the limit, increase `IngredientLimit` / `-ingredientlimit`, or edit the registration manually on event.shbf.se after importing.

## Validation (Go program)

The Go program checks that the recipe fits the chosen style using the **STYLE** element in the BeerXML. If the style defines min/max for a dimension, the recipe’s values are compared; if any value is outside the range, the run stops with an error.

| Check   | Recipe value (from BeerXML) | Style bounds (STYLE)   |
|--------|-----------------------------|------------------------|
| OG     | OG / EST_OG                 | OG_MIN, OG_MAX         |
| FG     | FG / EST_FG                 | FG_MIN, FG_MAX         |
| IBU    | IBU / EST_IBU               | IBU_MIN, IBU_MAX       |
| EBC (color) | EST_COLOR               | COLOR_MIN, COLOR_MAX   |
| ABV    | ABV / EST_ABV               | ABV_MIN, ABV_MAX       |

If min/max are missing for a dimension in the style, that check is skipped. The recipe must have a valid style (CATEGORY_NUMBER and STYLE_LETTER).

---

## Installing Go (for the Go program)

To build and run the Go version you need Go 1.21 or later.

- **Windows:** Download the [installer from go.dev](https://go.dev/dl/) (e.g. `go1.21.x.windows-amd64.msi`) and run it. Alternatively: [Chocolatey](https://chocolatey.org/) – `choco install golang`.
- **Linux:** e.g. `sudo apt install golang-go` (Debian/Ubuntu) or `sudo dnf install golang` (Fedora). Check version with `go version`.
- **macOS:** `brew install go`.

Verify the installation:

```powershell
go version
```

---

## Building the Go program

Source code lives under `go/src`. Build from that directory.

**Option 1 – Build script (recommended)**  
The script `go/build/build.ps1` builds for Windows, Linux, and macOS and places binaries in `go/build/bin/`:

```powershell
cd go\build
.\build.ps1
```

Default output: `go/build/bin/shbfsmreg.exe` (Windows), plus `shbfsmreg_linux.exe` and `shbfsmreg_darwin.exe` for other platforms. You can override version and paths via script parameters.

**Option 2 – Manual build (current platform only)**  
From `go/src`:

```powershell
cd go\src
go build -o ..\build\bin\shbfsmreg.exe .\cmd
```

Run the built binary:

```powershell
.\go\build\bin\shbfsmreg.exe -username "..." -password "..." -beerxmlpath "C:\path\to\recipe.xml" -brewername "Your name" -breweremail "your@email.com" -eventid 61 -fveventid 62
```

See [Parameters](#parameters) for all flags.

---

## Usage

### PowerShell

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

### Go (after building)

```powershell
.\go\build\bin\shbfsmreg.exe `
  -username "your_shbf_username" `
  -password "your_password" `
  -beerxmlpath "C:\path\to\recipe.xml" `
  -brewername "Your name" `
  -breweremail "your@email.com" `
  -eventid 61 `
  -fveventid 62
```

### Parameters

Same meaning across tools; names differ between PowerShell and Go.

| Description | PowerShell | Go (flag) | Required |
|-------------|------------|-----------|----------|
| SHBF login | `username` | `-username` | Yes |
| Password | `password` | `-password` | Yes |
| Path to BeerXML | `BeerXmlPath` | `-beerxmlpath` | Yes |
| Brewer's name | `brewersName` | `-brewername` | Yes |
| Brewer's email | `brewersEmail` | `-breweremail` | Yes |
| Event ID for SM | `shbfEventId` | `-eventid` (default 61) | Yes |
| FV Event ID | `shbfFvEventId` | `-fveventid` (default 62) | Yes |
| Judge competition | `CompDt` (1/0) | — | No (PowerShell only) |
| People's choice | `CompFv` (1/0) | — | No (PowerShell only) |
| Max rows malt/hops/other | `IngredientLimit` | `-ingredientlimit` (default 10) | No |

### Event IDs for 2026 SM

- **Event ID:** 61  
- **FV Event ID:** 62  

(Update these when using the script for other years.)

### After running

- Both the PowerShell script and the Go program print HTTP status and any validation errors from the server.
- PowerShell: `body.txt` and `response.txt` are saved in the current directory (for debugging).
- On validation errors, messages are shown in red and the run exits with code 1.

### Project structure

- `powershell/` – PowerShell script `shbf_sm_register.ps1`
- `go/src/` – Go source (cmd, internal/beerxml, internal/shbf, …)
- `go/build/` – Build script `build.ps1` and output in `go/build/bin/`
- `examples/` – Sample BeerXML files
