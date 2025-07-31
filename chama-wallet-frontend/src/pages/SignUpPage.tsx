import { useState } from 'react';
import { Link, Navigate, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { Eye, EyeOff } from 'lucide-react';
import { WalletIcon } from '@heroicons/react/24/outline';
import toast from 'react-hot-toast';

const SignUpPage = () => {
  const { user, signup, isLoading } = useAuth();
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    confirmPassword: '',
  });
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  if (user) {
    return <Navigate to="/dashboard" replace />;
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (formData.password !== formData.confirmPassword) {
      toast.error('Passwords do not match');
      return;
    }

    if (formData.password.length < 6) {
      toast.error('Password must be at least 6 characters');
      return;
    }

    try {
      await signup(formData.email, formData.password, formData.name);
      toast.success('Account created successfully!');
      navigate('/dashboard');
    } catch (error) {
      toast.error('Failed to create account');
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData(prev => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-[#e0e7ef] to-[#f5f6fa] font-inter px-2">
      <div className="w-full max-w-xl">
        <div className="rounded-3xl shadow-xl bg-white/90 backdrop-blur-md border border-gray-100 px-6 py-6 sm:px-8 sm:py-10 flex flex-col items-center animate-fade-in w-full">
          
          {/* Icon */}
          <div className="w-16 h-16 bg-gradient-to-br from-[#1a237e] to-[#2ecc71] rounded-2xl flex items-center justify-center mb-4 shadow-lg">
            <WalletIcon className="w-8 h-8 text-white" />
          </div>

          {/* Heading */}
          <h2 className="text-4xl font-bold text-[#1a237e] mb-2 text-center">
            Create your account
          </h2>
          <p className="mb-6 text-gray-600 text-lg text-center">
            Join the future of community savings
          </p>

          {/* Form */}
          <form className="w-full space-y-6" onSubmit={handleSubmit}>
            {/* Name */}
            <div>
              <label htmlFor="name" className="block text-lg font-semibold text-[#1a237e] mb-2">
                Full Name
              </label>
              <input
                id="name"
                name="name"
                type="text"
                required
                className="w-full rounded-2xl border border-gray-200 px-5 py-3 text-lg focus:ring-2 focus:ring-[#2ecc71] focus:outline-none bg-white/90"
                placeholder="Enter your full name"
                value={formData.name}
                onChange={handleChange}
              />
            </div>

            {/* Email */}
            <div>
              <label htmlFor="email" className="block text-lg font-semibold text-[#1a237e] mb-2">
                Email Address
              </label>
              <input
                id="email"
                name="email"
                type="email"
                required
                className="w-full rounded-2xl border border-gray-200 px-5 py-3 text-lg focus:ring-2 focus:ring-[#2ecc71] focus:outline-none bg-white/90"
                placeholder="Enter your email"
                value={formData.email}
                onChange={handleChange}
              />
            </div>

            {/* Password */}
            <div>
              <label htmlFor="password" className="block text-lg font-semibold text-[#1a237e] mb-2">
                Password
              </label>
              <div className="relative">
                <input
                  id="password"
                  name="password"
                  type={showPassword ? 'text' : 'password'}
                  required
                  className="w-full rounded-2xl border border-gray-200 px-5 py-3 text-lg pr-10 focus:ring-2 focus:ring-[#2ecc71] focus:outline-none bg-white/90"
                  placeholder="Create a password"
                  value={formData.password}
                  onChange={handleChange}
                />
                <button
                  type="button"
                  className="absolute inset-y-0 right-0 pr-3 flex items-center"
                  onClick={() => setShowPassword(!showPassword)}
                >
                  {showPassword ? <EyeOff className="h-5 w-5 text-gray-400" /> : <Eye className="h-5 w-5 text-gray-400" />}
                </button>
              </div>
            </div>

            {/* Confirm Password */}
            <div>
              <label htmlFor="confirmPassword" className="block text-lg font-semibold text-[#1a237e] mb-2">
                Confirm Password
              </label>
              <div className="relative">
                <input
                  id="confirmPassword"
                  name="confirmPassword"
                  type={showConfirmPassword ? 'text' : 'password'}
                  required
                  className="w-full rounded-2xl border border-gray-200 px-5 py-3 text-lg pr-10 focus:ring-2 focus:ring-[#2ecc71] focus:outline-none bg-white/90"
                  placeholder="Confirm your password"
                  value={formData.confirmPassword}
                  onChange={handleChange}
                />
                <button
                  type="button"
                  className="absolute inset-y-0 right-0 pr-3 flex items-center"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                >
                  {showConfirmPassword ? <EyeOff className="h-5 w-5 text-gray-400" /> : <Eye className="h-5 w-5 text-gray-400" />}
                </button>
              </div>
            </div>

            {/* Terms */}
            <div className="flex items-center text-base">
              <input
                id="terms"
                name="terms"
                type="checkbox"
                required
                className="h-4 w-4 text-[#2ecc71] focus:ring-[#2ecc71] border-gray-300 rounded"
              />
              <label htmlFor="terms" className="ml-2 text-gray-700">
                I agree to the{' '}
                <a href="#" className="text-[#2ecc71] hover:text-[#1a237e]">Terms of Service</a> and{' '}
                <a href="#" className="text-[#2ecc71] hover:text-[#1a237e]">Privacy Policy</a>
              </label>
            </div>

            {/* Submit */}
            <button
              type="submit"
              disabled={isLoading}
              className="w-full py-3 rounded-2xl bg-[#2ecc71] text-white font-bold text-xl shadow-md hover:bg-[#27ae60] transition duration-200"
            >
              {isLoading ? 'Creating account...' : 'Create Account'}
            </button>

            {/* Footer */}
            <div className="text-center mt-4">
              <span className="text-gray-600 text-lg">Already have an account? </span>
              <Link to="/login" className="font-bold text-[#1a237e] hover:text-[#2ecc71] transition-colors text-lg">Sign in</Link>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default SignUpPage;
