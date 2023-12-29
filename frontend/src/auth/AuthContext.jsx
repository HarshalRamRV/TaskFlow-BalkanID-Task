import { createContext, useState, useContext, useEffect } from "react";
import PropTypes from 'prop-types';

export const AuthContext = createContext();

export const useAuth = () => useContext(AuthContext);

export function AuthProvider({ children }) {
  const [token, setToken] = useState(null);

  const login = (tokenData) => {
    setToken(tokenData);
  };

  const logout = () => {
    setToken(null);
  };

  useEffect(() => {
    if (token) {
      console.log("Token updated:", token);
    }
  }, [token]);

  return (
    <AuthContext.Provider value={{ token, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

AuthProvider.propTypes = {
  children: PropTypes.node,
};
