const API_URL = '/api';
const FHIR_URL = '/fhir';

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

async function fhirRequest(endpoint, options = {}) {
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  const response = await fetch(`${FHIR_URL}${endpoint}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Request failed' }));
    throw new Error(error.error || 'Request failed');
  }

  return response.json();
}

export async function getPatients() {
  return request('/patients');
}

export async function getPractitioners() {
  const bundle = await fhirRequest('/Practitioner');

  if (!bundle.entry) {
    return [];
  }

  return bundle.entry.map(entry => {
    const resource = entry.resource;
    const practitioner = {
      id: resource.id?.value || '',
      firstName: '',
      lastName: '',
      middleName: '',
      specialization: ''
    };

    if (resource.name && resource.name.length > 0) {
      const name = resource.name[0];
      practitioner.lastName = name.family?.value || '';
      if (name.given && name.given.length > 0) {
        practitioner.firstName = name.given[0]?.value || '';
        if (name.given.length > 1) {
          practitioner.middleName = name.given[1]?.value || '';
        }
      }
    }

    if (resource.qualification && resource.qualification.length > 0) {
      practitioner.specialization = resource.qualification[0]?.code?.text?.value || '';
    }

    return practitioner;
  });
}

export async function createPractitioner(firstName, lastName, middleName, specialization) {
  const given = [{ value: firstName }];
  if (middleName) {
    given.push({ value: middleName });
  }

  const fhirPractitioner = {
    name: [
      {
        family: { value: lastName },
        given: given
      }
    ],
    qualification: [
      {
        code: {
          text: { value: specialization }
        }
      }
    ]
  };

  const response = await fhirRequest('/Practitioner', {
    method: 'POST',
    body: JSON.stringify(fhirPractitioner)
  });

  return {
    id: response.id?.value || '',
    firstName,
    lastName,
    middleName,
    specialization
  };
}
