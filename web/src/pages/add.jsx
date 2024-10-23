import React, { useState } from "react";
import axios from "axios";

const AddCompetence = () => {
  // State untuk menyimpan data form
  const [competenceName, setCompetenceName] = useState("");
  const [hardSkills, setHardSkills] = useState([
    { hardSkill_name: "", description: [{ unit_code: "", unit_title: "" }] },
  ]);
  const [softSkills, setSoftSkills] = useState([
    { softSkill_name: "", description: [{ unit_code: "", unit_title: "" }] },
  ]);

  // Fungsi untuk menangani perubahan input
  const handleCompetenceChange = (e) => {
    setCompetenceName(e.target.value);
  };

  // Fungsi untuk menangani perubahan pada hard skills
  const handleHardSkillsChange = (index, e) => {
    const { name, value } = e.target;
    const newHardSkills = [...hardSkills];
    newHardSkills[index][name] = value;
    setHardSkills(newHardSkills);
  };

  // Fungsi untuk menangani perubahan pada deskripsi hard skill
  const handleHardSkillDescriptionChange = (skillIndex, descIndex, e) => {
    const { name, value } = e.target;
    const newHardSkills = [...hardSkills];
    newHardSkills[skillIndex].description[descIndex][name] = value;
    setHardSkills(newHardSkills);
  };

  // Fungsi untuk menangani perubahan pada soft skills
  const handleSoftSkillsChange = (index, e) => {
    const { name, value } = e.target;
    const newSoftSkills = [...softSkills];
    newSoftSkills[index][name] = value;
    setSoftSkills(newSoftSkills);
  };

  // Fungsi untuk menangani perubahan pada deskripsi soft skill
  const handleSoftSkillDescriptionChange = (skillIndex, descIndex, e) => {
    const { name, value } = e.target;
    const newSoftSkills = [...softSkills];
    newSoftSkills[skillIndex].description[descIndex][name] = value;
    setSoftSkills(newSoftSkills);
  };

  // Fungsi untuk menambahkan hard skill baru
  const addHardSkill = () => {
    setHardSkills([
      ...hardSkills,
      { hardSkill_name: "", description: [{ unit_code: "", unit_title: "" }] },
    ]);
  };

  // Fungsi untuk menambahkan soft skill baru
  const addSoftSkill = () => {
    setSoftSkills([
      ...softSkills,
      { softSkill_name: "", description: [{ unit_code: "", unit_title: "" }] },
    ]);
  };

  // Fungsi untuk submit form
  const handleSubmit = async (e) => {
    e.preventDefault();

    const competenceData = {
      nama_kompetensi: competenceName,
      hard_skills: hardSkills,
      soft_skills: softSkills,
    };

    try {
      const response = await axios.post(
        "http://127.0.0.1:3000/api/competence",
        competenceData
      );
      console.log("Kompetensi berhasil ditambahkan:", response.data);
    } catch (error) {
      console.error("Error saat menambahkan kompetensi:", error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <label>
        Nama Kompetensi:
        <input
          type="text"
          value={competenceName}
          onChange={handleCompetenceChange}
        />
      </label>

      <h3>Hard Skills</h3>
      {hardSkills.map((skill, skillIndex) => (
        <div key={skillIndex}>
          <label>
            Nama Hard Skill:
            <input
              type="text"
              name="hardSkill_name"
              value={skill.hardSkill_name}
              onChange={(e) => handleHardSkillsChange(skillIndex, e)}
            />
          </label>

          {skill.description.map((desc, descIndex) => (
            <div key={descIndex}>
              <label>
                Unit Code:
                <input
                  type="text"
                  name="unit_code"
                  value={desc.unit_code}
                  onChange={(e) =>
                    handleHardSkillDescriptionChange(skillIndex, descIndex, e)
                  }
                />
              </label>
              <label>
                Unit Title:
                <input
                  type="text"
                  name="unit_title"
                  value={desc.unit_title}
                  onChange={(e) =>
                    handleHardSkillDescriptionChange(skillIndex, descIndex, e)
                  }
                />
              </label>
            </div>
          ))}
        </div>
      ))}
      <button type="button" onClick={addHardSkill}>
        Tambah Hard Skill
      </button>

      <h3>Soft Skills</h3>
      {softSkills.map((skill, skillIndex) => (
        <div key={skillIndex}>
          <label>
            Nama Soft Skill:
            <input
              type="text"
              name="softSkill_name"
              value={skill.softSkill_name}
              onChange={(e) => handleSoftSkillsChange(skillIndex, e)}
            />
          </label>

          {skill.description.map((desc, descIndex) => (
            <div key={descIndex}>
              <label>
                Unit Code:
                <input
                  type="text"
                  name="unit_code"
                  value={desc.unit_code}
                  onChange={(e) =>
                    handleSoftSkillDescriptionChange(skillIndex, descIndex, e)
                  }
                />
              </label>
              <label>
                Unit Title:
                <input
                  type="text"
                  name="unit_title"
                  value={desc.unit_title}
                  onChange={(e) =>
                    handleSoftSkillDescriptionChange(skillIndex, descIndex, e)
                  }
                />
              </label>
            </div>
          ))}
        </div>
      ))}
      <button type="button" onClick={addSoftSkill}>
        Tambah Soft Skill
      </button>

      <button type="submit">Submit Kompetensi</button>
    </form>
  );
};

export default AddCompetence;
