
import { Outlet } from 'react-router-dom';
import DashboardLayout from './DashboardLayout';


const Layout = () => {
  return (
    <DashboardLayout>
      <Outlet />
    </DashboardLayout>
  );
}

export default Layout