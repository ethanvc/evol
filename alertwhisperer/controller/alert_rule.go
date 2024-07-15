package controller

import (
	"context"
)

type AlertRuleController struct{}

func NewAlertRuleController() *AlertRuleController {
	return &AlertRuleController{}
}

type AlertRule struct {
	AlertRuleId int64 `json:"alert_rule_id"`
}

func (controller *AlertRuleController) CreateAlertRule(c context.Context, req *AlertRule) (*AlertRule, error) {
	return &AlertRule{
		AlertRuleId: 1,
	}, nil
}
