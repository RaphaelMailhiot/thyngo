import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import api from '../../services/api';

interface Media {
    id: number;
    slug: string;
    title: string;
    type: string;
    url: string;
    file_name: string;
    file_size: number;
    mime_type: string;
    visibility: string;
    created_at: string;
}

const MediaList: React.FC = () => {
    const [media, setMedia] = useState<Media[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchMedia = async () => {
        try {
            const response = await api.get('/media');
            if (response.data && response.data.data) {
                setMedia(response.data.data);
            }
        } catch (error) {
            console.error('Failed to fetch media', error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchMedia();
    }, []);

    const handleDelete = async (slug: string) => {
        if (window.confirm('Are you sure you want to delete this media item?')) {
            try {
                await api.delete(`/media/${slug}`);
                fetchMedia(); // Refresh list
            } catch (error) {
                console.error('Failed to delete media', error);
                alert('Could not delete media.');
            }
        }
    };

    return (
        <div>
            <div className="d-flex justify-content-between align-items-center mb-4 border-bottom border-secondary border-opacity-25 pb-3">
                <h2 style={{ fontFamily: 'var(--font-heading)' }}>
                    <i className="bi bi-images me-2 text-primary"></i>
                    Media Library
                </h2>
                <Link to="/media/new" className="btn btn-primary shadow-sm hover-lift">
                    <i className="bi bi-cloud-upload me-2"></i>Upload Media
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
                                    <th className="ps-4">Resource</th>
                                    <th>Type</th>
                                    <th>Slug</th>
                                    <th>Visibility</th>
                                    <th>Created At</th>
                                    <th className="text-end pe-4">Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {media.length === 0 ? (
                                    <tr>
                                        <td colSpan={6} className="text-center py-5 text-muted">
                                            <i className="bi bi-image-fill fs-1 d-block mb-2 opacity-25"></i>
                                            Your media library is empty.
                                        </td>
                                    </tr>
                                ) : (
                                    media.map(m => (
                                        <tr key={m.id}>
                                            <td className="ps-4">
                                                <div className="d-flex align-items-center">
                                                    <div className="bg-dark rounded overflow-hidden me-3 d-flex align-items-center justify-content-center" style={{ width: '48px', height: '48px', border: '1px solid rgba(255,255,255,0.1)' }}>
                                                        {m.type === 'image' && m.url ? (
                                                            <img src={m.url} alt={m.title} style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
                                                        ) : (
                                                            <i className={`bi bi-${m.type === 'image' ? 'image' : m.type === 'video' ? 'play-btn' : 'file-earmark-text'} text-primary fs-5`}></i>
                                                        )}
                                                    </div>
                                                    <div>
                                                        <div className="fw-semibold">{m.title}</div>
                                                        <div className="text-secondary" style={{ fontSize: '0.75rem' }}>{m.file_name} ({(m.file_size / 1024).toFixed(1)} KB)</div>
                                                    </div>
                                                </div>
                                            </td>
                                            <td>
                                                <span className="text-secondary small text-capitalize">{m.type}</span>
                                            </td>
                                            <td><code className="bg-dark text-info px-2 py-1 rounded small">{m.slug}</code></td>
                                            <td>
                                                <span className={`badge bg-${m.visibility === 'public' ? 'success' : m.visibility === 'private' ? 'secondary' : 'info'} bg-opacity-25 text-${m.visibility === 'public' ? 'success' : m.visibility === 'private' ? 'secondary' : 'info'} border border-${m.visibility === 'public' ? 'success' : m.visibility === 'private' ? 'secondary' : 'info'} border-opacity-25`}>
                                                    {m.visibility}
                                                </span>
                                            </td>
                                            <td className="text-secondary" style={{ fontSize: '0.9rem' }}>{new Date(m.created_at).toLocaleDateString()}</td>
                                            <td className="text-end pe-4">
                                                <div className="btn-group btn-group-sm">
                                                    <Link to={`/media/edit/${m.slug}`} className="btn btn-outline-secondary border-0" title="Edit">
                                                        <i className="bi bi-pencil"></i>
                                                    </Link>
                                                    <button onClick={() => handleDelete(m.slug)} className="btn btn-outline-danger border-0" title="Delete">
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

export default MediaList;
