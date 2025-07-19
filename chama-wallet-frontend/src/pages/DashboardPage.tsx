import { useAuth } from '../contexts/AuthContext'
import { useGroups } from '../hooks/useGroups'
import { useBalance } from '../hooks/useWallet'
import { Link } from 'react-router-dom'
import { Users, Wallet, TrendingUp, Plus, ArrowUpRight, ArrowDownRight } from 'lucide-react'
import LoadingSpinner from '../components/LoadingSpinner'

const DashboardPage = () => {
  const { user } = useAuth()
  const { data: groups, isLoading: groupsLoading } = useGroups()
  const { data: balance, isLoading: balanceLoading } = useBalance(user?.wallet || '')

  const stats = [
    {
      name: 'Total Groups',
      value: groups?.length || 0,
      icon: Users,
      color: 'bg-blue-500',
    },
    {
      name: 'Wallet Balance',
      value: balanceLoading ? '...' : balance?.data?.balances?.[0]?.split(': ')[1] || '0 XLM',
      icon: Wallet,
      color: 'bg-green-500',
    },
    {
      name: 'Total Savings',
      value: '0 XLM',
      icon: TrendingUp,
      color: 'bg-purple-500',
    },
  ]

  const recentActivity = [
    {
      type: 'contribution',
      description: 'Contributed to Alpha Chama',
      amount: '+50 XLM',
      time: '2 hours ago',
      icon: ArrowUpRight,
      color: 'text-green-600',
    },
    {
      type: 'withdrawal',
      description: 'Withdrew from Beta Group',
      amount: '-25 XLM',
      time: '1 day ago',
      icon: ArrowDownRight,
      color: 'text-red-600',
    },
  ]

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">
          Welcome, {user?.name}!
        </h1>
        <p className="text-gray-600 mt-2">
          Here's what's happening with your savings groups today.
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        {stats.map((stat, index) => {
          const Icon = stat.icon
          return (
            <div key={index} className="card">
              <div className="flex items-center">
                <div className={`w-12 h-12 ${stat.color} rounded-lg flex items-center justify-center`}>
                  <Icon className="w-6 h-6 text-white" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">{stat.name}</p>
                  <p className="text-2xl font-bold text-gray-900">{stat.value}</p>
                </div>
              </div>
            </div>
          )
        })}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* My Groups */}
        <div className="card">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-semibold text-gray-900">My Groups</h2>
            <Link to="/groups/create" className="btn btn-primary">
              <Plus className="w-4 h-4 mr-2" />
              Create Group
            </Link>
          </div>
          
          {groupsLoading ? (
            <div className="flex justify-center py-8">
              <LoadingSpinner />
            </div>
          ) : groups?.length === 0 ? (
            <div className="text-center py-8">
              <Users className="w-12 h-12 text-gray-400 mx-auto mb-4" />
              <p className="text-gray-600 mb-4">You haven't joined any groups yet</p>
              <Link to="/groups" className="btn btn-outline">
                Browse Groups
              </Link>
            </div>
          ) : (
            <div className="space-y-4">
              {groups?.slice(0, 3).map((group) => (
                <Link
                  key={group.ID}
                  to={`/groups/${group.ID}`}
                  className="block p-4 border border-gray-200 rounded-lg hover:border-primary-300 hover:bg-primary-50 transition-colors"
                >
                  <div className="flex items-center justify-between">
                    <div>
                      <h3 className="font-medium text-gray-900">{group.Name}</h3>
                      <p className="text-sm text-gray-600">{group.Description}</p>
                    </div>
                    <ArrowUpRight className="w-5 h-5 text-gray-400" />
                  </div>
                </Link>
              ))}
              {groups && groups.length > 3 && (
                <Link to="/groups" className="block text-center text-primary-600 hover:text-primary-700 font-medium">
                  View all groups
                </Link>
              )}
            </div>
          )}
        </div>

        {/* Recent Activity */}
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 mb-6">Recent Activity</h2>
          
          {recentActivity.length === 0 ? (
            <div className="text-center py-8">
              <TrendingUp className="w-12 h-12 text-gray-400 mx-auto mb-4" />
              <p className="text-gray-600">No recent activity</p>
            </div>
          ) : (
            <div className="space-y-4">
              {recentActivity.map((activity, index) => {
                const Icon = activity.icon
                return (
                  <div key={index} className="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
                    <div className="flex items-center">
                      <div className={`w-10 h-10 bg-gray-100 rounded-lg flex items-center justify-center`}>
                        <Icon className={`w-5 h-5 ${activity.color}`} />
                      </div>
                      <div className="ml-3">
                        <p className="font-medium text-gray-900">{activity.description}</p>
                        <p className="text-sm text-gray-600">{activity.time}</p>
                      </div>
                    </div>
                    <span className={`font-semibold ${activity.color}`}>
                      {activity.amount}
                    </span>
                  </div>
                )
              })}
            </div>
          )}
        </div>
      </div>

      {/* Quick Actions */}
      <div className="mt-8">
        <h2 className="text-xl font-semibold text-gray-900 mb-4">Quick Actions</h2>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Link to="/groups/create" className="btn btn-outline text-center">
            <Plus className="w-5 h-5 mx-auto mb-2" />
            Create Group
          </Link>
          <Link to="/groups" className="btn btn-outline text-center">
            <Users className="w-5 h-5 mx-auto mb-2" />
            Browse Groups
          </Link>
          <Link to="/wallet" className="btn btn-outline text-center">
            <Wallet className="w-5 h-5 mx-auto mb-2" />
            Manage Wallet
          </Link>
          <Link to="/transactions" className="btn btn-outline text-center">
            <TrendingUp className="w-5 h-5 mx-auto mb-2" />
            View Transactions
          </Link>
        </div>
      </div>
    </div>
  )
}

export default DashboardPage
