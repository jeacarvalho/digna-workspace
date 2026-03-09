module digna/accountant_dashboard

go 1.25.4

require github.com/mattn/go-sqlite3 v1.14.34

replace (
	github.com/providentia/digna/core_lume => ../core_lume
	github.com/providentia/digna/lifecycle => ../lifecycle
)
