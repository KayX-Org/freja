package freya

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"syscall"
	"testing"
	"time"
)

func TestApp(t *testing.T) {
	testCases := map[string]struct {
		initErr                   error
		runErr                    error
		stopErr                   error
		server                    *ServerMock
		expectedInitCalls         int
		expectedRunCalls          int
		expectedStopCalls         int
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
		"if the server is unable to start an error is expected": {
			expectedInitCalls:         1,
			expectedServerListenCalls: 1,
			expectedRunCalls:          1,
			server: &ServerMock{
				ListenAndServeFunc: func() error {
					return fmt.Errorf("test")
				},
				ShutdownFunc: func(context.Context) error {
					return nil
				},
			},
			expectedErr: fmt.Errorf("unable to run the server: test"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			app := App()
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
				app.AddServer(tc.server)
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

			if tc.server != nil {
				assert.Equal(t, len(tc.server.calls.ListenAndServe), tc.expectedServerListenCalls, "ListenAndServe called once")
				assert.Equal(t, len(tc.server.calls.Shutdown), tc.expectedServerStopCalls, "Shutdown called once")
			}
		})
	}
}
