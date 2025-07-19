# Chama Wallet Frontend

A modern React frontend for the Chama Wallet platform - a blockchain-powered savings and lending platform for Chamas (informal savings groups) built on the Stellar network.

## 🌟 Features

### Core Functionality
- **User Authentication**: Sign up, login, and logout functionality
- **Wallet Management**: Create wallets, view balances, send/receive XLM
- **Group Management**: Create, join, and manage savings groups
- **Contributions**: Make contributions to group savings
- **Transaction History**: View detailed transaction history
- **Stellar Integration**: Full integration with Stellar testnet

### User Interface
- **Responsive Design**: Works seamlessly on desktop, tablet, and mobile
- **Modern UI**: Clean, intuitive interface with Tailwind CSS
- **Real-time Updates**: Live balance and transaction updates
- **Loading States**: Smooth loading indicators and error handling
- **Toast Notifications**: User-friendly success/error messages

## 🚀 Getting Started

### Prerequisites
- Node.js 18+ 
- npm or yarn
- Running Chama Wallet Backend (see backend documentation)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd chama-wallet-frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Environment Setup**
   ```bash
   cp .env.example .env
   ```
   
   Update `.env` with your configuration:
   ```env
   VITE_API_BASE_URL=http://localhost:3000
   VITE_APP_NAME=Chama Wallet
   VITE_STELLAR_NETWORK=testnet
   ```

4. **Start the development server**
   ```bash
   npm run dev
   ```

5. **Open your browser**
   Navigate to `http://localhost:5173`

## 📁 Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── Layout.tsx      # Main layout wrapper
│   ├── Navbar.tsx      # Navigation component
│   ├── ProtectedRoute.tsx # Route protection
│   └── LoadingSpinner.tsx # Loading indicator
├── contexts/           # React contexts
│   └── AuthContext.tsx # Authentication context
├── hooks/              # Custom React hooks
│   ├── useWallet.ts    # Wallet operations
│   └── useGroups.ts    # Group operations
├── pages/              # Page components
│   ├── HomePage.tsx    # Landing page
│   ├── LoginPage.tsx   # Login form
│   ├── SignUpPage.tsx  # Registration form
│   ├── DashboardPage.tsx # User dashboard
│   ├── GroupsPage.tsx  # Groups listing
│   ├── GroupDetailPage.tsx # Group details
│   ├── CreateGroupPage.tsx # Create group form
│   ├── WalletPage.tsx  # Wallet management
│   └── TransactionsPage.tsx # Transaction history
├── services/           # API services
│   └── api.ts          # API client and endpoints
├── types/              # TypeScript type definitions
│   └── index.ts        # Shared types
├── App.tsx             # Main app component
├── main.tsx            # App entry point
└── index.css           # Global styles
```

## 🔧 API Integration

The frontend integrates with the following backend endpoints:

### Wallet Endpoints
- `POST /create-wallet` - Create new wallet
- `GET /balance/:address` - Get wallet balance
- `POST /transfer` - Transfer funds
- `GET /generate-keypair` - Generate new keypair
- `POST /fund/:address` - Fund account (testnet)
- `GET /transactions/:address` - Get transaction history

### Group Endpoints
- `POST /group/create` - Create new group
- `GET /groups` - Get all groups
- `GET /group/:id/balance` - Get group balance
- `POST /group/:id/join` - Join group
- `POST /group/:id/contribute` - Contribute to group

## 🎨 UI Components

### Design System
- **Colors**: Primary (blue), Stellar (cyan), semantic colors
- **Typography**: Inter font family with consistent sizing
- **Spacing**: 8px grid system
- **Components**: Reusable button, input, and card components

### Key Components
- **Navbar**: Responsive navigation with user menu
- **Cards**: Consistent card layout for content sections
- **Buttons**: Primary, secondary, and outline variants
- **Forms**: Styled form inputs with validation
- **Modals**: Overlay modals for actions like transfers

## 🔐 Authentication

The app uses a context-based authentication system:

- **AuthContext**: Manages user state and authentication methods
- **ProtectedRoute**: Wraps protected pages to ensure authentication
- **Local Storage**: Persists user session (development only)

### Authentication Flow
1. User signs up/logs in
2. User data stored in context and localStorage
3. Protected routes check authentication status
4. Automatic redirect to login if not authenticated

## 💰 Wallet Features

### Wallet Management
- **Balance Display**: Real-time XLM balance
- **Address Display**: Copy wallet address to clipboard
- **Fund Account**: Get testnet XLM from Friendbot
- **Generate Keypair**: Create new Stellar keypairs

### Transactions
- **Send XLM**: Transfer funds to other wallets
- **Transaction History**: View past transactions
- **Stellar Explorer**: Links to view transactions on Stellar Explorer

## 👥 Group Features

### Group Management
- **Create Groups**: Start new savings groups
- **Join Groups**: Join existing groups
- **Group Dashboard**: View group details and members
- **Group Wallet**: Each group has its own Stellar wallet

### Contributions
- **Make Contributions**: Send XLM to group wallet
- **View Balance**: See total group savings
- **Member List**: View all group members

## 🛠️ Development

### Available Scripts
- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

### Code Quality
- **TypeScript**: Full type safety
- **ESLint**: Code linting and formatting
- **Prettier**: Code formatting (recommended)

### State Management
- **React Query**: Server state management and caching
- **React Context**: Client state management
- **Local Storage**: Persistence for user sessions

## 🌐 Deployment

### Build for Production
```bash
npm run build
```

### Environment Variables
Ensure all environment variables are set for production:
- `VITE_API_BASE_URL`: Backend API URL
- `VITE_APP_NAME`: Application name
- `VITE_STELLAR_NETWORK`: Stellar network (testnet/mainnet)

### Deployment Platforms
The app can be deployed to:
- Vercel
- Netlify
- AWS S3 + CloudFront
- Any static hosting service

## 🔒 Security Considerations

### Development vs Production
- **Secret Keys**: Never store secret keys in frontend code
- **Environment Variables**: Use secure environment variable management
- **HTTPS**: Always use HTTPS in production
- **API Security**: Implement proper API authentication

### Best Practices
- Secret keys are only used for transaction signing
- All sensitive operations require user input
- Clear security warnings and tips throughout the UI
- Testnet usage for development and testing

## 🐛 Troubleshooting

### Common Issues

1. **Backend Connection**
   - Ensure backend is running on correct port
   - Check CORS configuration
   - Verify API_BASE_URL in .env

2. **Stellar Network**
   - Confirm using testnet for development
   - Check wallet addresses are valid
   - Ensure accounts are funded

3. **Build Issues**
   - Clear node_modules and reinstall
   - Check Node.js version compatibility
   - Verify all environment variables

## 📚 Additional Resources

- [Stellar Documentation](https://developers.stellar.org/)
- [React Documentation](https://react.dev/)
- [Tailwind CSS](https://tailwindcss.com/)
- [React Query](https://tanstack.com/query/latest)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Check the troubleshooting section
- Review the backend API documentation

---

Built with ❤️ for the Chama community using React, TypeScript, and Stellar blockchain technology.