package buffer

import "sync"

/*
A lot of our object pools are encapsulated inside of project
environments. This has a lot of benefits. It simplifies the
code, minimizes contention, and further isolates projects.

But for buffers, which are used both to generate SQL and
read results, having these per-env would be memory-
inefficient. Our buffers need to be relatively large (for
responses), so sharing a large pool across projects is likely
to result in much better usage (so we need far less of them).

Also, these objects have a lifecycle that's different than
an env. Specifically, it has to interact with the http framework
in a way that minimize the amount of copying we need to do.
*/

var buffers = sync.Pool{
	New: func() any {
		// max size doesn't really matter, we're going to
		// reset it to a per-project value on checkout.
		return New(65536, 65536)
	},
}

func Checkout(maxSize int) *Buffer {
	b := buffers.Get().(*Buffer)
	b.max = maxSize
	return b
}

func Release(b *Buffer) {
	b.Reset()
	buffers.Put(b)
}
