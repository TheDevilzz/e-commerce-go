import { API_BASE_URL, getErrorMessage } from "./client";

type AuthPayload = {
    token?: string;
    access_token?: string;
    refresh_token?: string;
    expires_in?: number;
    expires_at?: string;
    user?: {
        id: number;
        username: string;
        email: string;
        name: string;
        phone: string;
        date_of_birth: string;
        role: string;
    };
};

type AuthResponse = AuthPayload & {
    message?: string;
    data?: AuthPayload;
};

export const login = async (email: string, username: string, password: string) => {
    try {
        const response = await fetch(`${API_BASE_URL}/login`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ email, username, password }),
        });
        if (!response.ok) {
            throw new Error(await getErrorMessage(response, "Login failed"));
        }
        const data = (await response.json()) as AuthResponse;
        return data.data ?? data;
    } catch (error) {
        console.error("Error during login:", error);
        throw error;
    }
};

export const refreshAuthToken = async (refreshToken: string) => {
    const response = await fetch(`${API_BASE_URL}/refresh`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!response.ok) {
        throw new Error(await getErrorMessage(response, "Refresh token failed"));
    }

    const data = (await response.json()) as AuthResponse;
    return data.data ?? data;
};

export const logout = async (refreshToken: string) => {
    const response = await fetch(`${API_BASE_URL}/logout`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ refresh_token: refreshToken }),
    });

    if (!response.ok) {
        throw new Error(await getErrorMessage(response, "Logout failed"));
    }

    return response.json();
};

export const register = async (
    email: string,
    username: string,
    name: string,
    phone: string,
    dateOfBirth: string,
    password: string,
) => {
    try {
        const response = await fetch(`${API_BASE_URL}/register`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                email,
                username,
                name,
                phone,
                date_of_birth: dateOfBirth,
                password,
            }),
        });
        if (!response.ok) {
            throw new Error(await getErrorMessage(response, "Registration failed"));
        }
        const data = await response.json();
        return data.data ?? data;
    }
    catch (error) {
        console.error("Error during registration:", error);
        throw error;
    }
};
