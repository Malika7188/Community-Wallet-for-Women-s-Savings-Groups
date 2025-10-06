import { useAuth } from '../contexts/AuthContext'
import { useTransactions } from '../hooks/useWallet'
import { ExternalLink, ArrowUpRight, ArrowDownRight, Clock } from 'lucide-react'
import LoadingSpinner from '../components/LoadingSpinner'

const TransactionsPage = () => {
  const { user } = useAuth()
  const { data: transactions, isLoading } = useTransactions(user?.wallet || '')

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  const getTransactionIcon = () => {
    // This is a simplified way to determine transaction type
    // In a real app, you'd analyze the transaction operations
    return Math.random() > 0.5 ? ArrowUpRight : ArrowDownRight
  }

  const getTransactionColor = () => {
    return Math.random() > 0.5 ? 'text-green-600' : 'text-red-600'
  }

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Transaction History</h1>
        <p className="text-gray-600 mt-2">
          View all your wallet transactions on the Stellar network
        </p>
      </div>

      <div className="card">
        {isLoading ? (
          <div className="flex justify-center py-12">
            <LoadingSpinner size="lg" />
          </div>
        ) : transactions?.data?.transactions?.length === 0 ? (
          <div className="text-center py-12">
            <Clock className="w-16 h-16 text-gray-400 mx-auto mb-4" />
            <h3 className="text-xl font-semibold text-gray-900 mb-2">No Transactions Yet</h3>
            <p className="text-gray-600 mb-6">
              Your transaction history will appear here once you start using your wallet.
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-lg font-semibold text-gray-900">
                Recent Transactions
              </h2>
              <span className="text-sm text-gray-500">
                {transactions?.data?.transactions?.length || 0} transactions
              </span>
            </div>

            <div className="space-y-3">
              {transactions?.data?.transactions?.map((tx) => {
                const Icon = getTransactionIcon()
                const colorClass = getTransactionColor()
                
                return (
                  <div
                    key={tx.hash}
                    className="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
                  >
                    <div className="flex items-center">
                      <div className={`w-10 h-10 bg-gray-100 rounded-lg flex items-center justify-center`}>
                        <Icon className={`w-5 h-5 ${colorClass}`} />
                      </div>
                      <div className="ml-4">
                        <div className="flex items-center space-x-2">
                          <p className="font-medium text-gray-900">
                            Transaction
                          </p>
                          <span className={`px-2 py-1 text-xs rounded-full ${
                            tx.successful 
                              ? 'bg-green-100 text-green-800' 
                              : 'bg-red-100 text-red-800'
                          }`}>
                            {tx.successful ? 'Success' : 'Failed'}
                          </span>
                        </div>
                        <div className="flex items-center space-x-4 text-sm text-gray-600">
                          <span>Ledger: {tx.ledger}</span>
                          <span>{formatDate(tx.created_at)}</span>
                          <span>Fee: {tx.fee_charged} stroops</span>
                        </div>
                      </div>
                    </div>
                    
                    <div className="flex items-center space-x-3">
                      <div className="text-right">
                        <code className="text-xs text-gray-500 block">
                          {tx.hash.slice(0, 8)}...{tx.hash.slice(-8)}
                        </code>
                        {tx.memo && (
                          <p className="text-xs text-gray-600 mt-1">
                            Memo: {tx.memo}
                          </p>
                        )}
                      </div>
                      <a
                        href={`https://stellar.expert/explorer/testnet/tx/${tx.hash}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="p-2 text-gray-400 hover:text-gray-600"
                      >
                        <ExternalLink className="w-4 h-4" />
                      </a>
                    </div>
                  </div>
                )
              })}
            </div>

            {/* Pagination could be added here */}
            <div className="text-center pt-6">
              <p className="text-sm text-gray-500">
                Showing recent transactions. Visit Stellar Explorer for complete history.
              </p>
            </div>
          </div>
        )}
      </div>

      {/* Transaction Stats */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-8">
        <div className="card text-center">
          <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center mx-auto mb-3">
            <Clock className="w-6 h-6 text-blue-600" />
          </div>
          <h3 className="font-semibold text-gray-900">Total Transactions</h3>
          <p className="text-2xl font-bold text-gray-900 mt-1">
            {transactions?.data?.transactions?.length || 0}
          </p>
        </div>

        <div className="card text-center">
          <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center mx-auto mb-3">
            <ArrowUpRight className="w-6 h-6 text-green-600" />
          </div>
          <h3 className="font-semibold text-gray-900">Successful</h3>
          <p className="text-2xl font-bold text-gray-900 mt-1">
            {transactions?.data?.transactions?.filter(tx => tx.successful).length || 0}
          </p>
        </div>

        <div className="card text-center">
          <div className="w-12 h-12 bg-red-100 rounded-lg flex items-center justify-center mx-auto mb-3">
            <ArrowDownRight className="w-6 h-6 text-red-600" />
          </div>
          <h3 className="font-semibold text-gray-900">Failed</h3>
          <p className="text-2xl font-bold text-gray-900 mt-1">
            {transactions?.data?.transactions?.filter(tx => !tx.successful).length || 0}
          </p>
        </div>
      </div>
    </div>
  )
}

export default TransactionsPage