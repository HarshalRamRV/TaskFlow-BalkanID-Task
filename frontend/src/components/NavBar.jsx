import { useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import { useAuth } from "../auth/AuthContext";
import PropTypes from "prop-types";

const Navbar = ({ user }) => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const { token, logout } = useAuth();
  const navigate = useNavigate();

  const handleToggleMenu = () => {
    setIsMenuOpen(!isMenuOpen);
  };

  const handleDeactivateAccount = async () => {
    try {
      const response = await axios.delete("http://localhost:3001/deactivate", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      if (response.status === 200) {
        logout();
        navigate("/");
      } else {
        console.error("Deactivation failed");
      }
    } catch (error) {
      console.error(error);
    }
  };

  const handleDeleteAccount = async () => {
    try {
      const response = await axios.delete("http://localhost:3001/delete", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      if (response.status === 200) {
        logout();
        navigate("/");
      } else {
        console.error("Deletion failed");
      }
    } catch (error) {
      console.error(error);
    }
  };

  const handleLogout = () => {
    console.log("Logged out");
    navigate("/");
  };

  return (
    <div className="bg-black w-full overflow-hidden">
      <div className={`sm:p-4 p-2 flex justify-between items-center`}>
        <h1 className="text-2xl text-white font-bold">TaskFlow</h1>

        <div className="flex space-x-4 items-center">
          {user ? (
            <>
              <span className="text-white">Hello, {user.name}</span>{" "}
              <div className="cursor-pointer" onClick={handleToggleMenu}>
                <div className="h-5 w-5 text-white">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M4 6h16M4 12h16M4 18h16"
                    />
                  </svg>
                </div>
              </div>
            </>
          ) : (
            ""
          )}

          {isMenuOpen && (
            <div style={sideMenuStyles} className="z-10 hidden">
              <span className="text-black w-full block px-4 py-2">
                Role: {user.role}
              </span>
              <span className="text-black w-full block px-4 py-2">
                Group: {user.group}
              </span>
              <button
                className=" w-full block px-4 py-2 hover:bg-slate-100"
                onClick={handleDeactivateAccount}
              >
                Deactivate Account
              </button>
              <button
                className="w-full block px-4 py-2 hover:bg-slate-100"
                onClick={handleDeleteAccount}
              >
                Delete Account
              </button>
              <button
                className="w-full block px-4 py-2 hover:bg-slate-100"
                onClick={handleLogout}
              >
                Logout
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

const sideMenuStyles = {
  position: "absolute",
  top: "3rem",
  right: "1rem",
  background: "white",
  borderRadius: "8px",
  padding: "8px",
  boxShadow: "0 2px 4px rgba(0, 0, 0, 0.2)",
  display: "flex",
  flexDirection: "column",
  alignItems: "flex-end",
};

Navbar.propTypes = {
  user: PropTypes.object.isRequired,
};

export default Navbar;
