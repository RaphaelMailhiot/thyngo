import React, { useState, useEffect } from 'react';
import api from '../../../services/api';

interface Block {
    id: number;
    type: string;
    data: any;
}

interface Content {
    id: number;
    parent_id: number;
    order: number;
    type: string;
    block: Block;
}

interface BlockEditorProps {
    postSlug: string;
}

const BlockEditor: React.FC<BlockEditorProps> = ({ postSlug }) => {
    const [contents, setContents] = useState<Content[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchContents = async () => {
        try {
            setLoading(true);
            const response = await api.get(`/posts/${postSlug}/contents`);
            if (response.data && response.data.data) {
                // Ensure they are sorted by order
                const sorted = response.data.data.sort((a: Content, b: Content) => a.order - b.order);
                setContents(sorted);
            }
        } catch (error) {
            console.error("Failed to fetch contents", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (postSlug) {
            fetchContents();
        }
    }, [postSlug]);

    const addBlock = async (type: string) => {
        try {
            // First, create the block data based on type
            let blockData = {};
            if (type === 'title') blockData = { level: 2, content: '' };
            if (type === 'text') blockData = { content: '' };
            if (type === 'code') blockData = { language: 'javascript', content: '' };
            if (type === 'image') blockData = { url: '', title: '' };

            const blockPayload = { type, data: blockData };
            const blockRes = await api.post('/posts/blocks', blockPayload);
            const newBlock = blockRes.data.data;

            // Then link it as content
            const nextOrder = contents.length;
            const contentPayload = {
                order: nextOrder,
                type: type,
                block_id: newBlock.id
            };
            await api.post(`/posts/${postSlug}/contents`, contentPayload);

            // Refresh
            fetchContents();
        } catch (err) {
            console.error('Failed to add block', err);
            alert('Error adding block.');
        }
    };

    const handleLocalUpdate = (contentId: number, newData: any) => {
        setContents(contents.map(c => {
            if (c.id === contentId && c.block) {
                return { ...c, block: { ...c.block, data: newData } };
            }
            return c;
        }));
    };

    const handleContentSave = async (blockId: number, type: string, newData: any) => {
        try {
            await api.put(`/posts/blocks/${blockId}`, { type, data: newData });
        } catch (err) {
            console.error('Failed to update block', err);
        }
    };

    const moveContent = async (index: number, direction: 'up' | 'down') => {
        if (direction === 'up' && index === 0) return;
        if (direction === 'down' && index === contents.length - 1) return;

        const newContents = [...contents];
        const targetIndex = direction === 'up' ? index - 1 : index + 1;

        // Swap orders
        const tempOrder = newContents[index].order;
        newContents[index].order = newContents[targetIndex].order;
        newContents[targetIndex].order = tempOrder;

        // Swap positions in array for immediate UI update
        const tempContent = newContents[index];
        newContents[index] = newContents[targetIndex];
        newContents[targetIndex] = tempContent;

        setContents(newContents);

        // Save to API
        try {
            await Promise.all([
                api.put(`/posts/contents/${newContents[index].id}`, { order: newContents[index].order }),
                api.put(`/posts/contents/${newContents[targetIndex].id}`, { order: newContents[targetIndex].order })
            ]);
        } catch (err) {
            console.error('Failed to swap orders', err);
            fetchContents(); // revert on fail
        }
    };

    const deleteContent = async (contentId: number) => {
        if (!window.confirm("Delete this block?")) return;
        try {
            await api.delete(`/posts/contents/${contentId}`);
            fetchContents();
        } catch (err) {
            console.error("Failed to delete", err);
        }
    };

    if (loading) return <div className="text-center p-4">Loading Editor...</div>;

    return (
        <div className="mt-5">
            <h4 className="mb-4" style={{ fontFamily: 'var(--font-heading)' }}>
                <i className="bi bi-layout-wtf me-2 text-primary"></i>
                Block Editor
            </h4>

            <div className="block-list mb-4">
                {contents.length === 0 ? (
                    <div className="text-center p-5 border rounded glass-panel text-muted">
                        <i className="bi bi-magic fs-1 d-block mb-2"></i>
                        No blocks yet. Add your first block below!
                    </div>
                ) : (
                    contents.map((content, index) => (
                        <div key={content.id} className="card glass-panel border-0 mb-3 hover-lift relative group">
                            <div className="card-body">
                                <div className="d-flex justify-content-between mb-2 pb-2 border-bottom border-secondary border-opacity-25">
                                    <span className="badge bg-secondary bg-opacity-25 text-primary text-uppercase" style={{ letterSpacing: '1px' }}>
                                        {content.type}
                                    </span>
                                    <div className="btn-group btn-group-sm">
                                        <button
                                            onClick={() => moveContent(index, 'up')}
                                            disabled={index === 0}
                                            className="btn btn-outline-secondary border-0"
                                            title="Move Up"
                                        ><i className="bi bi-arrow-up"></i></button>
                                        <button
                                            onClick={() => moveContent(index, 'down')}
                                            disabled={index === contents.length - 1}
                                            className="btn btn-outline-secondary border-0"
                                            title="Move Down"
                                        ><i className="bi bi-arrow-down"></i></button>
                                        <button
                                            onClick={() => deleteContent(content.id)}
                                            className="btn btn-outline-danger border-0 ms-2"
                                            title="Delete Block"
                                        ><i className="bi bi-trash"></i></button>
                                    </div>
                                </div>

                                {/* Dynamic Block Rendering */}
                                <div className="block-content-editor mt-3">
                                    {content.type === 'title' && (
                                        <input
                                            type="text"
                                            className="form-control form-control-lg fw-bold border-0 bg-transparent fs-2 px-0"
                                            placeholder="Heading..."
                                            value={content.block?.data?.content || ''}
                                            onChange={(e) => handleLocalUpdate(content.id, { ...content.block.data, content: e.target.value })}
                                            onBlur={() => handleContentSave(content.block.id, content.type, content.block.data)}
                                        />
                                    )}
                                    {content.type === 'text' && (
                                        <textarea
                                            className="form-control border-0 bg-transparent px-0"
                                            placeholder="Write your text or paragraph here..."
                                            rows={4}
                                            style={{ resize: 'none' }}
                                            value={content.block?.data?.content || ''}
                                            onChange={(e) => handleLocalUpdate(content.id, { ...content.block.data, content: e.target.value })}
                                            onBlur={() => handleContentSave(content.block.id, content.type, content.block.data)}
                                        />
                                    )}
                                    {content.type === 'code' && (
                                        <div className="bg-dark p-3 rounded position-relative">
                                            <div className="position-absolute" style={{ top: '10px', right: '10px' }}>
                                                <input
                                                    type="text"
                                                    className="form-control form-control-sm bg-black text-white border-0"
                                                    placeholder="language"
                                                    value={content.block?.data?.language || ''}
                                                    onChange={(e) => handleLocalUpdate(content.id, { ...content.block.data, language: e.target.value })}
                                                    onBlur={() => handleContentSave(content.block.id, content.type, content.block.data)}
                                                />
                                            </div>
                                            <textarea
                                                className="form-control border-0 bg-transparent text-success font-monospace"
                                                placeholder="// Write code here..."
                                                rows={5}
                                                value={content.block?.data?.content || ''}
                                                onChange={(e) => handleLocalUpdate(content.id, { ...content.block.data, content: e.target.value })}
                                                onBlur={() => handleContentSave(content.block.id, content.type, content.block.data)}
                                            />
                                        </div>
                                    )}
                                    {content.type === 'image' && (
                                        <div className="d-flex flex-column gap-2">
                                            <input
                                                type="text"
                                                className="form-control bg-transparent"
                                                placeholder="https://image.url/here.jpg"
                                                value={content.block?.data?.url || ''}
                                                onChange={(e) => handleLocalUpdate(content.id, { ...content.block.data, url: e.target.value })}
                                                onBlur={() => handleContentSave(content.block.id, content.type, content.block.data)}
                                            />
                                            {content.block?.data?.url && (
                                                <img src={content.block.data.url} alt="Preview" className="img-fluid rounded mt-2 shadow-sm" style={{ maxHeight: '300px', objectFit: 'cover' }} />
                                            )}
                                        </div>
                                    )}
                                </div>
                            </div>
                        </div>
                    ))
                )}
            </div>

            <div className="d-flex gap-2 justify-content-center p-3 glass-panel rounded-pill sticky-bottom shadow-lg mx-auto" style={{ maxWidth: '500px', bottom: '20px' }}>
                <span className="text-secondary align-self-center me-2 fw-semibold">Add:</span>
                <button onClick={() => addBlock('title')} className="btn btn-sm btn-outline-primary rounded-pill px-3">
                    <i className="bi bi-type-h1 me-1"></i> Title
                </button>
                <button onClick={() => addBlock('text')} className="btn btn-sm btn-outline-primary rounded-pill px-3">
                    <i className="bi bi-justify me-1"></i> Text
                </button>
                <button onClick={() => addBlock('image')} className="btn btn-sm btn-outline-primary rounded-pill px-3">
                    <i className="bi bi-image me-1"></i> Image
                </button>
                <button onClick={() => addBlock('code')} className="btn btn-sm btn-outline-primary rounded-pill px-3">
                    <i className="bi bi-code-slash me-1"></i> Code
                </button>
            </div>
        </div>
    );
};

export default BlockEditor;
