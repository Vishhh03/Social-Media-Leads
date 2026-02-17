import { useState, useEffect } from 'react'
import { getContacts } from '../api'

export default function Contacts() {
    const [contacts, setContacts] = useState([])
    const [loading, setLoading] = useState(true)
    const [search, setSearch] = useState('')

    useEffect(() => {
        async function load() {
            try {
                const data = await getContacts(200, 0)
                setContacts(data.contacts || [])
            } catch (e) { console.error(e) }
            finally { setLoading(false) }
        }
        load()
    }, [])

    const filtered = contacts.filter(c =>
        (c.name || '').toLowerCase().includes(search.toLowerCase()) ||
        (c.email || '').toLowerCase().includes(search.toLowerCase()) ||
        (c.platform_user_id || '').includes(search)
    )

    function formatDate(ts) {
        if (!ts) return '—'
        return new Date(ts).toLocaleDateString([], { month: 'short', day: 'numeric', year: 'numeric' })
    }

    return (
        <>
            <div className="page-header">
                <div>
                    <div className="page-title">Contacts</div>
                    <div className="page-subtitle">{contacts.length} total leads</div>
                </div>
            </div>
            <div className="page-body">
                <div style={{ marginBottom: 14 }}>
                    <input className="input" style={{ maxWidth: 320 }} placeholder="Search contacts…" value={search} onChange={e => setSearch(e.target.value)} />
                </div>

                {loading ? (
                    <div className="loading-center"><div className="spinner" /></div>
                ) : filtered.length === 0 ? (
                    <div className="empty-state">
                        <div className="empty-state-icon">👥</div>
                        <div className="empty-state-title">{search ? 'No matches found' : 'No contacts yet'}</div>
                        <div className="empty-state-text">{search ? 'Try a different search' : 'Contacts will appear when leads message you'}</div>
                    </div>
                ) : (
                    <div className="table-wrap">
                        <table className="table">
                            <thead>
                                <tr>
                                    <th>Name</th>
                                    <th>Platform</th>
                                    <th>Phone</th>
                                    <th>Email</th>
                                    <th>Status</th>
                                    <th>Added</th>
                                </tr>
                            </thead>
                            <tbody>
                                {filtered.map(c => (
                                    <tr key={c.id}>
                                        <td style={{ fontWeight: 500 }}>{c.name || c.platform_user_id}</td>
                                        <td><span className={`inbox-item-platform ${c.platform}`}>{c.platform}</span></td>
                                        <td style={{ color: c.phone ? 'inherit' : 'var(--text-muted)' }}>{c.phone || '—'}</td>
                                        <td style={{ color: c.email ? 'inherit' : 'var(--text-muted)' }}>{c.email || '—'}</td>
                                        <td>{c.is_hot_lead ? <span className="badge badge-warning">🔥 Hot</span> : <span className="badge badge-info">Lead</span>}</td>
                                        <td style={{ color: 'var(--text-muted)', fontSize: 'var(--text-xs)' }}>{formatDate(c.created_at)}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </>
    )
}
