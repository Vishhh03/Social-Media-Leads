import { useState, useEffect } from 'react';
import { getProfile, updateProfile, changePassword, logout } from '../api';
import { useToast } from '../components/Toast';

export default function Settings() {
    const toast = useToast();
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);

    // Profile form
    const [fullName, setFullName] = useState('');
    const [email, setEmail] = useState('');
    const [companyName, setCompanyName] = useState('');

    // Password form
    const [currentPassword, setCurrentPassword] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');

    useEffect(() => {
        loadProfile();
    }, []);

    async function loadProfile() {
        try {
            const data = await getProfile();
            setUser(data.user);
            setFullName(data.user.full_name || '');
            setEmail(data.user.email || '');
            setCompanyName(data.user.company_name || '');
        } catch (err) {
            toast.error(err.message);
        } finally {
            setLoading(false);
        }
    }

    async function handleProfileSave(e) {
        e.preventDefault();
        setSaving(true);
        try {
            const data = await updateProfile({ full_name: fullName, email, company_name: companyName });
            setUser(data.user);
            toast.success('Profile updated successfully');
        } catch (err) {
            toast.error(err.message);
        } finally {
            setSaving(false);
        }
    }

    async function handlePasswordChange(e) {
        e.preventDefault();
        if (newPassword !== confirmPassword) {
            toast.error('New passwords do not match');
            return;
        }
        if (newPassword.length < 8) {
            toast.error('Password must be at least 8 characters');
            return;
        }
        setSaving(true);
        try {
            await changePassword(currentPassword, newPassword);
            toast.success('Password changed successfully');
            setCurrentPassword('');
            setNewPassword('');
            setConfirmPassword('');
        } catch (err) {
            toast.error(err.message);
        } finally {
            setSaving(false);
        }
    }

    if (loading) return <div className="page-content"><p>Loading...</p></div>;

    return (
        <div className="page-content">
            <h1>Settings</h1>

            {/* Profile Section */}
            <div className="card" style={{ marginBottom: '2rem' }}>
                <h2 style={{ marginBottom: '0.25rem' }}>Profile</h2>
                <p style={{ color: 'var(--text-secondary)', marginBottom: '1.5rem', fontSize: '0.875rem' }}>
                    Update your personal information
                </p>

                <form onSubmit={handleProfileSave}>
                    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                        <div className="form-group">
                            <label>Full Name</label>
                            <input
                                type="text"
                                value={fullName}
                                onChange={e => setFullName(e.target.value)}
                                placeholder="John Doe"
                            />
                        </div>
                        <div className="form-group">
                            <label>Email</label>
                            <input
                                type="email"
                                value={email}
                                onChange={e => setEmail(e.target.value)}
                                placeholder="you@company.com"
                            />
                        </div>
                        <div className="form-group">
                            <label>Company Name</label>
                            <input
                                type="text"
                                value={companyName}
                                onChange={e => setCompanyName(e.target.value)}
                                placeholder="Acme Inc."
                            />
                        </div>
                        <div className="form-group">
                            <label>Plan</label>
                            <input type="text" value={user?.plan || 'starter'} disabled
                                style={{ background: 'var(--bg-secondary)', cursor: 'not-allowed' }} />
                        </div>
                    </div>
                    <div style={{ marginTop: '1rem' }}>
                        <button type="submit" className="btn btn-primary" disabled={saving}>
                            {saving ? 'Saving...' : 'Save Changes'}
                        </button>
                    </div>
                </form>
            </div>

            {/* Password Section */}
            <div className="card" style={{ marginBottom: '2rem' }}>
                <h2 style={{ marginBottom: '0.25rem' }}>Change Password</h2>
                <p style={{ color: 'var(--text-secondary)', marginBottom: '1.5rem', fontSize: '0.875rem' }}>
                    Update your password to keep your account secure
                </p>

                <form onSubmit={handlePasswordChange}>
                    <div className="form-group">
                        <label>Current Password</label>
                        <input
                            type="password"
                            value={currentPassword}
                            onChange={e => setCurrentPassword(e.target.value)}
                            placeholder="••••••••"
                            required
                        />
                    </div>
                    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                        <div className="form-group">
                            <label>New Password</label>
                            <input
                                type="password"
                                value={newPassword}
                                onChange={e => setNewPassword(e.target.value)}
                                placeholder="Min 8 characters"
                                required
                                minLength={8}
                            />
                        </div>
                        <div className="form-group">
                            <label>Confirm New Password</label>
                            <input
                                type="password"
                                value={confirmPassword}
                                onChange={e => setConfirmPassword(e.target.value)}
                                placeholder="••••••••"
                                required
                            />
                        </div>
                    </div>
                    <div style={{ marginTop: '1rem' }}>
                        <button type="submit" className="btn btn-primary" disabled={saving}>
                            {saving ? 'Updating...' : 'Change Password'}
                        </button>
                    </div>
                </form>
            </div>

            {/* Account Info Section */}
            <div className="card" style={{ marginBottom: '2rem' }}>
                <h2 style={{ marginBottom: '0.25rem' }}>Account</h2>
                <p style={{ color: 'var(--text-secondary)', marginBottom: '1.5rem', fontSize: '0.875rem' }}>
                    Manage your account
                </p>

                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '1rem 0', borderTop: '1px solid var(--border)' }}>
                    <div>
                        <p style={{ fontWeight: 600 }}>Member since</p>
                        <p style={{ color: 'var(--text-secondary)', fontSize: '0.875rem' }}>
                            {user?.created_at ? new Date(user.created_at).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' }) : '—'}
                        </p>
                    </div>
                </div>

                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '1rem 0', borderTop: '1px solid var(--border)' }}>
                    <div>
                        <p style={{ fontWeight: 600, color: '#dc2626' }}>Sign Out</p>
                        <p style={{ color: 'var(--text-secondary)', fontSize: '0.875rem' }}>
                            Sign out of your account on this device
                        </p>
                    </div>
                    <button className="btn" onClick={logout}
                        style={{ background: '#fef2f2', color: '#dc2626', border: '1px solid #fecaca' }}>
                        Sign Out
                    </button>
                </div>
            </div>
        </div>
    );
}
