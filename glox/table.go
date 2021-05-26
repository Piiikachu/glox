package glox

type Table struct {
	table map[Hash]Value
}

func (t *Table) init() {
	t.table = make(map[Hash]Value)
}

func (t *Table) free() {
	t.table = make(map[Hash]Value)
}

func (t *Table) tableSet(key ObjString, value Value) {
	t.table[key.hash] = value
}

func (t *Table) tableGet(key ObjString) (Value, bool) {
	value, ok := t.table[key.hash]
	return value, ok
}

func (t *Table) tableDelete(key ObjString) {
	delete(t.table, key.hash)
}
func (t *Table) tableFind(str string) bool {
	_, ok := t.table[hashString(str)]
	return ok
}

// type Table struct {
// 	count    int
// 	capacity int
// 	entries  []Entry
// }

// type Entry struct {
// 	key   *ObjString
// 	value Value
// }

// func (t *Table) free() {
// 	t = new(Table)
// }

// func (t *Table) tableSet(key ObjString, value Value) bool{

// 	entry:=t.findEntry(key)
// 	isNewKey:=entry.key==nil
// 	if isNewKey {
// 		t.count++
// 	}
// 	entry.key=&key
// 	entry.value=value
// 	return isNewKey
// }

// func (t *Table) findEntry(key ObjString) Entry{

// }
