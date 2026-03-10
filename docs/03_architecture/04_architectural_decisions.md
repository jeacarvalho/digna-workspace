# Architectural Decisions Record - Sprint 12

## ADR-001: Integration via ui_web/main.go

### Context
The Sprint 12 prompt suggested creating a new entry point at `cmd/digna/main.go` for the Accountant Dashboard module. However, the project already has a well-established web interface architecture centered around the `ui_web` module.

### Decision
Integrate the Accountant Dashboard through the existing `modules/ui_web/main.go` instead of creating a new entry point.

### Consequences
**Positive:**
- Maintains architectural consistency across all web interfaces
- Centralizes HTTP route management
- Simplifies deployment (single binary for all web interfaces)
- Follows DRY principle by reusing existing infrastructure

**Negative:**
- Deviates from the original prompt specification
- Creates tighter coupling between `accountant_dashboard` and `ui_web` modules

### Status
Accepted and implemented.

---

## ADR-002: Embedded Templates vs Separate HTML Files

### Context
The prompt suggested creating separate HTML template files (`layout.html`, `dashboard.html`). The project has examples of both embedded templates and separate template files.

### Decision
Use embedded Go templates within the handler code instead of separate HTML files.

### Consequences
**Positive:**
- Simplifies deployment (fewer files to manage)
- Improves performance (templates compiled with binary)
- Enhances coherency (HTML close to handler logic)
- Easier testing (templates tested with handler code)

**Negative:**
- Less separation of concerns between logic and presentation
- Harder for non-developers to modify templates

### Status
Accepted and implemented.

---

## ADR-003: No Separate templates/ Directory

### Context
The prompt suggested creating a `templates/` directory within the `accountant_dashboard` module.

### Decision
Do not create a separate `templates/` directory since templates are embedded in the code.

### Consequences
**Positive:**
- Reduces project complexity
- Follows YAGNI principle (no unnecessary structure)
- Consistent with other modules using embedded templates

**Negative:**
- Deviates from prompt specification
- Less conventional for web development

### Status
Accepted and implemented.

---

## ADR-004: Public API Package Structure

### Context
The prompt focused on internal implementation but didn't specify external API design.

### Decision
Create a public API package (`pkg/dashboard/`) for external consumption of the Accountant Dashboard functionality.

### Consequences
**Positive:**
- Enables integration with other systems
- Provides clean separation between internal and external APIs
- Follows Go best practices for package design

**Negative:**
- Additional complexity beyond prompt requirements
- More interfaces to maintain

### Status
Accepted and implemented.

---

## Principles Applied

### 1. KISS (Keep It Simple)
- Chose the simplest solution that meets all requirements
- Avoided unnecessary architectural complexity

### 2. YAGNI (You Ain't Gonna Need It)
- Didn't implement structures not immediately needed
- Templates directory not created since embedded templates suffice

### 3. DRY (Don't Repeat Yourself)
- Reused existing `ui_web` infrastructure instead of creating new
- Leveraged established patterns from other modules

### 4. Consistency
- Maintained architectural patterns established in previous sprints
- Followed project conventions for module structure

---

## Validation

All architectural decisions have been validated through:

1. **Functional Testing:** All features work as specified
2. **Performance Testing:** Embedded templates show no performance degradation
3. **Integration Testing:** Seamless integration with existing `ui_web` module
4. **Code Quality:** High test coverage (69% total, 93.9% core packages)

---

## References

- [Session Log 012](../06_roadmap/05_session_log.md#session-log-012---sprint-12-painel-do-contador-social---decisões-arquiteturais)
- [Sprint 12 Status](../06_roadmap/04_status.md#sprint-12-painel-do-contador-social-accountant-dashboard--complete)
- [System Architecture](../01_system.md)