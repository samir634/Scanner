'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';

export default function Home() {
  const [file, setFile] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setFile(e.target.files[0]);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!file) return;

    setLoading(true);
    const formData = new FormData();
    formData.append('file', file);

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/upload`, {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        throw new Error('Upload failed');
      }

      const data = await response.json();
      router.push(`/results/${data.id}`);
    } catch (error) {
      console.error('Error:', error);
      alert('Failed to upload file. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className="min-h-screen p-8">
      <div className="max-w-2xl mx-auto">
        <h1 className="text-3xl font-bold mb-8 tracking-wide text-gradient">Code Analysis</h1>
        <div className="bg-white/10 backdrop-blur-sm p-6 rounded-lg shadow-md border border-white/20">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-white mb-2">
                Upload Code File
              </label>
              <input
                type="file"
                onChange={handleFileChange}
                accept=".zip,.js,.jsx,.ts,.tsx,.py,.java,.cpp,.c,.cs,.go,.rb,.php"
                className="w-full p-2 bg-white/5 border border-white/20 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                required
              />
              <p className="mt-1 text-sm text-gray-300">
                Supported formats: ZIP, JavaScript, TypeScript, Python, Java, C++, C, C#, Go, Ruby, PHP
              </p>
            </div>
            <button
              type="submit"
              disabled={!file || loading}
              className={`w-full py-2 px-4 rounded-md text-white font-medium transition-colors ${
                !file || loading
                  ? 'bg-gray-600 cursor-not-allowed'
                  : 'bg-blue-600 hover:bg-blue-700'
              }`}
            >
              {loading ? 'Uploading...' : 'Analyze Code'}
            </button>
          </form>
        </div>
      </div>
    </main>
  );
} 