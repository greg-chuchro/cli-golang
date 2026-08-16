package data_access

import (
	"reflect"

	"github.com/greg-chuchro/cli-golang/data_access/exceptions"
)

// IReadOnlyKeyValueStore provides read-only access to a key-value store for retrieving values from class members.
type IReadOnlyKeyValueStore[TKey comparable, TValue any] interface {
	Get(key TKey) TValue
	ContainsKey(key TKey) bool
	GetValueType(key TKey) reflect.Type
}

// IKeyValueStore provides read-write access to a key-value store for accessing class members.
type IKeyValueStore[TKey comparable, TValue any] interface {
	IReadOnlyKeyValueStore[TKey, TValue]
	Set(key TKey, value TValue)
	GetValueOrInitialize(key TKey) TValue
}

// IMediatedKeyValueStore provides mediated access to a key-value store through accessors with internal value conversion.
type IMediatedKeyValueStore[TKey comparable, TValue any, TAccessor comparable, TInternalValue any] interface {
	Accessors() []TAccessor
	GetByAccessor(accessor TAccessor) TInternalValue
	SetByAccessor(accessor TAccessor, value TInternalValue)
	ContainsAccessor(accessor TAccessor) bool
	GetAccessor(key TKey) TAccessor
	GetValueTypeByAccessor(accessor TAccessor) reflect.Type
	GetValueOrInitializeByAccessor(accessor TAccessor) TInternalValue
	TryGetAccessor(key TKey) (TAccessor, bool)
	ConvertFromInternalValue(value TInternalValue, accessor TAccessor) TValue
	ConvertToInternalValue(value TValue, accessor TAccessor) TInternalValue
}

// IKeyValueStoreMediated provides key-value store with mediated access through accessors and internal value conversion.
type IKeyValueStoreMediated[TKey comparable, TValue any, TAccessor comparable, TInternalValue any] interface {
	IKeyValueStore[TKey, TValue]
	IMediatedKeyValueStore[TKey, TValue, TAccessor, TInternalValue]
}

// IKeyValueStoreMediatedSimple provides key-value store with mediated access where internal value type matches external value type.
type IKeyValueStoreMediatedSimple[TKey comparable, TValue any, TAccessor comparable] interface {
	IKeyValueStoreMediated[TKey, TValue, TAccessor, TValue]
}

// IFunctionExecutor provides a key-value store that can execute functions (methods) with parameters.
type IFunctionExecutor[TParameterKey comparable, TParameterValue any] interface {
	IKeyValueStore[TParameterKey, TParameterValue]
	AddArgument(value TParameterValue)
	Invoke() (any, error)
}

// KeyValueStoreBase implements the mediated key-value store with internal value conversion.
// It is embedded by concrete stores and delegates abstract operations to the mediated interface.
type KeyValueStoreBase[TKey comparable, TValue any, TAccessor comparable, TInternalValue any] struct {
	store IKeyValueStoreMediated[TKey, TValue, TAccessor, TInternalValue]
}

func NewKeyValueStoreBase[TKey comparable, TValue any, TAccessor comparable, TInternalValue any](store IKeyValueStoreMediated[TKey, TValue, TAccessor, TInternalValue]) *KeyValueStoreBase[TKey, TValue, TAccessor, TInternalValue] {
	return &KeyValueStoreBase[TKey, TValue, TAccessor, TInternalValue]{store: store}
}

func (this *KeyValueStoreBase[TKey, TValue, TAccessor, TInternalValue]) Accessors() []TAccessor {
	return this.store.Accessors()
}

func (this *KeyValueStoreBase[TKey, TValue, TAccessor, TInternalValue]) Get(key TKey) TValue {
	accessor := this.store.GetAccessor(key)
	return this.store.ConvertFromInternalValue(this.store.GetByAccessor(accessor), accessor)
}

func (this *KeyValueStoreBase[TKey, TValue, TAccessor, TInternalValue]) Set(key TKey, value TValue) {
	accessor := this.store.GetAccessor(key)
	this.store.SetByAccessor(accessor, this.store.ConvertToInternalValue(value, accessor))
}

func (this *KeyValueStoreBase[TKey, TValue, TAccessor, TInternalValue]) GetByAccessor(accessor TAccessor) TInternalValue {
	return this.store.GetByAccessor(accessor)
}

