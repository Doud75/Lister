export type ValidationResult = {
    success: boolean;
    error?: string;
};

export const validateUsername = (username: string): ValidationResult => {
    const trimmed = username.trim();
    if (trimmed.length < 3) {
        return { success: false, error: "Le nom d'utilisateur doit contenir au moins 3 caractères." };
    }
    if (trimmed.length > 50) {
        return { success: false, error: "Le nom d'utilisateur ne peut pas dépasser 50 caractères." };
    }
    const regex = /^[a-zA-Z0-9_]+$/;
    if (!regex.test(trimmed)) {
        return { success: false, error: "Le nom d'utilisateur ne peut contenir que des lettres, des chiffres et des underscores." };
    }
    return { success: true };
};

export const validatePassword = (password: string): ValidationResult => {
    if (password.length < 8) {
        return { success: false, error: "Le mot de passe doit contenir au moins 8 caractères." };
    }
    if (!/[A-Z]/.test(password)) {
        return { success: false, error: "Le mot de passe doit contenir au moins une majuscule." };
    }
    if (!/[0-9]/.test(password)) {
        return { success: false, error: "Le mot de passe doit contenir au moins un chiffre." };
    }
    if (!/[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]/.test(password)) {
        return { success: false, error: "Le mot de passe doit contenir au moins un caractère spécial." };
    }
    return { success: true };
};

export const sanitizeText = (text: string): string => {
    return text.trim();
};
