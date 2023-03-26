package mocks

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

//go:generate minimock -i route256/checkout/internal/clients/productservice.Client -o ./mocks/client_minimock.go -n ClientMock

import (
	"context"
	"route256/checkout/internal/model"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// ClientMock implements productservice.Client
type ClientMock struct {
	t minimock.Tester

	funcGetProduct          func(ctx context.Context, SKU uint32) (pp1 *model.Product, err error)
	inspectFuncGetProduct   func(ctx context.Context, SKU uint32)
	afterGetProductCounter  uint64
	beforeGetProductCounter uint64
	GetProductMock          mClientMockGetProduct
}

// NewClientMock returns a mock for productservice.Client
func NewClientMock(t minimock.Tester) *ClientMock {
	m := &ClientMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetProductMock = mClientMockGetProduct{mock: m}
	m.GetProductMock.callArgs = []*ClientMockGetProductParams{}

	return m
}

type mClientMockGetProduct struct {
	mock               *ClientMock
	defaultExpectation *ClientMockGetProductExpectation
	expectations       []*ClientMockGetProductExpectation

	callArgs []*ClientMockGetProductParams
	mutex    sync.RWMutex
}

// ClientMockGetProductExpectation specifies expectation struct of the Client.GetProduct
type ClientMockGetProductExpectation struct {
	mock    *ClientMock
	params  *ClientMockGetProductParams
	results *ClientMockGetProductResults
	Counter uint64
}

// ClientMockGetProductParams contains parameters of the Client.GetProduct
type ClientMockGetProductParams struct {
	ctx context.Context
	SKU uint32
}

// ClientMockGetProductResults contains results of the Client.GetProduct
type ClientMockGetProductResults struct {
	pp1 *model.Product
	err error
}

// Expect sets up expected params for Client.GetProduct
func (mmGetProduct *mClientMockGetProduct) Expect(ctx context.Context, SKU uint32) *mClientMockGetProduct {
	if mmGetProduct.mock.funcGetProduct != nil {
		mmGetProduct.mock.t.Fatalf("ClientMock.GetProduct mock is already set by Set")
	}

	if mmGetProduct.defaultExpectation == nil {
		mmGetProduct.defaultExpectation = &ClientMockGetProductExpectation{}
	}

	mmGetProduct.defaultExpectation.params = &ClientMockGetProductParams{ctx, SKU}
	for _, e := range mmGetProduct.expectations {
		if minimock.Equal(e.params, mmGetProduct.defaultExpectation.params) {
			mmGetProduct.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetProduct.defaultExpectation.params)
		}
	}

	return mmGetProduct
}

// Inspect accepts an inspector function that has same arguments as the Client.GetProduct
func (mmGetProduct *mClientMockGetProduct) Inspect(f func(ctx context.Context, SKU uint32)) *mClientMockGetProduct {
	if mmGetProduct.mock.inspectFuncGetProduct != nil {
		mmGetProduct.mock.t.Fatalf("Inspect function is already set for ClientMock.GetProduct")
	}

	mmGetProduct.mock.inspectFuncGetProduct = f

	return mmGetProduct
}

// Return sets up results that will be returned by Client.GetProduct
func (mmGetProduct *mClientMockGetProduct) Return(pp1 *model.Product, err error) *ClientMock {
	if mmGetProduct.mock.funcGetProduct != nil {
		mmGetProduct.mock.t.Fatalf("ClientMock.GetProduct mock is already set by Set")
	}

	if mmGetProduct.defaultExpectation == nil {
		mmGetProduct.defaultExpectation = &ClientMockGetProductExpectation{mock: mmGetProduct.mock}
	}
	mmGetProduct.defaultExpectation.results = &ClientMockGetProductResults{pp1, err}
	return mmGetProduct.mock
}

// Set uses given function f to mock the Client.GetProduct method
func (mmGetProduct *mClientMockGetProduct) Set(f func(ctx context.Context, SKU uint32) (pp1 *model.Product, err error)) *ClientMock {
	if mmGetProduct.defaultExpectation != nil {
		mmGetProduct.mock.t.Fatalf("Default expectation is already set for the Client.GetProduct method")
	}

	if len(mmGetProduct.expectations) > 0 {
		mmGetProduct.mock.t.Fatalf("Some expectations are already set for the Client.GetProduct method")
	}

	mmGetProduct.mock.funcGetProduct = f
	return mmGetProduct.mock
}

