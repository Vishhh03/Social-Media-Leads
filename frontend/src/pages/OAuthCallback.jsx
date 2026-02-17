import { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

export default function OAuthCallback() {
    const navigate = useNavigate();
    const [params] = useSearchParams();

    useEffect(() => {
        const token = params.get('token');
        const userB64 = params.get('user');

        if (token && userB64) {
            try {
                const userJSON = atob(userB64);
                const user = JSON.parse(userJSON);
                localStorage.setItem('token', token);
                localStorage.setItem('user', JSON.stringify(user));
                navigate('/', { replace: true });
            } catch {
                navigate('/login?error=parse_failed', { replace: true });
            }
        } else {
            navigate('/login?error=missing_token', { replace: true });
        }
    }, [params, navigate]);

    return (
        <div className="auth-page">
            <div className="auth-box" style={{ textAlign: 'center' }}>
                <div className="spinner" style={{ margin: '0 auto 16px' }}></div>
                <p style={{ color: 'var(--text-secondary)' }}>Signing you in...</p>
            </div>
        </div>
    );
}
