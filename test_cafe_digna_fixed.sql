-- Script corrigido para popular o database da cafe_digna
-- Conecte-se ao database: sqlite3 /home/s015533607/Documentos/desenv/digna-workspace/data/entities/cafe_digna.db

-- Limpar dados existentes (se necessário)
DELETE FROM purchase_items;
DELETE FROM purchases;
DELETE FROM stock_items;
DELETE FROM suppliers;

-- 1. Adicionar um fornecedor
INSERT INTO suppliers (id, name, contact_info, created_at) VALUES 
('supplier_001', 'Fazenda Café Bom', 'João Silva - (11) 99999-9999', strftime('%s', 'now'));

-- 2. Adicionar itens ao estoque (café em grão e café moído)
-- Valores em centavos: R$ 45,00 = 4500 centavos
INSERT INTO stock_items (id, name, type, unit, quantity, min_quantity, unit_cost, created_at) VALUES 
('item_001', 'Café em Grão Arábica', 'INSUMO', 'KG', 50, 10, 4500, strftime('%s', 'now')),
('item_002', 'Café Moído para Venda', 'PRODUTO', 'KG', 20, 5, 8000, strftime('%s', 'now')),
('item_003', 'Café Torrado em Grão', 'MERCADORIA', 'KG', 30, 8, 7000, strftime('%s', 'now'));

-- 3. Adicionar uma compra (total: 50kg * R$45 + 20kg * R$80 = R$2250 + R$1600 = R$3850 = 385000 centavos)
INSERT INTO purchases (id, supplier_id, total_value, date, created_at) VALUES 
('purchase_001', 'supplier_001', 385000, strftime('%s', 'now'), strftime('%s', 'now'));

-- 4. Adicionar itens da compra
INSERT INTO purchase_items (id, purchase_id, stock_item_id, quantity, unit_cost, total_cost) VALUES 
('pi_001', 'purchase_001', 'item_001', 50, 4500, 225000),
('pi_002', 'purchase_001', 'item_002', 20, 8000, 160000);

-- Verificar os dados inseridos
SELECT '=== FORNECEDORES ===';
SELECT id, name, contact_info FROM suppliers;

SELECT '=== ESTOQUE ===';
SELECT id, name, type, quantity, unit_cost, 
       unit_cost/100.0 as preco_unitario_reais,
       (quantity * unit_cost)/100.0 as valor_total_reais
FROM stock_items;

SELECT '=== COMPRAS ===';
SELECT p.id, s.name as fornecedor, 
       p.total_value/100.0 as total_reais,
       datetime(p.date, 'unixepoch') as data_compra
FROM purchases p
JOIN suppliers s ON p.supplier_id = s.id;

SELECT '=== ITENS DA COMPRA ===';
SELECT pi.id, si.name as item, 
       pi.quantity, 
       pi.unit_cost/100.0 as preco_unitario,
       pi.total_cost/100.0 as total_item
FROM purchase_items pi
JOIN stock_items si ON pi.stock_item_id = si.id;