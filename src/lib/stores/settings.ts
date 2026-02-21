// Settings store
import { writable } from 'svelte/store';
import type { GlobalSettings } from '../types/settings';
import * as settingsAPI from '../api/settings';
import { addNotification } from './notifications';

function createSettingsStore() {
	const { subscribe, set, update } = writable<GlobalSettings | null>(null);

	return {
		subscribe,
		set,
		update,

		/**
		 * Load settings from API
		 */
		async load(): Promise<void> {
			const response = await settingsAPI.getSettings();

			if (response.error) {
				addNotification({
					type: 'error',
					message: `Failed to load settings: ${response.error.error}`,
					duration: 5000
				});
				return;
			}

			if (response.data) {
				set(response.data);
			}
		},

		/**
		 * Save settings to API
		 */
		async save(settings: GlobalSettings): Promise<boolean> {
			const response = await settingsAPI.updateSettings(settings);

			if (response.error) {
				addNotification({
					type: 'error',
					message: `Failed to save settings: ${response.error.error}`,
					duration: 5000
				});
				return false;
			}

			set(settings);
			addNotification({
				type: 'success',
				message: 'Settings saved successfully',
				duration: 3000
			});
			return true;
		}
	};
}

export const settings = createSettingsStore();
