'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';

interface AnalysisResult {
  id: string;
  data: any;
  status: 'processing' | 'completed' | 'error';
}

export default function ResultsPage() {
  const params = useParams();
  const [result, setResult] = useState<AnalysisResult | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchResults = async () => {
      try {
        const response = await fetch(
          `${process.env.NEXT_PUBLIC_API_URL}/results/${params.id}`
        );
        
        if (!response.ok) {
          throw new Error('Failed to fetch results');
        }

        const data = await response.json();
        setResult(data);
        
        // If the status is still processing, poll again after a delay
        if (data.status === 'processing') {
          setTimeout(fetchResults, 2000); // Poll every 2 seconds
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An error occurred');
      }
    };

    fetchResults();
  }, [params.id]);

  if (error) {
    return (
      <div className="min-h-screen flex flex-col">
        <header className="bg-gray-800 border-b border-gray-700">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="h-16 flex items-center justify-between">
              <div className="flex items-center">
                <h1 className="text-xl font-semibold text-gray-100">Code Analysis</h1>
              </div>
            </div>
          </div>
        </header>
        <main className="flex-1 bg-gray-900">
          <div className="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
            <div className="bg-red-900/50 border border-red-800 text-red-200 px-4 py-3 rounded-lg">
              <p>{error}</p>
            </div>
          </div>
        </main>
      </div>
    );
  }

  if (!result) {
    return (
      <div className="min-h-screen flex flex-col">
        <header className="bg-gray-800 border-b border-gray-700">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="h-16 flex items-center justify-between">
              <div className="flex items-center">
                <h1 className="text-xl font-semibold text-gray-100">Code Analysis</h1>
              </div>
            </div>
          </div>
        </header>
        <main className="flex-1 bg-gray-900">
          <div className="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
            <div className="animate-pulse">
              <div className="h-8 bg-gray-800 rounded-lg w-1/3 mb-4"></div>
              <div className="space-y-3">
                <div className="h-4 bg-gray-800 rounded-lg"></div>
                <div className="h-4 bg-gray-800 rounded-lg"></div>
                <div className="h-4 bg-gray-800 rounded-lg"></div>
              </div>
            </div>
          </div>
        </main>
      </div>
    );
  }

  if (result.status === 'processing') {
    return (
      <div className="min-h-screen flex flex-col">
        <header className="bg-gray-800 border-b border-gray-700">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="h-16 flex items-center justify-between">
              <div className="flex items-center">
                <h1 className="text-xl font-semibold text-gray-100">Code Analysis</h1>
              </div>
            </div>
          </div>
        </header>
        <main className="flex-1 bg-gray-900">
          <div className="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
            <div className="bg-gray-800 shadow-lg rounded-lg p-6">
              <h2 className="text-lg font-medium leading-6 text-gray-100 mb-6">Analyzing your code...</h2>
              <div className="animate-pulse">
                <div className="flex space-x-4 items-center">
                  <div className="rounded-full bg-gray-700 h-12 w-12"></div>
                  <div className="h-4 bg-gray-700 rounded-lg w-3/4"></div>
                </div>
                <div className="mt-6 space-y-3">
                  <div className="h-4 bg-gray-700 rounded-lg"></div>
                  <div className="h-4 bg-gray-700 rounded-lg"></div>
                  <div className="h-4 bg-gray-700 rounded-lg"></div>
                </div>
              </div>
            </div>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex flex-col">
      <header className="bg-gray-800 border-b border-gray-700">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="h-16 flex items-center justify-between">
            <div className="flex items-center">
              <h1 className="text-xl font-semibold text-gray-100">Code Analysis</h1>
              <p className="ml-4 text-sm text-gray-400">Results ID: {result.id}</p>
            </div>
          </div>
        </div>
      </header>

      <main className="flex-1 bg-gray-900">
        <div className="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
          <div className="bg-gray-800 shadow-lg rounded-lg overflow-hidden">
            <div className="px-4 py-5 sm:p-6">
              <h2 className="text-lg font-medium leading-6 text-gray-100 mb-6">Security Analysis Results</h2>

              {result.status === 'error' ? (
                <div className="bg-red-900/50 border border-red-800 text-red-200 px-4 py-3 rounded-lg">
                  <p className="flex items-center">
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    Error: {typeof result.data === 'object' && result.data.error ? result.data.error : 'An error occurred during analysis'}
                  </p>
                </div>
              ) : (
                <>
                  <div className="overflow-x-auto">
                    <div 
                      className="prose prose-invert max-w-none"
                      dangerouslySetInnerHTML={{ 
                        __html: typeof result.data === 'string' 
                          ? result.data
                            .replace(/```html\n?/g, '')  // Remove opening ```html
                            .replace(/```\s*$/g, '')     // Remove closing ```
                            .replace(/<table/g, '<table class="min-w-full divide-y divide-gray-700"')
                            .replace(/<thead/g, '<thead class="bg-gray-800"')
                            .replace(/<th/g, '<th class="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider"')
                            .replace(/<td/g, '<td class="px-6 py-4 text-sm text-gray-300"')
                            .replace(/<tr/g, '<tr class="hover:bg-gray-700/50"')
                            .replace(/High/g, '<span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-red-900/50 text-red-200">High</span>')
                            .replace(/Medium/g, '<span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-yellow-900/50 text-yellow-200">Medium</span>')
                            .replace(/Low/g, '<span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-900/50 text-green-200">Low</span>')
                          : 'No results found'
                      }} 
                    />
                  </div>
                  
                  <div className="mt-6 flex justify-end">
                    <button 
                      onClick={() => window.print()} 
                      className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 focus:ring-offset-gray-800"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2z" />
                      </svg>
                      Print Results
                    </button>
                  </div>
                </>
              )}
            </div>
          </div>
        </div>
      </main>
    </div>
  );
} 