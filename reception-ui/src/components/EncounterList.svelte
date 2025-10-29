<script>
  import { onMount, onDestroy } from 'svelte';
  import { encounters, loadingEncounters, encounterError } from '../stores/encounters.js';
  import { getEncounters, updateEncounterStatus } from '../services/api.js';
  import { WebSocketService } from '../services/websocket.js';

  let ws;
  let copiedId = null;
  let updatingStatus = {};

  const statuses = [
    { value: 'planned', label: 'Planned', color: 'gray' },
    { value: 'arrived', label: 'Arrived', color: 'green' },
    { value: 'in-progress', label: 'In Progress', color: 'blue' },
    { value: 'completed', label: 'Completed', color: 'purple' },
    { value: 'cancelled', label: 'Cancelled', color: 'red' }
  ];

  function getValue(field) {
    if (typeof field === 'string') return field;
    return field?.value || '';
  }

  onMount(async () => {
    await loadEncounters();

    const wsUrl = import.meta.env.DEV
      ? 'wss://localhost:8080/ws'
      : `wss://${window.location.host}/ws`;

    ws = new WebSocketService(wsUrl);

    ws.on('encounter_created', (encounter) => {
      encounters.update(e => [encounter, ...e]);
    });

    ws.on('encounter_status_updated', (encounterData) => {
      const encounterId = getValue(encounterData.id);

      let statusValue = getValue(encounterData.status);
      if (statusValue) {
        statusValue = statusValue.toLowerCase().replace(/_/g, '-');
        if (statusValue === 'finished') {
          statusValue = 'completed';
        }
      }

      encounters.update(e => e.map(enc =>
        getValue(enc.id) === encounterId ? { ...enc, status: { value: statusValue } } : enc
      ));
    });

    ws.connect();
  });

  async function loadEncounters() {
    loadingEncounters.set(true);
    encounterError.set(null);

    try {
      const data = await getEncounters();
      console.log('Loaded encounters:', data);
      encounters.set(data);
    } catch (e) {
      console.error('Failed to load encounters:', e);
      encounterError.set(e.message);
    } finally {
      loadingEncounters.set(false);
    }
  }

  onDestroy(() => {
    if (ws) {
      ws.disconnect();
    }
  });

  function formatDateTime(field) {
    const dateString = getValue(field);
    const date = new Date(dateString);
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const year = date.getFullYear();
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${day}.${month}.${year}, ${hours}:${minutes}`;
  }

  function formatDateOnly(dateString) {
    return new Date(dateString).toLocaleDateString();
  }

  async function copyUUID(uuid) {
    try {
      await navigator.clipboard.writeText(uuid);
      copiedId = uuid;
      setTimeout(() => {
        copiedId = null;
      }, 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  }

  async function handleStatusChange(encounter, newStatus) {
    const encounterId = getValue(encounter.id);
    updatingStatus[encounterId] = true;
    try {
      await updateEncounterStatus(encounterId, newStatus);
      encounters.update(e => e.map(enc =>
        getValue(enc.id) === encounterId ? { ...enc, status: { value: newStatus } } : enc
      ));
    } catch (err) {
      console.error('Failed to update status:', err);
      alert('Failed to update status: ' + err.message);
    } finally {
      updatingStatus[encounterId] = false;
    }
  }

  function getStatusColor(status) {
    const statusObj = statuses.find(s => s.value === status);
    return statusObj ? statusObj.color : 'green';
  }

  function getStatusLabel(status) {
    const statusObj = statuses.find(s => s.value === status);
    return statusObj ? statusObj.label : status;
  }

  $: todayVisits = $encounters.filter(e => {
    const visitDate = new Date(e.start_time);
    const today = new Date();
    return visitDate.toDateString() === today.toDateString();
  });
</script>

<div class="card">
  <div class="header-row">
    <div>
      <h2>Recent Visits</h2>
      <p class="subtitle">{todayVisits.length} visit{todayVisits.length !== 1 ? 's' : ''} today</p>
    </div>
  </div>

  {#if $encounterError}
    <div class="error">{$encounterError}</div>
  {/if}

  {#if $loadingEncounters && $encounters.length === 0}
    <p class="loading">Loading...</p>
  {:else if $encounters.length === 0}
    <div class="empty-state">
      <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
      </svg>
      <h3>No visits yet</h3>
      <p>Register the first patient visit</p>
    </div>
  {:else}
    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th>Visit ID</th>
            <th>Patient</th>
            <th>Doctor</th>
            <th>Time</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {#each $encounters as encounter, i}
            <tr style="animation-delay: {i * 50}ms">
              <td>
                <button
                  class="uuid-badge"
                  class:copied={copiedId === getValue(encounter.id)}
                  on:click={() => copyUUID(getValue(encounter.id))}
                  title="Click to copy full UUID"
                >
                  {#if copiedId === getValue(encounter.id)}
                    Copied
                  {:else}
                    {getValue(encounter.id).slice(0, 8)}
                  {/if}
                </button>
              </td>
              <td>
                <div class="name-cell">
                  <span class="name">
                    {encounter.patient?.last_name || 'Unknown'}
                    {encounter.patient?.first_name || ''}
                  </span>
                  {#if encounter.patient?.middle_name}
                    <span class="middle-name">{encounter.patient.middle_name}</span>
                  {/if}
                </div>
              </td>
              <td>
                <div class="doctor-cell">
                  <span class="doctor-name">
                    {encounter.practitioner?.LastName || 'Unknown'}
                    {encounter.practitioner?.FirstName || ''}
                    {#if encounter.practitioner?.MiddleName}
                      {encounter.practitioner.MiddleName}
                    {/if}
                  </span>
                  <span class="specialization">{encounter.practitioner?.Specialization || ''}</span>
                </div>
              </td>
              <td class="time-cell">{formatDateTime(encounter.start_time)}</td>
              <td>
                <select
                  class="status-select status-{getStatusColor(getValue(encounter.status))}"
                  value={getValue(encounter.status)}
                  on:change={(e) => handleStatusChange(encounter, e.target.value)}
                  disabled={updatingStatus[getValue(encounter.id)]}
                >
                  {#each statuses as status}
                    <option value={status.value}>{status.label}</option>
                  {/each}
                </select>
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
    background: var(--card-bg);
    border-radius: var(--radius);
    padding: 1.75rem;
    border: 1px solid var(--border);
    box-shadow: var(--shadow-sm);
    min-width: 750px;
  }

  .header-row {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 1.5rem;
    padding-bottom: 1rem;
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

  .table-container {
    overflow-x: auto;
    border-radius: var(--radius);
    min-width: 600px;
  }

  table {
    width: 100%;
    border-collapse: separate;
    border-spacing: 0;
    background: var(--card-bg);
    border-radius: var(--radius);
    overflow: hidden;
    border: 1px solid var(--border);
  }

  th, td {
    padding: 0.875rem 1rem;
    text-align: left;
    border-bottom: 1px solid var(--border);
  }

  th {
    background: var(--hover);
    font-weight: 600;
    font-size: 0.8125rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--text-secondary);
  }

  td {
    font-size: 0.9375rem;
    color: var(--text);
  }

  tr:last-child td {
    border-bottom: none;
  }

  tbody tr {
    transition: background 0.15s ease;
    animation: fadeIn 0.3s ease-out forwards;
    opacity: 0;
  }

  tbody tr:hover {
    background: var(--hover);
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

  .uuid-badge {
    font-family: 'Courier New', monospace;
    font-size: 0.75rem;
    font-weight: 600;
    color: var(--primary);
    background: linear-gradient(135deg, rgba(20, 184, 166, 0.1) 0%, rgba(6, 182, 212, 0.1) 100%);
    padding: 0.25rem 0.625rem;
    border-radius: 4px;
    border: 1px solid rgba(20, 184, 166, 0.2);
    cursor: pointer;
    transition: all 0.2s ease;
    outline: none;
    min-width: 80px;
    text-align: center;
  }

  .uuid-badge:hover {
    background: linear-gradient(135deg, rgba(20, 184, 166, 0.2) 0%, rgba(6, 182, 212, 0.2) 100%);
    border-color: rgba(20, 184, 166, 0.4);
    transform: translateY(-1px);
    box-shadow: 0 2px 4px rgba(20, 184, 166, 0.15);
  }

  .uuid-badge.copied {
    pointer-events: none;
  }

  .name-cell, .doctor-cell {
    display: flex;
    flex-direction: column;
    gap: 0.125rem;
  }

  .name, .doctor-name {
    font-weight: 500;
    color: var(--text);
  }

  .middle-name, .specialization {
    font-size: 0.8125rem;
    color: var(--text-light);
  }

  .time-cell {
    color: var(--text-secondary);
    font-size: 0.875rem;
  }

  .status-select {
    padding: 0.375rem 0.625rem;
    border-radius: 6px;
    font-size: 0.8125rem;
    font-weight: 500;
    border: 1px solid;
    cursor: pointer;
    transition: all 0.2s ease;
    outline: none;
  }

  .status-select:hover:not(:disabled) {
    opacity: 0.85;
  }

  .status-select:disabled {
    cursor: not-allowed;
    opacity: 0.6;
  }

  .status-select.status-gray {
    background: rgba(107, 114, 128, 0.1);
    color: #4b5563;
    border-color: rgba(107, 114, 128, 0.3);
  }

  .status-select.status-green {
    background: rgba(16, 185, 129, 0.1);
    color: #059669;
    border-color: rgba(16, 185, 129, 0.3);
  }

  .status-select.status-blue {
    background: rgba(59, 130, 246, 0.1);
    color: #2563eb;
    border-color: rgba(59, 130, 246, 0.3);
  }

  .status-select.status-purple {
    background: rgba(139, 92, 246, 0.1);
    color: #7c3aed;
    border-color: rgba(139, 92, 246, 0.3);
  }

  .status-select.status-red {
    background: rgba(239, 68, 68, 0.1);
    color: #dc2626;
    border-color: rgba(239, 68, 68, 0.3);
  }

  .empty-state {
    text-align: center;
    padding: 4rem 2rem;
    min-width: 600px;
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

  .error {
    color: var(--error);
    background: #fef2f2;
    padding: 0.75rem;
    border-radius: var(--radius-sm);
    font-size: 0.875rem;
    margin-bottom: 1rem;
  }

  .loading {
    text-align: center;
    padding: 2rem;
    color: var(--text-light);
  }
</style>
