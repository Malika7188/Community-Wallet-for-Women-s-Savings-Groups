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
    <div className="max-w-5xl mx-auto px-2 sm:px-6 lg:px-8 py-10">
      <div className="mb-10">
        <h1 className="text-4xl font-black text-[#1a237e] tracking-tight mb-1" style={{ fontFamily: 'Inter, Roboto, sans-serif' }}>My Wallet</h1>
        <p className="text-gray-600 text-lg font-medium">Manage your Stellar wallet and transactions</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-10">
        {/* Wallet Overview */}
        <div className="lg:col-span-2 space-y-8">
          {/* Balance Card */}
          <div className="relative rounded-3xl p-10 bg-white/80 backdrop-blur-lg shadow-2xl border border-gray-100 flex flex-col items-center overflow-hidden min-h-[260px]">
            <div className="absolute inset-0 pointer-events-none opacity-40" style={{background: 'linear-gradient(135deg, #e0f7fa 0%, #e8f5e9 100%)'}} />
            <div className="relative flex w-full items-center justify-between mb-6 z-10">
              <h2 className="text-2xl font-extrabold text-[#1a237e]">Wallet Balance</h2>
              <button
                onClick={() => refetch()}
                className="p-2 text-[#2ecc71] hover:text-[#1a237e] bg-white/70 rounded-full shadow"
                disabled={balanceLoading}
              >
                <RefreshCw className={`w-6 h-6 ${balanceLoading ? 'animate-spin' : ''}`} />
              </button>
            </div>
            <div className="relative text-center py-6 z-10">
              <div className="w-24 h-24 bg-gradient-to-br from-[#1a237e] to-[#2ecc71] rounded-full flex items-center justify-center mx-auto mb-4 shadow-lg">
                <Wallet className="w-12 h-12 text-white" />
              </div>
              {balanceLoading ? (
                <LoadingSpinner />
              ) : (
                <>
                  <p className="text-5xl font-black text-[#1a237e] mb-2" style={{ fontFamily: 'Inter, Roboto, sans-serif' }}>
                    {walletBalance}
                  </p>
                  <p className="text-gray-600 text-lg font-medium">Available Balance</p>
                </>
              )}
            </div>
          </div>

          {/* Wallet Address */}
          <div className="rounded-2xl bg-white/90 border border-gray-100 shadow p-6 flex flex-col gap-2">
            <h3 className="text-lg font-bold text-[#1a237e] mb-2">Wallet Address</h3>
            <div className="flex items-center space-x-3 p-3 bg-gray-50 rounded-xl">
              <code className="flex-1 text-base font-mono break-all text-[#1a237e]">
                {user?.wallet}
              </code>
              <button
                onClick={() => copyToClipboard(user?.wallet || '')}
                className="p-2 text-[#2ecc71] hover:text-[#1a237e]"
              >
                <Copy className="w-5 h-5" />
              </button>
              <a
                href={`https://stellar.expert/explorer/testnet/account/${user?.wallet}`}
                target="_blank"
                rel="noopener noreferrer"
                className="p-2 text-[#2ecc71] hover:text-[#1a237e]"
              >
                <ExternalLink className="w-5 h-5" />
              </a>
            </div>
          </div>

          {/* User Secret Key */}
          <div className="rounded-2xl bg-yellow-50 border border-yellow-200 shadow p-6">
            <h3 className="text-lg font-bold text-[#b58900] mb-2">Your Secret Key</h3>
            <div className="flex items-center space-x-2 mb-3">
              <code className="flex-1 text-xs bg-gray-100 p-2 rounded font-mono break-all text-[#b58900]">
                {user?.secret_key || 'Not available'}
              </code>
              <button
                onClick={() => copyToClipboard(user?.secret_key || '')}
                className="p-2 text-[#b58900] hover:text-yellow-700"
              >
                <Copy className="w-4 h-4" />
              </button>
            </div>
            <p className="text-sm text-yellow-700">
              ⚠️ Keep this secret key safe. You need it to make transactions.
            </p>
          </div>

          {/* Quick Actions */}
          <div className="rounded-2xl bg-white/90 border border-gray-100 shadow p-6">
            <h3 className="text-lg font-bold text-[#1a237e] mb-2">Quick Actions</h3>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <button
                onClick={() => setShowTransferModal(true)}
                className="inline-flex items-center justify-center px-6 py-4 rounded-xl bg-[#2ecc71] text-white font-bold text-lg shadow hover:bg-[#27ae60] transition-colors duration-200"
              >
                <Send className="w-5 h-5 mr-2" />
                Send XLM
              </button>
              <button
                onClick={handleFundAccount}
                disabled={fundAccount.isPending}
                className="inline-flex items-center justify-center px-6 py-4 rounded-xl border border-[#2ecc71] text-[#1a237e] font-bold text-lg shadow hover:bg-[#2ecc71]/10 transition-colors duration-200"
              >
                <Plus className="w-5 h-5 mr-2" />
                {fundAccount.isPending ? 'Funding...' : 'Fund Account'}
              </button>
              <button
                onClick={handleGenerateKeypair}
                disabled={generateKeypair.isPending}
                className="inline-flex items-center justify-center px-6 py-4 rounded-xl border border-[#2ecc71] text-[#1a237e] font-bold text-md shadow hover:bg-[#2ecc71]/10 transition-colors duration-200"
              >
                <Wallet className="w-8 h-8 mr-2 text-bold" />
                {generateKeypair.isPending ? 'Generating...' : ' Generate New Secret Key'}
              </button>
            </div>
          </div>
        </div>

        {/* Sidebar */}
        <div className="space-y-8">
          {/* Wallet Info */}
          <div className="rounded-2xl bg-white/90 border border-gray-100 shadow p-6">
            <h3 className="text-lg font-bold text-[#1a237e] mb-2">Wallet Info</h3>
            <div className="space-y-3 text-base">
              <div className="flex justify-between">
                <span className="text-gray-600">Network:</span>
                <span className="font-bold">Testnet</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Asset:</span>
                <span className="font-bold">XLM (Stellar Lumens)</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Status:</span>
                <span className="font-bold text-green-600">Active</span>
              </div>
            </div>
          </div>

          {/* Security Tips */}
          <div className="rounded-2xl bg-white/90 border border-gray-100 shadow p-6">
            <h3 className="text-lg font-bold text-[#1a237e] mb-2">Security Tips</h3>
            <ul className="space-y-2 text-base text-gray-600">
              <li>• Never share your secret key with anyone</li>
              <li>• Always verify recipient addresses</li>
              <li>• Keep your secret key backed up safely</li>
              <li>• Use testnet for learning and testing</li>
            </ul>
          </div>

          {/* Support */}
          <div className="rounded-2xl bg-white/90 border border-gray-100 shadow p-6">
            <h3 className="text-lg font-bold text-[#1a237e] mb-2">Need Help?</h3>
            <p className="text-base text-gray-600 mb-4">
              Having trouble with your wallet? Check out our resources.
            </p>
            <div className="space-y-2">
              <a href="#" className="block text-base text-[#2ecc71] hover:text-[#1a237e]">
                Wallet Guide
              </a>
              <a href="#" className="block text-base text-[#2ecc71] hover:text-[#1a237e]">
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
