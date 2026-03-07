module github.com/providentia/digna/legal_facade

go 1.25.4

require (
	github.com/providentia/digna/core_lume v0.0.0
	github.com/providentia/digna/lifecycle v0.0.0
)

replace (
	github.com/providentia/digna/core_lume => ../core_lume
	github.com/providentia/digna/lifecycle => ../lifecycle
)
