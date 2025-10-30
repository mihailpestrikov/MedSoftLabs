<script>
  import { onMount } from 'svelte';
  import { createEncounter, getPractitioners } from '../services/api.js';
  import { patients } from '../stores/patients.js';
  import { practitioners, loadingPractitioners, practitionerError } from '../stores/practitioners.js';

  let selectedPatientId = '';
  let selectedPractitionerId = '';
  let startTime = '';
  let error = '';
  let success = '';
  let loading = false;
  let patientSearch = '';
  let showPatientDropdown = false;

  $: filteredPatients = $patients.filter(p => {
    const searchLower = patientSearch.toLowerCase();
    const fullName = `${p.last_name} ${p.first_name} ${p.middle_name || ''}`.toLowerCase();
    return fullName.includes(searchLower);
  });

  onMount(async () => {
    await loadPractitioners();

    const now = new Date();
    const year = now.getFullYear();
    const month = String(now.getMonth() + 1).padStart(2, '0');
    const day = String(now.getDate()).padStart(2, '0');
    const hours = String(now.getHours()).padStart(2, '0');
    const minutes = String(now.getMinutes()).padStart(2, '0');
    startTime = `${year}-${month}-${day}T${hours}:${minutes}`;
  });

  async function loadPractitioners() {
    loadingPractitioners.set(true);
    practitionerError.set(null);

    try {
      const data = await getPractitioners();
      practitioners.set(data);
    } catch (e) {
      practitionerError.set(e.message);
    } finally {
      loadingPractitioners.set(false);
    }
  }

  function selectPatient(patient) {
    selectedPatientId = patient.id;
    patientSearch = `${patient.last_name} ${patient.first_name}${patient.middle_name ? ' ' + patient.middle_name : ''}`;
    showPatientDropdown = false;
  }

  function handlePatientInput() {
    showPatientDropdown = true;
    selectedPatientId = '';
  }

  async function handleSubmit() {
    error = '';
    success = '';

    if (!selectedPatientId) {
      error = 'Please select a patient';
      return;
    }

    if (!selectedPractitionerId) {
      error = 'Please select a doctor';
      return;
    }

    if (!startTime) {
      error = 'Please select date and time';
      return;
    }

    loading = true;

    try {
      const isoTime = new Date(startTime).toISOString();
      await createEncounter(selectedPatientId, selectedPractitionerId, isoTime);

      // Encounter will be added via WebSocket notification, no need to manually update store

      selectedPatientId = '';
      selectedPractitionerId = '';
      patientSearch = '';
      const now = new Date();
      const year = now.getFullYear();
      const month = String(now.getMonth() + 1).padStart(2, '0');
      const day = String(now.getDate()).padStart(2, '0');
      const hours = String(now.getHours()).padStart(2, '0');
      const minutes = String(now.getMinutes()).padStart(2, '0');
      startTime = `${year}-${month}-${day}T${hours}:${minutes}`;

      success = 'Visit registered successfully';
      setTimeout(() => success = '', 3000);
    } catch (e) {
      error = e.message;
    } finally {
      loading = false;
    }
  }
</script>

