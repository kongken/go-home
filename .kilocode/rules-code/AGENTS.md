# Project Coding Rules (Non-Obvious Only)

## Backend (Go)

- **Wire DI**: After adding new handlers/services/repositories, run `wire gen cmd/server/wire.go` to regenerate wire_gen.go. Add providers to [`internal/wire/providers.go`](internal/wire/providers.go).

- **Repository Pattern**: Each entity has interface in `internal/repository/xxx_repository.go` and Butterfly implementation in `xxx_repository_butterfly.go`. Use `NewXxxRepositoryButterfly()` constructor.

- **MongoDB Client**: Access via `mongo.GetClient("primary")` - the client name matches config key `store.mongo.primary`.

- **Database Name**: Hardcoded as `gohome` in repository implementations (not configurable).

- **Config Access**: Use `config.Get().JWT` or `config.Get()` for global config. Config loaded from [`config.yaml`](config.yaml).

- **Logger**: Use `log.FromContext(ctx)` from Butterfly framework, not standard log package.

- **Protobuf**: API definitions in `proto/`, generates to `pkg/proto/`. Run `make tproto` after proto changes. Do NOT edit generated code.

## Frontend (React)

- **Package Manager**: Use `yarn` (not npm) - project has yarn.lock.

- **API Client**: Axios instance configured in [`front/src/api/client.ts`](front/src/api/client.ts).

- **Auth State**: Managed by Zustand in [`front/src/stores/authStore.ts`](front/src/stores/authStore.ts).