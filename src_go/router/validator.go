package router

import (
	"reflect"
	"sync"

	"github.com/go-playground/validator/v10"
)

type ginValidator struct {
	once     sync.Once
	validate *validator.Validate
}

func (v *ginValidator) ValidateStruct(obj interface{}) error {
	value := reflect.ValueOf(obj)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	if valueType == reflect.Struct {
		v.lazyinit()
		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}

func (v *ginValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *ginValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")
	})
}
