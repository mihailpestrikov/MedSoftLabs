<script>
  import { onMount, onDestroy } from 'svelte';
  import { patients, loading, error } from '../stores/patients.js';
  import { getPatients } from '../services/api.js';
  import { WebSocketService } from '../services/websocket.js';

  let ws;

  onMount(async () => {
    await loadPatients();

    const wsUrl = import.meta.env.DEV
      ? 'ws://localhost:8081/ws'
      : `ws://${window.location.host}/ws`;

    ws = new WebSocketService(wsUrl);

    ws.on('patient_created', (patient) => {
      patients.update(p => [patient, ...p]);
    });

    ws.on('patient_deleted', ({ id }) => {
      patients.update(p => p.filter(patient => patient.id !== id));
    });

    ws.connect();
  });

  onDestroy(() => {
    if (ws) {
      ws.disconnect();
    }
  });

  async function loadPatients() {
    loading.set(true);
    error.set(null);

    try {
      const data = await getPatients();
      const sorted = data.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
      patients.set(sorted);
    } catch (e) {
      error.set(e.message);
    } finally {
      loading.set(false);
    }
  }

  let copiedId = null;

  function formatDate(date) {
    return new Date(date).toLocaleDateString();
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
</script>

<div class="card">
  <div class="header-row">
    <div>
      <h2>Registered Patients</h2>
      <p class="subtitle">{$patients.length} patient{$patients.length !== 1 ? 's' : ''} in system</p>
    </div>
    <div class="status-badge">
      <span class="status-dot"></span>
      Live
    </div>
  </div>

  {#if $error}
    <div class="error">{$error}</div>
  {/if}

  {#if $loading}
    <p class="loading">Loading patients...</p>
  {:else if $patients.length === 0}
    <div class="empty-state">
      <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
        <circle cx="9" cy="7" r="4"></circle>
        <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
        <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
      </svg>
      <h3>No patients yet</h3>
      <p>Waiting for first patient registration</p>
    </div>
  {:else}
    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th>
              <div class="th-content">
                <span>Patient ID</span>
              </div>
            </th>
            <th>
              <div class="th-content">
                <span>Name</span>
              </div>
            </th>
            <th>
              <div class="th-content">
                <span>Date of Birth</span>
              </div>
            </th>
            <th>
              <div class="th-content">
                <span>Gender</span>
              </div>
            </th>
            <th>
              <div class="th-content">
                <span>Registered</span>
              </div>
            </th>
          </tr>
        </thead>
        <tbody>
          {#each $patients as patient, i}
            <tr style="animation-delay: {i * 50}ms">
              <td>
                <div class="uuid-cell">
                  <button
                    class="uuid-badge"
                    class:copied={copiedId === patient.id}
                    on:click={() => copyUUID(patient.id)}
                    title="Click to copy full UUID"
                  >
                    {#if copiedId === patient.id}
                      Copied
                    {:else}
                      {patient.id.slice(0, 8)}
                    {/if}
                  </button>
                  <span class="uuid-full">{patient.id}</span>
                </div>
              </td>
              <td>
                <div class="name-cell">
                  <span class="name">{patient.last_name} {patient.first_name}</span>
                  {#if patient.middle_name}
                    <span class="middle-name">{patient.middle_name}</span>
                  {/if}
                </div>
              </td>
              <td>{formatDate(patient.date_of_birth)}</td>
              <td>
                <span class="gender-badge gender-{patient.gender}">{patient.gender}</span>
              </td>
              <td class="date-cell">{formatDate(patient.created_at)}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

<style>
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

  .status-badge {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    background: linear-gradient(135deg, var(--primary) 0%, var(--accent) 100%);
    color: white;
    border-radius: 100px;
    font-size: 0.8125rem;
    font-weight: 600;
    box-shadow: 0 4px 12px rgba(20, 184, 166, 0.25);
  }

  .status-dot {
    width: 8px;
    height: 8px;
    background: white;
    border-radius: 50%;
    animation: pulse 2s ease-in-out infinite;
  }

  @keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.5; }
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

  .th-content {
    display: flex;
    align-items: center;
    gap: 0.5rem;
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

  .uuid-cell {
    position: relative;
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
    display: inline-block;
    cursor: pointer;
    transition: all 0.2s ease;
    outline: none;
    min-width: 80px;
    text-align: center;
  }

  .uuid-badge:focus-visible {
    box-shadow: 0 0 0 3px rgba(20, 184, 166, 0.2);
  }

  .uuid-badge.copied {
    pointer-events: none;
  }

  .check-icon {
    font-weight: bold;
    font-size: 0.875rem;
  }

  .uuid-badge:hover {
    background: linear-gradient(135deg, rgba(20, 184, 166, 0.2) 0%, rgba(6, 182, 212, 0.2) 100%);
    border-color: rgba(20, 184, 166, 0.4);
    transform: translateY(-1px);
    box-shadow: 0 2px 4px rgba(20, 184, 166, 0.15);
  }

  .uuid-badge:active {
    transform: translateY(0);
    box-shadow: none;
  }

  .uuid-full {
    display: none;
  }

  .name-cell {
    display: flex;
    flex-direction: column;
    gap: 0.125rem;
  }

  .name {
    font-weight: 500;
    color: var(--text);
  }

  .middle-name {
    font-size: 0.8125rem;
    color: var(--text-light);
  }

  .gender-badge {
    display: inline-block;
    padding: 0.25rem 0.75rem;
    border-radius: 100px;
    font-size: 0.8125rem;
    font-weight: 500;
    text-transform: capitalize;
  }

  .gender-male {
    background: rgba(59, 130, 246, 0.1);
    color: #2563eb;
  }

  .gender-female {
    background: rgba(236, 72, 153, 0.1);
    color: #db2777;
  }

  .date-cell {
    color: var(--text-secondary);
  }

  .error {
    color: var(--error);
    padding: 1rem;
    background: #fef2f2;
    border-radius: var(--radius-sm);
    margin-bottom: 1rem;
  }
</style>
