import React from 'react';
import { Outlet, Navigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import Sidebar from '../components/Sidebar';

const Layout: React.FC = () => {
    const { isAuthenticated } = useAuth();

    if (!isAuthenticated) {
        return <Navigate to="/login" replace />;
    }

    return (
        <div className="d-flex" style={{ minHeight: '100vh' }}>
            <Sidebar />
            <div className="flex-grow-1" style={{ marginLeft: '280px', padding: '2rem' }}>
                <div className="container-fluid glass-panel p-4" style={{ minHeight: 'calc(100vh - 4rem)' }}>
                    <Outlet />
                </div>
            </div>
        </div>
    );
};

export default Layout;
