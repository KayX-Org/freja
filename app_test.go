package freja

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kayx-org/freja/healthcheck"
	"github.com/stretchr/testify/assert"
	"syscall"
	"testing"
	"time"
)

func TestAppMiddleware(t *testing.T) {
	testCases := map[string]struct {
		initErr                   error
		runErr                    error
		stopErr                   error
		server                    *ServerMock
		expectedInitCalls         int
		expectedRunCalls          int
		expectedStopCalls         int
		expectedFatalLogCalls     int
		expectedServerListenCalls int
		expectedServerStopCalls   int
		expectedErr               error
	}{
		"it should run correctly without server": {
			expectedInitCalls: 1,
			expectedRunCalls:  1,
			expectedStopCalls: 1,
		},
		"it should run correctly with server": {
			expectedInitCalls:         1,
			expectedRunCalls:          1,
			expectedStopCalls:         1,
			expectedServerListenCalls: 1,
			expectedServerStopCalls:   1,
			server: &ServerMock{
				ListenAndServeFunc: func() error {
					return nil
				},
				ShutdownFunc: func(context.Context) error {
					return nil
				},
			},
		},
		"if there is an error during init an error is expected": {
			initErr:           fmt.Errorf("test"),
			expectedInitCalls: 1,
			expectedErr:       fmt.Errorf("unable to run Init(): test"),
		},
		"if the server is unable to start the fatal log should be called": {
			expectedInitCalls:         1,
			expectedServerListenCalls: 1,
			expectedServerStopCalls:   1,
			expectedStopCalls:         1,
			expectedRunCalls:          1,
			expectedFatalLogCalls:     1,
			server: &ServerMock{
				ListenAndServeFunc: func() error {
					return fmt.Errorf("test")
				},
				ShutdownFunc: func(context.Context) error {
					return nil
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			logger := &LoggerMock{
				FatalfFunc: func(string, ...interface{}) {},
				InfoFunc:   func(...interface{}) {},
				InfofFunc:  func(string, ...interface{}) {},
			}
			app := NewApp(NewHealthCalculator(), logger)
			mid := &MiddlewareMock{
				InitFunc: func() error {
					return tc.initErr
				},
				RunFunc: func(context.Context) error {
					return tc.runErr
				},
				StopFunc: func(context.Context) error {
					return tc.stopErr
				},
			}

			if tc.server != nil {
				app.WithServer(tc.server)
			}
			app.AddMiddleware(mid)

			go func() {
				time.Sleep(time.Millisecond * 300)
				app.osSignal <- syscall.SIGTERM
			}()
			err := app.Start(context.Background())
			if fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tc.expectedErr) {
				t.Errorf("expected error %s, got %s", tc.expectedErr, err)
			}

			assert.Equal(t, len(mid.calls.Init), tc.expectedInitCalls, "init called once")
			assert.Equal(t, len(mid.calls.Run), tc.expectedRunCalls, "run called once")
			assert.Equal(t, len(mid.calls.Stop), tc.expectedStopCalls, "stop called once")

			assert.Equal(t, len(logger.calls.Fatalf), tc.expectedFatalLogCalls, "fatal log called once")

			if tc.server != nil {
				assert.Equal(t, len(tc.server.calls.ListenAndServe), tc.expectedServerListenCalls, "ListenAndServe called once")
				assert.Equal(t, len(tc.server.calls.Shutdown), tc.expectedServerStopCalls, "Shutdown called once")
			}
		})
	}
}

func TestAppHealthCheck(t *testing.T) {
	testCases := map[string]struct {
		middleware             Middleware
		healthCheck            healthcheck.HealthChecker
		healthCalculator       bool
		expectedAddCalls       int
		expectedCalculateCalls int
		expectedStatus         bool
		expectedSummary        string
	}{
		"if there is no health calculator to process it should return true": {
			healthCalculator:       false,
			expectedAddCalls:       0,
			expectedCalculateCalls: 0,
			expectedSummary:        "",
			expectedStatus:         true,
		},
		"it should process the result of the calculator": {
			//middleware: mockHCWithMiddleware{name: "foo"},
			healthCalculator:       true,
			expectedAddCalls:       0,
			expectedCalculateCalls: 1,
			expectedSummary: `[{"name":"foo","status":"up"},{"name":"bar","status":"down"}]
`,
			expectedStatus: false,
		},
		"it should add the middleware if implements the HealthChecker interface": {
			middleware:             mockHCWithMiddleware{name: "foo"},
			healthCalculator:       true,
			expectedAddCalls:       1,
			expectedCalculateCalls: 1,
			expectedSummary: `[{"name":"foo","status":"up"},{"name":"bar","status":"down"}]
`,
			expectedStatus: false,
		},
		"it should not add the middleware even if implements the HealthChecker interface if the calculator is not there": {
			middleware:             mockHCWithMiddleware{name: "foo"},
			healthCalculator:       false,
			expectedAddCalls:       0,
			expectedCalculateCalls: 0,
			expectedSummary:        "",
			expectedStatus:         true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			healthCalculator := &healthCalculatorMock{
				AddFunc: func(healthcheck.HealthChecker) {},
				CalculateFunc: func() (bool, []Status) {
					return false, []Status{{Name: "foo", Status: "up"}, {Name: "bar", Status: "down"}}
				}}

			var app *ServiceApp
			if tc.healthCalculator {
				app = NewApp(healthCalculator, &DummyLogger{})
			} else {
				app = NewApp(nil, &DummyLogger{})
			}

			if tc.middleware != nil {
				app.AddMiddleware(tc.middleware)
			}
			if tc.healthCheck != nil {
				app.AddHealthCheck(tc.healthCheck)
			}

			var buff bytes.Buffer
			status, err := app.HealthCheck(&buff)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, status)
			assert.Equal(t, tc.expectedSummary, buff.String())
			assert.Equal(t, tc.expectedAddCalls, len(healthCalculator.calls.Add))
			assert.Equal(t, tc.expectedCalculateCalls, len(healthCalculator.calls.Calculate))
		})
	}
}

type mockHCWithMiddleware struct {
	name string
}

func (m mockHCWithMiddleware) Init() error {
	return nil
}
func (m mockHCWithMiddleware) Run(context.Context) error {
	return nil
}
func (m mockHCWithMiddleware) Stop(context.Context) error {
	return nil
}
func (m mockHCWithMiddleware) Name() string {
	return m.name
}
func (m mockHCWithMiddleware) Status() healthcheck.ServiceStatus {
	return healthcheck.UP
}
