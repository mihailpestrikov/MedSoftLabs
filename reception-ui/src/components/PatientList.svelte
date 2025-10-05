<script>
  import { onMount, onDestroy } from 'svelte';
  import { patients, loading, error } from '../stores/patients.js';
  import { getPatients, deletePatient } from '../services/api.js';

  let refreshInterval;

  onMount(async () => {
    await loadPatients();
    refreshInterval = setInterval(loadPatients, 5000);
  });

  onDestroy(() => {
    if (refreshInterval) clearInterval(refreshInterval);
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

  function getStatusColor(hisId) {
    return hisId ? 'var(--success)' : 'var(--text-light)';
  }
</script>

<div class="card">
  <h2>Patients</h2>

  {#if $error}
    <div class="error">{$error}</div>
  {/if}

  {#if $loading && $patients.length === 0}
    <p>Loading...</p>
  {:else if $patients.length === 0}
    <p class="empty">No patients yet</p>
  {:else}
    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Date of Birth</th>
            <th>Gender</th>
            <th>HIS ID</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each $patients as patient}
            <tr>
              <td>{patient.id}</td>
              <td>{patient.last_name} {patient.first_name} {patient.middle_name || ''}</td>
              <td>{formatDate(patient.date_of_birth)}</td>
              <td>{patient.gender}</td>
              <td style="color: {getStatusColor(patient.his_patient_id)}">
                {patient.his_patient_id || 'Pending...'}
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
  h2 {
    margin-bottom: 1.75rem;
    color: var(--text);
    font-size: 1.125rem;
    font-weight: 600;
    letter-spacing: -0.01em;
  }

  .empty {
    text-align: center;
    color: var(--text-light);
    padding: 3rem 2rem;
    font-size: 0.9375rem;
  }

  .table-container {
    overflow-x: auto;
    border-radius: var(--radius);
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
</style>
