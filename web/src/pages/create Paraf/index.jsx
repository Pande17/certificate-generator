import { useEffect, useState } from "react";
import { useForm, Controller } from "react-hook-form";
import { Form, Input, Button, message } from "antd";
import { BackwardFilled } from "@ant-design/icons";
import MainLayout from "../MainLayout/Layout";
import { useNavigate } from "react-router-dom";
import { Signature } from "../api middleware";

const CreateParaf = () => {
  const [data, setData] = useState([]);
  const { control, handleSubmit, reset } = useForm();
  const navigate = useNavigate();

  useEffect(() => {
    const fetchSignature = async () => {
      try {
        const response = await Signature.get("/");
        setData(response.data.data);
      } catch (err) {
        message.error("Error fetching signature data");
        console.error(err);
      }
    };
    fetchSignature();
  }, []);

  // Function to navigate back to the previous page
  const backHandle = () => {
    navigate("/signature");
  };

  // Form submission handler
  const onSubmit = async (formData) => {
    try {
      // Format the request data
      const requestData = {
        ...formData,
        signature: formData.ttd, // The link to the signature image
        stamp: formData.Cap, // The link to the stamp image
        name: formData.atasNama,
        config_name: formData.atasNama,
        role: formData.jabatan,
      };

      const response = await Signature.post("/", requestData);
      if (response.status === 200) {
        message.success("Data saved successfully!");
        reset(); // Reset form after successful submission
      } else {
        message.error("Failed to save data. Please try again.");
      }
    } catch (error) {
      message.error("An error occurred while submitting the form.");
      console.error("Error submitting form:", error);
    }
  };

  return (
    <MainLayout>
      <div className="m-2">
        <Button
          style={{ width: "50px", height: "50px" }}
          icon={<BackwardFilled />}
          onClick={backHandle}
        />
      </div>
      <Form
        layout="vertical"
        onFinish={handleSubmit(onSubmit)}
        style={{
          width: "95%",
          maxHeight: "100vh",
          overflowY: "scroll",
          backgroundColor: "white",
          padding: "40px",
          borderRadius: "20px",
        }}
      >
        <h3 className="text-center font-Poppins text-2xl font-bold p-6">
          Buat Paraf
        </h3>

        {/* Sertifikat Name */}
        <Form.Item label="Display Nama" required>
          <Controller
            name="displayNama"
            control={control}
            rules={{ required: "Wajib mengisi display nama" }}
            render={({ field }) => (
              <Input {...field} placeholder="Masukkan nama display" />
            )}
          />
        </Form.Item>

        <Form.Item label="Nama Penandatangan" required>
          <Controller
            name="atasNama"
            control={control}
            rules={{ required: "Wajib mengisi Nama" }}
            render={({ field }) => (
              <Input {...field} placeholder="Masukkan nama Penandatangan" />
            )}
          />
        </Form.Item>

        <Form.Item label="Jabatan Penandatangan" required>
          <Controller
            name="jabatan"
            control={control}
            rules={{ required: "Wajib mengisi jabatan" }}
            render={({ field }) => (
              <Input {...field} placeholder="Masukkan Jabatan Penandatangan" />
            )}
          />
        </Form.Item>

        <Form.Item label="Link Gambar Tanda Tangan" required>
          <Controller
            name="ttd"
            control={control}
            rules={{ required: "Wajib mengisi Link" }}
            render={({ field }) => (
              <Input
                {...field}
                placeholder="Masukkan Link Gambar Tanda Tangan"
              />
            )}
          />
        </Form.Item>

        <Form.Item label="Link Gambar Cap Perusahaan" required>
          <Controller
            name="Cap"
            control={control}
            rules={{ required: "Wajib mengisi Link" }}
            render={({ field }) => (
              <Input {...field} placeholder="Masukkan Link Gambar Cap" />
            )}
          />
        </Form.Item>

        <Form.Item label="Link logo Perusahaan" required>
          <Controller
            name="linkLogo"
            control={control}
            defaultValue="https://res.cloudinary.com/dektxbmmb/image/upload/v1727833019/aset%20pdf/pnu45hydtyftsfxlqaxm.png"
            rules={{ required: "Link logo perusahaan diperlukan" }}
            render={({ field }) => (
              <>
                <Input
                  {...field}
                  placeholder="Masukkan Link Gambar"
                  style={{ width: "100%", height: "50px" }}
                />
                {/* Menampilkan gambar dari link yang dimasukkan */}
                {field.value && (
                  <div style={{ marginTop: "10px" }}>
                    <img
                      src={field.value}
                      alt="Logo Perusahaan"
                      style={{
                        width: "200px",
                        height: "200px",
                        border: "solid",
                        borderColor: "black",
                      }}
                    />
                  </div>
                )}
              </>
            )}
          />
        </Form.Item>

        <Form.Item>
          <Button
            type="primary"
            htmlType="submit"
            style={{ width: "100%", height: "50px" }}
          >
            Simpan
          </Button>
        </Form.Item>
      </Form>
    </MainLayout>
  );
};

export default CreateParaf;
