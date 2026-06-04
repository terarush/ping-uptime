# Development Guidelines

## Build Commands
- **Backend (Go)**: `go build -o tmp/main main.go`
- **Frontend (Vue)**: `npm run build --prefix web`
- **Full Clean Build**: `npm run build --prefix web && go build -o tmp/main main.go`

## Running Locally
- **Backend Server**: `go run main.go`
- **Frontend Dev Server**: `npm run dev --prefix web` (runs on `http://localhost:5173`)

## Testing and Type Checking
- **Backend Tests**: `go test ./...`
- **Frontend Unit Tests**: `npm run test:unit --prefix web`
- **Frontend Type Check**: `npm run type-check --prefix web`

## Code Guidelines and Standards

### Go Backend
- **Architecture**: Modular layout (`/modules/<name>`) with interfaces for handlers, services, and repositories.
- **Routing**: Group endpoints, apply authorization middleware at the group level where possible.
- **Middleware**: Use `middleware.Auth` for standard session validation and `middleware.Admin` for routes requiring administrative privileges.
- **Response Format**: Use the standard JSON response helpers (`h.r.SuccessResponse`, `h.r.ErrorResponse`).

### Vue Frontend
- **Design system**: Tailored using Tailwind CSS v4 and shadcn-vue components. Avoid direct styling when predefined tokens exist.
- **Views**: Flattened structure directly under `@/views` (e.g. `/` is `Index.vue`, `/app` is `App.vue`). Do not nest in subfolders.
- **State Management**: Access and manipulate authentication or global parameters inside Pinia stores (`@/stores`) wrapped in reactive composables (`@/composables`).
- **Form Validation**: Build forms using vee-validate `<Form>` / `<FormField>` structure integrated with Zod validation schemas (`@/validations`).
- **Token Management**: Bearer token auth via `accessToken` (15m expiry) and `refreshToken` (7d expiry) cookies. Direct endpoint validation is required.
