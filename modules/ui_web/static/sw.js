// Service Worker para Digna PWA
// Cache First strategy para templates e assets estáticos

const CACHE_NAME = 'digna-v1';
const STATIC_ASSETS = [
  '/',
  '/pdv',
  '/social',
  '/dashboard',
  '/static/manifest.json'
];

// Instalação: cache assets estáticos
self.addEventListener('install', (event) => {
  console.log('[SW] Installing...');
  
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then((cache) => {
        console.log('[SW] Caching static assets');
        return cache.addAll(STATIC_ASSETS);
      })
      .catch((err) => {
        console.log('[SW] Cache failed:', err);
      })
  );
  
  // Ativar imediatamente
  self.skipWaiting();
});

// Ativação: limpar caches antigos
self.addEventListener('activate', (event) => {
  console.log('[SW] Activating...');
  
  event.waitUntil(
    caches.keys()
      .then((cacheNames) => {
        return Promise.all(
          cacheNames
            .filter((name) => name !== CACHE_NAME)
            .map((name) => caches.delete(name))
        );
      })
  );
  
  self.clients.claim();
});

// Fetch: estratégia Cache First para GET, Network First para API
self.addEventListener('fetch', (event) => {
  const { request } = event;
  const url = new URL(request.url);
  
  // Não interceptar requests de API (POST, etc)
  if (request.method !== 'GET') {
    return;
  }
  
  // Não cachear chamadas de API dinâmicas
  if (url.pathname.startsWith('/api/')) {
    return;
  }
  
  event.respondWith(
    caches.match(request)
      .then((cachedResponse) => {
        if (cachedResponse) {
          // Retorna do cache mas atualiza em background
          fetch(request)
            .then((networkResponse) => {
              caches.open(CACHE_NAME).then((cache) => {
                cache.put(request, networkResponse.clone());
              });
            })
            .catch(() => {
              // Falha na rede, mantém cache
            });
          
          return cachedResponse;
        }
        
        // Não está no cache, busca na rede
        return fetch(request)
          .then((networkResponse) => {
            // Cacheia a resposta
            const responseClone = networkResponse.clone();
            caches.open(CACHE_NAME).then((cache) => {
              cache.put(request, responseClone);
            });
            return networkResponse;
          })
          .catch(() => {
            // Offline e não está no cache
            return new Response(
              '<h1>Offline</h1><p>Digna está offline. Algumas funcionalidades podem não estar disponíveis.</p>',
              { headers: { 'Content-Type': 'text/html' } }
            );
          });
      })
  );
});

// Sync em background (quando online novamente)
self.addEventListener('sync', (event) => {
  if (event.tag === 'sync-sales') {
    console.log('[SW] Background sync: sales');
    // Aqui seria implementada a sincronização de vendas pendentes
  }
});

// Push notifications (futuro)
self.addEventListener('push', (event) => {
  const options = {
    body: event.data.text(),
    icon: '/static/icon-192x192.png',
    badge: '/static/icon-72x72.png'
  };
  
  event.waitUntil(
    self.registration.showNotification('Digna', options)
  );
});
