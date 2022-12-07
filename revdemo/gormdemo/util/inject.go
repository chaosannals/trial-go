package util

import (
	"reflect"
)

type Injectable interface {
	GetType() reflect.Type
	GetValue() reflect.Value
	GetPointerValue() reflect.Value
}
