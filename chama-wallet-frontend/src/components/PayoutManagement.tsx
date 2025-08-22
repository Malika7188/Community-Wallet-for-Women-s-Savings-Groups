import React, { useState } from 'react'
import { DollarSign, Check, X, Clock, Users } from 'lucide-react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { groupApi, payoutApi } from '../services/api'
import type { Group, User, PayoutRequest } from '../types'
import toast from 'react-hot-toast'

interface PayoutManagementProps {
  group: Group
  currentUser: User
}

const PayoutManagement: React.FC<PayoutManagementProps> = ({ group, currentUser }) => {
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [payoutData, setPayoutData] = useState({
    recipient_id: '',
    amount: 0,
    round: 1
  })
  const queryClient = useQueryClient()

  const { data: payoutRequestsResponse } = useQuery({
    queryKey: ['payout-requests', group.ID],
    queryFn: () => groupApi.getPayoutRequests(group.ID)
  })

  const payoutRequests = payoutRequestsResponse?.data || []

  const createPayoutMutation = useMutation({
    mutationFn: (data: any) => groupApi.createPayoutRequest(group.ID, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['payout-requests', group.ID] })
      toast.success('Payout request created!')
      setShowCreateModal(false)
      setPayoutData({ recipient_id: '', amount: 0, round: 1 })
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to create payout request')
    }
  })

  const approvePayoutMutation = useMutation({
    mutationFn: ({ id, approved }: { id: string, approved: boolean }) => 
      payoutApi.approvePayout(id, { approved }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['payout-requests', group.ID] })
      toast.success('Payout decision recorded!')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to process payout decision')
    }
  })

  const currentUserMember = group.Members?.find(m => m.UserID === currentUser.id)
  const isAdmin = currentUserMember && ['creator', 'admin'].includes(currentUserMember.Role)
  const approvedMembers = group.Members?.filter(m => m.Status === 'approved') || []

  const handleCreatePayout = (e: React.FormEvent) => {
    e.preventDefault()
    createPayoutMutation.mutate(payoutData)
  }

  const handleApprovePayout = (payoutId: string, approved: boolean) => {
    approvePayoutMutation.mutate({ id: payoutId, approved })
  }

  const getPayoutStatus = (payout: PayoutRequest) => {
    const approvals = payout.Approvals?.filter(a => a.Approved).length || 0
    const rejections = payout.Approvals?.filter(a => !a.Approved).length || 0

    if (payout.Status === 'approved') return { text: 'Approved', color: 'green' }
    if (payout.Status === 'rejected') return { text: 'Rejected', color: 'red' }
    if (rejections > 0) return { text: 'Rejected', color: 'red' }
    if (approvals >= 1) return { text: 'Approved', color: 'green' }
    return { text: `Pending (${approvals}/1 approval)`, color: 'yellow' }
  }

  const hasUserVoted = (payout: PayoutRequest) => {
    return payout.Approvals?.some(a => a.Admin.id === currentUser.id)
  }

  return (
    <div className="space-y-6">
      {/* Create Payout Request */}
      {isAdmin && group.Status === 'active' && (
        <div className="card">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold">Payout Management</h3>
            <button
              onClick={() => setShowCreateModal(true)}
              className="btn btn-primary"
            >
              <DollarSign className="w-4 h-4 mr-2" />
              Create Payout Request
            </button>
          </div>
          
          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <h4 className="font-medium text-blue-900 mb-2">Payout Process</h4>
            <ul className="text-sm text-blue-800 space-y-1">
              <li>• Only admins can create payout requests</li>
              <li>• At least 1 admin must approve each payout</li>
              <li>• All members are notified when payouts are approved</li>
              <li>• Payouts follow the predetermined order set during group activation</li>
            </ul>
          </div>
        </div>
      )}

      {/* Payout Requests List */}
      <div className="card">
        <h3 className="text-lg font-semibold mb-4">Payout Requests</h3>
        
        {payoutRequests.length === 0 ? (
          <div className="text-center py-8 text-gray-500">
            No payout requests yet
          </div>
        ) : (
          <div className="space-y-4">
            {payoutRequests.map((payout: PayoutRequest) => {
              const status = getPayoutStatus(payout)
              const userHasVoted = hasUserVoted(payout)
              
              return (
                <div key={payout.ID} className="border border-gray-200 rounded-lg p-4">
                  <div className="flex items-center justify-between mb-3">
                    <div>
                      <h4 className="font-medium">
                        Payout to {payout.Recipient.name}
                      </h4>
                      <p className="text-sm text-gray-600">
                        Amount: {payout.Amount} XLM • Round {payout.Round}
                      </p>
                    </div>
                    <span className={`px-3 py-1 rounded-full text-sm font-medium ${
                      status.color === 'green' ? 'bg-green-100 text-green-800' :
                      status.color === 'red' ? 'bg-red-100 text-red-800' :
                      'bg-yellow-100 text-yellow-800'
                    }`}>
                      {status.text}
                    </span>
                  </div>

                  {/* Approval Status */}
                  <div className="mb-3">
                    <h5 className="text-sm font-medium text-gray-700 mb-2">Approvals:</h5>
                    <div className="flex flex-wrap gap-2">
                      {payout.Approvals?.map((approval) => (
                        <div key={approval.ID} className="flex items-center space-x-1">
                          <span className="text-sm">{approval.Admin.name}</span>
                          {approval.Approved ? (
                            <Check className="w-4 h-4 text-green-500" />
                          ) : (
                            <X className="w-4 h-4 text-red-500" />
                          )}
                        </div>
                      ))}
                    </div>
                  </div>

                  {/* Admin Actions */}
                  {isAdmin && payout.Status === 'pending' && !userHasVoted && (
                    <div className="flex space-x-2">
                      <button
                        onClick={() => handleApprovePayout(payout.ID, true)}
                        className="btn btn-primary btn-sm"
                        disabled={approvePayoutMutation.isPending}
                      >
                        <Check className="w-4 h-4 mr-1" />
                        Approve
                      </button>
                      <button
                        onClick={() => handleApprovePayout(payout.ID, false)}
                        className="btn btn-secondary btn-sm"
                        disabled={approvePayoutMutation.isPending}
                      >
                        <X className="w-4 h-4 mr-1" />
                        Reject
                      </button>
                    </div>
                  )}

                  {userHasVoted && (
                    <p className="text-sm text-gray-600">
                      You have already voted on this request
                    </p>
                  )}
                </div>
              )
            })}
          </div>
        )}
      </div>

      {/* Create Payout Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold mb-4">Create Payout Request</h3>
            <form onSubmit={handleCreatePayout}>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Recipient
                  </label>
                  <select
                    value={payoutData.recipient_id}
                    onChange={(e) => setPayoutData({
                      ...payoutData,
                      recipient_id: e.target.value
                    })}
                    className="input"
                    required
                  >
                    <option value="">Select recipient...</option>
                    {approvedMembers.map((member) => (
                      <option key={member.ID} value={member.UserID}>
                        {member.User.name}
                      </option>
                    ))}
                  </select>
                </div>
                
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Amount (XLM)
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    value={payoutData.amount}
                    onChange={(e) => setPayoutData({
                      ...payoutData,
                      amount: parseFloat(e.target.value)
                    })}
                    className="input"
                    required
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Round
                  </label>
                  <input
                    type="number"
                    value={payoutData.round}
                    onChange={(e) => setPayoutData({
                      ...payoutData,
                      round: parseInt(e.target.value)
                    })}
                    className="input"
                    required
                  />
                </div>
              </div>

              <div className="flex space-x-3 mt-6">
                <button
                  type="button"
                  onClick={() => setShowCreateModal(false)}
                  className="btn btn-secondary flex-1"
                >
                  Cancel
                </button>
                <button 
                  type="submit" 
                  className="btn btn-primary flex-1"
                  disabled={createPayoutMutation.isPending}
                >
                  {createPayoutMutation.isPending ? 'Creating...' : 'Create Request'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}

export default PayoutManagement
