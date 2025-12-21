package features

import (
	"os"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
)

// Flags holds all feature flags for the application
type Flags struct {
	autoApproval bool
	mu           sync.RWMutex
	logger       *logrus.Logger
}

var flags *Flags

// Initialize sets up feature flags
// To integrate with CloudBees Feature Management:
// See the integration guide at the bottom of this file
func Initialize(apiKey string, logger *logrus.Logger) (*Flags, error) {
	flags = &Flags{
		logger: logger,
	}

	// Load feature flags from environment variables
	// api.advancedFilters (default: false) - enable complex filtering
	advancedFiltersStr := os.Getenv("FEATURE_ADVANCED_FILTERS")
	if advancedFiltersStr != "" {
		advancedFilters, err := strconv.ParseBool(advancedFiltersStr)
		if err == nil {
			flags.advancedFilters = advancedFilters
		}
	}

	logger.WithFields(logrus.Fields{
		"advancedFilters": flags.advancedFilters,
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

// IsAdvancedFiltersEnabled returns whether advanced filters are enabled
func (f *Flags) IsAdvancedFiltersEnabled() bool {
	if f == nil {
		return false
	}
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.advancedFilters
}

// SetAdvancedFilters sets the advanced filters flag (for testing/admin purposes)
func (f *Flags) SetAdvancedFilters(enabled bool) {
	if f == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.advancedFilters = enabled
	f.logger.WithField("advancedFilters", enabled).Info("Feature flag updated")
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
   go get github.com/rollout/rox-go/core

2. Update imports:
   import (
       "github.com/rollout/rox-go/core/model"
       "github.com/rollout/rox-go/core/roxx"
   )

3. Replace the Flags struct:
   type Flags struct {
       AdvancedFilters model.RoxFlag
       logger          *logrus.Logger
   }

4. Update Initialize function:
   func Initialize(apiKey string, logger *logrus.Logger) (*Flags, error) {
       flags = &Flags{
           logger: logger,
       }

       // Register feature flag: api.advancedFilters (default: false)
       flags.AdvancedFilters = model.NewRoxFlag(false)

       // Register with CloudBees
       roxx.Register("api", flags)

       // Setup Rox with API key
       options := roxx.NewRoxOptions(roxx.RoxOptionsBuilder{})
       <-roxx.Setup(apiKey, options)

       logger.Info("CloudBees Feature Management initialized successfully")

       // Fetch latest feature flags
       go func() {
           roxx.Fetch()
           logger.Info("Initial feature flags fetched")
       }()

       return flags, nil
   }

5. Update IsAdvancedFiltersEnabled:
   func (f *Flags) IsAdvancedFiltersEnabled() bool {
       if f == nil || f.AdvancedFilters == nil {
           return false
       }
       return f.AdvancedFilters.IsEnabled(nil)
   }

6. Update Shutdown:
   func Shutdown() {
       if flags != nil {
           roxx.Shutdown()
           flags.logger.Info("CloudBees Feature Management shutdown complete")
       }
   }

For more information, see: https://docs.cloudbees.com/docs/cloudbees-feature-management/latest/
*/
