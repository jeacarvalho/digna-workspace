package domain

import "time"

type Supplier struct {
	ID          string
	Name        string
	ContactInfo string // Telefone, email ou contato simples
	CreatedAt   time.Time
}

func (s *Supplier) Validate() error {
	if s.Name == "" {
		return ErrInvalidSupplierName
	}
	return nil
}

var (
	ErrInvalidSupplierName = newDomainError("nome do fornecedor inválido")
)
