import { Link, Outlet } from 'react-router-dom';
import { useAuthStore } from './stores/authStore';

export default function AppShell() {
  const { isAuthenticated, logout } = useAuthStore();

  return (
    <div className="app-shell">
      <header className="topbar">
        <Link to="/" className="brand">
          HabitFlow
        </Link>
        <nav className="nav">
          <Link to="/app">Dashboard</Link>
          <Link to="/stats">Stats</Link>
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
