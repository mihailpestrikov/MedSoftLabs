import { writable } from 'svelte/store';

export const practitioners = writable([]);
export const loadingPractitioners = writable(false);
export const practitionerError = writable(null);
