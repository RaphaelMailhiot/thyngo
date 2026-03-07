import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import api from '../../services/api';

const MediaForm: React.FC = () => {
    const { slug } = useParams<{ slug: string }>();
    const navigate = useNavigate();
    const isEditing = !!slug;

    const [file, setFile] = useState<File | null>(null);
    const [title, setTitle] = useState('');
    const [newSlug, setNewSlug] = useState('');
    const [type, setType] = useState('image');
    const [visibility, setVisibility] = useState('public');
    const [loading, setLoading] = useState(isEditing);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState('');

    useEffect(() => {
        if (isEditing) {
            const fetchMedia = async () => {
                try {
                    const response = await api.get(`/media/${slug}`);
                    if (response.data && response.data.data) {
                        const m = response.data.data;
                        setTitle(m.title);
                        setNewSlug(m.slug);
                        setType(m.type);
                        setVisibility(m.visibility);
                    }
                } catch (err: any) {
                    setError('Failed to load media data.');
                } finally {
                    setLoading(false);
                }
            };
            fetchMedia();
        }
    }, [slug, isEditing]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSaving(true);
        setError('');

        try {
            if (isEditing) {
                // Media edit endpoint only accepts title and visibility
                await api.put(`/media/${slug}`, { title, visibility });
            } else {
                if (!file) {
                    setError('Please select a file to upload.');
                    setSaving(false);
                    return;
                }
                const formData = new FormData();
                formData.append('file', file);
                formData.append('slug', newSlug);
                formData.append('title', title);
                formData.append('type', type);
                formData.append('visibility', visibility);

                await api.post('/media', formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                });
            }
            navigate('/media');
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to save media.');
            setSaving(false);
        }
    };

    return (
        <div>
            <div className="d-flex justify-content-between align-items-center mb-4 border-bottom border-secondary border-opacity-25 pb-3">
                <h2 style={{ fontFamily: 'var(--font-heading)' }}>
                    <i className="bi bi-images me-2 text-primary"></i>
                    {isEditing ? 'Edit Media' : 'Upload New Media'}
                </h2>
                <Link to="/media" className="btn btn-outline-secondary border-0">
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
                                    placeholder="Media Title"
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
                                        placeholder="e.g. hero-banner"
                                    />
                                    {isEditing && <div className="form-text mt-2 opacity-50">Slugs cannot be changed after creation.</div>}
                                </div>

                                <div className="col-md-6">
                                    <label className="form-label text-secondary fw-semibold">Resource Type</label>
                                    <select
                                        className="form-select bg-transparent"
                                        value={type}
                                        onChange={(e) => setType(e.target.value)}
                                        disabled={isEditing}
                                    >
                                        <option value="image" className="bg-dark text-white">Image</option>
                                        <option value="video" className="bg-dark text-white">Video</option>
                                        <option value="document" className="bg-dark text-white">Document</option>
                                        <option value="other" className="bg-dark text-white">Other</option>
                                    </select>
                                </div>
                            </div>

                            {!isEditing && (
                                <div className="mb-4">
                                    <label className="form-label text-secondary fw-semibold">Media File</label>
                                    <input
                                        type="file"
                                        className="form-control bg-transparent"
                                        onChange={(e) => setFile(e.target.files ? e.target.files[0] : null)}
                                        required
                                    />
                                    <div className="form-text mt-2 opacity-50">Choose an image, video, or document to upload.</div>
                                </div>
                            )}

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
                                        <><i className="bi bi-check-circle me-2"></i>{isEditing ? 'Update Media' : 'Save Media'}</>
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

export default MediaForm;
