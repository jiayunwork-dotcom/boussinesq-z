# boussinesq-z

boussinesq-z is a Go calculator for vertical stress in an elastic half-space
under a point load or a uniformly loaded circular area. For a point load it
computes `sigma_z = 3*P*z^3/(2*pi*R^5)` with `R = sqrt(r^2+z^2)`, the influence
coefficient, and companion radial/tangential/shear components. Multiple point
loads can be superposed algebraically. Uniform circular loads use a fixed CSV
influence-coefficient table with bilinear interpolation; the table file is
re-read so edits are reflected in later results. The service is available
through HTTP JSON endpoints and CLI subcommands with no web page.

## Usage

Run the HTTP server:

```bash
go run . serve -addr :8080
```

Evaluate from the command line:

```bash
go run . stress -P 100000 -z 2 -r 0
go run . superpose -file example/superpose.json
go run . circular -q 100000 -a 1 -z 2 -r 0
```

Run the point-load scenario:

```bash
go run . example -file example/point-100kn.json
```

The example applies 100 kN at `z = 2 m`, `r = 0`; the resulting vertical
stress is about 11.94 kPa.

## HTTP API

```text
POST /api/stress     {"P":100000,"z":2,"r":0,"poisson":0.3}
POST /api/superpose  {"forces":[{"P":100000,"z":2,"r":0},{"P":50000,"z":3,"r":1}]}
POST /api/circular   {"q":100000,"a":1,"z":2,"r":0}
GET  /api/table-snapshot
GET  /health
```

Invalid inputs return an error body with HTTP 400. This includes negative
loads under the compression-positive convention, `z < 0`, and the `R = 0`
singularity with a positive load at the surface origin. A zero load is valid
and returns zero stress.

## Influence Table

`data/circular-influence.csv` stores rows for `r/a` and columns for `z/a`.
Regenerate it from the polar integral:

```bash
go run ./cmd/gentable -out data/circular-influence.csv -radial 256 -angle 96
```

## Code Layout

```text
internal/point     point-load formula, influence, horizontal components
internal/super     algebraic superposition of point loads
internal/circular  circular-load table, CSV IO, interpolation
internal/server    HTTP handlers and JSON responses
internal/cli       subcommand parsing and terminal output
example/           offline scenario JSON files
data/              influence coefficient CSV
cmd/gentable       table generator
```

## Build and Test

```bash
export GOTOOLCHAIN=local CGO_ENABLED=0
go build ./...
go test ./...
```

The Dockerfile builds the server binary and starts it on port 8080.
