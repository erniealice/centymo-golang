// subscription-drawer.js — centymo subscription (enrollment / agreement) drawer wiring.
//
// EXTERNALIZED from two inline <script nonce="{{$.Nonce}}"> blocks in
// domain/subscription/subscription/templates/subscription-drawer-form.html. The
// drawer is loaded into #sheetContent by HTMX; HTMX re-creates its inline scripts
// under the PAGE's CSP, whose nonce differs from the partial response's, so the
// enforcing script-src blocked them (MMIS prod 2026-10-23): the escalation inputs
// stayed `required` + enabled (Save silently did nothing) and the typed start date
// never reached date_time_start_iso. An external <script src> needs no nonce under
// script-src 'self'. Copied to apps/<app>/assets/js/centymo/ at startup.
//
// The <script src> re-executes on every drawer swap; each form is wired once
// (data-subscription-drawer-ready), and htmx:load covers a form swapped in after
// the script already ran.
//
// Form hooks (all on <form data-subscription-drawer>):
//   data-default-tz            IANA zone for the ISO offset (fallback Asia/Manila)
//   data-subscription-date-sync  present when the date→ISO / plan-search wiring applies
//                              (absent for the commencement-pricing SuggestPlanTerm form)
//   data-billing-currency, data-client-id, data-client-name  plan-search filters
//   data-plan-placeholder      text restored into the plan picker when cleared
(function () {
    'use strict';

    var FIXED = 'ESCALATION_MODE_FIXED_PERCENTAGE';
    var WITHIN = 'ESCALATION_SCOPE_WITHIN_AGREEMENT';

    // Escalation intervals: shown (and submitted) only for a fixed percentage
    // applied within the agreement. The fixed-fields section itself is handled by
    // pyeza's data-lf-show-when helper; this runs on the next tick so it always
    // has the last word over that helper for the nested interval inputs.
    function wireEscalation(form) {
        var mode = form.querySelector('[name="escalation_mode"]');
        var scope = form.querySelector('[name="escalation_scope"]');
        var intervals = form.querySelector('[data-escalation-interval-fields]');
        if (!mode || !scope || !intervals) { return; }
        function apply() {
            var on = mode.value === FIXED && scope.value === WITHIN;
            intervals.hidden = !on;
            intervals.querySelectorAll('input, select').forEach(function (input) { input.disabled = !on; });
        }
        function later() { setTimeout(apply, 0); }
        mode.addEventListener('change', later);
        scope.addEventListener('change', later);
        later();
    }

    // Compute the offset (e.g. "+08:00") for the chosen IANA tz at the given
    // wall-clock moment. Fallback to local browser offset if the Intl API
    // can't render the named zone (very rare; the OS tzdata is missing it).
    function offsetFor(tz, year, month, day, hour, minute) {
        try {
            var dtf = new Intl.DateTimeFormat('en-US', {
                timeZone: tz, hourCycle: 'h23',
                year: 'numeric', month: '2-digit', day: '2-digit',
                hour: '2-digit', minute: '2-digit', second: '2-digit'
            });
            // Build a UTC instant for the wall-clock time, then ask what time
            // it is in tz; the difference IS the offset for that wall clock.
            var asUTC = Date.UTC(year, month - 1, day, hour, minute, 0);
            var parts = dtf.formatToParts(new Date(asUTC));
            var lookup = {};
            parts.forEach(function (p) { lookup[p.type] = p.value; });
            var asTZ = Date.UTC(
                Number(lookup.year), Number(lookup.month) - 1, Number(lookup.day),
                Number(lookup.hour) % 24, Number(lookup.minute), Number(lookup.second)
            );
            // Positive RFC3339 offsets mean "ahead of UTC"; the delta is
            // wall-clock(UTC) - wall-clock(tz), so flip the sign.
            var offsetMinutes = -Math.round((asUTC - asTZ) / 60000);
            var sign = offsetMinutes >= 0 ? '+' : '-';
            var abs = Math.abs(offsetMinutes);
            return sign + String(Math.floor(abs / 60)).padStart(2, '0') + ':' + String(abs % 60).padStart(2, '0');
        } catch (e) {
            var local = new Date(year, month - 1, day, hour, minute, 0);
            var offMin = -local.getTimezoneOffset();
            var sgn = offMin >= 0 ? '+' : '-';
            var a = Math.abs(offMin);
            return sgn + String(Math.floor(a / 60)).padStart(2, '0') + ':' + String(a % 60).padStart(2, '0');
        }
    }

    function buildISO(tz, dateInput, timeInput) {
        var d = dateInput && dateInput.value;
        if (!d) { return ''; }
        var t = (timeInput && timeInput.value) || '00:00';
        var parts = d.split('-');
        if (parts.length !== 3) { return ''; }
        var y = Number(parts[0]), m = Number(parts[1]), day = Number(parts[2]);
        var tparts = t.split(':');
        var hh = Number(tparts[0]) || 0;
        var mm = Number(tparts[1]) || 0;
        var pad = function (n) { return String(n).padStart(2, '0'); };
        return parts[0] + '-' + pad(m) + '-' + pad(day) + 'T' + pad(hh) + ':' + pad(mm) + ':00' + offsetFor(tz, y, m, day, hh, mm);
    }

    // Date+time → ISO hidden fields, and the plan picker's search URL (client /
    // currency / validity filters). Date→ISO is required regardless of plan picker
    // mode; the picker is absent when the plan is locked.
    function wireDates(form) {
        var startDate = form.querySelector('[name="date_start_date"]');
        var startTime = form.querySelector('[name="date_start_time"]');
        var endDate = form.querySelector('[name="date_end_date"]');
        var endTime = form.querySelector('[name="date_end_time"]');
        var startISO = form.querySelector('[name="date_time_start_iso"]');
        var endISO = form.querySelector('[name="date_time_end_iso"]');
        if (!startDate || !endDate || !startISO || !endISO) { return; }

        var ds = form.dataset;
        var tz = ds.defaultTz || 'Asia/Manila';
        var planContainer = form.querySelector('#plan-search');
        var planHidden = planContainer ? planContainer.querySelector('.auto-complete-input') : null;
        var planTrigger = planContainer ? planContainer.querySelector('.auto-complete-trigger') : null;
        var baseURL = planContainer ? (planContainer.dataset.searchUrl || '') : '';

        function refreshHidden() {
            startISO.value = buildISO(tz, startDate, startTime);
            endISO.value = buildISO(tz, endDate, endTime);
        }

        function updatePlanSearchURL() {
            if (!planContainer) { return; }
            var params = [
                ['date_time_start_iso', startISO.value],
                ['date_time_end_iso', endISO.value],
                ['billing_currency', ds.billingCurrency],
                ['client_id', ds.clientId],
                ['client_name', ds.clientName]
            ];
            var url = baseURL;
            var sep = '?';
            params.forEach(function (p) {
                if (p[1]) {
                    url = url + sep + p[0] + '=' + encodeURIComponent(p[1]);
                    sep = '&';
                }
            });
            planContainer.dataset.searchUrl = url;
        }

        function clearPlanSelection() {
            if (!planContainer || !planHidden || !planTrigger) { return; }
            planHidden.value = '';
            planHidden.removeAttribute('value');
            var display = planTrigger.querySelector('.auto-complete-display');
            if (display) {
                var placeholder = document.createElement('span');
                placeholder.className = 'auto-complete-placeholder';
                placeholder.textContent = ds.planPlaceholder || '';
                display.replaceChildren(placeholder);
            }
            var clearBtn = planContainer.querySelector('.auto-complete-clear');
            if (clearBtn) { clearBtn.style.display = 'none'; }
        }

        function onDateTimeChange() {
            refreshHidden();
            updatePlanSearchURL();
            clearPlanSelection();
            // If the dropdown is open, force a refetch by toggling closed→open:
            // auto-complete's openDropdown() re-runs the search with the new URL.
            if (planContainer && planContainer.classList.contains('open') && planTrigger) {
                planTrigger.click();
                planTrigger.click();
            }
        }

        [startDate, startTime, endDate, endTime].forEach(function (el) {
            if (el) { el.addEventListener('change', onDateTimeChange); }
        });
        refreshHidden();
        updatePlanSearchURL();
    }

    function init(root) {
        (root || document).querySelectorAll('form[data-subscription-drawer]').forEach(function (form) {
            if (form.dataset.subscriptionDrawerReady) { return; }
            form.dataset.subscriptionDrawerReady = '1';
            wireEscalation(form);
            if (form.hasAttribute('data-subscription-date-sync')) { wireDates(form); }
        });
    }

    init(document);
    if (!window.__lfSubscriptionDrawerBound) {
        window.__lfSubscriptionDrawerBound = true;
        document.addEventListener('htmx:load', function () { init(document); });
        document.addEventListener('DOMContentLoaded', function () { init(document); });
    }
}());
