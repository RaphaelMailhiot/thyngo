import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, Link } from 'react-router-dom';
import api from '../../services/api';
import type { MusicAlbum, Media } from './types';

const MusicForm: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const isEditing = !!id;

    const [title, setTitle] = useState('');
    const [artist, setArtist] = useState('');
    const [releaseDate, setReleaseDate] = useState('');
    const [coverId, setCoverId] = useState<number | ''>('');
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
                    setMediaList(response.data.data.filter((m: Media) => m.type === 'image'));
                }
            } catch (err) {
                console.error('Failed to fetch image media', err);
            }
        };

        const fetchAlbum = async () => {
            if (!isEditing) return;
            try {
                const response = await api.get(`/media/albums`);
                const album = response.data.data.find((a: MusicAlbum) => a.id === Number(id));
                if (album) {
                    setTitle(album.title);
                    setArtist(album.artist || '');
                    setReleaseDate(album.release_date ? album.release_date.split('T')[0] : '');
                    setCoverId(album.cover_media_id || '');
                }
            } catch (err) {
                setError('Failed to load album data.');
            } finally {
                setLoading(false);
            }
        };

        fetchMedia();
        fetchAlbum();
    }, [id, isEditing]);

    const handleSearch = async () => {
        if (!searchQuery) return;
        setSearching(true);
        try {
            const response = await api.get(`/media/external/music/search?query=${encodeURIComponent(searchQuery)}`);
            setSearchResults(response.data.data || []);
        } catch (err) {
            console.error('Search failed', err);
        } finally {
            setSearching(false);
        }
    };

    const selectResult = async (a: any) => {
        setTitle(a.title);
        setArtist(a.artist);
        setReleaseDate(a.date);
        setSearchResults([]);
        setSearchQuery('');

        // Try to fetch artwork
        try {
            const artworkResp = await api.get(`/media/external/music/albums/${a.id}/artwork`);
            if (artworkResp.data.data) {
                setArtworkResults(artworkResp.data.data);
            }
        } catch (err) {
            console.error('Artwork fetch failed', err);
        }
    };

    const handleUseArtwork = async (url: string) => {
        setImportingArtwork(true);
        try {
            const response = await api.post('/media/external/import-external', {
                title: `Cover for ${title}`,
                url,
                visibility: 'public'
            });
            if (response.data.success) {
                const newMedia = response.data.data;
                // Add to mediaList
                setMediaList(prev => [newMedia, ...prev]);
                setCoverId(newMedia.id);
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

        const albumData = {
            title,
            artist,
            release_date: releaseDate ? new Date(releaseDate).toISOString() : null,
            cover_media_id: coverId === '' ? null : Number(coverId)
        };

        try {
            if (isEditing) {
                await api.put(`/media/albums/${id}`, albumData);
            } else {
                await api.post('/media/albums', albumData);
            }
            navigate('/media');
        } catch (err: any) {
            setError(err.response?.data?.error || 'Failed to save album.');
            setSaving(false);
        }
    };

    return (
        <div>
            <div className="d-flex justify-content-between align-items-center mb-4 border-bottom border-secondary border-opacity-25 pb-3">
                <h2 style={{ fontFamily: 'var(--font-heading)' }}>
                    <i className="bi bi-disc me-2 text-primary"></i>
                    {isEditing ? 'Edit Album' : 'Add New Album'}
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
                                <label className="form-label text-secondary fw-semibold small text-uppercase">Import from MusicBrainz</label>
                                <div className="input-group">
                                    <input 
                                        type="text" 
                                        className="form-control bg-transparent" 
                                        placeholder="Search for an album..." 
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
                                                <div className="fw-bold">{res.title}</div>
                                                <div className="opacity-50 small">{res.artist} • {res.date}</div>
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
                                            <img src={url} className="img-fluid rounded border border-secondary" style={{ height: '120px', width: '120px', objectFit: 'cover' }} />
                                            <button 
                                                type="button"
                                                className="btn btn-primary btn-sm position-absolute bottom-0 start-50 translate-middle-x mb-2 shadow-sm"
                                                style={{ fontSize: '10px' }}
                                                onClick={() => handleUseArtwork(url)}
                                                disabled={importingArtwork}
                                            >
                                                {importingArtwork ? '...' : 'Use as Cover'}
                                            </button>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )}

                        <form onSubmit={handleSubmit}>
                            <div className="mb-4">
                                <label className="form-label text-secondary fw-semibold">Album Title</label>
                                <input
                                    type="text"
                                    className="form-control form-control-lg bg-transparent"
                                    value={title}
                                    onChange={(e) => setTitle(e.target.value)}
                                    required
                                />
                            </div>

                            <div className="mb-4">
                                <label className="form-label text-secondary fw-semibold">Artist</label>
                                <input
                                    type="text"
                                    className="form-control bg-transparent"
                                    value={artist}
                                    onChange={(e) => setArtist(e.target.value)}
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
                                    <label className="form-label text-secondary fw-semibold">Cover Image</label>
                                    <select
                                        className="form-select bg-transparent"
                                        value={coverId}
                                        onChange={(e) => setCoverId(e.target.value ? Number(e.target.value) : '')}
                                    >
                                        <option value="" className="bg-dark text-white">No cover image</option>
                                        {mediaList.map(m => (
                                            <option key={m.id} value={m.id} className="bg-dark text-white">{m.title} ({m.slug})</option>
                                        ))}
                                    </select>
                                </div>
                            </div>

                            <div className="d-flex justify-content-end mt-5 pt-3 border-top border-secondary border-opacity-10">
                                <button type="submit" className="btn btn-primary px-5 py-2 shadow-sm hover-lift" disabled={saving}>
                                    {saving ? 'Saving...' : 'Save Album'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default MusicForm;
