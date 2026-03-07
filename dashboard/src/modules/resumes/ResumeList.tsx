import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import api from '../../services/api';

interface Resume {
    id: number;
    name: string;
    job: string;
    visibility: string;
    created_at: string;
}

const ResumeList: React.FC = () => {
    const [resumes, setResumes] = useState<Resume[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchResumes = async () => {
        try {
            const response = await api.get('/resumes');
            if (response.data && response.data.data) {
                setResumes(response.data.data);
            }
        } catch (error) {
            console.error('Failed to fetch resumes', error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchResumes();
    }, []);

    const handleDelete = async (id: number) => {
        if (window.confirm('Are you sure you want to delete this resume?')) {
            try {
                await api.delete(`/resumes/${id}`);
                fetchResumes();
            } catch (error) {
                console.error('Failed to delete resume', error);
                alert('Could not delete resume.');
            }
        }
    };

    return (
        <div>
            <div className="d-flex justify-content-between align-items-center mb-4 border-bottom border-secondary border-opacity-25 pb-3">
                <h2 style={{ fontFamily: 'var(--font-heading)' }}>
                    <i className="bi bi-file-earmark-person me-2 text-primary"></i>
                    Resumes
                </h2>
                <Link to="/resumes/new" className="btn btn-primary shadow-sm hover-lift">
                    <i className="bi bi-plus-lg me-2"></i>Create Resume
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
                                    <th className="ps-4">Candidate</th>
                                    <th>Job Title</th>
                                    <th>Visibility</th>
                                    <th>Created At</th>
                                    <th className="text-end pe-4">Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {resumes.length === 0 ? (
                                    <tr>
                                        <td colSpan={5} className="text-center py-5 text-muted">
                                            <i className="bi bi-person-badge fs-1 d-block mb-2 opacity-25"></i>
                                            No resumes found. Start by creating one!
                                        </td>
                                    </tr>
                                ) : (
                                    resumes.map(resume => (
                                        <tr key={resume.id}>
                                            <td className="ps-4">
                                                <div className="d-flex align-items-center">
                                                    <div className="bg-primary bg-opacity-10 rounded-circle p-2 me-3 d-flex align-items-center justify-content-center" style={{ width: '40px', height: '40px' }}>
                                                        <span className="text-primary fw-bold">{resume.name.charAt(0)}</span>
                                                    </div>
                                                    <span className="fw-semibold">{resume.name}</span>
                                                </div>
                                            </td>
                                            <td className="text-secondary small">{resume.job || 'Not specified'}</td>
                                            <td>
                                                <span className={`badge bg-${resume.visibility === 'public' ? 'success' : resume.visibility === 'private' ? 'secondary' : 'info'} bg-opacity-25 text-${resume.visibility === 'public' ? 'success' : resume.visibility === 'private' ? 'secondary' : 'info'} border border-${resume.visibility === 'public' ? 'success' : resume.visibility === 'private' ? 'secondary' : 'info'} border-opacity-25`}>
                                                    {resume.visibility}
                                                </span>
                                            </td>
                                            <td className="text-secondary" style={{ fontSize: '0.9rem' }}>{new Date(resume.created_at).toLocaleDateString()}</td>
                                            <td className="text-end pe-4">
                                                <div className="btn-group btn-group-sm">
                                                    <Link to={`/resumes/edit/${resume.id}`} className="btn btn-outline-secondary border-0" title="Edit">
                                                        <i className="bi bi-pencil"></i>
                                                    </Link>
                                                    <button onClick={() => handleDelete(resume.id)} className="btn btn-outline-danger border-0" title="Delete">
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

export default ResumeList;
