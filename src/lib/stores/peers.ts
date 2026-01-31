// Peers store for state management
import { writable } from 'svelte/store';
import type { Peer, PeerFormData, PeerCreateResponse } from '../types/peer';
import * as peersAPI from '../api/peers';
import { addNotification } from './notifications';

// HANDSHAKE_TIMEOUT: peer is online if lastHandshake within 120 seconds
const HANDSHAKE_TIMEOUT = 120000; // 120 seconds in milliseconds

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

	return {
		subscribe,
		set,

		/**
		 * Load all peers from API
		 */
		async load(): Promise<void> {
			const response = await peersAPI.listPeers();

			if (response.error) {
				addNotification({
					type: 'error',
					message: `Failed to load peers: ${response.error.error}`,
					duration: 5000
				});
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
		}
	};
}

export const peers = createPeersStore();
