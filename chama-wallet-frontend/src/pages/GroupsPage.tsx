import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useGroups } from '../hooks/useGroups'
import { Users, Plus, Search, Filter } from 'lucide-react'
import LoadingSpinner from '../components/LoadingSpinner'
import GroupCard from '../components/GroupCard'

const GroupsPage = () => {
  const { data: groups, isLoading } = useGroups()
  const [searchTerm, setSearchTerm] = useState('')
  const [filterType, setFilterType] = useState('all')

  // Add debug logs
  console.log('Groups data:', groups)
  console.log('Groups type:', typeof groups)
  console.log('Is array:', Array.isArray(groups))

  const filteredGroups = groups?.filter(group => {
    const matchesSearch = group.Name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         group.Description.toLowerCase().includes(searchTerm.toLowerCase())
    
    // Add more filter logic here based on filterType
    return matchesSearch
  }) || []

  return (
  <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold text-[#1a237e]">Savings Groups</h1>
          <p className="text-gray-600 mt-2">
            Discover and join savings groups in your community
          </p>
        </div>
        <Link to="/groups/create" className="inline-flex items-center px-6 py-3 rounded-lg border border-[#2ecc71] bg-[#2ecc71] text-white font-semibold shadow hover:bg-[#27ae60] transition mt-4 sm:mt-0">
          <Plus className="w-4 h-4 mr-2" />
          Create Group
        </Link>
      </div>

      {/* Search and Filter */}
      <div className="flex flex-col sm:flex-row gap-4 mb-8">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-[#2ecc71] w-5 h-5" />
          <input
            type="text"
            placeholder="Search groups..."
            className="w-full pl-10 px-4 py-3 rounded-lg border border-gray-200 focus:border-[#2ecc71] focus:ring-2 focus:ring-[#2ecc71] bg-white/90 text-gray-900 placeholder-gray-400 shadow"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
        <div className="relative">
          <Filter className="absolute left-3 top-1/2 transform -translate-y-1/2 text-[#2ecc71] w-5 h-5" />
          <select
            className="w-full pl-10 pr-8 px-4 py-3 rounded-lg border border-gray-200 focus:border-[#2ecc71] focus:ring-2 focus:ring-[#2ecc71] bg-white/90 text-gray-900 shadow"
            value={filterType}
            onChange={(e) => setFilterType(e.target.value)}
          >
            <option value="all">All Groups</option>
            <option value="active">Active</option>
            <option value="new">New</option>
          </select>
        </div>
      </div>

      {/* Groups Grid */}
      {isLoading ? (
        <div className="flex justify-center py-12">
          <LoadingSpinner size="lg" />
        </div>
      ) : filteredGroups.length === 0 ? (
        <div className="text-center py-12">
          <Users className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-xl font-semibold text-gray-900 mb-2">
            {searchTerm ? 'No groups found' : 'No groups available'}
          </h3>
          <p className="text-gray-600 mb-6">
            {searchTerm 
              ? 'Try adjusting your search terms'
              : 'Be the first to create a savings group!'
            }
          </p>
          <Link to="/groups/create" className="btn btn-primary">
            <Plus className="w-4 h-4 mr-2" />
            Create First Group
          </Link>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredGroups.map((group) => (
            <div className="bg-white/80 backdrop-blur-lg rounded-2xl shadow-xl p-6 border border-gray-100">
              <GroupCard key={group.ID} group={group} />
            </div>
          ))}
        </div>
      )}

      {/* Pagination could be added here */}
      {filteredGroups.length > 0 && (
        <div className="mt-8 text-center">
          <p className="text-gray-600">
            Showing {filteredGroups.length} of {groups?.length || 0} groups
          </p>
        </div>
      )}
    </div>
  )
}

export default GroupsPage
