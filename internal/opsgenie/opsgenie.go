package opsgenie

import (
	"context"
	"fmt"
	"time"

	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

const (
	NextOnCallPeopleCacheKey = "nextOnCallPeople"
	OnCallPeopleCacheKey     = "onCallPeople"
)

type Service struct {
	Responders      string
	LinearUserAgent string

	alertClient    *alert.Client
	scheduleClient *schedule.Client
	cache          *cache.Cache
}

func NewService(apiKey string, responders string, linearUserAgent string) (*Service, error) {
	alertClient, err := alert.NewClient(&client.Config{
		Logger: log.StandardLogger(),
		ApiKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create OpsGenie service: %v", err)
	}

	scheduleClient, err := schedule.NewClient(&client.Config{
		Logger: log.StandardLogger(),
		ApiKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create OpsGenie service: %v", err)
	}

	return &Service{
		LinearUserAgent: linearUserAgent,
		Responders:      responders,
		alertClient:     alertClient,
		scheduleClient:  scheduleClient,
		cache:           cache.New(10*time.Minute, 20*time.Minute),
	}, nil
}

func (s *Service) CreateAlert(ctx context.Context, id string, title string, description string, url string) error {
	_, err := s.alertClient.Create(ctx, &alert.CreateAlertRequest{
		Message:     title,
		Description: description,
		Responders: []alert.Responder{
			{
				Type: alert.TeamResponder,
				Name: s.Responders,
			},
		},
		Details: map[string]string{
			"url": url,
		},
		Alias:    id,
		Priority: alert.P1,
		User:     s.LinearUserAgent,
	})
	return err
}

func (s *Service) CloseAlert(ctx context.Context, id string) error {
	_, err := s.alertClient.Close(ctx, &alert.CloseAlertRequest{
		IdentifierType:  alert.ALIAS,
		IdentifierValue: id,
		User:            s.LinearUserAgent,
	})
	return err
}

func (s *Service) WhoIsOnCall(ctx context.Context, scheduleID string) (onCallPeople []string, err error) {
	onCallPeopleCached, found := s.cache.Get(OnCallPeopleCacheKey)
	if found {
		return onCallPeopleCached.([]string), nil
	}

	flat := true
	scheduleResult, err := s.scheduleClient.GetOnCalls(ctx, &schedule.GetOnCallsRequest{
		Flat:                   &flat,
		ScheduleIdentifierType: schedule.Name,
		ScheduleIdentifier:     scheduleID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get on-calls: %v", err)
	}

	onCallPeople = scheduleResult.OnCallRecipients
	s.cache.Set(OnCallPeopleCacheKey, onCallPeople, cache.DefaultExpiration)
	return onCallPeople, nil
}

func (s *Service) NextOnCall(ctx context.Context, scheduleID string) (nextOnCallPeople []string, err error) {
	nextOnCallPeopleCached, found := s.cache.Get(NextOnCallPeopleCacheKey)
	if found {
		return nextOnCallPeopleCached.([]string), nil
	}

	flat := true
	scheduleResult, err := s.scheduleClient.GetNextOnCall(ctx, &schedule.GetNextOnCallsRequest{
		Flat:                   &flat,
		ScheduleIdentifierType: schedule.Name,
		ScheduleIdentifier:     scheduleID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get next on-calls: %v", err)
	}

	nextOnCallPeople = scheduleResult.NextOncallParticipants
	s.cache.Set(NextOnCallPeopleCacheKey, nextOnCallPeople, cache.DefaultExpiration)
	return nextOnCallPeople, err
}
