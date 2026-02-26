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
                    <>
                        <div className="stats-grid">
                            <div className="stat-card" style={{ borderTop: '4px solid var(--accent)', boxShadow: 'var(--shadow-sm)' }}>
                                <div className="stat-label">Visits Booked This Week</div>
                                <div className="stat-value">0</div>
                                <div className="stat-change">↑ 0% this week</div>
                            </div>
                            <div className="stat-card" style={{ borderTop: '4px solid #10b981', boxShadow: 'var(--shadow-sm)' }}>
                                <div className="stat-label">Upcoming Visits</div>
                                <div className="stat-value">0</div>
                                <div className="stat-change">Syncing to calendar</div>
                            </div>
                            <div className="stat-card" style={{ borderTop: '4px solid #f59e0b', boxShadow: 'var(--shadow-sm)' }}>
                                <div className="stat-label">Conversion Rate</div>
                                <div className="stat-value">
                                    {stats.conversations > 0 ? ((0 / stats.conversations) * 100).toFixed(1) : '0.0'}%
                                </div>
                                <div className="stat-change">Visits / Conversations</div>
                            </div>
                            <div className="stat-card" style={{ borderTop: '4px solid #6366f1', boxShadow: 'var(--shadow-sm)' }}>
                                <div className="stat-label">Conversations Started</div>
                                <div className="stat-value">{stats.conversations}</div>
                                <div className="stat-change">Across IG/WA</div>
                            </div>
                        </div>

                        <div style={{ display: 'grid', gridTemplateColumns: '1.5fr 1fr', gap: 24 }}>
                            <div className="card" style={{ boxShadow: 'var(--shadow-sm)', border: '1px solid var(--accent)' }}>
                                <h3 style={{ fontWeight: 800, marginBottom: 16, display: 'flex', alignItems: 'center', gap: 10, color: 'var(--accent)' }}>
                                    <div style={{ width: 32, height: 32, borderRadius: 8, background: 'var(--accent-light)', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>🚀</div>
                                    Launch Property Visit Booking
                                </h3>
                                <div style={{ fontSize: '15px', color: 'var(--text-secondary)', lineHeight: 1.6, marginBottom: 24 }}>
                                    Activate your AI Site-Visit Setter in under 2 minutes. We will automatically qualify leads on Instagram and WhatsApp, and book site visits directly for your projects.
                                </div>
                                <a href="/template-wizard" className="btn btn-primary" style={{ display: 'inline-flex', alignItems: 'center', gap: 8, textDecoration: 'none', padding: '12px 24px', fontWeight: 600, fontSize: '15px' }}>
                                    Activate Automation
                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M5 12h14M12 5l7 7-7 7" /></svg>
                                </a>
                            </div>

                            <div className="card" style={{ boxShadow: 'var(--shadow-sm)', backgroundColor: 'var(--bg-secondary)', borderStyle: 'dashed' }}>
                                <div style={{ height: '100%', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', textAlign: 'center', padding: 20 }}>
                                    <div style={{ width: 40, height: 40, borderRadius: '50%', background: 'var(--white)', display: 'flex', alignItems: 'center', justifyContent: 'center', marginBottom: 12, color: 'var(--text-muted)' }}>
                                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" /><polyline points="22 4 12 14.01 9 11.01" /></svg>
                                    </div>
                                    <div style={{ fontWeight: 600, marginBottom: 4 }}>System Status</div>
                                    <div style={{ fontSize: 'var(--text-xs)', color: 'var(--text-secondary)' }}>All systems operational. No issues detected in the last 24h.</div>
                                </div>
                            </div>
                        </div>
                    </>
                )}
            </div>
        </>
    )
}
