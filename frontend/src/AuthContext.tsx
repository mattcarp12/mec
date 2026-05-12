import React, {
    createContext,
    useContext,
    useEffect,
    useState,
} from 'react';

import { api } from './api';

interface User {
    sub: string;
}

interface AuthState {
    isAuthenticated: boolean;
    user: User | null;
    loading: boolean;

    refreshSession: () => Promise<void>;
    logout: () => Promise<void>;
}

const AuthContext = createContext<AuthState>({
    isAuthenticated: false,
    user: null,
    loading: true,

    refreshSession: async () => { },
    logout: async () => { },
});

export function AuthProvider({
    children,
}: {
    children: React.ReactNode;
}) {

    const [loading, setLoading] = useState(true);

    const [user, setUser] = useState<User | null>(null);

    async function refreshSession() {
        try {
            const res = await api.get('/api/v1/users/me');

            setUser(res.data);
        } catch {
            setUser(null);
        }
    }

    async function logout() {
        try {
            await api.post('/api/v1/auth/logout');
        } finally {
            setUser(null);
        }
    }

    useEffect(() => {
        refreshSession()
            .finally(() => setLoading(false));
    }, []);

    return (
        <AuthContext.Provider
            value={{
                isAuthenticated: !!user,
                user,
                loading,
                refreshSession,
                logout,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    return useContext(AuthContext);
}