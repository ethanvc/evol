package controller

import (
	"context"
	"github.com/ethanvc/evol/alertwhisperer/repo"
)

type AlertRuleController struct {
	alertRuleRepo *repo.AlertRuleRepository
}

func NewAlertRuleController(alertRuleRepo *repo.AlertRuleRepository) *AlertRuleController {
	return &AlertRuleController{
		alertRuleRepo: alertRuleRepo,
	}
}

type AlertRule struct {
	AlertRuleId     int64   `json:"alert_rule_id"`
	Name            string  `json:"name"`
	Rule            string  `json:"rule"`
	Version         int     `json:"version"`
	CreateTime      int64   `json:"create_time"`
	UpdateTime      int64   `json:"update_time"`
	Threshold       float64 `json:"threshold"`
	DurationSeconds int64   `json:"duration_seconds"`
}

func (controller *AlertRuleController) CreateAlertRule(c context.Context, req *AlertRule) (*AlertRule, error) {
	return &AlertRule{
		AlertRuleId: 1,
	}, nil
}