import React, { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { groupApi } from '../services/api'
import type { Group, User } from '../types'
import { toast } from 'react-hot-toast'
import { CheckCircle, Clock, DollarSign, Users } from 'lucide-react'

interface RoundContributionsProps {
  group: Group
  currentUser: User
}

const RoundContributions: React.FC<RoundContributionsProps> = ({ group, currentUser }) => {
  const [selectedRound, setSelectedRound] = useState(group.CurrentRound || 1)
  const [showContributeModal, setShowContributeModal] = useState(false)
  const [secretKey, setSecretKey] = useState('')
  
  const queryClient = useQueryClient()

  // Get current user's member info
  const currentMember = group.Members?.find(m => m.UserID === currentUser.id)
  const isAdmin = currentMember?.Role === 'admin' || currentMember?.Role === 'creator'

  // Fetch round status
  const { data: roundStatus, isLoading } = useQuery({
    queryKey: ['round-status', group.ID, selectedRound],
    queryFn: () => groupApi.getRoundStatus(group.ID, selectedRound),
    enabled: !!group.ID && group.Status === 'active'
  })

  // Contribute mutation
  const contributeMutation = useMutation({
    mutationFn: (data: { round: number, amount: number, secret: string }) =>
      groupApi.contributeToRound(group.ID, data),
    onSuccess: () => {
      toast.success('Contribution successful!')
      setShowContributeModal(false)
      setSecretKey('')
      queryClient.invalidateQueries({ queryKey: ['round-status', group.ID, selectedRound] })
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Contribution failed')
    }
  })

  // Authorize payout mutation
  const authorizeMutation = useMutation({
    mutationFn: (round: number) => groupApi.authorizeRoundPayout(group.ID, { round }),
    onSuccess: () => {
      toast.success('Payout authorized successfully!')
      queryClient.invalidateQueries({ queryKey: ['round-status', group.ID, selectedRound] })
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Authorization failed')
    }
  })

  const handleContribute = (e: React.FormEvent) => {
    e.preventDefault()
    contributeMutation.mutate({
      round: selectedRound,
      amount: group.ContributionAmount || 0,
      secret: secretKey
    })
  }

  const handleAuthorize = () => {
    authorizeMutation.mutate(selectedRound)
  }

  if (group.Status !== 'active') {
    return (
      <div className="bg-gray-50 rounded-lg p-6 text-center">
        <p className="text-gray-600">Group must be active to view contributions</p>
      </div>
    )
  }

  if (isLoading) {
    return <div className="text-center py-4">Loading round status...</div>
  }

  const roundData = roundStatus?.data
  const currentUserStatus = roundData?.member_status?.find(ms => ms.member.UserID === currentUser.id)
  const allPaid = roundData?.paid_members === roundData?.total_members
  const canAuthorize = isAdmin && allPaid && !roundData?.round_status?.PayoutAuthorized

  return (
    <div className="space-y-6">
      {/* Round Selector */}
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold">Round Contributions</h3>
        <div className="flex items-center space-x-2">
          <label className="text-sm font-medium">Round:</label>
          <select
            value={selectedRound}
            onChange={(e) => setSelectedRound(parseInt(e.target.value))}
            className="px-3 py-1 border border-gray-300 rounded-md text-sm"
          >
            {Array.from({ length: group.Members?.filter(m => m.Status === 'approved').length || 1 }, (_, i) => (
              <option key={i + 1} value={i + 1}>Round {i + 1}</option>
            ))}
          </select>
        </div>
      </div>

      {/* Round Summary */}
      {roundData && (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div className="bg-blue-50 p-4 rounded-lg">
            <div className="flex items-center space-x-2">
              <DollarSign className="w-5 h-5 text-blue-600" />
              <div>
                <p className="text-sm text-blue-600">Total Collected</p>
                <p className="font-semibold">{roundData.round_status?.TotalReceived || 0} XLM</p>
              </div>
            </div>
          </div>
          
          <div className="bg-green-50 p-4 rounded-lg">
            <div className="flex items-center space-x-2">
              <Users className="w-5 h-5 text-green-600" />
              <div>
                <p className="text-sm text-green-600">Contributors</p>
                <p className="font-semibold">{roundData.paid_members}/{roundData.total_members}</p>
              </div>
            </div>
          </div>
          
          <div className="bg-yellow-50 p-4 rounded-lg">
            <div className="flex items-center space-x-2">
              <Clock className="w-5 h-5 text-yellow-600" />
              <div>
                <p className="text-sm text-yellow-600">Status</p>
                <p className="font-semibold capitalize">{roundData.round_status?.Status?.replace('_', ' ')}</p>
              </div>
            </div>
          </div>
          
          <div className="bg-purple-50 p-4 rounded-lg">
            <div className="flex items-center space-x-2">
              <CheckCircle className="w-5 h-5 text-purple-600" />
              <div>
                <p className="text-sm text-purple-600">Payout Status</p>
                <p className="font-semibold">
                  {roundData.round_status?.PayoutAuthorized ? 'Authorized' : 'Pending'}
                </p>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* User Action Buttons */}
      <div className="flex space-x-3">
        {currentUserStatus && !currentUserStatus.has_paid && (
          <button
            onClick={() => setShowContributeModal(true)}
            className="btn btn-primary"
          >
            Contribute {group.ContributionAmount} XLM
          </button>
        )}
        
        {canAuthorize && (
          <button
            onClick={handleAuthorize}
            disabled={authorizeMutation.isPending}
            className="btn btn-success"
          >
            {authorizeMutation.isPending ? 'Authorizing...' : 'Authorize Payout'}
          </button>
        )}
      </div>

      {/* Member Status List */}
      {roundData && (
        <div className="bg-white border rounded-lg overflow-hidden">
          <div className="px-4 py-3 bg-gray-50 border-b">
            <h4 className="font-medium">Member Contributions - Round {selectedRound}</h4>
          </div>
          <div className="divide-y">
            {roundData.member_status.map((memberStatus) => (
              <div key={memberStatus.member.ID} className="px-4 py-3 flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <div className={`w-3 h-3 rounded-full ${
                    memberStatus.has_paid ? 'bg-green-500' : 'bg-gray-300'
                  }`} />
                  <div>
                    <p className="font-medium">{memberStatus.member.User.name}</p>
                    <p className="text-sm text-gray-500">{memberStatus.member.User.email}</p>
                  </div>
                </div>
                <div className="text-right">
                  {memberStatus.has_paid ? (
                    <div>
                      <p className="text-green-600 font-medium">✓ Paid</p>
                      <p className="text-xs text-gray-500">
                        {memberStatus.contribution?.Amount} XLM
                      </p>
                    </div>
                  ) : (
                    <p className="text-gray-500">Pending</p>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Contribute Modal */}
      {showContributeModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold mb-4">Contribute to Round {selectedRound}</h3>
            <form onSubmit={handleContribute}>
              <div className="space-y-4">
                <div>
                  <p className="text-sm text-gray-600 mb-2">
                    Amount: <span className="font-medium">{group.ContributionAmount} XLM</span>
                  </p>
                </div>
                
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Secret Key
                  </label>
                  <input
                    type="password"
                    value={secretKey}
                    onChange={(e) => setSecretKey(e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="Enter your secret key"
                    required
                  />
                </div>
              </div>
              
              <div className="flex space-x-3 mt-6">
                <button
                  type="button"
                  onClick={() => {
                    setShowContributeModal(false)
                    setSecretKey('')
                  }}
                  className="btn btn-secondary flex-1"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={contributeMutation.isPending}
                  className="btn btn-primary flex-1"
                >
                  {contributeMutation.isPending ? 'Contributing...' : 'Contribute'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}

export default RoundContributions
