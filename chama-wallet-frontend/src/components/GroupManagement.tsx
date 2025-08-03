import React, { useState, useEffect } from 'react'
import { Settings, UserPlus, CheckCircle, Copy } from 'lucide-react'
import { useMutation, useQueryClient, useQuery } from '@tanstack/react-query'
import { groupApi } from '../services/api'
import toast from 'react-hot-toast'
import type { Group, User } from '../types'
import RoundContributions from './RoundContributions'

interface GroupManagementProps {
  group: Group
  currentUser: User
}

const GroupManagement: React.FC<GroupManagementProps> = ({
  group,
  currentUser,
}) => {
  const [showInviteModal, setShowInviteModal] = useState(false)
  const [showApproveModal, setShowApproveModal] = useState(false)
  const [showActivateModal, setShowActivateModal] = useState(false)
  const [inviteEmail, setInviteEmail] = useState('')
  const [availableUsers, setAvailableUsers] = useState<User[]>([])
  const [groupSettings, setGroupSettings] = useState({
    contribution_amount: 0,
    contribution_period: 30,
    payout_order: [] as string[]
  })
  const queryClient = useQueryClient()

  const inviteUserMutation = useMutation({
    mutationFn: (email: string) => groupApi.inviteToGroup(group.ID, { email }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['groups'] })
      toast.success('Invitation sent successfully!')
      setShowInviteModal(false)
      setInviteEmail('')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to send invitation')
    }
  })

  const approveGroupMutation = useMutation({
    mutationFn: () => groupApi.approveGroup(group.ID),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['groups'] })
      toast.success('Group approved successfully!')
      setShowApproveModal(false)
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to approve group')
    }
  })

  const activateGroupMutation = useMutation({
    mutationFn: (settings: any) => groupApi.activateGroup(group.ID, settings),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['groups'] })
      toast.success('Group activated successfully!')
      setShowActivateModal(false)
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to activate group')
    }
  })

  // Fetch available users when invite modal opens
  useEffect(() => {
    if (showInviteModal) {
      groupApi.getNonGroupMembers(group.ID)
        .then(response => setAvailableUsers(response.data))
        .catch(error => console.error('Failed to fetch available users:', error))
    }
  }, [showInviteModal, group.ID])

  const isAdmin = group.Members?.find(m =>
    m.UserID === currentUser.id && ['creator', 'admin'].includes(m.Role)
  )

  const isCreator = group.Members?.find(m =>
    m.UserID === currentUser.id && m.Role === 'creator'
  )

  const approvedMembers = group.Members?.filter(m => m.Status === 'approved') || []
  const pendingMembers = group.Members?.filter(m => m.Status === 'pending') || []
  const hasMinimumMembers = approvedMembers.length >= (group.MinMembers || 3)
  const isGroupFull = approvedMembers.length >= (group.MaxMembers || 20)

  const handleInvite = (e: React.FormEvent) => {
    e.preventDefault()
    inviteUserMutation.mutate(inviteEmail)
  }

  const handleApproveGroup = () => {
    approveGroupMutation.mutate()
  }

  const handleActivate = (e: React.FormEvent) => {
    e.preventDefault()
    console.log('Activating group with settings:', groupSettings)

    // Ensure payout order is not empty
    if (groupSettings.payout_order.length === 0) {
      toast.error('Payout order cannot be empty')
      return
    }

    activateGroupMutation.mutate({
      contribution_amount: groupSettings.contribution_amount,
      contribution_period: groupSettings.contribution_period,
      payout_order: groupSettings.payout_order,
    })
  }

  // Initialize payout order with member IDs in default order
  useEffect(() => {
    if (showActivateModal && approvedMembers.length > 0) {
      const defaultPayoutOrder = approvedMembers.map(member => member.UserID)
      console.log('Initializing payout order:', defaultPayoutOrder)
      console.log('Approved members:', approvedMembers.map(m => ({ id: m.UserID, name: m.User.name })))
      
      setGroupSettings(prev => ({
        ...prev,
        contribution_amount: prev.contribution_amount || 0,
        contribution_period: prev.contribution_period || 30,
        payout_order: defaultPayoutOrder
      }))
    }
  }, [showActivateModal, approvedMembers])

  useEffect(() => {
    if (group.PayoutOrder) {
      console.log('Raw PayoutOrder:', group.PayoutOrder)
      try {
        const parsed = JSON.parse(group.PayoutOrder)
        console.log('Parsed PayoutOrder:', parsed)
        console.log('Approved Members:', approvedMembers.map(m => ({ id: m.UserID, name: m.User.name })))
      } catch (error) {
        console.error('Error parsing PayoutOrder:', error)
      }
    }
  }, [group.PayoutOrder, approvedMembers])

  // Add this query to fetch the payout schedule
  const { data: payoutScheduleResponse } = useQuery({
    queryKey: ['payout-schedule', group.ID],
    queryFn: () => groupApi.getPayoutSchedule(group.ID),
    enabled: group.Status === 'active'
  })

  const payoutSchedule = payoutScheduleResponse?.data || []

  return (
    <div className="space-y-6">
      {/* Group Status */}
      <div className="card">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold">Group Status</h3>
          <span className={`px-3 py-1 rounded-full text-sm font-medium ${group.Status === 'active' ? 'bg-green-100 text-green-800' :
              group.Status === 'pending' ? 'bg-yellow-100 text-yellow-800' :
                'bg-gray-100 text-gray-800'
            }`}>
            {group.Status.charAt(0).toUpperCase() + group.Status.slice(1)}
          </span>
        </div>

        {/* Group Ready - Has minimum members for approval */}
        {group.Status === 'pending' && !group.IsApproved && isCreator && hasMinimumMembers && (
          <div className="bg-green-50 border border-green-200 rounded-lg p-4 mb-4">
            <p className="text-green-800 mb-3">
              ✅ Your group has {approvedMembers.length} members and meets the minimum requirement! You can now approve it for activation.
            </p>
            <button
              onClick={() => setShowApproveModal(true)}
              className="btn btn-primary"
            >
              <CheckCircle className="w-4 h-4 mr-2" />
              Approve Group
            </button>
          </div>
        )}

        {/* Group Approved - Ready for Activation */}
        {group.Status === 'pending' && group.IsApproved && isAdmin && (
          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <p className="text-blue-800 mb-3">
              Your group is approved! Set contribution rules and activate the group.
            </p>
            <button
              onClick={() => setShowActivateModal(true)}
              className="btn btn-primary"
            >
              <Settings className="w-4 h-4 mr-2" />
              Activate Group
            </button>
          </div>
        )}

        {/* Group Status Info */}
        <div className="bg-gray-50 rounded-lg p-4">
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <span className="text-gray-600">Members:</span>
              <span className="ml-2 font-medium">{approvedMembers.length}/{group.MaxMembers || 20}</span>
              <span className="text-xs text-gray-500 block">Min: {group.MinMembers || 3}</span>
            </div>
            <div>
              <span className="text-gray-600">Approved:</span>
              <span className="ml-2 font-medium">{group.IsApproved ? 'Yes' : 'No'}</span>
            </div>
          </div>
        </div>
      </div>

      {/* Group Wallet Info */}
      <div className="card">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Group Wallet</h3>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Group Wallet Address
            </label>
            <div className="flex items-center space-x-2">
              <code className="flex-1 text-xs bg-gray-100 p-2 rounded font-mono break-all">
                {group.Wallet}
              </code>
              <button
                onClick={() => navigator.clipboard.writeText(group.Wallet)}
                className="p-2 text-gray-400 hover:text-gray-600"
              >
                <Copy className="w-4 h-4" />
              </button>
            </div>
          </div>
          
          {isAdmin && (
            <div className="bg-red-50 border border-red-200 rounded-lg p-4">
              <label className="block text-sm font-medium text-red-700 mb-1">
                Group Secret Key (Admin Only)
              </label>
              <div className="flex items-center space-x-2 mb-2">
                <code className="flex-1 text-xs bg-gray-100 p-2 rounded font-mono break-all">
                  {group.SecretKey || 'Not available'}
                </code>
                <button
                  onClick={() => navigator.clipboard.writeText(group.SecretKey || '')}
                  className="p-2 text-gray-400 hover:text-gray-600"
                >
                  <Copy className="w-4 h-4" />
                </button>
              </div>
              <p className="text-sm text-red-700">
                ⚠️ This key controls the group wallet. Keep it secure!
              </p>
            </div>
          )}
        </div>
      </div>

      {/* Members Section */}
      <div className="card">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold">
            Members ({approvedMembers.length})
          </h3>
          {isAdmin && !isGroupFull && (
            <button
              onClick={() => setShowInviteModal(true)}
              className="btn btn-outline"
            >
              <UserPlus className="w-4 h-4 mr-2" />
              Invite Member
            </button>
          )}
        </div>

        {/* Approved Members */}
        <div className="space-y-3">
          {approvedMembers.map((member) => (
            <div key={member.ID} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
              <div>
                <p className="font-medium">{member.User.name}</p>
                <p className="text-sm text-gray-600">{member.User.email}</p>
              </div>
              <div className="flex items-center space-x-2">
                <span className={`px-2 py-1 rounded text-xs font-medium ${member.Role === 'creator' ? 'bg-purple-100 text-purple-800' :
                    member.Role === 'admin' ? 'bg-blue-100 text-blue-800' :
                      'bg-gray-100 text-gray-800'
                  }`}>
                  {member.Role}
                </span>
                <CheckCircle className="w-4 h-4 text-green-500" />
              </div>
            </div>
          ))}
        </div>

        {/* Pending Members */}
        {pendingMembers.length > 0 && (
          <div className="mt-6">
            <h4 className="font-medium text-gray-700 mb-3">Pending Approval</h4>
            <div className="space-y-2">
              {pendingMembers.map((member) => (
                <div key={member.ID} className="flex items-center justify-between p-3 bg-yellow-50 rounded-lg">
                  <div>
                    <p className="font-medium">{member.User.name}</p>
                    <p className="text-sm text-gray-600">{member.User.email}</p>
                  </div>
                  <span className="px-2 py-1 bg-yellow-100 text-yellow-800 rounded text-xs font-medium">
                    Pending
                  </span>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>

      {/* Payout Schedule Section */}
      {group.Status === 'active' && group.PayoutOrder && (
        <div className="card">
          <h3 className="text-lg font-semibold mb-4">Payout Schedule</h3>

          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-4">
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <span className="text-blue-700 font-medium">Contribution Amount:</span>
                <span className="ml-2">{group.ContributionAmount} XLM</span>
              </div>
              <div>
                <span className="text-blue-700 font-medium">Contribution Period:</span>
                <span className="ml-2">Every {group.ContributionPeriod} days</span>
              </div>
              <div>
                <span className="text-blue-700 font-medium">Current Round:</span>
                <span className="ml-2">{group.CurrentRound || 1}</span>
              </div>
              <div>
                <span className="text-blue-700 font-medium">Next Contribution:</span>
                <span className="ml-2">
                  {group.NextContributionDate
                    ? new Date(group.NextContributionDate).toLocaleDateString()
                    : 'Not set'
                  }
                </span>
              </div>
            </div>
          </div>

          <div className="space-y-3">
            <h4 className="font-medium text-gray-700">Payout Schedule:</h4>
            {payoutSchedule.length > 0 ? (
              payoutSchedule.map((schedule: any) => {
                const isCurrentRound = (group.CurrentRound || 1) === schedule.Round
                const isPastRound = (group.CurrentRound || 1) > schedule.Round
                const dueDate = new Date(schedule.DueDate)

                return (
                  <div
                    key={schedule.ID}
                    className={`flex items-center justify-between p-4 rounded-lg border ${isCurrentRound ? 'bg-green-50 border-green-200' :
                        isPastRound ? 'bg-gray-50 border-gray-200' :
                          'bg-white border-gray-200'
                      }`}
                  >
                    <div className="flex items-center space-x-3">
                      <div className={`w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold ${isCurrentRound ? 'bg-green-500 text-white' :
                          isPastRound ? 'bg-gray-400 text-white' :
                            'bg-blue-100 text-blue-800'
                        }`}>
                        {schedule.Round}
                      </div>
                      <div>
                        <p className="font-medium text-lg">{schedule.Member.User.name}</p>
                        <p className="text-sm text-gray-600">{schedule.Member.User.email}</p>
                      </div>
                    </div>

                    <div className="text-center">
                      <p className="font-bold text-lg text-green-600">{schedule.Amount.toFixed(2)} XLM</p>
                      <p className="text-sm text-gray-600">
                        {schedule.Status === 'paid' ? 'Paid' :
                          schedule.Status === 'pending' ? 'Processing' : 'Scheduled'}
                      </p>
                    </div>

                    <div className="text-right">
                      <p className="font-medium text-lg">
                        {dueDate.toLocaleDateString('en-US', {
                          weekday: 'short',
                          year: 'numeric',
                          month: 'short',
                          day: 'numeric'
                        })}
                      </p>
                      <p className="text-sm text-gray-600">
                        {schedule.PaidAt ? `Paid ${new Date(schedule.PaidAt).toLocaleDateString()}` :
                          isCurrentRound ? 'Current Round' :
                            isPastRound ? 'Overdue' : 'Scheduled'}
                      </p>
                    </div>

                    {schedule.Status === 'paid' && (
                      <div className="ml-3">
                        <CheckCircle className="w-6 h-6 text-green-500" />
                      </div>
                    )}
                  </div>
                )
              })
            ) : (
              <div className="text-center py-4 text-gray-500">
                No payout schedule available
              </div>
            )}
          </div>

          <div className="mt-6 grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="p-4 bg-gray-50 rounded-lg">
              <h5 className="text-sm font-medium text-gray-700 mb-2">Schedule Summary:</h5>
              <ul className="text-xs text-gray-600 space-y-1">
                <li>• Total Rounds: {approvedMembers.length}</li>
                <li>• Contribution per Member: {group.ContributionAmount} XLM</li>
                <li>• Total Pool per Round: {((group.ContributionAmount || 0) * approvedMembers.length).toFixed(2)} XLM</li>
                <li>• Cycle Duration: {((group.ContributionPeriod || 30) * approvedMembers.length)} days</li>
              </ul>
            </div>

            <div className="p-4 bg-blue-50 rounded-lg">
              <h5 className="text-sm font-medium text-blue-700 mb-2">How it works:</h5>
              <ul className="text-xs text-blue-600 space-y-1">
                <li>• Each member contributes every {group.ContributionPeriod} days</li>
                <li>• One member receives the full pool each round</li>
                <li>• Payouts follow the predetermined order</li>
                <li>• Cycle completes when everyone has received a payout</li>
              </ul>
            </div>
          </div>
        </div>
      )}

      {/* Round Contributions Section */}
      {group.Status === 'active' && (
        <div className="mt-8">
          <RoundContributions group={group} currentUser={currentUser} />
        </div>
      )}

      {/* Invite Modal */}
      {showInviteModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold mb-4">Invite Member</h3>
            <form onSubmit={handleInvite}>
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Select User to Invite
                </label>
                <select
                  value={inviteEmail}
                  onChange={(e) => setInviteEmail(e.target.value)}
                  className="input"
                  required
                >
                  <option value="">Choose a user...</option>
                  {availableUsers.map((user) => (
                    <option key={user.id} value={user.email}>
                      {user.name} ({user.email})
                    </option>
                  ))}
                </select>
              </div>
              <div className="flex space-x-3">
                <button
                  type="button"
                  onClick={() => setShowInviteModal(false)}
                  className="btn btn-secondary flex-1"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="btn btn-primary flex-1"
                  disabled={inviteUserMutation.isPending}
                >
                  {inviteUserMutation.isPending ? 'Sending...' : 'Send Invite'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Approve Group Modal */}
      {showApproveModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold mb-4">Approve Group</h3>
            <div className="mb-6">
              <p className="text-gray-600 mb-4">
                By approving this group, you confirm that:
              </p>
              <ul className="text-sm text-gray-600 space-y-2">
                <li>• The group has reached its maximum capacity</li>
                <li>• All members are verified and trusted</li>
                <li>• Members can now nominate admins</li>
                <li>• The group can be activated for contributions</li>
              </ul>
            </div>
            <div className="flex space-x-3">
              <button
                type="button"
                onClick={() => setShowApproveModal(false)}
                className="btn btn-secondary flex-1"
              >
                Cancel
              </button>
              <button
                onClick={handleApproveGroup}
                className="btn btn-primary flex-1"
                disabled={approveGroupMutation.isPending}
              >
                {approveGroupMutation.isPending ? 'Approving...' : 'Approve Group'}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Activate Group Modal */}
      {showActivateModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold mb-4">Activate Group</h3>
            <form onSubmit={handleActivate}>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Contribution Amount (XLM)
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    min="0"
                    value={groupSettings.contribution_amount || ''}
                    onChange={(e) => {
                      const value = e.target.value === '' ? 0 : parseFloat(e.target.value)
                      setGroupSettings(prev => ({
                        ...prev,
                        contribution_amount: isNaN(value) ? 0 : value
                      }))
                    }}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="Enter amount"
                    required
                  />
                </div>
                
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Contribution Period (days)
                  </label>
                  <select
                    value={groupSettings.contribution_period}
                    onChange={(e) => setGroupSettings(prev => ({
                      ...prev,
                      contribution_period: parseInt(e.target.value)
                    }))}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value={7}>Weekly (7 days)</option>
                    <option value={14}>Bi-weekly (14 days)</option>
                    <option value={30}>Monthly (30 days)</option>
                  </select>
                </div>
                
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Payout Order (Set the order members will receive payouts)
                  </label>
                  <div className="space-y-2 max-h-40 overflow-y-auto border rounded-lg p-3 bg-gray-50">
                    {groupSettings.payout_order.length > 0 ? (
                      groupSettings.payout_order.map((userId, index) => {
                        const member = approvedMembers.find(m => m.UserID === userId)
                        if (!member) return null

                        return (
                          <div key={`${userId}-${index}`} className="flex items-center justify-between p-2 bg-white rounded border">
                            <div className="flex items-center space-x-2">
                              <span className="w-6 h-6 bg-blue-100 text-blue-800 rounded-full flex items-center justify-center text-xs font-medium">
                                {index + 1}
                              </span>
                              <span className="text-sm font-medium">{member.User.name}</span>
                            </div>
                            <div className="flex items-center space-x-1">
                              <button
                                type="button"
                                onClick={() => {
                                  console.log('Move up clicked, current index:', index)
                                  console.log('Current payout order:', groupSettings.payout_order)
                                  
                                  if (index > 0) {
                                    const newOrder = [...groupSettings.payout_order]
                                    const temp = newOrder[index]
                                    newOrder[index] = newOrder[index - 1]
                                    newOrder[index - 1] = temp
                                    
                                    console.log('New payout order:', newOrder)
                                    
                                    setGroupSettings(prev => ({
                                      ...prev,
                                      payout_order: newOrder
                                    }))
                                  }
                                }}
                                disabled={index === 0}
                                className="text-xs px-2 py-1 bg-blue-100 text-blue-800 rounded disabled:opacity-50 hover:bg-blue-200 disabled:cursor-not-allowed"
                              >
                                ↑
                              </button>
                              <button
                                type="button"
                                onClick={() => {
                                  console.log('Move down clicked, current index:', index)
                                  console.log('Current payout order:', groupSettings.payout_order)
                                  
                                  if (index < groupSettings.payout_order.length - 1) {
                                    const newOrder = [...groupSettings.payout_order]
                                    const temp = newOrder[index]
                                    newOrder[index] = newOrder[index + 1]
                                    newOrder[index + 1] = temp
                                    
                                    console.log('New payout order:', newOrder)
                                    
                                    setGroupSettings(prev => ({
                                      ...prev,
                                      payout_order: newOrder
                                    }))
                                  }
                                }}
                                disabled={index === groupSettings.payout_order.length - 1}
                                className="text-xs px-2 py-1 bg-blue-100 text-blue-800 rounded disabled:opacity-50 hover:bg-blue-200 disabled:cursor-not-allowed"
                              >
                                ↓
                              </button>
                            </div>
                          </div>
                        )
                      })
                    ) : (
                      <div className="text-center py-2 text-gray-500 text-sm">
                        No members available for payout order
                      </div>
                    )}
                  </div>
                  <p className="text-xs text-gray-500 mt-2">
                    Use the arrows to reorder. Member #1 will receive the first payout, #2 the second, etc.
                  </p>
                </div>
                
                {/* Preview section */}
                <div className="bg-blue-50 border border-blue-200 rounded-lg p-3">
                  <h5 className="text-sm font-medium text-blue-700 mb-2">Preview:</h5>
                  <div className="text-xs text-blue-600 space-y-1">
                    <div>• Contribution per member: {groupSettings.contribution_amount || 0} XLM</div>
                    <div>• Total pool per round: {((groupSettings.contribution_amount || 0) * approvedMembers.length).toFixed(2)} XLM</div>
                    <div>• Payment every: {groupSettings.contribution_period} days</div>
                    <div>• Total cycle duration: {(groupSettings.contribution_period * approvedMembers.length)} days</div>
                  </div>
                </div>
              </div>
              
              <div className="flex space-x-3 mt-6">
                <button
                  type="button"
                  onClick={() => {
                    setShowActivateModal(false)
                    // Reset form when closing
                    setGroupSettings({
                      contribution_amount: 0,
                      contribution_period: 30,
                      payout_order: []
                    })
                  }}
                  className="btn btn-secondary flex-1"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="btn btn-primary flex-1"
                  disabled={activateGroupMutation.isPending || groupSettings.contribution_amount <= 0 || groupSettings.payout_order.length === 0}
                >
                  {activateGroupMutation.isPending ? 'Activating...' : 'Activate Group'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}

export default GroupManagement
