
// import React from 'react';
import { HomeIcon, UsersIcon, WalletIcon, ChartBarIcon, Cog6ToothIcon, ArrowLeftOnRectangleIcon } from '@heroicons/react/24/outline';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const navItems = [
  { name: 'Dashboard', icon: HomeIcon, to: '/dashboard' },
  { name: 'Groups', icon: UsersIcon, to: '/groups' },
  { name: 'Wallet', icon: WalletIcon, to: '/wallet' },
  { name: 'Activity', icon: ChartBarIcon, to: '/transactions' },
  { name: 'Settings', icon: Cog6ToothIcon, to: '/settings' },
];

const Sidebar = () => {
  const location = useLocation();
  const { logout } = useAuth();
  return (
    <aside className="h-screen w-64 bg-[#1a237e] flex flex-col py-8 px-4 shadow-lg">
      <div className="mb-10 flex items-center justify-center">
        <span className="text-2xl font-bold text-white tracking-wide">Chama Wallet</span>
      </div>
      <nav className="flex-1">
        <ul className="space-y-2">
          {navItems.map(({ name, icon: Icon, to }) => (
            <li key={name}>
              <Link
                to={to}
                className={`flex items-center px-4 py-3 rounded-lg transition-colors duration-200 text-white hover:bg-[#2ecc71] hover:text-[#1a237e] ${location.pathname === to ? 'bg-[#2ecc71] text-[#1a237e]' : ''}`}
              >
                <Icon className="h-6 w-6 mr-3" />
                <span className="font-medium">{name}</span>
              </Link>
            </li>
          ))}
        </ul>
      </nav>
      <div className="mt-8 flex flex-col items-center gap-3">
        <button
          onClick={logout}
          className="flex items-center w-full px-4 py-3 rounded-lg text-white hover:bg-red-100 hover:text-red-700 transition-colors duration-200 font-medium focus:outline-none"
        >
          <ArrowLeftOnRectangleIcon className="h-6 w-6 mr-3" />
          Logout
        </button>
        <span className="text-xs text-gray-300">Secured by Chama Wallet</span>
      </div>
    </aside>
  );
};

export default Sidebar;
