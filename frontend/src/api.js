const API_BASE = '/api/v1';

function getToken() {
    return localStorage.getItem('token');
}

function setToken(token) {
    localStorage.setItem('token', token);
}

function clearAuth() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
}

async function request(endpoint, options = {}) {
    const token = getToken();
    const headers = { 'Content-Type': 'application/json', ...options.headers };
    if (token) headers['Authorization'] = `Bearer ${token}`;

    const res = await fetch(`${API_BASE}${endpoint}`, { ...options, headers });

    if (res.status === 401) {
        clearAuth();
        window.location.href = '/login';
        throw new Error('Session expired. Please log in again.');
    }

    const data = await res.json();
    if (!res.ok) {
        throw new Error(data.error || `Request failed (${res.status})`);
    }
    return data;
}

// ---- Auth ----
export async function login(email, password) {
    const data = await request('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password }),
    });
    setToken(data.token);
    localStorage.setItem('user', JSON.stringify(data.user));
    return data;
}

export async function signup(email, password, fullName, companyName) {
    const data = await request('/auth/signup', {
        method: 'POST',
        body: JSON.stringify({ email, password, full_name: fullName, company_name: companyName }),
    });
    setToken(data.token);
    localStorage.setItem('user', JSON.stringify(data.user));
    return data;
}

export function logout() {
    clearAuth();
    window.location.href = '/login';
}

export function getUser() {
    try { return JSON.parse(localStorage.getItem('user')); }
    catch { return null; }
}

export function isAuthenticated() {
    return !!getToken();
}

// ---- Profile ----
export async function getProfile() {
    return request('/me');
}

export async function updateProfile(data) {
    const result = await request('/me', {
        method: 'PUT',
        body: JSON.stringify(data),
    });
    // Update cached user
    if (result.user) {
        localStorage.setItem('user', JSON.stringify(result.user));
    }
    return result;
}

export async function changePassword(currentPassword, newPassword) {
    return request('/me/password', {
        method: 'PUT',
        body: JSON.stringify({ current_password: currentPassword, new_password: newPassword }),
    });
}

// ---- Inbox ----
export async function getConversations() {
    return request('/inbox/conversations');
}

export async function getMessages(contactId) {
    return request(`/inbox/messages/${contactId}`);
}

export async function sendMessage(contactId, content) {
    return request(`/inbox/messages/${contactId}`, {
        method: 'POST',
        body: JSON.stringify({ content, message_type: 'text' }),
    });
}

export async function getContacts() {
    return request('/inbox/contacts');
}

// ---- Channels ----
export async function connectChannel(platform, accountId, accountName, accessToken) {
    return request('/channels', {
        method: 'POST',
        body: JSON.stringify({ platform, account_id: accountId, account_name: accountName, access_token: accessToken }),
    });
}

export async function getChannels() {
    return request('/channels');
}

export async function disconnectChannel(channelId) {
    return request(`/channels/${channelId}`, { method: 'DELETE' });
}

// ---- Automations ----
export async function getAutomations() {
    return request('/automations');
}

export async function createAutomation(data) {
    return request('/automations', {
        method: 'POST',
        body: JSON.stringify(data),
    });
}

export async function deleteAutomation(id) {
    return request(`/automations/${id}`, { method: 'DELETE' });
}

// ---- Broadcasts ----
export async function getBroadcasts() {
    return request('/broadcasts');
}

export async function createBroadcast(data) {
    return request('/broadcasts', {
        method: 'POST',
        body: JSON.stringify(data),
    });
}

export async function sendBroadcast(id) {
    return request(`/broadcasts/${id}/send`, { method: 'POST' });
}

// ---- AI Workflows ----
export async function generateWorkflow(prompt) {
    return request('/workflows/generate', {
        method: 'POST',
        body: JSON.stringify({ prompt }),
    });
}

export async function getWorkflows() {
    return request('/workflows');
}

export async function createWorkflow(data) {
    return request('/workflows', {
        method: 'POST',
        body: JSON.stringify(data),
    });
}

export async function updateWorkflow(id, data) {
    return request(`/workflows/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
    });
}

export async function deleteWorkflow(id) {
    return request(`/workflows/${id}`, { method: 'DELETE' });
}
