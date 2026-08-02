import { createBrowserRouter, Navigate } from 'react-router-dom';
import AppShell from './AppShell';
import AuthPage from './pages/AuthPage';
import DashboardPage from './pages/DashboardPage';
import HabitEditorPage from './pages/HabitEditorPage';
import HomePage from './pages/HomePage';
import StatsPage from './pages/StatsPage';
import { useAuthStore } from './stores/authStore';

const ProtectedRoute = ({ children }: { children: React.ReactElement }) => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  return isAuthenticated ? children : <Navigate to="/login" replace />;
};

export const router = createBrowserRouter([
  {
    path: '/',
    element: <AppShell />,
    children: [
      { index: true, element: <HomePage /> },
      { path: 'login', element: <AuthPage mode="login" /> },
      { path: 'register', element: <AuthPage mode="register" /> },
      {
        path: 'app',
        element: (
          <ProtectedRoute>
            <DashboardPage />
          </ProtectedRoute>
        ),
      },
      {
        path: 'habits/new',
        element: (
          <ProtectedRoute>
            <HabitEditorPage />
          </ProtectedRoute>
        ),
      },
      {
        path: 'habits/:id/edit',
        element: (
          <ProtectedRoute>
            <HabitEditorPage />
          </ProtectedRoute>
        ),
      },
      {
        path: 'stats',
        element: (
          <ProtectedRoute>
            <StatsPage />
          </ProtectedRoute>
        ),
      },
    ],
  },
]);
