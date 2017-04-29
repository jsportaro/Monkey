package object

import (
	"hash/fnv"
)

type Hashable interface {
	HashKey() HashKey
}

//HashKey hash
type HashKey struct {
	Type  ObjectType
	Value uint64
}

//HashKey booleans
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

//HashKey ints
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

//HashKey strings
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}
