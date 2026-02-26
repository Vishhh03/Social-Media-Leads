import { useState, useEffect, useRef } from 'react';
import { getConversations, getMessages, sendMessage } from '../api';
import { useToast } from '../components/Toast';

export default function Inbox() {
    const toast = useToast();
    const [conversations, setConversations] = useState([]);
    const [selected, setSelected] = useState(null);
    const [messages, setMessages] = useState([]);
    const [newMsg, setNewMsg] = useState('');
    const [loading, setLoading] = useState(true);
    const [sending, setSending] = useState(false);
    const messagesEnd = useRef(null);

    useEffect(() => { loadConversations(); }, []);
    useEffect(() => { if (selected) loadMessages(selected.id); }, [selected]);
    useEffect(() => { messagesEnd.current?.scrollIntoView({ behavior: 'smooth' }); }, [messages]);

    async function loadConversations() {
        try {
            const data = await getConversations();
            setConversations(data.conversations || []);
        } catch (err) {
            toast.error(err.message);
        } finally {
            setLoading(false);
        }
    }

    async function loadMessages(contactId) {
        try {
            const data = await getMessages(contactId);
            setMessages(data.messages || []);
        } catch (err) {
            toast.error(err.message);
        }
    }

    async function handleSend(e) {
        e.preventDefault();
        if (!newMsg.trim() || !selected) return;
        setSending(true);
        try {
            await sendMessage(selected.id, newMsg.trim());
            setNewMsg('');
            loadMessages(selected.id);
        } catch (err) {
            toast.error(err.message);
        } finally {
            setSending(false);
        }
    }

    const platformBadge = (platform) => {
        const colors = {
            whatsapp: { bg: '#dcfce7', text: '#15803d' },
            instagram: { bg: '#fce7f3', text: '#be185d' },
            facebook: { bg: '#dbeafe', text: '#1d4ed8' }
        };
        const color = colors[platform.toLowerCase()] || { bg: 'var(--bg-secondary)', text: 'var(--text-secondary)' };
        return (
            <span style={{
                fontSize: '9px', padding: '2px 6px', borderRadius: '10px',
                fontWeight: 700, textTransform: 'uppercase', letterSpacing: '0.04em',
                backgroundColor: color.bg, color: color.text,
                display: 'inline-block', marginTop: '4px'
            }}>
                {platform}
            </span>
        );
    };

    if (loading) return <div className="loading-center"><div className="spinner"></div></div>;

    return (
        <div className="inbox-layout" style={{ background: 'var(--bg-primary)' }}>
            {/* Conversation List */}
            <div className="inbox-list" style={{ boxShadow: 'var(--shadow-sm)', zIndex: 10 }}>
                <div className="inbox-list-header" style={{ padding: '20px 16px', borderBottom: '1px solid var(--border)' }}>
                    <div style={{ fontSize: 'var(--text-lg)', fontWeight: 700 }}>Messages</div>
                    <div style={{ fontSize: 'var(--text-xs)', color: 'var(--text-muted)', marginTop: 4 }}>{conversations.length} active threads</div>
                </div>
                <div className="inbox-list-items">
                    {conversations.length === 0 ? (
                        <div style={{ padding: '60px 20px', textAlign: 'center' }}>
                            <div style={{ fontSize: '32px', marginBottom: 16 }}>🌙</div>
                            <div style={{ color: 'var(--text-muted)', fontSize: 'var(--text-sm)' }}>No messages yet</div>
                        </div>
                    ) : conversations.map(c => (
                        <div
                            key={c.id}
                            className={`inbox-item ${selected?.id === c.id ? 'active' : ''}`}
                            onClick={() => setSelected(c)}
                            style={{
                                transition: 'all 0.15s ease',
                                borderLeft: selected?.id === c.id ? '4px solid var(--accent)' : '4px solid transparent'
                            }}
                        >
                            <div className="inbox-item-avatar" style={{
                                background: selected?.id === c.id ? 'var(--accent)' : 'var(--accent-light)',
                                color: selected?.id === c.id ? 'white' : 'var(--accent)',
                                fontSize: '14px', fontWeight: 700
                            }}>
                                {(c.name || '?')[0].toUpperCase()}
                            </div>
                            <div className="inbox-item-info">
                                <div className="inbox-item-name" style={{ fontWeight: 650 }}>{c.name || 'Unknown'}</div>
                                {platformBadge(c.platform)}
                            </div>
                            <div className="inbox-item-time" style={{ fontSize: '10px', color: 'var(--text-muted)' }}>
                                {c.updated_at ? new Date(c.updated_at).toLocaleDateString([], { month: 'short', day: 'numeric' }) : ''}
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            {/* Chat Thread */}
            {selected ? (
                <div className="inbox-chat">
                    <div className="inbox-chat-header" style={{ height: '70px', padding: '0 24px', boxShadow: 'var(--shadow-sm)' }}>
                        <div className="inbox-item-avatar" style={{ width: 36, height: 36, fontSize: '14px', background: 'var(--accent-light)', color: 'var(--accent)' }}>
                            {(selected.name || '?')[0].toUpperCase()}
                        </div>
                        <div style={{ flex: 1 }}>
                            <div style={{ fontWeight: 700, fontSize: '16px' }}>{selected.name || 'Unknown'}</div>
                            <div style={{ fontSize: 'var(--text-xs)', color: 'var(--text-muted)', display: 'flex', alignItems: 'center', gap: 6 }}>
                                <span className="platform-dot" style={{ width: 6, height: 6, borderRadius: '50%', background: '#10b981' }}></span>
                                {selected.platform} · ID: {selected.platform_user_id?.slice(-8)}
                            </div>
                        </div>
                        {selected.is_hot_lead && (
                            <div style={{
                                backgroundColor: 'var(--danger-light)', color: 'var(--danger)',
                                padding: '4px 12px', borderRadius: '50px', fontSize: '11px', fontWeight: 700,
                                border: '1px solid rgba(220, 38, 38, 0.1)'
                            }}>
                                🔥 HOT LEAD
                            </div>
                        )}
                    </div>
                    <div className="inbox-chat-messages" style={{ padding: '24px', gap: '12px', background: 'var(--bg-primary)' }}>
                        {messages.length === 0 ? (
                            <div style={{ flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--text-muted)', fontSize: 'var(--text-sm)' }}>
                                Beginning of conversation with {selected.name}
                            </div>
                        ) : messages.map(m => (
                            <div key={m.id} className={`chat-msg ${m.direction}`} style={{
                                padding: '12px 16px',
                                borderRadius: m.direction === 'inbound' ? '16px 16px 16px 4px' : '16px 16px 4px 16px',
                                boxShadow: 'var(--shadow-xs)',
                                fontSize: 'var(--text-base)',
                                border: m.direction === 'inbound' ? '1px solid var(--border)' : 'none'
                            }}>
                                <div style={{ marginBottom: 4 }}>{m.content}</div>
                                <div className="chat-msg-time" style={{ fontSize: '10px', display: 'flex', alignItems: 'center', gap: 4 }}>
                                    {new Date(m.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                                    {m.is_automated && <span style={{
                                        padding: '1px 4px', borderRadius: '4px', background: 'rgba(0,0,0,0.05)',
                                        fontSize: '9px', fontWeight: 600, textTransform: 'uppercase'
                                    }}>AI Bot</span>}
                                </div>
                            </div>
                        ))}
                        <div ref={messagesEnd} />
                    </div>
                    <div style={{ padding: '20px 24px', background: 'var(--white)', borderTop: '1px solid var(--border)' }}>
                        <form className="inbox-chat-input" onSubmit={handleSend} style={{ background: 'var(--bg-secondary)', borderRadius: '12px', padding: '4px 4px 4px 16px', border: '1px solid var(--border)' }}>
                            <input
                                className="input"
                                style={{ background: 'transparent', border: 'none', boxShadow: 'none' }}
                                value={newMsg}
                                onChange={e => setNewMsg(e.target.value)}
                                placeholder={`Reply as LeadPilot...`}
                                disabled={sending}
                            />
                            <button className="btn btn-primary" type="submit" disabled={sending || !newMsg.trim()} style={{ borderRadius: '10px', padding: '8px 20px' }}>
                                {sending ? '...' : (
                                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5"><line x1="22" y1="2" x2="11" y2="13" /><polyline points="22 2 15 22 11 13 2 9 22 2" /></svg>
                                )}
                            </button>
                        </form>
                    </div>
                </div>
            ) : (
                <div className="inbox-empty" style={{ background: 'var(--bg-primary)' }}>
                    <div style={{ textAlign: 'center' }}>
                        <div style={{
                            width: 80, height: 80, background: 'var(--white)', borderRadius: '24px',
                            display: 'flex', alignItems: 'center', justifyContent: 'center',
                            fontSize: '32px', margin: '0 auto 24px auto', boxShadow: 'var(--shadow-md)'
                        }}>💬</div>
                        <h2 style={{ fontSize: '20px', fontWeight: 700, color: 'var(--text-primary)', marginBottom: 8 }}>Select a conversation</h2>
                        <p style={{ color: 'var(--text-secondary)', fontSize: '14px' }}>Choose a thread from the list to view the message history<br />and reply in real-time.</p>
                    </div>
                </div>
            )}
        </div>
    );
}
