# boussinesq-z — Go 弹性半空间附加应力核算 HTTP 服务

本弹性半空间附加应力核算 HTTP 服务：给定荷载与坐标，计算点荷载或圆形荷载应力并可叠加；非法几何或参数须报错。

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
