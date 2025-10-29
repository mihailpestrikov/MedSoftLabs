<script>
  import EncounterList from './components/EncounterList.svelte';
  import './styles/global.css';
</script>

<div class="app">
  <header>
    <div class="header-content">
      <h1>Doctor Dashboard</h1>
    </div>
  </header>

  <main class="container">
    <EncounterList />
  </main>
</div>

<style>
  .app {
    min-height: 100vh;
  }

  header {
    background: var(--card-bg);
    border-bottom: 1px solid var(--border);
    padding: 1.25rem 2rem;
    margin: 2rem auto 2.5rem;
    max-width: 1000px;
    box-shadow: var(--shadow-sm);
  }

  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
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
