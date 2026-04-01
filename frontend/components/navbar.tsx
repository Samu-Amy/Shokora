import Link from "next/link";
import { Button } from "./ui/button";
import { LucideHome, LucideMenuSquare, LucideUser } from "lucide-react";

// TODO: fai mobile first

// TODO: imposta valori tailwind (palette colori, spazi, ecc.) e imposta tema chiaro/scuro (?)

export default function Navbar() {
  return (
    <nav className="w-full flex justify-center align-center p-4 order-2 sm:order-1 bg-neutral-900 sm:bg-neutral-50 text-neutral-200 sm:text-neutral-800">

      {/* - Mobile - */}
      <ul className="sm:hidden flex justify-center align-center gap-4">
        <li><Link href="/"><LucideHome /></Link></li>
        <li><Link href="/menu"><LucideMenuSquare /></Link></li>

        {/* 
        TODO: metti in sezione in un menu
        <li><Link href="/about">Chi siamo</Link></li>
        <li><Link href="/contacts">Contatti</Link></li>
        */}
        {/* <li><Link href="/auth/login" className="bg-neutral-100 sm:bg-neutral-900 hover:bg-neutral-200 sm:hover:bg-neutral-800 text-neutral-800 sm:text-neutral-200 font-medium sm:font-normal px-4 py-2 rounded-md">Accedi</Link></li> */}
        <li><Link href="/auth/login">
          <LucideUser />
        </Link></li>
      </ul>

      {/* - Desktop - */}
      <ul className="hidden sm:flex justify-center align-center gap-4">
        <li><Link href="/">Home</Link></li>
        <li><Link href="/menu">Menu</Link></li>
        <li><Link href="/about">Chi siamo</Link></li>
        <li><Link href="/contacts">Contatti</Link></li>
        {/* <li><Link href="/auth/login" className="bg-neutral-100 sm:bg-neutral-900 hover:bg-neutral-200 sm:hover:bg-neutral-800 text-neutral-800 sm:text-neutral-200 font-medium sm:font-normal px-4 py-2 rounded-md">Accedi</Link></li> */}
        <li><Link href="/auth/login">
          <Button size="lg">
            Accedi
          </Button>
        </Link></li>
      </ul>
    </nav>
  );
}