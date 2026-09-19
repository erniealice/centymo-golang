package form

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
)

var percentPattern = regexp.MustCompile(`^(\d{1,3})(?:\.(\d{1,2}))?$`)

func ParseEscalationPercentBPS(raw string) (int32, error) {
	value := strings.TrimSpace(raw)
	match := percentPattern.FindStringSubmatch(value)
	if match == nil {
		return 0, errors.New("percentage must use at most two decimal places")
	}
	whole, _ := strconv.Atoi(match[1])
	fraction := match[2]
	if len(fraction) == 1 {
		fraction += "0"
	}
	frac := 0
	if fraction != "" {
		frac, _ = strconv.Atoi(fraction)
	}
	bps := whole*100 + frac
	if bps < 1 || bps > 10000 {
		return 0, errors.New("percentage must be between 0.01 and 100.00")
	}
	return int32(bps), nil
}

func FormatEscalationPercentBPS(bps int32) string {
	if bps <= 0 {
		return ""
	}
	return strconv.FormatInt(int64(bps/100), 10) + "." + strconv.FormatInt(int64((bps%100)/10), 10) + strconv.FormatInt(int64(bps%10), 10)
}

func PopulateDefaultEscalation(data *Data, pp *priceplanpb.PricePlan) {
	if data == nil || pp == nil {
		return
	}
	if pp.DefaultEscalationMode != nil {
		data.DefaultEscalationMode = pp.GetDefaultEscalationMode().String()
	}
	if pp.DefaultEscalationScope != nil {
		data.DefaultEscalationScope = pp.GetDefaultEscalationScope().String()
	}
	if pp.DefaultEscalationRateBps != nil {
		data.DefaultEscalationRate = FormatEscalationPercentBPS(pp.GetDefaultEscalationRateBps())
	}
	if pp.DefaultEscalationFirstAfterMonths != nil {
		data.DefaultEscalationFirstAfterMonths = strconv.FormatInt(int64(pp.GetDefaultEscalationFirstAfterMonths()), 10)
	}
	if pp.DefaultEscalationEveryMonths != nil {
		data.DefaultEscalationEveryMonths = strconv.FormatInt(int64(pp.GetDefaultEscalationEveryMonths()), 10)
	}
}

func ApplyDefaultEscalation(pp *priceplanpb.PricePlan, values url.Values) error {
	if pp == nil {
		return errors.New("price plan is required")
	}
	modeRaw := values.Get("default_escalation_mode")
	if modeRaw == "" {
		mode := priceplanpb.EscalationMode_ESCALATION_MODE_UNSPECIFIED
		pp.DefaultEscalationMode = &mode
		return nil
	}
	modeNumber, ok := priceplanpb.EscalationMode_value[modeRaw]
	if !ok {
		return errors.New("invalid escalation mode")
	}
	mode := priceplanpb.EscalationMode(modeNumber)
	pp.DefaultEscalationMode = &mode
	if mode == priceplanpb.EscalationMode_ESCALATION_MODE_NONE {
		return nil
	}
	if mode != priceplanpb.EscalationMode_ESCALATION_MODE_FIXED_PERCENTAGE {
		return errors.New("invalid escalation mode")
	}

	scopeRaw := values.Get("default_escalation_scope")
	scopeNumber, ok := priceplanpb.EscalationScope_value[scopeRaw]
	if !ok {
		return errors.New("invalid escalation scope")
	}
	scope := priceplanpb.EscalationScope(scopeNumber)
	pp.DefaultEscalationScope = &scope
	rate, err := ParseEscalationPercentBPS(values.Get("default_escalation_rate"))
	if err != nil {
		return err
	}
	pp.DefaultEscalationRateBps = &rate

	if scope == priceplanpb.EscalationScope_ESCALATION_SCOPE_ON_RENEWAL {
		return nil
	}
	if scope != priceplanpb.EscalationScope_ESCALATION_SCOPE_WITHIN_AGREEMENT {
		return errors.New("invalid escalation scope")
	}
	first, err := parsePositiveMonth(values.Get("default_escalation_first_after_months"))
	if err != nil {
		return err
	}
	every, err := parsePositiveMonth(values.Get("default_escalation_every_months"))
	if err != nil {
		return err
	}
	pp.DefaultEscalationFirstAfterMonths = &first
	pp.DefaultEscalationEveryMonths = &every
	return nil
}

func parsePositiveMonth(raw string) (int32, error) {
	value, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 32)
	if err != nil || value < 1 || value > 1200 {
		return 0, errors.New("month interval must be between 1 and 1200")
	}
	return int32(value), nil
}
