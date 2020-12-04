package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
)

// HashKey represents a hash key object and holds the Type of Object
// hashed and its hash value in Value
type HashKey struct {
	Type  Type
	Value uint64
}

// HashKey returns a HashKey object
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

// HashKey returns a HashKey object
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// HashKey returns a HashKey object
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// HashKey returns a hash key for the given object.
func (f *Float) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(f.Inspect()))
	return HashKey{Type: f.Type(), Value: h.Sum64()}
}

// HashPair is an object that holds a key and value of type Object
type HashPair struct {
	Key   Object
	Value Object
}

// Hash is a hash map and holds a map of HashKey to HashPair(s)
type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Len() int {
	return len(h.Pairs)
}

func (h *Hash) Bool() bool {
	return len(h.Pairs) > 0
}

func (h *Hash) Compare(other Object) int {
	if obj, ok := other.(*Hash); ok {
		if len(h.Pairs) != len(obj.Pairs) {
			return -1
		}
		for _, pair := range h.Pairs {
			left := pair.Value
			hashed := left.(Hashable)
			right, ok := obj.Pairs[hashed.HashKey()]
			if !ok {
				return -1
			}
			cmp, ok := left.(Comparable)
			if !ok {
				return -1
			}
			if cmp.Compare(right.Value) != 0 {
				return cmp.Compare(right.Value)
			}
		}

		return 0
	}
	return -1
}

func (h *Hash) String() string {
	return h.Inspect()
}

// Type returns the type of the object
func (h *Hash) Type() Type { return HASH }

// Inspect returns a stringified version of the object for debugging
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// ToInterface converts this object to a go-interface, which will allow
// it to be used naturally in our sprintf/printf primitives.
//
// It might also be helpful for embedded users.
func (h *Hash) ToInterface() interface{} {
	return "<HASH>"
}
