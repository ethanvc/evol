package controller

import (
	"context"
	"github.com/ethanvc/evol/alertwhisperer/domain"
)

type AlertRuleController struct {
	alertRuleRepo *domain.AlertRuleRepository
}

func NewAlertRuleController(alertRuleRepo *domain.AlertRuleRepository) *AlertRuleController {
	return &AlertRuleController{
		alertRuleRepo: alertRuleRepo,
	}
}

func (controller *AlertRuleController) CreateAlertRule(c context.Context, req *domain.AlertRule) (*domain.AlertRule, error) {
	err := controller.alertRuleRepo.Create(c, req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
