package detail

// eligibility_rules_tab.go — Stage 2 (treasury-domain-rebuild) Eligibility
// Rules tab handler for the Collection Method detail page (pages.md §B-5 tab 2).
//
// The tab lists collection_method_eligibility_rule rows and provides a
// create/edit/delete drawer-form (entities.md §E-3 field set). All CRUD
// closures are nil-safe; when the espyna use cases are absent the tab renders
// an empty state rather than panicking.

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	eligrulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection_method_eligibility_rule"
)

// EligibilityRuleTabDeps holds the closures needed by the Eligibility Rules tab.
// All closures are optional (nil-safe): the tab renders an empty/disabled state
// when a closure is absent rather than panicking.
type EligibilityRuleTabDeps struct {
	Routes       centymo.CollectionMethodRoutes
	Labels       centymo.CollectionMethodLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ListCollectionMethodEligibilityRules  func(ctx context.Context, req *eligrulepb.ListCollectionMethodEligibilityRulesRequest) (*eligrulepb.ListCollectionMethodEligibilityRulesResponse, error)
	ReadCollectionMethodEligibilityRule   func(ctx context.Context, req *eligrulepb.ReadCollectionMethodEligibilityRuleRequest) (*eligrulepb.ReadCollectionMethodEligibilityRuleResponse, error)
	CreateCollectionMethodEligibilityRule func(ctx context.Context, req *eligrulepb.CreateCollectionMethodEligibilityRuleRequest) (*eligrulepb.CreateCollectionMethodEligibilityRuleResponse, error)
	UpdateCollectionMethodEligibilityRule func(ctx context.Context, req *eligrulepb.UpdateCollectionMethodEligibilityRuleRequest) (*eligrulepb.UpdateCollectionMethodEligibilityRuleResponse, error)
	DeleteCollectionMethodEligibilityRule func(ctx context.Context, req *eligrulepb.DeleteCollectionMethodEligibilityRuleRequest) (*eligrulepb.DeleteCollectionMethodEligibilityRuleResponse, error)
}

// EligibilityRuleRow is a flat representation of one rule for the template.
type EligibilityRuleRow struct {
	ID               string
	Name             string
	BearerMode       string
	ValidFromDate    string
	ValidUntilDate   string
	ExpiryDays       string
	MinAmount        string
	MaxAmount        string
	StackingPolicy   string
	JurisdictionCode string
	MaxPerInstance   string
	MaxPerClient     string
	TermsURL         string
	TermsSummary     string

	EditURL   string
	DeleteURL string
}

// EligibilityRuleFormData holds the data for the create/edit drawer-form.
type EligibilityRuleFormData struct {
	Labels       centymo.CollectionMethodEligibilityRuleLabels
	CommonLabels pyeza.CommonLabels

	IsEdit     bool
	ID         string
	MethodID   string
	FormAction string

	Name                      string
	BearerMode                string
	ValidFromDate             string
	ValidUntilDate            string
	ExpiryDaysAfterIssuance   string
	MinAmountCentavos         string
	MaxAmountCentavos         string
	StackingPolicy            string
	JurisdictionCode          string
	MaxRedemptionsPerInstance string
	MaxRedemptionsPerClient   string
	TermsURL                  string
	TermsSummary              string

	BearerModeOptions []types.SelectOption
	StackingOptions   []types.SelectOption
}

// EligibilityRuleTabData is passed to the eligibility-rules-tab template.
type EligibilityRuleTabData struct {
	Labels       centymo.CollectionMethodEligibilityRuleLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	MethodID string
	Rules    []EligibilityRuleRow
	AddURL   string
}

// eligibilityRuleTableID is the HTMX target used by CRUD success responses.
const eligibilityRuleTableID = "eligibility-rules-table"

// NewEligibilityRuleTabView returns the tab-content partial for the Eligibility
// Rules tab. The outer detail page delegates here when activeTab == "eligibility".
func NewEligibilityRuleTabView(deps *EligibilityRuleTabDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection_method_eligibility_rule", "list") {
			return centymo.HTMXError(fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "collection_method_eligibility_rule:list"))
		}

		methodID := viewCtx.Request.PathValue("id")
		if methodID == "" {
			return view.Error(fmt.Errorf("missing method id"))
		}

		rules, err := loadEligibilityRules(ctx, deps, methodID)
		if err != nil {
			log.Printf("ListCollectionMethodEligibilityRules method=%s: %v", methodID, err)
			return view.Error(fmt.Errorf("failed to load eligibility rules: %w", err))
		}

		addURL := route.ResolveURL(deps.Routes.EligibilityRuleAddURL, "method_id", methodID)

		td := &EligibilityRuleTabData{
			Labels:       deps.Labels.EligibilityRule,
			CommonLabels: deps.CommonLabels,
			TableLabels:  deps.TableLabels,
			MethodID:     methodID,
			Rules:        rules,
			AddURL:       addURL,
		}
		return view.OK("collection-method-eligibility-rules-tab", td)
	})
}

