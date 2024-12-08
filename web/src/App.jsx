import React from "react";
import { Routes, Route } from "react-router-dom";
import LoginPage from "./pages/login/index.jsx"; 
import Dashboard from "./pages/dashboard/index.jsx";
import CreatePage from "./pages/create certif/index.jsx";
import AddPage from "./pages/Competence/index.jsx";
import Layout from "./pages/MainLayout/Layout.jsx"
import Side from "./pages/Side/index.jsx"
import CreateKompetensi from "./pages/create competence/index.jsx"
import SignaturePage from "./pages/Signature page/index.jsx";
import CreateParaf from "./pages/create Paraf/index.jsx";
import CertificateTable from "./pages/qr page/index.jsx";
// import errorHandle from "./pages/Error.jsx";



const App = () => {
  return (
    <Routes>
      <Route path="/" element={<LoginPage />} />
      <Route path="/dashboard" element={<Dashboard />} />
      <Route path="/create" element={<CreatePage />} />
      <Route path="/competence" element={<AddPage />} />
      <Route
        path="/competence/create-competence"
        element={<CreateKompetensi />}
      />
      {/* <Route path="*" component={errorHandle} /> */}
      <Route path="/layout" element={<Layout />} />
      <Route path="/side" element={<Side />} />
      <Route path="/qrPage/:id" element={<CertificateTable />} />
      <Route path="/signature" element={<SignaturePage />} />
      <Route path="/createParaf" element={<CreateParaf />} />
    </Routes>
  );
};

export default App;
