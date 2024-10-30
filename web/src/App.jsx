import React from "react";
import { Routes, Route } from "react-router-dom";
import LoginPage from "./pages/login/index.jsx"; 
import Dashboard from "./pages/dashboard/index.jsx";
import CreatePage from "./pages/create page/index.jsx";
import AddPage from "./pages/tool page/index.jsx";
import Tes from "./pages/Tes.jsx"

const App = () => {
  return (
    <Routes>
      <Route path="/" element={<LoginPage />} />
      <Route path="/dashboard" element={<Dashboard />}/>
      <Route path="/create" element={<CreatePage />}/>
      <Route path="/tool" element={<AddPage />} />
      <Route path="/tes" element={<Tes />} />
    </Routes>
  );
};

export default App;
