import { Link, Outlet } from 'react-router-dom';
import { useEffect, useState } from 'react';
import { useAuthStore } from './stores/authStore';
import { useHabitsStore } from './stores/habitsStore';
import InstallPrompt from './pages/InstallPrompt';

export default function AppShell() {
  const { isAuthenticated, logout } = useAuthStore();
  const { queuedMutations } = useHabitsStore();
  const [updateReady, setUpdateReady] = useState(false);

  useEffect(() => {
    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.addEventListener('controllerchange', () => {
        setUpdateReady(true);
      });
    }
  }, []);

  return (
    <div className="app-shell">
      <header className="topbar">
        <Link to="/" className="brand">
          HabitFlow
        </Link>
        {updateReady && (
          <span className="pill">Update ready — refresh to apply</span>
        )}
        {queuedMutations.length > 0 && (
          <span className="pill">Offline queue: {queuedMutations.length}</span>
        )}
        <nav className="nav">
          <Link to="/app">Dashboard</Link>
          <Link to="/stats">Stats</Link>
          <InstallPrompt />
          {!isAuthenticated ? (
            <>
              <Link to="/login">Login</Link>
              <Link to="/register">Register</Link>
            </>
          ) : (
            <button type="button" onClick={() => logout()}>
              Logout
            </button>
          )}
        </nav>
      </header>
      <main className="content">
        <Outlet />
      </main>
    </div>
  );
}
