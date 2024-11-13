import React, { useEffect, useState } from "react";
import { useForm, Controller, useFieldArray } from "react-hook-form";
import { Form, Input, Button, Select } from "antd";
import MainLayout from "../MainLayout/Layout";
import axios from "axios";

function MyForm() {
  const [data, setData] = useState([]); // Competence data
  const [competenceData, setCompetenceData] = useState(null); // Selected competence data
  const { control, handleSubmit, reset } = useForm({
    defaultValues: {
      hardSkill: [],
      softSkill: [],
      selectedCompetenceId: "",
    },
  });

  const { fields: hardSkillFields, append: appendHardSkill } = useFieldArray({
    control,
    name: "hardSkill",
  });

  const { fields: softSkillFields, append: appendSoftSkill } = useFieldArray({
    control,
    name: "softSkill",
  });

  const onSubmit = (data) => {
    console.log("Data submitted:", data);
    reset(); // Reset after form submission
  };

  const { Option } = Select;

  useEffect(() => {
    const fetchApi = async () => {
      try {
        const response = await axios.get(
          "http://127.0.0.1:3000/api/competence"
        );
        setData(response.data.data); // Store competence data
      } catch (error) {
        console.log(error);
      }
    };
    fetchApi();
  }, []);

  const fetchCompetence = async (competenceId) => {
    const type = "id";
    const url = `http://127.0.0.1:3000/api/competence?type=${type}&s=${competenceId}`;
    try {
      const response = await axios.get(url);
      setCompetenceData(response.data.data);

      if (response.data.data.hard_skills) {
        response.data.data.hard_skills.forEach((hardSkill) =>
          appendHardSkill({
            skill_name: hardSkill.hardskill_name || "",
            combined_units: hardSkill.description
              .map((unit) => `${unit.unit_code} - ${unit.unit_title}`)
              .join("\n"), // Combines unit code and title in a single TextArea format
          })
        );
      }

      if (response.data.data.soft_skills) {
        response.data.data.soft_skills.forEach((softSkill) =>
          appendSoftSkill({
            skill_name: softSkill.softskill_name || "",
            combined_units: softSkill.description
              .map((unit) => `${unit.unit_code} - ${unit.unit_title}`)
              .join("\n"), // Combines unit code and title in a single TextArea format
          })
        );
      }
    } catch (err) {
      console.log(err);
    }
  };

  const handleCompetence = (value) => {
    fetchCompetence(value);
  };

  return (
    <MainLayout>
      <Form
        layout="vertical"
        style={{
          width: "95%",
          maxHeight: "100vh",
          overflowY: "scroll",
          backgroundColor: "white",
          padding: "40px",
          borderRadius: "20px",
          margin: "auto",
        }}
        onFinish={handleSubmit(onSubmit)}
      >
        <div className="text-center font-Poppins font-bold text-xl">
          Buat Sertifikat
        </div>

        <Form.Item label="Pilih Kompetensi" required>
          <Controller
            name="selectedCompetenceId"
            control={control}
            render={({ field }) => (
              <Select
                placeholder="Pilih kompetensi"
                {...field}
                style={{ width: "100%", height: "50px" }}
                onChange={(value) => {
                  field.onChange(value);
                  handleCompetence(value);
                }}
              >
                <Option value="" disabled>
                  Tambah Kompetensi Baru
                </Option>
                {data.length > 0 ? (
                  data.map((competence) => (
                    <Option key={competence._id} value={competence._id}>
                      {competence.nama_kompetensi || ""}
                    </Option>
                  ))
                ) : (
                  <Option disabled>Tidak ada kompetensi tersedia</Option>
                )}
              </Select>
            )}
          />
        </Form.Item>

        {/* Hard Skills */}
        {hardSkillFields.map((skill, index) => (
          <div key={index}>
            <Form.Item label={`Hardskill ${index + 1} Name`}>
              <Controller
                name={`hardSkill[${index}].skill_name`}
                control={control}
                render={({ field }) => (
                  <Input {...field} placeholder="Skill Name" />
                )}
              />
            </Form.Item>
            <Form.Item label="Units">
              <Controller
                name={`hardSkill[${index}].combined_units`}
                control={control}
                render={({ field }) => (
                  <Input.TextArea
                    {...field}
                    placeholder="Unit Code and Title"
                    rows={4}
                  />
                )}
              />
            </Form.Item>
          </div>
        ))}

        {/* Soft Skills */}
        {softSkillFields.map((skill, index) => (
          <div key={index}>
            <Form.Item label={`Softskill ${index + 1} Name`}>
              <Controller
                name={`softSkill[${index}].skill_name`}
                control={control}
                render={({ field }) => (
                  <Input {...field} placeholder="Skill Name" />
                )}
              />
            </Form.Item>
            <Form.Item label="Units">
              <Controller
                name={`softSkill[${index}].combined_units`}
                control={control}
                render={({ field }) => (
                  <Input.TextArea
                    {...field}
                    placeholder="Unit Code and Title"
                    rows={4}
                  />
                )}
              />
            </Form.Item>
          </div>
        ))}

        <Form.Item>
          <Button type="primary" htmlType="submit">
            Submit
          </Button>
        </Form.Item>
      </Form>
    </MainLayout>
  );
}

export default MyForm;
