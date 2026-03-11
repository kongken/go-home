# Project Debug Rules (Non-Obvious Only)

## Backend Debugging

- **MongoDB Connection**: If "mongo client 'primary' not found" error, check [`config.yaml`](config.yaml) has `store.mongo.primary.uri` correctly configured.

- **Wire Regeneration**: If DI fails after code changes, run `wire gen cmd/server/wire.go` - wire_gen.go must be regenerated.

- **Hot Reload**: Use `make dev` (requires air) for auto-reload during development.

- **Health Check**: Server health at `/health` and `/ready` endpoints - useful for debugging container/deployment issues.

## Frontend Debugging

- **API Base URL**: Configured in [`front/src/api/client.ts`](front/src/api/client.ts) - check if backend URL is correct.

- **Auth Token**: Stored in localStorage via Zustand persist - clear localStorage if auth issues occur.