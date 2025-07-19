import React, { useState } from 'react'
import { Users, Settings, UserPlus, CheckCircle } from 'lucide-react'
import type { Group, User } from '../types'

interface GroupManagementProps {
  group: Group
  currentUser: User
  onInviteUser: (email: string) => void
  onActivateGroup: (settings: any) => void
}

const GroupManagement: React.FC<GroupManagementProps> = ({
  group,
  currentUser,
  onInviteUser,
  onActivateGroup
}) => {
  const [showInviteModal, setShowInviteModal] = useState(false)
  const [showActivateModal, setShowActivateModal] = useState(false)
  const [inviteEmail, setInviteEmail] = useState('')
  const [groupSettings, setGroupSettings] = useState({
    contribution_amount: 0,
    contribution_period: 30,
    payout_order: ''
  })

  const isAdmin = group.Members?.find(m => 
    m.UserID === currentUser.id && ['creator', 'admin'].includes(m.Role)
  )

  const approvedMembers = group.Members?.filter(m => m.Status === 'approved') || []
  const pendingMembers = group.Members?.filter(m => m.Status === 'pending') || []

  const handleInvite = (e: React.FormEvent) => {
    e.preventDefault()
    onInviteUser(inviteEmail)
    setInviteEmail('')
    setShowInviteModal(false)
  }

  const handleActivate = (e: React.FormEvent) => {
    e.preventDefault()
    const memberIds = approvedMembers.map(m => m.UserID)
    const shuffledOrder = [...memberIds].sort(() => Math.random() - 0.5)
    
    onActivateGroup({
      ...groupSettings,
      payout_order: JSON.stringify(shuffledOrder)
    })
    setShowActivateModal(false)
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
        
        {group.Status === 'pending' && isAdmin && (
          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <p className="text-blue-800 mb-3">
              Your group is ready to be activated. Set contribution rules and activate the group.
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
      </div>

      {/* Members Section */}
      <div className="card">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold">
            Members ({approvedMembers.length})
          </h3>
          {isAdmin && (
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
                  Email Address
                </label>
                <input
                  type="email"
                  value={inviteEmail}
                  onChange={(e) => setInviteEmail(e.target.value)}
                  className="input"
                  placeholder="Enter email address"
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
              </div>
              <div className="flex space-x-3 mt-6">
                <button
                  type="button"
                  onClick={() => setShowActivateModal(false)}
                  className="btn btn-secondary flex-1"
                >
                  Cancel
                </button>
                <button type="submit" className="btn btn-primary flex-1">
                  Activate Group
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