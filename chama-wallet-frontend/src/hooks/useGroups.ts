import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { groupApi } from '../services/api'
import toast from 'react-hot-toast'

export const useGroups = () => {
  return useQuery({
    queryKey: ['groups'],
    queryFn: () => groupApi.getAllGroups(),
  })
}

export const useGroupBalance = (id: string) => {
  return useQuery({
    queryKey: ['group-balance', id],
    queryFn: () => groupApi.getGroupBalance(id),
    enabled: !!id,
  })
}

export const useGroupMutations = () => {
  const queryClient = useQueryClient()

  const createGroup = useMutation({
    mutationFn: groupApi.createGroup,
    onSuccess: () => {
      toast.success('Group created successfully!')
      queryClient.invalidateQueries({ queryKey: ['groups'] })
    },
    onError: () => {
      toast.error('Failed to create group')
    },
  })

  const joinGroup = useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => 
      groupApi.joinGroup(id, data),
    onSuccess: () => {
      toast.success('Joined group successfully!')
      queryClient.invalidateQueries({ queryKey: ['groups'] })
    },
    onError: () => {
      toast.error('Failed to join group')
    },
  })

  const contributeToGroup = useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => 
      groupApi.contributeToGroup(id, data),
    onSuccess: () => {
      toast.success('Contribution successful!')
      queryClient.invalidateQueries({ queryKey: ['group-balance'] })
    },
    onError: () => {
      toast.error('Contribution failed')
    },
  })

  return {
    createGroup,
    joinGroup,
    contributeToGroup,
  }
}