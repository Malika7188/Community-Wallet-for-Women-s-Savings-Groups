import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { useGroups, useGroupBalance, useGroupMutations } from '../hooks/useGroups'
import { useAuth } from '../contexts/AuthContext'
import { 
  Users, 
  Wallet, 
  Plus, 
  ArrowLeft, 
  Copy, 
  ExternalLink,
  TrendingUp,
  UserPlus
} from 'lucide-react'
import LoadingSpinner from '../components/LoadingSpinner'
import toast from 'react-hot-toast'

const GroupDetailPage = () => {
  const { id } = useParams<{ id: string }>()
  const { user } = useAuth()
  const { data: groups, isLoading: groupsLoading } = useGroups()
  const { data: balance, isLoading: balanceLoading } = useGroupBalance(id!)
  const { joinGroup, contributeToGroup } = useGroupMutations()
  
  const [showContributeModal, setShowContributeModal] = useState(false)
  const [showJoinModal, setShowJoinModal] = useState(false)
  const [contributeAmount, setContributeAmount] = useState('')
  const [secretKey, setSecretKey] = useState('')

  const group = groups?.data?.find(g => g.ID === id)

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    toast.success('Copied to clipboard!')
  }

  const handleJoinGroup = async () => {
    if (!user?.wallet) return
    
    try {
      await joinGroup.mutateAsync({
        id: id!,
        data: { wallet: user.wallet }
      })
      setShowJoinModal(false)
    } catch (error) {
      // Error handled by mutation
    }
  }

  const handleContribute = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!user?.wallet || !contributeAmount || !secretKey) return

    try {
      await contributeToGroup.mutateAsync({
        id: id!,
        data: {
          from: user.wallet,
          secret: secretKey,
          amount: contributeAmount
        }
      })
      setShowContributeModal(false)
      setContributeAmount('')
      setSecretKey('')
    } catch (error) {
      // Error handled by mutation
    }
  }

  if (groupsLoading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <LoadingSpinner size="lg" />
      </div>
    )
  }

  if (!group) {
    return (
      <div className="max-w-2xl mx-auto px-4 py-16 text-center">
        <h1 className="text-2xl font-bold text-gray-900 mb-4">Group Not Found</h1>
        <p className="text-gray-600 mb-8">The group you're looking for doesn't exist.</p>
        <Link to="/groups" className="btn btn-primary">
          Back to Groups
        </Link>
      </div>
    )
  }

  const isMember = group.Members?.some(member => member.Wallet === user?.wallet)
  const groupBalance = balance?.data?.balance || '0'

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <Link 
          to="/groups" 
          className="inline-flex items-center text-gray-600 hover:text-gray-900 mb-4"
        >
          <ArrowLeft className="w-4 h-4 mr-2" />
          Back to Groups
        </Link>
      </div>

      {/* Group Header */}
      <div className="card mb-8">
        <div className="flex flex-col md:flex-row md:items-center md:justify-between">
          <div className="flex items-center mb-4 md:mb-0">
            <div className="w-16 h-16 bg-gradient-to-r from-stellar-500 to-primary-600 rounded-2xl flex items-center justify-center mr-4">
              <Users className="w-8 h-8 text-white" />
            </div>
            <div>
              <h1 className="text-3xl font-bold text-gray-900">{group.Name}</h1>
              <p className="text-gray-600">{group.Description}</p>
            </div>
          </div>
          
          <div className="flex flex-col sm:flex-row gap-3">
            {!isMember ? (
              <button
                onClick={() => setShowJoinModal(true)}
                className="btn btn-primary"
              >
                <UserPlus className="w-4 h-4 mr-2" />
                Join Group
              </button>
            ) : (
              <button
                onClick={() => setShowContributeModal(true)}
                className="btn btn-primary"
              >
                <Plus className="w-4 h-4 mr-2" />
                Contribute
              </button>
            )}
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Group Stats */}
        <div className="lg:col-span-2 space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="card">
              <div className="flex items-center">
                <div className="w-12 h-12 bg-blue-500 rounded-lg flex items-center justify-center">
                  <Users className="w-6 h-6 text-white" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">Members</p>
                  <p className="text-2xl font-bold text-gray-900">
                    {group.Members?.length || 0}
                  </p>
                </div>
              </div>
            </div>

            <div className="card">
              <div className="flex items-center">
                <div className="w-12 h-12 bg-green-500 rounded-lg flex items-center justify-center">
                  <Wallet className="w-6 h-6 text-white" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">Total Balance</p>
                  <p className="text-2xl font-bold text-gray-900">
                    {balanceLoading ? '...' : `${groupBalance} XLM`}
                  </p>
                </div>
              </div>
            </div>

            <div className="card">
              <div className="flex items-center">
                <div className="w-12 h-12 bg-purple-500 rounded-lg flex items-center justify-center">
                  <TrendingUp className="w-6 h-6 text-white" />
                </div>
                <div className="ml-4">
                  <p className="text-sm font-medium text-gray-600">Contributions</p>
                  <p className="text-2xl font-bold text-gray-900">
                    {group.Contributions?.length || 0}
                  </p>
                </div>
              </div>
            </div>
          </div>

          {/* Members List */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Members</h3>
            {group.Members && group.Members.length > 0 ? (
              <div className="space-y-3">
                {group.Members.map((member, index) => (
                  <div key={member.ID} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                    <div className="flex items-center">
                      <div className="w-10 h-10 bg-primary-100 rounded-full flex items-center justify-center">
                        <span className="text-primary-600 font-medium">
                          {index + 1}
                        </span>
                      </div>
                      <div className="ml-3">
                        <p className="font-medium text-gray-900">
                          Member {index + 1}
                        </p>
                        <p className="text-sm text-gray-600 font-mono">
                          {member.Wallet.slice(0, 8)}...{member.Wallet.slice(-8)}
                        </p>
                      </div>
                    </div>
                    <button
                      onClick={() => copyToClipboard(member.Wallet)}
                      className="p-2 text-gray-400 hover:text-gray-600"
                    >
                      <Copy className="w-4 h-4" />
                    </button>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-8">
                <Users className="w-12 h-12 text-gray-400 mx-auto mb-4" />
                <p className="text-gray-600">No members yet</p>
              </div>
            )}
          </div>
        </div>

        {/* Group Info Sidebar */}
        <div className="space-y-6">
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Group Wallet</h3>
            <div className="space-y-3">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Wallet Address
                </label>
                <div className="flex items-center space-x-2">
                  <code className="flex-1 text-xs bg-gray-100 p-2 rounded font-mono break-all">
                    {group.Wallet}
                  </code>
                  <button
                    onClick={() => copyToClipboard(group.Wallet)}
                    className="p-2 text-gray-400 hover:text-gray-600"
                  >
                    <Copy className="w-4 h-4" />
                  </button>
                </div>
              </div>
              
              <a
                href={`https://stellar.expert/explorer/testnet/account/${group.Wallet}`}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center text-sm text-primary-600 hover:text-primary-700"
              >
                View on Stellar Explorer
                <ExternalLink className="w-4 h-4 ml-1" />
              </a>
            </div>
          </div>

          {isMember && (
            <div className="card">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">Quick Actions</h3>
              <div className="space-y-3">
                <button
                  onClick={() => setShowContributeModal(true)}
                  className="btn btn-primary w-full"
                >
                  <Plus className="w-4 h-4 mr-2" />
                  Make Contribution
                </button>
                <Link to="/transactions" className="btn btn-outline w-full">
                  <TrendingUp className="w-4 h-4 mr-2" />
                  View Transactions
                </Link>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Join Group Modal */}
      {showJoinModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg max-w-md w-full p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Join Group</h3>
            <p className="text-gray-600 mb-6">
              Are you sure you want to join "{group.Name}"? Your wallet address will be added to the group.
            </p>
            <div className="flex gap-3">
              <button
                onClick={() => setShowJoinModal(false)}
                className="btn btn-secondary flex-1"
              >
                Cancel
              </button>
              <button
                onClick={handleJoinGroup}
                disabled={joinGroup.isPending}
                className="btn btn-primary flex-1"
              >
                {joinGroup.isPending ? 'Joining...' : 'Join Group'}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Contribute Modal */}
      {showContributeModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg max-w-md w-full p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Make Contribution</h3>
            <form onSubmit={handleContribute} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Amount (XLM)
                </label>
                <input
                  type="number"
                  step="0.0000001"
                  min="0"
                  required
                  className="input"
                  placeholder="Enter amount"
                  value={contributeAmount}
                  onChange={(e) => setContributeAmount(e.target.value)}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Your Secret Key
                </label>
                <input
                  type="password"
                  required
                  className="input"
                  placeholder="Enter your secret key"
                  value={secretKey}
                  onChange={(e) => setSecretKey(e.target.value)}
                />
                <p className="text-xs text-gray-500 mt-1">
                  Your secret key is needed to sign the transaction
                </p>
              </div>
              <div className="flex gap-3">
                <button
                  type="button"
                  onClick={() => setShowContributeModal(false)}
                  className="btn btn-secondary flex-1"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={contributeToGroup.isPending}
                  className="btn btn-primary flex-1"
                >
                  {contributeToGroup.isPending ? 'Contributing...' : 'Contribute'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}

export default GroupDetailPage