package monitoring

import (
	"context"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/execution"
)

func TestStrategyMonitor_SignalHandling(t *testing.T) {
	monitor := NewStrategyMonitor(nil)

	signal := execution.Signal{
		StrategyID: "strat1",
	}
	monitor.OnSignalReceived(signal)

	status := monitor.GetStatus("strat1")
	if status.SignalsHandled != 1 {
		t.Errorf("expected 1 signal handled, got %d", status.SignalsHandled)
	}

	select {
	case event := <-monitor.Events():
		if event.EventType != "signal_received" {
			t.Errorf("unexpected event type: %s", event.EventType)
		}
	default:
		t.Error("expected event emitted")
	}
}

func TestStrategyMonitor_OrderHandling(t *testing.T) {
	monitor := NewStrategyMonitor(nil)

	order := execution.Order{
		StrategyID: "strat2",
	}
	monitor.OnOrderPlaced(order)

	status := monitor.GetStatus("strat2")
	if status.OrdersPlaced != 1 {
		t.Errorf("expected 1 order placed, got %d", status.OrdersPlaced)
	}

	select {
	case event := <-monitor.Events():
		if event.EventType != "order_placed" {
			t.Errorf("unexpected event type: %s", event.EventType)
		}
	default:
		t.Error("expected event emitted")
	}
}

func TestStrategyMonitor_StatusUpdate(t *testing.T) {
	monitor := NewStrategyMonitor(nil)

	monitor.OnStatusUpdate("strat3", "running", nil)
	status := monitor.GetStatus("strat3")
	if status.Status != "running" {
		t.Errorf("expected status 'running', got %s", status.Status)
	}

	select {
	case event := <-monitor.Events():
		if event.EventType != "status_update" {
			t.Errorf("unexpected event type: %s", event.EventType)
		}
	default:
		t.Error("expected event emitted")
	}

	// Simulate error update
	testErr := context.DeadlineExceeded
	monitor.OnStatusUpdate("strat3", "error", testErr)

	status = monitor.GetStatus("strat3")
	if status.LastError != testErr {
		t.Errorf("expected last error to be set")
	}

	select {
	case alert := <-monitor.Alerts():
		if alert.Level != "warning" {
			t.Errorf("expected warning alert, got %s", alert.Level)
		}
	default:
		t.Error("expected alert emitted")
	}
}

func TestStrategyMonitor_ListStatuses(t *testing.T) {
	monitor := NewStrategyMonitor(nil)

	monitor.OnSignalReceived(execution.Signal{StrategyID: "s1"})
	monitor.OnOrderPlaced(execution.Order{StrategyID: "s2"})
	monitor.OnStatusUpdate("s3", "running", nil)

	statuses := monitor.ListStatuses()
	if len(statuses) != 3 {
		t.Errorf("expected 3 statuses, got %d", len(statuses))
	}
}

func TestStrategyMonitor_Shutdown(t *testing.T) {
	monitor := NewStrategyMonitor(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := monitor.Shutdown(ctx)
	if err != nil {
		t.Errorf("unexpected shutdown error: %v", err)
	}
}
