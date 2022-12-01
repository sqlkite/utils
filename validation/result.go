package validation

/*
Technically a result can contain a list anything, so long as it
can be serialized to JSON. In reality, we expect it to be a list
of:
  - Invalid,
  - InvalidField, and/or
  - InvalidIndexedField

Invalid is a base and includes an integer code, an english description
and arbitrary meta data (imagine a string with a min and max length,
we'd expect our data to say the min and max).

And InvalidField includes an Invalid and adds a "fields". This
is the name of the field which is invalid. It's  a string array to support nested
objects, e.g. ["name"] or ["user", "name"].

Things go off the rails with InvalidIndexedField. This adds an "indexes" []int
to an InvalidField and is meant to represent the index in an array where validation
fails. The problem with indexes is that they're dynamic, and everything else
in our validation is static. The ["user", "name"] field is created on startup
and never changes. But with array indexes, if we did something like
["user", 23, "name"], the `23` would only be known at runtime. This would require
us to create a copy of fields on each error.

Instead, we keep our fields static and define it as ["user", "#", "name"] and
then add an indexes, which would be: [23]. This works for nested arrays too:
	fields: ["entries", "#", "users", "#", "name"]
	indexes: [2, 23]

It's pretty far from great, but it's efficient and it's something that clients
will be able to handle.
*/

type Result struct {
	len          uint64
	errors       []any
	pool         *Pool
	arrayIndexes []int
	arrayCount   int
}

func NewResult(maxErrors uint16) *Result {
	return &Result{
		arrayCount:   -1,
		errors:       make([]any, maxErrors),
		arrayIndexes: make([]int, 10),
	}
}

func (r Result) Errors() []any {
	return r.errors[:r.len]
}

func (r Result) IsValid() bool {
	return r.len == 0
}

func (r Result) Len() uint64 {
	return r.len
}

func (r *Result) Invalid(meta Meta, data any) {
	r.add(Invalid{
		Data:  data,
		Code:  meta.Code,
		Error: meta.Error,
	})
}

func (r *Result) InvalidField(fields []string, meta Meta, data any) {
	r.addField(InvalidField{
		Fields: fields,
		Invalid: Invalid{
			Data:  data,
			Code:  meta.Code,
			Error: meta.Error,
		},
	})
}

func (r *Result) addField(err InvalidField) {
	arrayCount := r.arrayCount
	if arrayCount == -1 {
		r.add(err)
	} else {
		r.add(InvalidIndexedField{
			InvalidField: err,
			Indexes:      r.arrayIndexes[:arrayCount+1],
		})
	}
}

func (r *Result) add(error any) {
	l := r.len
	errors := r.errors
	if l < uint64(len(errors)) {
		errors[l] = error
		r.len = l + 1
	}
}

func (r *Result) beginArray() {
	r.arrayCount += 1
}

func (r *Result) arrayIndex(i int) {
	r.arrayIndexes[r.arrayCount] = i
}

func (r *Result) endArray() {
	r.arrayCount -= 1
}

func (r *Result) Release() {
	if pool := r.pool; pool != nil {
		r.len = 0
		r.arrayCount = -1
		pool.list <- r
	}
}
