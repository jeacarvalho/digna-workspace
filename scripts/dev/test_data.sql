-- Dados de teste para smoke tests e E2E
-- Este arquivo é carregado automaticamente em ambientes de teste

-- Usuário de teste padrão
INSERT INTO users (id, email, name, role, status, created_at) VALUES 
('test-user-001', 'test@digna.local', 'Usuário de Teste', 'COORDINATOR', 'ACTIVE', strftime('%s', 'now'));

-- Entidade de teste
INSERT INTO entities (id, name, cnpj, status, created_at) VALUES 
('test-entity-001', 'Entidade de Teste', '00000000000191', 'ACTIVE', strftime('%s', 'now'));

-- Vínculo usuário-entidade
INSERT INTO user_entities (user_id, entity_id, role, created_at) VALUES 
('test-user-001', 'test-entity-001', 'COORDINATOR', strftime('%s', 'now'));

-- Membro de teste
INSERT INTO members (id, entity_id, name, email, role, status, created_at) VALUES 
('test-member-001', 'test-entity-001', 'Coordenador Teste', 'coord@test.local', 'COORDINATOR', 'ACTIVE', strftime('%s', 'now'));
