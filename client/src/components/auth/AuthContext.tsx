// authContext.tsx
import React, { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import * as apiHandlers from '../../api/apihandlers';

interface AuthContextType {
    isAuthenticated: boolean;
    user: any | null;
    checkAuth: () => Promise<void>;
}

interface AuthProviderProps {
    children: ReactNode;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
    const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
    const [user, setUser] = useState<any | null>(null);

    const checkAuth = async () => {
        const user = await apiHandlers.fetchUserInfo();
        setUser(user);
        setIsAuthenticated(user != null);
    };

    useEffect(() => {
        checkAuth(); // Initial auth check on component mount
    }, []);

    return (
        <AuthContext.Provider value={{ isAuthenticated, user, checkAuth }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) throw new Error("useAuth must be used within an AuthProvider");
    return context;
};
