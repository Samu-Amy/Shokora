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

  // TODO: aggiungi api call a "/api/v1/auth/me" per cercare di ottenere lo user (se autenticato), salvarlo nel context e mostrare "accedi" o l'icona utente (?)

  // TODO: https://www.youtube.com/watch?v=ieECVME5ZLU (redirect a protected routes a cui si prova ad accedere senza essere loggati (dopo il login si va a quella pagina))

  // TODO: per la ricerca aggiungere un timer per cercare i prodotti (o altro) solo dopo tot (es. 300ms) che l'utente ha smesso di scrivere
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
