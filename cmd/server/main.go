package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/vetyy/linear-opsgenie/internal/config"
	"github.com/vetyy/linear-opsgenie/internal/linear"
	"github.com/vetyy/linear-opsgenie/internal/opsgenie"
	"github.com/vetyy/linear-opsgenie/internal/server"
)

const (
	configLogLevel = "log.Level"

	configBindHTTP      = "bind.HTTP"
	configBindURLPrefix = "bind.URL.Prefix"

	configOpsGenieAPIKey     = "opsGenie.API.Key"
	configOpsGenieResponders = "opsGenie.Responders"

	configLinearWebhookID       = "linear.Webhook.ID"
	configLinearUserAgent       = "linear.UserAgent"
	configLinearUnstartedStates = "linear.Unstarted.States"
	configLinearCompletedStates = "linear.Completed.States"
)

var (
	// Version information related variables set by build command.
	version         string
	commitID        string
	commitTimestamp string
)

func init() {
	viper.SetDefault(configLogLevel, "info")
	viper.SetDefault(configBindURLPrefix, "/")
	viper.SetDefault(configBindHTTP, ":8080")
	viper.SetDefault(configLinearUserAgent, "Linear")
	viper.SetDefault(configLinearUnstartedStates, "Reported")
	viper.SetDefault(configLinearCompletedStates, "Resolved,Postmortem,Rejected")
}

func main() {
	err := config.SetupViper()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	err = config.SetupLogrus(viper.GetString(configLogLevel))
	if err != nil {
		log.Fatalf("failed to setup logrus: %v", err)
	}

	log.Infoln("Version:", version)
	log.Infoln("Commit ID:", commitID)
	log.Infoln("Commit timestamp:", commitTimestamp)

	requiredConfigVariables := []string{configLinearWebhookID, configOpsGenieAPIKey, configOpsGenieResponders}
	for _, v := range requiredConfigVariables {
		if !viper.IsSet(v) {
			log.Fatalf("required config variable '%s' not found", v)
		}
	}

	opsGenieService, err := opsgenie.NewService(
		viper.GetString(configOpsGenieAPIKey),
		viper.GetString(configOpsGenieResponders),
		viper.GetString(configLinearUserAgent),
	)
	if err != nil {
		log.Fatalf("failed to create OpsGenie service: %v", err)
	}

	linearCompletedStates := strings.Split(viper.GetString(configLinearCompletedStates), ",")
	linearUnstartedStates := strings.Split(viper.GetString(configLinearUnstartedStates), ",")
	linearService := linear.NewService(linearUnstartedStates, linearCompletedStates)
	serverService, err := server.NewService(opsGenieService, linearService)
	if err != nil {
		log.Fatalf("failed to create Server service: %v", err)
	}

	r := chi.NewRouter()
	if log.GetLevel() == log.DebugLevel {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Heartbeat("/health"))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Route(viper.GetString(configBindURLPrefix), func(r chi.Router) {
		r.Post(fmt.Sprintf("/webhook/%s", viper.GetString(configLinearWebhookID)), serverService.LinearWebhook)
		r.Get("/on-call/{schedule_id}", serverService.OnCallSchedule)
	})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		bindHTTP := viper.GetString(configBindHTTP)
		log.Infof("Running server at: %s", bindHTTP)
		if err = http.ListenAndServe(bindHTTP, r); err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
		wg.Done()
	}()
	wg.Wait()
}
