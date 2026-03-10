/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./modules/**/templates/**/*.html",
    "./modules/**/templates/**/*.js"
  ],
  theme: {
    extend: {
      colors: {
        // Paleta de Cores Oficial Digna - "Soberania e Suor"
        'digna-green': '#2E7D32',    // Verde Trabalho - ações positivas, ganhos, valor do trabalho
        'digna-ocher': '#F57F17',    // Ocre Comunidade - fundos de reserva, alertas, elementos de destaque
        'digna-blue': '#1565C0',     // Azul Soberania - links, botões institucionais, insumos
        'digna-bg': '#F9F9F6',       // Fundo Off-white - fundo global para evitar fadiga visual
        'digna-text': '#212121',     // Texto Principal - contraste suave (não preto absoluto)
        
        // Cores semânticas baseadas na paleta Digna
        'primary': '#2E7D32',        // Verde Trabalho como cor primária
        'secondary': '#F57F17',      // Ocre Comunidade como cor secundária
        'accent': '#1565C0',         // Azul Soberania como cor de destaque
      },
      fontFamily: {
        'sans': ['Inter', 'Ubuntu', 'system-ui', 'sans-serif'],
      },
      spacing: {
        'touch': '44px',             // Tamanho mínimo para área de toque mobile
      },
      minHeight: {
        'touch': '44px',             // Altura mínima para botões mobile
      },
    },
  },
  plugins: [],
}