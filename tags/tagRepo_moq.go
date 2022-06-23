// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package tags

import (
	"context"
	"github.com/google/uuid"
	"sync"
)

// Ensure, that RepositoryMock does implement Repository.
// If this is not the case, regenerate this file with moq.
var _ Repository = &RepositoryMock{}

// RepositoryMock is a mock implementation of Repository.
//
// 	func TestSomethingThatUsesRepository(t *testing.T) {
//
// 		// make and configure a mocked Repository
// 		mockedRepository := &RepositoryMock{
// 			CountFunc: func(contextMoqParam context.Context, uUID uuid.UUID) (int, error) {
// 				panic("mock out the Count method")
// 			},
// 			DeleteFunc: func(contextMoqParam context.Context, uUID uuid.UUID) error {
// 				panic("mock out the Delete method")
// 			},
// 			GetAllFunc: func(contextMoqParam context.Context) ([]Tags, error) {
// 				panic("mock out the GetAll method")
// 			},
// 			GetByIdFunc: func(contextMoqParam context.Context, uUID uuid.UUID) (*Tags, error) {
// 				panic("mock out the GetById method")
// 			},
// 			GetByIdsFunc: func(contextMoqParam context.Context, uUIDs []uuid.UUID) ([]Tags, error) {
// 				panic("mock out the GetByIds method")
// 			},
// 			GetByNameFunc: func(ctx context.Context, name string) (Tags, error) {
// 				panic("mock out the GetByName method")
// 			},
// 			GetByNamesFunc: func(contextMoqParam context.Context, strings ...string) ([]Tags, error) {
// 				panic("mock out the GetByNames method")
// 			},
// 			SaveFunc: func(contextMoqParam context.Context, tags Tags) error {
// 				panic("mock out the Save method")
// 			},
// 			UpdateFunc: func(contextMoqParam context.Context, tags Tags) error {
// 				panic("mock out the Update method")
// 			},
// 		}
//
// 		// use mockedRepository in code that requires Repository
// 		// and then make assertions.
//
// 	}
type RepositoryMock struct {
	// CountFunc mocks the Count method.
	CountFunc func(contextMoqParam context.Context, uUID uuid.UUID) (int, error)

	// DeleteFunc mocks the Delete method.
	DeleteFunc func(contextMoqParam context.Context, uUID uuid.UUID) error

	// GetAllFunc mocks the GetAll method.
	GetAllFunc func(contextMoqParam context.Context) ([]Tags, error)

	// GetByIdFunc mocks the GetById method.
	GetByIdFunc func(contextMoqParam context.Context, uUID uuid.UUID) (*Tags, error)

	// GetByIdsFunc mocks the GetByIds method.
	GetByIdsFunc func(contextMoqParam context.Context, uUIDs []uuid.UUID) ([]Tags, error)

	// GetByNameFunc mocks the GetByName method.
	GetByNameFunc func(ctx context.Context, name string) (Tags, error)

	// GetByNamesFunc mocks the GetByNames method.
	GetByNamesFunc func(contextMoqParam context.Context, strings ...string) ([]Tags, error)

	// SaveFunc mocks the Save method.
	SaveFunc func(contextMoqParam context.Context, tags Tags) error

	// UpdateFunc mocks the Update method.
	UpdateFunc func(contextMoqParam context.Context, tags Tags) error

	// calls tracks calls to the methods.
	calls struct {
		// Count holds details about calls to the Count method.
		Count []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// UUID is the uUID argument value.
			UUID uuid.UUID
		}
		// Delete holds details about calls to the Delete method.
		Delete []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// UUID is the uUID argument value.
			UUID uuid.UUID
		}
		// GetAll holds details about calls to the GetAll method.
		GetAll []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
		}
		// GetById holds details about calls to the GetById method.
		GetById []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// UUID is the uUID argument value.
			UUID uuid.UUID
		}
		// GetByIds holds details about calls to the GetByIds method.
		GetByIds []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// UUIDs is the uUIDs argument value.
			UUIDs []uuid.UUID
		}
		// GetByName holds details about calls to the GetByName method.
		GetByName []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Name is the name argument value.
			Name string
		}
		// GetByNames holds details about calls to the GetByNames method.
		GetByNames []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// Strings is the strings argument value.
			Strings []string
		}
		// Save holds details about calls to the Save method.
		Save []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// Tags is the tags argument value.
			Tags Tags
		}
		// Update holds details about calls to the Update method.
		Update []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// Tags is the tags argument value.
			Tags Tags
		}
	}
	lockCount      sync.RWMutex
	lockDelete     sync.RWMutex
	lockGetAll     sync.RWMutex
	lockGetById    sync.RWMutex
	lockGetByIds   sync.RWMutex
	lockGetByName  sync.RWMutex
	lockGetByNames sync.RWMutex
	lockSave       sync.RWMutex
	lockUpdate     sync.RWMutex
}

