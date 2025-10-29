<script>
  import { user } from '../stores/auth.js';
  import { logout } from '../services/api.js';
  import PatientForm from './PatientForm.svelte';
  import PatientList from './PatientList.svelte';
  import EncounterForm from './EncounterForm.svelte';
  import EncounterList from './EncounterList.svelte';

  let activeTab = 'patients';

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
    <div class="tabs">
      <button
        class="tab"
        class:active={activeTab === 'patients'}
        on:click={() => activeTab = 'patients'}
      >
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
          <circle cx="9" cy="7" r="4"></circle>
          <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
          <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
        </svg>
        Patients
      </button>
      <button
        class="tab"
        class:active={activeTab === 'visits'}
        on:click={() => activeTab = 'visits'}
      >
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
        </svg>
        Visits
      </button>
    </div>

    <div class="tab-content">
      <div class="tab-panel" class:hidden={activeTab !== 'patients'}>
        <div class="section">
          <PatientForm />
        </div>
        <div class="section">
          <PatientList />
        </div>
      </div>

      <div class="tab-panel" class:hidden={activeTab !== 'visits'}>
        <div class="section">
          <EncounterForm />
        </div>
        <div class="section">
          <EncounterList />
        </div>
      </div>
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

  .tabs {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 2rem;
    border-bottom: 2px solid var(--border);
    padding-bottom: 0;
    position: relative;
    z-index: 200;
  }

  .tab {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.875rem 1.5rem;
    background: transparent;
    border: none;
    border-bottom: 3px solid transparent;
    color: var(--text-light);
    font-size: 0.9375rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    position: relative;
    margin-bottom: -2px;
    outline: none;
    pointer-events: auto;
    z-index: 1000;
  }

  .tab svg {
    transition: all 0.2s ease;
  }

  .tab:hover {
    color: var(--text);
    background: rgba(20, 184, 166, 0.05);
  }

  .tab.active {
    color: var(--primary);
    border-bottom-color: var(--primary);
    font-weight: 600;
  }

  .tab.active svg {
    stroke: var(--primary);
  }

  .tab-content {
    position: relative;
  }

  .tab-panel {
    transition: opacity 0.2s ease;
  }

  .tab-panel.hidden {
    display: none !important;
    visibility: hidden !important;
    pointer-events: none !important;
  }
</style>
