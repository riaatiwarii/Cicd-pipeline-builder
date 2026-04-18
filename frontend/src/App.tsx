import React, { useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider, createTheme, CssBaseline, AppBar, Toolbar, Box, Button } from '@mui/material';
import LogoutIcon from '@mui/icons-material/Logout';

import { LoginPage } from './pages/LoginPage';
import { PipelineListPage } from './pages/PipelineListPage';
import { PipelineDetailPage } from './pages/PipelineDetailPage';
import { BuildDetailPage } from './pages/BuildDetailPage';
import { useAuthStore } from './store/store';

const theme = createTheme({
  palette: {
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
    mode: 'light',
  },
  typography: {
    fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
  },
});

interface ProtectedRouteProps {
  children: React.ReactNode;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
};

const Layout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated, logout } = useAuthStore();

  const handleLogout = () => {
    logout();
    window.location.href = '/login';
  };

  if (!isAuthenticated) {
    return <>{children}</>;
  }

  return (
    <>
      <AppBar position="static">
        <Toolbar>
          <Box sx={{ flexGrow: 1 }}>
            CI/CD Pipeline Builder
          </Box>
          <Button color="inherit" onClick={handleLogout} endIcon={<LogoutIcon />}>
            Logout
          </Button>
        </Toolbar>
      </AppBar>
      {children}
    </>
  );
};

export const App: React.FC = () => {
  useEffect(() => {
    // Initialize auth state from localStorage
  }, []);

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <BrowserRouter>
        <Layout>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route
              path="/pipelines"
              element={
                <ProtectedRoute>
                  <PipelineListPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/pipelines/:id"
              element={
                <ProtectedRoute>
                  <PipelineDetailPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/builds/:buildId"
              element={
                <ProtectedRoute>
                  <BuildDetailPage />
                </ProtectedRoute>
              }
            />
            <Route path="/" element={<Navigate to="/pipelines" replace />} />
          </Routes>
        </Layout>
      </BrowserRouter>
    </ThemeProvider>
  );
};

export default App;
