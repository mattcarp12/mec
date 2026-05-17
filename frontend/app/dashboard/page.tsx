// --- app/dashboard/page.tsx ---
"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { pb } from "@/lib/pb";

export default function DashboardPage() {
  const router = useRouter();
  const [userData, setUserData] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // 1. Check if the user is locally authenticated
    if (!pb.authStore.isValid) {
      router.push("/");
      return;
    }

    // 2. Fetch data from your Go backend
    const fetchGoData = async () => {
      try {
        const goRes = await fetch("/api/v1/me", {
          headers: {
            // Attach the JWT directly from the PocketBase store
            Authorization: `Bearer ${pb.authStore.token}`,
          },
        });

        if (!goRes.ok) {
          throw new Error("Token rejected by Go API");
        }

        const data = await goRes.json();
        setUserData(data);
      } catch (err) {
        // If Go rejects the token (e.g., expired), clear the local store and redirect
        pb.authStore.clear();
        router.push("/");
      } finally {
        setLoading(false);
      }
    };

    fetchGoData();
  }, [router]);

  // Show a simple loading state while fetching
  if (loading) {
    return (
      <div className="min-h-screen bg-slate-50 p-8 flex justify-center items-center">
        Loading...
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-50 p-8">
      <div className="max-w-2xl mx-auto bg-white p-6 border border-slate-200 rounded-lg shadow-sm">
        <div className="flex justify-between items-center mb-4">
          <h1 className="text-2xl font-bold text-slate-900">
            Auth Contract Successful
          </h1>
          <button
            onClick={() => {
              pb.authStore.clear();
              router.push("/");
            }}
            className="text-sm bg-slate-200 hover:bg-slate-300 text-slate-800 px-3 py-1 rounded"
          >
            Sign Out
          </button>
        </div>
        <p className="text-slate-600 mb-4">
          The Browser SPA successfully authenticated with PocketBase and passed
          the token to Go.
        </p>

        <div className="bg-slate-900 text-green-400 p-4 rounded font-mono text-sm overflow-x-auto">
          <pre>{JSON.stringify(userData, null, 2)}</pre>
        </div>
      </div>
    </div>
  );
}
