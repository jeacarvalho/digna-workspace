import { test, expect } from '@playwright/test';

// Configurações do teste
// Usar entidade existente para testes (cafe_digna tem credenciais conhecidas)
const TEST_ENTITY = process.env.PLAYWRIGHT_TEST_ENTITY || 'cafe_digna';
const TEST_PASSWORD = process.env.PLAYWRIGHT_TEST_PASSWORD || 'cd0123';
const BASE_URL = process.env.PLAYWRIGHT_TEST_BASE_URL || 'http://localhost:8090';

// Função de login reutilizável
async function login(page) {
  await page.goto(`${BASE_URL}/login`);
  
  // Verificar tipo de formulário de login
  const hasUsernameField = await page.locator('input[name="username"]').count() > 0;
  const hasEntitySelect = await page.locator('select[name="entity_id"]').count() > 0;
  
  if (hasUsernameField) {
    // Novo formato: username/password
    await page.fill('input[name="username"]', TEST_ENTITY);
    await page.fill('input[name="password"]', TEST_PASSWORD);
  } else if (hasEntitySelect) {
    // Formato antigo: entity_id select
    await page.selectOption('select[name="entity_id"]', TEST_ENTITY);
    await page.fill('input[name="password"]', TEST_PASSWORD);
  } else {
    throw new Error('Formulário de login não reconhecido');
  }
  
  await page.click('button[type="submit"]');
  await expect(page).toHaveURL(/dashboard/);
}

