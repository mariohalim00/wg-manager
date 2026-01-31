import adapter from '@sveltejs/adapter-static';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: {
		// adapter-static configuration for SPA mode
		adapter: adapter({
			fallback: 'index.html', // Enable SPA mode
			strict: true
		})
	}
};

export default config;
