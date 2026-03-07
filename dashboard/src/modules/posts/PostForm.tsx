import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import api from '../../services/api';
import BlockEditor from './components/BlockEditor';

const PostForm: React.FC = () => {
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
            const fetchPost = async () => {
                try {
                    const response = await api.get(`/posts/${slug}`);
                    if (response.data && response.data.data) {
                        const post = response.data.data;
                        setTitle(post.title);
                        setNewSlug(post.slug);
                        setVisibility(post.visibility);
                    }
                } catch (err: any) {
                    setError('Failed to load post data.');
                } finally {
                    setLoading(false);
                }
            };
            fetchPost();
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
                await api.put(`/posts/${slug}`, payload);
            } else {
                await api.post('/posts', payload);
            }
            navigate('/posts');
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to save post.');
            setSaving(false);
        }
    };

    return (
        <div>
            <div className="d-flex justify-content-between align-items-center mb-4 border-bottom pb-2">
                <h2>{isEditing ? 'Edit Post' : 'Create New Post'}</h2>
                <Link to="/posts" className="btn btn-outline-secondary">
                    Cancel
                </Link>
            </div>

            {loading ? (
                <div className="text-center mt-5"><div className="spinner-border text-primary" /></div>
            ) : (
                <>
                    <div className="card glass-panel border-0" style={{ maxWidth: '800px', margin: '0 auto' }}>
                        <div className="card-body p-4">
                            {error && <div className="alert alert-danger border-0 bg-danger bg-opacity-10 text-danger">{error}</div>}

                            <form onSubmit={handleSubmit}>
                                <div className="mb-4">
                                    <label className="form-label text-secondary fw-semibold">Title</label>
                                    <input
                                        type="text"
                                        className="form-control form-control-lg bg-transparent"
                                        value={title}
                                        onChange={(e) => setTitle(e.target.value)}
                                        required
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
                                            placeholder="e.g. my-first-post"
                                        />
                                        {isEditing && <div className="form-text">Slugs cannot be changed after creation.</div>}
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

                                <div className="d-flex justify-content-end">
                                    <button type="submit" className="btn btn-primary px-4 py-2" disabled={saving}>
                                        {saving ? 'Saving...' : (isEditing ? 'Save Post Meta' : 'Create Post & Start Writing')}
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>

                    {/* Show Block Editor only when editing an existing post */}
                    {isEditing && slug && (
                        <div style={{ maxWidth: '800px', margin: '0 auto' }}>
                            <BlockEditor postSlug={slug} />
                        </div>
                    )}
                </>
            )}
        </div>
    );
};

export default PostForm;
