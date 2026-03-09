module github.com/providentia/digna/budget

go 1.25.4

require (
	github.com/providentia/digna/cash_flow v0.0.0-00010101000000-000000000000
	github.com/providentia/digna/lifecycle v0.0.0
)

require (
	github.com/mattn/go-sqlite3 v1.14.34 // indirect
	github.com/providentia/digna/core_lume v0.0.0-00010101000000-000000000000 // indirect
)

replace github.com/providentia/digna/lifecycle => ../lifecycle

replace github.com/providentia/digna/cash_flow => ../cash_flow

replace github.com/providentia/digna/core_lume => ../core_lume
