# Product Requirements Document (PRD) - Ping Uptime

## 1. Product Overview
Ping Uptime is an elegant, high-performance, self-hosted service status dashboard and uptime monitor. It tracks the health of registered web endpoints and alerts administrators on downtimes.

## 2. Target Audience & Deployment
- **Audience**: Developers, Sysadmins, and DevOps engineers.
- **Deployment Model**: Self-hosted installation (on-premise or cloud virtual machines) with a single-instance database.

## 3. Core Features & Specifications

### 3.1 Initial Setup & Registrations
- **First-Run Setup**: On a fresh install (database contains 0 user records), the application requires a setup stage. The root page (`/`) automatically renders the **Setup Admin Account** screen instead of the login form.
- **Auto-Admin Assignment**: The first user registered during setup is automatically assigned the `'admin'` role.
- **Lockdown Security**: Once setup is completed, the registration endpoints (`/auth/register` or `/auth/setup`) are locked down. Subsequent attempts to call setup or register will return a `430 StatusForbidden` ("Setup already completed"). No public signup page exists.
- **Thread-Safety**: Concurrency checks are enforced using a mutex lock in the backend during registration to prevent double-admin setup race conditions.

### 3.2 Authentication & Session Flow (JWT Access/Refresh Tokens)
- **Token Dual-Storage**: Upon successful login or setup, the system generates:
  - **Access Token**: Short-lived JWT (15-minute expiration) used for request authorization.
  - **Refresh Token**: Long-lived JWT (7-day expiration) used to request new access tokens.
- **Silent Refresh Interceptor**: The frontend axios wrapper interceptor catches `401 Unauthorized` responses on expired access tokens and silently requests a new one from `/api/auth/refresh` using the refresh token, retrying the original request automatically without disrupting the user.
- **Server-Side Token Verification**: Session validation does not rely on local cookies. On page load, the frontend makes an API call to `/api/users/verify` to confirm that the server still recognizes the active token and user account in the database.

### 3.3 Dashboard Layout
- **App View (`/app`)**: Displays real-time operational status (latency, uptime percentage, status indicator) of registered services.
- **Responsive Layout**: Adapts gracefully to all screen resolutions (mobile, tablet, desktop) and includes a collapsible sidebar.
- **Theme Support**: Provides toggleable Light and Dark mode options. The custom OKLCH emerald dark theme offers soft, green-tinted, premium backgrounds.

## 4. Technical Architecture

### 4.1 Frontend Stack
- **Framework**: Vue 3 (Vite compile setup)
- **State Management**: Pinia Store (`@/stores/auth.ts`) wrapped in a reactive Composable (`@/composables/useAuth.ts`)
- **Validations**: Zod schemas (`@/validations/auth.ts`) integrated with VeeValidate (`<Form>` / `<FormField>`)
- **Styling**: Tailwind CSS v4 + shadcn-vue component layouts
- **Animations**: GSAP (Timeline card shake and entry transitions)

### 4.2 Backend Stack
- **Language**: Go 1.22+
- **HTTP Server**: Echo Web Framework
- **ORM & DB**: GORM connecting to MySQL
- **Authentication**: JWT token validator and parsing middleware

## 5. Security & Authorization Matrix
- **`/auth/setup` (POST)**: Public only when `users` table is empty.
- **`/users/verify` (GET)**: Guarded by `middleware.Auth`. Validates the access token in database.
- **`/users/*` (GET/POST/PUT/DELETE)**: Guarded by `middleware.Auth` and `middleware.Admin`. Enforces role validation so standard users cannot manage accounts.
