package class_member_stores

import (
	"reflect"

	"github.com/greg-chuchro/cli-golang/cli/data_access"
	"github.com/greg-chuchro/cli-golang/cli/exceptions"
	cms "github.com/greg-chuchro/cli-golang/data_access/class_member_stores"
)

// MethodInfoStore provides read-only access to CLI subcommands (methods) available on an object.
type MethodInfoStore struct {
	targetType      reflect.Type
	accessorsByName map[string]reflect.Method
	stringifier     data_access.IClassMemberStringifier
	accessValidator data_access.IAccessValidator
}

func NewMethodInfoStore(targetObject any, bindingFlags cms.BindingFlags, stringifier data_access.IClassMemberStringifier, accessValidator data_access.IAccessValidator) *MethodInfoStore {
	t := reflect.TypeOf(targetObject)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	methodType := t
	if methodType.Kind() != reflect.Ptr {
		ptrType := reflect.PointerTo(methodType)
		if ptrType.NumMethod() > 0 {
			methodType = ptrType
		}
	}
	s := &MethodInfoStore{
		targetType:      methodType,
		stringifier:     stringifier,
		accessValidator: accessValidator,
	}
	s.accessorsByName = s.buildAccessorsByName()
	return s
}

func (this *MethodInfoStore) buildAccessorsByName() map[string]reflect.Method {
	all := map[string]reflect.Method{}
	for i := 0; i < this.targetType.NumMethod(); i++ {
		m := this.targetType.Method(i)
		if !this.isValidMethod(m) {
			continue
		}
		all[m.Name] = m
	}

	byRequired := map[string]reflect.Method{}
	for _, m := range all {
		stField := reflect.StructField{Name: m.Name, Type: m.Type}
		for _, key := range this.stringifier.GetRequiredNames(stField) {
			if _, exists := byRequired[key]; !exists {
				byRequired[key] = m
			} else {
				panic(exceptions.CliErrors.NameCollision(m.Name, byRequired[key].Name))
			}
		}
	}

	byAlt := map[string]reflect.Method{}
	collisions := map[string]bool{}
	for _, m := range all {
		stField := reflect.StructField{Name: m.Name, Type: m.Type}
		for _, key := range this.stringifier.GetAlternativeNames(stField) {
			if _, exists := byAlt[key]; exists {
				collisions[key] = true
			} else {
				byAlt[key] = m
			}
		}
	}
	for key := range collisions {
		delete(byAlt, key)
	}

	for key, v := range byAlt {
		if _, exists := byRequired[key]; !exists {
			byRequired[key] = v
		}
	}
	return byRequired
}

func (this *MethodInfoStore) isValidMethod(m reflect.Method) bool {
	if m.Name == "" || m.Name[0] >= 'a' && m.Name[0] <= 'z' {
		return false
	}
	return true
}

func (this *MethodInfoStore) Get(key string) reflect.Method {
	m, ok := this.accessorsByName[key]
	if !ok {
		panic(exceptions.CliErrors.InvalidSubcommand(key))
	}
	return m
}

func (this *MethodInfoStore) ContainsKey(key string) bool {
	_, ok := this.accessorsByName[key]
	return ok
}

func (this *MethodInfoStore) GetValueType(key string) reflect.Type {
	m := this.Get(key)
	if m.Type.NumOut() == 0 {
		return nil
	}
	return m.Type.Out(0)
}

func (this *MethodInfoStore) GetAccessorKeysPairs() []data_access.AccessorKeysPair[reflect.StructField] {
	// Group by distinct method: accessorsByName maps each name (including aliases) to a
	// method, so iterate the methods and collect all their names as the key set.
	byMethod := map[reflect.Method][]string{}
	for name, m := range this.accessorsByName {
		byMethod[m] = append(byMethod[m], name)
	}
	result := make([]data_access.AccessorKeysPair[reflect.StructField], 0, len(byMethod))
	for m, keys := range byMethod {
		stField := reflect.StructField{Name: m.Name, Type: m.Type}
		result = append(result, *data_access.NewAccessorKeysPair(stField, keys))
	}
	return result
}
