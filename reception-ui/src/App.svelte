<script>
  import { onMount } from 'svelte';
  import { isAuthenticated } from './stores/auth.js';
  import { checkAuth } from './services/api.js';
  import Login from './components/Login.svelte';
  import Dashboard from './components/Dashboard.svelte';
  import './styles/global.css';

  let loading = true;

  onMount(async () => {
    await checkAuth();
    loading = false;
  });
</script>

{#if loading}
  <div class="loading">Loading...</div>
{:else if $isAuthenticated}
  <Dashboard />
{:else}
  <Login />
{/if}

<style>
  .loading {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    font-size: 1.2rem;
    color: var(--text-light);
  }
</style>
