# Project Architecture Rules (Non-Obvious Only)

## Layered Architecture

- **Handler → Service → Repository**: Strict layered pattern. Handlers depend on services, services depend on repositories. No skipping layers.

- **Wire Provider Sets**: Defined in [`internal/wire/providers.go`](internal/wire/providers.go):
  - `RepositorySet` → `CacheSet` → `ServiceSet` → `HandlerSet`
  - Order matters for dependency resolution.

## Dependency Injection

- **Google Wire**: Compile-time DI. After adding new components:
  1. Add provider to appropriate Set in [`internal/wire/providers.go`](internal/wire/providers.go)
  2. Run `wire gen cmd/server/wire.go`
  3. Commit the regenerated [`cmd/server/wire_gen.go`](cmd/server/wire_gen.go)

## Data Layer

- **MongoDB Client**: Named "primary" - must match config key `store.mongo.primary`.
- **Database Name**: Hardcoded `gohome` in repositories - changing requires code modification.
- **Redis**: Two instances configured - "cache" (db 0) and "session" (db 1).

## Frontend Architecture

- **State Management**: Zustand for global state (auth in [`front/src/stores/authStore.ts`](front/src/stores/authStore.ts)).
- **API Layer**: Centralized Axios in [`front/src/api/client.ts`](front/src/api/client.ts).
- **UI Components**: Radix UI primitives in `front/src/components/ui/`.