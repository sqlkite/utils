module src.sqlkite.com/utils

go 1.19

replace src.sqlkite.com/sqlite => ../sqlite

require (
	github.com/goccy/go-json v0.9.11
	github.com/google/uuid v1.3.0
	github.com/jackc/pgx/v5 v5.0.4
	github.com/valyala/fasthttp v1.41.0
	golang.org/x/crypto v0.1.0
	golang.org/x/sync v0.1.0
	src.sqlkite.com/sqlite v0.0.4
	src.sqlkite.com/tests v0.0.0-20221128084111-4f87425f94ef
)

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/puddle/v2 v2.0.0 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/text v0.4.0 // indirect
)
