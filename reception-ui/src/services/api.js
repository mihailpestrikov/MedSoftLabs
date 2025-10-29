import { get } from 'svelte/store';
import { accessToken, clearAuth, setAuth } from '../stores/auth.js';

const API_URL = '/api';

async function request(endpoint, options = {}) {
  const token = get(accessToken);

  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  let response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers,
    credentials: 'include',
  });

  if (response.status === 401 && token) {
    const refreshed = await refresh();
    if (refreshed) {
      headers.Authorization = `Bearer ${get(accessToken)}`;
      response = await fetch(`${API_URL}${endpoint}`, {
        ...options,
        headers,
        credentials: 'include',
      });
    } else {
      clearAuth();
      throw new Error('Session expired');
    }
  }

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Request failed' }));
    throw new Error(error.error || 'Request failed');
  }

  return response.json();
}

export async function login(username, password) {
  const data = await request('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  });
  setAuth(data.access_token, username);
  return data;
}

export async function register(username, password) {
  return request('/auth/register', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  });
}

export async function logout() {
  await request('/auth/logout', { method: 'POST' });
  clearAuth();
}

export async function checkAuth() {
  const savedUsername = localStorage.getItem('username');
  if (!savedUsername) return;

  try {
    const response = await fetch(`${API_URL}/auth/refresh`, {
      method: 'POST',
      credentials: 'include',
    });

    if (response.ok) {
      const data = await response.json();
      if (data.access_token) {
        setAuth(data.access_token, savedUsername);
      }
    } else {
      localStorage.removeItem('username');
    }
  } catch {
    localStorage.removeItem('username');
  }
}

async function refresh() {
  try {
    const data = await fetch(`${API_URL}/auth/refresh`, {
      method: 'POST',
      credentials: 'include',
    }).then(r => r.json());

    if (data.access_token) {
      accessToken.set(data.access_token);
      return true;
    }
    return false;
  } catch {
    return false;
  }
}

export async function getPatients() {
  return request('/patients');
}

export async function createPatient(patient) {
  return request('/patients', {
    method: 'POST',
    body: JSON.stringify(patient),
  });
}

export async function deletePatient(id) {
  return request(`/patients/${id}`, {
    method: 'DELETE',
  });
}

export async function getPractitioners() {
  return request('/practitioners');
}

export async function createEncounter(patientId, practitionerId, startTime) {
  return request('/encounters', {
    method: 'POST',
    body: JSON.stringify({
      patient_id: patientId,
      practitioner_id: practitionerId,
      start_time: startTime,
    }),
  });
}

export async function getEncounters() {
  return request('/encounters');
}

export async function updateEncounterStatus(encounterId, status) {
  return request(`/encounters/${encounterId}`, {
    method: 'PATCH',
    body: JSON.stringify({ status }),
  });
}
