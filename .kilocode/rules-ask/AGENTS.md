# Project Documentation Rules (Non-Obvious Only)

## Architecture Overview

- **Backend**: Go with Butterfly framework (custom internal framework), Gin HTTP router, MongoDB for storage, Redis for caching.

- **Frontend**: React 19 in `front/` directory - separate project with its own package.json. Uses Vite, TypeScript, Tailwind CSS, Radix UI components.

- **API Definitions**: Protobuf files in `proto/` generate Go code to `pkg/proto/`. Do not edit generated code.

## Key File Locations

- **Config**: [`config.yaml`](config.yaml) at project root - MongoDB URI, Redis, JWT settings.
- **Wire Providers**: [`internal/wire/providers.go`](internal/wire/providers.go) - all DI providers.
- **Main Entry**: [`cmd/server/main.go`](cmd/server/main.go) - app initialization and route setup.

## Naming Conventions

- Repository files: `xxx_repository.go` (interface) + `xxx_repository_butterfly.go` (implementation).
- Handler/Service: Standard `xxx_handler.go`, `xxx_service.go` naming.