// Stats store for interface statistics
import { writable } from 'svelte/store';
import type { InterfaceStats } from '../types/stats';
import * as statsAPI from '../api/stats';
import { addNotification } from './notifications';

/**
 * Stats writable store
 */
function createStatsStore() {
	const { subscribe, set } = writable<InterfaceStats | null>(null);

	return {
		subscribe,
		set,

		/**
		 * Load interface statistics from API
		 */
		async load(): Promise<void> {
			const response = await statsAPI.getStats();

			if (response.error) {
				addNotification({
					type: 'error',
					message: `Failed to load stats: ${response.error.error}`,
					duration: 5000
				});
				return;
			}

			if (response.data) {
				set(response.data);
			}
		}
	};
}

export const stats = createStatsStore();
