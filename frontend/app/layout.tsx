import type { Metadata } from "next";
import { Roboto } from "next/font/google";
import "./globals.css";



const roboto = Roboto({
  variable: "--font-roboto",
  subsets: ["latin"],
  weight: ["400", "500", "700"],
})


export const metadata: Metadata = {
  title: {
    default: "Daily Social",
    template: "Daily Social - %s",
  },
  icons: {
    icon: "/logo.png",
  }
};
export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={roboto.className}>
      <body
        suppressHydrationWarning
        className={`${roboto.variable}  antialiased`}
      >
        <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          { children }
        </main>
      </body>
    </html>
  );
}
