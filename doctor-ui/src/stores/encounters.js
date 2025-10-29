import { writable } from 'svelte/store';

export const encounters = writable([]);
export const loading = writable(false);
export const error = writable(null);
