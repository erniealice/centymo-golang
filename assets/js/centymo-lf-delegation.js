// centymo-lf-delegation.js — centymo CSP-safe delegated event handlers.
//
// Plan 6 (CSP) onclick= -> lf.on() sweep. Replaces inline onclick= attributes on
// centymo templates with data-* hooks + document-level delegation, so a strict
// script-src 'self' policy is no longer blocked by inline event-handler attributes.
//
// EXTERNALIZED (stage-2 nonce wave): formerly an inline <script> inside
// views/revenue/templates/lf-delegation.html ({{define "centymo-lf-delegation"}}).
// That partial is a tail include used by ~34 templates whose template dot varies
// (detail pages embedding types.PageData vs. flat form.Data drawer structs), so a
// single nonce="{{.Nonce}}" / nonce="{{$.Nonce}}" could not resolve uniformly —
// the field is absent on many of the includer dots and would panic the render.
// Moving the body to an external /assets/js file referenced via <script src> needs
// NO nonce under strict script-src 'self' and works for every includer regardless
// of its dot. Copied to apps/service-admin/assets/js/centymo/ by
// centymo.CopyStaticAssets at startup (container.go).
//
// This module is provenance-owned by centymo (lf.centymo.*, see
// js-css-architecture plan §"Provenance invariant"). It installs ONCE per session:
// lf.on() binds to document (which survives HTMX OOB swaps), and an internal flag
// short-circuits re-installation when an including partial re-renders inside an
// HTMX-swapped region (the external <script src> re-executes on swap, but the flag
// makes that a no-op).
//
// Hooks (each replaces a former onclick=):
//   data-lf-sheet="open"  + data-lf-sheet-title="<title>"  -> lf.ui.Sheet.open(title)
//   data-lf-sheet="open"  (no title attr)                  -> lf.ui.Sheet.open(triggerEl)
//   data-lf-sheet="close"                                  -> lf.ui.Sheet.close()
//   data-lf-nav-href="<url>"                               -> location.href = url
//   data-lf-stop-propagation                               -> event.stopPropagation()
//   data-lf-mv-clear                                       -> clear movements filters + re-apply
//   data-lf-mv-export (with data-lf-export-url)            -> build query string + navigate to export
//   data-lf-pricelist-product (change)                     -> lf.centymo.pricelist.PriceProductForm.onProductChange(el)
(function () {
    window.lf = window.lf || {};
    if (typeof lf.ns !== 'function' || typeof lf.on !== 'function') { return; }
    lf.ns('centymo');
    // Idempotent: bind the document-level delegated handlers exactly once.
    if (lf.centymo.__delegationInstalled) { return; }
    lf.centymo.__delegationInstalled = true;

    // --- Sheet open / close ------------------------------------------------
    // NOTE: do NOT preventDefault — the original onclick= did not, and these
    // elements also carry hx-get/hx-post that must fire natively. Sheet.open()
    // only sets the drawer title; HTMX performs the content fetch/swap.
    lf.on('click', '[data-lf-sheet="open"]', function () {
        if (!lf.ui || !lf.ui.Sheet) { return; }
        var title = this.getAttribute('data-lf-sheet-title');
        // Title present -> open by title string; absent -> open by trigger element
        // (preserves the former lf.ui.Sheet.open(this) call shape).
        lf.ui.Sheet.open(title !== null ? title : this);
    });
    lf.on('click', '[data-lf-sheet="close"]', function () {
        if (lf.ui && lf.ui.Sheet) { lf.ui.Sheet.close(); }
    });

    // --- Row navigation (former onclick="location.href='...'") -------------
    lf.on('click', '[data-lf-nav-href]', function () {
        var href = this.getAttribute('data-lf-nav-href');
        if (href) { window.location.href = href; }
    });

    // --- Stop propagation (former onclick="event.stopPropagation();") ------
    // Registered in the CAPTURE phase (4th arg true): the document handler then
    // runs BEFORE any ancestor's bubble-phase handler, so stopPropagation() still
    // prevents the click from reaching ancestor listeners — matching the original
    // inline-onclick semantics. The target element (the link) still acts (default
    // navigation is unaffected by stopPropagation).
    lf.on('click', '[data-lf-stop-propagation]', function (e) {
        e.stopPropagation();
    }, true);

    // --- Inventory movements: clear filters --------------------------------
    lf.on('click', '[data-lf-mv-clear]', function () {
        ['filter-date-from', 'filter-date-to', 'filter-location', 'filter-type', 'filter-search']
            .forEach(function (id) {
                var el = document.getElementById(id);
                if (el) { el.value = ''; }
            });
        var apply = document.getElementById('apply-filters-btn');
        if (apply && window.htmx) { window.htmx.trigger(apply, 'click'); }
    });

    // --- Inventory movements: export CSV with current filters --------------
    lf.on('click', '[data-lf-mv-export]', function () {
        var base = this.getAttribute('data-lf-export-url');
        if (!base) { return; }
        var params = new URLSearchParams();
        var map = {
            date_from: 'filter-date-from',
            date_to: 'filter-date-to',
            location: 'filter-location',
            type: 'filter-type',
            search: 'filter-search'
        };
        Object.keys(map).forEach(function (param) {
            var el = document.getElementById(map[param]);
            var v = el ? el.value : '';
            if (v) { params.set(param, v); }
        });
        var qs = params.toString();
        window.location.href = base + (qs ? '?' + qs : '');
    });

    // --- Price-list product select (former onchange=) ----------------------
    // The price-product drawer's product <select> previously carried
    //   onchange="lf.centymo.pricelist.PriceProductForm.onProductChange(this)"
    // The function is defined by that drawer's own inline <script nonce> when
    // the drawer renders; here we only delegate the binding. Scoped to the
    // dedicated data hook (NOT the generic #product_id, which collides with the
    // PO line-item drawer's product select). Resolved lazily at change time, so
    // the drawer's inline script has already defined PriceProductForm by then.
    lf.on('change', '[data-lf-pricelist-product]', function () {
        var ns = lf.centymo && lf.centymo.pricelist && lf.centymo.pricelist.PriceProductForm;
        if (ns && typeof ns.onProductChange === 'function') {
            ns.onProductChange(this);
        }
    });
})();
