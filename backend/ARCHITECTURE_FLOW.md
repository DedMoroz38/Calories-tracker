# Backend Structure Guide

This project uses a practical MVC-style layering on top of `Fiber` and `GORM`.
The important thing is not just the folders, but the wiring order:

```text
main -> server -> route registration -> middleware -> controller -> service -> model -> schema/database
```

The app does not use a DI container. Dependencies are passed manually, mostly through `Fiber` request context via `c.Locals(...)`.

## 1. Startup Chain

### `cmd/app/main.go`

Process entrypoint:

1. Prints `"Starting."`
2. Sets `security.LS.StorageSecret`
3. Calls `server.Serve()`

This file does not build controllers or services directly. Its job is only process startup.

### `internal/server/server.go`

`server.Serve()` does this:

1. Calls `NewApp()`
2. `NewApp()` creates the `Fiber` app
3. Registers global middleware such as `cors`
4. Calls `routes.RegisterApi(app)`
5. Returns the configured app
6. `Serve()` runs `app.Listen(":2000")`

So the real application composition starts inside `routes.RegisterApi(app)`.

## 2. The Composition Root

### `internal/MVC/routes/api_routes.go`

This is the main wiring function for the application.

It creates shared dependencies once:

```go
ykClient := external_api.NewYKClient()
dbService := db.NewDBService()
```

Then it installs them into the app:

1. A middleware stores `ykClient` in `c.Locals("ykClient")`
2. `/api` and `/api/v1` groups are created
3. `db.InitAdmin(dbService.DB)` ensures the admin user exists
4. `v1.Use(middleware.PGMiddleware(dbService))` makes DB access available to every `/api/v1/...` request
5. Route groups are registered:
   - `/auth`
   - `/products`
   - `/product_types`
   - `/product_platforms`
   - `/product_licenses`
   - `/bundles`
   - `/basket`
   - `/search`
   - `/pages`
   - `/orders`

This file is effectively the app's composition root.

### Route registration pattern

Each domain is registered the same way:

1. Create a sub-router from `v1.Group(...)`
2. Pass that router into a domain-specific `Register...Routes(...)` function

Example:

```go
authRoutes := v1.Group("/auth")
RegisterAuthRoutes(authRoutes)
```

That means route files do not create the app or the version group themselves. They only receive a scoped router and attach endpoints to it.

## 3. How the Database Handle Is Created

### `internal/db/postgres_service.go`

`db.NewDBService()`:

1. Builds the PostgreSQL DSN
2. Opens a single `*gorm.DB` connection with `gorm.Open(...)`
3. Runs `AutoMigrate(...)` for all schema structs
4. Returns:

```go
type DBService struct {
    DB *gorm.DB
}
```

That means the shared DB handle starts here:

```text
db.NewDBService() -> &DBService{DB: db}
```

## 4. How `*gorm.DB` Reaches a Service

This is the most important dependency flow in the repo.

### Step 1: DB is created once at app startup

In `RegisterApi`:

```go
dbService := db.NewDBService()
```

### Step 2: DB is attached to each request

`v1.Use(middleware.PGMiddleware(dbService))`

Inside `internal/middleware/db_middleware.go`:

```go
c.Locals("gorm", dbService.DB)
```

So every request under `/api/v1` gets the same shared `*gorm.DB` handle in request-local context.

### Step 3: Controller reads the DB from request context

Example from `AuthController.Register`:

```go
as := services.AuthService{}
as.DB = c.Locals("gorm").(*gorm.DB)
```

### Step 4: Service passes the same DB into a model

Example from `AuthService.Register`:

```go
am := models.AuthModel{}
am.DB = as.DB
```

### Step 5: Model executes queries

Example from `AuthModel.Register`:

```go
tx := am.DB.Create(&user)
```

So the actual chain is:

```text
db.NewDBService()
  -> PGMiddleware(dbService)
  -> c.Locals("gorm")
  -> controller.DB assignment
  -> service.DB assignment
  -> model.DB assignment
  -> GORM query
```

## 5. Why `AuthService` Can Access `.DB`

`AuthService` is declared as:

```go
type AuthService struct {
    *gorm.DB
}
```

This is an embedded field. Because of embedding, `as.DB` is valid.

Other services use the more explicit form:

```go
type ProductService struct {
    DB *gorm.DB
}
```

Both patterns work, and in both cases controllers assign the DB the same way:

```go
service.DB = c.Locals("gorm").(*gorm.DB)
```

So in this project:

- `AuthService` and `BasketService` embed `*gorm.DB`
- `ProductService` and `OrderService` declare `DB *gorm.DB`
- the controller usage is still effectively the same

## 6. Request Lifecycle: Auth Example

Example endpoint:

```text
POST /api/v1/auth/register
```

### Route layer

In `internal/MVC/routes/auth_routes.go`:

```go
user.Post("/register", middleware.BodyValidation(&dto.AuthRegisterRequest{}), ac.Register)
```

That means the request passes through:

1. `BodyValidation(...)`
2. `AuthController.Register`

### Validation middleware

`BodyValidation(...)`:

1. Parses JSON into the DTO struct
2. Calls `req.Validate()`
3. Stores the parsed DTO in:

```go
c.Locals("validatedBody", req)
```

### Controller layer

`AuthController.Register`:

1. Reads the validated DTO from `c.Locals("validatedBody")`
2. Creates `AuthService{}`
3. Reads `*gorm.DB` from `c.Locals("gorm")`
4. Calls `as.Register(...)`
5. Converts service result into HTTP response

The controller is HTTP-focused. It does not query the database directly.

### Service layer

`AuthService.Register`:

1. Creates `AuthModel{}`
2. Sets `am.DB = as.DB`
3. Runs business rules:
   - check whether email exists
   - decide whether to reject, update, or create user
   - create/send verification code
