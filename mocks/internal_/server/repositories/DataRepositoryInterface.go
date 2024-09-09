// Code generated by mockery v2.43.2. DO NOT EDIT.

package repositories

import (
	models "github.com/ShukinDmitriy/GophKeeper/internal/common/models"
	mock "github.com/stretchr/testify/mock"

	requests "github.com/ShukinDmitriy/GophKeeper/internal/common/models/requests"
)

// DataRepositoryInterface is an autogenerated mock type for the DataRepositoryInterface type
type DataRepositoryInterface struct {
	mock.Mock
}

type DataRepositoryInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *DataRepositoryInterface) EXPECT() *DataRepositoryInterface_Expecter {
	return &DataRepositoryInterface_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: dataCreate
func (_m *DataRepositoryInterface) Create(dataCreate requests.DataModel) (*models.DataInfo, error) {
	ret := _m.Called(dataCreate)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *models.DataInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(requests.DataModel) (*models.DataInfo, error)); ok {
		return rf(dataCreate)
	}
	if rf, ok := ret.Get(0).(func(requests.DataModel) *models.DataInfo); ok {
		r0 = rf(dataCreate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DataInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(requests.DataModel) error); ok {
		r1 = rf(dataCreate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DataRepositoryInterface_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type DataRepositoryInterface_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - dataCreate requests.DataModel
func (_e *DataRepositoryInterface_Expecter) Create(dataCreate interface{}) *DataRepositoryInterface_Create_Call {
	return &DataRepositoryInterface_Create_Call{Call: _e.mock.On("Create", dataCreate)}
}

func (_c *DataRepositoryInterface_Create_Call) Run(run func(dataCreate requests.DataModel)) *DataRepositoryInterface_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(requests.DataModel))
	})
	return _c
}

func (_c *DataRepositoryInterface_Create_Call) Return(_a0 *models.DataInfo, _a1 error) *DataRepositoryInterface_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DataRepositoryInterface_Create_Call) RunAndReturn(run func(requests.DataModel) (*models.DataInfo, error)) *DataRepositoryInterface_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: id, userID
func (_m *DataRepositoryInterface) Delete(id uint, userID uint) error {
	ret := _m.Called(id, userID)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uint, uint) error); ok {
		r0 = rf(id, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DataRepositoryInterface_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type DataRepositoryInterface_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - id uint
//   - userID uint
func (_e *DataRepositoryInterface_Expecter) Delete(id interface{}, userID interface{}) *DataRepositoryInterface_Delete_Call {
	return &DataRepositoryInterface_Delete_Call{Call: _e.mock.On("Delete", id, userID)}
}

func (_c *DataRepositoryInterface_Delete_Call) Run(run func(id uint, userID uint)) *DataRepositoryInterface_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint), args[1].(uint))
	})
	return _c
}

func (_c *DataRepositoryInterface_Delete_Call) Return(_a0 error) *DataRepositoryInterface_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DataRepositoryInterface_Delete_Call) RunAndReturn(run func(uint, uint) error) *DataRepositoryInterface_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Find provides a mock function with given fields: id, userID
func (_m *DataRepositoryInterface) Find(id uint, userID uint) (*models.DataInfo, error) {
	ret := _m.Called(id, userID)

	if len(ret) == 0 {
		panic("no return value specified for Find")
	}

	var r0 *models.DataInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(uint, uint) (*models.DataInfo, error)); ok {
		return rf(id, userID)
	}
	if rf, ok := ret.Get(0).(func(uint, uint) *models.DataInfo); ok {
		r0 = rf(id, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DataInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(uint, uint) error); ok {
		r1 = rf(id, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DataRepositoryInterface_Find_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Find'
type DataRepositoryInterface_Find_Call struct {
	*mock.Call
}

// Find is a helper method to define mock.On call
//   - id uint
//   - userID uint
func (_e *DataRepositoryInterface_Expecter) Find(id interface{}, userID interface{}) *DataRepositoryInterface_Find_Call {
	return &DataRepositoryInterface_Find_Call{Call: _e.mock.On("Find", id, userID)}
}

func (_c *DataRepositoryInterface_Find_Call) Run(run func(id uint, userID uint)) *DataRepositoryInterface_Find_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint), args[1].(uint))
	})
	return _c
}

func (_c *DataRepositoryInterface_Find_Call) Return(_a0 *models.DataInfo, _a1 error) *DataRepositoryInterface_Find_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DataRepositoryInterface_Find_Call) RunAndReturn(run func(uint, uint) (*models.DataInfo, error)) *DataRepositoryInterface_Find_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: request
func (_m *DataRepositoryInterface) List(request requests.DataList) ([]*models.DataInfo, error) {
	ret := _m.Called(request)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []*models.DataInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(requests.DataList) ([]*models.DataInfo, error)); ok {
		return rf(request)
	}
	if rf, ok := ret.Get(0).(func(requests.DataList) []*models.DataInfo); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.DataInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(requests.DataList) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DataRepositoryInterface_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type DataRepositoryInterface_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - request requests.DataList
func (_e *DataRepositoryInterface_Expecter) List(request interface{}) *DataRepositoryInterface_List_Call {
	return &DataRepositoryInterface_List_Call{Call: _e.mock.On("List", request)}
}

func (_c *DataRepositoryInterface_List_Call) Run(run func(request requests.DataList)) *DataRepositoryInterface_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(requests.DataList))
	})
	return _c
}

func (_c *DataRepositoryInterface_List_Call) Return(_a0 []*models.DataInfo, _a1 error) *DataRepositoryInterface_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DataRepositoryInterface_List_Call) RunAndReturn(run func(requests.DataList) ([]*models.DataInfo, error)) *DataRepositoryInterface_List_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: id, request
func (_m *DataRepositoryInterface) Update(id uint, request requests.DataModel) (*models.DataInfo, error) {
	ret := _m.Called(id, request)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *models.DataInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(uint, requests.DataModel) (*models.DataInfo, error)); ok {
		return rf(id, request)
	}
	if rf, ok := ret.Get(0).(func(uint, requests.DataModel) *models.DataInfo); ok {
		r0 = rf(id, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DataInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(uint, requests.DataModel) error); ok {
		r1 = rf(id, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DataRepositoryInterface_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type DataRepositoryInterface_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - id uint
//   - request requests.DataModel
func (_e *DataRepositoryInterface_Expecter) Update(id interface{}, request interface{}) *DataRepositoryInterface_Update_Call {
	return &DataRepositoryInterface_Update_Call{Call: _e.mock.On("Update", id, request)}
}

func (_c *DataRepositoryInterface_Update_Call) Run(run func(id uint, request requests.DataModel)) *DataRepositoryInterface_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint), args[1].(requests.DataModel))
	})
	return _c
}

func (_c *DataRepositoryInterface_Update_Call) Return(_a0 *models.DataInfo, _a1 error) *DataRepositoryInterface_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DataRepositoryInterface_Update_Call) RunAndReturn(run func(uint, requests.DataModel) (*models.DataInfo, error)) *DataRepositoryInterface_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewDataRepositoryInterface creates a new instance of DataRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDataRepositoryInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *DataRepositoryInterface {
	mock := &DataRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
