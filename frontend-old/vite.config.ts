import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	// TODO: togli in production
	server: {
		host: '0.0.0.0', // Esplicita invece di true
		port: 5173,
		strictPort: false,
		proxy: {
			'/api': {
				target: 'http://localhost:8080', // TODO: da settare poi il reverse proxy su nginx
				changeOrigin: true
			}
		}
	}
});
