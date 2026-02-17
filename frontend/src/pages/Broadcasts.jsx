import { useState, useEffect } from 'react';
import { getBroadcasts, createBroadcast, sendBroadcast } from '../api';
import { useToast } from '../components/Toast';

export default function Broadcasts() {
    const toast = useToast();
    const [broadcasts, setBroadcasts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [showModal, setShowModal] = useState(false);
    const [form, setForm] = useState({ name: '', content: '' });

    useEffect(() => { loadBroadcasts(); }, []);

    async function loadBroadcasts() {
        try {
            const data = await getBroadcasts();
            setBroadcasts(data.broadcasts || []);
        } catch (err) {
            toast.error(err.message);
        } finally {
            setLoading(false);
        }
    }

    async function handleCreate(e) {
        e.preventDefault();
        try {
            await createBroadcast({ name: form.name, content: form.content });
            toast.success('Broadcast draft created');
            setShowModal(false);
            setForm({ name: '', content: '' });
            loadBroadcasts();
        } catch (err) {
            toast.error(err.message);
        }
    }

    async function handleSend(id) {
        if (!confirm('Send this broadcast to ALL contacts? This cannot be undone.')) return;
        try {
            await sendBroadcast(id);
            toast.success('Broadcast sending started');
            loadBroadcasts();
        } catch (err) {
            toast.error(err.message);
        }
    }

    const statusBadge = (status) => {
        const map = {
            draft: 'badge-info',
            sending: 'badge-warning',
            sent: 'badge-success',
            failed: 'badge-danger',
        };
        return <span className={`badge ${map[status] || 'badge-info'}`}>{status}</span>;
    };

    if (loading) return <div className="loading-center"><div className="spinner"></div></div>;

    return (
        <div className="page-content">
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
                <div>
                    <h1 style={{ marginBottom: '4px' }}>Broadcasts</h1>
                    <p style={{ color: 'var(--text-secondary)', fontSize: 'var(--text-sm)' }}>Send bulk messages to all your contacts</p>
                </div>
                <button className="btn btn-primary" onClick={() => setShowModal(true)}>+ New Broadcast</button>
            </div>

            {broadcasts.length === 0 ? (
                <div className="empty-state">
                    <div className="empty-state-icon">📢</div>
                    <div className="empty-state-title">No broadcasts yet</div>
                    <div className="empty-state-text">Create a broadcast to message all your contacts at once</div>
                </div>
            ) : (
                <div className="broadcast-list">
                    {broadcasts.map(b => (
                        <div className="broadcast-item" key={b.id}>
                            <div className="broadcast-info">
                                <div className="broadcast-name">{b.name}</div>
                                <div className="broadcast-meta">
                                    {b.content.substring(0, 80)}{b.content.length > 80 ? '...' : ''}
                                    {b.total_sent > 0 && ` · ${b.total_sent} sent`}
                                    {b.total_failed > 0 && ` · ${b.total_failed} failed`}
                                </div>
                            </div>
                            <div className="broadcast-actions">
                                {statusBadge(b.status)}
                                {b.status === 'draft' && (
                                    <button className="btn btn-sm btn-primary" onClick={() => handleSend(b.id)}>Send Now</button>
                                )}
                            </div>
                        </div>
                    ))}
                </div>
            )}

            {showModal && (
                <div className="modal-backdrop" onClick={() => setShowModal(false)}>
                    <div className="modal" onClick={e => e.stopPropagation()}>
                        <h2>Create Broadcast</h2>
                        <form onSubmit={handleCreate}>
                            <div className="form-group">
                                <label>Campaign Name</label>
                                <input className="input" value={form.name} onChange={e => setForm({ ...form, name: e.target.value })}
                                    placeholder="e.g., New Year Promo" required />
                            </div>
                            <div className="form-group">
                                <label>Message</label>
                                <textarea className="input" rows={4} value={form.content} onChange={e => setForm({ ...form, content: e.target.value })}
                                    placeholder="Write your broadcast message..." required />
                            </div>
                            <div className="modal-actions">
                                <button type="button" className="btn" onClick={() => setShowModal(false)}>Cancel</button>
                                <button type="submit" className="btn btn-primary">Create Draft</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
}
