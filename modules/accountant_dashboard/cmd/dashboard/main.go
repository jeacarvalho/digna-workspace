package main

import (
	"fmt"
	"log"
	"net/http"

	"digna/accountant_dashboard/internal/domain"
	"digna/accountant_dashboard/internal/handler"
	"digna/accountant_dashboard/internal/repository"
	"digna/accountant_dashboard/internal/service"
)

func main() {
	repo := repository.NewSQLiteFiscalAdapter()
	mapper := domain.NewDefaultAccountMapper()
	translator := service.NewTranslatorService(repo, mapper)
	dashboardHandler := handler.NewDashboardHandler(translator, mapper)

	http.HandleFunc("/accountant/dashboard", dashboardHandler.Dashboard)
	http.HandleFunc("/accountant/export", dashboardHandler.ExportFiscal)

	fmt.Println("🚀 Accountant Dashboard starting on http://localhost:8081/accountant/dashboard")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
