<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();

  //TODO: add metadata
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<nav class="py-4 bg-slate-800 text-white flex gap-4 justify-center items-center">
  <a class="hover:text-orange-400" href="/">Home</a>
  <a class="hover:text-orange-400" href="/about">About</a>
  <a class="hover:text-orange-400" href="/contacts">Contacts</a>
  <a class="hover:text-orange-400" href="/menu">Menu</a>
</nav>
<main>
  {@render children()}
</main>
