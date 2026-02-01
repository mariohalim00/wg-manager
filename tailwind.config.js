import forms from '@tailwindcss/forms';
import typography from '@tailwindcss/typography';
import daisyui from 'daisyui';
/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			backdropBlur: {
				xs: '2px',
				sm: '4px',
				md: '12px',
				lg: '16px',
				xl: '24px',
				'2xl': '40px'
			},
			colors: {
				primary: '#137fec',
				secondary: '#22c55e',
				glass: {
					bg: 'rgba(255, 255, 255, 0.1)',
					'bg-secondary': 'rgba(255, 255, 255, 0.05)',
					'bg-hover': 'rgba(255, 255, 255, 0.15)',
					'bg-modal': 'rgba(255, 255, 255, 0.2)',
					border: 'rgba(255, 255, 255, 0.1)',
					'border-subtle': 'rgba(255, 255, 255, 0.08)'
				}
			}
		}
	},
	plugins: [forms, typography, daisyui]
};
