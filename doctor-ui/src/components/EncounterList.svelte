<script>
  import { onMount, onDestroy } from 'svelte';
  import { encounters, loading, error } from '../stores/encounters.js';
  import { practitioners, selectedPractitionerId } from '../stores/practitioners.js';
  import { getPractitioners, getEncountersByPractitioner, updateEncounterStatus } from '../services/api.js';
  import { WebSocketService } from '../services/websocket.js';

  let ws;
  let updatingStatus = {};

  onMount(async () => {
    await loadPractitioners();

    const wsUrl = import.meta.env.DEV
      ? 'wss://localhost:8081/ws'
      : `wss://${window.location.host}/ws`;

    ws = new WebSocketService(wsUrl);

    ws.on('encounter_created', (encounterData) => {
      if (!$selectedPractitionerId) return;

      const encounter = parseEncounter(encounterData);
      if (encounter.practitionerId === $selectedPractitionerId) {
        encounters.update(e => [encounter, ...e]);
      }
    });

    ws.on('encounter_status_updated', (encounterData) => {
      if (!$selectedPractitionerId) return;

      const encounter = parseEncounter(encounterData);

      if (encounter.practitionerId === $selectedPractitionerId) {
        encounters.update(e => {
          const updated = e.map(enc => enc.id === encounter.id ? encounter : enc);
          return updated;
        });
      }
    });

    ws.connect();
  });

  onDestroy(() => {
    if (ws) {
      ws.disconnect();
    }
  });

  async function loadPractitioners() {
    try {
      const data = await getPractitioners();
      practitioners.set(data);

      if (data.length > 0 && !$selectedPractitionerId) {
        selectedPractitionerId.set(data[0].id);
        await loadEncounters(data[0].id);
      }
    } catch (e) {
      error.set(e.message);
    }
  }

  async function loadEncounters(practitionerId) {
    if (!practitionerId) return;

    loading.set(true);
    error.set(null);

    try {
      const data = await getEncountersByPractitioner(practitionerId);
      const encounters_data = data || [];
      const parsed = encounters_data.map(e => parseEncounter(e));
      const sorted = parsed.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));
      encounters.set(sorted);
    } catch (e) {
      error.set(e.message);
      encounters.set([]);
    } finally {
      loading.set(false);
    }
  }

  async function handlePractitionerChange(event) {
    const practitionerId = event.target.value;
    selectedPractitionerId.set(practitionerId);
    await loadEncounters(practitionerId);
  }

  async function handleStatusChange(encounterId, newStatus) {
    updatingStatus[encounterId] = true;
    error.set(null);

    try {
      await updateEncounterStatus(encounterId, newStatus);

      encounters.update(e =>
        e.map(enc => enc.id === encounterId ? { ...enc, status: newStatus } : enc)
      );
    } catch (e) {
      error.set(e.message);
    } finally {
      updatingStatus[encounterId] = false;
    }
  }

  function getValue(field) {
    if (typeof field === 'string') return field;
    return field?.value || '';
  }

  function parseEncounter(encounterData) {
    const id = getValue(encounterData.id) || encounterData.id || '';
    const status = getValue(encounterData.status) || encounterData.status || 'planned';

    let normalizedStatus = status.toLowerCase().replace(/_/g, '-');
    if (normalizedStatus === 'finished') {
      normalizedStatus = 'completed';
    }

    let practitionerId = encounterData.practitionerId || '';
    let patientId = encounterData.patientId || '';
    let patientName = '';
    let patientGender = '';
    let practitionerName = encounterData.practitionerName || encounterData.practitionerDisplay || '';
    let createdAt = encounterData.createdAt || new Date().toISOString();

    if (encounterData.participant && Array.isArray(encounterData.participant)) {
      const participant = encounterData.participant[0];
      if (participant?.individual) {
        practitionerName = getValue(participant.individual.display) || practitionerName;
        const ref = getValue(participant.individual.reference);
        if (ref && ref.includes('/')) {
          practitionerId = ref.split('/').pop();
        }
      }
    }

    if (encounterData.subject) {
      let fullDisplay = getValue(encounterData.subject.display) || encounterData.patientName || encounterData.patientDisplay || '';
      const genderMatch = fullDisplay.match(/\s*\[(male|female)\]$/i);
      if (genderMatch) {
        patientGender = genderMatch[1].toLowerCase();
        patientName = fullDisplay.replace(/\s*\[(male|female)\]$/i, '');
      } else {
        patientName = fullDisplay;
      }

      const ref = getValue(encounterData.subject.reference);
      if (ref && ref.includes('/')) {
        patientId = ref.split('/').pop();
      }
    } else {
      let fallbackName = encounterData.patientName || encounterData.patientDisplay || '';
      const genderMatch = fallbackName.match(/\s*\[(male|female)\]$/i);
      if (genderMatch) {
        patientGender = genderMatch[1].toLowerCase();
        patientName = fallbackName.replace(/\s*\[(male|female)\]$/i, '');
      } else {
        patientName = fallbackName;
      }
    }

    if (encounterData.period?.start) {
      const valueUs = getValue(encounterData.period.start.valueUs);
      if (valueUs) {
        const microseconds = parseInt(valueUs);
        const milliseconds = microseconds / 1000;
        createdAt = new Date(milliseconds).toISOString();
      }
    }

    return {
      id,
      patientId,
      patientName,
      patientGender,
      practitionerId,
      practitionerName,
      status: normalizedStatus,
      createdAt
    };
  }

  function formatDate(date) {
    const d = new Date(date);
    const day = String(d.getDate()).padStart(2, '0');
    const month = String(d.getMonth() + 1).padStart(2, '0');
    const year = d.getFullYear();
    const hours = String(d.getHours()).padStart(2, '0');
    const minutes = String(d.getMinutes()).padStart(2, '0');
    return `${day}.${month}.${year}, ${hours}:${minutes}`;
  }

  function getStatusColor(status) {
    switch (status) {
      case 'planned':
        return 'status-planned';
      case 'arrived':
        return 'status-arrived';
      case 'in-progress':
        return 'status-in-progress';
      case 'completed':
        return 'status-completed';
      case 'cancelled':
        return 'status-cancelled';
      default:
        return 'status-planned';
    }
  }

  function getNextStatus(currentStatus) {
    switch (currentStatus) {
      case 'arrived':
        return 'in-progress';
      case 'in-progress':
        return 'completed';
      default:
        return null;
    }
  }

  function getNextStatusLabel(currentStatus) {
    switch (currentStatus) {
      case 'arrived':
        return 'Start Examination';
      case 'in-progress':
        return 'Complete';
      default:
        return null;
    }
  }
