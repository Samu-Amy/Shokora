<script lang="ts">
	import type { FetchStatus } from "$lib/types";
	import { onMount } from "svelte";

  let status = $state<FetchStatus>("loading");
  let productId = $state<number | undefined>();

  // TODO: usare load in +page.ts invece che onMount qua (?)
  onMount(async () => {
    try {
      const res = await fetch("/api/v1/menu/products/12", {
        method: "GET",
        headers: {
          "Content-Type": "application/json"
        }
      });
  
      if (!res.ok) {
        const text = await res.text
        console.log(`Error: ${res.status}, ${text}`);
        status = "error";
      }
  
      let product = await res.json();
      productId = product.product;
      status = "success";
    } catch (err) {
      // TODO: gestisci (?)
      console.log("Errore")
      status = "error";
    }
  });

</script>

<svelte:head>
  <title>Il nostro menu | Shokora</title>
</svelte:head>

<div class="min-w-full min-h-screen flex flex-col justify-center items-center">
  <h1>Menu</h1>
  <p>Test page (data)</p>
  
  {#if status == "success"}
    <p>Id prodotto: {productId}</p>
  {:else if status == "loading"}
    <p>Caricamento...</p>
  {:else}
    <p>Errore durante l'ottenimento dei dati</p>
  {/if}
</div>