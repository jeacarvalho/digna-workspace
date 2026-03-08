package domain

import (
	"errors"
	"fmt"
	"time"
)

// Member representa um cooperado/membro da entidade
type Member struct {
	ID        string
	EntityID  string
	Name      string
	Email     string
	Phone     string
	CPF       string // Opcional (LGPD)
	Role      MemberRole
	Status    MemberStatus
	JoinedAt  time.Time
	Skills    []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MemberRole define o papel do membro na entidade
type MemberRole string

const (
	RoleCoordinator MemberRole = "COORDINATOR"
	RoleMember      MemberRole = "MEMBER"
	RoleAdvisor     MemberRole = "ADVISOR"
)

// MemberStatus define o status do membro
type MemberStatus string

const (
	StatusActive   MemberStatus = "ACTIVE"
	StatusInactive MemberStatus = "INACTIVE"
)

var (
	ErrMemberInvalidName    = errors.New("member name is required")
	ErrMemberInvalidEmail   = errors.New("member email is invalid")
	ErrMemberInvalidRole    = errors.New("member role is invalid")
	ErrMemberInvalidStatus  = errors.New("member status is invalid")
	ErrMemberDuplicateEmail = errors.New("member email already exists")
)

// Validate verifica se o membro está válido
func (m *Member) Validate() error {
	if m.Name == "" {
		return ErrMemberInvalidName
	}

	if m.Email == "" {
		return ErrMemberInvalidEmail
	}

	if !isValidRole(m.Role) {
		return ErrMemberInvalidRole
	}

	if !isValidStatus(m.Status) {
		return ErrMemberInvalidStatus
	}

	return nil
}

// CanVote verifica se o membro tem direito a voto
func (m *Member) CanVote() bool {
	return m.Status == StatusActive && (m.Role == RoleCoordinator || m.Role == RoleMember)
}

// IsCoordinator verifica se o membro é coordenador
func (m *Member) IsCoordinator() bool {
	return m.Role == RoleCoordinator && m.Status == StatusActive
}

// CanManage verifica se o membro pode gerenciar a entidade
func (m *Member) CanManage() bool {
	return m.IsCoordinator()
}

// Deactivate desativa o membro
func (m *Member) Deactivate() {
	m.Status = StatusInactive
	m.UpdatedAt = time.Now()
}

// Activate ativa o membro
func (m *Member) Activate() {
	m.Status = StatusActive
	m.UpdatedAt = time.Now()
}

// AddSkill adiciona uma habilidade ao membro
func (m *Member) AddSkill(skill string) {
	for _, s := range m.Skills {
		if s == skill {
			return
		}
	}
	m.Skills = append(m.Skills, skill)
	m.UpdatedAt = time.Now()
}

// RemoveSkill remove uma habilidade do membro
func (m *Member) RemoveSkill(skill string) {
	for i, s := range m.Skills {
		if s == skill {
			m.Skills = append(m.Skills[:i], m.Skills[i+1:]...)
			m.UpdatedAt = time.Now()
			return
		}
	}
}

// String retorna representação textual do membro
func (m *Member) String() string {
	return fmt.Sprintf("Member{ID: %s, Name: %s, Role: %s, Status: %s}",
		m.ID, m.Name, m.Role, m.Status)
}

// MemberStats contém estatísticas do membro
type MemberStats struct {
	MemberID            string
	Name                string
	TotalWorkMinutes    int64
	TotalWorkHours      float64
	NumberOfWorkLogs    int
	AverageWorkPerDay   float64
	LastWorkDate        *time.Time
	ContributionPercent float64
}

// isValidRole verifica se o papel é válido
func isValidRole(role MemberRole) bool {
	switch role {
	case RoleCoordinator, RoleMember, RoleAdvisor:
		return true
	}
	return false
}

// isValidStatus verifica se o status é válido
func isValidStatus(status MemberStatus) bool {
	switch status {
	case StatusActive, StatusInactive:
		return true
	}
	return false
}
