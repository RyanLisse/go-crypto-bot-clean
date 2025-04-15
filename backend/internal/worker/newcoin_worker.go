package worker

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
)

// NewCoinWorker periodically checks for new coin listings
type NewCoinWorker struct {
	newCoinUC *usecase.NewCoinUseCase
	cfg       *config.Config
	logger    zerolog.Logger
	stopCh    chan struct{}
	stopped   bool
}

// NewNewCoinWorker creates a new NewCoinWorker instance
func NewNewCoinWorker(newCoinUC *usecase.NewCoinUseCase, cfg *config.Config, logger zerolog.Logger) *NewCoinWorker {
	return &NewCoinWorker{
		newCoinUC: newCoinUC,
		cfg:       cfg,
		logger:    logger.With().Str("component", "NewCoinWorker").Logger(),
		stopCh:    make(chan struct{}),
	}
}

// Start begins the worker's periodic execution
func (w *NewCoinWorker) Start(ctx context.Context) {
	w.logger.Info().Msg("Starting NewCoinWorker")
	// TODO: Add worker polling interval configuration to config.go and config file(s)
	// Using hardcoded default for now.
	interval := 5 * time.Minute
	w.logger.Warn().Msgf("Worker polling interval configuration missing, using default: %v. Please add 'workers.new_coin_detection.polling_interval_sec' to your config.", interval)
	// Example of how it would look if config existed:
	// if w.cfg != nil && w.cfg.Workers != nil && w.cfg.Workers.NewCoinDetection.PollingIntervalSec > 0 {
	// 	interval = time.Duration(w.cfg.Workers.NewCoinDetection.PollingIntervalSec) * time.Second
	// } else {
	// 	w.logger.Warn().Msgf("Invalid or missing NewCoinDetection polling interval in config, using default: %v", interval)
	// }

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run once immediately on start
	w.runDetection(ctx)

	for {
		select {
		case <-ticker.C:
			w.runDetection(ctx)
		case <-w.stopCh:
			w.logger.Info().Msg("Stopping NewCoinWorker")
			w.stopped = true
			return
		case <-ctx.Done():
			w.logger.Info().Msg("Context cancelled, stopping NewCoinWorker")
			w.stopped = true
			return
		}
	}
}

// Stop signals the worker to stop
func (w *NewCoinWorker) Stop() {
	if !w.stopped {
		close(w.stopCh)
	}
}

// runDetection executes a single detection cycle
func (w *NewCoinWorker) runDetection(ctx context.Context) {
	w.logger.Info().Msg("Running new coin detection cycle")
	if err := w.newCoinUC.DetectNewCoins(); err != nil {
		w.logger.Error().Err(err).Msg("Error during new coin detection")
	} else {
		w.logger.Info().Msg("New coin detection cycle completed successfully")
	}
}
