module github.com/providentia/digna/ui_web

go 1.25.4

require (
	github.com/providentia/digna/budget v0.0.0
	github.com/providentia/digna/cash_flow v0.0.0
	github.com/providentia/digna/lifecycle v0.0.0
	github.com/providentia/digna/pdv_ui v0.0.0
	github.com/providentia/digna/reporting v0.0.0
	github.com/providentia/digna/supply v0.0.0
)

require (
	github.com/deckarep/golang-set/v2 v2.8.0 // indirect
	github.com/go-jose/go-jose/v3 v3.0.4 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.34 // indirect
	github.com/playwright-community/playwright-go v0.5700.1 // indirect
	github.com/providentia/digna/core_lume v0.0.0 // indirect
)

replace (
	github.com/providentia/digna/accountant_dashboard => ../accountant_dashboard
	github.com/providentia/digna/budget => ../budget
	github.com/providentia/digna/cash_flow => ../cash_flow
	github.com/providentia/digna/core_lume => ../core_lume
	github.com/providentia/digna/lifecycle => ../lifecycle
	github.com/providentia/digna/pdv_ui => ../pdv_ui
	github.com/providentia/digna/reporting => ../reporting
	github.com/providentia/digna/supply => ../supply
)