<div class="card">
  <h2>Register Visit</h2>

  <form on:submit|preventDefault={handleSubmit}>
    <div class="form-grid">
      <div class="field patient-autocomplete">
        <label for="patient">Patient</label>
        <input
          id="patient"
          type="text"
          bind:value={patientSearch}
          on:input={handlePatientInput}
          on:focus={() => showPatientDropdown = true}
          placeholder="Search patient..."
          disabled={loading}
          required
          autocomplete="off"
        />
        {#if showPatientDropdown && patientSearch && filteredPatients.length > 0}
          <div class="dropdown">
            {#each filteredPatients.slice(0, 5) as patient}
              <button
                type="button"
                class="dropdown-item"
                on:click={() => selectPatient(patient)}
              >
                <div class="patient-info">
                  <span class="patient-name">{patient.last_name} {patient.first_name}</span>
                  {#if patient.middle_name}
                    <span class="patient-middle">{patient.middle_name}</span>
                  {/if}
                </div>
                <span class="patient-dob">{new Date(patient.date_of_birth).toLocaleDateString()}</span>
              </button>
            {/each}
          </div>
        {/if}
      </div>

      <div class="field">
        <label for="practitioner">Doctor</label>
        <select
          id="practitioner"
          bind:value={selectedPractitionerId}
          disabled={loading || $loadingPractitioners}
          required
        >
          <option value="">Select doctor...</option>
          {#each $practitioners as practitioner}
            <option value={practitioner.ID}>
              {practitioner.LastName} {practitioner.FirstName} - {practitioner.Specialization}
            </option>
          {/each}
        </select>
      </div>

      <div class="field">
        <label for="startTime">Date & Time</label>
        <input
          id="startTime"
          type="datetime-local"
          bind:value={startTime}
          disabled={loading}
          required
        />
      </div>
    </div>

    {#if error}
      <div class="error">{error}</div>
    {/if}

    {#if success}
      <div class="success">{success}</div>
    {/if}

    {#if $practitionerError}
      <div class="error">Failed to load doctors: {$practitionerError}</div>
    {/if}

    <button type="submit" class="btn-primary" disabled={loading || $loadingPractitioners}>
      {loading ? 'Registering...' : 'Register Visit'}
    </button>
  </form>
</div>

<svelte:window on:click={() => showPatientDropdown = false} />

<style>
  .card {
    background: var(--card-bg);
    border-radius: var(--radius);
    padding: 1.75rem;
    border: 1px solid var(--border);
    box-shadow: var(--shadow-sm);
    min-width: 600px;
  }

  h2 {
    color: var(--text);
    font-size: 1.25rem;
    font-weight: 600;
    margin-bottom: 1.5rem;
    letter-spacing: -0.01em;
  }

  .form-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    position: relative;
  }

  .field.patient-autocomplete {
    grid-column: 1 / -1;
  }

  label {
    color: var(--text-secondary);
    font-size: 0.875rem;
    font-weight: 500;
  }

  input, select {
    padding: 0.75rem;
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    font-size: 0.9375rem;
    color: var(--text);
    background: white;
    transition: all 0.2s ease;
  }

  input:focus, select:focus {
    outline: none;
    border-color: var(--primary);
    box-shadow: 0 0 0 3px rgba(20, 184, 166, 0.1);
  }

  input:disabled, select:disabled {
    background: var(--hover);
    cursor: not-allowed;
    opacity: 0.6;
  }

  .dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background: white;
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    margin-top: 0.25rem;
    box-shadow: var(--shadow-md);
    z-index: 100;
    max-height: 300px;
    overflow-y: auto;
  }

  .dropdown-item {
    width: 100%;
    padding: 0.75rem 1rem;
    border: none;
    background: white;
    text-align: left;
    cursor: pointer;
    transition: background 0.15s ease;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .dropdown-item:hover {
    background: var(--hover);
  }

  .dropdown-item:not(:last-child) {
    border-bottom: 1px solid var(--border);
  }

  .patient-info {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .patient-name {
    font-weight: 500;
    color: var(--text);
  }

  .patient-middle {
    font-size: 0.8125rem;
    color: var(--text-light);
  }

  .patient-dob {
    font-size: 0.8125rem;
    color: var(--text-secondary);
  }

  .btn-primary {
    width: 100%;
    padding: 0.875rem;
    background: linear-gradient(135deg, var(--primary) 0%, var(--accent) 100%);
    color: white;
    border: none;
    border-radius: var(--radius-sm);
    font-size: 0.9375rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .btn-primary:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 6px 16px rgba(20, 184, 166, 0.3);
  }

  .btn-primary:disabled {
    opacity: 0.6;
    cursor: not-allowed;
    transform: none;
  }

  .error {
    color: var(--error);
    background: #fef2f2;
    padding: 0.75rem;
    border-radius: var(--radius-sm);
    font-size: 0.875rem;
    margin-bottom: 1rem;
  }

  .success {
    color: var(--success);
    background: #f0fdf4;
    padding: 0.75rem;
    border-radius: var(--radius-sm);
    font-size: 0.875rem;
    margin-bottom: 1rem;
  }
</style>
