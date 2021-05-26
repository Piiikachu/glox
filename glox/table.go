package glox

type Hash uint32

type Hashable interface {
	hashCode() Hash
}

type Table struct {
	table map[Hash]Value
}

func (t *Table) init() {
	t.table = make(map[Hash]Value)
}

func (t *Table) free() {
	t.table = make(map[Hash]Value)
}

func (t *Table) tableSet(key Hashable, value Value) {
	t.table[key.hashCode()] = value
}

func (t *Table) tableGet(key Hashable) (Value, bool) {
	value, ok := t.table[key.hashCode()]
	return value, ok
}

func (t *Table) tableDelete(key Hashable) {
	delete(t.table, key.hashCode())
}
func (t *Table) tableFind(str string) bool {
	_, ok := t.table[hashString(str)]
	return ok
}
