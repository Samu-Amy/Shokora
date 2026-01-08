<script lang="ts">
	import { onMount } from "svelte";

  // TODO: ricorda di fare fetch in onMount

  let productId = $state<number | undefined>();

  onMount(async () => {
    const res = await fetch("/api/v1/menu/products/12", {
      method: "GET",
      headers: {
        "Content-Type": "application/json"
      }
    });

    if (!res.ok) {
      const text = await res.text
      console.log(`Error: ${res.status}, ${text}`);
    }

    let product = await res.json();
    productId = product.product;
  });

</script>

<svelte:head>
  <title>Il nostro menu | Shokora</title>
</svelte:head>

<div class="min-w-full min-h-screen flex flex-col justify-center items-center">
  <h1>Menu</h1>
  <p>Test page (data)</p>
  {#if productId}
    <p>Id prodotto: {productId}</p>
  {:else}
    <p>Loading...</p>
  {/if}
</div>