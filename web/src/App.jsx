import React from "react";
import { Routes, Route } from "react-router-dom";
import LoginPage from "./pages/login/index.jsx"; 
import Dashboard from "./pages/dashboard/index.jsx";
import CreatePage from "./pages/create page/index.jsx";
import AddPage from "./pages/tool page/index.jsx";
import Layout from "./pages/MainLayout/Layout.jsx"
import Side from "./pages/Side/index.jsx"
import Tes from "./pages/tes/Index.jsx"

const App = () => {
  return (
    <Routes>
      <Route path="/" element={<LoginPage />} />
      <Route path="/dashboard" element={<Dashboard />}/>
      <Route path="/create" element={<CreatePage />}/>
      <Route path="/tool" element={<AddPage />} />
      <Route path="/layout" element={<Layout />} />
      <Route path="/side" element={<Side />} />
      <Route path="/Tes" element={<Tes />} />
    </Routes>
  );
};

export default App;
