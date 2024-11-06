import React, { useEffect, useState } from "react";
import { useForm, Controller } from "react-hook-form";
import { Form, Input, DatePicker, Button, InputNumber, Select } from "antd";
import MainLayout from "../MainLayout/Layout";
import axios from "axios";

function MyForm() {
  const [data, setData] = useState([]);
  const { control, handleSubmit, reset, watch } = useForm();
  console.log(watch());

  const onSubmit = (data) => {
    console.log("Data submitted:", data);
    reset(); // Reset form setelah submit
  };

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
  }, []);

  return (
    <MainLayout>
      <Form
        layout="vertical"
        onFinish={handleSubmit(onSubmit)}
        style={{
          maxHeight: "100vh",
          overflowY: "scroll",
          backgroundColor: "white",
          padding: "40px",
          borderRadius: "20px",
        }}
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
                style={{ width: "400px", height: "50px" }}
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
                style={{ width: "400px", height: "50px" }}
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
                style={{ width: "400px", height: "50px" }}
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
                style={{ width: "400px", height: "50px" }}
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
                style={{ width: "400px", height: "50px" }}
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
                style={{ width: "400px", height: "50px" }}
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
                style={{ width: "400px", height: "50px" }}
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
                style={{ width: "400px", height: "50px" }}
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
                style={{ width: "400px", height: "50px" }}
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
                style={{ width: "400px", height: "50px" }}
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
                style={{ width: "400px", height: "50px" }}
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
                style={{ width: "400px", height: "50px" }}
              />
            )}
          />
        </Form.Item>

        <Form.Item>
          <h1 className="font-Poppins text-xl font-semibold">
            Pilih kompetensi
          </h1>
          <Controller
            name="Kompetensi"
            control={control}
            rules={{ required: "Please Chose One of Competence" }}
            render={({ field }) => (
              <select
                {...field}
                placeholder="Kompetensi"
                style={{
                  width: "400px",
                  height: "50px",
                  border: "2px",
                  borderStyle: "solid",
                  borderRadius: "5px",
                  padding: "10px",
                  borderColor: "gray",
                  opacity: "40%",
                }}
              >
                {data.map((item) => (
                  <option key={item.kompetensi_id} value={item.kompetensi_id}>
                    {item.nama_kompetensi}
                  </option>
                ))}
              </select>
            )}
          />
        </Form.Item>

        <Form.Item></Form.Item>

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
