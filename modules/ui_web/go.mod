module github.com/providentia/digna/ui_web

go 1.25.4

require (
	github.com/providentia/digna/core_lume v0.0.0
	github.com/providentia/digna/lifecycle v0.0.0
	github.com/providentia/digna/pdv_ui v0.0.0
	github.com/providentia/digna/reporting v0.0.0
)

replace (
	github.com/providentia/digna/core_lume => ../core_lume
	github.com/providentia/digna/lifecycle => ../lifecycle
	github.com/providentia/digna/pdv_ui => ../pdv_ui
	github.com/providentia/digna/reporting => ../reporting
)
