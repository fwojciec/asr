// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/fwojciec/asr"
	"io"
	"sync"
)

// Ensure, that GetterMock does implement asr.Getter.
// If this is not the case, regenerate this file with moq.
var _ asr.Getter = &GetterMock{}

// GetterMock is a mock implementation of asr.Getter.
//
// 	func TestSomethingThatUsesGetter(t *testing.T) {
//
// 		// make and configure a mocked asr.Getter
// 		mockedGetter := &GetterMock{
// 			GetFunc: func(ctx context.Context, url string) (io.ReadCloser, error) {
// 				panic("mock out the Get method")
// 			},
// 		}
//
// 		// use mockedGetter in code that requires asr.Getter
// 		// and then make assertions.
//
// 	}
type GetterMock struct {
	// GetFunc mocks the Get method.
	GetFunc func(ctx context.Context, url string) (io.ReadCloser, error)

	// calls tracks calls to the methods.
	calls struct {
		// Get holds details about calls to the Get method.
		Get []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// URL is the url argument value.
			URL string
		}
	}
	lockGet sync.RWMutex
}

// Get calls GetFunc.
func (mock *GetterMock) Get(ctx context.Context, url string) (io.ReadCloser, error) {
	if mock.GetFunc == nil {
		panic("GetterMock.GetFunc: method is nil but Getter.Get was just called")
	}
	callInfo := struct {
		Ctx context.Context
		URL string
	}{
		Ctx: ctx,
		URL: url,
	}
	mock.lockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	mock.lockGet.Unlock()
	return mock.GetFunc(ctx, url)
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//     len(mockedGetter.GetCalls())
func (mock *GetterMock) GetCalls() []struct {
	Ctx context.Context
	URL string
} {
	var calls []struct {
		Ctx context.Context
		URL string
	}
	mock.lockGet.RLock()
	calls = mock.calls.Get
	mock.lockGet.RUnlock()
	return calls
}