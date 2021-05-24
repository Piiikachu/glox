package main

import "fmt"

type Obj interface {
	ObjType() ObjType
}

type ObjString struct {
	length int
	str    string
}

type ObjType byte

const (
	OBJ_STRING ObjType = iota
)

func (os *ObjString) ObjType() ObjType {
	return OBJ_STRING
}

func (v Value) ObjType() ObjType {
	return v.asObj().ObjType()
}

func (v Value) isString() bool {
	return v.isObjType(OBJ_STRING)
}

func (v Value) isObjType(t ObjType) bool {
	return v.isType(VAL_OBJ) && v.ObjType() == t
}

func (v Value) asString() ObjString {
	obj := v.asObj()
	objstr := obj.(*ObjString)
	return *objstr
}

func (v Value) asCString() []byte {
	return []byte(v.asString().str)
}

func newObjString(str string) *ObjString {
	return &ObjString{
		length: len(str),
		str:    str,
	}
}

func (v *Value) printObject(){
	switch v.ObjType(){
	case OBJ_STRING:
		fmt.Printf(v.asString().str)
	}
}