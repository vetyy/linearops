package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"

	"github.com/vetyy/linear-opsgenie/internal/linear"
	"github.com/vetyy/linear-opsgenie/internal/opsgenie"
)

type Service struct {
	opsgenie *opsgenie.Service
	linear   *linear.Service
}

func NewService(opsGenieService *opsgenie.Service, linearService *linear.Service) (*Service, error) {
	return &Service{
		opsgenie: opsGenieService,
		linear:   linearService,
	}, nil
}

func (s *Service) LinearWebhook(w http.ResponseWriter, r *http.Request) {
	linearAction := &linear.Action{}
	err := json.NewDecoder(r.Body).Decode(&linearAction)
	if err != nil {
		handleInternalError(w, err, "")
		return
	}

	if s.linear.IsUnstarted(linearAction) && s.linear.IsUrgent(linearAction) {
		err = s.opsgenie.CreateAlert(r.Context(), linearAction.Data.ID, linearAction.Data.Title, linearAction.Data.Description, linearAction.URL)
	} else if s.linear.IsCompleted(linearAction) && s.linear.IsUrgent(linearAction) {
		err = s.opsgenie.CloseAlert(r.Context(), linearAction.Data.ID)
	} else {
		log.Infof("received linear webhook without action: %v", linearAction)
	}

	if err != nil {
		handleInternalError(w, err, "failed OpsGenie request")
		return
	}
	handleSuccess(w, nil)
}

type OnCallScheduleResponse struct {
	OnCall     []string `json:"on_call"`
	NextOnCall []string `json:"next_on_call"`
}

func (s *Service) OnCallSchedule(w http.ResponseWriter, r *http.Request) {
	scheduleID := chi.URLParam(r, "schedule_id")

	onCall, err := s.opsgenie.WhoIsOnCall(r.Context(), scheduleID)
	if err != nil {
		handleInternalError(w, err, "")
		return
	}

	nextOnCall, err := s.opsgenie.NextOnCall(r.Context(), scheduleID)
	if err != nil {
		handleInternalError(w, err, "")
		return
	}
	handleSuccess(w, OnCallScheduleResponse{OnCall: onCall, NextOnCall: nextOnCall})
}
