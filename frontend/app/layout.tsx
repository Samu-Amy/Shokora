import type { Metadata } from "next";
import { Geist, Geist_Mono, Inter } from "next/font/google";
import "./globals.css";
import { cn } from "@/lib/utils";
import Navbar from "@/components/navbar";

const inter = Inter({ subsets: ['latin'], variable: '--font-sans' });


// TODO: fai anche robot.txt ed altre cose utili

// TODO: fai manifest per pwa (e capisci come farla funzionare con nginx)

// TODO: controlla corretto utilizzo di html semantics (tags)

// TODO: fai pagina 404 (adatta anche a mobile, con tasto per tornare indietro)

// const geistSans = Geist({
//   variable: "--font-geist-sans",
//   subsets: ["latin"],
// });

// const geistMono = Geist_Mono({
//   variable: "--font-geist-mono",
//   subsets: ["latin"],
// });

const interSans = Geist({
  variable: "--font-inter-sans",
  subsets: ["latin"],
});

const interMono = Geist_Mono({
  variable: "--font-inter-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Shokora",
  description: "Bar pasticceria e cioccolateria", // TODO: cambia
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="en"
      className={cn("h-full", "antialiased", interSans.variable, interMono.variable, "font-sans", inter.variable)}
    >
      {/* TODO: usa tag html bene (tipo nav, aside, ecc.) */}
      <body className="min-h-full flex flex-col justify-between">
        <Navbar />
        <main className="order-1 sm:order-2 flex-1 px-2 sm:px-64">
          {children}
        </main>
      </body>
    </html>
  );
}
