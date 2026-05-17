// --- app/page.tsx ---
"use client"; // This makes it a pure browser component

import { useState } from "react";
import { useRouter } from "next/navigation";
import { pb } from "@/lib/pb";

export default function LoginPage() {
  const router = useRouter();
  const [error, setError] = useState("");

  const handleLogin = async (e: React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault(); // Prevent the default browser form submission

    const formData = new FormData(e.currentTarget);
    const email = formData.get("email") as string;
    const password = formData.get("password") as string;

    try {
      // 1. Authenticate directly from the browser
      // PocketBase will automatically save the JWT to localStorage!
      await pb.collection("users").authWithPassword(email, password);

      // 2. Redirect to dashboard
      router.push("/dashboard");
    } catch (err) {
      setError("Invalid credentials. Please try again.");
    }
  };

  return (
    <div className="flex h-screen items-center justify-center bg-slate-50">
      <form
        onSubmit={handleLogin}
        className="flex flex-col gap-4 border border-slate-200 bg-white p-8 rounded-lg shadow-sm w-96"
      >
        <h1 className="text-2xl font-bold text-slate-900">Sign In</h1>
        <p className="text-sm text-slate-500 mb-4">
          Log into your Amorita Solutions account.
        </p>

        {error && <div className="text-red-500 text-sm mb-2">{error}</div>}

        <input
          name="email"
          type="email"
          placeholder="name@example.com"
          className="border border-slate-300 p-2 rounded focus:outline-none focus:ring-2 focus:ring-slate-900"
          required
        />
        <input
          name="password"
          type="password"
          placeholder="••••••••"
          className="border border-slate-300 p-2 rounded focus:outline-none focus:ring-2 focus:ring-slate-900"
          required
        />

        <button
          type="submit"
          className="bg-slate-900 text-white font-medium p-2 rounded hover:bg-slate-800 transition-colors"
        >
          Sign In
        </button>
      </form>
    </div>
  );
}
