// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package tag

import (
	"github.com/Iiqbal2000/bareknews/domain"
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
// 			DeleteFunc: func(s string) error {
// 				panic("mock out the Delete method")
// 			},
// 			GetAllFunc: func() ([]bareknews.Tags, error) {
// 				panic("mock out the GetAll method")
// 			},
// 			GetByIdFunc: func(s string) (*bareknews.Tags, error) {
// 				panic("mock out the GetById method")
// 			},
// 			GetByNameFunc: func(names string) (bareknews.Tags, error) {
// 				panic("mock out the GetByName method")
// 			},
// 			GetByNamesFunc: func(names ...string) ([]bareknews.Tags, error) {
// 				panic("mock out the GetByNames method")
// 			},
// 			GetByNewsIdFunc: func(id string) ([]bareknews.Tags, error) {
// 				panic("mock out the GetByNewsId method")
// 			},
// 			SaveFunc: func(tags bareknews.Tags) error {
// 				panic("mock out the Save method")
// 			},
// 			UpdateFunc: func(tags bareknews.Tags) error {
// 				panic("mock out the Update method")
// 			},
// 		}
//
// 		// use mockedRepository in code that requires Repository
// 		// and then make assertions.
//
// 	}
type RepositoryMock struct {
	// DeleteFunc mocks the Delete method.
	DeleteFunc func(s string) error

	// GetAllFunc mocks the GetAll method.
	GetAllFunc func() ([]domain.Tags, error)

	// GetByIdFunc mocks the GetById method.
	GetByIdFunc func(s string) (*domain.Tags, error)

	// GetByNameFunc mocks the GetByName method.
	GetByNameFunc func(names string) (domain.Tags, error)

	// GetByNamesFunc mocks the GetByNames method.
	GetByNamesFunc func(names ...string) ([]domain.Tags, error)

	// GetByNewsIdFunc mocks the GetByNewsId method.
	GetByNewsIdFunc func(id string) ([]domain.Tags, error)

	// SaveFunc mocks the Save method.
	SaveFunc func(tags domain.Tags) error

	// UpdateFunc mocks the Update method.
	UpdateFunc func(tags domain.Tags) error

	// calls tracks calls to the methods.
	calls struct {
		// Delete holds details about calls to the Delete method.
		Delete []struct {
			// S is the s argument value.
			S string
		}
		// GetAll holds details about calls to the GetAll method.
		GetAll []struct {
		}
		// GetById holds details about calls to the GetById method.
		GetById []struct {
			// S is the s argument value.
			S string
		}
		// GetByName holds details about calls to the GetByName method.
		GetByName []struct {
			// Names is the names argument value.
			Names string
		}
		// GetByNames holds details about calls to the GetByNames method.
		GetByNames []struct {
			// Names is the names argument value.
			Names []string
		}
		// GetByNewsId holds details about calls to the GetByNewsId method.
		GetByNewsId []struct {
			// ID is the id argument value.
			ID string
		}
		// Save holds details about calls to the Save method.
		Save []struct {
			// Tags is the tags argument value.
			Tags domain.Tags
		}
		// Update holds details about calls to the Update method.
		Update []struct {
			// Tags is the tags argument value.
			Tags domain.Tags
		}
	}
	lockDelete      sync.RWMutex
	lockGetAll      sync.RWMutex
	lockGetById     sync.RWMutex
	lockGetByName   sync.RWMutex
	lockGetByNames  sync.RWMutex
	lockGetByNewsId sync.RWMutex
	lockSave        sync.RWMutex
	lockUpdate      sync.RWMutex
}

// Delete calls DeleteFunc.
func (mock *RepositoryMock) Delete(s string) error {
	if mock.DeleteFunc == nil {
		panic("RepositoryMock.DeleteFunc: method is nil but Repository.Delete was just called")
	}
	callInfo := struct {
		S string
	}{
		S: s,
	}
	mock.lockDelete.Lock()
	mock.calls.Delete = append(mock.calls.Delete, callInfo)
	mock.lockDelete.Unlock()
	return mock.DeleteFunc(s)
}

// DeleteCalls gets all the calls that were made to Delete.
// Check the length with:
//     len(mockedRepository.DeleteCalls())
func (mock *RepositoryMock) DeleteCalls() []struct {
	S string
} {
	var calls []struct {
		S string
	}
	mock.lockDelete.RLock()
	calls = mock.calls.Delete
	mock.lockDelete.RUnlock()
	return calls
}

// GetAll calls GetAllFunc.
func (mock *RepositoryMock) GetAll() ([]domain.Tags, error) {
	if mock.GetAllFunc == nil {
		panic("RepositoryMock.GetAllFunc: method is nil but Repository.GetAll was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetAll.Lock()
	mock.calls.GetAll = append(mock.calls.GetAll, callInfo)
	mock.lockGetAll.Unlock()
	return mock.GetAllFunc()
}

// GetAllCalls gets all the calls that were made to GetAll.
// Check the length with:
//     len(mockedRepository.GetAllCalls())
func (mock *RepositoryMock) GetAllCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetAll.RLock()
	calls = mock.calls.GetAll
	mock.lockGetAll.RUnlock()
	return calls
}

// GetById calls GetByIdFunc.
func (mock *RepositoryMock) GetById(s string) (*domain.Tags, error) {
	if mock.GetByIdFunc == nil {
		panic("RepositoryMock.GetByIdFunc: method is nil but Repository.GetById was just called")
	}
	callInfo := struct {
		S string
	}{
		S: s,
	}
	mock.lockGetById.Lock()
	mock.calls.GetById = append(mock.calls.GetById, callInfo)
	mock.lockGetById.Unlock()
	return mock.GetByIdFunc(s)
}

// GetByIdCalls gets all the calls that were made to GetById.
// Check the length with:
//     len(mockedRepository.GetByIdCalls())
func (mock *RepositoryMock) GetByIdCalls() []struct {
	S string
} {
	var calls []struct {
		S string
	}
	mock.lockGetById.RLock()
	calls = mock.calls.GetById
	mock.lockGetById.RUnlock()
	return calls
}

// GetByName calls GetByNameFunc.
func (mock *RepositoryMock) GetByName(names string) (domain.Tags, error) {
	if mock.GetByNameFunc == nil {
		panic("RepositoryMock.GetByNameFunc: method is nil but Repository.GetByName was just called")
	}
	callInfo := struct {
		Names string
	}{
		Names: names,
	}
	mock.lockGetByName.Lock()
	mock.calls.GetByName = append(mock.calls.GetByName, callInfo)
	mock.lockGetByName.Unlock()
	return mock.GetByNameFunc(names)
}

// GetByNameCalls gets all the calls that were made to GetByName.
// Check the length with:
//     len(mockedRepository.GetByNameCalls())
func (mock *RepositoryMock) GetByNameCalls() []struct {
	Names string
} {
	var calls []struct {
		Names string
	}
	mock.lockGetByName.RLock()
	calls = mock.calls.GetByName
	mock.lockGetByName.RUnlock()
	return calls
}

// GetByNames calls GetByNamesFunc.
func (mock *RepositoryMock) GetByNames(names ...string) ([]domain.Tags, error) {
	if mock.GetByNamesFunc == nil {
		panic("RepositoryMock.GetByNamesFunc: method is nil but Repository.GetByNames was just called")
	}
	callInfo := struct {
		Names []string
	}{
		Names: names,
	}
	mock.lockGetByNames.Lock()
	mock.calls.GetByNames = append(mock.calls.GetByNames, callInfo)
	mock.lockGetByNames.Unlock()
	return mock.GetByNamesFunc(names...)
}

// GetByNamesCalls gets all the calls that were made to GetByNames.
// Check the length with:
//     len(mockedRepository.GetByNamesCalls())
func (mock *RepositoryMock) GetByNamesCalls() []struct {
	Names []string
} {
	var calls []struct {
		Names []string
	}
	mock.lockGetByNames.RLock()
	calls = mock.calls.GetByNames
	mock.lockGetByNames.RUnlock()
	return calls
}

// GetByNewsId calls GetByNewsIdFunc.
func (mock *RepositoryMock) GetByNewsId(id string) ([]domain.Tags, error) {
	if mock.GetByNewsIdFunc == nil {
		panic("RepositoryMock.GetByNewsIdFunc: method is nil but Repository.GetByNewsId was just called")
	}
	callInfo := struct {
		ID string
	}{
		ID: id,
	}
	mock.lockGetByNewsId.Lock()
	mock.calls.GetByNewsId = append(mock.calls.GetByNewsId, callInfo)
	mock.lockGetByNewsId.Unlock()
	return mock.GetByNewsIdFunc(id)
}

// GetByNewsIdCalls gets all the calls that were made to GetByNewsId.
// Check the length with:
//     len(mockedRepository.GetByNewsIdCalls())
func (mock *RepositoryMock) GetByNewsIdCalls() []struct {
	ID string
} {
	var calls []struct {
		ID string
	}
	mock.lockGetByNewsId.RLock()
	calls = mock.calls.GetByNewsId
	mock.lockGetByNewsId.RUnlock()
	return calls
}

// Save calls SaveFunc.
func (mock *RepositoryMock) Save(tags domain.Tags) error {
	if mock.SaveFunc == nil {
		panic("RepositoryMock.SaveFunc: method is nil but Repository.Save was just called")
	}
	callInfo := struct {
		Tags domain.Tags
	}{
		Tags: tags,
	}
	mock.lockSave.Lock()
	mock.calls.Save = append(mock.calls.Save, callInfo)
	mock.lockSave.Unlock()
	return mock.SaveFunc(tags)
}

// SaveCalls gets all the calls that were made to Save.
// Check the length with:
//     len(mockedRepository.SaveCalls())
func (mock *RepositoryMock) SaveCalls() []struct {
	Tags domain.Tags
} {
	var calls []struct {
		Tags domain.Tags
	}
	mock.lockSave.RLock()
	calls = mock.calls.Save
	mock.lockSave.RUnlock()
	return calls
}

// Update calls UpdateFunc.
func (mock *RepositoryMock) Update(tags domain.Tags) error {
	if mock.UpdateFunc == nil {
		panic("RepositoryMock.UpdateFunc: method is nil but Repository.Update was just called")
	}
	callInfo := struct {
		Tags domain.Tags
	}{
		Tags: tags,
	}
	mock.lockUpdate.Lock()
	mock.calls.Update = append(mock.calls.Update, callInfo)
	mock.lockUpdate.Unlock()
	return mock.UpdateFunc(tags)
}

// UpdateCalls gets all the calls that were made to Update.
// Check the length with:
//     len(mockedRepository.UpdateCalls())
func (mock *RepositoryMock) UpdateCalls() []struct {
	Tags domain.Tags
} {
	var calls []struct {
		Tags domain.Tags
	}
	mock.lockUpdate.RLock()
	calls = mock.calls.Update
	mock.lockUpdate.RUnlock()
	return calls
}
