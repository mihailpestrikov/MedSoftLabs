<script>
  import { onMount, onDestroy } from 'svelte';
  import { patients, loading, error } from '../stores/patients.js';
  import { getPatients, deletePatient } from '../services/api.js';
  import { WebSocketService } from '../services/websocket.js';

  let ws;
  let copiedId = null;

  onMount(async () => {
    await loadPatients();

    const wsUrl = import.meta.env.DEV
      ? 'wss://localhost:8080/ws'
      : `wss://${window.location.host}/ws`;

    ws = new WebSocketService(wsUrl);

    ws.on('patient_created', (patient) => {
      patients.update(p => [...p, patient]);
    });

    ws.on('patient_deleted', ({ id }) => {
      patients.update(p => p.filter(patient => patient.id !== id));
    });

    ws.on('patient_his_id_update', ({ id, his_patient_id }) => {
      patients.update(p => p.map(patient =>
        patient.id === id ? { ...patient, his_patient_id } : patient
      ));
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
      patients.set(data);
    } catch (e) {
      error.set(e.message);
    } finally {
      loading.set(false);
    }
  }

  async function handleDelete(id) {
    if (!confirm('Delete this patient?')) return;

    try {
      await deletePatient(id);
      patients.update(p => p.filter(patient => patient.id !== id));
    } catch (e) {
      alert(`Failed to delete: ${e.message}`);
    }
  }

  function formatDate(date) {
    return new Date(date).toLocaleDateString();
  }

  async function copyHisId(hisId) {
    try {
      await navigator.clipboard.writeText(hisId);
      copiedId = hisId;
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
      <h2>Patients</h2>
      <p class="subtitle">{$patients.length} patient{$patients.length !== 1 ? 's' : ''} registered</p>
    </div>
  </div>

  {#if $error}
    <div class="error">{$error}</div>
  {/if}

  {#if $loading && $patients.length === 0}
    <p class="loading">Loading...</p>
  {:else if $patients.length === 0}
    <div class="empty-state">
      <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
        <circle cx="9" cy="7" r="4"></circle>
        <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
        <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
      </svg>
      <h3>No patients yet</h3>
      <p>Add your first patient above</p>
    </div>
  {:else}
    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th>
              <div class="th-content">
                <span>ID</span>
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
                <span>HIS ID</span>
              </div>
            </th>
            <th>
              <div class="th-content">
                <span>Actions</span>
              </div>
            </th>
          </tr>
        </thead>
        <tbody>
          {#each $patients as patient, i}
            <tr style="animation-delay: {i * 50}ms">
              <td>
                <div class="id-badge">{patient.id}</div>
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
              <td>
                {#if patient.his_patient_id}
                  <button
                    class="his-id-badge"
                    class:copied={copiedId === patient.his_patient_id}
                    on:click={() => copyHisId(patient.his_patient_id)}
                    title="Click to copy HIS ID"
                  >
                    {#if copiedId === patient.his_patient_id}
                      Copied
                    {:else}
                      {patient.his_patient_id.slice(0, 8)}
                    {/if}
                  </button>
                {:else}
                  <span class="pending">Syncing...</span>
                {/if}
              </td>
              <td>
                <button
                  class="delete-btn"
                  on:click={() => handleDelete(patient.id)}
                >
                  Delete
                </button>
              </td>
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

  .id-badge {
    font-family: 'Courier New', monospace;
    font-size: 1rem;
    font-weight: 700;
    color: var(--text);
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

  .his-id-badge {
    font-family: 'Courier New', monospace;
    font-size: 0.75rem;
    font-weight: 600;
    color: var(--success);
    background: linear-gradient(135deg, rgba(16, 185, 129, 0.1) 0%, rgba(5, 150, 105, 0.1) 100%);
    padding: 0.25rem 0.625rem;
    border-radius: 4px;
    border: 1px solid rgba(16, 185, 129, 0.2);
    display: inline-block;
    cursor: pointer;
    transition: all 0.2s ease;
    outline: none;
    min-width: 80px;
    text-align: center;
  }

  .his-id-badge:hover {
    background: linear-gradient(135deg, rgba(16, 185, 129, 0.2) 0%, rgba(5, 150, 105, 0.2) 100%);
    border-color: rgba(16, 185, 129, 0.4);
    transform: translateY(-1px);
    box-shadow: 0 2px 4px rgba(16, 185, 129, 0.15);
  }

  .his-id-badge:active {
    transform: translateY(0);
    box-shadow: none;
  }

  .his-id-badge:focus-visible {
    box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.2);
  }

  .his-id-badge.copied {
    pointer-events: none;
  }

  .pending {
    color: var(--text-light);
    font-size: 0.875rem;
    font-style: italic;
  }

  .delete-btn {
    padding: 0.375rem 0.875rem;
    background: transparent;
    color: var(--error);
    border: 1px solid currentColor;
    font-size: 0.8125rem;
  }

  .delete-btn:hover {
    background: var(--error);
    color: white;
  }

  .loading {
    text-align: center;
    padding: 3rem 2rem;
    color: var(--text-light);
    font-size: 0.9375rem;
  }

  .error {
    color: var(--error);
    padding: 1rem;
    background: #fef2f2;
    border-radius: var(--radius-sm);
    margin-bottom: 1rem;
  }
</style>
