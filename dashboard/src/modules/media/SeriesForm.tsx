import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import api from '../../services/api';
import type { Series } from './types';

const SeriesForm: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const isEditing = !!id;

    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');
    const [loading, setLoading] = useState(isEditing);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState('');

    // Search states
    const [searchQuery, setSearchQuery] = useState('');
    const [searchResults, setSearchResults] = useState<any[]>([]);
    const [searching, setSearching] = useState(false);
    const [artworkResults, setArtworkResults] = useState<string[]>([]);
    const [importingArtwork, setImportingArtwork] = useState(false);

    useEffect(() => {
        if (isEditing) {
            const fetchSeries = async () => {
                try {
                    const response = await api.get(`/media/series`);
                    const series = response.data.data.find((s: Series) => s.id === Number(id));
                    if (series) {
                        setTitle(series.title);
                        setDescription(series.description || '');
                    }
                } catch (err) {
                    setError('Failed to load series data.');
                } finally {
                    setLoading(false);
                }
            };
            fetchSeries();
        }
    }, [id, isEditing]);

    const handleSearch = async () => {
        if (!searchQuery) return;
        setSearching(true);
        try {
            const response = await api.get(`/media/external/series/search?query=${encodeURIComponent(searchQuery)}`);
            setSearchResults(response.data.data || []);
        } catch (err) {
            console.error('Search failed', err);
        } finally {
            setSearching(false);
        }
    };

    const selectResult = (s: any) => {
        setTitle(s.name);
        setDescription(s.overview);
        setSearchResults([]);
        setSearchQuery('');

        // Fetch artwork from Fanart
        fetchArtwork(s.id);
    };

    const fetchArtwork = async (tmdbId: number) => {
        try {
            const response = await api.get(`/media/external/series/${tmdbId}/artwork`);
            if (response.data.data) {
                setArtworkResults(response.data.data);
            }
        } catch (err) {
            console.error('Failed to fetch series artwork', err);
        }
    };

    const handleUseArtwork = async (url: string) => {
        setImportingArtwork(true);
        try {
            const response = await api.post('/media/external/import-external', {
                title: `Poster for ${title}`,
                url,
                visibility: 'public'
            });
            if (response.data.success) {
                // For series, we don't have a direct resource image like Movie, 
                // but we might want to store it or just inform the user.
                // However, the user asked to "choose a cover", so for series 
                // we might need a "poster_media_id" in the schema if it existed.
                // Since it doesn't currently, we just import it.
                alert('Series poster imported into Media Library.');
                setArtworkResults([]);
            }
        } catch (err) {
            console.error('Failed to import artwork', err);
            setError('Failed to import artwork.');
        } finally {
            setImportingArtwork(false);
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSaving(true);
        setError('');

        try {
            if (isEditing) {
                await api.put(`/media/series/${id}`, { title, description });
            } else {
                await api.post('/media/series', { title, description });
            }
            navigate('/media');
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to save series.');
            setSaving(false);
        }
    };

    return (
        <div>
            <div className="d-flex justify-content-between align-items-center mb-4 border-bottom border-secondary border-opacity-25 pb-3">
                <h2 style={{ fontFamily: 'var(--font-heading)' }}>
                    <i className="bi bi-collection-play me-2 text-primary"></i>
                    {isEditing ? 'Edit Series' : 'Add New Series'}
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

                        {!isEditing && (
                            <div className="mb-4 p-3 bg-dark bg-opacity-25 rounded border border-secondary border-opacity-10">
                                <label className="form-label text-secondary fw-semibold small text-uppercase">Import from TMDB</label>
                                <div className="input-group">
                                    <input 
                                        type="text" 
                                        className="form-control bg-transparent" 
                                        placeholder="Search for a series..." 
                                        value={searchQuery}
                                        onChange={(e) => setSearchQuery(e.target.value)}
                                        onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), handleSearch())}
                                    />
                                    <button className="btn btn-outline-primary" type="button" onClick={handleSearch} disabled={searching}>
                                        {searching ? <span className="spinner-border spinner-border-sm" /> : <i className="bi bi-search"></i>}
                                    </button>
                                </div>
                                {searchResults.length > 0 && (
                                    <div className="mt-2 list-group list-group-flush border border-secondary border-opacity-10 rounded overflow-hidden shadow-sm">
                                        {searchResults.slice(0, 5).map(res => (
                                            <button 
                                                key={res.id} 
                                                type="button" 
                                                className="list-group-item list-group-item-action bg-dark text-white border-secondary border-opacity-10 py-2 small"
                                                onClick={() => selectResult(res)}
                                            >
                                                <div className="fw-bold">{res.name} <span className="text-secondary fw-normal">({res.first_air_date?.substring(0, 4)})</span></div>
                                                <div className="text-truncate opacity-50 small">{res.overview}</div>
                                            </button>
                                        ))}
                                    </div>
                                )}
                            </div>
                        )}

                        {artworkResults.length > 0 && (
                            <div className="mb-4">
                                <label className="form-label text-secondary fw-semibold small text-uppercase">Available Artwork (Fanart.tv)</label>
                                <div className="d-flex gap-2 overflow-auto pb-2" style={{ scrollbarWidth: 'thin' }}>
                                    {artworkResults.map((url, idx) => (
                                        <div key={idx} className="position-relative flex-shrink-0" style={{ width: '120px' }}>
                                            <img src={url} className="img-fluid rounded border border-secondary" style={{ height: '180px', objectFit: 'cover' }} />
                                            <button 
                                                type="button"
                                                className="btn btn-primary btn-sm position-absolute bottom-0 start-50 translate-middle-x mb-2 shadow-sm"
                                                style={{ fontSize: '10px' }}
                                                onClick={() => handleUseArtwork(url)}
                                                disabled={importingArtwork}
                                            >
                                                {importingArtwork ? '...' : 'Import Poster'}
                                            </button>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )}

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

                            <div className="mb-4">
                                <label className="form-label text-secondary fw-semibold">Description</label>
                                <textarea
                                    className="form-control bg-transparent"
                                    rows={5}
                                    value={description}
                                    onChange={(e) => setDescription(e.target.value)}
                                />
                            </div>

                            <div className="d-flex justify-content-end mt-5 pt-3 border-top border-secondary border-opacity-10">
                                <button type="submit" className="btn btn-primary px-5 py-2 shadow-sm hover-lift" disabled={saving}>
                                    {saving ? 'Saving...' : 'Save Series'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default SeriesForm;
