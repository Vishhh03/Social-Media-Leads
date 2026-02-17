import { useState, useEffect } from 'react'
import { getContacts, getConversations, getAutomations, getChannels } from '../api'

export default function Dashboard() {
    const [stats, setStats] = useState({ contacts: 0, conversations: 0, automations: 0, channels: 0 })
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        async function load() {
            try {
                const [c, conv, a, ch] = await Promise.all([
                    getContacts(1, 0),
                    getConversations(1, 0),
                    getAutomations(),
                    getChannels(),
                ])
                setStats({
                    contacts: c.count || 0,
                    conversations: conv.count || 0,
                    automations: a.count || 0,
                    channels: ch.count || 0,
                })
            } catch (e) {
                console.error(e)
            } finally {
                setLoading(false)
            }
        }
        load()
    }, [])

    return (
        <>
            <div className="page-header">
                <div>
                    <div className="page-title">Dashboard</div>
                    <div className="page-subtitle">Overview of your lead automation</div>
                </div>
            </div>
            <div className="page-body">
                {loading ? (
                    <div className="loading-center"><div className="spinner" /></div>
                ) : (
                    <div className="stats-grid">
                        <div className="stat-card">
                            <div className="stat-label">Total Contacts</div>
                            <div className="stat-value">{stats.contacts}</div>
                        </div>
                        <div className="stat-card">
                            <div className="stat-label">Conversations</div>
                            <div className="stat-value">{stats.conversations}</div>
                        </div>
                        <div className="stat-card">
                            <div className="stat-label">Active Automations</div>
                            <div className="stat-value">{stats.automations}</div>
                        </div>
                        <div className="stat-card">
                            <div className="stat-label">Connected Channels</div>
                            <div className="stat-value">{stats.channels}</div>
                        </div>
                    </div>
                )}

                <div className="card" style={{ marginTop: 8 }}>
                    <h3 style={{ fontWeight: 600, marginBottom: 12 }}>Quick Start</h3>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: 10, fontSize: 'var(--text-sm)', color: 'var(--text-secondary)' }}>
                        <div>1. <strong>Connect a channel</strong> — Go to Channels and link your WhatsApp, Instagram, or Facebook account.</div>
                        <div>2. <strong>Set up automations</strong> — Create keyword-based auto-reply rules in Automations.</div>
                        <div>3. <strong>Manage leads</strong> — Incoming messages will appear in your Inbox. Reply directly from there.</div>
                        <div>4. <strong>Send broadcasts</strong> — Blast messages to all your contacts at once.</div>
                    </div>
                </div>
            </div>
        </>
    )
}
