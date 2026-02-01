import { writable, get as getStoreValue } from 'svelte/store';
import type { Peer, PeerFormData, PeerCreateResponse, PeerUpdateRequest } from '../types/peer';
import * as peersAPI from '../api/peers';
import { addNotification } from './notifications';

// HANDSHAKE_TIMEOUT: peer is online if lastHandshake within 120 seconds
const HANDSHAKE_TIMEOUT = 120000; // 120 seconds in milliseconds
const POLLING_INTERVAL = 3000; // 3 seconds

/**
 * Determine peer status based on lastHandshake timestamp
 */
function deriveStatus(lastHandshake: string): 'online' | 'offline' {
	if (lastHandshake === '0' || !lastHandshake) {
		return 'offline';
	}

	const handshakeTime = new Date(lastHandshake).getTime();
	const now = Date.now();
	const diff = now - handshakeTime;

	return diff <= HANDSHAKE_TIMEOUT ? 'online' : 'offline';
}

/**
 * Peers writable store
 */
function createPeersStore() {
	const { subscribe, set, update } = writable<Peer[]>([]);
	let pollingTimer: ReturnType<typeof setInterval> | null = null;

	return {
		subscribe,
		set,

		/**
		 * Load all peers from API
		 */
		async load(silent = false): Promise<void> {
			const response = await peersAPI.listPeers();

			if (response.error) {
				if (!silent) {
					addNotification({
						type: 'error',
						message: `Failed to load peers: ${response.error.error}`,
						duration: 5000
					});
				}
				return;
			}

			if (response.data) {
				// Derive status for each peer
				const peersWithStatus = response.data.map((peer) => ({
					...peer,
					status: deriveStatus(peer.lastHandshake)
				}));
				set(peersWithStatus);
			}
		},

		/**
		 * Add a new peer
		 */
		async add(data: PeerFormData): Promise<PeerCreateResponse | null> {
			const response = await peersAPI.addPeer(data);

			if (response.error) {
				addNotification({
					type: 'error',
					message: `Failed to add peer: ${response.error.error}`,
					duration: 5000
				});
				return null;
			}

			if (response.data) {
				addNotification({
					type: 'success',
					message: `Peer "${data.name}" added successfully!`,
					duration: 3000
				});

				// Reload peers to get updated list
				await this.load();
				return response.data;
			}

			return null;
		},

		/**
		 * Update an existing peer
		 */
		async update(peerId: string, data: PeerUpdateRequest): Promise<boolean> {
			const response = await peersAPI.updatePeer(peerId, data);

			if (response.error) {
				addNotification({
					type: 'error',
					message: `Failed to update peer: ${response.error.error}`,
					duration: 5000
				});
				return false;
			}

			if (response.data) {
				addNotification({
					type: 'success',
					message: `Peer "${data.name || 'updated'}" saved successfully!`,
					duration: 3000
				});

				// Update local state
				const updatedPeer = {
					...response.data,
					status: deriveStatus(response.data.lastHandshake)
				};
				update((peers) => peers.map((p) => (p.id === peerId ? updatedPeer : p)));
				return true;
			}

			return false;
		},

		/**
		 * Regenerate keys for a peer
		 */
		async regenerateKeys(peerId: string): Promise<PeerCreateResponse | null> {
			const response = await peersAPI.regenerateKeys(peerId);

			if (response.error) {
				addNotification({
					type: 'error',
					message: `Failed to regenerate keys: ${response.error.error}`,
					duration: 5000
				});
				return null;
			}

			if (response.data) {
				addNotification({
					type: 'success',
					message: `Keys regenerated for peer. Please save the new configuration.`,
					duration: 5000
				});
				await this.load(); // Reload to ensure all stats/IDs are fresh
				return response.data;
			}

			return null;
		},

		async getConfig(peerId: string): Promise<string | null> {
			try {
				const response = await fetch(peersAPI.getConfigUrl(peerId));
				if (!response.ok) {
					const errorText = await response.text();
					addNotification({
						type: 'error',
						message: `Failed to fetch config: ${errorText || response.statusText}`,
						duration: 5000
					});
					return null;
				}
				return await response.text();
			} catch (err) {
				addNotification({
					type: 'error',
					message: `Network error while fetching config`,
					duration: 5000
				});
				return null;
			}
		},

		/**
		 * Remove a peer by ID
		 */
		async remove(peerId: string, peerName: string): Promise<boolean> {
			const response = await peersAPI.removePeer(peerId);

			if (response.error) {
				addNotification({
					type: 'error',
					message: `Failed to remove peer: ${response.error.error}`,
					duration: 5000
				});
				return false;
			}

			addNotification({
				type: 'success',
				message: `Peer "${peerName}" removed successfully!`,
				duration: 3000
			});

			// Remove from local state
			update((peers) => peers.filter((p) => p.id !== peerId));
			return true;
		},

		/**
		 * Start periodic polling for stats/status
		 */
		startPolling() {
			if (pollingTimer) return;
			// Initial load done by components usually, but start interval here
			pollingTimer = setInterval(() => {
				this.load(true); // silent load for polling
			}, POLLING_INTERVAL);
		},

		/**
		 * Stop periodic polling
		 */
		stopPolling() {
			if (pollingTimer) {
				clearInterval(pollingTimer);
				pollingTimer = null;
			}
		}
	};
}

export const peers = createPeersStore();
