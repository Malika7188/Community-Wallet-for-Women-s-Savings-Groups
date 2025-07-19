import React, { useState } from 'react'
import { Users, Crown, Vote } from 'lucide-react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { groupApi } from '../services/api'
import type { Group, Member, User } from '../types'
import toast from 'react-hot-toast'

interface AdminNominationProps {
  group: Group
  currentUser: User
}

const AdminNomination: React.FC<AdminNominationProps> = ({ group, currentUser }) => {
  const [showNominationModal, setShowNominationModal] = useState(false)
  const [selectedMember, setSelectedMember] = useState('')
  const queryClient = useQueryClient()

  const nominateAdminMutation = useMutation({
    mutationFn: (data: { nominee_id: string }) => 
      groupApi.nominateAdmin(group.ID, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['groups'] })
      toast.success('Admin nomination submitted!')
      setShowNominationModal(false)
      setSelectedMember('')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to nominate admin')
    }
  })

  const eligibleMembers = group.Members?.filter(m => 
    m.Status === 'approved' && 
    m.Role === 'member' && 
    m.UserID !== currentUser.id
  ) || []

  const admins = group.Members?.filter(m => 
    m.Status === 'approved' && 
    ['admin', 'creator'].includes(m.Role)
  ) || []

  const currentUserMember = group.Members?.find(m => m.UserID === currentUser.id)
  const canNominate = currentUserMember?.Status === 'approved'

  const handleNominate = (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedMember) return

    nominateAdminMutation.mutate({ nominee_id: selectedMember })
  }

  return (
    <div className="card">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-semibold flex items-center">
          <Crown className="w-5 h-5 mr-2 text-yellow-500" />
          Group Administration
        </h3>
        {canNominate && eligibleMembers.length > 0 && (
          <button
            onClick={() => setShowNominationModal(true)}
            className="btn btn-outline"
          >
            <Vote className="w-4 h-4 mr-2" />
            Nominate Admin
          </button>
        )}
      </div>

      {/* Current Admins */}
      <div className="mb-6">
        <h4 className="font-medium text-gray-700 mb-3">Current Administrators</h4>
        <div className="space-y-2">
          {admins.map((admin) => (
            <div key={admin.ID} className="flex items-center justify-between p-3 bg-yellow-50 rounded-lg">
              <div>
                <p className="font-medium">{admin.User.name}</p>
                <p className="text-sm text-gray-600">{admin.User.email}</p>
              </div>
              <span className={`px-2 py-1 rounded text-xs font-medium ${
                admin.Role === 'creator' ? 'bg-purple-100 text-purple-800' :
                'bg-yellow-100 text-yellow-800'
              }`}>
                {admin.Role === 'creator' ? 'Creator' : 'Admin'}
              </span>
            </div>
          ))}
        </div>
      </div>

      {/* Nomination Info */}
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
        <h4 className="font-medium text-blue-900 mb-2">Admin Nomination Process</h4>
        <ul className="text-sm text-blue-800 space-y-1">
          <li>• Any approved member can nominate another member for admin role</li>
          <li>• A member needs 2 nominations to become an admin</li>
          <li>• Admins can approve new members and manage group settings</li>
          <li>• At least 2 admins must approve payouts</li>
        </ul>
      </div>

      {/* Nomination Modal */}
      {showNominationModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold mb-4">Nominate Admin</h3>
            <form onSubmit={handleNominate}>
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Select Member to Nominate
                </label>
                <select
                  value={selectedMember}
                  onChange={(e) => setSelectedMember(e.target.value)}
                  className="input"
                  required
                >
                  <option value="">Choose a member...</option>
                  {eligibleMembers.map((member) => (
                    <option key={member.ID} value={member.UserID}>
                      {member.User.name} ({member.User.email})
                    </option>
                  ))}
                </select>
              </div>
              
              <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-3 mb-4">
                <p className="text-sm text-yellow-800">
                  <strong>Note:</strong> This member will become an admin if they receive 
                  2 nominations from different members.
                </p>
              </div>

              <div className="flex space-x-3">
                <button
                  type="button"
                  onClick={() => setShowNominationModal(false)}
                  className="btn btn-secondary flex-1"
                >
                  Cancel
                </button>
                <button 
                  type="submit" 
                  className="btn btn-primary flex-1"
                  disabled={nominateAdminMutation.isPending}
                >
                  {nominateAdminMutation.isPending ? 'Nominating...' : 'Nominate'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}

export default AdminNomination