"use client"

import { Button } from "@/components/ui/button";
import { useRouter, useSearchParams } from "next/navigation";

export default function LoginPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const error = searchParams.get("error");

  const handleGoogleLogin = async () => {
    try {
      const res = await fetch("/api/v1/auth/google");

      if (!res.ok) {
        // TODO: gestisci (throw error e gestione nel catch?)
        router.replace("/auth/login?error=oauth_failed");
      }

      const data = await res.json();

      router.push(data.data.url);

    } catch (err) {
      router.replace("/auth/login?error=oauth_failed");
    }
  }

  // TODO: guarda tutorial e usa shadcn-ui con TanStack Form , scarica anche zod e TanStack Query

  return (
    <div className="flex flex-col gap-8 justify-center align-center">
      <h1>Login Page</h1>
      {error === "oauth_failed" && <p className="p-4 rounded-md bg-red-50 text-red-600 border border-red-600">
        Errore autenticazione con Google
      </p>}
      <Button onClick={handleGoogleLogin}>Accedi con Google</Button> // TODO: blocca se non si accettano le condizioni
    </div>
  );
}