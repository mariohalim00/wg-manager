import { writable } from 'svelte/store';
import type { Peer } from '$lib/types';

export const peers = writable<Peer[]>([]);
