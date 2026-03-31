<script>
    import { onMount } from 'svelte'
    import { goto } from '$app/navigation'

    onMount(async () => {
      try {
        // Get params from url
        const params = new URLSearchParams(window.location.search)
        const code = params.get('code')
        const state = params.get('state')

        if (!code || !state) {
            goto('/login?error=oauth_failed') // TODO: aggiungi errore in quella pagina (da param)
            return
        }

        // Fetch google callback
        const res = await fetch('/api/v1/auth/google/callback', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify({ code, state })
        })

        // Check res and navigate
        if (res.ok) {
            goto('/')
        } else {
            goto('/login?error=oauth_failed')
        }
      } catch {
        goto('/login?error=oauth_failed')
      }
    })
</script>

<p>Accesso in corso...</p>