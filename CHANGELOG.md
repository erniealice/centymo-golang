# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0-alpha] - 2026-06-15

Commerce domain (largest domain package) — first published alpha.

### Added
- Products (variants, options, images), product lines, inventory (serials, depreciation, movements), price lists/plans, subscriptions, revenue/sales (line items, payments, invoices, email), expenditure (purchase orders, expenses, categories), collections (money in), disbursements (money out), and checkout with serial allocation.

### Changed
- `go.mod` now references published tags (`v0.1.0-alpha`) instead of local `replace` directives; local development continues via `go.work`.

[Unreleased]: https://github.com/erniealice/centymo-golang/compare/v0.1.0-alpha...HEAD
[0.1.0-alpha]: https://github.com/erniealice/centymo-golang/releases/tag/v0.1.0-alpha
