import React from "react";
import { Routes, Route } from "react-router-dom";
import LoginPage from "./pages/login.jsx"; 
import Dashboard from "./pages/dashboard.jsx";

const App = () => {
  return (
    <Routes>
      <Route path="/" element={<LoginPage />} />
      <Route path="/dashboard" element={<Dashboard />}/>
    </Routes>
  );
};

export default App;