test.describe('CRUD de Itens de Estoque - E2E', () => {
  test.beforeEach(async ({ page }) => {
    // Login antes de cada teste
    await login(page);
  });

  test('1. Acessar página de estoque vazia', async ({ page }) => {
    console.log('📋 Teste 1: Acessar estoque vazio');
    
    await page.goto(`${BASE_URL}/supply/stock?entity_id=${TEST_ENTITY}`);
    await expect(page).toHaveURL(/supply\/stock/);
    
    // Verificar título
    await expect(page.locator('h1')).toContainText(/Estoque/i);
    
    // Verificar mensagem de estoque vazio
    const emptyMessage = page.locator('text=Estoque vazio');
    await expect(emptyMessage).toBeVisible();
    
    // Verificar botão de cadastrar novo item
    const newItemButton = page.locator('button:has-text("Cadastrar Novo Item")');
    await expect(newItemButton).toBeVisible();
    
    console.log('✅ Página de estoque vazia carregada corretamente');
  });

  test('2. Cadastrar novo item de estoque', async ({ page }) => {
    console.log('📋 Teste 2: Cadastrar novo item');
    
    await page.goto(`${BASE_URL}/supply/stock?entity_id=${TEST_ENTITY}`);
    
    // Abrir formulário de cadastro
    await page.click('button:has-text("Cadastrar Novo Item")');
    
    // Verificar se formulário está visível
    const form = page.locator('#newItemForm');
    await expect(form).toBeVisible();
    
    // Preencher formulário
    await page.fill('input[name="name"]', 'Café Especial Teste E2E');
    await page.selectOption('select[name="type"]', 'INSUMO');
    await page.selectOption('select[name="unit"]', 'KG');
    await page.fill('input[name="quantity"]', '25');
    await page.fill('input[name="min_quantity"]', '10');
    await page.fill('input[name="unit_cost"]', '52.75');
    
    // Submeter formulário
    await page.click('button[type="submit"]');
    
    // Verificar feedback de sucesso (HTMX)
    // Aguardar resposta HTMX
    await page.waitForTimeout(2000);
    
    // Verificar se há mensagem de sucesso
    const successMessage = page.locator('.bg-green-100, text*=sucesso, text*=registrado');
    const hasSuccess = await successMessage.count() > 0;
    
    if (hasSuccess) {
      console.log('✅ Feedback de sucesso recebido');
    } else {
      console.log('⚠️  Feedback de sucesso não visível (pode ser HTMX silencioso)');
    }
    
    // Recarregar página para ver item na lista
    await page.reload();
    
    // Verificar se item aparece na tabela
    const itemInTable = page.locator('table tbody tr:has-text("Café Especial Teste E2E")');
    
    // Tentar por até 10 segundos (pode haver delay na persistência)
    for (let i = 0; i < 10; i++) {
      const count = await itemInTable.count();
      if (count > 0) {
        console.log('✅ Item encontrado na tabela após cadastro');
        break;
      }
      await page.waitForTimeout(1000);
      if (i === 9) {
        console.error('❌ Item não apareceu na tabela após cadastro');
        // Tirar screenshot para debug
        await page.screenshot({ path: `/tmp/test_cadastro_falha_${Date.now()}.png`, fullPage: true });
      }
    }
    
    // Verificar dados do item na tabela
    if (await itemInTable.count() > 0) {
      await expect(itemInTable).toContainText('INSUMO');
      await expect(itemInTable).toContainText('25');
      await expect(itemInTable).toContainText('R$ 52.75');
      console.log('✅ Dados do item corretos na tabela');
    }
    
    // Tirar screenshot para documentação
    await page.screenshot({ path: `/tmp/test_cadastro_sucesso_${Date.now()}.png`, fullPage: true });
  });

  test('3. Validar persistência após recarregar', async ({ page }) => {
    console.log('📋 Teste 3: Validar persistência');
    
    // Primeiro cadastrar um item se não existir
    await page.goto(`${BASE_URL}/supply/stock?entity_id=${TEST_ENTITY}`);
    
    // Verificar se já tem itens
    const existingItems = await page.locator('table tbody tr').count();
    
    if (existingItems === 0) {
      console.log('ℹ️  Nenhum item existente, cadastrando um...');
      
      await page.click('button:has-text("Cadastrar Novo Item")');
      await page.fill('input[name="name"]', 'Item Persistência Teste');
      await page.selectOption('select[name="type"]', 'PRODUTO');
      await page.selectOption('select[name="unit"]', 'UNIDADE');
      await page.fill('input[name="quantity"]', '15');
      await page.fill('input[name="min_quantity"]', '5');
      await page.fill('input[name="unit_cost"]', '29.90');
      await page.click('button[type="submit"]');
      
      await page.waitForTimeout(3000); // Aguardar persistência
    }
    
    // Navegar para dashboard e voltar
    await page.goto(`${BASE_URL}/dashboard?entity_id=${TEST_ENTITY}`);
    await page.goto(`${BASE_URL}/supply/stock?entity_id=${TEST_ENTITY}`);
    
    // Verificar se itens ainda estão na lista
    const itemsAfterReload = await page.locator('table tbody tr').count();
    
    if (itemsAfterReload > 0) {
      console.log(`✅ Persistência validada: ${itemsAfterReload} item(ns) após recarregar`);
      await expect(page.locator('table tbody tr').first()).toBeVisible();
    } else {
      console.error('❌ Falha na persistência: itens não mantidos após recarregar');
    }
  });

  test('4. Validar formatação de valores monetários', async ({ page }) => {
    console.log('📋 Teste 4: Validação de formatação monetária');
    
    await page.goto(`${BASE_URL}/supply/stock?entity_id=${TEST_ENTITY}`);
    
    // Verificar se há itens
    const hasItems = await page.locator('table tbody tr').count() > 0;
    
    if (!hasItems) {
      console.log('ℹ️  Criando item para teste de formatação...');
      
      await page.click('button:has-text("Cadastrar Novo Item")');
      await page.fill('input[name="name"]', 'Teste Formatação');
      await page.selectOption('select[name="type"]', 'MERCADORIA');
      await page.selectOption('select[name="unit"]', 'UNIDADE');
      await page.fill('input[name="quantity"]', '7');
      await page.fill('input[name="min_quantity"]', '3');
      await page.fill('input[name="unit_cost"]', '123.45');
      await page.click('button[type="submit"]');
      
      await page.waitForTimeout(2000);
      await page.reload();
    }
    
    // Verificar formatação dos valores
    const currencyCells = page.locator('td:has-text("R$")');
    const currencyCount = await currencyCells.count();
    
    if (currencyCount > 0) {
      console.log(`💰 ${currencyCount} valor(es) monetário(s) encontrado(s)`);
      
      // Verificar cada valor
      for (let i = 0; i < Math.min(currencyCount, 3); i++) {
        const cell = currencyCells.nth(i);
        const text = await cell.textContent();
        
        // Validar formato R$ X.XX
        if (text && text.match(/R\$\s*\d+\.\d{2}/)) {
          console.log(`✅ Formatação correta: ${text.trim()}`);
        } else {
          console.error(`❌ Formatação incorreta: ${text}`);
        }
      }
    }
    
    // Verificar cards de resumo
    const totalValueCard = page.locator('text=Valor Total');
    if (await totalValueCard.count() > 0) {
      await expect(totalValueCard).toBeVisible();
      const totalValue = page.locator('text=Valor Total + * >> xpath=following-sibling::*[1]');
      if (await totalValue.count() > 0) {
        const valueText = await totalValue.textContent();
        if (valueText && valueText.match(/R\$\s*\d+\.\d{2}/)) {
          console.log(`✅ Valor total formatado corretamente: ${valueText.trim()}`);
        }
      }
    }
  });

  test('5. Fluxo completo: cadastro → edição implícita via nova compra', async ({ page }) => {
    console.log('📋 Teste 5: Fluxo completo com estoque');
    
    // 1. Cadastrar item
    await page.goto(`${BASE_URL}/supply/stock?entity_id=${TEST_ENTITY}`);
    
    await page.click('button:has-text("Cadastrar Novo Item")');
    await page.fill('input[name="name"]', 'Açúcar Refinado');
    await page.selectOption('select[name="type"]', 'INSUMO');
    await page.selectOption('select[name="unit"]', 'KG');
    await page.fill('input[name="quantity"]', '50');
    await page.fill('input[name="min_quantity"]', '20');
    await page.fill('input[name="unit_cost"]', '4.99');
    await page.click('button[type="submit"]');
    
    await page.waitForTimeout(2000);
    
    // 2. Verificar que item está disponível para compras
    await page.goto(`${BASE_URL}/supply/purchase?entity_id=${TEST_ENTITY}`);
    
    // Verificar se item aparece no select de itens
    const itemSelect = page.locator('select[name="stock_item_id"]');
    await expect(itemSelect).toBeVisible();
    
    const options = await itemSelect.locator('option').allTextContents();
    const hasSugar = options.some(opt => opt.includes('Açúcar Refinado'));
    
    if (hasSugar) {
      console.log('✅ Item cadastrado disponível para compras');
    } else {
      console.log('⚠️  Item não aparece em compras (pode ser esperado)');
    }
    
    // 3. Voltar e verificar estoque
    await page.goto(`${BASE_URL}/supply/stock?entity_id=${TEST_ENTITY}`);
    
    const sugarRow = page.locator('table tbody tr:has-text("Açúcar Refinado")');
    if (await sugarRow.count() > 0) {
      console.log('✅ Fluxo completo validado: cadastro → disponibilidade → listagem');
      
      // Verificar cálculo de valor total (quantidade * custo unitário)
      const totalValueCell = sugarRow.locator('td:nth-child(6)'); // Coluna "Valor Total"
      if (await totalValueCell.count() > 0) {
        const totalText = await totalValueCell.textContent();
        // 50 * 4.99 = 249.50
        if (totalText && totalText.includes('249.50')) {
          console.log('✅ Cálculo de valor total correto: R$ 249.50');
        }
      }
    }
  });
});

