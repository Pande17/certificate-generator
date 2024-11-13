import React, { useEffect, useState } from "react";
import { useForm, Controller, useFieldArray } from "react-hook-form";
import { Form, Input, DatePicker, Button, InputNumber, Select } from "antd";
import MainLayout from "../MainLayout/Layout";
import axios from "axios";

function MyForm() {
  const [data, setData] = useState([]);
  const [cData, setCData] = useState(null)
  const { control, handleSubmit, reset} = useForm();
  const { field, replace} = useFieldArray({
    control,
    name: "Hardskill",
    name:"SoftSkill"
  });

  const onSubmit = (data) => {
    console.log("Data submitted:", data);
    reset(); t
  };

  const { Option } = Select;
  useEffect(() => {
    const fetchApi = async () => {
      try {
        const response = await axios.get(
          "http://127.0.0.1:3000/api/competence"
        );
        setData(response.data.data);
      } catch (Error) {
        console.log(Error);
      }
    };
    fetchApi();

  }, 
  []);

 const fetchcompetence = async (competenceId) => {
   const type = "id";
   const url = `http://127.0.0.1:3000/api/competence?type=${type}&s=${competenceId}`;

   try {
     const response = await axios.get(url);
     setCData(response.data.data);
     console.log(response.data.data)

   } catch (err) {
     console.log(err);
   }
 };

 const handleCompetence = (value) => {
   fetchcompetence(value);
 };

  return (
    <MainLayout>
      <Form
        layout="vertical"
        style={{
          width:"95%",
          maxHeight: "100vh",
          overflowY: "scroll",
          backgroundColor: "white",
          padding: "40px",
          borderRadius: "20px",
          margin:'auto',
        }}
        onFinish={handleSubmit(onSubmit)}
      >
        <div className="text-center font-Poppins font-bold text-xl">
          Buat Sertifikat
        </div>
        <Form.Item label="Nama" required>
          <Controller
            name="nama"
            control={control}
            rules={{ required: "Nama is required" }}
            render={({ field }) => (
              <Input
                {...field}
                placeholder="Masukkan nama"
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>

        <Form.Item label="Field of Study" required>
          <Controller
            name="fieldOfStudy"
            control={control}
            rules={{ required: "Field of Study is required" }}
            render={({ field }) => (
              <Input
                {...field}
                placeholder="Masukkan field of study"
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>

        <Form.Item label="Valid Time" required>
          <Controller
            name="validTime"
            control={control}
            rules={{ required: "Valid Time is required" }}
            render={({ field }) => (
              <DatePicker
                {...field}
                placeholder="Pilih valid time"
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>

        <Form.Item label="Expired Time (Start)" required>
          <Controller
            name="expiredTimeStart"
            control={control}
            rules={{ required: "Expired Time (Start) is required" }}
            render={({ field }) => (
              <DatePicker
                {...field}
                placeholder="Pilih expired time start"
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>

        <Form.Item label="Expired Time (End)" required>
          <Controller
            name="expiredTimeEnd"
            control={control}
            rules={{ required: "Expired Time (End) is required" }}
            render={({ field }) => (
              <DatePicker
                {...field}
                placeholder="Pilih expired time end"
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>

        <Form.Item label="Code Referral (Order)" required>
          <Controller
            name="codeReferralOrder"
            control={control}
            rules={{ required: "Code Referral (Order) is required" }}
            render={({ field }) => (
              <Input
                {...field}
                placeholder="Masukkan code referral (order)"
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>

        <Form.Item label="Code Referral (Field of Study)" required>
          <Controller
            name="codeReferralFieldOfStudy"
            control={control}
            rules={{ required: "Code Referral (Field of Study) is required" }}
            render={({ field }) => (
              <Input
                {...field}
                placeholder="Masukkan code referral (field of study)"
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>

        <Form.Item label="Code Referral (Month)" required>
          <Controller
            name="codeReferralMonth"
            control={control}
            rules={{ required: "Code Referral (Month) is required" }}
            render={({ field }) => (
              <Input
                {...field}
                placeholder="Masukkan code referral (month)"
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>

        <Form.Item label="Code Referral (Year)" required>
          <Controller
            name="codeReferralYear"
            control={control}
            rules={{ required: "Code Referral (Year) is required" }}
            render={({ field }) => (
              <Input
                {...field}
                placeholder="Masukkan code referral (year)"
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>

        <Form.Item label="SKKNI" required>
          <Controller
            name="skkni"
            control={control}
            rules={{ required: "SKKNI is required" }}
            render={({ field }) => (
              <Input
                {...field}
                placeholder="Masukkan SKKNI"
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>

        <Form.Item label="Total Meeting" required>
          <Controller
            name="totalMeeting"
            control={control}
            rules={{ required: "Total Meeting is required" }}
            render={({ field }) => (
              <InputNumber
                {...field}
                placeholder="Masukkan total meeting"
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>

        <Form.Item label="Meeting Time" required>
          <Controller
            name="meetingTime"
            control={control}
            rules={{ required: "Meeting Time is required" }}
            render={({ field }) => (
              <Input
                {...field}
                placeholder="Masukkan meeting time"
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>

        <h1 className="text-center font-Poppins text-2xl font-medium">
          Pilih kompetensi
        </h1>
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

        <h1 className="text-center font-Poppins text-2xl font-medium">
          Hard skill
        </h1>

        <Form.Item label="Hard skill">
          <Controller
            name="Hardskill"
            control={control}
            render={({ field }) => (
              <Input
                {...field}
                readOnly
                style={{ width: "100%", height: "50px" }}
                placeholder="Tes"
              />
            )}
          />
        </Form.Item>

        <Form.Item label="Code unit & code title ">
          <Controller
            name="Code&JudulUnitHS"
            control={control}
            render={({ field }) => (
              <Input
                {...field}
                readOnly
                style={{ width: "100%", height: "50px" }}
                placeholder="Tes"
              />
            )}
          />
        </Form.Item>

        <Form.Item label="JP" required>
          <Controller
            name="jpHardskill"
            control={control}
            rules={{ required: "Masukan nilai" }}
            render={({ field }) => (
              <Input
                {...field}
                style={{ width: "100%", height: "50px" }}
                placeholder="1-10"
              />
            )}
          />
        </Form.Item>
        <h1 className="text-center font-Poppins text-2xl font-medium">
          Soft skill
        </h1>

        <Form.Item label="Soft skill">
          <Controller
            name="Softskill"
            control={control}
            render={({ field }) => (
              <Input
                {...field}
                readOnly
                style={{ width: "100%", height: "50px" }}
                placeholder="Tes"
              />
            )}
          />
        </Form.Item>

        <Form.Item label="Code unit & code title ">
          <Controller
            name="Code&JudulUnitSS"
            control={control}
            render={({ field }) => (
              <Input
                {...field}
                readOnly
                style={{ width: "100%", height: "50px" }}
                placeholder="Tes"
              />
            )}
          />
        </Form.Item>

        <Form.Item label="JP" required>
          <Controller
            name="jpSoftskill"
            control={control}
            rules={{ required: "Masukan nilai" }}
            render={({ field }) => (
              <Input
                {...field}
                style={{ width: "100%", height: "50px" }}
                placeholder="1-10"
              />
            )}
          />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit" onSubmit={onSubmit}>
            Submit
          </Button>
        </Form.Item>
      </Form>
    </MainLayout>
  );
}

export default MyForm;
