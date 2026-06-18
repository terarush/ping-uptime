package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"ping-uptime/modules/subscribers/domain/entity"
	"ping-uptime/modules/subscribers/domain/repository"
)

var (
	ErrAlreadySubscribed = errors.New("already subscribed")
	ErrNotFound          = errors.New("subscriber not found")
)

type SubscriberService struct {
	repo repository.SubscriberRepository
}

func NewSubscriberService(repo repository.SubscriberRepository) *SubscriberService {
	return &SubscriberService{repo: repo}
}

func (s *SubscriberService) Subscribe(ctx context.Context, email string, pageID uint) (*entity.Subscriber, error) {
	existing, err := s.repo.FindByEmailAndPage(ctx, email, pageID)
	if err == nil {
		if existing.Verified {
			return nil, ErrAlreadySubscribed
		}
		// Resend verification
		return existing, nil
	}

	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return nil, err
	}
	sub := &entity.Subscriber{
		Email:        email,
		StatusPageID: pageID,
		Verified:     false,
		Token:        hex.EncodeToString(token),
	}
	err = s.repo.Create(ctx, sub)
	return sub, err
}

func (s *SubscriberService) Verify(ctx context.Context, token string) error {
	sub, err := s.repo.FindByToken(ctx, token)
	if err != nil {
		return ErrNotFound
	}
	sub.Verified = true
	return s.repo.Update(ctx, sub)
}

func (s *SubscriberService) Unsubscribe(ctx context.Context, email string, pageID uint) error {
	sub, err := s.repo.FindByEmailAndPage(ctx, email, pageID)
	if err != nil {
		return ErrNotFound
	}
	return s.repo.Delete(ctx, sub.ID)
}

func (s *SubscriberService) UnsubscribeByToken(ctx context.Context, token string) error {
	sub, err := s.repo.FindByToken(ctx, token)
	if err != nil {
		return ErrNotFound
	}
	return s.repo.Delete(ctx, sub.ID)
}

func (s *SubscriberService) GetSubscribers(ctx context.Context, pageID uint) ([]*entity.Subscriber, error) {
	return s.repo.FindByPageID(ctx, pageID)
}

func (s *SubscriberService) CountByPageID(ctx context.Context, pageID uint) (int64, error) {
	return s.repo.CountByPageID(ctx, pageID)
}
