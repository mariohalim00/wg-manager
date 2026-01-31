// Notifications store for toast notifications
import { writable } from 'svelte/store';
import type { Notification } from '../types/notification';

/**
 * Notifications writable store
 */
function createNotificationsStore() {
	const { subscribe, update } = writable<Notification[]>([]);

	return {
		subscribe,

		/**
		 * Add a notification to the stack
		 */
		add(notification: Omit<Notification, 'id'>): void {
			const id = crypto.randomUUID();
			const newNotification: Notification = {
				...notification,
				id
			};

			update((notifications) => [...notifications, newNotification]);

			// Auto-dismiss after duration (default 3000ms)
			const duration = notification.duration ?? 3000;
			if (duration > 0) {
				setTimeout(() => {
					this.remove(id);
				}, duration);
			}
		},

		/**
		 * Remove a notification by ID
		 */
		remove(id: string): void {
			update((notifications) => notifications.filter((n) => n.id !== id));
		}
	};
}

export const notifications = createNotificationsStore();

/**
 * Helper function to add notification from other modules
 */
export function addNotification(notification: Omit<Notification, 'id'>): void {
	notifications.add(notification);
}
