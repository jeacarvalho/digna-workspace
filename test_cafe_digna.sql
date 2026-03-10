-- Script para popular o database da cafe_digna com dados de teste
-- Conecte-se ao database: sqlite3 /home/s015533607/Documentos/desenv/digna-workspace/data/entities/cafe_digna.db

-- 1. Adicionar um fornecedor
INSERT INTO suppliers (id, name, contact, created_at) VALUES 
('supplier_001', 'Fazenda Café Bom', 'João Silva - (11) 99999-9999', strftime('%s', 'now'));

-- 2. Adicionar uma compra
INSERT INTO purchases (id, supplier_id, total_amount, status, purchase_date, created_at) VALUES 
('purchase_001', 'supplier_001', 5000, 'PAGO', strftime('%s', 'now'), strftime('%s', 'now'));

-- 3. Adicionar itens ao estoque (café em grão e café moído)
INSERT INTO stock_items (id, name, type, unit, quantity, min_quantity, unit_cost, created_at) VALUES 
('item_001', 'Café em Grão Arábica', 'INSUMO', 'KG', 50, 10, 4500, strftime('%s', 'now')),
('item_002', 'Café Moído para Venda', 'PRODUTO', 'KG', 20, 5, 8000, strftime('%s', 'now')),
('item_003', 'Café Torrado em Grão', 'MERCADORIA', 'KG', 30, 8, 7000, strftime('%s', 'now'));

-- 4. Adicionar itens da compra
INSERT INTO purchase_items (id, purchase_id, stock_item_id, quantity, unit_price, created_at) VALUES 
('pi_001', 'purchase_001', 'item_001', 50, 4500, strftime('%s', 'now')),
('pi_002', 'purchase_001', 'item_002', 20, 8000, strftime('%s', 'now'));

-- 5. Adicionar uma entrada no ledger para a compra
INSERT INTO entries (id, entry_date, description, reference, created_at) VALUES 
(1, strftime('%s', 'now'), 'Compra de café da Fazenda Café Bom', 'purchase_001', strftime('%s', 'now'));

-- 6. Adicionar postings (lançamentos contábeis)
-- Primeiro, verifique os IDs das contas
SELECT id, code, name FROM accounts;

-- Se não houver contas, crie algumas básicas
INSERT INTO accounts (id, code, name, type, created_at) VALUES 
('acc_001', '1', 'Caixa', 'ASSET', strftime('%s', 'now')),
('acc_002', '2', 'Estoque', 'ASSET', strftime('%s', 'now')),
('acc_003', '3', 'Receitas', 'REVENUE', strftime('%s', 'now')),
('acc_004', '4', 'Despesas', 'EXPENSE', strftime('%s', 'now'));

-- Agora adicione os postings para a compra
INSERT INTO postings (id, entry_id, account_id, amount, direction, created_at) VALUES 
(1, 1, 'acc_001', 5000, 'DEBIT', strftime('%s', 'now')),  -- Saída de caixa
(2, 1, 'acc_002', 5000, 'CREDIT', strftime('%s', 'now')); -- Entrada no estoque

-- Verificar os dados inseridos
SELECT '=== FORNECEDORES ===';
SELECT * FROM suppliers;

SELECT '=== ESTOQUE ===';
SELECT id, name, type, quantity, unit_cost FROM stock_items;

SELECT '=== COMPRAS ===';
SELECT p.id, s.name as supplier, p.total_amount, p.status FROM purchases p
JOIN suppliers s ON p.supplier_id = s.id;

SELECT '=== CONTAS ===';
SELECT code, name, type FROM accounts;