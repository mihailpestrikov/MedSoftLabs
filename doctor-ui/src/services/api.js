const API_URL = '/api';

async function request(endpoint, options = {}) {
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Request failed' }));
    throw new Error(error.error || 'Request failed');
  }

  return response.json();
}

export async function getPractitioners() {
  return request('/practitioners');
}

export async function getEncountersByPractitioner(practitionerId) {
  return request(`/encounters/${practitionerId}`);
}

export async function updateEncounterStatus(encounterId, status) {
  return request(`/encounters/${encounterId}`, {
    method: 'PATCH',
    body: JSON.stringify({ status })
  });
}
