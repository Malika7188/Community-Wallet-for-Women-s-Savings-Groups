import React, { useState } from 'react'
import { Bell, Check, X, Users, DollarSign, UserPlus } from 'lucide-react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api, { notificationApi } from '../services/api'
import type { Notification, GroupInvitation } from '../types'
import toast from 'react-hot-toast'

interface NotificationCenterProps {
  isCollapsed: boolean;
}

const NotificationCenter: React.FC<NotificationCenterProps> = ({ isCollapsed }) => {
  const [showNotifications, setShowNotifications] = useState(false)
  const [selected, setSelected] = useState<string[]>([])
  const queryClient = useQueryClient()
  const clearNotificationMutation = useMutation({
    mutationFn: (id: string) => notificationApi.clearNotification(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] })
      toast.success('Notification cleared')
    },
    onError: () => {
      toast.error('Failed to clear notification')
    }
  })

  const clearSelected = () => {
    selected.forEach(id => clearNotificationMutation.mutate(id))
    setSelected([])
  }

  const { data: notifications = [] } = useQuery({
    queryKey: ['notifications'],
    queryFn: () => {
      console.log('🔍 Fetching notifications...')
      return notificationApi.getNotifications().then((res: { data: Notification[] }) => {
        console.log('✅ Notifications received:', res.data)
        return res.data
      })
    },
    refetchInterval: 30000, // Refetch every 30 seconds
  })

  const { data: invitations = [] } = useQuery({
    queryKey: ['invitations'],
    queryFn: () => {
      console.log('🔍 Fetching invitations...')
      return notificationApi.getInvitations().then((res: { data: GroupInvitation[] }) => {
        console.log('✅ Invitations received:', res.data)
        return res.data
      })
    },
    refetchInterval: 30000, // Refetch every 30 seconds
  })

  const markAsReadMutation = useMutation({
    mutationFn: (id: string) => notificationApi.markAsRead(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] })
    }
  })

  const acceptInvitationMutation = useMutation({
    mutationFn: (id: string) => notificationApi.acceptInvitation(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['invitations'] })
      queryClient.invalidateQueries({ queryKey: ['groups'] })
      queryClient.invalidateQueries({ queryKey: ['userGroups'] }) // Add this if you have user-specific groups
      toast.success('Invitation accepted! You are now a member of the group.')
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to accept invitation')
    }
  })

  const rejectInvitationMutation = useMutation({
    mutationFn: (id: string) => notificationApi.rejectInvitation(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['invitations'] })
      toast.success('Invitation rejected')
    }
  })

  const unreadCount = notifications.filter(n => !n.Read).length + invitations.length

  const getNotificationIcon = (type: string) => {
    switch (type) {
      case 'contribution_reminder':
        return <DollarSign className="w-5 h-5 text-yellow-500" />
      case 'payout_approved':
        return <Check className="w-5 h-5 text-green-500" />
      case 'new_member_request':
        return <UserPlus className="w-5 h-5 text-blue-500" />
      case 'admin_promotion':
        return <Users className="w-5 h-5 text-purple-500" />
      default:
        return <Bell className="w-5 h-5 text-gray-500" />
    }
  }

  const handleMarkAsRead = (id: string) => {
    markAsReadMutation.mutate(id)
  }

  const handleAcceptInvitation = (id: string) => {
    acceptInvitationMutation.mutate(id)
  }

  const handleRejectInvitation = (id: string) => {
    rejectInvitationMutation.mutate(id)
  }

  return (
    <div className="relative">
      <button
        onClick={() => setShowNotifications(!showNotifications)}
        className={`flex items-center w-full px-4 py-3 rounded-lg transition-colors duration-200 text-white hover:bg-[#2ecc71] hover:text-[#1a237e] ${showNotifications ? 'bg-[#2ecc71] text-[#1a237e]' : ''}`}
        style={{ outline: 'none' }}
      >
        <Bell className="h-6 w-6" />
        {!isCollapsed && <span className="font-medium ml-3">Notifications</span>}
        {unreadCount > 0 && (
          <span className={`absolute top-2 right-6 bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center ${isCollapsed ? 'right-2' : ''}`}>
            {unreadCount > 9 ? '9+' : unreadCount}
          </span>
        )}
      </button>

      {showNotifications && (
  <div className="absolute left-0 top-12 w-80 backdrop-blur-lg bg-white/80 border border-gray-200 rounded-2xl shadow-2xl z-50 max-h-96 overflow-y-auto transition-all duration-300">
          <div className="p-4 border-b border-gray-200 flex items-center justify-between">
            <h3 className="text-lg font-semibold text-[#1a237e]">Notifications</h3>
            <span className="text-xs text-gray-500">{unreadCount > 0 ? `${unreadCount} unread` : 'All read'}</span>
          </div>

          {/* Invitations */}
          {invitations.length > 0 && (
            <div className="border-b border-gray-200">
              <div className="p-3 bg-blue-100 rounded-t-xl">
                <h4 className="font-medium text-blue-900">Group Invitations</h4>
              </div>
              {invitations.map((invitation) => (
                <div key={invitation.ID} className="p-4 border-b border-gray-100 last:border-b-0 bg-white/70 hover:bg-blue-50 transition rounded-xl">
                  <div className="flex items-start space-x-3">
                    <UserPlus className="w-5 h-5 text-blue-500 mt-1" />
                    <div className="flex-1">
                      <p className="text-sm font-medium text-[#1a237e]">
                        Invitation to join "{invitation.Group.Name}"
                      </p>
                      <p className="text-xs text-gray-600 mt-1">
                        From {invitation.Inviter.name}
                      </p>
                      <div className="flex space-x-2 mt-3">
                        <button
                          onClick={() => handleAcceptInvitation(invitation.ID)}
                          className="px-3 py-1 bg-green-600 text-white text-xs rounded-lg shadow hover:bg-green-700 focus:outline-none"
                          disabled={acceptInvitationMutation.isPending}
                        >
                          Accept
                        </button>
                        <button
                          onClick={() => handleRejectInvitation(invitation.ID)}
                          className="px-3 py-1 bg-gray-600 text-white text-xs rounded-lg shadow hover:bg-gray-700 focus:outline-none"
                          disabled={rejectInvitationMutation.isPending}
                        >
                          Reject
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}

          {/* Regular Notifications */}
          <div className="max-h-64 overflow-y-auto">
            {notifications.length === 0 && invitations.length === 0 ? (
              <div className="p-4 text-center text-gray-500">
                No notifications
              </div>
            ) : (
              notifications.map((notification) => (
                <div
                  key={notification.ID}
                  className={`p-4 border-b border-gray-100 last:border-b-0 rounded-xl transition flex items-start gap-3 ${!notification.Read ? 'bg-blue-50' : 'bg-white/70 hover:bg-blue-50'}`}
                >
                  <input
                    type="checkbox"
                    checked={selected.includes(notification.ID)}
                    onChange={e => {
                      if (e.target.checked) setSelected([...selected, notification.ID])
                      else setSelected(selected.filter(id => id !== notification.ID))
                    }}
                    className="mt-1 accent-[#2ecc71]"
                  />
                  {getNotificationIcon(notification.Type)}
                  <div className="flex-1">
                    <p className="text-sm font-medium text-[#1a237e]">{notification.Title}</p>
                    <p className="text-xs text-gray-600 mt-1">{notification.Message}</p>
                    <p className="text-xs text-gray-400 mt-2">
                      {new Date(notification.CreatedAt).toLocaleDateString()}
                    </p>
                    {!notification.Read && (
                      <button
                        onClick={() => handleMarkAsRead(notification.ID)}
                        className="text-xs text-blue-600 hover:text-blue-800 mt-2 font-semibold"
                      >
                        Mark as read
                      </button>
                    )}
                  </div>
                </div>
              ))
            )}
            {selected.length > 0 && (
              <div className="p-4 flex justify-end">
                <button
                  onClick={clearSelected}
                  className="px-4 py-2 bg-red-500 text-white rounded-lg shadow hover:bg-red-700 transition font-semibold"
                >
                  Clear Selected
                </button>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}

export default NotificationCenter
