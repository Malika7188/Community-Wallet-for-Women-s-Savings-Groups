import React from 'react';
import { ArrowUpRightIcon, ArrowDownRightIcon } from '@heroicons/react/24/outline';

interface Transaction {
  type: 'contribution' | 'withdrawal';
  description: string;
  amount: string;
  time: string;
  group: string;
}

interface TransactionHistoryProps {
  transactions: Transaction[];
}

const TransactionHistory: React.FC<TransactionHistoryProps> = ({ transactions }) => {
  return (
    <div className="bg-white rounded-2xl shadow-md p-6">
      <h2 className="text-xl font-semibold text-gray-900 mb-4">Recent Activity</h2>
      {transactions.length === 0 ? (
        <div className="text-center py-8 text-gray-500">No recent activity</div>
      ) : (
        <ul className="divide-y divide-gray-100">
          {transactions.map((tx, idx) => (
            <li key={idx} className="flex items-center justify-between py-3">
              <div className="flex items-center gap-3">
                <div className={`w-10 h-10 rounded-lg flex items-center justify-center ${tx.type === 'contribution' ? 'bg-emerald-100' : 'bg-red-100'}`}>
                  {tx.type === 'contribution' ? (
                    <ArrowUpRightIcon className="w-5 h-5 text-emerald-600" />
                  ) : (
                    <ArrowDownRightIcon className="w-5 h-5 text-red-600" />
                  )}
                </div>
                <div>
                  <div className="font-medium text-gray-900">{tx.description}</div>
                  <div className="text-xs text-gray-500">{tx.time} • {tx.group}</div>
                </div>
              </div>
              <span className={`font-semibold ${tx.type === 'contribution' ? 'text-emerald-600' : 'text-red-600'}`}>{tx.amount}</span>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default TransactionHistory;