func (this *KeyValueStoreBase[TKey, TValue, TAccessor, TInternalValue]) SetByAccessor(accessor TAccessor, value TInternalValue) {
	this.store.SetByAccessor(accessor, value)
}

func (this *KeyValueStoreBase[TKey, TValue, TAccessor, TInternalValue]) ContainsAccessor(accessor TAccessor) bool {
	return this.store.ContainsAccessor(accessor)
}

func (this *KeyValueStoreBase[TKey, TValue, TAccessor, TInternalValue]) ContainsKey(key TKey) bool {
	_, ok := this.store.TryGetAccessor(key)
	return ok
}

func (this *KeyValueStoreBase[TKey, TValue, TAccessor, TInternalValue]) GetAccessor(key TKey) TAccessor {
	accessor, ok := this.store.TryGetAccessor(key)
	if !ok {
		panic(exceptions.KeyNotFound(key))
	}
	return accessor
}

func (this *KeyValueStoreBase[TKey, TValue, TAccessor, TInternalValue]) GetValueType(key TKey) reflect.Type {
	accessor := this.store.GetAccessor(key)
	return this.store.GetValueTypeByAccessor(accessor)
}

func (this *KeyValueStoreBase[TKey, TValue, TAccessor, TInternalValue]) GetValueTypeByAccessor(accessor TAccessor) reflect.Type {
	return this.store.GetValueTypeByAccessor(accessor)
}

func (this *KeyValueStoreBase[TKey, TValue, TAccessor, TInternalValue]) GetValueOrInitialize(key TKey) TValue {
	accessor := this.store.GetAccessor(key)
	return this.store.ConvertFromInternalValue(this.store.GetValueOrInitializeByAccessor(accessor), accessor)
}

func (this *KeyValueStoreBase[TKey, TValue, TAccessor, TInternalValue]) GetValueOrInitializeByAccessor(accessor TAccessor) TInternalValue {
	return this.store.GetValueOrInitializeByAccessor(accessor)
}

func (this *KeyValueStoreBase[TKey, TValue, TAccessor, TInternalValue]) TryGetAccessor(key TKey) (TAccessor, bool) {
	return this.store.TryGetAccessor(key)
}

func (this *KeyValueStoreBase[TKey, TValue, TAccessor, TInternalValue]) ConvertFromInternalValue(value TInternalValue, accessor TAccessor) TValue {
	return this.store.ConvertFromInternalValue(value, accessor)
}

func (this *KeyValueStoreBase[TKey, TValue, TAccessor, TInternalValue]) ConvertToInternalValue(value TValue, accessor TAccessor) TInternalValue {
	return this.store.ConvertToInternalValue(value, accessor)
}

// DualKeyValueStoreBase combines two key-value stores with primary and secondary lookup priority.
type DualKeyValueStoreBase[TKey comparable, TValue any] struct {
	primaryStore   func() IKeyValueStore[TKey, TValue]
	secondaryStore func() IKeyValueStore[TKey, TValue]
}

func (this *DualKeyValueStoreBase[TKey, TValue]) PrimaryStore() IKeyValueStore[TKey, TValue] {
	return this.primaryStore()
}

func (this *DualKeyValueStoreBase[TKey, TValue]) SecondaryStore() IKeyValueStore[TKey, TValue] {
	return this.secondaryStore()
}

func (this *DualKeyValueStoreBase[TKey, TValue]) Get(key TKey) TValue {
	if this.primaryStore().ContainsKey(key) {
		return this.primaryStore().Get(key)
	}
	return this.secondaryStore().Get(key)
}

func (this *DualKeyValueStoreBase[TKey, TValue]) Set(key TKey, value TValue) {
	if this.primaryStore().ContainsKey(key) {
		this.primaryStore().Set(key, value)
	} else {
		this.secondaryStore().Set(key, value)
	}
}

func (this *DualKeyValueStoreBase[TKey, TValue]) ContainsKey(key TKey) bool {
	if this.primaryStore().ContainsKey(key) {
		return this.primaryStore().ContainsKey(key)
	}
	return this.secondaryStore().ContainsKey(key)
}