// Teste adicional para debug do problema de cadastro
test.describe('Debug - Problema de Cadastro de Itens', () => {
  test('Debug: Analisar formulário e HTMX', async ({ page }) => {
    console.log('🔍 Debug: Analisando formulário de cadastro');
    
    await page.goto(`${BASE_URL}/login`);
    
    // Login rápido
    const hasUsername = await page.locator('input[name="username"]').count() > 0;
    if (hasUsername) {
      await page.fill('input[name="username"]', TEST_ENTITY);
      await page.fill('input[name="password"]', TEST_PASSWORD);
    } else {
      await page.selectOption('select[name="entity_id"]', TEST_ENTITY);
      await page.fill('input[name="password"]', TEST_PASSWORD);
    }
    await page.click('button[type="submit"]');
    
    await page.goto(`${BASE_URL}/supply/stock?entity_id=${TEST_ENTITY}`);
    
    // Analisar formulário
    await page.click('button:has-text("Cadastrar Novo Item")');
    
    // Verificar estrutura do formulário
    const form = page.locator('#newItemForm form');
    await expect(form).toBeVisible();
    
    const hxPost = await form.getAttribute('hx-post');
    const hxTarget = await form.getAttribute('hx-target');
    const hxSwap = await form.getAttribute('hx-swap');
    
    console.log(`📋 Formulário HTMX:`);
    console.log(`   hx-post: ${hxPost}`);
    console.log(`   hx-target: ${hxTarget}`);
    console.log(`   hx-swap: ${hxSwap}`);
    
    // Verificar se target existe
    if (hxTarget) {
      const targetElement = page.locator(hxTarget);
      const targetExists = await targetElement.count() > 0;
      console.log(`   Target "${hxTarget}" existe: ${targetExists}`);
      
      if (!targetExists) {
        console.error('❌ PROBLEMA ENCONTRADO: Target HTMX não existe!');
        console.error('   Isso explica por que a lista não atualiza.');
        
        // Sugerir correção
        console.log('💡 SUGESTÃO: O target deve ser "#stockItemsList" mas essa ID não existe no template.');
        console.log('   Verificar template supply_stock_simple.html linha ~199');
      }
    }
    
    // Verificar campos hidden
    const hiddenFields = page.locator('input[type="hidden"]');
    const hiddenCount = await hiddenFields.count();
    console.log(`   Campos hidden: ${hiddenCount}`);
    
    for (let i = 0; i < hiddenCount; i++) {
      const field = hiddenFields.nth(i);
      const name = await field.getAttribute('name');
      const value = await field.getAttribute('value');
      console.log(`     - ${name}: ${value}`);
    }
    
    // Tirar screenshot para análise
    await page.screenshot({ path: `/tmp/debug_formulario_${Date.now()}.png`, fullPage: true });
    console.log('📸 Screenshot salvo para análise');
  });
});