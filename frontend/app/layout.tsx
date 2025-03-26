import type { Metadata } from 'next';
import React, { type ReactNode } from 'react';
import { Inter } from 'next/font/google';
import './globals.css';

const inter = Inter({
  subsets: ['latin'],
  display: 'swap',
  weight: ['400', '500', '600', '700'],
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
    <html lang="en" className={inter.className}>
      <body className="min-h-screen bg-gray-900 text-gray-100">
        {children}
      </body>
    </html>
  )
} 