</script>

<div class="card">
  <div class="header-row">
    <div>
      <h2>Patient Encounters</h2>
      <p class="subtitle">
        {#if $selectedPractitionerId}
          {$encounters.length} encounter{$encounters.length !== 1 ? 's' : ''}
        {:else}
          Select a practitioner to view encounters
        {/if}
      </p>
    </div>
    <div class="practitioner-selector">
      <label for="practitioner">Practitioner:</label>
      <select
        id="practitioner"
        value={$selectedPractitionerId}
        on:change={handlePractitionerChange}
      >
        {#each $practitioners as practitioner}
          <option value={practitioner.id}>
            {practitioner.lastName} {practitioner.firstName}
            {#if practitioner.middleName}
              {practitioner.middleName}
            {/if}
            {#if practitioner.specialization}
              ({practitioner.specialization})
            {/if}
          </option>
        {/each}
      </select>
    </div>
  </div>

  {#if $error}
    <div class="error">{$error}</div>
  {/if}

  {#if $loading}
    <p class="loading">Loading encounters...</p>
  {:else if !$selectedPractitionerId}
    <div class="empty-state">
      <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
        <circle cx="12" cy="7" r="4"></circle>
      </svg>
      <h3>No practitioner selected</h3>
      <p>Select a practitioner to view their encounters</p>
    </div>
  {:else if $encounters.length === 0}
    <div class="empty-state">
      <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M9 11l3 3L22 4"></path>
        <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"></path>
      </svg>
      <h3>No encounters yet</h3>
      <p>Waiting for patient encounters</p>
    </div>
  {:else}
    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th>Patient</th>
            <th>Gender</th>
            <th>Time</th>
            <th>Status</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each $encounters as encounter, i}
            <tr style="animation-delay: {i * 50}ms">
              <td>
                <div class="patient-cell">
                  <span class="patient-name">{encounter.patientName || encounter.patientId || 'Unknown'}</span>
                  <span class="patient-id">{(encounter.patientId || '').slice(0, 8) || 'N/A'}</span>
                </div>
              </td>
              <td>
                {#if encounter.patientGender}
                  <span class="gender-badge gender-{encounter.patientGender.toLowerCase()}">
                    {encounter.patientGender}
                  </span>
                {/if}
              </td>
              <td class="date-cell">{encounter.createdAt ? formatDate(encounter.createdAt) : 'N/A'}</td>
              <td>
                <span class="status-badge {getStatusColor(encounter.status || 'planned')}">
                  {encounter.status || 'planned'}
                </span>
              </td>
              <td>
                {#if encounter.id && getNextStatus(encounter.status)}
                  <button
                    class="action-button"
                    disabled={updatingStatus[encounter.id]}
                    on:click={() => handleStatusChange(encounter.id, getNextStatus(encounter.status))}
                  >
                    {updatingStatus[encounter.id] ? 'Updating...' : getNextStatusLabel(encounter.status)}
                  </button>
                {:else}
                  <span class="no-action">
                    {(encounter.status || 'planned') === 'completed' ? 'Completed' : 'No action'}
                  </span>
                {/if}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

<style>
  .card {
    max-width: 1000px;
    margin: 0 auto;
  }

  .header-row {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 2rem;
    padding-bottom: 1.25rem;
    border-bottom: 2px solid var(--border);
  }

  h2 {
    color: var(--text);
    font-size: 1.25rem;
    font-weight: 600;
    letter-spacing: -0.01em;
    margin-bottom: 0.25rem;
  }

  .subtitle {
    color: var(--text-light);
    font-size: 0.875rem;
  }

  .practitioner-selector {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .practitioner-selector label {
    color: var(--text-secondary);
    font-size: 0.9375rem;
    font-weight: 500;
  }

  .practitioner-selector select {
    padding: 0.625rem 1rem;
    border: 1px solid var(--border);
    border-radius: 8px;
    font-size: 0.9375rem;
    color: var(--text);
    background: white;
    cursor: pointer;
    transition: all 0.2s ease;
    min-width: 250px;
  }

  .practitioner-selector select:focus {
    outline: none;
    border-color: var(--primary);
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }

  .empty-state {
    text-align: center;
    padding: 4rem 2rem;
  }

  .empty-state svg {
    color: var(--text-light);
    margin-bottom: 1.5rem;
  }

  .empty-state h3 {
    color: var(--text);
    font-size: 1.125rem;
    font-weight: 600;
    margin-bottom: 0.5rem;
  }

  .empty-state p {
    color: var(--text-light);
    font-size: 0.9375rem;
  }

  .table-container {
    overflow-x: auto;
    border-radius: var(--radius);
  }

  tbody tr {
    animation: fadeIn 0.3s ease-out forwards;
    opacity: 0;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(-10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .patient-cell {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .patient-name {
    font-weight: 500;
    color: var(--text);
  }

  .patient-id {
    font-family: 'Courier New', monospace;
    font-size: 0.75rem;
    color: var(--text-light);
  }

  .status-badge {
    display: inline-block;
    padding: 0.375rem 0.875rem;
    border-radius: 100px;
    font-size: 0.8125rem;
    font-weight: 500;
    text-transform: capitalize;
  }

  .status-planned {
    background: rgba(107, 114, 128, 0.1);
    color: #4b5563;
  }

  .status-arrived {
    background: rgba(16, 185, 129, 0.1);
    color: #059669;
  }

  .status-in-progress {
    background: rgba(59, 130, 246, 0.1);
    color: #2563eb;
  }

  .status-completed {
    background: rgba(139, 92, 246, 0.1);
    color: #7c3aed;
  }

  .status-cancelled {
    background: rgba(239, 68, 68, 0.1);
    color: #dc2626;
  }

  .date-cell {
    color: var(--text-secondary);
    font-size: 0.875rem;
  }

  .action-button {
    padding: 0.5rem 1rem;
    background: linear-gradient(135deg, var(--primary) 0%, var(--accent) 100%);
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 0.875rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    min-width: 140px;
  }

  .action-button:hover:not(:disabled) {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(59, 130, 246, 0.25);
  }

  .action-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .no-action {
    color: var(--text-light);
    font-size: 0.875rem;
    font-style: italic;
  }

  .gender-badge {
    display: inline-block;
    padding: 0.375rem 0.875rem;
    border-radius: 100px;
    font-size: 0.8125rem;
    font-weight: 500;
    text-transform: capitalize;
  }

  .gender-badge.gender-male {
    background: rgba(59, 130, 246, 0.1);
    color: #2563eb;
  }

  .gender-badge.gender-female {
    background: rgba(236, 72, 153, 0.1);
    color: #db2777;
  }

  .error {
    color: var(--error);
    padding: 1rem;
    background: #fef2f2;
    border-radius: var(--radius-sm);
    margin-bottom: 1rem;
  }
</style>
