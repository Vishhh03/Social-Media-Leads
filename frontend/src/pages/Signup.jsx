import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { signup } from '../api';
import { useToast } from '../components/Toast';

export default function Signup() {
    const navigate = useNavigate();
    const toast = useToast();
    const [fullName, setFullName] = useState('');
    const [company, setCompany] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [loading, setLoading] = useState(false);

    async function handleSubmit(e) {
        e.preventDefault();
        setLoading(true);
        try {
            await signup(email, password, fullName, company);
            navigate('/dashboard');
        } catch (err) {
            toast.error(err.message);
        } finally {
            setLoading(false);
        }
    }

    function handleGoogleLogin() {
        window.location.href = '/api/v1/auth/google';
    }

    return (
        <div className="auth-page">
            <div className="auth-box">
                <h1>Get started</h1>
                <p>Create your LeadPilot account</p>

                <button
                    className="btn"
                    onClick={handleGoogleLogin}
                    style={{
                        width: '100%', justifyContent: 'center', padding: '11px',
                        marginBottom: '20px', gap: '10px', fontSize: 'var(--text-sm)',
                        border: '1px solid var(--border)', fontWeight: 500,
                    }}
                >
                    <svg width="18" height="18" viewBox="0 0 24 24">
                        <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 01-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" fill="#4285F4" />
                        <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853" />
                        <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05" />
                        <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335" />
                    </svg>
                    Continue with Google
                </button>

                <div style={{
                    display: 'flex', alignItems: 'center', gap: '12px',
                    marginBottom: '20px', color: 'var(--text-muted)', fontSize: 'var(--text-xs)',
                }}>
                    <hr style={{ flex: 1, border: 'none', borderTop: '1px solid var(--border)' }} />
                    OR
                    <hr style={{ flex: 1, border: 'none', borderTop: '1px solid var(--border)' }} />
                </div>

                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label>Full Name</label>
                        <input className="input" type="text" value={fullName} onChange={e => setFullName(e.target.value)}
                            placeholder="John Doe" required />
                    </div>
                    <div className="form-group">
                        <label>Company</label>
                        <input className="input" type="text" value={company} onChange={e => setCompany(e.target.value)}
                            placeholder="Acme Inc. (optional)" />
                    </div>
                    <div className="form-group">
                        <label>Email</label>
                        <input className="input" type="email" value={email} onChange={e => setEmail(e.target.value)}
                            placeholder="you@company.com" required />
                    </div>
                    <div className="form-group">
                        <label>Password</label>
                        <input className="input" type="password" value={password} onChange={e => setPassword(e.target.value)}
                            placeholder="Min 8 characters" required minLength={8} />
                    </div>
                    <button className="btn btn-primary" type="submit" disabled={loading}>
                        {loading ? 'Creating account...' : 'Create Account'}
                    </button>
                </form>
                <div className="auth-footer">
                    Already have an account? <Link to="/login">Sign in</Link>
                </div>
            </div>
        </div>
    );
}
