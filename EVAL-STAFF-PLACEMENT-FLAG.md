# FLAG — Staff rating-column / employment-fields view has no settled home

**Plan:** `20260604-performance-evaluation` · **Raised:** 2026-06-14 · **Status:** BLOCKED on a human decision (do NOT build staff views until resolved).

## What this flag is about

The Performance-Evaluation backend keys every rating to `evaluation.subject_staff_id` (the leased human Identity = `entity/staff`). The UI surface that displays those ratings — the **Associates list Rating column** (`associate-rating-{staff_id}`) plus the staff availability/seniority/employment columns and the staff detail page — needs a **view package home**. There is **no `staff` view anywhere in the codebase today**, so this is a net-new ~12-file view set. The question is *which package it lives in*. This run did NOT build it.

## The conflict (two rules disagree)

1. **The user's literal placement rule (this task):** "subscription-related views (staff, product_price_plan, subscription_seat) → centymo." Read literally, `staff` → **centymo** (`domain/entity/staff/`).

2. **The entydad / 14-layer domain convention (wiki + espyna + lyngua reality):**
   - `staff` is an **`entity`-domain** entity: proto at `proto/v1/domain/entity/staff`, espyna use cases at `usecases/domain/entity/staff`.
   - Its domain peers — `client`, `supplier` — live in **entydad** at `domain/entity/party/`.
   - The staff lyngua labels live in the **entydad tier** (`general/staff.json`, validation/errors only today).
   - **centymo has no `domain/entity/` directory at all.** Placing staff in centymo would create the FIRST `entity`-domain view in centymo, splitting one proto domain across two packages (espyna usecases say `entity`, the view would say centymo) and orphaning the staff view from its client/supplier siblings.

`product_price_plan` is genuinely a `subscription`-domain entity, so its placement in centymo is uncontested (and is what THIS run extended). `staff` is the only contested member of the placement rule.

## Options for the human decision

- **Option A — centymo (`domain/entity/staff/`), literal placement rule.**
  Pro: honors the user's stated rule verbatim; one-time decision.
  Con: first `entity`-domain dir in centymo; cross-package domain split vs the espyna `entity` usecases + entydad client/supplier peers + entydad-tier lyngua; future structure-audit (STR) drift risk.

- **Option B — entydad (`domain/entity/party/staff/`), domain convention (view-delta.md RECOMMENDATION).**
  Pro: matches the proto/espyna `entity` domain, sits beside client/supplier, tracks the entydad-tier `staff.json` lyngua home; clean 14-layer domain map.
  Con: contradicts the literal "subscription → centymo" wording (though staff is an Identity primitive, not a subscription primitive — see entities.md §"the four primitives": staff = **I**dentity, not the **O**ffering).

**The file set is IDENTICAL either way — only the package import path changes.** Whichever package is chosen, the ~12 files are: `staff_module.go`, `staff/{labels,routes,embed,deps}.go`, `staff/list/page.go` (availability Status chip available/assigned/bench/offboarded + Seniority + employment_type columns + **Rating column via `GetLatestEvaluationScore`**, testid `associate-rating-{staff_id}`), `staff/detail/page.go`, `staff/actions.go`, `staff/form/form.go`, `staff/templates/{list,detail,staff-drawer-form}.html`.

**Note:** entydad has no worktree provisioned in this run (only `fayna-eval` + `centymo-eval` exist). If Option B is chosen, an entydad worktree must be created before the staff view can be built.

## What this run DID build (for the record)

- `product_price_plan` rate-band (`billing_amount_min` / `billing_amount_max`) form fields + table display column + labels — EXTEND only, in centymo (uncontested). See the structured return / the four edited files under `domain/subscription/product_price_plan/`.
- This flag note. **No staff views were created.**

## Recommendation

view-delta.md and the 14-layer domain map both recommend **Option B (entydad `domain/entity/party/staff/`)**, treating "subscription → centymo" as covering the Offering (product_price_plan) but not the Identity primitive (staff). But this is the user's call — the placement rule is theirs to interpret. Decide before Wave E2 staff fan-out runs.
