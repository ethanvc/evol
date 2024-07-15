package controller

import (
	"context"
	"github.com/ethanvc/evol/base"
	"google.golang.org/grpc/codes"
)

type AlertRuleController struct{}

func NewAlertRuleController() *AlertRuleController {
	return &AlertRuleController{}
}

type AlertRule struct {
	AlertRuleId int64 `json:"alert_rule_id"`
}

func (controller *AlertRuleController) CreateAlertRule(c context.Context, req *AlertRule) (*AlertRule, error) {
	return nil, base.New(codes.Unimplemented, "")
}
