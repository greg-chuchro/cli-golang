package class_member_stores

import (
	"reflect"

	"github.com/greg-chuchro/cli-golang/data_access"
)

// ParameterStore provides key-value access to method parameters by name.
type ParameterStore struct {
	*ParameterStoreBase
}

func NewParameterStore(method reflect.Value) *ParameterStore {
	return &ParameterStore{ParameterStoreBase: NewParameterStoreBase(method)}
}

var _ data_access.IKeyValueStoreMediated[string, any, int, any] = (*ParameterStore)(nil)
