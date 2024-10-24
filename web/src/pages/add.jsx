import React, { useState } from "react";
import axios from "axios";
import { Form, Input, Button, Space, message } from "antd";
import { PlusOutlined } from "@ant-design/icons";

const AddCompetence = () => {
  const [competenceName, setCompetenceName] = useState("");
  const [hardSkills, setHardSkills] = useState([
    { hardSkill_name: "", description: [{ unit_code: "", unit_title: "" }] },
  ]);
  const [softSkills, setSoftSkills] = useState([
    { softSkill_name: "", description: [{ unit_code: "", unit_title: "" }] },
  ]);

  const handleCompetenceChange = (e) => setCompetenceName(e.target.value);

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
      { hardSkill_name: "", description: [{ unit_code: "", unit_title: "" }] },
    ]);
  };

  const addSoftSkill = () => {
    setSoftSkills([
      ...softSkills,
      { softSkill_name: "", description: [{ unit_code: "", unit_title: "" }] },
    ]);
  };

  const handleSubmit = async () => {
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
      message.success("Kompetensi berhasil ditambahkan!");
    } catch (error) {
      message.error("Error saat menambahkan kompetensi!");
    }
  };

  return (
    <Form layout="vertical" onFinish={handleSubmit}>
      <Form.Item label="Nama Kompetensi" required>
        <Input
          value={competenceName}
          onChange={handleCompetenceChange}
          placeholder="Masukkan nama kompetensi"
        />
      </Form.Item>

      <h3>Hard Skills</h3>
      {hardSkills.map((skill, skillIndex) => (
        <div key={skillIndex}>
          <Form.Item label={`Nama Hard Skill ${skillIndex + 1}`}>
            <Input
              name="hardSkill_name"
              value={skill.hardSkill_name}
              onChange={(e) => handleHardSkillsChange(skillIndex, e)}
              placeholder="Masukkan nama hard skill"
            />
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
              name="softSkill_name"
              value={skill.softSkill_name}
              onChange={(e) => handleSoftSkillsChange(skillIndex, e)}
              placeholder="Masukkan nama soft skill"
            />
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

      <Button type="primary" style={{width: "200px", height:"40px"}} htmlType="submit" block>
        Submit Kompetensi
      </Button>
    </Form>
  );
};

export default AddCompetence;