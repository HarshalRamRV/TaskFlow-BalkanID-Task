import "./App.css";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import RegistrationForm from "./components/Register";
import LoginForm from "./components/Login";
import TaskBoard from "./components/TaskBoard";
import { AuthProvider } from "./auth/AuthContext"; 

const App = () => {
  return (
    <AuthProvider> 
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<LoginForm />} />
          <Route path="/register" element={<RegistrationForm />} />
          <Route path="/tasks" element={<TaskBoard />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
};

export default App;
