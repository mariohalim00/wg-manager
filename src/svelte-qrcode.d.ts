declare module 'svelte-qrcode' {
	import type { SvelteComponent } from 'svelte';

	export interface QRCodeProps {
		value: string;
		size?: number;
		level?: 'L' | 'M' | 'Q' | 'H';
		bgColor?: string;
		fgColor?: string;
	}

	export default class QRCode extends SvelteComponent<QRCodeProps> {}
}
