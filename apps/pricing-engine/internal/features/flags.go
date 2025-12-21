package features

import (
	"os"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
)

// Flags holds all feature flags for the application
type Flags struct {
	insightsV2    bool
	alertsEnabled bool
	mu            sync.RWMutex
	logger        *logrus.Logger
}

var flags *Flags

// Initialize sets up feature flags
// To integrate with CloudBees Feature Management:
// 1. Install: go get github.com/rollout/rox-go/v5/core
// 2. Import the SDK
// 3. Replace this implementation with CloudBees Rox SDK initialization
func Initialize(apiKey string, logger *logrus.Logger) (*Flags, error) {
	flags = &Flags{
		logger: logger,
	}

	// Load feature flags from environment variables
	// api.insightsV2 (default: false) - use new insights calculation algorithm
	insightsV2Str := os.Getenv("FEATURE_INSIGHTS_V2")
	if insightsV2Str != "" {
		insightsV2, err := strconv.ParseBool(insightsV2Str)
		if err == nil {
			flags.insightsV2 = insightsV2
		}
	}

	// api.alertsEnabled (default: true) - enable alert generation
	flags.alertsEnabled = true // Default to enabled
	alertsEnabledStr := os.Getenv("FEATURE_ALERTS_ENABLED")
	if alertsEnabledStr != "" {
		alertsEnabled, err := strconv.ParseBool(alertsEnabledStr)
		if err == nil {
			flags.alertsEnabled = alertsEnabled
		}
	}

	logger.WithFields(logrus.Fields{
		"insightsV2":    flags.insightsV2,
		"alertsEnabled": flags.alertsEnabled,
	}).Info("Feature flags initialized")

	if apiKey != "" && apiKey != "dev-mode" {
		logger.Warn("CloudBees Feature Management API key provided but SDK not integrated. See flags.go for integration instructions.")
	}

	return flags, nil
}

// GetFlags returns the global flags instance
func GetFlags() *Flags {
	return flags
}

// IsInsightsV2Enabled returns whether the V2 insights algorithm should be used
func (f *Flags) IsInsightsV2Enabled() bool {
	if f == nil {
		return false
	}
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.insightsV2
}

// IsAlertsEnabled returns whether alerts are enabled
func (f *Flags) IsAlertsEnabled() bool {
	if f == nil {
		return true // Default to enabled
	}
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.alertsEnabled
}

// SetInsightsV2 sets the insights V2 flag (for testing/admin purposes)
func (f *Flags) SetInsightsV2(enabled bool) {
	if f == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.insightsV2 = enabled
	f.logger.WithField("insightsV2", enabled).Info("Feature flag updated")
}

// SetAlertsEnabled sets the alerts enabled flag (for testing/admin purposes)
func (f *Flags) SetAlertsEnabled(enabled bool) {
	if f == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.alertsEnabled = enabled
	f.logger.WithField("alertsEnabled", enabled).Info("Feature flag updated")
}

// Shutdown gracefully shuts down the feature management system
func Shutdown() {
	if flags != nil {
		flags.logger.Info("Feature management shutdown complete")
	}
}

/*
CloudBees Feature Management Integration Guide:

To integrate with CloudBees Feature Management (Rox SDK), follow these steps:

1. Install the CloudBees Rox SDK:
   go get github.com/rollout/rox-go/v5/core

2. Update imports:
   import (
       "github.com/rollout/rox-go/v5/core"
   )

3. Replace the Flags struct:
   type Flags struct {
       InsightsV2    *core.RoxFlag
       AlertsEnabled *core.RoxFlag
       logger        *logrus.Logger
   }

4. Update Initialize function:
   func Initialize(apiKey string, logger *logrus.Logger) (*Flags, error) {
       flags = &Flags{
           logger: logger,
       }

       // Register feature flag: api.insightsV2 (default: false)
       flags.InsightsV2 = core.NewRoxFlag(false)

       // Register feature flag: api.alertsEnabled (default: true)
       flags.AlertsEnabled = core.NewRoxFlag(true)

       // Register with CloudBees
       core.Register("api", flags)

       // Setup Rox with API key
       options := core.NewRoxOptions(core.RoxOptionsBuilder{})
       <-core.Setup(apiKey, options)

       logger.Info("CloudBees Feature Management initialized successfully")

       // Fetch latest feature flags
       go func() {
           core.Fetch()
           logger.Info("Initial feature flags fetched")
       }()

       return flags, nil
   }

5. Update IsInsightsV2Enabled:
   func (f *Flags) IsInsightsV2Enabled() bool {
       if f == nil || f.InsightsV2 == nil {
           return false
       }
       return f.InsightsV2.IsEnabled(nil)
   }

6. Update IsAlertsEnabled:
   func (f *Flags) IsAlertsEnabled() bool {
       if f == nil || f.AlertsEnabled == nil {
           return true
       }
       return f.AlertsEnabled.IsEnabled(nil)
   }

7. Update Shutdown:
   func Shutdown() {
       if flags != nil {
           core.Shutdown()
           flags.logger.Info("CloudBees Feature Management shutdown complete")
       }
   }

Feature Flags:
- api.insightsV2 (default: false) - use new insights calculation algorithm
- api.alertsEnabled (default: true) - enable alert generation

Environment Variables (Current Implementation):
- FEATURE_INSIGHTS_V2: Set to "true" to enable V2 algorithm
- FEATURE_ALERTS_ENABLED: Set to "false" to disable alerts

For more information, see: https://docs.cloudbees.com/docs/cloudbees-feature-management/latest/
*/