func (this *DualKeyValueStoreBase[TKey, TValue]) GetValueType(key TKey) reflect.Type {
	if this.primaryStore().ContainsKey(key) {
		return this.primaryStore().GetValueType(key)
	}
	return this.secondaryStore().GetValueType(key)
}

func (this *DualKeyValueStoreBase[TKey, TValue]) GetValueOrInitialize(key TKey) TValue {
	if this.primaryStore().ContainsKey(key) {
		return this.primaryStore().GetValueOrInitialize(key)
	}
	return this.secondaryStore().GetValueOrInitialize(key)
}

// DualKeyValueStore combines two key-value stores with primary and secondary lookup priority.
type DualKeyValueStore[TKey comparable, TValue any] struct {
	DualKeyValueStoreBase[TKey, TValue]
	primary   IKeyValueStore[TKey, TValue]
	secondary IKeyValueStore[TKey, TValue]
}

func NewDualKeyValueStore[TKey comparable, TValue any](primaryStore IKeyValueStore[TKey, TValue], secondaryStore IKeyValueStore[TKey, TValue]) *DualKeyValueStore[TKey, TValue] {
	s := &DualKeyValueStore[TKey, TValue]{primary: primaryStore, secondary: secondaryStore}
	s.DualKeyValueStoreBase = DualKeyValueStoreBase[TKey, TValue]{
		primaryStore:   func() IKeyValueStore[TKey, TValue] { return s.primary },
		secondaryStore: func() IKeyValueStore[TKey, TValue] { return s.secondary },
	}
	return s
}

func (this *DualKeyValueStore[TKey, TValue]) PrimaryStore() IKeyValueStore[TKey, TValue] {
	return this.primary
}

func (this *DualKeyValueStore[TKey, TValue]) SecondaryStore() IKeyValueStore[TKey, TValue] {
	return this.secondary
}

// DistinctDualKeyValueStoreBase combines two key-value stores and throws on ambiguous keys when shadowing is disabled.
type DistinctDualKeyValueStoreBase[TKey comparable, TValue any] struct {
	primary         IKeyValueStore[TKey, TValue]
	secondary       IKeyValueStore[TKey, TValue]
	EnableShadowing bool
}

func NewDistinctDualKeyValueStoreBase[TKey comparable, TValue any](primaryStore IKeyValueStore[TKey, TValue], secondaryStore IKeyValueStore[TKey, TValue]) *DistinctDualKeyValueStoreBase[TKey, TValue] {
	return &DistinctDualKeyValueStoreBase[TKey, TValue]{primary: primaryStore, secondary: secondaryStore}
}

func (this *DistinctDualKeyValueStoreBase[TKey, TValue]) assertNoCollision(key TKey) {
	if !this.EnableShadowing {
		if this.primary.ContainsKey(key) && this.secondary.ContainsKey(key) {
			panic(exceptions.AmbiguousKey(key))
		}
	}
}

func (this *DistinctDualKeyValueStoreBase[TKey, TValue]) Get(key TKey) TValue {
	this.assertNoCollision(key)
	if this.primary.ContainsKey(key) {
		return this.primary.Get(key)
	}
	return this.secondary.Get(key)
}

func (this *DistinctDualKeyValueStoreBase[TKey, TValue]) Set(key TKey, value TValue) {
	this.assertNoCollision(key)
	if this.primary.ContainsKey(key) {
		this.primary.Set(key, value)
	} else {
		this.secondary.Set(key, value)
	}
}

func (this *DistinctDualKeyValueStoreBase[TKey, TValue]) ContainsKey(key TKey) bool {
	this.assertNoCollision(key)
	if this.primary.ContainsKey(key) {
		return this.primary.ContainsKey(key)
	}
	return this.secondary.ContainsKey(key)
}

func (this *DistinctDualKeyValueStoreBase[TKey, TValue]) GetValueType(key TKey) reflect.Type {
	this.assertNoCollision(key)
	if this.primary.ContainsKey(key) {
		return this.primary.GetValueType(key)
	}
	return this.secondary.GetValueType(key)
}

func (this *DistinctDualKeyValueStoreBase[TKey, TValue]) GetValueOrInitialize(key TKey) TValue {
	this.assertNoCollision(key)
	if this.primary.ContainsKey(key) {
		return this.primary.GetValueOrInitialize(key)
	}
	return this.secondary.GetValueOrInitialize(key)
}
