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
  Members?: Member[]
  Contributions?: Contribution[]
}

export interface Member {
  ID: string
  GroupID: string
  Wallet: string
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