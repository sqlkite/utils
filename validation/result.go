package validation

type Result struct {
	len    uint64
	errors []any
	pool   *Pool
}

func NewResult(maxErrors uint16) *Result {
	return &Result{
		errors: make([]any, maxErrors),
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

func (r *Result) InvalidField(field string, meta Meta, data any) {
	r.add(InvalidField{
		Field: field,
		Invalid: Invalid{
			Data:  data,
			Code:  meta.Code,
			Error: meta.Error,
		},
	})
}

func (r *Result) add(error any) {
	l := r.len
	errors := r.errors
	if l < uint64(len(errors)) {
		errors[l] = error
		r.len = l + 1
	}
}

func (r *Result) Release() {
	if pool := r.pool; pool != nil {
		r.len = 0
		pool.list <- r
	}
}
