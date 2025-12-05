
import React from 'react';
import { HomeIcon, UsersIcon, WalletIcon, ChartBarIcon, ArrowLeftOnRectangleIcon, Bars3Icon } from '@heroicons/react/24/outline';
import NotificationCenter from './NotificationCenter';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const navItems = [
  { name: 'Dashboard', icon: HomeIcon, to: '/dashboard' },
  { name: 'Groups', icon: UsersIcon, to: '/groups' },
  { name: 'Wallet', icon: WalletIcon, to: '/wallet' },
  { name: 'Activity', icon: ChartBarIcon, to: '/transactions' },
  // { name: 'Settings', icon: Cog6ToothIcon, to: '/settings' },
];

interface SidebarProps {
  isCollapsed: boolean;
  setIsCollapsed: React.Dispatch<React.SetStateAction<boolean>>;
}

const Sidebar: React.FC<SidebarProps> = ({ isCollapsed, setIsCollapsed }) => {
  const location = useLocation();
  const { user, logout } = useAuth();

  const toggleSidebar = () => {
    setIsCollapsed(!isCollapsed);
  };

  return (
    <aside className={`fixed top-0 left-0 h-screen bg-[#1a237e] flex flex-col py-8 px-4 shadow-lg z-40 transition-width duration-300 ${isCollapsed ? 'w-20' : 'w-64'}`}>
      <div className="mb-10 flex items-center justify-between">
  {!isCollapsed && <span className="text-2xl font-bold text-white tracking-wide">Umoja Wallet</span>}
        <button onClick={toggleSidebar} className="text-white hover:text-[#2ecc71]">
          <Bars3Icon className="h-8 w-8" />
        </button>
      </div>
      <nav className="flex-1">
        <ul className="space-y-2">
          {navItems.map(({ name, icon: Icon, to }) => (
            <li key={name}>
              <Link
                to={to}
                className={`flex items-center px-4 py-3 rounded-lg transition-colors duration-200 text-white hover:bg-[#2ecc71] hover:text-[#1a237e] ${location.pathname === to ? 'bg-[#2ecc71] text-[#1a237e]' : ''}`}
              >
                <Icon className="h-6 w-6" />
                {!isCollapsed && <span className="font-medium ml-3">{name}</span>}
              </Link>
            </li>
          ))}
          <li>
            <NotificationCenter isCollapsed={isCollapsed} />
          </li>
        </ul>
      </nav>
      <div className="mt-8 flex flex-col items-center gap-3">
        <div className="flex items-center w-full px-4 py-3">
          <img src={`https://ui-avatars.com/api/?name=${user?.name}&background=random`} alt="User" className="h-8 w-8 rounded-full mr-3" />
          {!isCollapsed && <span className="text-white font-medium">{user?.name}</span>}
        </div>
        <button
          onClick={logout}
          className="flex items-center w-full px-4 py-3 rounded-lg text-white hover:bg-red-100 hover:text-red-700 transition-colors duration-200 font-medium focus:outline-none"
        >
          <ArrowLeftOnRectangleIcon className="h-6 w-6" />
          {!isCollapsed && <span className="ml-3">Logout</span>}
        </button>
  {!isCollapsed && <span className="text-xs text-gray-300">Secured by Umoja Wallet</span>}
      </div>
    </aside>
  );
};

export default Sidebar;
