// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package freja

import (
	"github.com/diego1q2w/freja/healthcheck"
	"sync"
)

var (
	lockhealthCalculatorMockAdd       sync.RWMutex
	lockhealthCalculatorMockCalculate sync.RWMutex
)

// Ensure, that healthCalculatorMock does implement healthCalculator.
// If this is not the case, regenerate this file with moq.
var _ healthCalculator = &healthCalculatorMock{}

// healthCalculatorMock is a mock implementation of healthCalculator.
//
//     func TestSomethingThatUseshealthCalculator(t *testing.T) {
//
//         // make and configure a mocked healthCalculator
//         mockedhealthCalculator := &healthCalculatorMock{
//             AddFunc: func(in1 healthcheck.HealthChecker)  {
// 	               panic("mock out the Add method")
//             },
//             CalculateFunc: func() (bool, []Status) {
// 	               panic("mock out the Calculate method")
//             },
//         }
//
//         // use mockedhealthCalculator in code that requires healthCalculator
//         // and then make assertions.
//
//     }
type healthCalculatorMock struct {
	// AddFunc mocks the Add method.
	AddFunc func(in1 healthcheck.HealthChecker)

	// CalculateFunc mocks the Calculate method.
	CalculateFunc func() (bool, []Status)

	// calls tracks calls to the methods.
	calls struct {
		// Add holds details about calls to the Add method.
		Add []struct {
			// In1 is the in1 argument value.
			In1 healthcheck.HealthChecker
		}
		// Calculate holds details about calls to the Calculate method.
		Calculate []struct {
		}
	}
}

// Add calls AddFunc.
func (mock *healthCalculatorMock) Add(in1 healthcheck.HealthChecker) {
	if mock.AddFunc == nil {
		panic("healthCalculatorMock.AddFunc: method is nil but healthCalculator.Add was just called")
	}
	callInfo := struct {
		In1 healthcheck.HealthChecker
	}{
		In1: in1,
	}
	lockhealthCalculatorMockAdd.Lock()
	mock.calls.Add = append(mock.calls.Add, callInfo)
	lockhealthCalculatorMockAdd.Unlock()
	mock.AddFunc(in1)
}

// AddCalls gets all the calls that were made to Add.
// Check the length with:
//     len(mockedhealthCalculator.AddCalls())
func (mock *healthCalculatorMock) AddCalls() []struct {
	In1 healthcheck.HealthChecker
} {
	var calls []struct {
		In1 healthcheck.HealthChecker
	}
	lockhealthCalculatorMockAdd.RLock()
	calls = mock.calls.Add
	lockhealthCalculatorMockAdd.RUnlock()
	return calls
}

// Calculate calls CalculateFunc.
func (mock *healthCalculatorMock) Calculate() (bool, []Status) {
	if mock.CalculateFunc == nil {
		panic("healthCalculatorMock.CalculateFunc: method is nil but healthCalculator.Calculate was just called")
	}
	callInfo := struct {
	}{}
	lockhealthCalculatorMockCalculate.Lock()
	mock.calls.Calculate = append(mock.calls.Calculate, callInfo)
	lockhealthCalculatorMockCalculate.Unlock()
	return mock.CalculateFunc()
}

// CalculateCalls gets all the calls that were made to Calculate.
// Check the length with:
//     len(mockedhealthCalculator.CalculateCalls())
func (mock *healthCalculatorMock) CalculateCalls() []struct {
} {
	var calls []struct {
	}
	lockhealthCalculatorMockCalculate.RLock()
	calls = mock.calls.Calculate
	lockhealthCalculatorMockCalculate.RUnlock()
	return calls
}
