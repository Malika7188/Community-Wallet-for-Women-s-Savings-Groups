import axios from 'axios'
import { 
  Group, 
  WalletBalance, 
  Transaction, 
  TransferRequest, 
  CreateGroupRequest, 
  JoinGroupRequest, 
  ContributeRequest 
} from '../types'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Wallet APIs
export const walletApi = {
  createWallet: () => api.post('/create-wallet'),
  getBalance: (address: string) => api.get<WalletBalance>(`/balance/${address}`),
  transferFunds: (data: TransferRequest) => api.post('/transfer', data),
  generateKeypair: () => api.get('/generate-keypair'),
  fundAccount: (address: string) => api.post(`/fund/${address}`),
  getTransactions: (address: string) => api.get<{ transactions: Transaction[] }>(`/transactions/${address}`),
}

// Group APIs
export const groupApi = {
  createGroup: (data: CreateGroupRequest) => api.post<Group>('/group/create', data),
  getAllGroups: () => api.get<Group[]>('/groups'),
  getGroupBalance: (id: string) => api.get(`/group/${id}/balance`),
  joinGroup: (id: string, data: JoinGroupRequest) => api.post(`/group/${id}/join`, data),
  contributeToGroup: (id: string, data: ContributeRequest) => api.post(`/group/${id}/contribute`, data),
}

export default api