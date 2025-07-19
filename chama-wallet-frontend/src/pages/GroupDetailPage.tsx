import React, { useState } from 'react'
import { useParams } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { useGroups, useGroupBalance } from '../hooks/useGroups'
import LoadingSpinner from '../components/LoadingSpinner'
import GroupManagement from '../components/GroupManagement'
import AdminNomination from '../components/AdminNomination'
import PayoutManagement from '../components/PayoutManagement'
import { Users, Wallet, Settings, Crown, DollarSign } from 'lucide-react'

const GroupDetailPage = () => {
  const { id } = useParams<{ id: string }>()
  const { user } = useAuth()
  const { data: groups, isLoading: groupsLoading } = useGroups()
  const { data: balance, isLoading: balanceLoading } = useGroupBalance(id!)
  const [activeTab, setActiveTab] = useState('overview')

  const group = groups?.find(g => g.ID === id)

  if (groupsLoading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <LoadingSpinner size="lg" />
      </div>
    )
  }

  if (!group) {
    return (
      <div className="text-center py-8">
        <h2 className="text-2xl font-bold text-gray-900">Group not found</h2>
        <p className="text-gray-600 mt-2">You don't have access to this group.</p>
      </div>
    )
  }

  const currentUserMember = group.Members?.find(m => m.UserID === user?.id)
  const isAdmin = currentUserMember && ['creator', 'admin'].includes(currentUserMember.Role)

  const tabs = [
    { id: 'overview', label: 'Overview', icon: Users },
    { id: 'management', label: 'Management', icon: Settings },
    { id: 'nominations', label: 'Admin Nominations', icon: Crown },
    { id: 'payouts', label: 'Payouts', icon: DollarSign },
  ]

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      {/* Group Header */}
      <div className="card mb-8">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">{group.Name}</h1>
            <p className="text-gray-600 mt-2">{group.Description}</p>
          </div>
          <div className="text-right">
            <div className="flex items-center space-x-2 mb-2">
              <Wallet className="w-5 h-5 text-stellar-600" />
              <span className="text-lg font-semibold">
                {balanceLoading ? 'Loading...' : `${balance?.data?.balance || '0'} XLM`}
              </span>
            </div>
            <span className={`px-3 py-1 rounded-full text-sm font-medium ${
              group.Status === 'active' ? 'bg-green-100 text-green-800' :
              group.Status === 'pending' ? 'bg-yellow-100 text-yellow-800' :
              'bg-gray-100 text-gray-800'
            }`}>
              {group.Status.charAt(0).toUpperCase() + group.Status.slice(1)}
            </span>
          </div>
        </div>

        {/* Group Stats */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mt-6">
          <div className="bg-blue-50 p-4 rounded-lg">
            <h3 className="text-sm font-medium text-blue-900">Total Members</h3>
            <p className="text-2xl font-bold text-blue-600">
              {group.Members?.filter(m => m.Status === 'approved').length || 0}
            </p>
          </div>
          <div className="bg-green-50 p-4 rounded-lg">
            <h3 className="text-sm font-medium text-green-900">Contribution Amount</h3>
            <p className="text-2xl font-bold text-green-600">
              {group.ContributionAmount || 0} XLM
            </p>
          </div>
          <div className="bg-purple-50 p-4 rounded-lg">
            <h3 className="text-sm font-medium text-purple-900">Current Round</h3>
            <p className="text-2xl font-bold text-purple-600">
              {group.CurrentRound || 0}
            </p>
          </div>
          <div className="bg-yellow-50 p-4 rounded-lg">
            <h3 className="text-sm font-medium text-yellow-900">Contribution Period</h3>
            <p className="text-2xl font-bold text-yellow-600">
              {group.ContributionPeriod || 0} days
            </p>
          </div>
        </div>
      </div>

      {/* Navigation Tabs */}
      <div className="border-b border-gray-200 mb-8">
        <nav className="-mb-px flex space-x-8">
          {tabs.map((tab) => {
            const Icon = tab.icon
            return (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`py-2 px-1 border-b-2 font-medium text-sm flex items-center space-x-2 ${
                  activeTab === tab.id
                    ? 'border-primary-500 text-primary-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }`}
              >
                <Icon className="w-4 h-4" />
                <span>{tab.label}</span>
              </button>
            )
          })}
        </nav>
      </div>

      {/* Tab Content */}
      <div>
        {activeTab === 'overview' && (
          <div className="space-y-6">
            {/* Members List */}
            <div className="card">
              <h3 className="text-lg font-semibold mb-4">Group Members</h3>
              <div className="space-y-3">
                {group.Members?.filter(m => m.Status === 'approved').map((member) => (
                  <div key={member.ID} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                    <div>
                      <p className="font-medium">{member.User.name}</p>
                      <p className="text-sm text-gray-600">{member.User.email}</p>
                    </div>
                    <span className={`px-2 py-1 rounded text-xs font-medium ${
                      member.Role === 'creator' ? 'bg-purple-100 text-purple-800' :
                      member.Role === 'admin' ? 'bg-blue-100 text-blue-800' :
                      'bg-gray-100 text-gray-800'
                    }`}>
                      {member.Role}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}

        {activeTab === 'management' && user && (
          <GroupManagement 
            group={group} 
            currentUser={user}
            onInviteUser={() => {}} // Implement these handlers
            onActivateGroup={() => {}}
          />
        )}

        {activeTab === 'nominations' && user && (
          <AdminNomination group={group} currentUser={user} />
        )}

        {activeTab === 'payouts' && user && (
          <PayoutManagement group={group} currentUser={user} />
        )}
      </div>
    </div>
  )
}

export default GroupDetailPage
