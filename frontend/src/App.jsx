import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ToastProvider } from './components/Toast';
import { isAuthenticated } from './api';
import Layout from './components/Layout';
import Login from './pages/Login';
import Signup from './pages/Signup';
import Dashboard from './pages/Dashboard';
import Inbox from './pages/Inbox';
import Contacts from './pages/Contacts';
import Automations from './pages/Automations';
import Channels from './pages/Channels';
import Broadcasts from './pages/Broadcasts';
import Settings from './pages/Settings';
import OAuthCallback from './pages/OAuthCallback';

function ProtectedRoute({ children }) {
    return isAuthenticated() ? children : <Navigate to="/login" />;
}

export default function App() {
    return (
        <ToastProvider>
            <BrowserRouter>
                <Routes>
                    <Route path="/login" element={<Login />} />
                    <Route path="/signup" element={<Signup />} />
                    <Route path="/oauth/callback" element={<OAuthCallback />} />
                    <Route path="/" element={<ProtectedRoute><Layout /></ProtectedRoute>}>
                        <Route index element={<Dashboard />} />
                        <Route path="inbox" element={<Inbox />} />
                        <Route path="contacts" element={<Contacts />} />
                        <Route path="automations" element={<Automations />} />
                        <Route path="channels" element={<Channels />} />
                        <Route path="broadcasts" element={<Broadcasts />} />
                        <Route path="settings" element={<Settings />} />
                    </Route>
                </Routes>
            </BrowserRouter>
        </ToastProvider>
    );
}
