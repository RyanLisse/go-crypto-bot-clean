# Audit & Fix Plan

**Date:** 2025-04-17

## 1. Back‑end Dependency Fixes

1. Add missing decimal module:
   ```bash
   go get github.com/shopspring/decimal@latest
   go mod tidy
   ```
2. Re‑build the project:
   ```bash
   go build ./...
   ```

## 2. Break the Import Cycle

1. Reproduce cycle with `go build`, note import trace.
2. Enforce layering:
   - **domain/port**: all interfaces
   - **usecase**: business logic, depends on ports & models
   - **service**: implements ports, depends only on ports & models
   - **adapter**: framework wiring, sits at the edge
3. Refactor to remove service ↔ usecase mutual imports:
   - Move shared types/helpers into a new `domain/common` or extend ports
   - Ensure `service` imports only `domain/port` & `domain/model`
   - Ensure `usecase` imports only `domain/port`, `domain/model`, and needed service APIs
4. Re‑build to confirm the cycle is resolved.

## 3. Migrate All Status Types to `model.CoinStatus`

1. Search for `model.Status` across repos, services, use cases, mocks, tests.
2. Update signatures, GORM queries, casts to `model.CoinStatus`.
3. Re‑build and run tests after each batch.

## 4. Update & Align Mocks/Tests

1. Adjust mock methods to accept and return `model.CoinStatus`.
2. Fix all tests referencing the old `Status` type or mismatched signatures.
3. Run:
   ```bash
   go test ./...
   ```
   until all pass.

## 5. Front‑end Build & Type‑check

1. From `frontend_next/`, run:
   ```bash
   bun run build
   # or next build
   ```
2. Fix any TS/CSS errors (e.g. import paths in `globals.css`).
3. Verify `layout.tsx` changes (ClerkProvider nesting, fonts).
4. Run front‑end tests (Jest, etc.) if present.

## 6. Final End‑to‑End Smoke Test

1. Start back‑end:
   ```bash
   ./run_server_sqlite.sh
   ```
2. Start front‑end dev server, navigate through core flows:
   - Sign‑in
   - Dashboard charts & tables
   - New coin listings
3. Monitor console, network, and server logs for hidden errors.

---
Follow these steps in sequence, rerunning builds/tests after each major change, to isolate and resolve every compiler error and test failure.
