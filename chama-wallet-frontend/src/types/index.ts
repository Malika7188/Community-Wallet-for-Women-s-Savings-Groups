export interface User {
  id: string
  email: string
  name: string
  wallet: string
  createdAt: string
}

export interface Group {
  ID: string
  Name: string
  Description: string
  Wallet: string
  CreatorID: string
  Creator: User
  Members?: Member[]
  Contributions?: Contribution[]
  Status: 'pending' | 'active' | 'completed'
  ContributionAmount?: number
  ContributionPeriod?: number
  PayoutOrder?: string
  CurrentRound?: number
  MaxMembers?: number
  MinMembers?: number
  NextContributionDate?: string
  IsApproved?: boolean
  CreatedAt: string
  UpdatedAt: string
}

export interface Member {
  ID: string
  GroupID: string
  UserID: string
  User: User
  Wallet: string
  Role: 'member' | 'admin' | 'creator'
  JoinedAt: string
  Status: 'pending' | 'approved' | 'rejected'
}

export interface Contribution {
  ID: string
  GroupID: string
  MemberID: string
  Amount: number
}

export interface WalletBalance {
  balances: string[]
}

export interface Transaction {
  hash: string
  ledger: number
  memo: string
  successful: boolean
  created_at: string
  fee_charged: string
}

export interface TransferRequest {
  from_seed: string
  to_address: string
  amount: string
}

export interface CreateGroupRequest {
  name: string
  description: string
}

export interface JoinGroupRequest {
  wallet: string
}

export interface ContributeRequest {
  from: string
  secret: string
  amount: string
}

export interface GroupInvitation {
  ID: string
  GroupID: string
  Group: Group
  InviterID: string
  Inviter: User
  Email: string
  UserID?: string
  User?: User
  Status: 'pending' | 'accepted' | 'rejected'
  CreatedAt: string
  ExpiresAt: string
}

export interface Notification {
  ID: string
  UserID: string
  GroupID: string
  Group: Group
  Type: string
  Title: string
  Message: string
  Read: boolean
  CreatedAt: string
}

export interface InviteUserRequest {
  email: string
}

export interface ActivateGroupRequest {
  contribution_amount: number
  contribution_period: number
  payout_order: string[]
}

export interface PayoutRequest {
  ID: string
  GroupID: string
  Group: Group
  RecipientID: string
  Recipient: User
  Amount: number
  Round: number
  Status: 'pending' | 'approved' | 'rejected' | 'completed'
  Approvals?: PayoutApproval[]
  CreatedAt: string
}

export interface PayoutApproval {
  ID: string
  PayoutRequestID: string
  AdminID: string
  Admin: User
  Approved: boolean
  CreatedAt: string
}

export interface PayoutSchedule {
  ID: string
  GroupID: string
  MemberID: string
  Member: Member
  Round: number
  Amount: number
  DueDate: string
  Status: 'scheduled' | 'paid' | 'pending'
  PaidAt?: string
  TxHash?: string
  CreatedAt: string
  UpdatedAt: string
}

export interface RoundContribution {
  ID: string
  GroupID: string
  MemberID: string
  Member: Member
  Round: number
  Amount: number
  Status: 'pending' | 'confirmed' | 'failed'
  TxHash: string
  CreatedAt: string
  UpdatedAt: string
}

export interface RoundStatus {
  ID: string
  GroupID: string
  Round: number
  TotalRequired: number
  TotalReceived: number
  ContributorsCount: number
  RequiredCount: number
  Status: 'collecting' | 'ready_for_payout' | 'completed'
  PayoutAuthorized: boolean
  CreatedAt: string
  UpdatedAt: string
}

export interface MemberContributionStatus {
  member: Member
  has_paid: boolean
  contribution?: RoundContribution
}

export interface RoundStatusResponse {
  round: number
  round_status: RoundStatus
  member_status: MemberContributionStatus[]
  total_members: number
  paid_members: number
}
