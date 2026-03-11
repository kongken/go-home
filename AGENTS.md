# AGENTS.md

This file provides guidance to agents when working with code in this repository.


## Git

当前仓库 github repo: https://github.com/kongken/go-home
你可以调用 mcp 查看 ci 构建情况

## Project Structure

- **Root**: Go backend (Gin + Butterfly framework + MongoDB + Redis)
- **proto/**: Protobuf API definitions (use `make tproto` to generate)
- **front/**: React 19 frontend (Vite + TypeScript + Tailwind + Radix UI)
- **pkg/proto/**: Generated protobuf Go code (auto-generated, do not edit)

## Build Commands

```bash
# Backend
make build          # Build binary to build/go-home
make run            # Run server (go run cmd/server/main.go)
make test           # Run all tests (go test -v ./...)
make lint           # Run golangci-lint
make tproto         # Generate protobuf code (requires buf)
make dev            # Run with hot reload (requires air)

# Frontend (from front/ directory)
yarn dev            # Development server
yarn build          # Production build
```

## Architecture Patterns

- **Dependency Injection**: Uses Google Wire. All providers defined in [`internal/wire/providers.go`](internal/wire/providers.go). Run `wire gen cmd/server/wire.go` after adding new handlers/services/repositories.

- **Repository Pattern**: Each entity has interface + Butterfly implementation (e.g., `UserRepository` interface in [`internal/repository/user_repository.go`](internal/repository/user_repository.go), implementation in `user_repository_butterfly.go`). MongoDB client via `mongo.GetClient("primary")`.

- **Handler → Service → Repository**: Standard layered architecture. Handlers use Gin context, services contain business logic, repositories handle data access.

- **Butterfly Framework**: App lifecycle managed by `butterfly.orx.me/core/app`. Config loaded from [`config.yaml`](config.yaml) with MongoDB URI at `store.mongo.primary.uri`.

## Key Conventions

- Database name: `gohome` (hardcoded in repository implementations)
- Config access: `config.Get().JWT` or `config.Get()` for global config
- Logger: `log.FromContext(ctx)` from Butterfly
- Repository naming: `NewXxxRepositoryButterfly()` constructor pattern
- Protobuf: API definitions in `proto/`, generates to `pkg/proto/`