package glox

import "fmt"

type Hash uint32

type IObj interface {
	ObjType() ObjType
	next() IObj
	free()
	getHash() Hash
}

type Obj struct {
	nextObj IObj
	hash    Hash
}

type ObjString struct {
	Obj
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

func (obj *Obj) next() IObj {
	return obj.nextObj
}

func (obj *Obj) getHash() Hash {
	return obj.hash
}

func (os *ObjString) free() {
	os.str = ""
	os.length = 0
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
	obj := &ObjString{
		length: len(str),
		str:    str,
	}
	obj.nextObj = vm.objects
	obj.hash = hashString(str)
	vm.objects = obj

	vm.strings.tableSet(*obj,NIL_VAL())

	return obj
}

func hashString(str string) Hash {
	hash := 2166136261
	for _, b := range []byte(str) {
		hash ^= int(b)
		hash *= 16777619
	}
	return Hash(hash)
}

func (v *Value) printObject() {
	switch v.ObjType() {
	case OBJ_STRING:
		fmt.Printf(v.asString().str)
	}
}