// NewEligibilityRuleAddAction handles GET + POST for the add drawer-form.
func NewEligibilityRuleAddAction(deps *EligibilityRuleTabDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection_method_eligibility_rule", "create") {
			return centymo.HTMXError(fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "collection_method_eligibility_rule:create"))
		}
		methodID := viewCtx.Request.PathValue("method_id")
		if methodID == "" {
			return view.Error(fmt.Errorf("missing method_id"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			fd := buildEmptyEligibilityRuleFormData(deps, methodID)
			fd.FormAction = route.ResolveURL(deps.Routes.EligibilityRuleAddURL, "method_id", methodID)
			return view.OK("collection-method-eligibility-rule-drawer-form", fd)
		}

		if deps.CreateCollectionMethodEligibilityRule == nil {
			return centymo.HTMXError("Eligibility rule create is not wired yet.")
		}
		rule, err := parseEligibilityRuleForm(viewCtx.Request, "")
		if err != nil {
			return view.Error(err)
		}
		if _, err := deps.CreateCollectionMethodEligibilityRule(ctx, &eligrulepb.CreateCollectionMethodEligibilityRuleRequest{
			Data: rule,
		}); err != nil {
			log.Printf("CreateCollectionMethodEligibilityRule method=%s: %v", methodID, err)
			return view.Error(fmt.Errorf("failed to create eligibility rule: %w", err))
		}
		return centymo.HTMXSuccess(eligibilityRuleTableID)
	})
}

// NewEligibilityRuleEditAction handles GET + POST for the edit drawer-form.
func NewEligibilityRuleEditAction(deps *EligibilityRuleTabDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection_method_eligibility_rule", "update") {
			return centymo.HTMXError(fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "collection_method_eligibility_rule:update"))
		}
		methodID := viewCtx.Request.PathValue("method_id")
		ruleID := viewCtx.Request.PathValue("rule_id")
		if methodID == "" || ruleID == "" {
			return view.Error(fmt.Errorf("missing method_id or rule_id"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			if deps.ReadCollectionMethodEligibilityRule == nil {
				return centymo.HTMXError("Eligibility rule read is not wired yet.")
			}
			resp, err := deps.ReadCollectionMethodEligibilityRule(ctx, &eligrulepb.ReadCollectionMethodEligibilityRuleRequest{
				Data: &eligrulepb.CollectionMethodEligibilityRule{Id: ruleID},
			})
			if err != nil {
				return view.Error(fmt.Errorf("failed to read eligibility rule: %w", err))
			}
			data := resp.GetData()
			if len(data) == 0 {
				return view.Error(fmt.Errorf("eligibility rule not found"))
			}
			fd := buildFormDataFromRule(deps, methodID, data[0])
			fd.FormAction = route.ResolveURL(
				route.ResolveURL(deps.Routes.EligibilityRuleEditURL, "method_id", methodID),
				"rule_id", ruleID,
			)
			fd.IsEdit = true
			fd.ID = ruleID
			return view.OK("collection-method-eligibility-rule-drawer-form", fd)
		}

		rule, err := parseEligibilityRuleForm(viewCtx.Request, ruleID)
		if err != nil {
			return view.Error(err)
		}
		if deps.UpdateCollectionMethodEligibilityRule == nil {
			return centymo.HTMXError("Eligibility rule update is not wired yet.")
		}
		if _, err := deps.UpdateCollectionMethodEligibilityRule(ctx, &eligrulepb.UpdateCollectionMethodEligibilityRuleRequest{
			Data: rule,
		}); err != nil {
			log.Printf("UpdateCollectionMethodEligibilityRule %s: %v", ruleID, err)
			return view.Error(fmt.Errorf("failed to update eligibility rule: %w", err))
		}
		return centymo.HTMXSuccess(eligibilityRuleTableID)
	})
}