// Count calls CountFunc.
func (mock *RepositoryMock) Count(contextMoqParam context.Context, uUID uuid.UUID) (int, error) {
	if mock.CountFunc == nil {
		panic("RepositoryMock.CountFunc: method is nil but Repository.Count was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		UUID            uuid.UUID
	}{
		ContextMoqParam: contextMoqParam,
		UUID:            uUID,
	}
	mock.lockCount.Lock()
	mock.calls.Count = append(mock.calls.Count, callInfo)
	mock.lockCount.Unlock()
	return mock.CountFunc(contextMoqParam, uUID)
}

// CountCalls gets all the calls that were made to Count.
// Check the length with:
//     len(mockedRepository.CountCalls())
func (mock *RepositoryMock) CountCalls() []struct {
	ContextMoqParam context.Context
	UUID            uuid.UUID
} {
	var calls []struct {
		ContextMoqParam context.Context
		UUID            uuid.UUID
	}
	mock.lockCount.RLock()
	calls = mock.calls.Count
	mock.lockCount.RUnlock()
	return calls
}

// Delete calls DeleteFunc.
func (mock *RepositoryMock) Delete(contextMoqParam context.Context, uUID uuid.UUID) error {
	if mock.DeleteFunc == nil {
		panic("RepositoryMock.DeleteFunc: method is nil but Repository.Delete was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		UUID            uuid.UUID
	}{
		ContextMoqParam: contextMoqParam,
		UUID:            uUID,
	}
	mock.lockDelete.Lock()
	mock.calls.Delete = append(mock.calls.Delete, callInfo)
	mock.lockDelete.Unlock()
	return mock.DeleteFunc(contextMoqParam, uUID)
}

// DeleteCalls gets all the calls that were made to Delete.
// Check the length with:
//     len(mockedRepository.DeleteCalls())
func (mock *RepositoryMock) DeleteCalls() []struct {
	ContextMoqParam context.Context
	UUID            uuid.UUID
} {
	var calls []struct {
		ContextMoqParam context.Context
		UUID            uuid.UUID
	}
	mock.lockDelete.RLock()
	calls = mock.calls.Delete
	mock.lockDelete.RUnlock()
	return calls
}

// GetAll calls GetAllFunc.
func (mock *RepositoryMock) GetAll(contextMoqParam context.Context) ([]Tags, error) {
	if mock.GetAllFunc == nil {
		panic("RepositoryMock.GetAllFunc: method is nil but Repository.GetAll was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
	}{
		ContextMoqParam: contextMoqParam,
	}
	mock.lockGetAll.Lock()
	mock.calls.GetAll = append(mock.calls.GetAll, callInfo)
	mock.lockGetAll.Unlock()
	return mock.GetAllFunc(contextMoqParam)
}

// GetAllCalls gets all the calls that were made to GetAll.
// Check the length with:
//     len(mockedRepository.GetAllCalls())
func (mock *RepositoryMock) GetAllCalls() []struct {
	ContextMoqParam context.Context
} {
	var calls []struct {
		ContextMoqParam context.Context
	}
	mock.lockGetAll.RLock()
	calls = mock.calls.GetAll
	mock.lockGetAll.RUnlock()
	return calls
}

// GetById calls GetByIdFunc.
func (mock *RepositoryMock) GetById(contextMoqParam context.Context, uUID uuid.UUID) (*Tags, error) {
	if mock.GetByIdFunc == nil {
		panic("RepositoryMock.GetByIdFunc: method is nil but Repository.GetById was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		UUID            uuid.UUID
	}{
		ContextMoqParam: contextMoqParam,
		UUID:            uUID,
	}
	mock.lockGetById.Lock()
	mock.calls.GetById = append(mock.calls.GetById, callInfo)
	mock.lockGetById.Unlock()
	return mock.GetByIdFunc(contextMoqParam, uUID)
}

// GetByIdCalls gets all the calls that were made to GetById.
// Check the length with:
//     len(mockedRepository.GetByIdCalls())
func (mock *RepositoryMock) GetByIdCalls() []struct {
	ContextMoqParam context.Context
	UUID            uuid.UUID
} {
	var calls []struct {
		ContextMoqParam context.Context
		UUID            uuid.UUID
	}
	mock.lockGetById.RLock()
	calls = mock.calls.GetById
	mock.lockGetById.RUnlock()
	return calls
}

// GetByIds calls GetByIdsFunc.
func (mock *RepositoryMock) GetByIds(contextMoqParam context.Context, uUIDs []uuid.UUID) ([]Tags, error) {
	if mock.GetByIdsFunc == nil {
		panic("RepositoryMock.GetByIdsFunc: method is nil but Repository.GetByIds was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		UUIDs           []uuid.UUID
	}{
		ContextMoqParam: contextMoqParam,
		UUIDs:           uUIDs,
	}
	mock.lockGetByIds.Lock()
	mock.calls.GetByIds = append(mock.calls.GetByIds, callInfo)
	mock.lockGetByIds.Unlock()
	return mock.GetByIdsFunc(contextMoqParam, uUIDs)
}

// GetByIdsCalls gets all the calls that were made to GetByIds.
// Check the length with:
//     len(mockedRepository.GetByIdsCalls())
func (mock *RepositoryMock) GetByIdsCalls() []struct {
	ContextMoqParam context.Context
	UUIDs           []uuid.UUID
} {
	var calls []struct {
		ContextMoqParam context.Context
		UUIDs           []uuid.UUID
	}
	mock.lockGetByIds.RLock()
	calls = mock.calls.GetByIds
	mock.lockGetByIds.RUnlock()
	return calls
}

// GetByName calls GetByNameFunc.
func (mock *RepositoryMock) GetByName(ctx context.Context, name string) (Tags, error) {
	if mock.GetByNameFunc == nil {
		panic("RepositoryMock.GetByNameFunc: method is nil but Repository.GetByName was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Name string
	}{
		Ctx:  ctx,
		Name: name,
	}
	mock.lockGetByName.Lock()
	mock.calls.GetByName = append(mock.calls.GetByName, callInfo)
	mock.lockGetByName.Unlock()
	return mock.GetByNameFunc(ctx, name)
}

// GetByNameCalls gets all the calls that were made to GetByName.
// Check the length with:
//     len(mockedRepository.GetByNameCalls())
func (mock *RepositoryMock) GetByNameCalls() []struct {
	Ctx  context.Context
	Name string
} {
	var calls []struct {
		Ctx  context.Context
		Name string
	}
	mock.lockGetByName.RLock()
	calls = mock.calls.GetByName
	mock.lockGetByName.RUnlock()
	return calls
}

// GetByNames calls GetByNamesFunc.
func (mock *RepositoryMock) GetByNames(contextMoqParam context.Context, strings ...string) ([]Tags, error) {
	if mock.GetByNamesFunc == nil {
		panic("RepositoryMock.GetByNamesFunc: method is nil but Repository.GetByNames was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		Strings         []string
	}{
		ContextMoqParam: contextMoqParam,
		Strings:         strings,
	}
	mock.lockGetByNames.Lock()
	mock.calls.GetByNames = append(mock.calls.GetByNames, callInfo)
	mock.lockGetByNames.Unlock()
	return mock.GetByNamesFunc(contextMoqParam, strings...)
}

// GetByNamesCalls gets all the calls that were made to GetByNames.
// Check the length with:
//     len(mockedRepository.GetByNamesCalls())
func (mock *RepositoryMock) GetByNamesCalls() []struct {
	ContextMoqParam context.Context
	Strings         []string
} {
	var calls []struct {
		ContextMoqParam context.Context
		Strings         []string
	}
	mock.lockGetByNames.RLock()
	calls = mock.calls.GetByNames
	mock.lockGetByNames.RUnlock()
	return calls
}

// Save calls SaveFunc.
func (mock *RepositoryMock) Save(contextMoqParam context.Context, tags Tags) error {
	if mock.SaveFunc == nil {
		panic("RepositoryMock.SaveFunc: method is nil but Repository.Save was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		Tags            Tags
	}{
		ContextMoqParam: contextMoqParam,
		Tags:            tags,
	}
	mock.lockSave.Lock()
	mock.calls.Save = append(mock.calls.Save, callInfo)
	mock.lockSave.Unlock()
	return mock.SaveFunc(contextMoqParam, tags)
}

// SaveCalls gets all the calls that were made to Save.
// Check the length with:
//     len(mockedRepository.SaveCalls())
func (mock *RepositoryMock) SaveCalls() []struct {
	ContextMoqParam context.Context
	Tags            Tags
} {
	var calls []struct {
		ContextMoqParam context.Context
		Tags            Tags
	}
	mock.lockSave.RLock()
	calls = mock.calls.Save
	mock.lockSave.RUnlock()
	return calls
}

// Update calls UpdateFunc.
func (mock *RepositoryMock) Update(contextMoqParam context.Context, tags Tags) error {
	if mock.UpdateFunc == nil {
		panic("RepositoryMock.UpdateFunc: method is nil but Repository.Update was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		Tags            Tags
	}{
		ContextMoqParam: contextMoqParam,
		Tags:            tags,
	}
	mock.lockUpdate.Lock()
	mock.calls.Update = append(mock.calls.Update, callInfo)
	mock.lockUpdate.Unlock()
	return mock.UpdateFunc(contextMoqParam, tags)
}

// UpdateCalls gets all the calls that were made to Update.
// Check the length with:
//     len(mockedRepository.UpdateCalls())
func (mock *RepositoryMock) UpdateCalls() []struct {
	ContextMoqParam context.Context
	Tags            Tags
} {
	var calls []struct {
		ContextMoqParam context.Context
		Tags            Tags
	}
	mock.lockUpdate.RLock()
	calls = mock.calls.Update
	mock.lockUpdate.RUnlock()
	return calls
}
