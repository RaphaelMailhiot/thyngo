import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import api from '../../services/api';

const ResumeForm: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const isEditing = !!id;

    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [phone, setPhone] = useState('');
    const [job, setJob] = useState('');
    const [visibility, setVisibility] = useState('public');
    const [loading, setLoading] = useState(isEditing);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState('');

    useEffect(() => {
        if (isEditing) {
            const fetchResume = async () => {
                try {
                    const response = await api.get(`/resumes/${id}`);
                    if (response.data && response.data.data) {
                        const r = response.data.data;
                        setName(r.name);
                        setEmail(r.email);
                        setPhone(r.phone);
                        setJob(r.job);
                        setVisibility(r.visibility);
                    }
                } catch (err: any) {
                    setError('Failed to load resume data.');
                } finally {
                    setLoading(false);
                }
            };
            fetchResume();
        }
    }, [id, isEditing]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSaving(true);
        setError('');

        try {
            const payload = {
                name,
                email,
                phone,
                job,
                visibility
            };

            if (isEditing) {
                await api.put(`/resumes/${id}`, payload);
            } else {
                await api.post('/resumes', payload);
            }
            navigate('/resumes');
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to save resume.');
            setSaving(false);
        }
    };

    return (
        <div>
            <div className="d-flex justify-content-between align-items-center mb-4 border-bottom border-secondary border-opacity-25 pb-3">
                <h2 style={{ fontFamily: 'var(--font-heading)' }}>
                    <i className="bi bi-file-earmark-person me-2 text-primary"></i>
                    {isEditing ? 'Edit Resume' : 'Create New Resume'}
                </h2>
                <Link to="/resumes" className="btn btn-outline-secondary border-0">
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
                                <label className="form-label text-secondary fw-semibold">Full Name</label>
                                <input
                                    type="text"
                                    className="form-control form-control-lg bg-transparent"
                                    value={name}
                                    onChange={(e) => setName(e.target.value)}
                                    required
                                    placeholder="e.g. John Doe"
                                />
                            </div>

                            <div className="row">
                                <div className="col-md-6 mb-4">
                                    <label className="form-label text-secondary fw-semibold">Email</label>
                                    <input
                                        type="email"
                                        className="form-control bg-transparent"
                                        value={email}
                                        onChange={(e) => setEmail(e.target.value)}
                                        placeholder="john@example.com"
                                    />
                                </div>
                                <div className="col-md-6 mb-4">
                                    <label className="form-label text-secondary fw-semibold">Phone</label>
                                    <input
                                        type="text"
                                        className="form-control bg-transparent"
                                        value={phone}
                                        onChange={(e) => setPhone(e.target.value)}
                                        placeholder="+1 234 567 890"
                                    />
                                </div>
                            </div>

                            <div className="mb-4">
                                <label className="form-label text-secondary fw-semibold">Job Title</label>
                                <input
                                    type="text"
                                    className="form-control bg-transparent"
                                    value={job}
                                    onChange={(e) => setJob(e.target.value)}
                                    placeholder="e.g. Senior Software Engineer"
                                />
                            </div>

                            <div className="mb-4">
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

                            <div className="d-flex justify-content-end mt-5 pt-3 border-top border-secondary border-opacity-10">
                                <button type="submit" className="btn btn-primary px-5 py-2 shadow-sm hover-lift" disabled={saving}>
                                    {saving ? (
                                        <><span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true" /> Saving...</>
                                    ) : (
                                        <><i className="bi bi-check-circle me-2"></i>{isEditing ? 'Update Resume' : 'Create Resume'}</>
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

export default ResumeForm;
