import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { walletApi } from '../services/api'
import toast from 'react-hot-toast'

export const useWallet = () => {
  const queryClient = useQueryClient()

  const createWallet = useMutation({
    mutationFn: walletApi.createWallet,
    onSuccess: () => {
      toast.success('Wallet created successfully!')
      queryClient.invalidateQueries({ queryKey: ['wallet'] })
    },
    onError: () => {
      toast.error('Failed to create wallet')
    },
  })

  const generateKeypair = useMutation({
    mutationFn: walletApi.generateKeypair,
    onSuccess: () => {
      toast.success('Keypair generated successfully!')
    },
    onError: () => {
      toast.error('Failed to generate keypair')
    },
  })

  const fundAccount = useMutation({
    mutationFn: (address: string) => walletApi.fundAccount(address),
    onSuccess: () => {
      toast.success('Account funded successfully!')
      queryClient.invalidateQueries({ queryKey: ['balance'] })
    },
    onError: () => {
      toast.error('Failed to fund account')
    },
  })

  const transferFunds = useMutation({
    mutationFn: walletApi.transferFunds,
    onSuccess: () => {
      toast.success('Transfer completed successfully!')
      queryClient.invalidateQueries({ queryKey: ['balance'] })
      queryClient.invalidateQueries({ queryKey: ['transactions'] })
    },
    onError: () => {
      toast.error('Transfer failed')
    },
  })

  return {
    createWallet,
    generateKeypair,
    fundAccount,
    transferFunds,
  }
}

export const useBalance = (address: string) => {
  return useQuery({
    queryKey: ['balance', address],
    queryFn: () => walletApi.getBalance(address),
    enabled: !!address,
  })
}

export const useTransactions = (address: string) => {
  return useQuery({
    queryKey: ['transactions', address],
    queryFn: () => walletApi.getTransactions(address),
    enabled: !!address,
  })
}