import { writable } from 'svelte/store';

export const practitioners = writable([]);
export const loading = writable(false);
export const error = writable(null);
