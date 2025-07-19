import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useGroupMutations } from '../hooks/useGroups'
import { Users, ArrowLeft } from 'lucide-react'
import { Link } from 'react-router-dom'

const CreateGroupPage = () => {
  const navigate = useNavigate()
  const { createGroup } = useGroupMutations()
  const [formData, setFormData] = useState({
    name: '',
    description: '',
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    try {
      await createGroup.mutateAsync(formData)
      navigate('/groups')
    } catch (error) {
      // Error is handled by the mutation
    }
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setFormData(prev => ({
      ...prev,
      [e.target.name]: e.target.value
    }))
  }

  return (
    <div className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <Link 
          to="/groups" 
          className="inline-flex items-center text-gray-600 hover:text-gray-900 mb-4"
        >
          <ArrowLeft className="w-4 h-4 mr-2" />
          Back to Groups
        </Link>
        
        <div className="text-center">
          <div className="w-16 h-16 bg-gradient-to-r from-stellar-500 to-primary-600 rounded-2xl flex items-center justify-center mx-auto mb-4">
            <Users className="w-8 h-8 text-white" />
          </div>
          <h1 className="text-3xl font-bold text-gray-900">Create New Group</h1>
          <p className="text-gray-600 mt-2">
            Start a new savings group and invite your community
          </p>
        </div>
      </div>

      <div className="card">
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-2">
              Group Name *
            </label>
            <input
              type="text"
              id="name"
              name="name"
              required
              className="input"
              placeholder="Enter group name (e.g., Alpha Chama)"
              value={formData.name}
              onChange={handleChange}
            />
            <p className="text-sm text-gray-500 mt-1">
              Choose a memorable name for your savings group
            </p>
          </div>

          <div>
            <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-2">
              Description *
            </label>
            <textarea
              id="description"
              name="description"
              required
              rows={4}
              className="input resize-none"
              placeholder="Describe the purpose and goals of your group..."
              value={formData.description}
              onChange={handleChange}
            />
            <p className="text-sm text-gray-500 mt-1">
              Explain what your group is about and its savings goals
            </p>
          </div>

          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <h3 className="font-medium text-blue-900 mb-2">What happens next?</h3>
            <ul className="text-sm text-blue-800 space-y-1">
              <li>• A unique Stellar wallet will be created for your group</li>
              <li>• You'll become the group administrator</li>
              <li>• You can invite members to join your group</li>
              <li>• Members can start making contributions immediately</li>
            </ul>
          </div>

          <div className="flex flex-col sm:flex-row gap-4">
            <Link 
              to="/groups" 
              className="btn btn-secondary flex-1"
            >
              Cancel
            </Link>
            <button
              type="submit"
              disabled={createGroup.isPending}
              className="btn btn-primary flex-1"
            >
              {createGroup.isPending ? 'Creating Group...' : 'Create Group'}
            </button>
          </div>
        </form>
      </div>

      {/* Additional Info */}
      <div className="mt-8 text-center">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">
          Why Create a Savings Group?
        </h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="text-center">
            <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center mx-auto mb-2">
              <Users className="w-6 h-6 text-green-600" />
            </div>
            <h4 className="font-medium text-gray-900">Community Support</h4>
            <p className="text-sm text-gray-600">Save together with friends and family</p>
          </div>
          <div className="text-center">
            <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center mx-auto mb-2">
              <Users className="w-6 h-6 text-blue-600" />
            </div>
            <h4 className="font-medium text-gray-900">Transparency</h4>
            <p className="text-sm text-gray-600">All transactions are recorded on blockchain</p>
          </div>
          <div className="text-center">
            <div className="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center mx-auto mb-2">
              <Users className="w-6 h-6 text-purple-600" />
            </div>
            <h4 className="font-medium text-gray-900">Goal Achievement</h4>
            <p className="text-sm text-gray-600">Reach your savings goals faster together</p>
          </div>
        </div>
      </div>
    </div>
  )
}

export default CreateGroupPage