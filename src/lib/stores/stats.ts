// Stats store for interface statistics
import { writable } from 'svelte/store';
import type { InterfaceStats, StatsHistoryItem } from '../types/stats';
import * as statsAPI from '../api/stats';
import { addNotification } from './notifications';

/**
 * Stats writable store
 */
function createStatsStore() {
	const currentStats = writable<InterfaceStats | null>(null);
	const statsHistory = writable<StatsHistoryItem[]>([]);

	return {
		subscribe: currentStats.subscribe,
		history: { subscribe: statsHistory.subscribe },
		set: currentStats.set,

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
				currentStats.set(response.data);
			}
		},

		/**
		 * Load statistics history from API
		 */
		async loadHistory(): Promise<void> {
			const response = await statsAPI.getStatsHistory();

			if (response.error) {
				console.error('Failed to load stats history:', response.error);
				return;
			}

			if (response.data) {
				statsHistory.set(response.data);
			}
		}
	};
}

export const stats = createStatsStore();
