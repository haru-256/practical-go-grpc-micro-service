# Copilot Instructions

## Architecture Snapshot

- The active code lives in `workbench/`: CQRS split into `service/command` (write), `service/query` (read), and `service/client` (REST facade) fed by protobufs in `api/` and backed by dual MySQL instances under `db/`.
- Command uses SQLBoiler repositories (`internal/infrastructure/sqlboiler`), Query uses GORM (`internal/infrastructure/db`), Client calls both via Connect RPC plus Echo HTTP handlers (`internal/presentation/server`).
- Shared building blocks live in `workbench/pkg`: `connect/interceptor` (logging + protovalidate) and `log` (slog + Otel handler); reuse these instead of inventing ad-hoc middleware.
- `official/` mirrors the book reference implementation; treat it as read-only context, new work belongs under `workbench/`.

## Build & Run Workflow

- Run `mise install` once to provision tool versions (Go 1.25.1, buf, sqlboiler, ginkgo, swag, etc.).
- From `workbench/`: `make init` (tidy + octocov), `make lint`, `make fmt`, `make test`, `make test-all` (`-tags=integration`), and `make up`/`make down` for the full docker-compose stack in `compose.yaml`.
- Service Makefiles provide focused commands: `service/command` exposes `make ginkgo`, `make generate-db-models`, `make run-server`; `service/query` and `service/client` offer `make test[-all]`, `make run-server`, and buf-curl helpers; Client also has `make swag-init` for Swagger regeneration.
- Protobuf edits require `cd workbench/api && make generate` (buf lint/format/generate). Update Go modules afterwards with `go mod tidy` at the repo root or per service.

## Database Expectations

- MySQL lives in `workbench/db` with GTID replication from command→query. Bring databases up via `docker compose up -d command_db query_db db_admin` and seed with `make -C workbench/db create-data dump restore start-replication` when tests require fresh data.
- Integration tests tagged `integration` expect both DBs plus replication; CI-focused runs use `-tags=ci` to skip DB interactions.

## Coding Conventions

- Constructors follow "accept interfaces, return concrete"; register them with Uber Fx modules (`internal/*/module.go`) using `fx.Annotate(..., fx.As(new(interface)))`.
- Config is provided by Viper-backed structs (`internal/infrastructure/config`) reading each service's `config.toml`, overridable via env vars (`.` → `_`).
- Error flow is standardized: domain/application/CRUD/internal errors live in `workbench/pkg/errs`; presentation layers translate them to Connect error codes.
- Logging goes through injected `*slog.Logger` instances (usually from `pkg/log`). Use context-aware methods and structured attrs; Connect services should wrap handlers with `pkg/connect/interceptor` logger + validator interceptors.

## Testing & Tooling Nuances

- Command service unit specs use Ginkgo/Gomega (`go tool ginkgo run --race ./...`); other packages rely on `go test` + `testify`. Mocks are generated via `go generate ./...` using `gomock` definitions under `internal/mock`.
- Client REST handlers expect DTO validation tags enforced by Echo's `CustomValidator`; keep Swagger annotations in `cmd/server/main.go` synchronized with DTO fields before running `make swag-init`.
- When calling gRPC endpoints manually, prefer the buf-curl snippets baked into each service Makefile so payloads match the proto schema (they default to HTTP/2 prior knowledge on localhost ports 8083/8085).

## Collaboration Tips

- Keep shared protobuf models (`api/proto/common/v1`) as the single source of truth; regenerate both Go stubs and documentation (`api/gen`) before touching service-facing code.
- Any change that crosses layers (DTO ↔ domain ↔ persistence) needs updates in mapping helpers plus associated tests—search for existing converters in `internal/presentation/server/handler.go` before adding new ones.
- If you introduce new interceptors, register them centrally in `pkg/connect/interceptor` and wire them through the Fx server modules so all services pick them up uniformly.

## Utility

- Use `rg` instead of `grep` to search for symbols across the monorepo efficiently.
- Use `sd` instead of `sed` for in-place multi-file text replacements with regex support.
- Use `bat` instead of `cat` for syntax-highlighted file viewing in the terminal.
- Use `fd` instead of `find` for faster and more intuitive file searching.
