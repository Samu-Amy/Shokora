"use client"

// TODO: cambia (metti tipo caricamento o altro)

import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";

export default function GoogleCallbackPage() {
  const router = useRouter();
  const searchParams = useSearchParams();

  const state = searchParams.get("state");
  const code = searchParams.get("code");

  const [user, setUser] = useState<User>();


  useEffect(() => {
    const handleCallback = async () => {

      if (!state || !code) {
        router.replace("/auth/login?error=oauth_failed"); // TODO: fai pagina per errori auth (?)
        return;
      }

      try {
        const res = await fetch("/api/v1/auth/google/callback", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include", // TODO: usalo (oppure "same-origin" ma potrebbe dare problemi, almeno in dev) in tutte le richieste non pubbliche non legate all'auth
          body: JSON.stringify({ state, code }),
        });


        if (!res.ok) {
          res.json().then(data => console.log(data.error)); // TODO: per controllare i messaggi di errore devo cambiare la logica (ottenere il json comunque prima e poi leggere error se !res.ok e data altrimenti)
          router.replace("/auth/login?error=oauth_failed");
          return;
        }

        const data = await res.json();
        const body = data.data as LoginUserRes;

        if (body.verification_id) {
          // 2FA Required
          router.replace(`auth/2fa?verification_id=${body.verification_id}`);
        }

        if (body.user) {
          setUser(body.user); // [DEBUG]

          // User authenticated
          // router.replace("/");
          return;
        }

        router.replace("/auth/login?error=oauth_failed");

      } catch (err) {
        router.replace("/auth/login?error=oauth_failed");
      }
    };

    handleCallback();
  }, []);

  return (
    <div>
      <h1>Google Callback Page</h1>
      <p>{user?.first_name}</p>
      <p>{user?.last_name}</p>
      <p>{user?.email}</p>
      <p>{user?.birthday.toString()}</p>
      <p>{user?.is_verified}</p>
      <p>{user?.role}</p>
      <p>{user?.permissions}</p>
    </div>
  );
}