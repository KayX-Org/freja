package middleware

import (
	"context"
	"fmt"
	"github.com/kayx-org/freja/healthcheck"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewDB(t *testing.T) {
	testCases := map[string]struct {
		pingErr        error
		hcName         string
		expectedName   string
		expectedStatus healthcheck.ServiceStatus
	}{
		"if the pinging process return error then it should have the correct status": {
			pingErr:        fmt.Errorf("test"),
			hcName:         "foo",
			expectedName:   "foo",
			expectedStatus: healthcheck.DOWN,
		},
		"if the pinging process is fine then it should have the correct status": {
			pingErr:        nil,
			hcName:         "foo",
			expectedName:   "foo",
			expectedStatus: healthcheck.UP,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db := &dbMock{
				PingContextFunc: func(context.Context) error {
					return tc.pingErr
				},
			}
			midd := NewDB(db, OptionHealthCheckName(tc.hcName), OptionWindowCheck(time.Millisecond))
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*3)
			defer cancel()

			err := midd.Run(ctx)
			time.Sleep(time.Millisecond * 6)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, midd.status)
			assert.Equal(t, tc.expectedName, midd.name)
		})
	}
}
