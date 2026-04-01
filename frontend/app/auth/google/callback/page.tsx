"use client"

// TODO: cambia (metti tipo caricamento o altro)

import { useRouter, useSearchParams } from "next/navigation";
import { useEffect } from "react";

export default function GoogleCallbackPage() {
  const router = useRouter();
  const searchParams = useSearchParams();

  const state = searchParams.get("state");
  const code = searchParams.get("code");


  useEffect(() => {
    if (!state || !code) {
      router.replace("/auth/login?error=oauth_failed"); // TODO: fai pagina per errori auth (?)
      return;
    }

    try {
      fetch("/api/v1/auth/google/callback", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include", // TODO: usalo (oppure "same-origin" ma potrebbe dare problemi, almeno in dev) in tutte le richieste non pubbliche non legate all'auth
        body: JSON.stringify({ state, code }),
      }).catch(_ => {
        router.replace("/auth/login?error=oauth_failed");
      });

    } catch (err) {
      router.replace("/auth/login?error=oauth_failed");
    }

  }, []);

  return (
    <div>
      <h1>Google Callback Page</h1>
    </div>
  );
}