4. Calls helpers such as:
   - `authorization`
   - `security`
   - `email_builder`
   - external email API

The service is the use-case layer. It coordinates multiple components.

### Model layer

`AuthModel` does the actual persistence work:

- `Create`
- `Where(...).First(...)`
- `Save`
- `Delete`
- `Updates`

The model knows how to talk to the database. It does not know anything about HTTP.

### Schema layer

The model reads and writes DB structs from `internal/db`, for example:

- `db.AuthUser`
- `db.AuthVerifyCode`
- `db.AuthResetToken`

That is the final layer where GORM maps Go structs to SQL tables.

## 7. Request Lifecycle: Protected Route Example

Example endpoint:

```text
GET /api/v1/auth/me
```

Route:

```go
user.Get("/me", middleware.CheckJWTWithRole("user"), ac.GetUserInfo)
```

Flow:

1. `CheckJWTWithRole("user")` reads cookie `jwt_token`
2. It parses the JWT
3. It stores claims in:

```go
c.Locals("jwt_claims", claims)
```

4. It also uses `c.Locals("gorm")` to create `AuthModel` and verify password freshness
5. If the role is allowed, request reaches `AuthController.GetUserInfo`
6. The controller calls `middleware.GetUserIDFromJWTClaims(c)`
7. The controller creates `AuthService`, injects DB, and calls `GetEmail(userID)`

So middleware is not only auth checking here; it also prepares data that the controller will reuse.

## 8. How External Clients Are Passed

The database is not the only dependency passed through request context.

### Payment client example

In `RegisterApi`:

```go
ykClient := external_api.NewYKClient()
```

Then middleware stores it:

```go
c.Locals("ykClient", ykClient)
```

Later, `OrderController.CreateOrder` does:

```go
os.YKClient = c.Locals("ykClient").(*api.YKClient)
```

So the same general pattern is used for non-DB dependencies too:

```text
create shared dependency once -> put into c.Locals(...) -> read in controller -> pass into service
```

## 9. Responsibility of Each Folder

### `cmd/app`

Application entrypoint.

### `internal/server`

Creates and runs the `Fiber` app.

### `internal/MVC/routes`

Declares endpoints and attaches route-level middleware.

### `internal/middleware`

Adds request-scoped data:

- parsed body
- DB handle
- JWT claims

### `internal/MVC/controllers`

HTTP layer:

- reads params/body/cookies/context
- creates service instances
- maps service result to HTTP response

Most controllers also embed `BaseController`, which currently provides a shared fallback method:

```go
type AuthController struct {
    BaseController
}
```

That inheritance-like embedding is used for shared controller behavior, not for dependency injection.

### `internal/MVC/services`

Business/use-case layer:

- combines models
- runs application logic
- calls utilities and external APIs

### `internal/MVC/models`

Persistence layer:

- contains GORM queries
- reads/writes schema structs

The product domain is split into subpackages that mirror this same layering:

- `controllers/product_controllers`
- `services/product_services`
- `models/product_models`
- `db/product_schema`

### `internal/db`

Database bootstrap and schema definitions.

### `internal/common/dto`

Request/response objects and validation methods.

### `internal/util`

Cross-cutting helpers such as JWT, password hashing, image handling, converters, validators, email building.

### `internal/api`

External integrations such as payment and email.

## 10. Reusable Mental Model for Another Project

If you copy this structure into another project, keep this rule:

```text
Routes know which controller method handles a URL.
Controllers know how to speak HTTP and create services.
Services know business logic and orchestrate models.
Models know database queries.
Schemas define table shapes.
Middleware injects per-request dependencies and validated request data.
```

Minimal template:

```go
// route
userRoutes.Post("/", middleware.BodyValidation(&dto.UserCreateRequest{}), uc.Create)

// controller
func (uc UserController) Create(c *fiber.Ctx) error {
    req := c.Locals("validatedBody").(*dto.UserCreateRequest)

    us := services.UserService{}
    us.DB = c.Locals("gorm").(*gorm.DB)

    result, apiErr := us.Create(*req)
    if apiErr != nil {
        return c.Status(apiErr.StatusCode).JSON(dto.BaseResponse{Message: apiErr})
    }

    return c.Status(fiber.StatusCreated).JSON(dto.BaseResponse{Data: result})
}

// service
type UserService struct {
    DB *gorm.DB
}

func (us UserService) Create(req dto.UserCreateRequest) (*db.User, *errors.APIError) {
    um := models.UserModel{}
    um.DB = us.DB
    return um.Create(req)
}

// model
type UserModel struct {
    DB *gorm.DB
}
```

## 11. Important Note If You Reuse This Pattern

This repo uses request context as a lightweight manual DI system.

That is fine for a small-to-medium service, but keep these rules:

1. Put shared infrastructure in one place at startup
2. Inject request-scoped access through middleware
3. Keep controllers thin
4. Keep raw SQL and GORM calls inside models
5. If you start a transaction with `tx := db.Begin()`, pass `tx` down to models instead of continuing to use the original `DB`

That last point matters because several services in this repo start transactions, but the transaction handle is not consistently propagated into every model call. If you copy the architecture, copy the layering, not that mistake.

## 12. Short Summary

The real structure of this codebase is:

```text
main
  -> server.NewApp / server.Serve
  -> routes.RegisterApi
  -> startup dependencies created once (DB, payment client)
  -> middleware injects dependencies into request context
  -> route picks controller method
  -> controller pulls validated input + dependencies from context
  -> controller constructs service
  -> service constructs model(s)
  -> model runs GORM on schema structs
  -> controller returns HTTP response
```

If you preserve that flow, another project will follow the same structure cleanly.
