import { useAuth } from '../contexts/AuthContext';
import { useGroups } from '../hooks/useGroups';
import { useBalance } from '../hooks/useWallet';
import { Link } from 'react-router-dom';
import { UsersIcon, WalletIcon, ChartBarIcon, PlusIcon } from '@heroicons/react/24/outline';
import LoadingSpinner from '../components/LoadingSpinner';
import BankCard from '../components/BankCard';
import TransactionHistory from '../components/TransactionHistory';

const DashboardPage = () => {
  const { user } = useAuth()
  const { data: groups, isLoading: groupsLoading } = useGroups()
  const { data: balance, isLoading: balanceLoading } = useBalance(user?.wallet || '')

  const stats = [
    {
      name: 'Total Groups',
      value: groups?.length || 0,
      icon: <UsersIcon className="w-7 h-7" />,
      colorClass: 'bg-[#1a237e]',
    },
    {
      name: 'Wallet Balance',
      value: balanceLoading ? '...' : balance?.data?.balances?.[0]?.split(': ')[1] || '0 XLM',
      icon: <WalletIcon className="w-7 h-7" />,
      colorClass: 'bg-[#2ecc71]',
    },
    {
      name: 'Total Savings',
      value: '0 XLM',
      icon: <ChartBarIcon className="w-7 h-7" />,
      colorClass: 'bg-[#2ecc71]',
      progress: 60, // Example progress
    },
  ];

  const recentActivity: {
    type: 'contribution' | 'withdrawal';
    description: string;
    amount: string;
    time: string;
    group: string;
  }[] = [
    {
      type: 'contribution',
      description: 'Contributed to Alpha Chama',
      amount: '+50 XLM',
      time: '2 hours ago',
      group: 'Alpha Chama',
    },
    {
      type: 'withdrawal',
      description: 'Withdrew from Beta Group',
      amount: '-25 XLM',
      time: '1 day ago',
      group: 'Beta Group',
    },
  ];

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-[#1a237e]">Welcome, {user?.name}!</h1>
        <p className="text-gray-600 mt-2">Here's what's happening with your savings groups today.</p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        {stats.map((stat, idx) => (
          <BankCard key={idx} title={stat.name} value={stat.value} icon={stat.icon} colorClass={stat.colorClass} progress={stat.progress} />
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* My Groups */}
        <div className="bg-white rounded-2xl shadow-md p-6">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-semibold text-[#1a237e]">My Groups</h2>
            <Link to="/groups/create" className="inline-flex items-center px-4 py-2 bg-[#2ecc71] text-white rounded-lg font-semibold hover:bg-[#27ae60] transition">
              <PlusIcon className="w-5 h-5 mr-2" />
              Create Group
            </Link>
          </div>
          {groupsLoading ? (
            <div className="flex justify-center py-8">
              <LoadingSpinner />
            </div>
          ) : groups?.length === 0 ? (
            <div className="text-center py-8">
              <UsersIcon className="w-12 h-12 text-gray-400 mx-auto mb-4" />
              <p className="text-gray-600 mb-4">You haven't joined any groups yet</p>
              <Link to="/groups" className="inline-block px-4 py-2 border border-[#1a237e] text-[#1a237e] rounded-lg font-semibold hover:bg-[#f5f6fa] transition">Browse Groups</Link>
            </div>
          ) : (
            <div className="space-y-4">
              {groups?.slice(0, 3).map((group) => (
                <Link
                  key={group.ID}
                  to={`/groups/${group.ID}`}
                  className="block p-4 border border-gray-100 rounded-xl hover:border-[#2ecc71] hover:bg-emerald-50 transition-colors"
                >
                  <div className="flex items-center justify-between">
                    <div>
                      <h3 className="font-medium text-gray-900">{group.Name}</h3>
                      <p className="text-sm text-gray-600">{group.Description}</p>
                    </div>
                  </div>
                </Link>
              ))}
              {groups && groups.length > 3 && (
                <Link to="/groups" className="block text-center text-[#1a237e] hover:text-[#2ecc71] font-medium">View all groups</Link>
              )}
            </div>
          )}
        </div>

        {/* Recent Activity */}
        <TransactionHistory transactions={recentActivity} />
      </div>

      {/* Quick Actions */}
      <div className="mt-8">
        <h2 className="text-xl font-semibold text-[#1a237e] mb-4">Quick Actions</h2>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Link to="/groups/create" className="inline-flex flex-col items-center justify-center px-4 py-6 bg-white rounded-2xl shadow hover:bg-emerald-50 border border-gray-100 transition">
            <PlusIcon className="w-7 h-7 text-[#2ecc71] mb-2" />
            <span className="font-semibold text-[#1a237e]">Create Group</span>
          </Link>
          <Link to="/groups" className="inline-flex flex-col items-center justify-center px-4 py-6 bg-white rounded-2xl shadow hover:bg-emerald-50 border border-gray-100 transition">
            <UsersIcon className="w-7 h-7 text-[#1a237e] mb-2" />
            <span className="font-semibold text-[#1a237e]">Browse Groups</span>
          </Link>
          <Link to="/wallet" className="inline-flex flex-col items-center justify-center px-4 py-6 bg-white rounded-2xl shadow hover:bg-emerald-50 border border-gray-100 transition">
            <WalletIcon className="w-7 h-7 text-[#2ecc71] mb-2" />
            <span className="font-semibold text-[#1a237e]">Manage Wallet</span>
          </Link>
          <Link to="/transactions" className="inline-flex flex-col items-center justify-center px-4 py-6 bg-white rounded-2xl shadow hover:bg-emerald-50 border border-gray-100 transition">
            <ChartBarIcon className="w-7 h-7 text-[#1a237e] mb-2" />
            <span className="font-semibold text-[#1a237e]">View Transactions</span>
          </Link>
        </div>
      </div>
    </div>
  );
}

export default DashboardPage
