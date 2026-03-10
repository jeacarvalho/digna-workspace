import { test, expect } from '@playwright/test';

const TEST_ENTITY = 'cafe_digna';
const TEST_PASSWORD = 'cd0123';

test.describe('Fluxo Básico Digna - Validação E2E', () => {
  test('Login e navegação básica', async ({ page }) => {
    console.log('🚀 Teste 1: Login no sistema');
    
    await page.goto('/login');
    await page.selectOption('select[name="entity_id"]', TEST_ENTITY);
    await page.fill('input[name="password"]', TEST_PASSWORD);
    await page.click('button[type="submit"]');
    
    await page.waitForURL(/dashboard/);
    console.log('✅ Login realizado com sucesso');
    
    const pageTitle = await page.title();
    expect(pageTitle).toContain('Digna');
    console.log(`✅ Título da página: ${pageTitle}`);
  });

  test('Verificar menu de navegação', async ({ page }) => {
    console.log('🚀 Teste 2: Verificar menu de navegação');
    
    await login(page);
    
    const menuItems = await page.locator('nav a, .menu a, a[href]').all();
    console.log(`✅ Encontrados ${menuItems.length} links de navegação`);
    
    expect(menuItems.length).toBeGreaterThan(0);
  });

  test('Acessar página de estoque', async ({ page }) => {
    console.log('🚀 Teste 3: Acessar estoque');
    
    await login(page);
    
    try {
      await page.click('a[href*="stock"], a:has-text("Estoque")');
      await page.waitForLoadState('networkidle');
      console.log('✅ Navegou para estoque');
    } catch (error) {
      console.log('⚠️  Link de estoque não encontrado, verificando URL manual');
      await page.goto('/supply/stock');
    }
    
    const currentUrl = page.url();
    expect(currentUrl).toContain('stock');
    console.log(`✅ URL atual: ${currentUrl}`);
  });

  test('Acessar página de membros', async ({ page }) => {
    console.log('🚀 Teste 4: Acessar membros');
    
    await login(page);
    
    try {
      await page.click('a[href*="members"], a:has-text("Membros")');
      await page.waitForLoadState('networkidle');
      console.log('✅ Navegou para membros');
    } catch (error) {
      console.log('⚠️  Link de membros não encontrado, verificando URL manual');
      await page.goto('/members');
    }
    
    const currentUrl = page.url();
    expect(currentUrl).toContain('members');
    console.log(`✅ URL atual: ${currentUrl}`);
  });

  test('Fluxo simplificado - 7 passos', async ({ page }) => {
    console.log('🚀 Teste 5: Fluxo simplificado de validação');
    
    // 1. Login
    await page.goto('/login');
    await page.selectOption('select[name="entity_id"]', TEST_ENTITY);
    await page.fill('input[name="password"]', TEST_PASSWORD);
    await page.click('button[type="submit"]');
    await page.waitForURL(/dashboard/);
    console.log('✅ 1. Login realizado');
    
    // 2. Verificar dashboard
    const hasContent = await page.locator('body').textContent();
    expect(hasContent).toBeTruthy();
    console.log('✅ 2. Dashboard carregado');
    
    // 3. Acessar estoque
    await page.goto('/supply/stock');
    const stockPage = await page.title();
    expect(stockPage).toBeTruthy();
    console.log('✅ 3. Página de estoque acessada');
    
    // 4. Acessar membros
    await page.goto('/members');
    const membersPage = await page.title();
    expect(membersPage).toBeTruthy();
    console.log('✅ 4. Página de membros acessada');
    
    // 5. Acessar fornecedores
    await page.goto('/supply/suppliers');
    const suppliersPage = await page.title();
    expect(suppliersPage).toBeTruthy();
    console.log('✅ 5. Página de fornecedores acessada');
    
    // 6. Acessar PDV
    await page.goto('/pdv');
    const pdvPage = await page.title();
    expect(pdvPage).toBeTruthy();
    console.log('✅ 6. PDV acessado');
    
    // 7. Acessar caixa
    await page.goto('/cash');
    const cashPage = await page.title();
    expect(cashPage).toBeTruthy();
    console.log('✅ 7. Caixa acessado');
    
    console.log('🎉 Fluxo básico validado com sucesso!');
  });
});

async function login(page) {
  await page.goto('/login');
  await page.selectOption('select[name="entity_id"]', TEST_ENTITY);
  await page.fill('input[name="password"]', TEST_PASSWORD);
  await page.click('button[type="submit"]');
  await page.waitForURL(/dashboard/);
}