import { Link } from 'react-router-dom';
import { Users } from 'lucide-react';
import { useGroupBalance } from '../hooks/useGroups';
import type { Group } from '../types';

interface GroupCardProps {
  group: Group;
}

const GroupCard = ({ group }: GroupCardProps) => {
  const { data: balance, isLoading: balanceLoading } = useGroupBalance(group.ID);

  const groupBalance = balance?.data?.balance || '0';

  return (
    <Link
      key={group.ID}
      to={`/groups/${group.ID}`}
      className="card hover:shadow-md transition-shadow"
    >
      <div className="flex items-center justify-between mb-4">
        <div className="w-12 h-12 bg-gradient-to-r from-stellar-500 to-primary-600 rounded-lg flex items-center justify-center">
          <Users className="w-6 h-6 text-white" />
        </div>
        <span className="text-sm text-gray-500">
          {group.Members?.length || 0} members
        </span>
      </div>
      
      <h3 className="text-xl font-semibold text-gray-900 mb-2">
        {group.Name}
      </h3>
      
      <p className="text-gray-600 mb-4 line-clamp-2">
        {group.Description}
      </p>
      
      <div className="flex items-center justify-between text-sm">
        <span className="text-gray-500">
          Total Savings: {balanceLoading ? '...' : `${groupBalance} XLM`}
        </span>
        <span className="text-primary-600 font-medium">
          View Details →
        </span>
      </div>
    </Link>
  );
};

export default GroupCard;
