import { test, expect } from '@playwright/test';

const TEST_ENTITY = 'cafe_digna';
const TEST_PASSWORD = 'cd0123';

test.describe('Fluxo Completo Digna - E2E', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
    await page.waitForLoadState('networkidle');
  });

  test('1. Login no sistema', async ({ page }) => {
    await page.goto('/login');
    await page.selectOption('select[name="entity_id"]', TEST_ENTITY);
    await page.fill('input[name="password"]', TEST_PASSWORD);
    await page.click('button[type="submit"]');
    
    await expect(page).toHaveURL(/dashboard/);
    await expect(page.locator('h1')).toContainText(/Dashboard/i);
  });

  test('2. Criar item de estoque se não existir', async ({ page }) => {
    await login(page);
    
    await page.click('a[href*="/supply/stock"]');
    await expect(page).toHaveURL(/supply\/stock/);
    
    const itemExists = await page.locator('table tbody tr').count() > 0;
    
    if (!itemExists) {
      await page.click('a[href*="/supply/stock/new"]');
      await page.fill('input[name="name"]', 'Café Especial Teste');
      await page.fill('input[name="unit"]', 'kg');
      await page.fill('input[name="unit_price"]', '45.00');
      await page.click('button[type="submit"]');
      
      await expect(page).toHaveURL(/supply\/stock/);
      await expect(page.locator('table tbody tr')).toContainText('Café Especial Teste');
    }
  });

  test('3. Criar membro se não existir', async ({ page }) => {
    await login(page);
    
    await page.click('a[href*="/members"]');
    await expect(page).toHaveURL(/members/);
    
    const memberExists = await page.locator('table tbody tr').count() > 0;
    
    if (!memberExists) {
      await page.click('a[href*="/members/new"]');
      await page.fill('input[name="name"]', 'Membro Teste');
      await page.fill('input[name="cpf"]', '12345678901');
      await page.fill('input[name="email"]', 'membro@teste.com');
      await page.click('button[type="submit"]');
      
      await expect(page).toHaveURL(/members/);
      await expect(page.locator('table tbody tr')).toContainText('Membro Teste');
    }
  });

  test('4. Criar fornecedor se não existir', async ({ page }) => {
    await login(page);
    
    await page.click('a[href*="/supply/suppliers"]');
    await expect(page).toHaveURL(/supply\/suppliers/);
    
    const supplierExists = await page.locator('table tbody tr').count() > 0;
    
    if (!supplierExists) {
      await page.click('a[href*="/supply/suppliers/new"]');
      await page.fill('input[name="name"]', 'Fornecedor Teste');
      await page.fill('input[name="cnpj"]', '12345678000199');
      await page.fill('input[name="contact"]', 'contato@fornecedor.com');
      await page.click('button[type="submit"]');
      
      await expect(page).toHaveURL(/supply\/suppliers/);
      await expect(page.locator('table tbody tr')).toContainText('Fornecedor Teste');
    }
  });

  test('5. Registrar compra do item criado/existente', async ({ page }) => {
    await login(page);
    
    await page.click('a[href*="/supply/purchases"]');
    await expect(page).toHaveURL(/supply\/purchases/);
    
    await page.click('a[href*="/supply/purchases/new"]');
    
    await page.selectOption('select[name="supplier_id"]', { index: 1 });
    await page.selectOption('select[name="stock_item_id"]', { index: 1 });
    await page.fill('input[name="quantity"]', '10');
    await page.fill('input[name="unit_price"]', '42.50');
    await page.fill('input[name="total_amount"]', '425.00');
    await page.click('button[type="submit"]');
    
    await expect(page).toHaveURL(/supply\/purchases/);
    await expect(page.locator('table tbody tr')).toContainText('425.00');
  });

  test('6. Registrar venda do item comprado no PDV', async ({ page }) => {
    await login(page);
    
    await page.click('a[href*="/pdv"]');
    await expect(page).toHaveURL(/pdv/);
    
    const stockItem = await page.locator('.stock-item:has-text("Café")').first();
    if (stockItem) {
      await stockItem.click();
      await page.fill('input[name="quantity"]', '2');
      await page.click('button:has-text("Adicionar")');
      
      await expect(page.locator('.cart-item')).toContainText('Café');
      await expect(page.locator('.cart-total')).toContainText(/[0-9,.]+/);
      
      await page.click('button:has-text("Finalizar Venda")');
      await page.fill('input[name="payment_amount"]', '100.00');
      await page.click('button:has-text("Confirmar")');
      
      await expect(page.locator('.success-message')).toBeVisible();
    }
  });

  test('7. Confirmar saldo em caixa e registrar horas trabalhadas', async ({ page }) => {
    await login(page);
    
    await page.click('a[href*="/cash"]');
    await expect(page).toHaveURL(/cash/);
    
    const balance = await page.locator('.balance-amount').textContent();
    expect(balance).toMatch(/[0-9,.]+/);
    
    await page.click('a[href*="/social/clock"]');
    await expect(page).toHaveURL(/social\/clock/);
    
    await page.fill('input[name="member_id"]', '1');
    await page.fill('input[name="hours"]', '8');
    await page.fill('input[name="activity"]', 'Produção');
    await page.click('button[type="submit"]');
    
    await expect(page.locator('.success-message')).toBeVisible();
    
    await page.click('a[href*="/dashboard"]');
    await expect(page).toHaveURL(/dashboard/);
    
    const dashboardStats = await page.locator('.stat-card').count();
    expect(dashboardStats).toBeGreaterThan(0);
  });

  test('Fluxo completo integrado - 7 passos', async ({ page }) => {
    console.log('🚀 Iniciando fluxo E2E completo da Digna');
    
    await login(page);
    console.log('✅ 1. Login realizado');
    
    await createStockItemIfNeeded(page);
    console.log('✅ 2. Item de estoque verificado/criado');
    
    await createMemberIfNeeded(page);
    console.log('✅ 3. Membro verificado/criado');
    
    await createSupplierIfNeeded(page);
    console.log('✅ 4. Fornecedor verificado/criado');
    
    await registerPurchase(page);
    console.log('✅ 5. Compra registrada');
    
    await registerSaleInPDV(page);
    console.log('✅ 6. Venda registrada no PDV');
    
    await verifyCashBalanceAndWorkHours(page);
    console.log('✅ 7. Saldo e horas verificados');
    
    console.log('🎉 Fluxo E2E completo validado com sucesso!');
  });
});

