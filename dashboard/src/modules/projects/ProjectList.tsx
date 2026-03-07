import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import api from '../../services/api';

interface Project {
    id: number;
    slug: string;
    title: string;
    visibility: string;
    created_at: string;
}

const ProjectList: React.FC = () => {
    const [projects, setProjects] = useState<Project[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchProjects = async () => {
        try {
            const response = await api.get('/projects');
            if (response.data && response.data.data) {
                setProjects(response.data.data);
            }
        } catch (error) {
            console.error('Failed to fetch projects', error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchProjects();
    }, []);

    const handleDelete = async (slug: string) => {
        if (window.confirm('Are you sure you want to delete this project?')) {
            try {
                await api.delete(`/projects/${slug}`);
                fetchProjects(); // Refresh list
            } catch (error) {
                console.error('Failed to delete project', error);
                alert('Could not delete project.');
            }
        }
    };

    return (
        <div>
            <div className="d-flex justify-content-between align-items-center mb-4 border-bottom border-secondary border-opacity-25 pb-3">
                <h2 style={{ fontFamily: 'var(--font-heading)' }}>
                    <i className="bi bi-stack me-2 text-primary"></i>
                    Projects
                </h2>
                <Link to="/projects/new" className="btn btn-primary shadow-sm hover-lift">
                    <i className="bi bi-plus-lg me-2"></i>Create Project
                </Link>
            </div>

            {loading ? (
                <div className="text-center mt-5">
                    <div className="spinner-border text-primary" role="status">
                        <span className="visually-hidden">Loading...</span>
                    </div>
                </div>
            ) : (
                <div className="card glass-panel border-0">
                    <div className="table-responsive">
                        <table className="table table-hover mb-0 align-middle">
                            <thead>
                                <tr className="text-secondary" style={{ fontSize: '0.85rem', textTransform: 'uppercase', letterSpacing: '1px' }}>
                                    <th className="ps-4">ID</th>
                                    <th>Title</th>
                                    <th>Slug</th>
                                    <th>Visibility</th>
                                    <th>Created At</th>
                                    <th className="text-end pe-4">Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {projects.length === 0 ? (
                                    <tr>
                                        <td colSpan={6} className="text-center py-5 text-muted">
                                            <i className="bi bi-folder-x fs-1 d-block mb-2 opacity-25"></i>
                                            No projects found. Create your first one!
                                        </td>
                                    </tr>
                                ) : (
                                    projects.map(project => (
                                        <tr key={project.id}>
                                            <td className="ps-4 text-secondary">#{project.id}</td>
                                            <td className="fw-semibold">{project.title}</td>
                                            <td><code className="bg-dark text-info px-2 py-1 rounded small">{project.slug}</code></td>
                                            <td>
                                                <span className={`badge bg-${project.visibility === 'public' ? 'success' : project.visibility === 'private' ? 'secondary' : 'info'} bg-opacity-25 text-${project.visibility === 'public' ? 'success' : project.visibility === 'private' ? 'secondary' : 'info'} border border-${project.visibility === 'public' ? 'success' : project.visibility === 'private' ? 'secondary' : 'info'} border-opacity-25`}>
                                                    {project.visibility}
                                                </span>
                                            </td>
                                            <td className="text-secondary" style={{ fontSize: '0.9rem' }}>{new Date(project.created_at).toLocaleDateString()}</td>
                                            <td className="text-end pe-4">
                                                <div className="btn-group btn-group-sm">
                                                    <Link to={`/projects/edit/${project.slug}`} className="btn btn-outline-secondary border-0" title="Edit">
                                                        <i className="bi bi-pencil"></i>
                                                    </Link>
                                                    <button onClick={() => handleDelete(project.slug)} className="btn btn-outline-danger border-0" title="Delete">
                                                        <i className="bi bi-trash"></i>
                                                    </button>
                                                </div>
                                            </td>
                                        </tr>
                                    ))
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>
            )}
        </div>
    );
};

export default ProjectList;
