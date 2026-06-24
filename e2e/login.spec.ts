import { test, expect } from '@playwright/test'

test.describe('Login Flow', () => {
  test('should show login page', async ({ page }) => {
    await page.goto('/login')
    await expect(page.getByText('Visual Mesin')).toBeVisible()
    await expect(page.getByText('Silakan login untuk melanjutkan')).toBeVisible()
  })

  test('should login with valid admin credentials', async ({ page }) => {
    await page.goto('/login')
    await page.getByPlaceholder('Email / NIP').fill('admin@admin.com')
    await page.getByPlaceholder('Password').fill('password')
    await page.getByRole('button', { name: /login/i }).click()
    await expect(page).toHaveURL(/\/dashboard/, { timeout: 15000 })
  })

  test('should show error with invalid credentials', async ({ page }) => {
    await page.goto('/login')
    await page.getByPlaceholder('Email / NIP').fill('wrong@email.com')
    await page.getByPlaceholder('Password').fill('wrongpass')
    await page.getByRole('button', { name: /login/i }).click()
    await expect(page.getByText(/gagal|error|salah/i)).toBeVisible({ timeout: 10000 })
  })

  test('should logout successfully', async ({ page }) => {
    await page.goto('/login')
    await page.getByPlaceholder('Email / NIP').fill('admin@admin.com')
    await page.getByPlaceholder('Password').fill('password')
    await page.getByRole('button', { name: /login/i }).click()
    await expect(page).toHaveURL(/\/dashboard/, { timeout: 15000 })

    await page.getByTestId('user-menu').click()
    await page.getByText('Logout').click()
    await expect(page).toHaveURL('/login', { timeout: 10000 })
  })
})

test.describe('RBAC - Non-admin user', () => {
  test('should have limited menu for prod user', async ({ page }) => {
    await page.goto('/login')
    await page.getByPlaceholder('Email / NIP').fill('user@visualmesin.com')
    await page.getByPlaceholder('Password').fill('user123')
    await page.getByRole('button', { name: /login/i }).click()
    await expect(page).toHaveURL(/\/dashboard/, { timeout: 15000 })

    await expect(page.getByText('Administrasi')).not.toBeVisible()
  })
})
