// Code generated by mockery v2.46.0. DO NOT EDIT.

package mockservice

import (
	model "avito-shop-service/internal/model"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

type Repository_Expecter struct {
	mock *mock.Mock
}

func (_m *Repository) EXPECT() *Repository_Expecter {
	return &Repository_Expecter{mock: &_m.Mock}
}

// Balance provides a mock function with given fields: ctx, username
func (_m *Repository) Balance(ctx context.Context, username string) (int64, error) {
	ret := _m.Called(ctx, username)

	if len(ret) == 0 {
		panic("no return value specified for Balance")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (int64, error)); ok {
		return rf(ctx, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) int64); ok {
		r0 = rf(ctx, username)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_Balance_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Balance'
type Repository_Balance_Call struct {
	*mock.Call
}

// Balance is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
func (_e *Repository_Expecter) Balance(ctx interface{}, username interface{}) *Repository_Balance_Call {
	return &Repository_Balance_Call{Call: _e.mock.On("Balance", ctx, username)}
}

func (_c *Repository_Balance_Call) Run(run func(ctx context.Context, username string)) *Repository_Balance_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Repository_Balance_Call) Return(_a0 int64, _a1 error) *Repository_Balance_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_Balance_Call) RunAndReturn(run func(context.Context, string) (int64, error)) *Repository_Balance_Call {
	_c.Call.Return(run)
	return _c
}

// CreateUser provides a mock function with given fields: ctx, username, passwordHash
func (_m *Repository) CreateUser(ctx context.Context, username string, passwordHash string) (model.User, error) {
	ret := _m.Called(ctx, username, passwordHash)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (model.User, error)); ok {
		return rf(ctx, username, passwordHash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) model.User); ok {
		r0 = rf(ctx, username, passwordHash)
	} else {
		r0 = ret.Get(0).(model.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, username, passwordHash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_CreateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUser'
type Repository_CreateUser_Call struct {
	*mock.Call
}

// CreateUser is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
//   - passwordHash string
func (_e *Repository_Expecter) CreateUser(ctx interface{}, username interface{}, passwordHash interface{}) *Repository_CreateUser_Call {
	return &Repository_CreateUser_Call{Call: _e.mock.On("CreateUser", ctx, username, passwordHash)}
}

func (_c *Repository_CreateUser_Call) Run(run func(ctx context.Context, username string, passwordHash string)) *Repository_CreateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *Repository_CreateUser_Call) Return(_a0 model.User, _a1 error) *Repository_CreateUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_CreateUser_Call) RunAndReturn(run func(context.Context, string, string) (model.User, error)) *Repository_CreateUser_Call {
	_c.Call.Return(run)
	return _c
}

// GetUser provides a mock function with given fields: ctx, username
func (_m *Repository) GetUser(ctx context.Context, username string) (model.User, error) {
	ret := _m.Called(ctx, username)

	if len(ret) == 0 {
		panic("no return value specified for GetUser")
	}

	var r0 model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (model.User, error)); ok {
		return rf(ctx, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) model.User); ok {
		r0 = rf(ctx, username)
	} else {
		r0 = ret.Get(0).(model.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_GetUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUser'
type Repository_GetUser_Call struct {
	*mock.Call
}

// GetUser is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
func (_e *Repository_Expecter) GetUser(ctx interface{}, username interface{}) *Repository_GetUser_Call {
	return &Repository_GetUser_Call{Call: _e.mock.On("GetUser", ctx, username)}
}

func (_c *Repository_GetUser_Call) Run(run func(ctx context.Context, username string)) *Repository_GetUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Repository_GetUser_Call) Return(_a0 model.User, _a1 error) *Repository_GetUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_GetUser_Call) RunAndReturn(run func(context.Context, string) (model.User, error)) *Repository_GetUser_Call {
	_c.Call.Return(run)
	return _c
}

// Inventory provides a mock function with given fields: ctx, username
func (_m *Repository) Inventory(ctx context.Context, username string) ([]model.Inventory, error) {
	ret := _m.Called(ctx, username)

	if len(ret) == 0 {
		panic("no return value specified for Inventory")
	}

	var r0 []model.Inventory
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]model.Inventory, error)); ok {
		return rf(ctx, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []model.Inventory); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Inventory)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_Inventory_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Inventory'
type Repository_Inventory_Call struct {
	*mock.Call
}

