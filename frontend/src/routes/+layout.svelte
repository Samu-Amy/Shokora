<script lang="ts">
	import './layout.css';

	let { children } = $props();

  // TODO: add metadata e sistema titles
  // TODO: fai funzione per le requests che accetta i metodi giusti (oppure funzioni tipo getUser, postProduct, ecc. che hanno già il metodo "incorporato")
  // TODO: per gli sconti accettare sia percentuale che prezzo scontato (mostrando il prezzo scontato/percentuale risultante) -> nel caso si scelga il prezzo scontato si mostra l'equivalente in percentuale e viceversa, sul db si salva in versione percentuale - ATTENZIONE: se il prezzo base cambia, lo sconto rimane uguale (e quindi il prezzo scontato cambia -> al cambio di prezzo mostra anche lo sconto (e il prezzo scontato) per poterlo aggiustare)
  // TODO: usa Tanstack Query (?) e una libreria per la validazione degli input
  // TODO: come evitare che qualcuno senza autorizaazione possa accedere all'admin panel (alle pagine del pannello (non possono in ogni caso ottenere o modificare i dati dal backend, ma devono essere reindirizzati alle pagine accessibili (fetch on mount -> backend invia una richiesta di redirect (?))))
  // TODO: update dati db -> attenzione a null/undefined vs "" -> la prima significa che non si modificano, la seconda che si cancella il contenuto
  // TODO: aggiungi debounce (per evitare di creare cose più volte, inviare troppe email, ecc.)
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
