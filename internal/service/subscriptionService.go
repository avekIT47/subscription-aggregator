package service

import (
	"subscription-aggregator/internal/model"
	"subscription-aggregator/internal/repository"
	"time"

	"github.com/google/uuid"
)

type SubscriptionService interface {
	Create(sub *model.Subscription) error
	GetByID(id uuid.UUID) (*model.Subscription, error)
	Update(sub *model.Subscription) error
	Delete(id uuid.UUID) error
	GetList(filter model.Subscription) ([]model.Subscription, error)
	GetTotal(userID uuid.UUID, serviceName string, from, to *time.Time) (uint, error)
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) Create(sub *model.Subscription) error {
	return s.repo.Create(sub)
}

func (s *subscriptionService) GetByID(id uuid.UUID) (*model.Subscription, error) {
	return s.repo.GetByID(id)
}

func (s *subscriptionService) Update(sub *model.Subscription) error {
	return s.repo.Update(sub)
}

func (s *subscriptionService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *subscriptionService) GetList(filter model.Subscription) ([]model.Subscription, error) {
	return s.repo.GetList(filter)
}

func (s *subscriptionService) GetTotal(userID uuid.UUID, serviceName string, from, to *time.Time) (uint, error) {
	return s.repo.CalcTotal(userID, serviceName, from, to)
}
