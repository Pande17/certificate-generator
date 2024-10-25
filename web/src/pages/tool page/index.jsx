import React, { useEffect, useState } from "react";
import axios from "axios";
import { Form, Input, Button, Space, message, Select } from "antd";
import { PlusOutlined, MinusCircleOutlined } from "@ant-design/icons";

const { Option } = Select;

const AddCompetence = () => {
  const [competenceName, setCompetenceName] = useState("");
  const [hardSkills, setHardSkills] = useState([
    { skill_name: "", description: [{ unit_code: "", unit_title: "" }] },
  ]);
  const [softSkills, setSoftSkills] = useState([
    { skill_name: "", description: [{ unit_code: "", unit_title: "" }] },
  ]);
  const [competencies, setCompetencies] = useState([]);
  const [selectedCompetenceId, setSelectedCompetenceId] = useState(null);

  // Fetch competencies from the API
  useEffect(() => {
    const fetchCompetencies = async () => {
      try {
        const response = await axios.get(
          "http://127.0.0.1:3000/api/competence"
        );
        if (response.data && Array.isArray(response.data.data)) {
          setCompetencies(response.data.data);
        } else {
          message.error("Data kompetensi tidak valid!");
        }
      } catch (error) {
        console.error("Error fetching competencies:", error);
        message.error("Error fetching competencies!");
      }
    };

    fetchCompetencies();
  }, []);

  // Mengisi otomatis hard skills dan soft skills saat kompetensi dipilih
  useEffect(() => {
    const selectedCompetence = competencies.find(
      (c) => c._id === selectedCompetenceId
    );
    if (selectedCompetence) {
      setCompetenceName(selectedCompetence.nama_kompetensi);
      setHardSkills(selectedCompetence.hard_skills || []);
      setSoftSkills(selectedCompetence.soft_skills || []);
    } else {
      resetForm();
    }
  }, [selectedCompetenceId, competencies]);

  const resetForm = () => {
    setCompetenceName("");
    setHardSkills([
      { skill_name: "", description: [{ unit_code: "", unit_title: "" }] },
    ]);
    setSoftSkills([
      { skill_name: "", description: [{ unit_code: "", unit_title: "" }] },
    ]);
    setSelectedCompetenceId(null);
  };

  const handleCompetenceChange = (value) => {
    setSelectedCompetenceId(value);
  };

  const handleHardSkillsChange = (index, e) => {
    const { name, value } = e.target;
    const newHardSkills = [...hardSkills];
    newHardSkills[index][name] = value;
    setHardSkills(newHardSkills);
  };

  const handleHardSkillDescriptionChange = (skillIndex, descIndex, e) => {
    const { name, value } = e.target;
    const newHardSkills = [...hardSkills];
    newHardSkills[skillIndex].description[descIndex][name] = value;
    setHardSkills(newHardSkills);
  };

  const handleSoftSkillsChange = (index, e) => {
    const { name, value } = e.target;
    const newSoftSkills = [...softSkills];
    newSoftSkills[index][name] = value;
    setSoftSkills(newSoftSkills);
  };

  const handleSoftSkillDescriptionChange = (skillIndex, descIndex, e) => {
    const { name, value } = e.target;
    const newSoftSkills = [...softSkills];
    newSoftSkills[skillIndex].description[descIndex][name] = value;
    setSoftSkills(newSoftSkills);
  };

  const addHardSkill = () => {
    setHardSkills([
      ...hardSkills,
      { skill_name: "", description: [{ unit_code: "", unit_title: "" }] },
    ]);
  };

  const addSoftSkill = () => {
    setSoftSkills([
      ...softSkills,
      { skill_name: "", description: [{ unit_code: "", unit_title: "" }] },
    ]);
  };

  // Fungsi untuk menghapus hard skill
  const removeHardSkill = (index) => {
    const newHardSkills = hardSkills.filter((_, i) => i !== index);
    setHardSkills(newHardSkills);
  };

  // Fungsi untuk menghapus soft skill
  const removeSoftSkill = (index) => {
    const newSoftSkills = softSkills.filter((_, i) => i !== index);
    setSoftSkills(newSoftSkills);
  };

  const handleSubmit = async () => {
    const competenceData = {
      nama_kompetensi: competenceName,
      hard_skills: hardSkills,
      soft_skills: softSkills,
    };

    try {
      if (selectedCompetenceId) {
        // Jika mengedit, perbarui kompetensi
        await axios.put(
          `http://127.0.0.1:3000/api/competence/${selectedCompetenceId}`,
          competenceData
        );
        message.success("Kompetensi berhasil diperbarui!");
      } else {
        // Jika menambah, buat kompetensi baru
        await axios.post(
          "http://127.0.0.1:3000/api/competence",
          competenceData
        );
        message.success("Kompetensi berhasil ditambahkan!");
      }
      resetForm(); // Reset form setelah submit
    } catch (error) {
      console.error("Error saat menyimpan kompetensi:", error);
      message.error("Error saat menyimpan kompetensi!");
    }
  };

  return (
    <Form layout="vertical" onFinish={handleSubmit}>
      <Form.Item label="Pilih Kompetensi" required>
        <Select
          placeholder="Pilih kompetensi"
          onChange={handleCompetenceChange}
        >
          <Option value={null}>Tambah Kompetensi Baru</Option>
          {competencies.length > 0 ? (
            competencies.map((competence) => (
              <Option key={competence._id} value={competence._id}>
                {competence.nama_kompetensi}
              </Option>
            ))
          ) : (
            <Option disabled>Tidak ada kompetensi tersedia</Option>
          )}
        </Select>
      </Form.Item>

      <Form.Item label="Nama Kompetensi" required>
        <Input
          value={competenceName}
          onChange={(e) => setCompetenceName(e.target.value)}
          placeholder="Masukkan nama kompetensi"
        />
      </Form.Item>

      <h3>Hard Skills</h3>
      {hardSkills.map((skill, skillIndex) => (
        <div key={skillIndex}>
          <Form.Item label={`Nama Hard Skill ${skillIndex + 1}`}>
            <Input
              name="skill_name"
              value={skill.skill_name}
              onChange={(e) => handleHardSkillsChange(skillIndex, e)}
              placeholder="Masukkan nama hard skill"
            />
            <Button
              type="text"
              danger
              icon={<MinusCircleOutlined />}
              onClick={() => removeHardSkill(skillIndex)}
            >
              Hapus
            </Button>
          </Form.Item>

          {skill.description.map((desc, descIndex) => (
            <Space key={descIndex} direction="vertical">
              <Form.Item label="Unit Code">
                <Input
                  name="unit_code"
                  value={desc.unit_code}
                  onChange={(e) =>
                    handleHardSkillDescriptionChange(skillIndex, descIndex, e)
                  }
                  placeholder="Masukkan unit code"
                />
              </Form.Item>
              <Form.Item label="Unit Title">
                <Input
                  name="unit_title"
                  value={desc.unit_title}
                  onChange={(e) =>
                    handleHardSkillDescriptionChange(skillIndex, descIndex, e)
                  }
                  placeholder="Masukkan unit title"
                />
              </Form.Item>
            </Space>
          ))}
        </div>
      ))}
      <Button
        type="dashed"
        onClick={addHardSkill}
        block
        icon={<PlusOutlined />}
      >
        Tambah Hard Skill
      </Button>

      <h3>Soft Skills</h3>
      {softSkills.map((skill, skillIndex) => (
        <div key={skillIndex}>
          <Form.Item label={`Nama Soft Skill ${skillIndex + 1}`}>
            <Input
              name="skill_name"
              value={skill.skill_name}
              onChange={(e) => handleSoftSkillsChange(skillIndex, e)}
              placeholder="Masukkan nama soft skill"
            />
            <Button
              type="text"
              danger
              icon={<MinusCircleOutlined />}
              onClick={() => removeSoftSkill(skillIndex)}
            >
              Hapus
            </Button>
          </Form.Item>

          {skill.description.map((desc, descIndex) => (
            <Space key={descIndex} direction="vertical">
              <Form.Item label="Unit Code">
                <Input
                  name="unit_code"
                  value={desc.unit_code}
                  onChange={(e) =>
                    handleSoftSkillDescriptionChange(skillIndex, descIndex, e)
                  }
                  placeholder="Masukkan unit code"
                />
              </Form.Item>
              <Form.Item label="Unit Title">
                <Input
                  name="unit_title"
                  value={desc.unit_title}
                  onChange={(e) =>
                    handleSoftSkillDescriptionChange(skillIndex, descIndex, e)
                  }
                  placeholder="Masukkan unit title"
                />
              </Form.Item>
            </Space>
          ))}
        </div>
      ))}
      <Button
        type="dashed"
        onClick={addSoftSkill}
        block
        icon={<PlusOutlined />}
      >
        Tambah Soft Skill
      </Button>

      <Button
        type="primary"
        style={{ width: "200px", height: "40px" }}
        htmlType="submit"
        block
      >
        Submit Kompetensi
      </Button>
    </Form>
  );
};

export default AddCompetence;