// When sets expectation for the Client.GetProduct which will trigger the result defined by the following
// Then helper
func (mmGetProduct *mClientMockGetProduct) When(ctx context.Context, SKU uint32) *ClientMockGetProductExpectation {
	if mmGetProduct.mock.funcGetProduct != nil {
		mmGetProduct.mock.t.Fatalf("ClientMock.GetProduct mock is already set by Set")
	}

	expectation := &ClientMockGetProductExpectation{
		mock:   mmGetProduct.mock,
		params: &ClientMockGetProductParams{ctx, SKU},
	}
	mmGetProduct.expectations = append(mmGetProduct.expectations, expectation)
	return expectation
}

// Then sets up Client.GetProduct return parameters for the expectation previously defined by the When method
func (e *ClientMockGetProductExpectation) Then(pp1 *model.Product, err error) *ClientMock {
	e.results = &ClientMockGetProductResults{pp1, err}
	return e.mock
}

// GetProduct implements productservice.Client
func (mmGetProduct *ClientMock) GetProduct(ctx context.Context, SKU uint32) (pp1 *model.Product, err error) {
	mm_atomic.AddUint64(&mmGetProduct.beforeGetProductCounter, 1)
	defer mm_atomic.AddUint64(&mmGetProduct.afterGetProductCounter, 1)

	if mmGetProduct.inspectFuncGetProduct != nil {
		mmGetProduct.inspectFuncGetProduct(ctx, SKU)
	}

	mm_params := &ClientMockGetProductParams{ctx, SKU}

	// Record call args
	mmGetProduct.GetProductMock.mutex.Lock()
	mmGetProduct.GetProductMock.callArgs = append(mmGetProduct.GetProductMock.callArgs, mm_params)
	mmGetProduct.GetProductMock.mutex.Unlock()

	for _, e := range mmGetProduct.GetProductMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.pp1, e.results.err
		}
	}

	if mmGetProduct.GetProductMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetProduct.GetProductMock.defaultExpectation.Counter, 1)
		mm_want := mmGetProduct.GetProductMock.defaultExpectation.params
		mm_got := ClientMockGetProductParams{ctx, SKU}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetProduct.t.Errorf("ClientMock.GetProduct got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetProduct.GetProductMock.defaultExpectation.results
		if mm_results == nil {
			mmGetProduct.t.Fatal("No results are set for the ClientMock.GetProduct")
		}
		return (*mm_results).pp1, (*mm_results).err
	}
	if mmGetProduct.funcGetProduct != nil {
		return mmGetProduct.funcGetProduct(ctx, SKU)
	}
	mmGetProduct.t.Fatalf("Unexpected call to ClientMock.GetProduct. %v %v", ctx, SKU)
	return
}

// GetProductAfterCounter returns a count of finished ClientMock.GetProduct invocations
func (mmGetProduct *ClientMock) GetProductAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetProduct.afterGetProductCounter)
}

// GetProductBeforeCounter returns a count of ClientMock.GetProduct invocations
func (mmGetProduct *ClientMock) GetProductBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetProduct.beforeGetProductCounter)
}

// Calls returns a list of arguments used in each call to ClientMock.GetProduct.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetProduct *mClientMockGetProduct) Calls() []*ClientMockGetProductParams {
	mmGetProduct.mutex.RLock()

	argCopy := make([]*ClientMockGetProductParams, len(mmGetProduct.callArgs))
	copy(argCopy, mmGetProduct.callArgs)

	mmGetProduct.mutex.RUnlock()

	return argCopy
}

// MinimockGetProductDone returns true if the count of the GetProduct invocations corresponds
// the number of defined expectations
func (m *ClientMock) MinimockGetProductDone() bool {
	for _, e := range m.GetProductMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetProductMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetProductCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetProduct != nil && mm_atomic.LoadUint64(&m.afterGetProductCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetProductInspect logs each unmet expectation
func (m *ClientMock) MinimockGetProductInspect() {
	for _, e := range m.GetProductMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ClientMock.GetProduct with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetProductMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetProductCounter) < 1 {
		if m.GetProductMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to ClientMock.GetProduct")
		} else {
			m.t.Errorf("Expected call to ClientMock.GetProduct with params: %#v", *m.GetProductMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetProduct != nil && mm_atomic.LoadUint64(&m.afterGetProductCounter) < 1 {
		m.t.Error("Expected call to ClientMock.GetProduct")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ClientMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockGetProductInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ClientMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *ClientMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetProductDone()
}