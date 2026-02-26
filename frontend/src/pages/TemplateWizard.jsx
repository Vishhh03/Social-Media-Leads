import { useState, useEffect } from 'react';
import { activatePropertyVisit, getPropertyVisitConfig } from '../api';
import { useToast } from '../components/Toast';

export default function TemplateWizard() {
    const { addToast } = useToast();
    const [status, setStatus] = useState('idle'); // idle -> activating -> active
    const [loading, setLoading] = useState(true);
    const [config, setConfig] = useState({
        projectName: '',
        brochureUrl: '',
        agentPhone: '',
    });

    // Load existing wizard config on mount so re-visiting pre-fills the form
    useEffect(() => {
        getPropertyVisitConfig()
            .then(data => {
                if (data?.active && data?.config) {
                    setConfig({
                        projectName: data.config.project_name || '',
                        brochureUrl: data.config.brochure_url || '',
                        agentPhone: data.config.agent_phone || '',
                    });
                    setStatus('active');
                }
            })
            .catch(() => { /* No config yet — stay in idle */ })
            .finally(() => setLoading(false));
    }, []);

    const handleActivate = async () => {
        if (!config.projectName || !config.agentPhone) {
            addToast('error', 'Project Name and Agent WhatsApp Number are required.');
            return;
        }

        setStatus('activating');
        try {
            await activatePropertyVisit(config.projectName, config.brochureUrl, config.agentPhone);
            setStatus('active');
            addToast('success', '🚀 Automation is live for ' + config.projectName + '!');
        } catch (err) {
            addToast('error', err.message || 'Failed to activate. Please try again.');
            setStatus('idle');
        }
    };

    if (loading) {
        return <div className="loading-center"><div className="spinner" /></div>;
    }

    return (
        <>
            <div className="page-header">
                <div>
                    <div className="page-title">Property Visit Automation</div>
                    <div className="page-subtitle">Your AI Site-Visit Setter — configure once, books visits forever.</div>
                </div>
            </div>

            <div className="page-body">
                <div style={{ maxWidth: 800, margin: '0 auto' }}>
                    {status === 'active' ? (
                        <div className="card" style={{ textAlign: 'center', padding: '48px 24px', boxShadow: 'var(--shadow-sm)' }}>
                            <div style={{ width: 80, height: 80, borderRadius: '50%', background: 'var(--accent-light)', color: 'var(--accent)', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 24px', fontSize: 40 }}>
                                🎉
                            </div>
                            <h2 style={{ fontSize: 26, fontWeight: 800, marginBottom: 12 }}>
                                {config.projectName} is Live!
                            </h2>
                            <p style={{ color: 'var(--text-secondary)', fontSize: 16, maxWidth: 500, margin: '0 auto 32px', lineHeight: 1.6 }}>
                                Your AI assistant is now qualifying leads and booking visits on Instagram and WhatsApp automatically.
                            </p>

                            {/* Activation Simulator */}
                            <div style={{ padding: 24, background: 'var(--bg-secondary)', borderRadius: 12, border: '1px dashed var(--border)', display: 'inline-block', textAlign: 'left', maxWidth: 440 }}>
                                <p style={{ fontWeight: 700, marginBottom: 8, fontSize: 15 }}>🧪 Test it right now</p>
                                <p style={{ color: 'var(--text-secondary)', marginBottom: 16, fontSize: 14, lineHeight: 1.5 }}>
                                    Message your connected account:<br />
                                    <strong style={{ color: 'var(--text-primary)' }}>
                                        "Hi, I want to know more about {config.projectName}"
                                    </strong>
                                </p>
                                <div style={{ display: 'flex', gap: 12 }}>
                                    <a href="/dashboard/inbox" className="btn btn-secondary" style={{ textDecoration: 'none', fontSize: 14 }}>
                                        Watch in Inbox →
                                    </a>
                                    <button
                                        className="btn btn-ghost"
                                        onClick={() => setStatus('idle')}
                                        style={{ fontSize: 14 }}
                                    >
                                        Edit Configuration
                                    </button>
                                </div>
                            </div>
                        </div>
                    ) : (
                        <div className="card" style={{ boxShadow: 'var(--shadow-sm)' }}>
                            <div style={{ marginBottom: 32 }}>
                                <h3 style={{ fontWeight: 700, fontSize: 18, marginBottom: 8 }}>Setup Your Property Visit Bot</h3>
                                <p style={{ color: 'var(--text-secondary)', fontSize: 14, lineHeight: 1.6 }}>
                                    The conversational flow is pre-built. Just tell us about your project and where to send booking alerts.
                                </p>
                            </div>

                            <div style={{ display: 'flex', flexDirection: 'column', gap: 24 }}>
                                <div className="form-group">
                                    <label className="form-label" style={{ fontWeight: 600 }}>Project Name *</label>
                                    <input
                                        type="text"
                                        className="form-input"
                                        value={config.projectName}
                                        onChange={e => setConfig({ ...config, projectName: e.target.value })}
                                        placeholder="e.g. Prestige Lakeside Habitat"
                                    />
                                </div>

                                <div className="form-group">
                                    <label className="form-label" style={{ fontWeight: 600 }}>Project Brochure Link <span style={{ color: 'var(--text-muted)', fontWeight: 400 }}>(optional)</span></label>
                                    <input
                                        type="text"
                                        className="form-input"
                                        value={config.brochureUrl}
                                        onChange={e => setConfig({ ...config, brochureUrl: e.target.value })}
                                        placeholder="https://example.com/brochure.pdf"
                                    />
                                    <p style={{ fontSize: 12, color: 'var(--text-muted)', marginTop: 6 }}>
                                        The bot will send this to qualified leads automatically.
                                    </p>
                                </div>

                                <div className="form-group">
                                    <label className="form-label" style={{ fontWeight: 600 }}>Available Visit Slots</label>
                                    <div style={{ display: 'flex', gap: 10, flexWrap: 'wrap' }}>
                                        {['10:00 AM', '2:00 PM', '4:00 PM'].map(slot => (
                                            <div key={slot} style={{ padding: '8px 16px', borderRadius: 20, background: 'var(--accent-light)', color: 'var(--accent)', fontWeight: 600, fontSize: 13 }}>
                                                {slot}
                                            </div>
                                        ))}
                                        <div style={{ padding: '8px 16px', borderRadius: 20, background: 'var(--bg-secondary)', color: 'var(--text-muted)', fontSize: 13, border: '1px dashed var(--border)' }}>
                                            + Custom (Pro)
                                        </div>
                                    </div>
                                    <p style={{ fontSize: 12, color: 'var(--text-muted)', marginTop: 6 }}>
                                        Double-booking is prevented automatically.
                                    </p>
                                </div>

                                <div className="form-group" style={{ borderTop: '1px solid var(--border)', paddingTop: 24 }}>
                                    <label className="form-label" style={{ fontWeight: 600 }}>Your WhatsApp Number (for booking alerts) *</label>
                                    <input
                                        type="text"
                                        className="form-input"
                                        value={config.agentPhone}
                                        onChange={e => setConfig({ ...config, agentPhone: e.target.value })}
                                        placeholder="+91 98765 43210"
                                    />
                                    <p style={{ fontSize: 12, color: 'var(--text-muted)', marginTop: 6 }}>
                                        You'll get an instant structured alert when a lead confirms a visit time.
                                    </p>
                                </div>

                                <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: 8 }}>
                                    <button
                                        className="btn btn-primary"
                                        onClick={handleActivate}
                                        disabled={status === 'activating'}
                                        style={{ padding: '12px 32px', fontSize: 16, fontWeight: 700 }}
                                    >
                                        {status === 'activating' ? 'Activating...' : 'Activate Automation 🚀'}
                                    </button>
                                </div>
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </>
    );
}
