import { useState } from 'react'
import { useAuth } from '../contexts/AuthContext'
import { useBalance, useWallet } from '../hooks/useWallet'
import { 
  Wallet, 
  Copy, 
  ExternalLink, 
  Plus, 
  Send, 
  RefreshCw,
  Eye,
  EyeOff
} from 'lucide-react'
import LoadingSpinner from '../components/LoadingSpinner'
import toast from 'react-hot-toast'

const WalletPage = () => {
  const { user } = useAuth()
  const { data: balance, isLoading: balanceLoading, refetch } = useBalance(user?.wallet || '')
  const { generateKeypair, fundAccount, transferFunds } = useWallet()
  
  const [showTransferModal, setShowTransferModal] = useState(false)
  const [showKeypairModal, setShowKeypairModal] = useState(false)
  const [showSecretKey, setShowSecretKey] = useState(false)
  const [generatedKeypair, setGeneratedKeypair] = useState<any>(null)
  
  const [transferData, setTransferData] = useState({
    to_address: '',
    amount: '',
    from_seed: ''
  })

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    toast.success('Copied to clipboard!')
  }

  const handleGenerateKeypair = async () => {
    try {
      const result = await generateKeypair.mutateAsync()
      setGeneratedKeypair(result.data)
      setShowKeypairModal(true)
    } catch (error) {
      // Error handled by mutation
    }
  }

  const handleFundAccount = async () => {
    if (!user?.wallet) return
    try {
      await fundAccount.mutateAsync(user.wallet)
      refetch()
    } catch (error) {
      // Error handled by mutation
    }
  }

  const handleTransfer = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await transferFunds.mutateAsync(transferData)
      setShowTransferModal(false)
      setTransferData({ to_address: '', amount: '', from_seed: '' })
      refetch()
    } catch (error) {
      // Error handled by mutation
    }
  }

  const walletBalance = balance?.data?.balances?.[0]?.split(': ')[1] || '0 XLM'

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">My Wallet</h1>
        <p className="text-gray-600 mt-2">
          Manage your Stellar wallet and transactions
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Wallet Overview */}
        <div className="lg:col-span-2 space-y-6">
          {/* Balance Card */}
          <div className="card">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-semibold text-gray-900">Wallet Balance</h2>
              <button
                onClick={() => refetch()}
                className="p-2 text-gray-400 hover:text-gray-600"
                disabled={balanceLoading}
              >
                <RefreshCw className={`w-5 h-5 ${balanceLoading ? 'animate-spin' : ''}`} />
              </button>
            </div>
            
            <div className="text-center py-8">
              <div className="w-20 h-20 bg-gradient-to-r from-stellar-500 to-primary-600 rounded-full flex items-center justify-center mx-auto mb-4">
                <Wallet className="w-10 h-10 text-white" />
              </div>
              
              {balanceLoading ? (
                <LoadingSpinner />
              ) : (
                <>
                  <p className="text-4xl font-bold text-gray-900 mb-2">
                    {walletBalance}
                  </p>
                  <p className="text-gray-600">Available Balance</p>
                </>
              )}
            </div>
          </div>

          {/* Wallet Address */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Wallet Address</h3>
            <div className="flex items-center space-x-3 p-4 bg-gray-50 rounded-lg">
              <code className="flex-1 text-sm font-mono break-all">
                {user?.wallet}
              </code>
              <button
                onClick={() => copyToClipboard(user?.wallet || '')}
                className="p-2 text-gray-400 hover:text-gray-600"
              >
                <Copy className="w-5 h-5" />
              </button>
              <a
                href={`https://stellar.expert/explorer/testnet/account/${user?.wallet}`}
                target="_blank"
                rel="noopener noreferrer"
                className="p-2 text-gray-400 hover:text-gray-600"
              >
                <ExternalLink className="w-5 h-5" />
              </a>
            </div>
          </div>

          {/* Quick Actions */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Quick Actions</h3>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <button
                onClick={() => setShowTransferModal(true)}
                className="btn btn-primary"
              >
                <Send className="w-4 h-4 mr-2" />
                Send XLM
              </button>
              <button
                onClick={handleFundAccount}
                disabled={fundAccount.isPending}
                className="btn btn-outline"
              >
                <Plus className="w-4 h-4 mr-2" />
                {fundAccount.isPending ? 'Funding...' : 'Fund Account'}
              </button>
              <button
                onClick={handleGenerateKeypair}
                disabled={generateKeypair.isPending}
                className="btn btn-outline"
              >
                <Wallet className="w-4 h-4 mr-2" />
                {generateKeypair.isPending ? 'Generating...' : 'New Keypair'}
              </button>
            </div>
          </div>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Wallet Info */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Wallet Info</h3>
            <div className="space-y-3 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-600">Network:</span>
                <span className="font-medium">Testnet</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Asset:</span>
                <span className="font-medium">XLM (Stellar Lumens)</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Status:</span>
                <span className="font-medium text-green-600">Active</span>
              </div>
            </div>
          </div>

          {/* Security Tips */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Security Tips</h3>
            <ul className="space-y-2 text-sm text-gray-600">
              <li>• Never share your secret key with anyone</li>
              <li>• Always verify recipient addresses</li>
              <li>• Keep your secret key backed up safely</li>
              <li>• Use testnet for learning and testing</li>
            </ul>
          </div>

          {/* Support */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Need Help?</h3>
            <p className="text-sm text-gray-600 mb-4">
              Having trouble with your wallet? Check out our resources.
            </p>
            <div className="space-y-2">
              <a href="#" className="block text-sm text-primary-600 hover:text-primary-700">
                Wallet Guide
              </a>
              <a href="#" className="block text-sm text-primary-600 hover:text-primary-700">
                Contact Support
              </a>
            </div>
          </div>
        </div>
      </div>

      {/* Transfer Modal */}
      {showTransferModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg max-w-md w-full p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Send XLM</h3>
            <form onSubmit={handleTransfer} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Recipient Address
                </label>
                <input
                  type="text"
                  required
                  className="input"
                  placeholder="Enter recipient's wallet address"
                  value={transferData.to_address}
                  onChange={(e) => setTransferData(prev => ({ ...prev, to_address: e.target.value }))}
                />
              </div>
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
                  value={transferData.amount}
                  onChange={(e) => setTransferData(prev => ({ ...prev, amount: e.target.value }))}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Your Secret Key
                </label>
                <div className="relative">
                  <input
                    type={showSecretKey ? 'text' : 'password'}
                    required
                    className="input pr-10"
                    placeholder="Enter your secret key"
                    value={transferData.from_seed}
                    onChange={(e) => setTransferData(prev => ({ ...prev, from_seed: e.target.value }))}
                  />
                  <button
                    type="button"
                    onClick={() => setShowSecretKey(!showSecretKey)}
                    className="absolute inset-y-0 right-0 pr-3 flex items-center"
                  >
                    {showSecretKey ? (
                      <EyeOff className="h-5 w-5 text-gray-400" />
                    ) : (
                      <Eye className="h-5 w-5 text-gray-400" />
                    )}
                  </button>
                </div>
              </div>
              <div className="flex gap-3">
                <button
                  type="button"
                  onClick={() => setShowTransferModal(false)}
                  className="btn btn-secondary flex-1"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={transferFunds.isPending}
                  className="btn btn-primary flex-1"
                >
                  {transferFunds.isPending ? 'Sending...' : 'Send'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Keypair Modal */}
      {showKeypairModal && generatedKeypair && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg max-w-md w-full p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">New Keypair Generated</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Public Key
                </label>
                <div className="flex items-center space-x-2">
                  <code className="flex-1 text-xs bg-gray-100 p-2 rounded font-mono break-all">
                    {generatedKeypair.public_key}
                  </code>
                  <button
                    onClick={() => copyToClipboard(generatedKeypair.public_key)}
                    className="p-2 text-gray-400 hover:text-gray-600"
                  >
                    <Copy className="w-4 h-4" />
                  </button>
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Secret Key
                </label>
                <div className="flex items-center space-x-2">
                  <code className="flex-1 text-xs bg-gray-100 p-2 rounded font-mono break-all">
                    {generatedKeypair.secret_seed}
                  </code>
                  <button
                    onClick={() => copyToClipboard(generatedKeypair.secret_seed)}
                    className="p-2 text-gray-400 hover:text-gray-600"
                  >
                    <Copy className="w-4 h-4" />
                  </button>
                </div>
              </div>
              <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-3">
                <p className="text-sm text-yellow-800">
                  <strong>Important:</strong> Save your secret key securely. You'll need it to access this wallet.
                </p>
              </div>
            </div>
            <button
              onClick={() => setShowKeypairModal(false)}
              className="btn btn-primary w-full mt-4"
            >
              Close
            </button>
          </div>
        </div>
      )}
    </div>
  )
}

export default WalletPage