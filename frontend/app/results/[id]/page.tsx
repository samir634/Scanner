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
      <div className="min-h-screen bg-gradient-to-b from-blue-600 to-blue-900">
        <header className="bg-gradient-to-r from-indigo-700 to-blue-800 text-white shadow-md">
          <div className="max-w-6xl mx-auto py-6 px-8">
            <h1 className="text-3xl font-bold tracking-wide text-gradient">Code Analysis</h1>
          </div>
        </header>
        <main className="max-w-2xl mx-auto p-8">
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg shadow">
            <p>{error}</p>
          </div>
        </main>
      </div>
    );
  }

  if (!result) {
    return (
      <div className="min-h-screen bg-gradient-to-b from-blue-600 to-blue-900">
        <header className="bg-gradient-to-r from-indigo-700 to-blue-800 text-white shadow-md">
          <div className="max-w-6xl mx-auto py-6 px-8">
            <h1 className="text-3xl font-bold tracking-wide text-gradient">Code Analysis</h1>
          </div>
        </header>
        <main className="max-w-2xl mx-auto p-8">
          <div className="animate-pulse">
            <div className="h-8 bg-white bg-opacity-30 rounded-lg w-1/3 mb-4"></div>
            <div className="space-y-3">
              <div className="h-4 bg-white bg-opacity-30 rounded-lg"></div>
              <div className="h-4 bg-white bg-opacity-30 rounded-lg"></div>
              <div className="h-4 bg-white bg-opacity-30 rounded-lg"></div>
            </div>
          </div>
        </main>
      </div>
    );
  }

  if (result.status === 'processing') {
    return (
      <div className="min-h-screen bg-gradient-to-b from-blue-600 to-blue-900">
        <header className="bg-gradient-to-r from-indigo-700 to-blue-800 text-white shadow-md">
          <div className="max-w-6xl mx-auto py-6 px-8">
            <h1 className="text-3xl font-bold tracking-wide text-gradient">Code Analysis</h1>
          </div>
        </header>
        <main className="max-w-2xl mx-auto p-8">
          <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-lg shadow-md p-6 text-white">
            <h2 className="text-2xl font-semibold mb-6">Analyzing your code...</h2>
            <div className="animate-pulse">
              <div className="flex space-x-4 items-center">
                <div className="rounded-full bg-blue-300 h-12 w-12"></div>
                <div className="h-4 bg-blue-300 rounded-lg w-3/4"></div>
              </div>
              <div className="mt-6 space-y-3">
                <div className="h-4 bg-white bg-opacity-20 rounded-lg"></div>
                <div className="h-4 bg-white bg-opacity-20 rounded-lg"></div>
                <div className="h-4 bg-white bg-opacity-20 rounded-lg"></div>
              </div>
            </div>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-b from-blue-600 to-blue-900 flex flex-col">
      <header className="bg-gradient-to-r from-indigo-700 to-blue-800 text-white shadow-md">
        <div className="max-w-6xl mx-auto py-6 px-8">
          <h1 className="text-3xl font-bold tracking-wide text-gradient">Code Analysis</h1>
          <p className="mt-2 opacity-90">Results ID: {result.id}</p>
        </div>
      </header>

      <main className="max-w-6xl mx-auto p-8 flex-grow text-white">
        <div className="mb-8">
          <h2 className="text-2xl font-semibold">Security Analysis Results</h2>
        </div>

        {result.status === 'error' ? (
          <div className="bg-red-900 bg-opacity-70 border border-red-800 text-red-100 px-6 py-4 rounded-lg shadow-md">
            <p className="flex items-center">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Error: {typeof result.data === 'object' && result.data.error ? result.data.error : 'An error occurred during analysis'}
            </p>
          </div>
        ) : (
          <>
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-lg shadow-md overflow-hidden">
              <div className="p-6">
                <style jsx global>{`
                  table {
                    width: 100%;
                    border-collapse: separate;
                    border-spacing: 0;
                    margin-bottom: 1rem;
                    border: 1px solid rgba(255, 255, 255, 0.1);
                    border-radius: 0.375rem;
                    overflow: hidden;
                    font-family: inherit;
                  }
                  thead {
                    background-color: rgba(255, 255, 255, 0.1);
                  }
                  th {
                    padding: 1rem;
                    text-align: left;
                    font-weight: 600;
                    border-bottom: 2px solid rgba(255, 255, 255, 0.1);
                    color: white;
                    text-transform: uppercase;
                    font-size: 0.875rem;
                    letter-spacing: 0.05em;
                    white-space: nowrap;
                  }
                  td {
                    padding: 1rem;
                    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
                    vertical-align: top;
                    color: rgba(255, 255, 255, 0.9);
                    line-height: 1.5;
                  }
                  tr:hover td {
                    background-color: rgba(255, 255, 255, 0.05);
                  }
                  tr:last-child td {
                    border-bottom: none;
                  }
                  .severity-high {
                    color: #ef4444;
                    font-weight: 600;
                    background-color: rgba(239, 68, 68, 0.2);
                    padding: 0.25rem 0.5rem;
                    border-radius: 0.25rem;
                    display: inline-block;
                    margin: 0.125rem 0;
                  }
                  .severity-medium {
                    color: #f59e0b;
                    font-weight: 600;
                    background-color: rgba(245, 158, 11, 0.2);
                    padding: 0.25rem 0.5rem;
                    border-radius: 0.25rem;
                    display: inline-block;
                    margin: 0.125rem 0;
                  }
                  .severity-low {
                    color: #10b981;
                    font-weight: 600;
                    background-color: rgba(16, 185, 129, 0.2);
                    padding: 0.25rem 0.5rem;
                    border-radius: 0.25rem;
                    display: inline-block;
                    margin: 0.125rem 0;
                  }
                  @media print {
                    body { 
                      print-color-adjust: exact; 
                      -webkit-print-color-adjust: exact;
                      background: white !important;
                      color: black !important;
                    }
                    .min-h-screen { 
                      min-height: auto; 
                    }
                    header { 
                      padding: 1rem 0; 
                      margin-bottom: 1rem;
                      background: white !important;
                      color: black !important;
                    }
                    th, td {
                      color: black !important;
                      border-color: #e2e8f0 !important;
                    }
                    thead {
                      background-color: #f1f5f9 !important;
                    }
                    .severity-high,
                    .severity-medium,
                    .severity-low {
                      background-color: transparent !important;
                      color: black !important;
                    }
                    @page { 
                      margin: 1.5cm; 
                    }
                  }
                `}</style>
                <div 
                  className="analysis-results"
                  dangerouslySetInnerHTML={{ 
                    __html: typeof result.data === 'string' 
                      ? result.data
                        .replace(/\bHigh\b/g, '<span class="severity-high">High</span>')
                        .replace(/\bMedium\b/g, '<span class="severity-medium">Medium</span>')
                        .replace(/\bLow\b/g, '<span class="severity-low">Low</span>')
                      : 'No results found'
                  }} 
                />
              </div>
            </div>
            
            <div className="mt-6 flex justify-end">
              <button 
                onClick={() => window.print()} 
                className="bg-indigo-500 hover:bg-indigo-600 text-white font-medium py-2 px-6 rounded-md shadow-sm transition-colors flex items-center"
              >
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2z" />
                </svg>
                Print Results
              </button>
            </div>
          </>
        )}
      </main>
      
      <footer className="border-t border-blue-500 bg-blue-800 bg-opacity-30 py-4 text-sm text-blue-100 w-full">
        <div className="flex justify-center items-center">
          <p>© {new Date().getFullYear()} Deviate, LLC. All rights reserved.</p>
        </div>
      </footer>
    </div>
  );
} 