// NewEligibilityRuleDeleteAction handles POST for the delete action.
func NewEligibilityRuleDeleteAction(deps *EligibilityRuleTabDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection_method_eligibility_rule", "delete") {
			return centymo.HTMXError(fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "collection_method_eligibility_rule:delete"))
		}
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.FormValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		if deps.DeleteCollectionMethodEligibilityRule == nil {
			return centymo.HTMXError("Eligibility rule delete is not wired yet.")
		}
		if _, err := deps.DeleteCollectionMethodEligibilityRule(ctx, &eligrulepb.DeleteCollectionMethodEligibilityRuleRequest{
			Data: &eligrulepb.CollectionMethodEligibilityRule{Id: id},
		}); err != nil {
			log.Printf("DeleteCollectionMethodEligibilityRule %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to delete eligibility rule: %w", err))
		}
		return centymo.HTMXSuccess(eligibilityRuleTableID)
	})
}

// --- helpers -----------------------------------------------------------------

func loadEligibilityRules(ctx context.Context, deps *EligibilityRuleTabDeps, methodID string) ([]EligibilityRuleRow, error) {
	if deps.ListCollectionMethodEligibilityRules == nil {
		return nil, nil
	}
	resp, err := deps.ListCollectionMethodEligibilityRules(ctx, &eligrulepb.ListCollectionMethodEligibilityRulesRequest{})
	if err != nil {
		return nil, err
	}
	var rows []EligibilityRuleRow
	for _, r := range resp.GetData() {
		editURL := route.ResolveURL(
			route.ResolveURL(deps.Routes.EligibilityRuleEditURL, "method_id", methodID),
			"rule_id", r.GetId(),
		)
		deleteURL := route.ResolveURL(deps.Routes.EligibilityRuleDeleteURL, "method_id", methodID)
		row := EligibilityRuleRow{
			ID:               r.GetId(),
			Name:             r.GetName(),
			BearerMode:       eligBearerModeShort(r.GetBearerMode().String()),
			ValidFromDate:    r.GetValidFromDate(),
			ValidUntilDate:   r.GetValidUntilDate(),
			StackingPolicy:   eligStackingShort(r.GetStackingPolicy().String()),
			JurisdictionCode: r.GetJurisdictionCode(),
			TermsURL:         r.GetTermsUrl(),
			TermsSummary:     r.GetTermsSummary(),
			EditURL:          editURL,
			DeleteURL:        deleteURL,
		}
		if v := r.GetExpiryDaysAfterIssuance(); v > 0 {
			row.ExpiryDays = strconv.FormatInt(int64(v), 10)
		}
		if v := r.GetMinAmountCentavos(); v > 0 {
			row.MinAmount = formatCentavos(v)
		}
		if v := r.GetMaxAmountCentavos(); v > 0 {
			row.MaxAmount = formatCentavos(v)
		}
		if v := r.GetMaxRedemptionsPerInstance(); v > 0 {
			row.MaxPerInstance = strconv.FormatInt(int64(v), 10)
		}
		if v := r.GetMaxRedemptionsPerClient(); v > 0 {
			row.MaxPerClient = strconv.FormatInt(int64(v), 10)
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func buildEmptyEligibilityRuleFormData(deps *EligibilityRuleTabDeps, methodID string) *EligibilityRuleFormData {
	l := deps.Labels.EligibilityRule
	fd := &EligibilityRuleFormData{
		Labels:       l,
		CommonLabels: deps.CommonLabels,
		MethodID:     methodID,
	}
	fd.BearerModeOptions = []types.SelectOption{
		{Value: "COLLECTION_METHOD_ELIGIBILITY_RULE_BEARER_MODE_HOLDER_BOUND", Label: l.Form.BearerModeHolderBound},
		{Value: "COLLECTION_METHOD_ELIGIBILITY_RULE_BEARER_MODE_HOLDER_TRANSFERABLE", Label: l.Form.BearerModeTransferable},
	}
	fd.StackingOptions = []types.SelectOption{
		{Value: "COLLECTION_METHOD_ELIGIBILITY_RULE_STACKING_POLICY_EXCLUSIVE", Label: l.Form.StackingExclusive},
		{Value: "COLLECTION_METHOD_ELIGIBILITY_RULE_STACKING_POLICY_STACKABLE", Label: l.Form.StackingStackable},
		{Value: "COLLECTION_METHOD_ELIGIBILITY_RULE_STACKING_POLICY_FIRST_ONLY", Label: l.Form.StackingFirstOnly},
	}
	return fd
}

func buildFormDataFromRule(deps *EligibilityRuleTabDeps, methodID string, r *eligrulepb.CollectionMethodEligibilityRule) *EligibilityRuleFormData {
	fd := buildEmptyEligibilityRuleFormData(deps, methodID)
	fd.Name = r.GetName()
	fd.BearerMode = r.GetBearerMode().String()
	fd.ValidFromDate = r.GetValidFromDate()
	fd.ValidUntilDate = r.GetValidUntilDate()
	fd.StackingPolicy = r.GetStackingPolicy().String()
	fd.JurisdictionCode = r.GetJurisdictionCode()
	fd.TermsURL = r.GetTermsUrl()
	fd.TermsSummary = r.GetTermsSummary()
	if v := r.GetExpiryDaysAfterIssuance(); v > 0 {
		fd.ExpiryDaysAfterIssuance = strconv.FormatInt(int64(v), 10)
	}
	if v := r.GetMinAmountCentavos(); v > 0 {
		fd.MinAmountCentavos = formatCentavos(v)
	}
	if v := r.GetMaxAmountCentavos(); v > 0 {
		fd.MaxAmountCentavos = formatCentavos(v)
	}
	if v := r.GetMaxRedemptionsPerInstance(); v > 0 {
		fd.MaxRedemptionsPerInstance = strconv.FormatInt(int64(v), 10)
	}
	if v := r.GetMaxRedemptionsPerClient(); v > 0 {
		fd.MaxRedemptionsPerClient = strconv.FormatInt(int64(v), 10)
	}
	return fd
}

func parseEligibilityRuleForm(r *http.Request, id string) (*eligrulepb.CollectionMethodEligibilityRule, error) {
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("parse form: %w", err)
	}
	rule := &eligrulepb.CollectionMethodEligibilityRule{
		Id:   id,
		Name: r.FormValue("name"),
		BearerMode: eligrulepb.CollectionMethodEligibilityRuleBearerMode(
			eligrulepb.CollectionMethodEligibilityRuleBearerMode_value[r.FormValue("bearer_mode")],
		),
		StackingPolicy: eligrulepb.CollectionMethodEligibilityRuleStackingPolicy(
			eligrulepb.CollectionMethodEligibilityRuleStackingPolicy_value[r.FormValue("stacking_policy")],
		),
	}
	if v := r.FormValue("valid_from_date"); v != "" {
		rule.ValidFromDate = &v
	}
	if v := r.FormValue("valid_until_date"); v != "" {
		rule.ValidUntilDate = &v
	}
	if v := r.FormValue("expiry_days_after_issuance"); v != "" {
		if d, err := strconv.ParseInt(v, 10, 32); err == nil && d > 0 {
			dd := int32(d)
			rule.ExpiryDaysAfterIssuance = &dd
		}
	}
	if v := r.FormValue("min_amount_centavos"); v != "" {
		if cents := parseCentavosFromForm(v); cents > 0 {
			rule.MinAmountCentavos = &cents
		}
	}
	if v := r.FormValue("max_amount_centavos"); v != "" {
		if cents := parseCentavosFromForm(v); cents > 0 {
			rule.MaxAmountCentavos = &cents
		}
	}
	if v := r.FormValue("jurisdiction_code"); v != "" {
		rule.JurisdictionCode = &v
	}
	if v := r.FormValue("max_redemptions_per_instance"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 32); err == nil && n > 0 {
			nn := int32(n)
			rule.MaxRedemptionsPerInstance = &nn
		}
	}
	if v := r.FormValue("max_redemptions_per_client"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 32); err == nil && n > 0 {
			nn := int32(n)
			rule.MaxRedemptionsPerClient = &nn
		}
	}
	if v := r.FormValue("terms_url"); v != "" {
		rule.TermsUrl = &v
	}
	if v := r.FormValue("terms_summary"); v != "" {
		rule.TermsSummary = &v
	}
	return rule, nil
}

func parseCentavosFromForm(s string) int64 {
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	cents := int64(f * 100)
	if cents < 0 {
		return 0
	}
	return cents
}

// eligBearerModeShort strips the long proto prefix from a bearer mode enum string.
func eligBearerModeShort(s string) string {
	const pfx = "COLLECTION_METHOD_ELIGIBILITY_RULE_BEARER_MODE_"
	if len(s) > len(pfx) {
		return s[len(pfx):]
	}
	return s
}

// eligStackingShort strips the long proto prefix from a stacking policy enum string.
func eligStackingShort(s string) string {
	const pfx = "COLLECTION_METHOD_ELIGIBILITY_RULE_STACKING_POLICY_"
	if len(s) > len(pfx) {
		return s[len(pfx):]
	}
	return s
}
