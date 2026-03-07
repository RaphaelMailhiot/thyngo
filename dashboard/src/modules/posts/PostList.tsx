import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import api from '../../services/api';

interface Post {
    id: number;
    slug: string;
    title: string;
    visibility: string;
    created_at: string;
}

const PostList: React.FC = () => {
    const [posts, setPosts] = useState<Post[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchPosts = async () => {
        try {
            const response = await api.get('/posts');
            if (response.data && response.data.data) {
                setPosts(response.data.data);
            }
        } catch (error) {
            console.error('Failed to fetch posts', error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchPosts();
    }, []);

    const handleDelete = async (slug: string) => {
        if (window.confirm('Are you sure you want to delete this post?')) {
            try {
                await api.delete(`/posts/${slug}`);
                fetchPosts(); // Refresh list
            } catch (error) {
                console.error('Failed to delete post', error);
                alert('Could not delete post.');
            }
        }
    };

    return (
        <div>
            <div className="d-flex justify-content-between align-items-center mb-4 border-bottom pb-2">
                <h2>Posts (Blogs)</h2>
                <Link to="/posts/new" className="btn btn-primary">
                    <i className="bi bi-plus-lg me-2"></i>Create New
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
                            <thead className="table-light">
                                <tr>
                                    <th>ID</th>
                                    <th>Title</th>
                                    <th>Slug</th>
                                    <th>Visibility</th>
                                    <th>Created At</th>
                                    <th className="text-end">Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {posts.length === 0 ? (
                                    <tr>
                                        <td colSpan={6} className="text-center py-4">No posts found. Create one!</td>
                                    </tr>
                                ) : (
                                    posts.map(post => (
                                        <tr key={post.id}>
                                            <td>{post.id}</td>
                                            <td className="fw-semibold">{post.title}</td>
                                            <td><code>{post.slug}</code></td>
                                            <td>
                                                <span className={`badge bg-${post.visibility === 'public' ? 'success' : post.visibility === 'private' ? 'secondary' : 'info'}`}>
                                                    {post.visibility}
                                                </span>
                                            </td>
                                            <td>{new Date(post.created_at).toLocaleDateString()}</td>
                                            <td className="text-end">
                                                <Link to={`/posts/edit/${post.slug}`} className="btn btn-sm btn-outline-secondary me-2">
                                                    <i className="bi bi-pencil"></i> Edit
                                                </Link>
                                                <button onClick={() => handleDelete(post.slug)} className="btn btn-sm btn-outline-danger">
                                                    <i className="bi bi-trash"></i> Delete
                                                </button>
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

export default PostList;