// Inventory is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
func (_e *Repository_Expecter) Inventory(ctx interface{}, username interface{}) *Repository_Inventory_Call {
	return &Repository_Inventory_Call{Call: _e.mock.On("Inventory", ctx, username)}
}

func (_c *Repository_Inventory_Call) Run(run func(ctx context.Context, username string)) *Repository_Inventory_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Repository_Inventory_Call) Return(_a0 []model.Inventory, _a1 error) *Repository_Inventory_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_Inventory_Call) RunAndReturn(run func(context.Context, string) ([]model.Inventory, error)) *Repository_Inventory_Call {
	_c.Call.Return(run)
	return _c
}

// SwapBalance provides a mock function with given fields: ctx, fromUser, toUser, amount
func (_m *Repository) SwapBalance(ctx context.Context, fromUser string, toUser string, amount int) error {
	ret := _m.Called(ctx, fromUser, toUser, amount)

	if len(ret) == 0 {
		panic("no return value specified for SwapBalance")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int) error); ok {
		r0 = rf(ctx, fromUser, toUser, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Repository_SwapBalance_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SwapBalance'
type Repository_SwapBalance_Call struct {
	*mock.Call
}

// SwapBalance is a helper method to define mock.On call
//   - ctx context.Context
//   - fromUser string
//   - toUser string
//   - amount int
func (_e *Repository_Expecter) SwapBalance(ctx interface{}, fromUser interface{}, toUser interface{}, amount interface{}) *Repository_SwapBalance_Call {
	return &Repository_SwapBalance_Call{Call: _e.mock.On("SwapBalance", ctx, fromUser, toUser, amount)}
}

func (_c *Repository_SwapBalance_Call) Run(run func(ctx context.Context, fromUser string, toUser string, amount int)) *Repository_SwapBalance_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(int))
	})
	return _c
}

func (_c *Repository_SwapBalance_Call) Return(_a0 error) *Repository_SwapBalance_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Repository_SwapBalance_Call) RunAndReturn(run func(context.Context, string, string, int) error) *Repository_SwapBalance_Call {
	_c.Call.Return(run)
	return _c
}

// Transaction provides a mock function with given fields: ctx, username
func (_m *Repository) Transaction(ctx context.Context, username string) ([]model.Transaction, error) {
	ret := _m.Called(ctx, username)

	if len(ret) == 0 {
		panic("no return value specified for Transaction")
	}

	var r0 []model.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]model.Transaction, error)); ok {
		return rf(ctx, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []model.Transaction); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_Transaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Transaction'
type Repository_Transaction_Call struct {
	*mock.Call
}

// Transaction is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
func (_e *Repository_Expecter) Transaction(ctx interface{}, username interface{}) *Repository_Transaction_Call {
	return &Repository_Transaction_Call{Call: _e.mock.On("Transaction", ctx, username)}
}

func (_c *Repository_Transaction_Call) Run(run func(ctx context.Context, username string)) *Repository_Transaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Repository_Transaction_Call) Return(_a0 []model.Transaction, _a1 error) *Repository_Transaction_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_Transaction_Call) RunAndReturn(run func(context.Context, string) ([]model.Transaction, error)) *Repository_Transaction_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateBalance provides a mock function with given fields: ctx, username, itemName
func (_m *Repository) UpdateBalance(ctx context.Context, username string, itemName string) error {
	ret := _m.Called(ctx, username, itemName)

	if len(ret) == 0 {
		panic("no return value specified for UpdateBalance")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, username, itemName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Repository_UpdateBalance_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateBalance'
type Repository_UpdateBalance_Call struct {
	*mock.Call
}

// UpdateBalance is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
//   - itemName string
func (_e *Repository_Expecter) UpdateBalance(ctx interface{}, username interface{}, itemName interface{}) *Repository_UpdateBalance_Call {
	return &Repository_UpdateBalance_Call{Call: _e.mock.On("UpdateBalance", ctx, username, itemName)}
}

func (_c *Repository_UpdateBalance_Call) Run(run func(ctx context.Context, username string, itemName string)) *Repository_UpdateBalance_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *Repository_UpdateBalance_Call) Return(_a0 error) *Repository_UpdateBalance_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Repository_UpdateBalance_Call) RunAndReturn(run func(context.Context, string, string) error) *Repository_UpdateBalance_Call {
	_c.Call.Return(run)
	return _c
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
