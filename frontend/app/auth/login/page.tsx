"use client"

import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";

export default function LoginPage() {
  const router = useRouter();

  const handleGoogleLogin = async () => {
    try {
      const res = await fetch("/api/v1/auth/google");

      if (!res.ok) {
        // TODO: gestisci (throw error e gestione nel catch?)
        console.log("Errore login google");
      }

      const data = await res.json();

      router.push(data.data.url);

    } catch (err) {
      // TODO: gestisci
      console.log(err);
    }
  }

  return (
    <div className="flex flex-col gap-8 justify-center align-center">
      <h1>Login Page</h1>
      <Button onClick={handleGoogleLogin}>Accedi con Google</Button>
    </div>
  );
}