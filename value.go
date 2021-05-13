package main

type Value float64

type ValueArray struct {
	capacity int
	count    int
	values   []Value
}

func (array *ValueArray) write(value Value) {
	array.values = append(array.values, value)
	array.count = len(array.values)
	array.capacity = cap(array.values)
}

func (array *ValueArray) free() {
	array = new(ValueArray)
}
