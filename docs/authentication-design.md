# Authentication & Session Management Implementation

This document describes the end-to-end implementation of authentication, JWT access tokens, and persistent session management for the Org Mem GDG backend.

## 1. Architectural Overview

The system uses a **Dual-Token Strategy**:
- **Short-lived Access Tokens (JWT)**: Used for stateless authorization on protected routes.
- **Long-lived Refresh Tokens (Opaque)**: Used to obtain new access tokens and manage sessions.

This implementation ensures that access can be revoked immediately (via session invalidation) while maintaining the performance benefits of JWTs.

## 2. Token Lifecycle

### 2.1 Login & Initial Issuance
When a user logs in via `/api/login`:
1. **Validation**: Credentials are verified against the `members` or `executives` tables.
2. **Session Creation**: A new session record is created in the `sessions` table (PostgreSQL).
3. **Token Pair Generation**:
   - **Refresh Token**: A 32-byte random hex string. Its hash (SHA-256) is stored in the DB.
   - **Access Token**: A JWT containing the `user_id` (`uid`) and the `refresh_token_id` (`sid`).
4. **Response**: The client receives the JWT and a combined refresh string: `refresh_token_id.refresh_token_secret`.

### 2.2 Token Rotation (Refresh)
When the access token expires, the client calls `/api/refresh`:
1. **Verification**: The server parses the combined refresh string and lookups the session in the DB using `refresh_token_id`.
2. **Security Checks**: It ensures the session is not expired and `revoked_at` is null.
3. **Rotation**:
   - A **new random refresh token secret** is generated.
   - The database record is updated with the new hash (invalidating the previous refresh token).
   - The session expiry is extended.
4. **New Pair**: A new JWT and new combined refresh string are returned.

### 2.3 Logout & Revocation
When `/api/logout` is called:
1. **Mark Revoked**: The `revoked_at` column for that specific `refresh_token_id` is set to the current timestamp.
2. **Immediate Effect**:
   - Any further attempts to use that Refresh Token will fail.
   - The Middleware will reject any Access Tokens linked to that `refresh_token_id`.

## 3. Middleware Enforcement

The `Auth` middleware ([backend/internal/middleware/auth.go](backend/internal/middleware/auth.go)) protects routes by:
1. Validating the JWT signature and expiry.
2. Extracting the `refresh_token_id` (`sid`) from the JWT claims.
3. **Session Verification**: Querying the `sessions` table by `refresh_token_id` to ensure the specific session has not been revoked. This allows for immediate global logout across all devices if a session is compromised.

## 4. Key Security Features

- **SHA-256 Hashing**: Refresh tokens are never stored in plain text.
- **Constant-Time Comparison**: Used when validating refresh token hashes to prevent timing attacks.
- **Session-JWT Linkage**: By including the `session_id` in the JWT, we bridge the gap between stateless tokens and stateful revocation.
- **Refresh Rotation**: Prevents "replay" attacks if a refresh token is stolen, as each token is valid for only a single use.

## 5. File References

- **Logic**: [backend/internal/services/auth_service.go](backend/internal/services/auth_service.go)
- **Middleware**: [backend/internal/middleware/auth.go](backend/internal/middleware/auth.go)
- **Token Management**: [backend/internal/utils/jwt.go](backend/internal/utils/jwt.go) & [backend/internal/utils/session.go](backend/internal/utils/session.go)
- **Schema**: [backend/internal/database/migrations/001_init_sessions_polymorphic.sql](backend/internal/database/migrations/001_init_sessions_polymorphic.sql)
