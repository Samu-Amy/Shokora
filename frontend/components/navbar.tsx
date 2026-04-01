import Link from "next/link";

// TODO: fai mobile first

// TODO: imposta valori tailwind (palette colori, spazi, ecc.) e imposta tema chiaro/scuro (?)

export default function Navbar() {
  return (
    <nav className="w-full flex justify-center align-center p-4">
      <ul className="flex justify-center align-center gap-4">
        <li><Link href="/">Home</Link></li>
        <li><Link href="/menu">Menu</Link></li>
        <li><Link href="/about">Chi siamo</Link></li>
        <li><Link href="/contacts">Contatti</Link></li>
        <li><Link href="/auth/login" className="bg-neutral-900 hover:bg-neutral-800 text-neutral-200 px-4 py-2 rounded-md">Accedi</Link></li>
      </ul>
    </nav>
  );
}