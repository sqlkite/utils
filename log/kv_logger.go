package log

import (
	"io"
	"strconv"
	"strings"
	"time"
)

type KvLogger struct {
	out io.Writer

	// the position in buffer to write to next
	pos uint64

	// reference back into our pool
	pool *Pool

	// buffer that we write our message to
	buffer []byte

	// A logger can have a fixed piece of data which is
	// always included (e.g pid=$PROJECT_ID for a project-owned
	// logger). Once our fixed data is set, pos will never be
	// less than fixedLen.
	fixedLen uint64

	// A logger can also have temporary repeated data
	// (e.g. rid=$REQUEST_ID for an env-owned logger).
	// After logging a message, pos == multiUseLen. Only
	// on reset/release will pos == fixedLen
	multiUseLen uint64
}

func KvFactory(maxSize uint32, out io.Writer) Factory {
	return func(pool *Pool) Logger {
		return &KvLogger{
			out:    out,
			pool:   pool,
			buffer: make([]byte, maxSize),
		}
	}
}

// Get the bytes from the logger. This is only valid before Log is called (after
// log is called, you'll get an empty slice). Only really useful for testing.
func (l *KvLogger) Bytes() []byte {
	return l.buffer[:l.pos]
}

// Logger will _always_ include this data. Meant to be used with the Field builder.
// Even once released to the pool and re-checked out, this data will still be in the logger.
// For checkout-specific data, see MultiUse().
func (l *KvLogger) Fixed() {
	l.fixedLen = l.pos
}

// Similar to Fixed, but exists only while checked out
func (l *KvLogger) MultiUse() Logger {
	l.multiUseLen = l.pos
	return l
}

// Add a field (key=value) where value is a string
func (l *KvLogger) String(key string, value string) Logger {
	l.writeKeyValue(key, value, false)
	return l
}

// Add a field (key=value) where value is an int
func (l *KvLogger) Int(key string, value int) Logger {
	return l.Int64(key, int64(value))
}

// Add a field (key=value) where value is an int
func (l *KvLogger) Int64(key string, value int64) Logger {
	l.writeKeyValue(key, strconv.FormatInt(value, 10), true)
	return l
}

// Add a field (key=value) where value is an error
func (l *KvLogger) Err(err error) Logger {
	se, ok := err.(*StructuredError)
	if !ok {
		return l.String("err", err.Error())
	}

	l.Int("code", se.Code).String("err", se.Err.Error())
	for key, value := range se.Data {
		switch v := value.(type) {
		case string:
			l.String(key, v)
		case int:
			l.Int(key, v)
		}
	}
	return l
}

// Write the log to our globally configured writer
func (l *KvLogger) Log() {
	l.LogTo(l.out)
}

func (l *KvLogger) LogTo(out io.Writer) {
	pos := l.pos
	buffer := l.buffer

	// no length check, if we did everything right, there should
	// always be at least 1 space in our buffer
	buffer[pos] = '\n'
	out.Write(buffer[:pos+1])

	if l.multiUseLen == 0 {
		l.Release()
	}
}

func (l *KvLogger) Reset() {
	l.pos = l.fixedLen
}

func (l *KvLogger) Release() {
	l.pos = l.fixedLen // Reset()
	if pool := l.pool; pool != nil {
		pool.list <- l
	}
}

// Log an info-level message. Every message must have a [hopefully] unique context
func (l *KvLogger) Info(ctx string) Logger {
	return l.start(ctx, []byte("l=info t="))
}

// Log an warn-level message. Every message must have a [hopefully] unique context
func (l *KvLogger) Warn(ctx string) Logger {
	return l.start(ctx, []byte("l=warn t="))
}

// Log an error-level message. Every message must have a [hopefully] unique context
func (l *KvLogger) Error(ctx string) Logger {
	return l.start(ctx, []byte("l=error t="))
}

// Log an fatal-level message. Every message must have a [hopefully] unique context
func (l *KvLogger) Fatal(ctx string) Logger {
	return l.start(ctx, []byte("l=fatal t="))
}

func (l *KvLogger) Field(field Field) Logger {
	pos := l.pos
	buffer := l.buffer
	bl := uint64(len(buffer))

	// might already have data
	if pos != 0 && pos < bl {
		buffer[pos] = ' '
		pos += 1
	}

	if pos < uint64(len(buffer)) {
		data := field.kv
		copy(buffer[pos:], data)
		l.pos = pos + uint64(len(data))
	}
	return l
}

// "starts" a new log message. Every message always contains a timestamp (t) a
// context (c) and a level (l).
func (l *KvLogger) start(ctx string, meta []byte) Logger {
	pos := l.pos
	buffer := l.buffer

	bl := uint64(len(buffer))

	// pos > 0 when MultiUse is enabled
	if pos > 0 && pos < bl {
		buffer[pos] = ' '
		pos = pos + 1
	}

	copy(buffer[pos:], meta)
	pos += uint64(len(meta))

	t := strconv.FormatInt(time.Now().Unix(), 10)
	copy(buffer[pos:], t)
	pos += uint64(len(t))

	// we always expect the ctx to be safe and to outlive this log
	copy(buffer[pos:], []byte(" c="))
	pos += 3

	copy(buffer[pos:], ctx)
	pos += uint64(len(ctx))

	l.pos = pos
	return l
}

// When safe, we're being told that value 100% does not need
// to be escaped (e.g. we know the value is an int), so we don't need to check/encode it.
func (l *KvLogger) writeKeyValue(key string, value string, safe bool) {
	l.pos = writeKeyValue(key, value, safe, l.pos, l.buffer)
}

// We expect key to always be safe to write as-is.
// We only encode newline and quotes. If either is present, the value is quote encoded.
func writeKeyValue(key string, value string, safe bool, pos uint64, buffer []byte) uint64 {
	bl := uint64(len(buffer))

	// Need at least enough room for:
	// space sperator + equal separator + trailing newline
	// + our key + our value
	if bl-pos < uint64(len(key)+len(value))+3 {
		return pos
	}

	if pos > 0 {
		buffer[pos] = ' '
		pos += 1
	}

	copy(buffer[pos:], key)
	pos += uint64(len(key))

	buffer[pos] = '='
	pos += 1

	if safe || !strings.ContainsAny(value, " =\"\n") {
		copy(buffer[pos:], value)
		return pos + uint64(len(value))
	}

	buffer[pos] = '"'
	pos += 1

	// -2 because we need enough space for our quote and final newline
	var i int
	for ; i < len(value) && pos < bl-5; i++ {
		c := value[i]
		switch c {
		case '\n':
			buffer[pos] = '\\'
			buffer[pos+1] = 'n'
			pos += 2
		case '"':
			buffer[pos] = '\\'
			buffer[pos+1] = '"'
			pos += 2
		default:
			buffer[pos] = c
			pos += 1
		}
	}

	if pos == bl-5 && i < len(value) {
		copy(buffer[pos:], "...")
		pos += 3
	}

	buffer[pos] = '"'
	return pos + 1
}
