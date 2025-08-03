import { Routes, Route } from 'react-router-dom'
import { AuthProvider } from './contexts/AuthContext'
import Layout from './components/Layout';
import HomePage from './pages/HomePage'
import LoginPage from './pages/LoginPage'
import SignUpPage from './pages/SignUpPage'
import DashboardPage from './pages/DashboardPage'
import GroupsPage from './pages/GroupsPage'
import GroupDetailPage from './pages/GroupDetailPage'
import CreateGroupPage from './pages/CreateGroupPage'
import WalletPage from './pages/WalletPage'
import TransactionsPage from './pages/TransactionsPage'
import ProtectedRoute from './components/ProtectedRoute'

import PublicLayout from './components/PublicLayout';

function App() {
  return (
    <AuthProvider>
      <Routes>
        {/* Public pages */}
        <Route element={<PublicLayout />}>
          <Route path="/" element={<HomePage />} />
          <Route path="login" element={<LoginPage />} />
          <Route path="signup" element={<SignUpPage />} />
        </Route>
        {/* Protected pages with dashboard layout */}
        <Route element={<Layout />}>
          <Route path="dashboard" element={<ProtectedRoute><DashboardPage /></ProtectedRoute>} />
          <Route path="groups" element={<ProtectedRoute><GroupsPage /></ProtectedRoute>} />
          <Route path="groups/create" element={<ProtectedRoute><CreateGroupPage /></ProtectedRoute>} />
          <Route path="groups/:id" element={<ProtectedRoute><GroupDetailPage /></ProtectedRoute>} />
          <Route path="wallet" element={<ProtectedRoute><WalletPage /></ProtectedRoute>} />
          <Route path="transactions" element={<ProtectedRoute><TransactionsPage /></ProtectedRoute>} />
        </Route>
        {/* Catch all route for 404 */}
        <Route path="*" element={<div className="min-h-screen flex items-center justify-center">
          <div className="text-center">
            <h1 className="text-4xl font-bold text-gray-900 mb-4">404</h1>
            <p className="text-gray-600 mb-8">Page not found</p>
            <a href="/" className="btn btn-primary">Go Home</a>
          </div>
        </div>} />
      </Routes>
    </AuthProvider>
  );
}

export default App