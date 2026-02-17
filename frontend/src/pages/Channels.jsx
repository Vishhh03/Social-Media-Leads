import { useState, useEffect } from 'react';
import { getChannels, connectChannel, disconnectChannel } from '../api';
import { useToast } from '../components/Toast';

export default function Channels() {
    const toast = useToast();
    const [channels, setChannels] = useState([]);
    const [loading, setLoading] = useState(true);
    const [showModal, setShowModal] = useState(false);
    const [form, setForm] = useState({ platform: 'whatsapp', account_id: '', account_name: '', access_token: '' });

    useEffect(() => { loadChannels(); }, []);

    async function loadChannels() {
        try {
            const data = await getChannels();
            setChannels(data.channels || []);
        } catch (err) {
            toast.error(err.message);
        } finally {
            setLoading(false);
        }
    }

    async function handleConnect(e) {
        e.preventDefault();
        try {
            await connectChannel(form.platform, form.account_id, form.account_name, form.access_token);
            toast.success('Channel connected successfully');
            setShowModal(false);
            setForm({ platform: 'whatsapp', account_id: '', account_name: '', access_token: '' });
            loadChannels();
        } catch (err) {
            toast.error(err.message);
        }
    }

    async function handleDisconnect(id) {
        if (!confirm('Disconnect this channel?')) return;
        try {
            await disconnectChannel(id);
            toast.success('Channel disconnected');
            loadChannels();
        } catch (err) {
            toast.error(err.message);
        }
    }

    const platformColors = {
        whatsapp: { bg: '#dcfce7', color: '#15803d', label: '🟢 WhatsApp' },
        instagram: { bg: '#fce7f3', color: '#be185d', label: '📸 Instagram' },
        facebook: { bg: '#dbeafe', color: '#1d4ed8', label: '📘 Facebook' },
    };

    if (loading) return <div className="loading-center"><div className="spinner"></div></div>;

    return (
        <div className="page-content">
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
                <div>
                    <h1 style={{ marginBottom: '4px' }}>Channels</h1>
                    <p style={{ color: 'var(--text-secondary)', fontSize: 'var(--text-sm)' }}>Connect your WhatsApp, Instagram, and Facebook accounts</p>
                </div>
                <button className="btn btn-primary" onClick={() => setShowModal(true)}>+ Connect Channel</button>
            </div>

            {channels.length === 0 ? (
                <div className="empty-state">
                    <div className="empty-state-icon">📡</div>
                    <div className="empty-state-title">No channels connected</div>
                    <div className="empty-state-text">Connect a WhatsApp, Instagram, or Facebook account to start receiving messages</div>
                </div>
            ) : (
                <div className="channel-grid">
                    {channels.map(ch => {
                        const pc = platformColors[ch.platform] || {};
                        return (
                            <div className="channel-card" key={ch.id}>
                                <div className="channel-card-platform" style={{ color: pc.color }}>{pc.label || ch.platform}</div>
                                <div className="channel-card-name">{ch.account_name || 'Unnamed'}</div>
                                <div className="channel-card-id">{ch.account_id}</div>
                                <div style={{ marginTop: '10px' }}>
                                    <span className={`badge ${ch.is_active ? 'badge-success' : 'badge-danger'}`}>
                                        {ch.is_active ? 'Active' : 'Inactive'}
                                    </span>
                                </div>
                                <div className="channel-card-actions">
                                    <button className="btn btn-sm btn-danger" onClick={() => handleDisconnect(ch.id)}>Disconnect</button>
                                </div>
                            </div>
                        );
                    })}
                </div>
            )}

            {showModal && (
                <div className="modal-backdrop" onClick={() => setShowModal(false)}>
                    <div className="modal" onClick={e => e.stopPropagation()}>
                        <h2>Connect Channel</h2>
                        <form onSubmit={handleConnect}>
                            <div className="form-group">
                                <label>Platform</label>
                                <select className="input" value={form.platform} onChange={e => setForm({ ...form, platform: e.target.value })}>
                                    <option value="whatsapp">WhatsApp</option>
                                    <option value="instagram">Instagram</option>
                                    <option value="facebook">Facebook</option>
                                </select>
                            </div>
                            <div className="form-group">
                                <label>Account ID</label>
                                <input className="input" value={form.account_id} onChange={e => setForm({ ...form, account_id: e.target.value })}
                                    placeholder={form.platform === 'whatsapp' ? 'Phone Number ID' : 'Page ID'} required />
                            </div>
                            <div className="form-group">
                                <label>Account Name</label>
                                <input className="input" value={form.account_name} onChange={e => setForm({ ...form, account_name: e.target.value })}
                                    placeholder="My Business Page" />
                            </div>
                            <div className="form-group">
                                <label>Access Token</label>
                                <textarea className="input" rows={3} value={form.access_token} onChange={e => setForm({ ...form, access_token: e.target.value })}
                                    placeholder="Paste your Meta access token here" required style={{ fontFamily: 'monospace', fontSize: '12px' }} />
                            </div>
                            <div className="modal-actions">
                                <button type="button" className="btn" onClick={() => setShowModal(false)}>Cancel</button>
                                <button type="submit" className="btn btn-primary">Connect</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
}
