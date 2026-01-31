// Notification types for toast messages

export interface Notification {
	id: string; // Unique ID for dismissal
	type: 'success' | 'error' | 'warning' | 'info';
	message: string;
	duration?: number; // Auto-dismiss duration (ms), undefined = manual dismiss
}
