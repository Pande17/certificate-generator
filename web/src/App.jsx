import React from "react";
import { Routes, Route } from "react-router-dom";
import LoginPage from "./pages/login/index.jsx"; 
import Dashboard from "./pages/dashboard/index.jsx";
import CreatePage from "./pages/create certif/index.jsx";
import AddPage from "./pages/Competence/index.jsx";
import Layout from "./pages/MainLayout/Layout.jsx"
import Side from "./pages/Side/index.jsx"
import Tes from "./pages/tes/Index.jsx"
import CreateKompetensi from "./pages/create competence/index.jsx"

const App = () => {
  return (
    <Routes>
      <Route path="/" element={<LoginPage />} />
      <Route path="/dashboard" element={<Dashboard />}/>
      <Route path="/create" element={<CreatePage />}/>
      <Route path="/competence" element={<AddPage />} />
      <Route path="/competence/create-competence" element={<CreateKompetensi />}/>
      <Route path="/layout" element={<Layout />} />
      <Route path="/side" element={<Side />} />
      <Route path="/Tes" element={<Tes />} />
    </Routes>
  );
};

export default App;
