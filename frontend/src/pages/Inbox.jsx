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

    const platformBadge = (platform) => (
        <span className={`inbox-item-platform ${platform}`}>{platform}</span>
    );

    if (loading) return <div className="loading-center"><div className="spinner"></div></div>;

    return (
        <div className="inbox-layout">
            {/* Conversation List */}
            <div className="inbox-list">
                <div className="inbox-list-header">
                    Conversations
                    <span style={{ float: 'right', fontSize: 'var(--text-xs)', color: 'var(--text-muted)' }}>{conversations.length}</span>
                </div>
                <div className="inbox-list-items">
                    {conversations.length === 0 ? (
                        <div style={{ padding: '40px 20px', textAlign: 'center', color: 'var(--text-muted)', fontSize: 'var(--text-sm)' }}>
                            No conversations yet
                        </div>
                    ) : conversations.map(c => (
                        <div
                            key={c.id}
                            className={`inbox-item ${selected?.id === c.id ? 'active' : ''}`}
                            onClick={() => setSelected(c)}
                        >
                            <div className="inbox-item-avatar">{(c.name || '?')[0].toUpperCase()}</div>
                            <div className="inbox-item-info">
                                <div className="inbox-item-name">{c.name || 'Unknown'}</div>
                                <div className="inbox-item-msg">{platformBadge(c.platform)}</div>
                            </div>
                            <div className="inbox-item-time">
                                {c.updated_at ? new Date(c.updated_at).toLocaleDateString() : ''}
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            {/* Chat Thread */}
            {selected ? (
                <div className="inbox-chat">
                    <div className="inbox-chat-header">
                        <div className="inbox-item-avatar" style={{ width: 30, height: 30, fontSize: '12px' }}>
                            {(selected.name || '?')[0].toUpperCase()}
                        </div>
                        <div>
                            <div style={{ fontWeight: 600 }}>{selected.name || 'Unknown'}</div>
                            <div style={{ fontSize: 'var(--text-xs)', color: 'var(--text-muted)' }}>{selected.platform} · {selected.platform_user_id}</div>
                        </div>
                        {selected.is_hot_lead && <span className="badge badge-danger" style={{ marginLeft: 'auto' }}>🔥 Hot Lead</span>}
                    </div>
                    <div className="inbox-chat-messages">
                        {messages.map(m => (
                            <div key={m.id} className={`chat-msg ${m.direction}`}>
                                <div>{m.content}</div>
                                <div className="chat-msg-time">
                                    {new Date(m.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                                    {m.is_automated && <span className="chat-msg-auto"> · Bot</span>}
                                </div>
                            </div>
                        ))}
                        <div ref={messagesEnd} />
                    </div>
                    <form className="inbox-chat-input" onSubmit={handleSend}>
                        <input
                            className="input"
                            value={newMsg}
                            onChange={e => setNewMsg(e.target.value)}
                            placeholder="Type a message..."
                            disabled={sending}
                        />
                        <button className="btn btn-primary" type="submit" disabled={sending || !newMsg.trim()}>
                            {sending ? '...' : 'Send'}
                        </button>
                    </form>
                </div>
            ) : (
                <div className="inbox-empty">
                    <div>
                        <div style={{ fontSize: '48px', marginBottom: '12px' }}>💬</div>
                        <div>Select a conversation to start messaging</div>
                    </div>
                </div>
            )}
        </div>
    );
}
