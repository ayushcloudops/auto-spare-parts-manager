# PROJECT HANDOFF — Auto Spare Parts Shop Management (Desktop App)

> Living document. Update the **Status** section as steps complete.

## What this project is
A production-grade, fully OFFLINE desktop application for an auto spare parts
shop, intended to be sold to non-technical shopkeepers. Primary target: Windows
(macOS/Linux are a bonus). It is NOT a SaaS. No internet required at runtime.

Core modules (planned): Dashboard, Product/Inventory, Billing, Invoice History,
Customers (with credit), Suppliers, Purchase Entry, Reports (+CSV), ESC/POS
thermal-printer receipts (80mm), India GST.

## Tech stack
- Desktop shell: Wails v2 (v2.12.0)
- Backend: Go 1.26 (module name: `autoshop`)
- Frontend: React 18 + TypeScript + Tailwind (Vite)
- DB: SQLite via GORM
- Migrations: gormigrate (versioned, reversible)
- Architecture: Clean Architecture + Repository pattern + Service layer +
  manual Dependency Injection.

Layered dependency direction (inward):
  frontend (React/TS)
    -> Wails bindings (thin Go structs exposed to JS)
      -> service layer (business logic, transactions)
        -> repository interfaces (in domain) + GORM implementations
          -> SQLite + Printer interface

## LOCKED engineering decisions (do not change without strong reason)
1. SQLite driver = pure-Go `github.com/glebarez/sqlite` (modernc backend),
   NOT cgo `mattn/go-sqlite3`. Reason: Windows cross-compile with no mingw/gcc.
2. Money = `internal/pkg/money.Money` = `int64` PAISE (1 rupee = 100 paise).
   NEVER float. Persists as INTEGER (Valuer/Scanner + GormDataType "integer").
   All rounding happens in ONE place: `Money.Percent()`.
3. Tax model = India GST. Intra-state = CGST + SGST (half each); inter-state =
   IGST. Decided by shop StateCode vs customer. Invoices store DENORMALISED
   totals + per-line SNAPSHOTS so historical bills are immutable.
4. Invoice numbering = atomic, gap-free, per-financial-year counter via
   `InvoiceSequence`, incremented inside the invoice's DB transaction.
5. DB file = per-OS app data dir via `os.UserConfigDir()` +
   `AutoShopManager/shop.db`. Env `AUTOSHOP_DB` overrides (tests/portable).
6. DI = manual constructor injection in one container at startup. No framework.
7. Stock = append-only `StockMovement` ledger (delta + balanceAfter per change).
8. Soft deletes (GORM DeletedAt) so a deleted product still resolves on old bills.

## Conventions
- Money fields are `money.Money`; GST rate is a `float64` percent (0/5/12/18/28).
- Repository INTERFACES in `internal/domain`; GORM IMPLEMENTATIONS in
  `internal/repository` (built on a generic `repository.Base[T]`).
- Transactions propagate via context: `TxManager.Do` puts a *gorm.DB tx in ctx;
  repos pick it up automatically (unit-of-work). Services own tx boundaries.
- Services own business logic; Wails bindings stay thin (validate, map DTO,
  call service, translate errors via `apperr.KindOf`).
- Migrations append-only: never edit a shipped migration, add a new one.
- Every module ships as a full vertical slice (repo->service->binding->UI),
  kept buildable, with tests.

## How to build / test
- Full build: `wails build`  (outputs build/bin/AutoShopManager.app on macOS)
- Dev mode: `wails dev`
- Backend tests: `go test ./internal/...`
- Toolchain: Go 1.26, Node 26, Wails CLI v2.12.0. Pure-Go SQLite = no cgo.

## Status (step-by-step build plan)
- [x] 1. Architecture designed
- [x] 2. Wails project scaffolded — builds & packages
- [x] 3. Database layer — config paths, money type, models, gormigrate, seed; tested
- [x] 4. Backend core — generic repository, tx manager, DI container, bootstrap
- [x] 5. Frontend shell — Tailwind, layout, nav, theme, routing, API wrapper
- [x] 6. FE<->BE smoke test — SystemHandler.Health() round-trip, health dot in topbar
- [x] 7. Product module — repo/service/handler/UI, stock ledger, search/filter/low-stock; tested
- [x] 8. Billing module + India GST engine — GST engine, FY invoice numbering, credit, stock ledger; tested
- [x] 9. Invoice History, Customers, Suppliers, Purchase, Reports, Printing, Dashboard — all shipped
      · Invoice History: search/date filter, view + ESC/POS receipt preview, print/reprint
      · Customers: CRUD + outstanding credit; Suppliers: CRUD
      · Purchase Entry: receive stock (increments inventory + ledger, updates cost)
      · Reports: sales/profit/top-products/low-stock + CSV export
      · Printing: ESC/POS 80mm receipt builder + CUPS `lp` sender + file fallback (Windows raw = TODO)
      · Dashboard: live today's sales/bills/products/low-stock/credit + recent bills
      · Settings: shop profile (drives GST state logic + receipt header) + printer name

ALL 9 STEPS COMPLETE. Backend: `go test ./internal/...` all green; `wails build` green.
Handlers bound: System, Products, Billing, Customers, Suppliers, Purchases, Reports,
Dashboard, Settings, Print.

## Known follow-ups (not blocking)
- Windows raw thermal printing: add a winspool `WritePrinter` implementation of the
  printer.Printer interface (internal/printer). Today Windows falls back to writing the
  raw receipt to <appdata>/receipts; macOS/Linux print via CUPS `lp`. Needs on-device test.
- App icon / installer branding (build/ assets) before distribution.

## Future features (architecture allows; DO NOT implement now)
Cloud Sync, User Login, Multiple Shops, Mobile App, Barcode Scanner, QR/UPI
payments, WhatsApp/Email invoice, Vehicle Compatibility DB, AI Search, OCR
purchase bills, Online Backup. Each maps onto existing layers (e.g. Cloud Sync =
a second repository implementation; Login = binding-layer middleware).
