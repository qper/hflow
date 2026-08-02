import React from 'react';
import ReactDOM from 'react-dom/client';
import { RouterProvider } from 'react-router-dom';
import './styles.css';
import { router } from './router';
import { registerServiceWorker } from './lib/pwa';
import { syncQueuedMutations } from './lib/sync';

registerServiceWorker();
window.addEventListener('online', () => {
  void syncQueuedMutations();
});
void syncQueuedMutations();

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
);
