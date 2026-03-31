<script lang="ts">
	import { onMount } from "svelte";
	import type { PageData } from "./$types";
	import type { FetchStatus } from "$lib/types";

  let { data }: {data: PageData} = $props();

  let status = $state<FetchStatus>("loading");

  onMount(async () => {
    // TODO: implementa debouncing (magari crea una funzione per gestire le query con metodi, debouncing ed altro)
    try {
      const res = await fetch(
        `/api/v1/auth/verify-email/${data.token}`,
        {
          method: "POST",
        }
      );
      
      console.log(res); // TODO: togli
      
      // TODO: gestisci meglio
      if (!res.ok) {
        status = "error";
      }

      // TODO: fai redirect con messaggio di successo (?)
      
      status = "success";
    } catch (err) {
      // TODO: gestisci
      console.log(err)
      status = "error";
    }
  });

  // TODO: per il reset della password mettere opzione "disconnetti da tutti i dispositivi" -> elimina tutti i Refresh Token (di tutte le sessioni) di quell'utente

  // TODO: apri pagina otp solo se c'è verification_id nel payload (altrimenti non si può verificare con l'otp)
</script>

<svelte:head>
  <title>Verifica email | Shokora</title>
</svelte:head>

<div class="min-w-full min-h-screen flex flex-col justify-center items-center">
  <h1>Verifica dell'email</h1>

  {#if status == "success"}
    <p>Email verificata con successo!</p>
  {:else if status == "loading"}
    <p>Stiamo verificando la tua email...</p>
  {:else}
    <p>Errore durante la verifica dell'email</p>
  {/if}
</div>