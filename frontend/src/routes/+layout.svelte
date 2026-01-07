<script lang="ts">
	import './layout.css';

	let { children } = $props();

  //TODO: add metadata e sistema titles
  //TODO: fai funzione per le requests che accetta i metodi giusti (oppure funzioni tipo getUser, postProduct, ecc. che hanno già il metodo "incorporato")
</script>

<svelte:head>
	<!-- <link rel="icon" href={favicon} /> -->
  <title>Shokora - La pasticceria</title>
</svelte:head>

<nav class="py-4 bg-slate-800 text-white flex gap-4 justify-center items-center">
  <a class="hover:text-orange-400" href="/">Home</a>
  <a class="hover:text-orange-400" href="/menu">Menu</a>
  <a class="hover:text-orange-400" href="/about">About</a>
  <a class="hover:text-orange-400" href="/contacts">Contacts</a>
</nav>
<main>
  {@render children()}
</main>
