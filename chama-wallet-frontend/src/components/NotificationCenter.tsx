import React, { useState } from 'react'
import { Bell, Check, X, Users, DollarSign, UserPlus } from 'lucide-react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '../services/api'
import type { Notification, GroupInvitation } from '../types'
import toast from 'react-hot-toast'

const NotificationCenter = () => {
  const [showNotifications, setShowNotifications] = useState(false)
  const queryClient = useQueryClient()

  const { data: notifications = [] } = useQuery({
    queryKey: ['notifications'],
    queryFn: () => api.get<Notification[]>('/notifications').then(res => res.data)
  })

  const { data: invitations = [] } = useQuery({
    queryKey: ['invitations'],
    queryFn: () => api.get<GroupInvitation[]>('/invitations').then(res => res.data)
  })

  const markAsReadMutation = useMutation({
    mutationFn: (id: string) => api.put(`/notifications/${id}/read`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] })
    }
  })

  const acceptInvitationMutation = useMutation({
    mutationFn: (id: string) => api.post(`/invitations/${id}/accept`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['invitations'] })
      queryClient.invalidateQueries({ queryKey: ['groups'] })
      toast.success('Invitation accepted!')
    }
  })

  const rejectInvitationMutation = useMutation({
    mutationFn: (id: string) => api.post(`/invitations/${id}/reject`),
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
        className="relative p-2 text-gray-600 hover:text-gray-900 transition-colors"
      >
        <Bell className="w-6 h-6" />
        {unreadCount > 0 && (
          <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
            {unreadCount > 9 ? '9+' : unreadCount}
          </span>
        )}
      </button>

      {showNotifications && (
        <div className="absolute right-0 mt-2 w-96 bg-white rounded-lg shadow-lg border border-gray-200 z-50 max-h-96 overflow-y-auto">
          <div className="p-4 border-b border-gray-200">
            <h3 className="text-lg font-semibold">Notifications</h3>
          </div>

          {/* Invitations */}
          {invitations.length > 0 && (
            <div className="border-b border-gray-200">
              <div className="p-3 bg-blue-50">
                <h4 className="font-medium text-blue-900">Group Invitations</h4>
              </div>
              {invitations.map((invitation) => (
                <div key={invitation.ID} className="p-4 border-b border-gray-100 last:border-b-0">
                  <div className="flex items-start space-x-3">
                    <UserPlus className="w-5 h-5 text-blue-500 mt-1" />
                    <div className="flex-1">
                      <p className="text-sm font-medium">
                        Invitation to join "{invitation.Group.Name}"
                      </p>
                      <p className="text-xs text-gray-600 mt-1">
                        From {invitation.Inviter.name}
                      </p>
                      <div className="flex space-x-2 mt-3">
                        <button
                          onClick={() => handleAcceptInvitation(invitation.ID)}
                          className="px-3 py-1 bg-green-600 text-white text-xs rounded hover:bg-green-700"
                          disabled={acceptInvitationMutation.isPending}
                        >
                          Accept
                        </button>
                        <button
                          onClick={() => handleRejectInvitation(invitation.ID)}
                          className="px-3 py-1 bg-gray-600 text-white text-xs rounded hover:bg-gray-700"
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
                  className={`p-4 border-b border-gray-100 last:border-b-0 ${
                    !notification.Read ? 'bg-blue-50' : 'bg-white'
                  }`}
                >
                  <div className="flex items-start space-x-3">
                    {getNotificationIcon(notification.Type)}
                    <div className="flex-1">
                      <p className="text-sm font-medium">{notification.Title}</p>
                      <p className="text-xs text-gray-600 mt-1">{notification.Message}</p>
                      <p className="text-xs text-gray-400 mt-2">
                        {new Date(notification.CreatedAt).toLocaleDateString()}
                      </p>
                      {!notification.Read && (
                        <button
                          onClick={() => handleMarkAsRead(notification.ID)}
                          className="text-xs text-blue-600 hover:text-blue-800 mt-2"
                        >
                          Mark as read
                        </button>
                      )}
                    </div>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  )
}

export default NotificationCenter