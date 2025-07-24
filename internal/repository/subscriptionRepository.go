package repository

import (
	"subscription-aggregator/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(sub *model.Subscription) error
	GetByID(id uuid.UUID) (*model.Subscription, error)
	Update(sub *model.Subscription) error
	Delete(id uuid.UUID) error
	GetList(filter model.Subscription) ([]model.Subscription, error)
	CalcTotal(userID uuid.UUID, serviceName string, from, to *time.Time) (uint, error)
}

type subscriptionRepo struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepo{db: db}
}

func (r *subscriptionRepo) Create(sub *model.Subscription) error {
	return r.db.Create(sub).Error
}

func (r *subscriptionRepo) GetByID(id uuid.UUID) (*model.Subscription, error) {
	var sub model.Subscription
	if err := r.db.First(&sub, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepo) Update(sub *model.Subscription) error {
	return r.db.Save(sub).Error
}

func (r *subscriptionRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Subscription{}, "id = ?", id).Error
}

func (r *subscriptionRepo) GetList(filter model.Subscription) ([]model.Subscription, error) {
	var subs []model.Subscription
	query := r.db.Model(&model.Subscription{})

	if filter.ID.String() != "" {
		query = query.Where("id = ?", filter.ID)
	}

	if err := query.Find(&subs).Error; err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *subscriptionRepo) CalcTotal(userID uuid.UUID, serviceName string, from, to *time.Time) (uint, error) {
	var total uint
	query := r.db.Model(&model.Subscription{}).Where("id = ?", userID)
	if query.Error != nil {
		return 0, query.Error
	}
	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}
	if from != nil {
		query = query.Where("start_date >= ?", *from)
	}
	if to != nil {
		query = query.Where("start_date <= ?", *to)
	}

	err := query.Select("SUM(price)").Scan(&total).Error
	return total, err
}
