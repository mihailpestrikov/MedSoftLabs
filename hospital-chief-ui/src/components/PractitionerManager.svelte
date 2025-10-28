<script>
  import { onMount } from 'svelte';
  import { practitioners, loading, error } from '../stores/practitioners.js';
  import { getPractitioners, createPractitioner } from '../services/api.js';

  let showForm = false;
  let formData = {
    firstName: '',
    lastName: '',
    middleName: '',
    specialization: ''
  };
  let submitting = false;

  onMount(async () => {
    await loadPractitioners();
  });

  async function loadPractitioners() {
    loading.set(true);
    error.set(null);

    try {
      const data = await getPractitioners();
      practitioners.set(data);
    } catch (e) {
      error.set(e.message);
    } finally {
      loading.set(false);
    }
  }

  let copiedId = null;

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

  async function handleSubmit() {
    if (!formData.firstName || !formData.lastName || !formData.specialization) {
      return;
    }

    submitting = true;
    error.set(null);

    try {
      const newPractitioner = await createPractitioner(
        formData.firstName,
        formData.lastName,
        formData.middleName,
        formData.specialization
      );

      practitioners.update(p => [newPractitioner, ...p]);

      // Reset form
      formData = {
        firstName: '',
        lastName: '',
        middleName: '',
        specialization: ''
      };
      showForm = false;
    } catch (e) {
      error.set(e.message);
    } finally {
      submitting = false;
    }
  }
</script>

<div class="card">
  <div class="header-row">
    <div>
      <h2>Medical Practitioners</h2>
      <p class="subtitle">{$practitioners.length} practitioner{$practitioners.length !== 1 ? 's' : ''} registered</p>
    </div>
    <button class="add-button" on:click={() => showForm = !showForm}>
      {showForm ? '✕ Cancel' : '+ Add Practitioner'}
    </button>
  </div>

  {#if showForm}
    <div class="form-container">
      <h3>Add New Practitioner</h3>
      <form on:submit|preventDefault={handleSubmit}>
        <div class="form-row">
          <div class="form-group">
            <label for="firstName">First Name *</label>
            <input
              id="firstName"
              type="text"
              bind:value={formData.firstName}
              required
              placeholder="Иван"
            />
          </div>
          <div class="form-group">
            <label for="lastName">Last Name *</label>
            <input
              id="lastName"
              type="text"
              bind:value={formData.lastName}
              required
              placeholder="Петров"
            />
          </div>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label for="middleName">Middle Name</label>
            <input
              id="middleName"
              type="text"
              bind:value={formData.middleName}
              placeholder="Сергеевич"
            />
          </div>
          <div class="form-group">
            <label for="specialization">Specialization *</label>
            <input
              id="specialization"
              type="text"
              bind:value={formData.specialization}
              required
              placeholder="Терапевт"
            />
          </div>
        </div>
        <button type="submit" class="submit-button" disabled={submitting}>
          {submitting ? 'Creating...' : 'Create Practitioner'}
        </button>
      </form>
    </div>
  {/if}

  {#if $error}
    <div class="error">{$error}</div>
  {/if}

  {#if $loading}
    <p class="loading">Loading practitioners...</p>
  {:else if $practitioners.length === 0}
    <div class="empty-state">
      <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
        <circle cx="12" cy="7" r="4"></circle>
      </svg>
      <h3>No practitioners yet</h3>
      <p>Add your first medical practitioner to get started</p>
    </div>
  {:else}
    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th>
              <div class="th-content">
                <span>Practitioner ID</span>
              </div>
            </th>
            <th>
              <div class="th-content">
                <span>Name</span>
              </div>
            </th>
            <th>
              <div class="th-content">
                <span>Specialization</span>
              </div>
            </th>
          </tr>
        </thead>
        <tbody>
          {#each $practitioners as practitioner, i}
            <tr style="animation-delay: {i * 50}ms">
              <td>
                <div class="uuid-cell">
                  <button
                    class="uuid-badge"
                    class:copied={copiedId === practitioner.id}
                    on:click={() => copyUUID(practitioner.id)}
                    title="Click to copy full UUID"
                  >
                    {#if copiedId === practitioner.id}
                      Copied
                    {:else}
                      {practitioner.id.slice(0, 8)}
                    {/if}
                  </button>
                  <span class="uuid-full">{practitioner.id}</span>
                </div>
              </td>
              <td>
                <div class="name-cell">
                  <span class="name">{practitioner.lastName} {practitioner.firstName}</span>
                  {#if practitioner.middleName}
                    <span class="middle-name">{practitioner.middleName}</span>
                  {/if}
                </div>
              </td>
              <td>
                <span class="specialization-badge">{practitioner.specialization}</span>
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

  .add-button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.625rem 1.25rem;
    background: linear-gradient(135deg, var(--primary) 0%, var(--accent) 100%);
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
    box-shadow: 0 4px 12px rgba(20, 184, 166, 0.25);
  }

  .add-button:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 16px rgba(20, 184, 166, 0.3);
  }

  .add-button:active {
    transform: translateY(0);
  }

  .form-container {
    background: linear-gradient(135deg, rgba(20, 184, 166, 0.05) 0%, rgba(6, 182, 212, 0.05) 100%);
    padding: 1.5rem;
    border-radius: 12px;
    margin-bottom: 2rem;
    border: 1px solid rgba(20, 184, 166, 0.1);
  }

  .form-container h3 {
    color: var(--text);
    font-size: 1.125rem;
    font-weight: 600;
    margin-bottom: 1.5rem;
  }

  .form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  label {
    color: var(--text-secondary);
    font-size: 0.875rem;
    font-weight: 500;
  }

  input {
    padding: 0.75rem;
    border: 1px solid var(--border);
    border-radius: 8px;
    font-size: 0.9375rem;
    color: var(--text);
    background: white;
    transition: all 0.2s ease;
  }

  input:focus {
    outline: none;
    border-color: var(--primary);
    box-shadow: 0 0 0 3px rgba(20, 184, 166, 0.1);
  }

  .submit-button {
    width: 100%;
    padding: 0.875rem;
    background: linear-gradient(135deg, var(--primary) 0%, var(--accent) 100%);
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 0.9375rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
    margin-top: 0.5rem;
  }

  .submit-button:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 6px 16px rgba(20, 184, 166, 0.3);
  }

  .submit-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
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

  .specialization-badge {
    display: inline-block;
    padding: 0.375rem 0.875rem;
    border-radius: 6px;
    font-size: 0.875rem;
    font-weight: 500;
    background: linear-gradient(135deg, rgba(139, 92, 246, 0.1) 0%, rgba(168, 85, 247, 0.1) 100%);
    color: #7c3aed;
    border: 1px solid rgba(139, 92, 246, 0.2);
    min-width: 120px;
    text-align: center;
  }

  .error {
    color: var(--error);
    padding: 1rem;
    background: #fef2f2;
    border-radius: var(--radius-sm);
    margin-bottom: 1rem;
  }

  .loading {
    text-align: center;
    padding: 2rem;
    color: var(--text-light);
  }
</style>
