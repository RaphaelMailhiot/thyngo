import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { useAuth } from '../contexts/AuthContext';

const Login: React.FC = () => {
    const [identifier, setIdentifier] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const navigate = useNavigate();
    const { login } = useAuth();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        try {
            const response = await axios.post('/api/users/login', {
                identifier,
                password
            });

            if (response.data.success && response.data.token) {
                login(response.data.token);
                navigate('/dashboard');
            } else {
                setError(response.data.error || 'Login failed');
            }
        } catch (err: any) {
            setError(err.response?.data?.error || 'An error occurred during login');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="d-flex align-items-center justify-content-center bg-light" style={{ minHeight: '100vh' }}>
            <div className="card shadow-sm" style={{ width: '400px' }}>
                <div className="card-body p-4">
                    <h3 className="card-title text-center mb-4">ThynGo Admin</h3>

                    {error && (
                        <div className="alert alert-danger" role="alert">
                            {error}
                        </div>
                    )}

                    <form onSubmit={handleSubmit}>
                        <div className="mb-3">
                            <label className="form-label">Username or Email</label>
                            <input
                                type="text"
                                className="form-control"
                                value={identifier}
                                onChange={(e) => setIdentifier(e.target.value)}
                                required
                            />
                        </div>
                        <div className="mb-4">
                            <label className="form-label">Password</label>
                            <input
                                type="password"
                                className="form-control"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                required
                            />
                        </div>
                        <button
                            type="submit"
                            className="btn btn-primary w-100"
                            disabled={loading}
                        >
                            {loading ? 'Signing in...' : 'Sign in'}
                        </button>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default Login;
