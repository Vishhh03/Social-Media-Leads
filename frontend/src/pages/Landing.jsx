import React from 'react';
import { Link } from 'react-router-dom';
import { useToast } from '../components/Toast';

export default function Landing() {
    const toast = useToast();

    return (
        <div className="landing-page" style={{ backgroundColor: 'var(--bg-primary)', minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
            {/* Top Navigation */}
            <nav style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '16px 40px', backgroundColor: 'rgba(255, 255, 255, 0.8)', backdropFilter: 'blur(10px)', borderBottom: '1px solid var(--border)', position: 'sticky', top: 0, zIndex: 50 }}>
                <div style={{ fontSize: '20px', fontWeight: 'bold', display: 'flex', alignItems: 'center', gap: '8px', color: 'var(--text-primary)' }}>
                    <div style={{ width: 24, height: 24, background: 'var(--grad-primary)', borderRadius: 6, display: 'inline-block' }}></div>
                    Lead<span style={{ color: 'var(--accent)' }}>Pilot</span>
                </div>
                <div style={{ display: 'flex', gap: '24px', alignItems: 'center' }}>
                    <Link to="/login" style={{ color: 'var(--text-secondary)', fontWeight: 500, fontSize: '14px', textDecoration: 'none' }}>Sign In</Link>
                    <Link to="/signup" className="btn btn-primary" style={{ padding: '8px 20px', borderRadius: '8px', boxShadow: 'var(--shadow-sm)' }}>Get Started</Link>
                </div>
            </nav>

            {/* Hero Section */}
            <header style={{
                textAlign: 'center',
                padding: '120px 24px',
                flex: 1,
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                background: 'var(--grad-hero)',
                position: 'relative',
                overflow: 'hidden'
            }}>
                <div className="animate-up" style={{ maxWidth: '900px', margin: '0 auto', zIndex: 2 }}>
                    <div style={{
                        display: 'inline-block', padding: '6px 16px', borderRadius: '50px', fontSize: '13px', fontWeight: 600,
                        backgroundColor: 'var(--accent-light)', color: 'var(--accent)', marginBottom: '24px',
                        border: '1px solid rgba(37, 99, 235, 0.1)'
                    }}>
                        Automate your growth.
                    </div>
                    <h1 style={{ fontSize: '64px', fontWeight: 850, letterSpacing: '-0.04em', lineHeight: 1.1, marginBottom: '24px', color: 'var(--text-primary)' }}>
                        Turn conversations into <br />
                        <span style={{
                            background: 'var(--grad-primary)',
                            WebkitBackgroundClip: 'text',
                            WebkitTextFillColor: 'transparent'
                        }}>customers on autopilot.</span>
                    </h1>
                    <p style={{ fontSize: '20px', color: 'var(--text-secondary)', marginBottom: '48px', lineHeight: 1.6, maxWidth: '650px', margin: '0 auto 48px auto' }}>
                        LeadPilot unifies your social channels into one clean inbox, and uses intelligent workflows to automatically qualify leads 24/7.
                    </p>
                    <div style={{ display: 'flex', gap: '16px', justifyContent: 'center' }}>
                        <Link to="/signup" className="btn btn-primary btn-lg" style={{
                            padding: '16px 32px', borderRadius: '12px', fontWeight: 600,
                            boxShadow: '0 10px 15px -3px rgba(37, 99, 235, 0.3)',
                            fontSize: '16px'
                        }}>
                            Start Building Free
                        </Link>
                        <button onClick={() => toast.info('Book a demo coming soon')} className="btn btn-lg" style={{
                            padding: '16px 32px', borderRadius: '12px', fontWeight: 600,
                            backgroundColor: 'white', border: '1px solid var(--border)',
                            color: 'var(--text-primary)', fontSize: '16px',
                            boxShadow: 'var(--shadow-sm)'
                        }}>
                            Book a Demo
                        </button>
                    </div>
                </div>

                {/* Dashboard Mockup Frame */}
                <div className="animate-up delay-2" style={{
                    marginTop: '100px', maxWidth: '1080px', width: '100%',
                    border: '1px solid var(--border)', borderRadius: '24px', padding: '12px',
                    backgroundColor: 'rgba(255, 255, 255, 0.4)',
                    backdropFilter: 'blur(8px)',
                    boxShadow: 'var(--shadow-xl)',
                    position: 'relative'
                }}>
                    <div style={{ width: '100%', height: '480px', backgroundColor: 'var(--bg-primary)', borderRadius: '16px', border: '1px solid var(--border-light)', position: 'relative', overflow: 'hidden', display: 'flex' }}>
                        {/* Mock Sidebar */}
                        <div style={{ width: '220px', borderRight: '1px solid var(--border)', backgroundColor: 'var(--white)', padding: '24px 12px' }}>
                            <div style={{ height: 14, width: 80, backgroundColor: 'var(--border-light)', borderRadius: 4, marginBottom: 40, marginLeft: 10 }}></div>
                            <div style={{ height: 32, width: '100%', backgroundColor: 'var(--accent-light)', borderRadius: 8, marginBottom: 12 }}></div>
                            <div style={{ height: 32, width: '100%', backgroundColor: 'transparent', borderRadius: 8, marginBottom: 12 }}></div>
                            <div style={{ height: 32, width: '100%', backgroundColor: 'transparent', borderRadius: 8, marginBottom: 12 }}></div>
                        </div>
                        {/* Mock Content */}
                        <div style={{ flex: 1, padding: '40px', display: 'flex', flexDirection: 'column', gap: '32px' }}>
                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                <div style={{ height: 28, width: 140, backgroundColor: 'var(--border-light)', borderRadius: 6 }}></div>
                                <div style={{ height: 36, width: 100, backgroundColor: 'var(--accent-light)', borderRadius: 8 }}></div>
                            </div>
                            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 24 }}>
                                {[1, 2, 3].map(i => (
                                    <div key={i} style={{ height: 140, backgroundColor: 'var(--white)', border: '1px solid var(--border)', borderRadius: 12, padding: 24, boxShadow: 'var(--shadow-xs)' }}>
                                        <div style={{ height: 12, width: 80, backgroundColor: 'var(--border-light)', borderRadius: 4, marginBottom: 12 }}></div>
                                        <div style={{ height: 32, width: 50, backgroundColor: 'var(--text-primary)', borderRadius: 6 }}></div>
                                    </div>
                                ))}
                            </div>
                            <div style={{ flex: 1, backgroundColor: 'var(--white)', border: '1px solid var(--border)', borderRadius: 12, padding: 24, boxShadow: 'var(--shadow-xs)' }}>
                                <div style={{ height: 14, width: 120, backgroundColor: 'var(--border-light)', borderRadius: 4, marginBottom: 20 }}></div>
                                <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
                                    {[1, 2, 3].map(i => (
                                        <div key={i} style={{ height: 40, borderBottom: '1px solid var(--border-light)', display: 'flex', alignItems: 'center', gap: 12 }}>
                                            <div style={{ width: 24, height: 24, borderRadius: '50%', backgroundColor: 'var(--bg-secondary)' }}></div>
                                            <div style={{ height: 10, width: 100, backgroundColor: 'var(--border-light)', borderRadius: 4 }}></div>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </header>

            {/* Features Section */}
            <section style={{ padding: '100px 24px', backgroundColor: 'var(--white)', borderTop: '1px solid var(--border)' }}>
                <div style={{ maxWidth: '1200px', margin: '0 auto' }}>
                    <div style={{ textAlign: 'center', marginBottom: '80px' }}>
                        <h2 style={{ fontSize: '40px', fontWeight: 800, marginBottom: '16px', color: 'var(--text-primary)', letterSpacing: '-0.02em' }}>Everything you need to grow</h2>
                        <p style={{ fontSize: '18px', color: 'var(--text-secondary)', maxWidth: '600px', margin: '0 auto' }}>A unified platform to manage conversations and automate your sales pipeline.</p>
                    </div>

                    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))', gap: '32px' }}>
                        {/* Card 1 */}
                        <div className="card" style={{ padding: '40px', textAlign: 'left', border: '1px solid var(--border)', boxShadow: 'var(--shadow-sm)', transition: 'transform 0.2s', cursor: 'default' }}>
                            <div style={{ width: 48, height: 48, borderRadius: 12, background: 'var(--grad-primary)', color: 'white', display: 'flex', alignItems: 'center', justifyContent: 'center', marginBottom: 24, boxShadow: '0 4px 12px rgba(37, 99, 235, 0.2)' }}>
                                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" /></svg>
                            </div>
                            <h3 style={{ fontSize: '20px', fontWeight: 700, marginBottom: '12px', color: 'var(--text-primary)' }}>Unified Inbox</h3>
                            <p style={{ color: 'var(--text-secondary)', lineHeight: 1.6, fontSize: '15px' }}>Manage Facebook, Instagram, and WhatsApp messages from a single, clean interface.</p>
                        </div>
                        {/* Card 2 */}
                        <div className="card" style={{ padding: '40px', textAlign: 'left', border: '1px solid var(--border)', boxShadow: 'var(--shadow-sm)', transition: 'transform 0.2s', cursor: 'default' }}>
                            <div style={{ width: 48, height: 48, borderRadius: 12, background: 'linear-gradient(135deg, #059669 0%, #10b981 100%)', color: 'white', display: 'flex', alignItems: 'center', justifyContent: 'center', marginBottom: 24, boxShadow: '0 4px 12px rgba(5, 150, 105, 0.2)' }}>
                                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M4 7V4h16v3H4z" /><path d="M9 11v10H4V11h5z" /><path d="M20 11v10h-5V11h5z" /><path d="M9 16h6" /></svg>
                            </div>
                            <h3 style={{ fontSize: '20px', fontWeight: 700, marginBottom: '12px', color: 'var(--text-primary)' }}>Visual Workflows</h3>
                            <p style={{ color: 'var(--text-secondary)', lineHeight: 1.6, fontSize: '15px' }}>Build complex automation rules and time delays using our simple drag-and-drop builder.</p>
                        </div>
                        {/* Card 3 */}
                        <div className="card" style={{ padding: '40px', textAlign: 'left', border: '1px solid var(--border)', boxShadow: 'var(--shadow-sm)', transition: 'transform 0.2s', cursor: 'default' }}>
                            <div style={{ width: 48, height: 48, borderRadius: 12, background: 'linear-gradient(135deg, #d97706 0%, #f59e0b 100%)', color: 'white', display: 'flex', alignItems: 'center', justifyContent: 'center', marginBottom: 24, boxShadow: '0 4px 12px rgba(217, 119, 6, 0.2)' }}>
                                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6" /></svg>
                            </div>
                            <h3 style={{ fontSize: '20px', fontWeight: 700, marginBottom: '12px', color: 'var(--text-primary)' }}>Increase Conversion</h3>
                            <p style={{ color: 'var(--text-secondary)', lineHeight: 1.6, fontSize: '15px' }}>Never miss a lead. Instantly reply to inquiries and follow up consistently to maximize sales.</p>
                        </div>
                    </div>
                </div>
            </section>

            {/* Footer */}
            <footer style={{ padding: '48px 24px', borderTop: '1px solid var(--border)', backgroundColor: 'var(--bg-primary)', textAlign: 'center' }}>
                <div style={{ fontWeight: 'bold', fontSize: '20px', marginBottom: '12px', color: 'var(--text-primary)' }}>
                    Lead<span style={{ color: 'var(--accent)' }}>Pilot</span>
                </div>
                <p style={{ color: 'var(--text-secondary)', fontSize: '14px' }}>&copy; {new Date().getFullYear()} LeadPilot. All rights reserved.</p>
            </footer>
        </div>
    );
}
