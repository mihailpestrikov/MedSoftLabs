<script>
  import { user } from '../stores/auth.js';
  import { logout } from '../services/api.js';
  import PatientForm from './PatientForm.svelte';
  import PatientList from './PatientList.svelte';

  async function handleLogout() {
    try {
      await logout();
    } catch (e) {
      console.error('Logout failed:', e);
    }
  }
</script>

<div class="dashboard">
  <header>
    <div class="container header-content">
      <h1>Reception Dashboard</h1>
      <div class="user-info">
        <span>Welcome, {$user?.username}</span>
        <button on:click={handleLogout} class="logout-btn">Logout</button>
      </div>
    </div>
  </header>

  <main class="container">
    <div class="section">
      <PatientForm />
    </div>

    <div class="section">
      <PatientList />
    </div>
  </main>
</div>

<style>
  .dashboard {
    min-height: 100vh;
  }

  header {
    background: var(--card-bg);
    border-bottom: 1px solid var(--border);
    padding: 1.25rem 0;
    margin-bottom: 2.5rem;
    width: 100%;
    box-shadow: var(--shadow-sm);
  }

  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    max-width: 1600px;
    margin: 0 auto;
    padding: 0 2rem;
  }

  h1 {
    background: linear-gradient(135deg, var(--primary) 0%, var(--accent) 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    font-size: 1.5rem;
    font-weight: 700;
    letter-spacing: -0.02em;
  }

  .user-info {
    display: flex;
    align-items: center;
    gap: 1.25rem;
  }

  .user-info span {
    color: var(--text-secondary);
    font-size: 0.875rem;
    font-weight: 500;
  }

  .logout-btn {
    padding: 0.5rem 1rem;
    background: transparent;
    color: var(--text-secondary);
    border: 1px solid var(--border);
    font-size: 0.875rem;
  }

  .logout-btn:hover {
    background: var(--hover);
    border-color: var(--text-light);
    color: var(--text);
    box-shadow: none;
    transform: none;
  }

  .section {
    margin-bottom: 2rem;
  }
</style>
