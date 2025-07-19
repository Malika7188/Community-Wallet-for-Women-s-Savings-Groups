import { Link } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { Users, Shield, Zap, Globe, ArrowRight, Star } from 'lucide-react'

const HomePage = () => {
  const { user } = useAuth()

  const features = [
    {
      icon: Users,
      title: 'Community Savings',
      description: 'Create and join savings groups with friends, family, or community members.',
    },
    {
      icon: Shield,
      title: 'Blockchain Security',
      description: 'Built on Stellar blockchain for transparent and secure transactions.',
    },
    {
      icon: Zap,
      title: 'Instant Transfers',
      description: 'Send and receive funds instantly with low transaction fees.',
    },
    {
      icon: Globe,
      title: 'Global Access',
      description: 'Access your savings from anywhere in the world, 24/7.',
    },
  ]

  const testimonials = [
    {
      name: 'Sarah Mwangi',
      role: 'Chama Leader',
      content: 'This platform has revolutionized how our savings group operates. Everything is transparent and secure.',
      rating: 5,
    },
    {
      name: 'John Kimani',
      role: 'Small Business Owner',
      content: 'The instant transfers and low fees make it perfect for our business transactions.',
      rating: 5,
    },
  ]

  return (
    <div className="min-h-screen">
      {/* Hero Section */}
      <section className="bg-gradient-to-br from-stellar-50 via-primary-50 to-white py-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <h1 className="text-4xl md:text-6xl font-bold text-gray-900 mb-6">
              Empower Your
              <span className="text-transparent bg-clip-text bg-gradient-to-r from-stellar-600 to-primary-600">
                {' '}Community Savings
              </span>
            </h1>
            <p className="text-xl text-gray-600 mb-8 max-w-3xl mx-auto">
              Join the future of community savings with our blockchain-powered platform. 
              Create savings groups, contribute securely, and grow your wealth together.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              {user ? (
                <Link to="/dashboard" className="btn btn-primary text-lg px-8 py-3">
                  Go to Dashboard
                  <ArrowRight className="w-5 h-5 ml-2" />
                </Link>
              ) : (
                <>
                  <Link to="/signup" className="btn btn-primary text-lg px-8 py-3">
                    Get Started
                    <ArrowRight className="w-5 h-5 ml-2" />
                  </Link>
                  <Link to="/login" className="btn btn-outline text-lg px-8 py-3">
                    Sign In
                  </Link>
                </>
              )}
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              Why Choose Chama Wallet?
            </h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              Built specifically for African savings groups with modern blockchain technology
            </p>
          </div>
          
          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
            {features.map((feature, index) => {
              const Icon = feature.icon
              return (
                <div key={index} className="text-center">
                  <div className="w-16 h-16 bg-gradient-to-r from-stellar-500 to-primary-600 rounded-2xl flex items-center justify-center mx-auto mb-4">
                    <Icon className="w-8 h-8 text-white" />
                  </div>
                  <h3 className="text-xl font-semibold text-gray-900 mb-2">
                    {feature.title}
                  </h3>
                  <p className="text-gray-600">
                    {feature.description}
                  </p>
                </div>
              )
            })}
          </div>
        </div>
      </section>

      {/* How It Works Section */}
      <section className="py-20 bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              How It Works
            </h2>
            <p className="text-xl text-gray-600">
              Get started in three simple steps
            </p>
          </div>
          
          <div className="grid md:grid-cols-3 gap-8">
            <div className="text-center">
              <div className="w-12 h-12 bg-primary-600 text-white rounded-full flex items-center justify-center mx-auto mb-4 text-xl font-bold">
                1
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-2">
                Create Your Wallet
              </h3>
              <p className="text-gray-600">
                Sign up and get your secure Stellar wallet instantly
              </p>
            </div>
            
            <div className="text-center">
              <div className="w-12 h-12 bg-primary-600 text-white rounded-full flex items-center justify-center mx-auto mb-4 text-xl font-bold">
                2
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-2">
                Join or Create Groups
              </h3>
              <p className="text-gray-600">
                Start a new savings group or join existing ones
              </p>
            </div>
            
            <div className="text-center">
              <div className="w-12 h-12 bg-primary-600 text-white rounded-full flex items-center justify-center mx-auto mb-4 text-xl font-bold">
                3
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-2">
                Start Saving Together
              </h3>
              <p className="text-gray-600">
                Make contributions and watch your savings grow
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Testimonials Section */}
      <section className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              What Our Users Say
            </h2>
          </div>
          
          <div className="grid md:grid-cols-2 gap-8">
            {testimonials.map((testimonial, index) => (
              <div key={index} className="card">
                <div className="flex items-center mb-4">
                  {[...Array(testimonial.rating)].map((_, i) => (
                    <Star key={i} className="w-5 h-5 text-yellow-400 fill-current" />
                  ))}
                </div>
                <p className="text-gray-600 mb-4 italic">
                  "{testimonial.content}"
                </p>
                <div>
                  <p className="font-semibold text-gray-900">{testimonial.name}</p>
                  <p className="text-sm text-gray-500">{testimonial.role}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 bg-gradient-to-r from-stellar-600 to-primary-600">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="text-3xl md:text-4xl font-bold text-white mb-4">
            Ready to Transform Your Savings?
          </h2>
          <p className="text-xl text-stellar-100 mb-8 max-w-2xl mx-auto">
            Join thousands of users who are already saving smarter with blockchain technology
          </p>
          {!user && (
            <Link to="/signup" className="btn bg-white text-stellar-600 hover:bg-gray-100 text-lg px-8 py-3">
              Start Your Journey
              <ArrowRight className="w-5 h-5 ml-2" />
            </Link>
          )}
        </div>
      </section>
    </div>
  )
}

export default HomePage