import React, { Suspense } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import {ClientProvider} from "./context/ClientContext.tsx";
import {AuthProvider} from "./context/AuthContext.tsx";
import { NotifProvider } from "./context/NotifContext";
import FallbackLoading from "./pages/FallbackLoading";

const AppLazy = React.lazy(() => import('./App.tsx'))

createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
      <ClientProvider>
          <AuthProvider>
              <NotifProvider>
                  <Suspense fallback={<FallbackLoading />}>
                      <AppLazy />
                  </Suspense>
              </NotifProvider>
          </AuthProvider>
      </ClientProvider>
  </React.StrictMode>,
)