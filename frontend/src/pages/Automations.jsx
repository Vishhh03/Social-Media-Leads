import { useState, useEffect } from 'react';
import { getAutomations, createAutomation, deleteAutomation } from '../api';
import { useToast } from '../components/Toast';

export default function Automations() {
    const toast = useToast();
    const [automations, setAutomations] = useState([]);
    const [loading, setLoading] = useState(true);
    const [showModal, setShowModal] = useState(false);
    const [form, setForm] = useState({ name: '', trigger_type: 'keyword', keywords: '', reply_text: '', delay_ms: 0 });

    useEffect(() => { loadAutomations(); }, []);

    async function loadAutomations() {
        try {
            const data = await getAutomations();
            setAutomations(data.automations || []);
        } catch (err) {
            toast.error(err.message);
        } finally {
            setLoading(false);
        }
    }

    async function handleCreate(e) {
        e.preventDefault();
        try {
            await createAutomation({
                name: form.name,
                trigger_type: form.trigger_type,
                keywords: form.keywords.split(',').map(k => k.trim()).filter(Boolean),
                reply_text: form.reply_text,
                delay_ms: parseInt(form.delay_ms) || 0,
                is_active: true,
            });
            toast.success('Automation created');
            setShowModal(false);
            setForm({ name: '', trigger_type: 'keyword', keywords: '', reply_text: '', delay_ms: 0 });
            loadAutomations();
        } catch (err) {
            toast.error(err.message);
        }
    }

    async function handleDelete(id) {
        if (!confirm('Delete this automation?')) return;
        try {
            await deleteAutomation(id);
            toast.success('Automation deleted');
            loadAutomations();
        } catch (err) {
            toast.error(err.message);
        }
    }

    if (loading) return <div className="loading-center"><div className="spinner"></div></div>;

    return (
        <div className="page-content">
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
                <div>
                    <h1 style={{ marginBottom: '4px' }}>Automations</h1>
                    <p style={{ color: 'var(--text-secondary)', fontSize: 'var(--text-sm)' }}>Auto-reply rules for incoming messages</p>
                </div>
                <button className="btn btn-primary" onClick={() => setShowModal(true)}>+ New Rule</button>
            </div>

            {automations.length === 0 ? (
                <div className="empty-state">
                    <div className="empty-state-icon">⚡</div>
                    <div className="empty-state-title">No automations yet</div>
                    <div className="empty-state-text">Create keyword-based auto-reply rules to respond instantly</div>
                </div>
            ) : (
                <div className="auto-grid">
                    {automations.map(a => (
                        <div className="auto-card" key={a.id}>
                            <div className="auto-card-header">
                                <span className="auto-card-title">{a.name}</span>
                                <div style={{ display: 'flex', gap: '6px', alignItems: 'center' }}>
                                    <span className={`badge ${a.is_active ? 'badge-success' : 'badge-danger'}`}>
                                        {a.is_active ? 'Active' : 'Off'}
                                    </span>
                                    <button className="btn btn-sm btn-danger" onClick={() => handleDelete(a.id)}>Delete</button>
                                </div>
                            </div>
                            <div className="auto-card-trigger">Trigger: {a.trigger_type}</div>
                            <div className="auto-card-keywords">
                                {(a.keywords || []).map((kw, i) => <span className="tag" key={i}>{kw}</span>)}
                            </div>
                            <div className="auto-card-reply">↩ {a.reply_text}</div>
                            {a.delay_ms > 0 && (
                                <div style={{ fontSize: 'var(--text-xs)', color: 'var(--text-muted)', marginTop: '6px' }}>
                                    ⏱ {a.delay_ms}ms delay
                                </div>
                            )}
                        </div>
                    ))}
                </div>
            )}

            {showModal && (
                <div className="modal-backdrop" onClick={() => setShowModal(false)}>
                    <div className="modal" onClick={e => e.stopPropagation()}>
                        <h2>New Automation Rule</h2>
                        <form onSubmit={handleCreate}>
                            <div className="form-group">
                                <label>Rule Name</label>
                                <input className="input" value={form.name} onChange={e => setForm({ ...form, name: e.target.value })}
                                    placeholder="e.g., Price inquiry reply" required />
                            </div>
                            <div className="form-group">
                                <label>Trigger Type</label>
                                <select className="input" value={form.trigger_type} onChange={e => setForm({ ...form, trigger_type: e.target.value })}>
                                    <option value="keyword">Keyword Match</option>
                                    <option value="first_message">First Message</option>
                                </select>
                            </div>
                            {form.trigger_type === 'keyword' && (
                                <div className="form-group">
                                    <label>Keywords (comma-separated)</label>
                                    <input className="input" value={form.keywords} onChange={e => setForm({ ...form, keywords: e.target.value })}
                                        placeholder="price, cost, how much" required />
                                </div>
                            )}
                            <div className="form-group">
                                <label>Auto-Reply Message</label>
                                <textarea className="input" rows={3} value={form.reply_text} onChange={e => setForm({ ...form, reply_text: e.target.value })}
                                    placeholder="Thanks for reaching out! Our prices start at..." required />
                            </div>
                            <div className="form-group">
                                <label>Delay (milliseconds, 0 = instant)</label>
                                <input className="input" type="number" value={form.delay_ms} onChange={e => setForm({ ...form, delay_ms: e.target.value })}
                                    placeholder="0" min="0" />
                            </div>
                            <div className="modal-actions">
                                <button type="button" className="btn" onClick={() => setShowModal(false)}>Cancel</button>
                                <button type="submit" className="btn btn-primary">Create Rule</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
}
