<script lang="ts">
	import { onMount } from "svelte";
	import type { PageData } from "./$types";
	import type { FetchStatus } from "$lib/types";

  let { token }: PageData = $props();

  let status = $state<FetchStatus>("loading");

  onMount(async () => {
    try {
      const res = await fetch(`/api/v1/auth/verify-email/${token}`);
      
      console.log(res); // TODO: togli
      
      // TODO: gestisci meglio
      if (!res.ok) {
        status = "error";
      }
      
      status = "success";
    } catch (err) {
      // TODO: gestisci
      console.log(err)
      status = "error";
    }
  });
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