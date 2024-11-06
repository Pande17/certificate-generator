import React, { useState, useEffect } from "react";
import axios from "axios";

function CompetencyForm() {
  const [competencies, setCompetencies] = useState([]);
  const [selectedCompetency, setSelectedCompetency] = useState(null);
  const [unitDetails, setUnitDetails] = useState({
    hard_skills: [],
    soft_skills: [],
  });
  const [nilai, setNilai] = useState("");

  // Fetch data kompetensi
  useEffect(() => {
    axios
      .get("http://127.0.0.1:3000/api/competence")
      .then((response) => {
        setCompetencies(response.data.data);
      })
      .catch((error) => {
        console.error("Error fetching competencies:", error);
      });
  }, []);

  // Fetch data detail kompetensi berdasarkan pilihan
  const handleCompetencyChange = (kompetensi_id) => {
    setSelectedCompetency(kompetensi_id);
    axios
      .get(`http://127.0.0.1:3000/api/competence/${kompetensi_id}`)
      .then((response) => {
        setUnitDetails(response.data.data);
      })
      .catch((error) => {
        console.error("Error fetching competency details:", error);
      });
  };

  return (
    <div>
      <h2>Form Kompetensi</h2>

      {/* Dropdown untuk memilih mata pelajaran */}
      <label>Pilih Mata Pelajaran:</label>
      <select
        value={selectedCompetency || ""}
        onChange={(e) => handleCompetencyChange(e.target.value)}
      >
        <option value="" disabled>
          Pilih Kompetensi
        </option>
        {competencies.map((competency) => (
          <option key={competency._id} value={competency.kompetensi_id}>
            {competency.nama_kompetensi}
          </option>
        ))}
      </select>

      {/* Bagian unit detail */}
      {selectedCompetency && (
        <div>
          <h3>Detail Unit</h3>

          {/* Menampilkan hard skills */}
          <h4>Hard Skills:</h4>
          {unitDetails.hard_skills.map((skill, index) => (
            <div key={index}>
              <strong>{skill.hardSkill_name}</strong>
              {skill.description.map((desc, idx) => (
                <p key={idx}>
                  Kode Unit: {desc.unit_code}, Judul Unit: {desc.unit_title}
                </p>
              ))}
            </div>
          ))}

          {/* Menampilkan soft skills */}
          <h4>Soft Skills:</h4>
          {unitDetails.soft_skills.map((skill, index) => (
            <div key={index}>
              <strong>{skill.softSkill_name}</strong>
              {skill.description.map((desc, idx) => (
                <p key={idx}>
                  Kode Unit: {desc.unit_code}, Judul Unit: {desc.unit_title}
                </p>
              ))}
            </div>
          ))}
        </div>
      )}

      {/* Input nilai */}
      <div>
        <label>Nilai:</label>
        <input
          type="number"
          value={nilai}
          onChange={(e) => setNilai(e.target.value)}
          placeholder="Masukkan nilai"
        />
      </div>
    </div>
  );
}

export default CompetencyForm;
