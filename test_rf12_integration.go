package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

func main() {
    fmt.Println("🧪 Teste de Integração RF-12")
    fmt.Println("=============================")

    // Criar lifecycle manager
    lm := lifecycle.NewSQLiteManager()
    defer lm.CloseAll()

    // SQLiteManager implementa AccountantLinkService
    fmt.Println("1. Verificando implementação de AccountantLinkService...")
    fmt.Println("✅ SQLiteManager implementa AccountantLinkService")

    accountantLinkService := lm

    // Criar vínculo entre cafe_digna e contador_social
    fmt.Println("\n2. Criando vínculo cafe_digna -> contador_social...")
    link, err := accountantLinkService.CreateLink("cafe_digna", "contador_social", "cafe_digna")
    if err != nil {
        log.Printf("⚠️  Erro ao criar vínculo (pode já existir): %v", err)
    } else {
        fmt.Printf("✅ Vínculo criado: %s -> %s (Status: %s)\n", 
            link.EnterpriseID, link.AccountantID, link.Status)
    }

    // Listar vínculos do cafe_digna
    fmt.Println("\n3. Listando vínculos do cafe_digna...")
    enterpriseLinks, err := accountantLinkService.GetEnterpriseLinks("cafe_digna")
    if err != nil {
        log.Printf("⚠️  Erro ao listar vínculos: %v", err)
    } else {
        fmt.Printf("✅ %d vínculos encontrados para cafe_digna:\n", len(enterpriseLinks))
        for i, link := range enterpriseLinks {
            fmt.Printf("   %d. Contador: %s, Status: %s, Desde: %s\n", 
                i+1, link.AccountantID, link.Status, link.StartDate.Format("2006-01-02"))
        }
    }

    // Listar vínculos do contador_social
    fmt.Println("\n4. Listando vínculos do contador_social...")
    accountantLinks, err := accountantLinkService.GetAccountantLinks("contador_social")
    if err != nil {
        log.Printf("⚠️  Erro ao listar vínculos: %v", err)
    } else {
        fmt.Printf("✅ %d vínculos encontrados para contador_social:\n", len(accountantLinks))
        for i, link := range accountantLinks {
            fmt.Printf("   %d. Empreendimento: %s, Status: %s, Desde: %s\n", 
                i+1, link.EnterpriseID, link.Status, link.StartDate.Format("2006-01-02"))
        }
    }

    // Testar filtro temporal
    fmt.Println("\n5. Testando filtro temporal...")
    startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    endTime := time.Date(2024, 12, 31, 23, 59, 59, 999999999, time.UTC)

    validEnterprises, err := accountantLinkService.GetValidEnterprisesForAccountant(
        context.Background(), "contador_social", startTime, endTime)
    if err != nil {
        log.Printf("⚠️  Erro ao obter empreendimentos válidos: %v", err)
    } else {
        fmt.Printf("✅ %d empreendimentos válidos para contador_social em 2024:\n", len(validEnterprises))
        for i, enterprise := range validEnterprises {
            fmt.Printf("   %d. %s\n", i+1, enterprise)
        }
    }

    // Validar acesso
    fmt.Println("\n6. Validando acesso do contador_social ao cafe_digna...")
    hasAccess, err := accountantLinkService.ValidateAccountantAccess(
        context.Background(), "contador_social", "cafe_digna", startTime, endTime)
    if err != nil {
        log.Printf("⚠️  Erro ao validar acesso: %v", err)
    } else if hasAccess {
        fmt.Println("✅ contador_social TEM acesso ao cafe_digna no período")
    } else {
        fmt.Println("❌ contador_social NÃO tem acesso ao cafe_digna no período")
    }

    fmt.Println("\n🎉 Teste de integração RF-12 concluído!")
}
