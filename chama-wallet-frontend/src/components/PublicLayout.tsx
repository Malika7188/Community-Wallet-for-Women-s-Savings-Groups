import { Outlet } from 'react-router-dom';

const PublicLayout = () => {
  return (
    <div className="min-h-screen bg-gradient-to-br from-[#e0e7ef] to-[#f5f6fa] flex flex-col justify-center items-center">
      <Outlet />
    </div>
  );
};

export default PublicLayout;
