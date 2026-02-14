# plugin-morphe-crud-go-gin

Generates Go Gin HTTP handlers and route registration from Morphe model definitions.

## What it generates

Given Morphe model files (`.mod`), this plugin generates:

1. **Gin handlers** — HTTP handler structs that wire Gin request handling to repository method calls
2. **Route registration** — `RegisterRoutes` methods per handler and a top-level `RegisterAllRoutes` function
3. **Response DTO structs** — `.str` files that are compiled to Go structs by `plugin-morphe-go-struct`

> **Note:** Repository interfaces are generated separately by `plugin-morpherepo-go` (via the morpherepo pipeline). This plugin consumes those interfaces — it does not generate them.

## Derivation rules

**From identifiers:**
- `primary: ID` → `GetByID(ctx, id string)`
- `code: Code` → `GetByCode(ctx, code string)`

**From relationships:**
- `ForOne` → filter parameter on `GetAll` (e.g., `ForOne Organization` → `organizationID *string`)
- `ForOnePoly` → filter parameter on `GetAll` (e.g., `ForOnePoly Owner` → `ownerID *string`)
- `HasMany` / `HasOne` — no filter generated (inverse side)

**Standard CRUD:**
- `GetAll` — List with derived filters
- `GetByID` — Fetch by primary identifier
- `Create` — Accept model input, call repo
- `Update` — Accept model input, call repo
- `Delete` — By primary key

## What remains hand-written

1. **Repository implementations** — SQL queries, business logic (or use `plugin-morpherepo-go-psql` to generate PSQL implementations)
2. **Custom endpoints** — Non-CRUD operations (resolve, download, etc.)
3. **Middleware** — Auth, rate limiting, logging
4. **Server wiring** — Combining generated + custom handlers

## Configuration

```yaml
config:
  "@kalo-build/plugin-morphe-crud-go-gin":
    handlers:
      PackagePath: "github.com/your-org/your-app/internal/generated/handlers"
    repo:
      PackagePath: "github.com/your-org/your-app/internal/generated/repo"
    models:
      PackagePath: "github.com/your-org/your-app/internal/types/models"
    structures:
      PackagePath: "github.com/your-org/your-app/internal/types/structures"
    excludeModels:
      - User
```

## Pipeline context

This plugin is typically used alongside:

- `plugin-morphe-morpherepo` — generates morpherepo definitions from Morphe models
- `plugin-morpherepo-go` — generates Go repository interfaces from morpherepo
- `plugin-morpherepo-go-psql` — generates PSQL repository implementations from morpherepo
- `plugin-morphe-go-struct` — compiles the `.str` response DTOs into Go struct files
