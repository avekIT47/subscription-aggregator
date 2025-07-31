package service

import (
	"subscription-aggregator/internal/model"
	"subscription-aggregator/internal/repository"
	"subscription-aggregator/pkg/logger"
	"time"

	"github.com/google/uuid"
)

type SubscriptionService interface {
	Create(sub *model.Subscription) error
	GetByID(id uuid.UUID) (*model.Subscription, error)
	Update(sub *model.Subscription) error
	Delete(id uuid.UUID) error
	GetList(filter model.Subscription, offset, limit int) ([]model.Subscription, error)
	GetTotal(userID string, serviceName string, from, to *time.Time) (uint, error)
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	logger.Log.Info("Creating new SubscriptionService")
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) Create(sub *model.Subscription) error {
	logger.Log.Infof("Service: creating subscription for user %s, service %s", sub.UserID, sub.ServiceName)
	err := s.repo.Create(sub)
	if err != nil {
		logger.Log.Errorf("Service: error creating subscription: %v", err)
	} else {
		logger.Log.Infof("Service: subscription created with ID %s", sub.ID)
	}
	return err
}

func (s *subscriptionService) GetByID(id uuid.UUID) (*model.Subscription, error) {
	logger.Log.Infof("Service: getting subscription by ID %s", id)
	sub, err := s.repo.GetByID(id)
	if err != nil {
		logger.Log.Errorf("Service: subscription with ID %s not found: %v", id, err)
		return nil, err
	}
	logger.Log.Infof("Service: subscription with ID %s retrieved", id)
	return sub, nil
}

func (s *subscriptionService) Update(sub *model.Subscription) error {
	logger.Log.Infof("Service: updating subscription with ID %s", sub.ID)
	err := s.repo.Update(sub)
	if err != nil {
		logger.Log.Errorf("Service: error updating subscription ID %s: %v", sub.ID, err)
	} else {
		logger.Log.Infof("Service: subscription ID %s updated successfully", sub.ID)
	}
	return err
}

func (s *subscriptionService) Delete(id uuid.UUID) error {
	logger.Log.Infof("Service: deleting subscription with ID %s", id)
	err := s.repo.Delete(id)
	if err != nil {
		logger.Log.Errorf("Service: error deleting subscription ID %s: %v", id, err)
	} else {
		logger.Log.Infof("Service: subscription ID %s deleted successfully", id)
	}
	return err
}

func (s *subscriptionService) GetList(filter model.Subscription, offset, limit int) ([]model.Subscription, error) {
	logger.Log.Infof("Service: getting subscription list for user %s with offset %d and limit %d", filter.UserID, offset, limit)
	subs, err := s.repo.GetList(filter, offset, limit)
	if err != nil {
		logger.Log.Errorf("Service: error getting subscription list: %v", err)
		return nil, err
	}
	logger.Log.Infof("Service: retrieved %d subscriptions", len(subs))
	return subs, nil
}

func (s *subscriptionService) GetTotal(userID string, serviceName string, from, to *time.Time) (uint, error) {
	logger.Log.Infof("Service: calculating total for user %s, service %s, from %v to %v", userID, serviceName, from, to)
	total, err := s.repo.CalcTotal(userID, serviceName, from, to)
	if err != nil {
		logger.Log.Errorf("Service: error calculating total: %v", err)
		return 0, err
	}
	logger.Log.Infof("Service: total calculated: %d", total)
	return total, nil
}
