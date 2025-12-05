import { useEffect, useState } from 'react';

interface Post {
    slug: string;
    title: string;
    created_at: string;
    updated_at: string;
}

interface ApiResponse {
    data: Post[];
    success: boolean;
}

const GetPosts = () => {
    const [posts, setPosts] = useState<Post[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchPosts = async () => {
            try {
                const res = await fetch('http://api.localhost/api/posts', { cache: 'no-store' });
                const data: ApiResponse = await res.json();

                if (data.success && data.data) {
                    setPosts(data.data);
                } else {
                    setError('Erreur : impossible de récupérer les posts');
                }
            } catch (err) {
                setError(`Erreur de connexion : ${err}`);
            } finally {
                setLoading(false);
            }
        };

        fetchPosts();
    }, []);

    if (loading) return <p>Chargement...</p>;
    if (error) return <p style={{ color: 'red' }}>{error}</p>;
    if (posts.length === 0) return <p>Aucun post trouvé</p>;

    return (
        <table className='table table-bordered'>
            <thead>
            <tr>
                <th>slug</th>
                <th>title</th>
                <th>created_at</th>
                <th>updated_at</th>
            </tr>
            </thead>

            <tbody>
            {posts.map((post, index) => <tr key={index}>
                <td>{post.slug}</td>
                <td>{post.title}</td>
                <td>{post.created_at}</td>
                <td>{post.updated_at}</td>
            </tr>)}
            </tbody>
        </table>
    );
}

export default GetPosts;