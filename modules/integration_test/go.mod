module github.com/providentia/digna/integration_test

go 1.25.4

require (
	github.com/providentia/digna/cash_flow v0.0.0
	github.com/providentia/digna/core_lume v0.0.0
	github.com/providentia/digna/legal_facade v0.0.0
	github.com/providentia/digna/lifecycle v0.0.0
	github.com/providentia/digna/pdv_ui v0.0.0
	github.com/providentia/digna/reporting v0.0.0
	github.com/providentia/digna/sync_engine v0.0.0
)

require github.com/mattn/go-sqlite3 v1.14.34 // indirect

replace (
	github.com/providentia/digna/cash_flow => ../cash_flow
	github.com/providentia/digna/core_lume => ../core_lume
	github.com/providentia/digna/legal_facade => ../legal_facade
	github.com/providentia/digna/lifecycle => ../lifecycle
	github.com/providentia/digna/pdv_ui => ../pdv_ui
	github.com/providentia/digna/reporting => ../reporting
	github.com/providentia/digna/sync_engine => ../sync_engine
)
