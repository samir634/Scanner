import type { Metadata } from 'next';
import React, { type ReactNode } from 'react';
import { Montserrat } from 'next/font/google';
import './globals.css';

const montserrat = Montserrat({
  subsets: ['latin'],
  display: 'swap',
  weight: ['400', '600', '700'],
});

export const metadata: Metadata = {
  title: 'Code Analysis',
  description: 'Analyze your code for security vulnerabilities',
}

export default function RootLayout({
  children,
}: {
  children: ReactNode
}) {
  return (
    <html lang="en" className={montserrat.className}>
      <body className="min-h-screen bg-gradient-to-b from-blue-600 to-blue-900 text-white">
        {children}
      </body>
    </html>
  )
} 