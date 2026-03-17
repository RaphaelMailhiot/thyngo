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
import MovieForm from '../modules/media/MovieForm';
import SeriesForm from '../modules/media/SeriesForm';
import MusicForm from '../modules/media/MusicForm';
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
                        <Route path="/media/movies/new" element={<MovieForm />} />
                        <Route path="/media/movies/edit/:id" element={<MovieForm />} />
                        <Route path="/media/series/new" element={<SeriesForm />} />
                        <Route path="/media/series/edit/:id" element={<SeriesForm />} />
                        <Route path="/media/albums/new" element={<MusicForm />} />
                        <Route path="/media/albums/edit/:id" element={<MusicForm />} />

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
