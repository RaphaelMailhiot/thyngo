import React from 'react';

const Dashboard: React.FC = () => {
    return (
        <div>
            <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
                <h1 className="h2">Dashboard</h1>
            </div>

            <div className="row">
                <div className="col-md-4 mb-4">
                    <div className="card text-white bg-primary shadow-sm h-100">
                        <div className="card-body">
                            <h5 className="card-title">Manage Posts</h5>
                            <p className="card-text">Create and edit blog posts, manage block content and SEO visibility.</p>
                            <a href="/posts" className="btn btn-light text-primary mt-2">Go to Posts</a>
                        </div>
                    </div>
                </div>

                <div className="col-md-4 mb-4">
                    <div className="card text-white bg-success shadow-sm h-100">
                        <div className="card-body">
                            <h5 className="card-title">Manage Projects</h5>
                            <p className="card-text">Track and display your portfolio work with specialized project categories.</p>
                            <a href="/projects" className="btn btn-light text-success mt-2">Go to Projects</a>
                        </div>
                    </div>
                </div>

                <div className="col-md-4 mb-4">
                    <div className="card text-white bg-info shadow-sm h-100">
                        <div className="card-body">
                            <h5 className="card-title">Media Library</h5>
                            <p className="card-text">Upload and manage images and documents used across your site.</p>
                            <a href="/media" className="btn btn-light text-info mt-2">Go to Media</a>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Dashboard;
