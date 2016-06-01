// Copyright (C) 2015 The Protocol Authors.

package protocol

// The Vector type represents a version vector. The zero value is a usable
// version vector. The vector has slice semantics and some operations on it
// are "append-like" in that they may return the same vector modified, or v
// new allocated Vector with the modified contents.

// Counter represents a single counter in the version vector.

// Update returns a Vector with the index for the specific ID incremented by
// one. If it is possible, the vector v is updated and returned. If it is not,
// a copy will be created, updated and returned.
func (v Vector) Update(id ShortID) Vector {
	for i := range v.Counters {
		if v.Counters[i].ID == uint64(id) {
			// Update an existing index
			v.Counters[i].Value++
			return v
		} else if v.Counters[i].ID > uint64(id) {
			// Insert a new index
			nv := make([]Counter, len(v.Counters)+1)
			copy(nv, v.Counters[:i])
			nv[i].ID = uint64(id)
			nv[i].Value = 1
			copy(nv[i+1:], v.Counters[i:])
			return Vector{nv}
		}
	}
	// Append a new index
	return Vector{append(v.Counters, Counter{uint64(id), 1})}
}

// Merge returns the vector containing the maximum indexes from v and b. If it
// is possible, the vector v is updated and returned. If it is not, a copy
// will be created, updated and returned.
func (v Vector) Merge(b Vector) Vector {
	var vi, bi int
	for bi < len(b.Counters) {
		if vi == len(v.Counters) {
			// We've reach the end of v, all that remains are appends
			return Vector{append(v.Counters, b.Counters[bi:]...)}
		}

		if v.Counters[vi].ID > b.Counters[bi].ID {
			// The index from b should be inserted here
			n := make([]Counter, len(v.Counters)+1)
			copy(n, v.Counters[:vi])
			n[vi] = b.Counters[bi]
			copy(n[vi+1:], v.Counters[vi:])
			v.Counters = n
		}

		if v.Counters[vi].ID == b.Counters[bi].ID {
			if val := b.Counters[bi].Value; val > v.Counters[vi].Value {
				v.Counters[vi].Value = val
			}
		}

		if bi < len(b.Counters) && v.Counters[vi].ID == b.Counters[bi].ID {
			bi++
		}
		vi++
	}

	return v
}

// Copy returns an identical vector that is not shared with v.
func (v Vector) Copy() Vector {
	nv := make([]Counter, len(v.Counters))
	copy(nv, v.Counters)
	return Vector{nv}
}

// Equal returns true when the two vectors are equivalent.
func (v Vector) Equal(b Vector) bool {
	return v.Compare(b) == Equal
}

// LesserEqual returns true when the two vectors are equivalent or v is Lesser
// than b.
func (v Vector) LesserEqual(b Vector) bool {
	comp := v.Compare(b)
	return comp == Lesser || comp == Equal
}

// GreaterEqual returns true when the two vectors are equivalent or v is Greater
// than b.
func (v Vector) GreaterEqual(b Vector) bool {
	comp := v.Compare(b)
	return comp == Greater || comp == Equal
}

// Concurrent returns true when the two vectors are concurrent.
func (v Vector) Concurrent(b Vector) bool {
	comp := v.Compare(b)
	return comp == ConcurrentGreater || comp == ConcurrentLesser
}

// Counter returns the current value of the given counter ID.
func (v Vector) Counter(id ShortID) uint64 {
	for _, c := range v.Counters {
		if c.ID == uint64(id) {
			return c.Value
		}
	}
	return 0
}
