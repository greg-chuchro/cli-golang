package class_member_stores

import (
	"github.com/greg-chuchro/cli-golang/data_access"
)

// FieldStore provides key-value access to object fields by name.
type FieldStore struct {
	*FieldStoreBase
}

func NewFieldStore(obj any, bindingFlags BindingFlags) *FieldStore {
	return &FieldStore{FieldStoreBase: NewFieldStoreBase(obj, bindingFlags, nil)}
}

var _ data_access.IKeyValueStoreMediated[string, any, string, any] = (*FieldStore)(nil)
