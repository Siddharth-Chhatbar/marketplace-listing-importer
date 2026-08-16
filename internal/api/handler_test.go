package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIsServiceLive(t *testing.T) {
	checker := fakeDatabaseChecker{err: errors.New("unavailable")}
	handler := NewHandler(checker)
	request := httptest.NewRequest(http.MethodGet, "/livez", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	expectedStatus := http.StatusNoContent
	if response.Code != expectedStatus {
		t.Errorf("got status %v, want status %v", response.Code, expectedStatus)
	}
}

type fakeDatabaseChecker struct {
	err              error
	waitForContext   bool
	requireDeadline  bool
	maximumTimeLimit time.Duration
}

func (f fakeDatabaseChecker) PingContext(ctx context.Context) error {
	if f.requireDeadline {
		deadline, ok := ctx.Deadline()
		if !ok {
			return errors.New("context has no deadline")
		}
		if remaining := time.Until(deadline); remaining <= 0 || remaining > f.maximumTimeLimit {
			return fmt.Errorf("context deadline is outside expected limit: %s", remaining)
		}
	}
	if f.waitForContext {
		<-ctx.Done()
		return ctx.Err()
	}
	return f.err
}

func TestIsServiceReady(t *testing.T) {
	tests := []struct {
		name       string
		checkErr   error
		cancel     bool
		waitForCtx bool
		checkLimit bool
		wantStatus int
	}{
		{
			name:       "dependency available",
			checkErr:   nil,
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "dependency unavailable",
			checkErr:   errors.New("unavailable"),
			wantStatus: http.StatusServiceUnavailable,
		},
		{
			name:       "request cancelled",
			cancel:     true,
			waitForCtx: true,
			wantStatus: http.StatusServiceUnavailable,
		},
		{
			name:       "readiness limits database check duration",
			checkLimit: true,
			wantStatus: http.StatusNoContent,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			checker := fakeDatabaseChecker{
				err:              test.checkErr,
				waitForContext:   test.waitForCtx,
				requireDeadline:  test.checkLimit,
				maximumTimeLimit: 2 * time.Second,
			}
			handler := NewHandler(checker)

			request := httptest.NewRequest(http.MethodGet, "/readyz", nil)
			if test.cancel {
				ctx, cancel := context.WithCancel(request.Context())
				cancel()
				request = request.WithContext(ctx)
			}
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != test.wantStatus {
				t.Errorf("got status %v, want status %v", response.Code, test.wantStatus)
			}
		})
	}
}
