import React, { createContext, useContext } from "react";
import { Redirect } from "react-router-dom";

export const AuthContext = createContext();

const AuthProvider = ({ children }) => {
    const[token, setToken] = useState({});

    return {};
}

