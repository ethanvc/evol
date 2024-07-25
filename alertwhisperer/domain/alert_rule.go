package domain

import (
	"context"
	"gorm.io/gorm"
)

type AlertRule struct {
	AlertRuleId     int64   `json:"alert_rule_id" gorm:"primaryKey;autoIncrement"`
	Name            string  `json:"name" gorm:"uniqueIndex:i_u_name;type:varchar(255) not null"`
	Rule            string  `json:"rule" gorm:"type:varchar(4096)"`
	Version         int     `json:"version" gorm:"type:int not null"`
	CreateTime      int64   `json:"create_time" gorm:"type:bigint not null"`
	UpdateTime      int64   `json:"update_time" gorm:"type:bigint not null"`
	Threshold       float64 `json:"threshold" gorm:"type:double not null"`
	DurationSeconds int32   `json:"duration_seconds" gorm:"type:int not null"`
}

func (rule *AlertRule) TableName() string {
	return "alert_rule_tab"
}

type AlertRuleRepository struct {
	db *gorm.DB
}

func NewAlertRuleRepository(db *gorm.DB) *AlertRuleRepository {
	return &AlertRuleRepository{
		db: db,
	}
}

func (repo *AlertRuleRepository) Create(c context.Context, rule *AlertRule) error {
	err := repo.db.WithContext(c).Create(rule).Error
	if err != nil {
		return err
	}
	return nil
}
