<script>
  import { createPatient } from '../services/api.js';
  import { patients } from '../stores/patients.js';

  let firstName = '';
  let lastName = '';
  let middleName = '';
  let dateOfBirth = '';
  let gender = 'male';
  let error = '';
  let success = '';
  let loading = false;

  const today = new Date().toISOString().split('T')[0];
  const minDate = new Date(new Date().getFullYear() - 120, 0, 1).toISOString().split('T')[0];

  function validateName(value) {
    return value.replace(/[0-9]/g, '');
  }

  function handleNameInput(event, field) {
    const cleaned = validateName(event.target.value);
    if (field === 'firstName') firstName = cleaned;
    if (field === 'lastName') lastName = cleaned;
    if (field === 'middleName') middleName = cleaned;
  }

  async function handleSubmit() {
    error = '';
    success = '';
    loading = true;

    try {
      await createPatient({
        first_name: firstName,
        last_name: lastName,
        middle_name: middleName || undefined,
        date_of_birth: dateOfBirth,
        gender,
      });

      firstName = '';
      lastName = '';
      middleName = '';
      dateOfBirth = '';
      gender = 'male';

      success = 'Patient added successfully';
      setTimeout(() => success = '', 3000);
    } catch (e) {
      error = e.message;
    } finally {
      loading = false;
    }
  }
</script>

<div class="card">
  <h2>Add Patient</h2>

  <form on:submit|preventDefault={handleSubmit}>
    <div class="form-grid">
      <div class="field">
        <label for="lastName">Last Name</label>
        <input
          id="lastName"
          type="text"
          bind:value={lastName}
          on:input={(e) => handleNameInput(e, 'lastName')}
          on:keypress={(e) => /[0-9]/.test(e.key) && e.preventDefault()}
          disabled={loading}
          required
        />
      </div>

      <div class="field">
        <label for="firstName">First Name</label>
        <input
          id="firstName"
          type="text"
          bind:value={firstName}
          on:input={(e) => handleNameInput(e, 'firstName')}
          on:keypress={(e) => /[0-9]/.test(e.key) && e.preventDefault()}
          disabled={loading}
          required
        />
      </div>

      <div class="field">
        <label for="middleName">Middle Name</label>
        <input
          id="middleName"
          type="text"
          bind:value={middleName}
          on:input={(e) => handleNameInput(e, 'middleName')}
          on:keypress={(e) => /[0-9]/.test(e.key) && e.preventDefault()}
          disabled={loading}
        />
      </div>

      <div class="field">
        <label for="dateOfBirth">Date of Birth</label>
        <input
          id="dateOfBirth"
          type="date"
          bind:value={dateOfBirth}
          min={minDate}
          max={today}
          disabled={loading}
          required
        />
      </div>

      <div class="field">
        <label for="gender">Gender</label>
        <select id="gender" bind:value={gender} disabled={loading}>
          <option value="male">Male</option>
          <option value="female">Female</option>
        </select>
      </div>
    </div>

    {#if error}
      <div class="error">{error}</div>
    {/if}

    {#if success}
      <div class="success">{success}</div>
    {/if}

    <button type="submit" disabled={loading}>
      {loading ? 'Adding...' : 'Add Patient'}
    </button>
  </form>
</div>

<style>
  h2 {
    margin-bottom: 1.75rem;
    color: var(--text);
    font-size: 1.125rem;
    font-weight: 600;
    letter-spacing: -0.01em;
  }

  .form-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
    gap: 1.25rem;
    margin-bottom: 1.5rem;
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  label {
    font-weight: 500;
    font-size: 0.8125rem;
    color: var(--text-secondary);
    letter-spacing: 0.01em;
  }
</style>
