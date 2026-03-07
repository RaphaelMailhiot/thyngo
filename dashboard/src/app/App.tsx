import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from '../contexts/AuthContext';
import AdminLayout from '../layouts/AdminLayout';
import Login from '../pages/Login';
import Dashboard from '../pages/Dashboard';
import PostList from '../modules/posts/PostList';
import PostForm from '../modules/posts/PostForm';
import ProjectList from '../modules/projects/ProjectList';
import ProjectForm from '../modules/projects/ProjectForm';
import MediaList from '../modules/media/MediaList';
import MediaForm from '../modules/media/MediaForm';
import ResumeList from '../modules/resumes/ResumeList';
import ResumeForm from '../modules/resumes/ResumeForm';

function App() {
    return (
        <AuthProvider>
            <BrowserRouter>
                <Routes>
                    <Route path="/login" element={<Login />} />

                    <Route element={<AdminLayout />}>
                        <Route path="/" element={<Navigate to="/dashboard" replace />} />
                        <Route path="/dashboard" element={<Dashboard />} />

                        <Route path="/posts" element={<PostList />} />
                        <Route path="/posts/new" element={<PostForm />} />
                        <Route path="/posts/edit/:slug" element={<PostForm />} />

                        <Route path="/projects" element={<ProjectList />} />
                        <Route path="/projects/new" element={<ProjectForm />} />
                        <Route path="/projects/edit/:slug" element={<ProjectForm />} />

                        <Route path="/media" element={<MediaList />} />
                        <Route path="/media/new" element={<MediaForm />} />
                        <Route path="/media/edit/:slug" element={<MediaForm />} />

                        <Route path="/resumes" element={<ResumeList />} />
                        <Route path="/resumes/new" element={<ResumeForm />} />
                        <Route path="/resumes/edit/:id" element={<ResumeForm />} />
                    </Route>
                </Routes>
            </BrowserRouter>
        </AuthProvider>
    );
}

export default App;
