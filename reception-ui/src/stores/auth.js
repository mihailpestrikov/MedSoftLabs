import { writable } from 'svelte/store';

export const accessToken = writable(null);
export const user = writable(null);
export const isAuthenticated = writable(false);

export function setAuth(token, username) {
  accessToken.set(token);
  user.set({ username });
  isAuthenticated.set(true);
  localStorage.setItem('username', username);
}

export function clearAuth() {
  accessToken.set(null);
  user.set(null);
  isAuthenticated.set(false);
  localStorage.removeItem('username');
}
