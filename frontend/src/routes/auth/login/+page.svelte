<script lang="ts">
  async function loginWithGoogle() {
    try {
      const res = await fetch("/api/v1/auth/google");
      const { data } = await res.json();
      window.location.href = data.url; // Navigate to Google OAuth login
    } catch (err) {
      // TODO: gestisci (?)
      console.log("Errore")
      // status = "error";
    }
  }
</script>

<svelte:head>
  <title>Accedi | Shokora</title>
</svelte:head>

<div class="min-w-full min-h-screen flex flex-col justify-center items-center">
  <h1>Accedi</h1>
  <button onclick={loginWithGoogle}>Accedi con google</button>
</div>