package service

import (
	"context"
	"errors"
	"ping-uptime/modules/incidents/domain/entity"
	"ping-uptime/modules/incidents/domain/repository"
	"time"
)

var (
	ErrIncidentNotFound = errors.New("incident not found")
)

type IncidentService struct {
	incidentRepo repository.IncidentRepository
}

func NewIncidentService(incidentRepo repository.IncidentRepository) *IncidentService {
	return &IncidentService{
		incidentRepo: incidentRepo,
	}
}

func (s *IncidentService) CreateIncident(ctx context.Context, incident *entity.Incident) error {
	if incident.MonitorID == 0 {
		return errors.New("monitor_id cannot be 0")
	}
	return s.incidentRepo.Create(ctx, incident)
}

func (s *IncidentService) DeleteIncident(ctx context.Context, id uint) error {
	return s.incidentRepo.Delete(ctx, id)
}

func (s *IncidentService) GetAllIncidents(ctx context.Context) ([]*entity.Incident, error) {
	return s.incidentRepo.FindAll(ctx)
}

func (s *IncidentService) GetIncidentsByUserID(ctx context.Context, userID uint) ([]*entity.Incident, error) {
	return s.incidentRepo.FindByUserID(ctx, userID)
}

func (s *IncidentService) GetIncidentsByMonitorID(ctx context.Context, monitorID uint) ([]*entity.Incident, error) {
	return s.incidentRepo.FindByMonitorID(ctx, monitorID)
}

func (s *IncidentService) GetIncidentByID(ctx context.Context, id uint) (*entity.Incident, error) {
	return s.incidentRepo.FindByID(ctx, id)
}

func (s *IncidentService) UpdateIncident(ctx context.Context, incident *entity.Incident) error {
	return s.incidentRepo.Update(ctx, incident)
}

func (s *IncidentService) ResolveIncident(ctx context.Context, id uint) error {
	incident, err := s.incidentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if incident.Status != "active" {
		return errors.New("incident is not active")
	}
	now := time.Now()
	resolved := &entity.Incident{
		MonitorID:    incident.MonitorID,
		UserID:       incident.UserID,
		Status:       "resolved",
		ErrorMessage: incident.ErrorMessage,
		Latency:      incident.Latency,
		CreatedAt:    now,
		ResolvedAt:   &now,
	}
	return s.incidentRepo.Create(ctx, resolved)
}
