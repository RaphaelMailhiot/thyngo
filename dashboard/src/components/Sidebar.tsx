import React from 'react';
import { NavLink } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const Sidebar: React.FC = () => {
    const { logout } = useAuth();

    return (
        <div className="d-flex flex-column flex-shrink-0 p-3 glass-sidebar shadow-lg" style={{ width: '280px', height: '100vh', position: 'fixed', left: 0, top: 0, zIndex: 1000 }}>
            <a href="/" className="d-flex align-items-center mb-3 mb-md-0 me-md-auto text-decoration-none" style={{ color: 'var(--text-primary)' }}>
                <i className="bi bi-layers-half fs-3 me-2" style={{ color: 'var(--accent-primary)' }}></i>
                <span className="fs-4 fw-bold" style={{ fontFamily: 'var(--font-heading)' }}>ThynGo Admin</span>
            </a>
            <hr style={{ borderColor: 'var(--border-color)' }} />
            <ul className="nav nav-pills flex-column mb-auto gap-2">
                <li className="nav-item">
                    <NavLink to="/dashboard" className={({ isActive }) => `nav-link px-3 py-2 ${isActive ? 'active' : ''}`}>
                        <i className="bi bi-speedometer2 me-3"></i>
                        Dashboard
                    </NavLink>
                </li>
                <li>
                    <NavLink to="/posts" className={({ isActive }) => `nav-link px-3 py-2 ${isActive ? 'active' : ''}`}>
                        <i className="bi bi-file-earmark-richtext me-3"></i>
                        Posts
                    </NavLink>
                </li>
                <li>
                    <NavLink to="/projects" className={({ isActive }) => `nav-link px-3 py-2 ${isActive ? 'active' : ''}`}>
                        <i className="bi bi-kanban me-3"></i>
                        Projects
                    </NavLink>
                </li>
                <li>
                    <NavLink to="/media" className={({ isActive }) => `nav-link px-3 py-2 ${isActive ? 'active' : ''}`}>
                        <i className="bi bi-images me-3"></i>
                        Media
                    </NavLink>
                </li>
                <li>
                    <NavLink to="/resumes" className={({ isActive }) => `nav-link px-3 py-2 ${isActive ? 'active' : ''}`}>
                        <i className="bi bi-person-badge me-3"></i>
                        Resumes
                    </NavLink>
                </li>
            </ul>
            <hr style={{ borderColor: 'var(--border-color)' }} />
            <div className="dropdown mt-2">
                <button className="btn btn-outline-danger w-100 hover-lift" onClick={logout} style={{ backgroundColor: 'rgba(239, 68, 68, 0.1)', borderColor: 'rgba(239, 68, 68, 0.3)' }}>
                    <i className="bi bi-box-arrow-left me-2"></i> Sign out
                </button>
            </div>
        </div>
    );
};

export default Sidebar;
