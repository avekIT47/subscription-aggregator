package repository

import (
	"subscription-aggregator/internal/model"
	"subscription-aggregator/pkg/logger"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(sub *model.Subscription) error
	GetByID(id uuid.UUID) (*model.Subscription, error)
	Update(sub *model.Subscription) error
	Delete(id uuid.UUID) error
	GetList(filter model.Subscription, offset, limit int) ([]model.Subscription, error)
	CalcTotal(userID string, serviceName string, from, to *time.Time) (uint, error)
}

type subscriptionRepo struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	logger.Log.Info("Creating new SubscriptionRepository")
	return &subscriptionRepo{db: db}
}

func (r *subscriptionRepo) Create(sub *model.Subscription) error {
	logger.Log.Infof("Creating subscription for user %s, service %s", sub.UserID, sub.ServiceName)
	err := r.db.Create(sub).Error
	if err != nil {
		logger.Log.Errorf("Error creating subscription: %v", err)
	} else {
		logger.Log.Infof("Subscription created with ID %s", sub.ID)
	}
	return err
}

func (r *subscriptionRepo) GetByID(id uuid.UUID) (*model.Subscription, error) {
	logger.Log.Infof("Getting subscription by ID %s", id)
	var sub model.Subscription
	err := r.db.First(&sub, "id = ?", id).Error
	if err != nil {
		logger.Log.Errorf("Subscription with ID %s not found: %v", id, err)
		return nil, err
	}
	logger.Log.Infof("Subscription with ID %s retrieved", id)
	return &sub, nil
}

func (r *subscriptionRepo) Update(sub *model.Subscription) error {
	logger.Log.Infof("Updating subscription with ID %s", sub.ID)
	err := r.db.Save(sub).Error
	if err != nil {
		logger.Log.Errorf("Error updating subscription ID %s: %v", sub.ID, err)
	} else {
		logger.Log.Infof("Subscription ID %s updated successfully", sub.ID)
	}
	return err
}

func (r *subscriptionRepo) Delete(id uuid.UUID) error {
	logger.Log.Infof("Deleting subscription with ID %s", id)
	err := r.db.Delete(&model.Subscription{}, "id = ?", id).Error
	if err != nil {
		logger.Log.Errorf("Error deleting subscription ID %s: %v", id, err)
	} else {
		logger.Log.Infof("Subscription ID %s deleted successfully", id)
	}
	return err
}

func (r *subscriptionRepo) GetList(filter model.Subscription, offset, limit int) ([]model.Subscription, error) {
	logger.Log.Infof("Getting subscriptions list for user %s with offset %d and limit %d", filter.UserID, offset, limit)
	var subs []model.Subscription
	query := r.db.Model(&model.Subscription{})

	if filter.ID.String() != "" {
		query = query.Where("user_id = ?", filter.UserID)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset >= 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&subs).Error
	if err != nil {
		logger.Log.Errorf("Error retrieving subscriptions list: %v", err)
		return nil, err
	}

	logger.Log.Infof("Retrieved %d subscriptions", len(subs))
	return subs, nil
}

func diffMonths(from, to time.Time) int {
	yearDiff := to.Year() - from.Year()
	monthDiff := int(to.Month()) - int(from.Month())

	months := yearDiff*12 + monthDiff

	if to.Day() >= from.Day() {
		months++
	}

	logger.Log.Debugf("diffMonths: from %s to %s = %d months", from.Format("2006-01-02"), to.Format("2006-01-02"), months)
	return months
}

func (r *subscriptionRepo) CalcTotal(userID string, serviceName string, from, to *time.Time) (uint, error) {
	logger.Log.Infof("Calculating total subscription cost for user %s, service %s, from %v to %v", userID, serviceName, from, to)
	var subs []model.Subscription

	query := r.db.Model(&model.Subscription{}).Where("user_id = ?", userID)

	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	if from != nil && to != nil {
		query = query.Where(`
			start_date <= ? AND (end_date IS NULL OR end_date >= ?)`,
			*to, *from,
		)
	} else if from != nil {
		query = query.Where("end_date IS NULL OR end_date >= ?", *from)
	} else if to != nil {
		query = query.Where("start_date <= ?", *to)
	}

	err := query.Find(&subs).Error
	if err != nil {
		logger.Log.Errorf("Error querying subscriptions for total calculation: %v", err)
		return 0, err
	}

	logger.Log.Infof("Found %d subscriptions to process for total calculation", len(subs))

	var total uint
	for _, sub := range subs {
		start := sub.StartDate
		end := time.Now()
		if sub.EndDate != nil {
			end = *sub.EndDate
		}

		if from != nil && start.Before(*from) {
			start = *from
		}
		if to != nil && end.After(*to) {
			end = *to
		}

		if end.Before(start) {
			logger.Log.Warnf("Subscription ID %s: end date %s before start date %s after filter adjustment, skipping", sub.ID, end, start)
			continue
		}

		months := diffMonths(start, end)
		if months == 0 {
			months = 1
		}

		logger.Log.Debugf("Subscription ID %s: price %d x months %d = %d", sub.ID, sub.Price, months, sub.Price*uint(months))
		total += sub.Price * uint(months)
	}

	logger.Log.Infof("Total subscription cost calculated: %d", total)
	return total, nil
}
