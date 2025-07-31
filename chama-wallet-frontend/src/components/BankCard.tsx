import React from 'react';

interface BankCardProps {
  title: string;
  value: string | number;
  icon: React.ReactNode;
  progress?: number;
  colorClass?: string;
}

const BankCard: React.FC<BankCardProps> = ({ title, value, icon, progress, colorClass }) => {
  return (
    <div
      className="relative rounded-3xl p-10 flex flex-col gap-4 min-h-[180px] min-w-[340px] max-w-[420px] w-full bg-white/80 backdrop-blur-lg shadow-2xl border border-gray-100 transition-transform duration-300 ease-out hover:-translate-y-3 hover:shadow-3xl group overflow-hidden"
      style={{
        boxShadow: '0 12px 40px 0 rgba(31, 38, 135, 0.18)',
        fontFamily: 'Inter, Roboto, sans-serif',
      }}
    >
      <div className="absolute inset-0 pointer-events-none opacity-50 group-hover:opacity-70 transition-opacity duration-300" style={{background: 'linear-gradient(135deg, #e0f7fa 0%, #e8f5e9 100%)'}} />
      <div className="relative flex items-center gap-6 z-10">
        <div className={`w-16 h-16 rounded-2xl flex items-center justify-center text-white text-2xl shadow-lg ${colorClass || 'bg-[#1a237e]'} transition-transform duration-300 group-hover:scale-110`}>{icon}</div>
        <div>
          <div className="text-lg text-gray-500 font-semibold tracking-wide" style={{ fontFamily: 'Inter, Roboto, sans-serif' }}>{title}</div>
          <div className="text-4xl font-black text-[#1a237e] tracking-tight drop-shadow-sm" style={{ fontFamily: 'Inter, Roboto, sans-serif' }}>{value}</div>
        </div>
      </div>
      {typeof progress === 'number' && (
        <div className="mt-6 relative z-10">
          <div className="w-full bg-gray-200 rounded-full h-3">
            <div className="bg-[#2ecc71] h-3 rounded-full transition-all duration-500" style={{ width: `${progress}%` }}></div>
          </div>
        </div>
      )}
    </div>
  );
};

export default BankCard;