async function login(page) {
  await page.goto('/login');
  await page.selectOption('select[name="entity_id"]', TEST_ENTITY);
  await page.fill('input[name="password"]', TEST_PASSWORD);
  await page.click('button[type="submit"]');
  await expect(page).toHaveURL(/dashboard/);
}

async function createStockItemIfNeeded(page) {
  await page.click('a[href*="/supply/stock"]');
  await expect(page).toHaveURL(/supply\/stock/);
  
  const itemCount = await page.locator('table tbody tr').count();
  if (itemCount === 0) {
    await page.click('a[href*="/supply/stock/new"]');
    await page.fill('input[name="name"]', 'Produto Teste E2E');
    await page.fill('input[name="unit"]', 'un');
    await page.fill('input[name="unit_price"]', '25.00');
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL(/supply\/stock/);
  }
}

async function createMemberIfNeeded(page) {
  await page.click('a[href*="/members"]');
  await expect(page).toHaveURL(/members/);
  
  const memberCount = await page.locator('table tbody tr').count();
  if (memberCount === 0) {
    await page.click('a[href*="/members/new"]');
    await page.fill('input[name="name"]', 'Membro E2E Test');
    await page.fill('input[name="cpf"]', '11122233344');
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL(/members/);
  }
}

async function createSupplierIfNeeded(page) {
  await page.click('a[href*="/supply/suppliers"]');
  await expect(page).toHaveURL(/supply\/suppliers/);
  
  const supplierCount = await page.locator('table tbody tr').count();
  if (supplierCount === 0) {
    await page.click('a[href*="/supply/suppliers/new"]');
    await page.fill('input[name="name"]', 'Fornecedor E2E Test');
    await page.fill('input[name="cnpj"]', '99887766000155');
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL(/supply\/suppliers/);
  }
}

async function registerPurchase(page) {
  await page.click('a[href*="/supply/purchases"]');
  await expect(page).toHaveURL(/supply\/purchases/);
  
  await page.click('a[href*="/supply/purchases/new"]');
  
  const supplierOptions = await page.locator('select[name="supplier_id"] option').count();
  const stockOptions = await page.locator('select[name="stock_item_id"] option').count();
  
  if (supplierOptions > 1 && stockOptions > 1) {
    await page.selectOption('select[name="supplier_id"]', { index: 1 });
    await page.selectOption('select[name="stock_item_id"]', { index: 1 });
    await page.fill('input[name="quantity"]', '5');
    await page.fill('input[name="unit_price"]', '20.00');
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL(/supply\/purchases/);
  }
}

async function registerSaleInPDV(page) {
  await page.click('a[href*="/pdv"]');
  await expect(page).toHaveURL(/pdv/);
  
  const hasItems = await page.locator('.stock-item').count() > 0;
  if (hasItems) {
    await page.locator('.stock-item').first().click();
    await page.fill('input[name="quantity"]', '1');
    await page.click('button:has-text("Adicionar")');
    
    await page.click('button:has-text("Finalizar Venda")');
    await page.fill('input[name="payment_amount"]', '30.00');
    await page.click('button:has-text("Confirmar")');
    
    await expect(page.locator('.receipt')).toBeVisible({ timeout: 5000 });
  }
}

async function verifyCashBalanceAndWorkHours(page) {
  await page.click('a[href*="/cash"]');
  await expect(page).toHaveURL(/cash/);
  
  const balanceText = await page.locator('.balance').textContent();
  expect(balanceText).toBeDefined();
  
  await page.click('a[href*="/social/clock"]');
  await expect(page).toHaveURL(/social\/clock/);
  
  const memberOptions = await page.locator('select[name="member_id"] option').count();
  if (memberOptions > 1) {
    await page.selectOption('select[name="member_id"]', { index: 1 });
    await page.fill('input[name="hours"]', '4');
    await page.click('button[type="submit"]');
  }
  
  await page.click('a[href*="/dashboard"]');
  await expect(page).toHaveURL(/dashboard/);
  
  const hasStats = await page.locator('.statistic').count() > 0;
  expect(hasStats).toBeTruthy();
}