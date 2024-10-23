import React from "react";
import { Routes, Route } from "react-router-dom";
import LoginPage from "./pages/login.jsx"; 
import Dashboard from "./pages/dashboard.jsx";
import CreatePage from "./pages/Create.jsx";
import AddPage from "./pages/add.jsx";

const App = () => {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/dashboard" element={<Dashboard />}/>
      <Route path="/create" element={<CreatePage />}/>
      <Route path="/add" element={<AddPage />} />
    </Routes>
  );
};

export default App;
