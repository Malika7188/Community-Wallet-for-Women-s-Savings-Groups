import React, { useState } from 'react'
import { Users, Settings, UserPlus, CheckCircle } from 'lucide-react'
import type { Group, User } from '../types'

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
  const isGroupFull = approvedMembers.length >= (group.MaxMembers || 10)

  const handleInvite = (e: React.FormEvent) => {
    e.preventDefault()
    inviteUserMutation.mutate(inviteEmail)
  }

  const handleApproveGroup = () => {
    approveGroupMutation.mutate()
  }

  const handleActivate = (e: React.FormEvent) => {
    e.preventDefault()
    activateGroupMutation.mutate({
      ...groupSettings,
    })
  }

  return (
    <div className="space-y-6">
      {/* Group Status */}
      <div className="card">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold">Group Status</h3>
          <span className={`px-3 py-1 rounded-full text-sm font-medium ${
            group.Status === 'active' ? 'bg-green-100 text-green-800' :
            group.Status === 'pending' ? 'bg-yellow-100 text-yellow-800' :
            'bg-gray-100 text-gray-800'
          }`}>
            {group.Status.charAt(0).toUpperCase() + group.Status.slice(1)}
          </span>
        </div>
        
        {/* Group Full - Ready for Approval */}
        {group.Status === 'pending' && !group.IsApproved && isCreator && isGroupFull && (
          <div className="bg-green-50 border border-green-200 rounded-lg p-4 mb-4">
            <p className="text-green-800 mb-3">
              🎉 Your group is now full! You can approve it to allow admin nominations and activation.
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
              <span className="ml-2 font-medium">{approvedMembers.length}/{group.MaxMembers || 10}</span>
            </div>
            <div>
              <span className="text-gray-600">Approved:</span>
              <span className="ml-2 font-medium">{group.IsApproved ? 'Yes' : 'No'}</span>
            </div>
          </div>
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
                <span className={`px-2 py-1 rounded text-xs font-medium ${
                  member.Role === 'creator' ? 'bg-purple-100 text-purple-800' :
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
                    value={groupSettings.contribution_amount}
                    onChange={(e) => setGroupSettings({
                      ...groupSettings,
                      contribution_amount: parseFloat(e.target.value)
                    })}
                    className="input"
                  required
                />
              </div>
              <div className="flex space-x-3">
                <button
                  type="button"
                  onClick={() => setShowInviteModal(false)}
                  className="btn btn-secondary flex-1"
                >
                  Cancel
                </button>
                <button type="submit" className="btn btn-primary flex-1">
                  Send Invite
                </button>
              </div>
            </form>
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
                    value={groupSettings.contribution_amount}
                    onChange={(e) => setGroupSettings({
                      ...groupSettings,
                      contribution_amount: parseFloat(e.target.value)
                    })}
                    className="input"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Contribution Period (days)
                  </label>
                  <select
                    value={groupSettings.contribution_period}
                    onChange={(e) => setGroupSettings({
                      ...groupSettings,
                      contribution_period: parseInt(e.target.value)
                    })}
                    className="input"
                  >
                    <option value={7}>Weekly</option>
                    <option value={14}>Bi-weekly</option>
                    <option value={30}>Monthly</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Payout Order
                  </label>
                  <div className="space-y-2 max-h-32 overflow-y-auto">
                    {approvedMembers.map((member, index) => (
                      <div key={member.ID} className="flex items-center justify-between p-2 bg-gray-50 rounded">
                        <span className="text-sm">{member.User.name}</span>
                        <input
                          type="number"
                          min="1"
                          max={approvedMembers.length}
                          value={groupSettings.payout_order.indexOf(member.UserID) + 1 || index + 1}
                          onChange={(e) => {
                            const newOrder = [...groupSettings.payout_order]
                            const position = parseInt(e.target.value) - 1
                            newOrder[position] = member.UserID
                            setGroupSettings({
                              ...groupSettings,
                              payout_order: newOrder
                            })
                          }}
                          className="w-16 px-2 py-1 border rounded text-sm"
                        />
                      </div>
                    ))}
                  </div>
                </div>
              </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Contribution Period (days)
                  </label>
                  <select
                    value={groupSettings.contribution_period}
                    onChange={(e) => setGroupSettings({
                      ...groupSettings,
                      contribution_period: parseInt(e.target.value)
                    })}
                    className="input"
                  >
                    <option value={7}>Weekly</option>
                    <option value={14}>Bi-weekly</option>
                    <option value={30}>Monthly</option>
                  </select>
                </div>
              </div>
              <div className="flex space-x-3 mt-6">
                <button
                  type="button"
                  onClick={() => setShowActivateModal(false)}
                  className="btn btn-secondary flex-1"
                >
                  Cancel
                </button>
                <button 
                  type="submit" 
                  className="btn btn-primary flex-1"
                  disabled={activateGroupMutation.isPending}
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