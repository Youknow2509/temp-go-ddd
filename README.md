# temp-go-ddd

English | [Tiếng Việt](README.vi.md)

## Contact

- **Mail**: *lytranvinh.work@gmail.com*
- **Github**: *https://github.com/Youknow2509*

## Overview

This repository is a Go Domain-Driven Design (DDD) template. It is not a complete product, but a practical starter structure that you can extend for your own service.

## About this template

This template focuses on a clean DDD skeleton for Go backends: clear layer boundaries, centralized bootstrap, and extension points for business code.

### 1) Current project structure

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

### 2) Folder responsibilities

- `cmd/`: service entrypoint. `cmd/server/main.go` handles bootstrap and runtime lifecycle.
- `internal/domain/`: domain contracts and models shared across the system.
    - `domain/config`, `domain/logger`, `domain/cache`: port interfaces.
    - `domain/model`: domain-level data structures such as `SystemConfig`, `LoggerSetting`.
- `pkg/`: adapter implementations for domain ports.
    - `pkg/config/viper.go`: YAML config loading and merging.
    - `pkg/logger/zap.go`: logger implementation using Zap + lumberjack.
- `internal/initialize/`: system bootstrap wiring (config, logger, and future dependencies).
- `internal/global/`: shared runtime singletons/state (logger, config, waitgroup).
- `internal/application/`: application use cases/services (currently a scaffold).
- `internal/infrastructure/`: infrastructure adapters for DB/message broker/external cache (scaffold).
- `internal/interface/`: transport layer (HTTP/gRPC handlers, routing, DTO mapping) (scaffold).
- `internal/shared/`: internal shared utilities.
- `config/`: default and override configuration files.
- `environment/`: compose files and sample env for local/dev.
- `tests/`: unit and integration tests.

### 3) Container and CI files

- `Dockerfile`: multi-stage Docker build with a minimal `scratch` runtime image.
- `.dockerignore`: excludes unnecessary files from Docker build context.
- `.github/workflows/ci.yml`: basic CI pipeline for Go build and Docker build checks.

## How to use this template

1. Clone the repository.

```bash
git clone https://github.com/Youknow2509/temp-go-ddd.git
```

2. Update system name in `Makefile`.

```makefile
SYSTEM_NAME = ABCD
```

3. Prepare config and environment variables.

- Create `config/config.yaml` to override values from `config/config_default.yaml`.
- Create `.env` from `environment/.env.dev` and adjust values.

4. View available Make targets.

```bash
make help
```

5. Start implementing your business code.

## Build and Docker

### Local build

```bash
make build
```

### Docker build (Buildx)

```bash
make docker_build
```

Notes:

- Docker image build uses `docker buildx build` with `--load` for local usage.
- The runtime image is `scratch` (no OS packages) for minimal size.
- TLS CA certificates are copied into runtime image so HTTPS/TLS calls still work.

## CI/CD (basic)

GitHub Actions workflow at `.github/workflows/ci.yml` runs on push and pull request:

1. `go build ./...` (compile check)
2. Docker image build check using Buildx

This is a basic quality gate to ensure the project can be built successfully.

## DDD - Domain-Driven Design

This template follows a layered DDD style with ports/adapters.

### 1) Layers in this template

- **Domain layer (`internal/domain`)**
    - Contains pure business contracts and models.
    - Does not depend on framework-specific code.
    - Example interfaces: `IConfig`, `ILogger`, `IDistributedCache`, `ILocalCache`.

- **Application layer (`internal/application`)**
    - Orchestrates use cases and business workflows.
    - Depends on domain interfaces, not concrete implementations.
    - Currently scaffolded for your project-specific logic.

- **Infrastructure layer (`internal/infrastructure`, `pkg`)**
    - Implements domain ports.
    - Existing examples: `ViperConfig` (config adapter), `ZapLogger` (logger adapter).

- **Interface layer (`internal/interface`)**
    - Exposes system entry points (HTTP, gRPC, event consumers).
    - Maps transport data to application commands/queries.

### 2) Current startup flow

1. `cmd/server/main.go` initializes global `WaitGroup`.
2. Calls `initialize.Initialize()`.
3. `initializeConfig()`:
    - Creates config adapter via `pkg/config.NewViperConfig()`.
    - Loads `config/config_default.yaml`.
    - Loads and deep-merges `config/config.yaml` when present.
    - Stores result in `global.SystemConfig`.
4. `initializeLogger()`:
    - Builds `LoggerConfigPkg` from `global.SystemConfig`.
    - Creates logger with `pkg/logger.NewZapLogger(...)`.
    - Stores logger in `global.Logger`.
5. `main` blocks on `global.WaitGroup.Wait()`.

### 3) Dependency direction

- `interface` -> `application` -> `domain`
- `infrastructure/pkg` implement domain interfaces
- `domain` does not import `application/interface/infrastructure`

This keeps business rules stable and allows replacing adapters (logger/config/cache/DB) with minimal impact on core logic.

### 4) How to extend

- Add bounded contexts under `internal/domain` and matching use cases in `internal/application`.
- Define repository ports in `internal/domain/repository`.
- Implement repositories/adapters in `internal/infrastructure`.
- Expose use cases via `internal/interface`.
- Wire dependencies in `internal/initialize`.

The main goal is a clean domain with technical details (framework, I/O, infra) isolated outside business core.
