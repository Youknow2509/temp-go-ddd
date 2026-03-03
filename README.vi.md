# temp-go-ddd

[English](README.md) | Tiếng Việt

## Liên hệ

- **Mail**: *lytranvinh.work@gmail.com*
- **Github**: *https://github.com/Youknow2509*

## Tổng quan

Đây là template DDD (Domain-Driven Design) cho Go. Dự án chưa phải một sản phẩm hoàn chỉnh, mà là bộ khung thực tế để bạn mở rộng theo nghiệp vụ của riêng mình.

## About this template

Template này tập trung vào skeleton DDD cho backend Go: tách lớp rõ ràng, bootstrap tập trung, và để sẵn các điểm mở rộng cho business code.

### 1) Cấu trúc project hiện tại

```text
.
├── .github/
│   └── workflows/
│       └── ci.yml
├── cmd/
│   └── server/
│       └── main.go
├── .dockerignore
├── Dockerfile
├── config/
│   ├── config_default.yaml
│   └── config.yaml
├── docs/
├── environment/
│   ├── .env.dev
│   ├── docker-compose-dev.yml
│   └── docker-compose.yml
├── internal/
│   ├── application/
│   ├── constant/
│   │   ├── config.go
│   │   ├── logger.go
│   │   └── system.go
│   ├── domain/
│   │   ├── cache/
│   │   │   └── cache_domain.go
│   │   ├── config/
│   │   │   └── config_domain.go
│   │   ├── logger/
│   │   │   └── logger_domain.go
│   │   ├── model/
│   │   │   ├── config_model.go
│   │   │   └── logger_model.go
│   │   └── repository/
│   ├── global/
│   │   └── global.go
│   ├── infrastructure/
│   ├── initialize/
│   │   ├── initialize.go
│   │   ├── initialize_config.go
│   │   └── initialize_logger.go
│   ├── interface/
│   └── shared/
├── pkg/
│   ├── config/
│   │   └── viper.go
│   └── logger/
│       └── zap.go
├── proto/
├── tests/
├── Makefile
└── go.mod
```

### 2) Trách nhiệm từng thư mục

- `cmd/`: entrypoint của service. `cmd/server/main.go` xử lý bootstrap và lifecycle runtime.
- `internal/domain/`: chứa contract và model domain dùng xuyên suốt hệ thống.
    - `domain/config`, `domain/logger`, `domain/cache`: định nghĩa interface (port).
    - `domain/model`: cấu trúc dữ liệu domain-level như `SystemConfig`, `LoggerSetting`.
- `pkg/`: adapter implementation cho các interface domain.
    - `pkg/config/viper.go`: load và merge cấu hình YAML.
    - `pkg/logger/zap.go`: logger bằng Zap + lumberjack.
- `internal/initialize/`: tập trung wiring khởi tạo hệ thống (config, logger, dependency khác).
- `internal/global/`: shared runtime singleton/state (logger, config, waitgroup).
- `internal/application/`: use case/service nghiệp vụ (hiện đang là scaffold).
- `internal/infrastructure/`: adapter cho DB/message broker/cache external (scaffold).
- `internal/interface/`: tầng transport (HTTP/gRPC handlers, routing, DTO mapping) (scaffold).
- `internal/shared/`: utility dùng chung nội bộ.
- `config/`: file config mặc định và override.
- `environment/`: compose files và env mẫu cho local/dev.
- `tests/`: unit test và integration test.

### 3) File container và CI

- `Dockerfile`: build image theo multi-stage với runtime `scratch` tối giản.
- `.dockerignore`: loại các file/thư mục không cần thiết khỏi Docker build context.
- `.github/workflows/ci.yml`: pipeline CI cơ bản để check Go build và Docker build.

## Cách sử dụng template

1. Clone repository.

```bash
git clone https://github.com/Youknow2509/temp-go-ddd.git
```

2. Đổi tên hệ thống trong `Makefile`.

```makefile
SYSTEM_NAME = ABCD
```

3. Chuẩn bị config và biến môi trường.

- Tạo `config/config.yaml` để override `config/config_default.yaml`.
- Tạo `.env` từ `environment/.env.dev` rồi cập nhật giá trị phù hợp.

4. Xem các lệnh Make có sẵn.

```bash
make help
```

5. Bắt đầu triển khai business code.

## Build và Docker

### Build local

```bash
make build
```

### Docker build (Buildx)

```bash
make docker_build
```

Ghi chú:

- Build image dùng `docker buildx build` với `--load` để dùng local.
- Runtime image là `scratch` (không OS package) để nhẹ nhất.
- TLS CA certificates được copy vào runtime image để gọi HTTPS/TLS bình thường.

## CI/CD (cơ bản)

Workflow GitHub Actions tại `.github/workflows/ci.yml` chạy khi push và pull request:

1. `go build ./...` (kiểm tra compile)
2. Docker build check bằng Buildx

Đây là quality gate cơ bản để đảm bảo project vẫn build được.

## DDD - Domain-Driven Design

Template triển khai DDD theo mô hình layered + ports/adapters.

### 1) Các layer trong template

- **Domain layer (`internal/domain`)**
    - Chứa contract và model nghiệp vụ thuần.
    - Không phụ thuộc framework cụ thể.
    - Ví dụ interface: `IConfig`, `ILogger`, `IDistributedCache`, `ILocalCache`.

- **Application layer (`internal/application`)**
    - Điều phối use case và workflow nghiệp vụ.
    - Phụ thuộc interface của domain, không phụ thuộc implementation cụ thể.
    - Hiện đang là scaffold để bạn mở rộng.

- **Infrastructure layer (`internal/infrastructure`, `pkg`)**
    - Cung cấp implementation cho các port domain.
    - Ví dụ hiện tại: `ViperConfig` (config adapter), `ZapLogger` (logger adapter).

- **Interface layer (`internal/interface`)**
    - Cung cấp cổng vào hệ thống (HTTP, gRPC, event consumer).
    - Chuyển đổi dữ liệu transport thành command/query cho application.

### 2) Startup flow hiện tại

1. `cmd/server/main.go` khởi tạo global `WaitGroup`.
2. Gọi `initialize.Initialize()`.
3. `initializeConfig()`:

- Tạo config adapter qua `pkg/config.NewViperConfig()`.
- Load `config/config_default.yaml`.
- Nếu có, load và deep-merge `config/config.yaml`.
- Ghi kết quả vào `global.SystemConfig`.

4. `initializeLogger()`:

- Dùng `global.SystemConfig` tạo `LoggerConfigPkg`.
- Tạo logger bằng `pkg/logger.NewZapLogger(...)`.
- Ghi logger vào `global.Logger`.

5. `main` chờ goroutine bằng `global.WaitGroup.Wait()`.

### 3) Hướng phụ thuộc (dependency direction)

- `interface` -> `application` -> `domain`
- `infrastructure/pkg` implement interface của `domain`
- `domain` không import ngược `application/interface/infrastructure`

Nhờ đó bạn có thể thay adapter (logger/config/cache/DB) mà không phải sửa core business.

### 4) Cách mở rộng

- Thêm bounded context trong `internal/domain` và use case tương ứng ở `internal/application`.
- Định nghĩa repository port tại `internal/domain/repository`.
- Implement repository/adapters tại `internal/infrastructure`.
- Expose use case qua `internal/interface`.
- Wire dependency tại `internal/initialize`.

Mục tiêu chính là giữ domain sạch, còn chi tiết kỹ thuật (framework, I/O, infra) nằm ngoài business core.
