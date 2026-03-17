import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import api from '../../services/api';
import type { Movie, Media } from './types';

const MovieForm: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const isEditing = !!id;

    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');
    const [releaseDate, setReleaseDate] = useState('');
    const [duration, setDuration] = useState(0);
    const [mediaId, setMediaId] = useState<number | ''>('');
    const [mediaList, setMediaList] = useState<Media[]>([]);
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
        const fetchMedia = async () => {
            try {
                const response = await api.get('/media');
                if (response.data && response.data.data) {
                    setMediaList(response.data.data.filter((m: Media) => m.type === 'video'));
                }
            } catch (err) {
                console.error('Failed to fetch video media', err);
            }
        };

        const fetchMovie = async () => {
            if (!isEditing) return;
            try {
                const response = await api.get(`/media/movies`);
                const movie = response.data.data.find((m: Movie) => m.id === Number(id));
                if (movie) {
                    setTitle(movie.title);
                    setDescription(movie.description || '');
                    setReleaseDate(movie.release_date ? movie.release_date.split('T')[0] : '');
                    setDuration(movie.duration || 0);
                    setMediaId(movie.media_id);
                }
            } catch (err) {
                setError('Failed to load movie data.');
            } finally {
                setLoading(false);
            }
        };

        fetchMedia();
        fetchMovie();
    }, [id, isEditing]);

    const handleSearch = async () => {
        if (!searchQuery) return;
        setSearching(true);
        try {
            const response = await api.get(`/media/external/movies/search?query=${encodeURIComponent(searchQuery)}`);
            setSearchResults(response.data.data || []);
        } catch (err) {
            console.error('Search failed', err);
        } finally {
            setSearching(false);
        }
    };

    const selectResult = (m: any) => {
        setTitle(m.title);
        setDescription(m.overview);
        setReleaseDate(m.release_date);
        setSearchResults([]);
        setSearchQuery('');

        // Fetch artwork from Fanart
        fetchArtwork(m.id);
    };

    const fetchArtwork = async (tmdbId: number) => {
        try {
            const response = await api.get(`/media/external/movies/${tmdbId}/artwork`);
            if (response.data.data) {
                setArtworkResults(response.data.data);
            }
        } catch (err) {
            console.error('Failed to fetch movie artwork', err);
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
                const newMedia = response.data.data;
                // Add to mediaList so it can be selected
                setMediaList(prev => [newMedia, ...prev]);
                setMediaId(newMedia.id);
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

        const movieData = {
            title,
            description,
            release_date: releaseDate ? new Date(releaseDate).toISOString() : null,
            duration: Number(duration),
            media_id: Number(mediaId)
        };

        try {
            if (isEditing) {
                await api.put(`/media/movies/${id}`, movieData);
            } else {
                await api.post('/media/movies', movieData);
            }
            navigate('/media');
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to save movie.');
            setSaving(false);
        }
    };

    return (
        <div>
            <div className="d-flex justify-content-between align-items-center mb-4 border-bottom border-secondary border-opacity-25 pb-3">
                <h2 style={{ fontFamily: 'var(--font-heading)' }}>
                    <i className="bi bi-film me-2 text-primary"></i>
                    {isEditing ? 'Edit Movie' : 'Add New Movie'}
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
                                        placeholder="Search for a movie..." 
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
                                                <div className="fw-bold">{res.title} <span className="text-secondary fw-normal">({res.release_date?.substring(0, 4)})</span></div>
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
                                                {importingArtwork ? '...' : 'Use as Poster'}
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
                                    rows={3}
                                    value={description}
                                    onChange={(e) => setDescription(e.target.value)}
                                />
                            </div>

                            <div className="row mb-4">
                                <div className="col-md-6">
                                    <label className="form-label text-secondary fw-semibold">Release Date</label>
                                    <input
                                        type="date"
                                        className="form-control bg-transparent"
                                        value={releaseDate}
                                        onChange={(e) => setReleaseDate(e.target.value)}
                                    />
                                </div>
                                <div className="col-md-6">
                                    <label className="form-label text-secondary fw-semibold">Duration (seconds)</label>
                                    <input
                                        type="number"
                                        className="form-control bg-transparent"
                                        value={duration}
                                        onChange={(e) => setDuration(Number(e.target.value))}
                                    />
                                </div>
                            </div>

                            <div className="mb-4">
                                <label className="form-label text-secondary fw-semibold">Video Resource</label>
                                <select
                                    className="form-select bg-transparent"
                                    value={mediaId}
                                    onChange={(e) => setMediaId(Number(e.target.value))}
                                    required
                                >
                                    <option value="" className="bg-dark text-white">Select a video...</option>
                                    {mediaList.map(m => (
                                        <option key={m.id} value={m.id} className="bg-dark text-white">{m.title} ({m.slug})</option>
                                    ))}
                                </select>
                                <div className="form-text mt-2 opacity-50">Link this movie to a file in your Media Library.</div>
                            </div>

                            <div className="d-flex justify-content-end mt-5 pt-3 border-top border-secondary border-opacity-10">
                                <button type="submit" className="btn btn-primary px-5 py-2 shadow-sm hover-lift" disabled={saving}>
                                    {saving ? 'Saving...' : 'Save Movie'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default MovieForm;
