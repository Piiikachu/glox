package main

import "fmt"

type Value struct {
	valueType ValueType
	val       interface{}
}

type ValueArray struct {
	capacity int
	count    int
	values   []Value
}

type ValueType byte

const (
	VAL_BOOL ValueType = iota
	VAL_NIL
	VAL_NUMBER
	VAL_OBJ
)

func BOOL_VAL(value bool) Value {
	return Value{
		valueType: VAL_BOOL,
		val:       value,
	}

}
func NUMBER_VAL(value float64) Value {
	return Value{
		valueType: VAL_NUMBER,
		val:       value,
	}
}

func NIL_VAL() Value {
	return Value{
		valueType: VAL_NIL,
		val:       nil,
	}
}

func OBJ_VAL(value Obj) Value {
	return Value{
		valueType: VAL_OBJ,
		val:       value,
	}
}

func (value Value) asBool() bool {
	return value.val.(bool)
}

func (value Value) asNumber() float64 {
	return value.val.(float64)
}

func (value Value) asObj() Obj {
	return value.val.(Obj)
}

func (value Value) isType(valType ValueType) bool {
	return value.valueType == valType
}

func (v1 Value) equals(v2 Value) bool {
	if v1.valueType != v2.valueType {
		return false
	}
	switch v1.valueType {
	case VAL_BOOL:
		return v1.asBool() == v2.asBool()
	case VAL_NIL:
		return true
	case VAL_NUMBER:
		return v1.asNumber() == v2.asNumber()
	case VAL_OBJ:
		{
			a := v1.asString()
			b := v2.asString()
			return a == b
		}
	default:
		return false
	}
}

func (array *ValueArray) write(value Value) {
	array.values = append(array.values, value)
	array.count = len(array.values)
	array.capacity = cap(array.values)
}

func (array *ValueArray) free() {
	array = new(ValueArray)
}

func printValue(value Value) {
	switch value.valueType {
	case VAL_BOOL:
		if value.asBool() {
			fmt.Printf("true")
		} else {
			fmt.Printf("false")
		}
	case VAL_NIL:
		fmt.Printf("nil")
	case VAL_NUMBER:
		fmt.Printf("%g", value.asNumber())
	case VAL_OBJ:
		value.printObject()
	}
}
