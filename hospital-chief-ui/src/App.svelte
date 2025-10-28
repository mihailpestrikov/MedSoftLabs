<script>
  import PatientTable from './components/PatientTable.svelte';
  import PractitionerManager from './components/PractitionerManager.svelte';
  import './styles/global.css';

  let activeTab = 'patients';
</script>

<div class="app">
  <header>
    <div class="container header-content">
      <h1>Hospital Information System</h1>
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
        class:active={activeTab === 'practitioners'}
        on:click={() => activeTab = 'practitioners'}
      >
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
          <circle cx="12" cy="7" r="4"></circle>
        </svg>
        Practitioners
      </button>
    </div>

    <div class="tab-content">
      <div class="tab-panel" class:hidden={activeTab !== 'patients'}>
        <PatientTable />
      </div>
      <div class="tab-panel" class:hidden={activeTab !== 'practitioners'}>
        <PractitionerManager />
      </div>
    </div>
  </main>
</div>

<style>
  .app {
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

  .tabs {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 2rem;
    border-bottom: 2px solid var(--border);
    padding-bottom: 0;
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
    display: none;
  }
</style>
