import React, { useEffect, useState } from "react";
import { useForm, Controller, useFieldArray } from "react-hook-form";
import {
  Form,
  Input,
  DatePicker,
  Button,
  InputNumber,
  Select,
  message,
} from "antd";
import MainLayout from "../MainLayout/Layout";
import axios from "axios";

function MyForm() {
  const [competenceData, setCompetenceData] = useState([]);
  const [signatureData, setSignatureData] = useState([]); // State baru untuk tanda tangan
  const { control, handleSubmit, reset } = useForm({
    defaultValues: {
      hardSkill: [],
      softSkill: [],
      selectedCompetenceId: "",
      selectedSignatureId: "", // Field baru untuk tanda tangan
    },
  });

  const { Option } = Select;

  useEffect(() => {
    // Fetch competence data
    const fetchCompetence = async () => {
      try {
        const response = await axios.get(
          "http://127.0.0.1:3000/api/competence"
        );
        setCompetenceData(response.data.data);
      } catch (error) {
        console.log("Error fetching competence data:", error);
      }
    };

    // Fetch signature data
    const fetchSignature = async () => {
      try {
        const response = await axios.get("http://127.0.0.1:3000/api/signature");
        setSignatureData(response.data.data);
      } catch (error) {
        console.log("Error fetching signature data:", error);
      }
    };

    fetchCompetence();
    fetchSignature();
  }, []);

  const onSubmit = async (formData) => {
    console.log("Form submitted:", formData);
    // Add your submission logic here
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
        <h1 className="text-center font-Poppins text-2xl font-medium p-6">
          Pilih Kompetensi dan Tanda Tangan
        </h1>

        {/* Dropdown untuk Kompetensi */}
        <Form.Item required>
          <Controller
            name="selectedCompetenceId"
            control={control}
            render={({ field }) => (
              <Select
                placeholder="Pilih kompetensi"
                {...field}
                style={{ width: "100%", height: "50px" }}
              >
                <Option value="" disabled>
                  Pilih Kompetensi
                </Option>
                {competenceData.map((competence) => (
                  <Option key={competence._id} value={competence._id}>
                    {competence.nama_kompetensi || ""}
                  </Option>
                ))}
              </Select>
            )}
          />
        </Form.Item>

        {/* Dropdown untuk Tanda Tangan */}
        <Form.Item required>
          <Controller
            name="selectedSignatureId"
            control={control}
            render={({ field }) => (
              <Select
                placeholder="Pilih tanda tangan"
                {...field}
                style={{ width: "100%", height: "50px" }}
              >
                <Option value="" disabled>
                  Pilih Tanda Tangan
                </Option>
                {signatureData.map((signature) => (
                  <Option key={signature._id} value={signature._id}>
                    {signature.name}
                  </Option>
                ))}
              </Select>
            )}
          />
        </Form.Item>

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
