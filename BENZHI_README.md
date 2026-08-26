# boussinesq-z — Go 弹性半空间附加应力 HTTP 后端与命令行工具

The server starts on `:8080` by default and exposes `POST /api/stress` and
`POST /api/superpose`; CLI subcommands `stress`, `superpose`, and `circular`
provide the same calculations in a terminal.

## Build / Run / Test

```text
go build ./...
go run . serve -addr :8080
go run . example -file example/point-100kn.json
go test ./...
```

## Evaluation Image

Evaluation-specific files (do not overwrite project Dockerfile/README):

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md` (this file)

Build and verify in container:

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
