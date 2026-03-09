module github.com/providentia/digna/supply

go 1.25.4

require github.com/providentia/digna/lifecycle v0.0.0

require (
	github.com/mattn/go-sqlite3 v1.14.34 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
)

replace github.com/providentia/digna/lifecycle => ../lifecycle
