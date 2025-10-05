import { writable } from 'svelte/store';

export const patients = writable([]);
export const loading = writable(false);
export const error = writable(null);
