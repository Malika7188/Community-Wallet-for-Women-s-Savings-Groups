import axios from 'axios'
import type { 
  Group, 
  WalletBalance, 
  Transaction, 
  TransferRequest, 
  CreateGroupRequest, 
  JoinGroupRequest, 
  ContributeRequest,
  InviteUserRequest,
  ActivateGroupRequest,
  User,
  Notification,
  GroupInvitation,
  RoundStatusResponse
} from '../types'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add request interceptor to include auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Add response interceptor to handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Clear auth data and redirect to login
      localStorage.removeItem('user')
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// Auth APIs
export const authApi = {
  register: (data: { name: string; email: string; password: string }) => 
    api.post('/auth/register', data),
  login: (data: { email: string; password: string }) => 
    api.post('/auth/login', data),
  logout: () => api.post('/auth/logout'),
  getProfile: () => api.get('/auth/profile'),
  updateProfile: (data: { name: string }) => api.put('/auth/profile', data),
}

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
  inviteToGroup: (id: string, data: InviteUserRequest) => api.post(`/group/${id}/invite`, data),
  getNonGroupMembers: (id: string) => api.get<User[]>(`/group/${id}/non-members`),
  approveGroup: (id: string) => api.post(`/group/${id}/approve`),
  activateGroup: (id: string, data: ActivateGroupRequest) => api.post(`/group/${id}/activate`, data),
  nominateAdmin: (id: string, data: { nominee_id: string }) => api.post(`/group/${id}/nominate-admin`, data),
  approveMember: (id: string, data: { member_id: string, action: string }) => api.post(`/group/${id}/approve-member`, data),
  createPayoutRequest: (id: string, data: any) => api.post(`/group/${id}/payout-request`, data),
  getPayoutRequests: (id: string) => api.get(`/group/${id}/payout-requests`),
  getPayoutSchedule: (id: string) => api.get(`/group/${id}/payout-schedule`),
  contributeToRound: (id: string, data: { round: number, amount: number, secret: string }) => 
    api.post(`/group/${id}/contribute-round`, data),
  getRoundStatus: (id: string, round: number) => 
    api.get<RoundStatusResponse>(`/group/${id}/round-status?round=${round}`),
  authorizeRoundPayout: (id: string, data: { round: number }) => 
    api.post(`/group/${id}/authorize-payout`, data),
}

export const payoutApi = {
  approvePayout: (id: string, data: { approved: boolean }) => api.post(`/payout/${id}/approve`, data),
}

export const notificationApi = {
  getNotifications: () => api.get<Notification[]>('/notifications'),
  markAsRead: (id: string) => api.put(`/notifications/${id}/read`),
  getInvitations: () => api.get<GroupInvitation[]>('/invitations'),
  acceptInvitation: (id: string) => api.post(`/invitations/${id}/accept`),
  rejectInvitation: (id: string) => api.post(`/invitations/${id}/reject`),
}

export default api
