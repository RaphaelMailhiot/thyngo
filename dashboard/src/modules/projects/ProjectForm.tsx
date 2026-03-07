import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import api from '../../services/api';

const ProjectForm: React.FC = () => {
    const { slug } = useParams<{ slug: string }>();
    const navigate = useNavigate();
    const isEditing = !!slug;

    const [title, setTitle] = useState('');
    const [newSlug, setNewSlug] = useState('');
    const [visibility, setVisibility] = useState('public');
    const [loading, setLoading] = useState(isEditing);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState('');

    useEffect(() => {
        if (isEditing) {
            const fetchProject = async () => {
                try {
                    const response = await api.get(`/projects/${slug}`);
                    if (response.data && response.data.data) {
                        const project = response.data.data;
                        setTitle(project.title);
                        setNewSlug(project.slug);
                        setVisibility(project.visibility);
                    }
                } catch (err: any) {
                    setError('Failed to load project data.');
                } finally {
                    setLoading(false);
                }
            };
            fetchProject();
        }
    }, [slug, isEditing]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSaving(true);
        setError('');

        try {
            const payload = {
                title,
                slug: newSlug,
                visibility
            };

            if (isEditing) {
                await api.put(`/projects/${slug}`, payload);
            } else {
                await api.post('/projects', payload);
            }
            navigate('/projects');
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to save project.');
            setSaving(false);
        }
    };

    return (
        <div>
            <div className="d-flex justify-content-between align-items-center mb-4 border-bottom border-secondary border-opacity-25 pb-3">
                <h2 style={{ fontFamily: 'var(--font-heading)' }}>
                    <i className="bi bi-stack me-2 text-primary"></i>
                    {isEditing ? 'Edit Project' : 'Create New Project'}
                </h2>
                <Link to="/projects" className="btn btn-outline-secondary border-0">
                    <i className="bi bi-x-lg me-2"></i>Cancel
                </Link>
            </div>

            {loading ? (
                <div className="text-center mt-5"><div className="spinner-border text-primary" /></div>
            ) : (
                <div className="card glass-panel border-0 mx-auto" style={{ maxWidth: '800px' }}>
                    <div className="card-body p-4">
                        {error && <div className="alert alert-danger bg-danger bg-opacity-10 border-0 text-danger mb-4">{error}</div>}

                        <form onSubmit={handleSubmit}>
                            <div className="mb-4">
                                <label className="form-label text-secondary fw-semibold">Title</label>
                                <input
                                    type="text"
                                    className="form-control form-control-lg bg-transparent"
                                    value={title}
                                    onChange={(e) => setTitle(e.target.value)}
                                    required
                                    placeholder="Project Title"
                                />
                            </div>

                            <div className="row mb-4">
                                <div className="col-md-6">
                                    <label className="form-label text-secondary fw-semibold">Slug</label>
                                    <input
                                        type="text"
                                        className="form-control bg-transparent"
                                        value={newSlug}
                                        onChange={(e) => setNewSlug(e.target.value)}
                                        required
                                        disabled={isEditing}
                                        placeholder="e.g. awesome-app"
                                    />
                                    {isEditing && <div className="form-text mt-2 opacity-50">Slugs cannot be changed after creation.</div>}
                                </div>

                                <div className="col-md-6">
                                    <label className="form-label text-secondary fw-semibold">Visibility</label>
                                    <select
                                        className="form-select bg-transparent"
                                        value={visibility}
                                        onChange={(e) => setVisibility(e.target.value)}
                                    >
                                        <option value="public" className="bg-dark text-white">Public (Visible to everyone)</option>
                                        <option value="private" className="bg-dark text-white">Private (Only visible to you)</option>
                                        <option value="shared" className="bg-dark text-white">Shared (Visible to specific users)</option>
                                    </select>
                                </div>
                            </div>

                            <div className="d-flex justify-content-end mt-5 pt-3 border-top border-secondary border-opacity-10">
                                <button type="submit" className="btn btn-primary px-5 py-2 shadow-sm hover-lift" disabled={saving}>
                                    {saving ? (
                                        <><span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true" /> Saving...</>
                                    ) : (
                                        <><i className="bi bi-check-circle me-2"></i>{isEditing ? 'Update Project' : 'Create Project'}</>
                                    )}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default ProjectForm;
