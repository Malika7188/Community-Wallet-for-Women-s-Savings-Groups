import { useState } from 'react';
import { Link, Navigate, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { Eye, EyeOff } from 'lucide-react';
import { WalletIcon } from '@heroicons/react/24/outline';
import toast from 'react-hot-toast';

const LoginPage = () => {
  const { user, login, isLoading } = useAuth()
  const navigate = useNavigate()
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  })
  const [showPassword, setShowPassword] = useState(false)

  if (user) {
    return <Navigate to="/dashboard" replace />
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await login(formData.email, formData.password)
      toast.success('Login successful!')
      navigate('/dashboard')
    } catch (error) {
      toast.error('Invalid email or password')
    }
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData(prev => ({
      ...prev,
      [e.target.name]: e.target.value
    }))
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-[#e0e7ef] to-[#f5f6fa] font-inter">
      <div className="w-full max-w-xl">
        <div className="rounded-3xl shadow-2xl bg-white/80 backdrop-blur-lg border border-gray-100 px-16 py-16 flex flex-col items-center animate-fade-in">
          <div className="w-15 h-0 bg-gradient-to-br from-[#1a237e] to-[#2ecc71] rounded-3xl flex items-center justify-center mb-6 shadow-xl">
            <WalletIcon className="w-10 h-10 text-white" />
          </div>
          <h2 className="text-5xl font-black text-[#1a237e] tracking-tight mb-2" style={{ fontFamily: 'Inter, Roboto, sans-serif' }}>
            Welcome back
          </h2>
          <p className="mb-8 text-gray-600 font-semibold text-xl">Sign in to your Chama Wallet account</p>
          <form className="w-full space-y-7" onSubmit={handleSubmit}>
            <div>
              <label htmlFor="email" className="block text-lg font-bold text-[#1a237e] mb-2">Email address</label>
              <input
                id="email"
                name="email"
                type="email"
                required
                className="w-full rounded-2xl border border-gray-200 px-6 py-4 text-lg focus:ring-2 focus:ring-[#2ecc71] focus:outline-none bg-white/90 font-inter"
                placeholder="Enter your email"
                value={formData.email}
                onChange={handleChange}
              />
            </div>
            <div>
              <label htmlFor="password" className="block text-lg font-bold text-[#1a237e] mb-2">Password</label>
              <div className="relative">
                <input
                  id="password"
                  name="password"
                  type={showPassword ? 'text' : 'password'}
                  required
                  className="w-full rounded-2xl border border-gray-200 px-6 py-4 text-lg pr-12 focus:ring-2 focus:ring-[#2ecc71] focus:outline-none bg-white/90 font-inter"
                  placeholder="Enter your password"
                  value={formData.password}
                  onChange={handleChange}
                />
                <button
                  type="button"
                  className="absolute inset-y-0 right-0 pr-4 flex items-center"
                  onClick={() => setShowPassword(!showPassword)}
                >
                  {showPassword ? (
                    <EyeOff className="h-6 w-6 text-gray-400" />
                  ) : (
                    <Eye className="h-6 w-6 text-gray-400" />
                  )}
                </button>
              </div>
            </div>
            <button
              type="submit"
              disabled={isLoading}
              className="w-full py-4 rounded-2xl bg-[#2ecc71] text-white font-extrabold text-2xl shadow-lg hover:bg-[#27ae60] transition-colors duration-200 font-inter"
            >
              {isLoading ? 'Signing in...' : 'Sign in'}
            </button>
            <div className="text-center mt-4">
              <span className="text-gray-600 text-lg">Don't have an account? </span>
              <Link to="/signup" className="font-bold text-[#1a237e] hover:text-[#2ecc71] transition-colors text-lg">Sign up</Link>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}
export default LoginPage