import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import api from '../../services/api';
import type { Media, Movie, Series, MusicAlbum } from './types';

type Category = 'files' | 'movies' | 'series' | 'albums';

const MediaList: React.FC = () => {
    const [media, setMedia] = useState<Media[]>([]);
    const [movies, setMovies] = useState<Movie[]>([]);
    const [series, setSeries] = useState<Series[]>([]);
    const [albums, setAlbums] = useState<MusicAlbum[]>([]);
    const [category, setCategory] = useState<Category>('files');
    const [loading, setLoading] = useState(true);

    const fetchData = async () => {
        setLoading(true);
        try {
            let endpoint = '/media';
            if (category === 'movies') endpoint = '/media/movies';
            else if (category === 'series') endpoint = '/media/series';
            else if (category === 'albums') endpoint = '/media/albums';

            const response = await api.get(endpoint);
            if (response.data && response.data.data) {
                if (category === 'files') setMedia(response.data.data);
                else if (category === 'movies') setMovies(response.data.data);
                else if (category === 'series') setSeries(response.data.data);
                else if (category === 'albums') setAlbums(response.data.data);
            }
        } catch (error) {
            console.error(`Failed to fetch ${category}`, error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, [category]);

    const handleDelete = async (slug: string) => {
        if (window.confirm('Are you sure you want to delete this media item?')) {
            try {
                await api.delete(`/media/${slug}`);
                fetchData(); // Refresh list
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
                <div className="d-flex gap-2">
                    <div className="btn-group me-3">
                        <button 
                            className={`btn btn-sm ${category === 'files' ? 'btn-primary' : 'btn-outline-secondary'}`}
                            onClick={() => setCategory('files')}
                        >Files</button>
                        <button 
                            className={`btn btn-sm ${category === 'movies' ? 'btn-primary' : 'btn-outline-secondary'}`}
                            onClick={() => setCategory('movies')}
                        >Movies</button>
                        <button 
                            className={`btn btn-sm ${category === 'series' ? 'btn-primary' : 'btn-outline-secondary'}`}
                            onClick={() => setCategory('series')}
                        >Series</button>
                        <button 
                            className={`btn btn-sm ${category === 'albums' ? 'btn-primary' : 'btn-outline-secondary'}`}
                            onClick={() => setCategory('albums')}
                        >Music</button>
                    </div>
                    <Link to={category === 'files' ? "/media/new" : `/media/${category}/new`} className="btn btn-primary shadow-sm hover-lift">
                        <i className={`bi bi-${category === 'files' ? 'cloud-upload' : 'plus-lg'} me-2`}></i>
                        Add {category === 'files' ? 'Media' : category === 'albums' ? 'Album' : category.slice(0, -1)}
                    </Link>
                </div>
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
                                    {category === 'files' ? (
                                        <>
                                            <th className="ps-4">Resource</th>
                                            <th>Type</th>
                                            <th>Slug</th>
                                            <th>Visibility</th>
                                            <th>Created At</th>
                                        </>
                                    ) : (
                                        <>
                                            <th className="ps-4">Title</th>
                                            <th>{category === 'movies' ? 'Release Date' : category === 'series' ? 'Description' : 'Artist'}</th>
                                            <th>Created At</th>
                                        </>
                                    )}
                                    <th className="text-end pe-4">Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {category === 'files' ? (
                                    media.length === 0 ? (
                                        <tr><td colSpan={6} className="text-center py-5 text-muted">Empty</td></tr>
                                    ) : (
                                        media.map(m => (
                                            <tr key={m.id}>
                                                <td className="ps-4">
                                                    <div className="d-flex align-items-center">
                                                        <div className="bg-dark rounded overflow-hidden me-3 d-flex align-items-center justify-content-center" style={{ width: '48px', height: '48px' }}>
                                                            {m.type === 'image' && m.url ? <img src={m.url} style={{ width: '100%', height: '100%', objectFit: 'cover' }} /> : <i className="bi bi-file-earmark-text text-primary"></i>}
                                                        </div>
                                                        <div>{m.title}</div>
                                                    </div>
                                                </td>
                                                <td>{m.type}</td>
                                                <td><code>{m.slug}</code></td>
                                                <td>{m.visibility}</td>
                                                <td>{new Date(m.created_at).toLocaleDateString()}</td>
                                                <td className="text-end pe-4">
                                                    <Link to={`/media/edit/${m.slug}`} className="btn btn-sm btn-outline-secondary border-0"><i className="bi bi-pencil"></i></Link>
                                                    <button onClick={() => handleDelete(m.slug)} className="btn btn-sm btn-outline-danger border-0"><i className="bi bi-trash"></i></button>
                                                </td>
                                            </tr>
                                        ))
                                    )
                                ) : category === 'movies' ? (
                                    movies.length === 0 ? (
                                        <tr><td colSpan={4} className="text-center py-5 text-muted">No movies found.</td></tr>
                                    ) : (
                                        movies.map(m => (
                                            <tr key={m.id}>
                                                <td className="ps-4 font-semibold">{m.title}</td>
                                                <td>{m.release_date ? new Date(m.release_date).getFullYear() : 'N/A'}</td>
                                                <td>{new Date(m.created_at).toLocaleDateString()}</td>
                                                <td className="text-end pe-4">
                                                    <Link to={`/media/movies/edit/${m.id}`} className="btn btn-sm btn-outline-secondary border-0"><i className="bi bi-pencil"></i></Link>
                                                </td>
                                            </tr>
                                        ))
                                    )
                                ) : category === 'series' ? (
                                    series.length === 0 ? (
                                        <tr><td colSpan={4} className="text-center py-5 text-muted">No series found.</td></tr>
                                    ) : (
                                        series.map(s => (
                                            <tr key={s.id}>
                                                <td className="ps-4 font-semibold">{s.title}</td>
                                                <td className="text-truncate" style={{ maxWidth: '200px' }}>{s.description}</td>
                                                <td>{new Date(s.created_at).toLocaleDateString()}</td>
                                                <td className="text-end pe-4">
                                                    <Link to={`/media/series/edit/${s.id}`} className="btn btn-sm btn-outline-secondary border-0"><i className="bi bi-pencil"></i></Link>
                                                </td>
                                            </tr>
                                        ))
                                    )
                                ) : (
                                    albums.length === 0 ? (
                                        <tr><td colSpan={4} className="text-center py-5 text-muted">No albums found.</td></tr>
                                    ) : (
                                        albums.map(a => (
                                            <tr key={a.id}>
                                                <td className="ps-4 font-semibold">{a.title}</td>
                                                <td>{a.artist}</td>
                                                <td>{new Date(a.created_at).toLocaleDateString()}</td>
                                                <td className="text-end pe-4">
                                                    <Link to={`/media/albums/edit/${a.id}`} className="btn btn-sm btn-outline-secondary border-0"><i className="bi bi-pencil"></i></Link>
                                                </td>
                                            </tr>
                                        ))
                                    )
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
