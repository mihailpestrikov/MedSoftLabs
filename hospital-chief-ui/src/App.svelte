<script>
  import PatientTable from './components/PatientTable.svelte';
  import './styles/global.css';
</script>

<div class="app">
  <header>
    <div class="container header-content">
      <h1>Hospital Information System</h1>
    </div>
  </header>

  <main class="container">
    <PatientTable />
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
</style>
