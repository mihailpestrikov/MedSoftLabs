import { writable } from 'svelte/store';

export const encounters = writable([]);
export const loadingEncounters = writable(false);
export const encounterError = writable(null);
