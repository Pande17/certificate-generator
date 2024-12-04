import { useEffect, useState } from "react";
import { useForm, Controller, useFieldArray } from "react-hook-form";
import { Form, Input, Button, Space, message, Select } from "antd";
import {
  PlusOutlined,
  MinusCircleOutlined,
  BackwardFilled,
} from "@ant-design/icons";
import MainLayout from "../MainLayout/Layout";
import { useNavigate } from "react-router-dom";
import { Kompetensi } from "../api middleware";

const Tool = () => {
  const navigate = useNavigate();

  const backHandle = () => {
    navigate("/competence");
  };

  const { control, handleSubmit, reset, watch, getValues } = useForm({
    defaultValues: {
      skkni: "",
      divisi: "",
      competenceName: "",
      hardSkills: [
        {
          skill_name: "",
          description: [{ id: "", unit_code: "", unit_title: "" }],
        },
      ],
      softSkills: [
        { skill_name: "", description: [{ unit_code: "", unit_title: "" }] },
      ],
      selectedCompetenceId: null,
    },
  });

  const {
    fields: hardSkillsFields,
    append: addHardSkill,
    update: upHardSkill,
    remove: removeHardSkill,
  } = useFieldArray({ control, name: "hardSkills" });

  const {
    fields: softSkillsFields,
    append: addSoftSkill,
    update: upSoftSkill,
    remove: removeSoftSkill,
  } = useFieldArray({ control, name: "softSkills" });

  const [competencies, setCompetencies] = useState([]);
  // Fetch competencies from the API
  useEffect(() => {
    const fetchCompetencies = async () => {
      try {
        const response = await Kompetensi.get("/");
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

  console.log(hardSkillsFields, watch());
  const onSubmit = async (data) => {
    const competenceData = {
      nama_kompetensi: data.competenceName,
      skkni: data.skkni,
      divisi: data.devisi,
      hard_skills: data.hardSkills,
      soft_skills: data.softSkills,
    };

    try {
      if (data.selectedCompetenceId) {
        await Kompetensi.put(`/${data.selectedCompetenceId}`, competenceData);
        message.success("Kompetensi berhasil diperbarui!");
      } else {
        await Kompetensi.post("/", competenceData);
        message.success("Kompetensi berhasil ditambahkan!");
      }
      reset();
    } catch (error) {
      console.error("Error saat menyimpan kompetensi:", error);
      message.error("Error saat menyimpan kompetensi!");
    }
    navigate("/competence");
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
          Buat kompetensi{" "}
        </h3>

        <Form.Item label="Nama Kompetensi" required>
          <Controller
            name="competenceName"
            control={control}
            rules={{ required: "Nama Kompetensi wajib diisi!" }}
            render={({ field }) => (
              <Input
                placeholder="Masukkan nama kompetensi"
                {...field}
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>
        <Form.Item label="Skkni" required>
          <Controller
            name="skkni"
            control={control}
            render={({ field }) => (
              <Input
                placeholder="SKKNI No. 16 Th. 2016"
                {...field}
                style={{ width: "100%", height: "50px" }}
              />
            )}
          />
        </Form.Item>

        <Form.Item label="Divisi" required>
          <Controller
            name="devisi"
            control={control}
            rules={{
              required:
                "Input divisi harus diisi!",
              validate: (value) =>
                value.length <= 6 ||
                "Input divisi terlalu panjang! (maks 6 huruf)",
            }}
            render={({ field, fieldState: { error } }) => (
              <>
                <Input
                  placeholder="IT"
                  {...field}
                  style={{ width: "100%", height: "50px" }}
                />
                {error && (
                  <span style={{ color: "red", fontSize: "12px" }}>
                    {error.message}
                  </span>
                )}
              </>
            )}
          />
        </Form.Item>

        <h3 className="text-center font-Poppins text-2xl font-medium p-6">
          Hard Skills
        </h3>
        {hardSkillsFields.map((field, index) => {
          return (
            <div key={field.id}>
              <Form.Item label={`Nama Hard Skill ${index + 1}`}>
                <Controller
                  name={`hardSkills.${index}.skill_name`}
                  control={control}
                  render={({ field }) => (
                    <Input
                      placeholder="Masukkan nama hard skill"
                      {...field}
                      style={{ width: "100%", height: "50px" }}
                    />
                  )}
                />
                <Button
                  type="text"
                  danger
                  icon={<MinusCircleOutlined />}
                  onClick={() => removeHardSkill(index)}
                >
                  Hapus
                </Button>
              </Form.Item>

              {/* Tambahkan tombol dan field untuk deskripsi */}
              <Space direction="vertical">
                {field.description.map((descField, descIndex) => (
                  <div key={descField.id}>
                    <Form.Item label={`kode unit ${descIndex + 1}`}>
                      <Controller
                        name={`hardSkills.${index}.description.${descIndex}.unit_code`}
                        control={control}
                        render={({ field }) => (
                          <Input
                            placeholder="Masukkan unit code"
                            {...field}
                            style={{ width: "100%", height: "50px" }}
                          />
                        )}
                      />
                    </Form.Item>
                    <Form.Item label={`judul unit ${descIndex + 1}`}>
                      <Controller
                        name={`hardSkills.${index}.description.${descIndex}.unit_title`}
                        control={control}
                        render={({ field }) => (
                          <Input
                            placeholder="Masukkan unit title"
                            {...field}
                            style={{ width: "100%", height: "50px" }}
                          />
                        )}
                      />
                    </Form.Item>
                    <Button
                      type="text"
                      danger
                      icon={<MinusCircleOutlined />}
                      onClick={() => {
                        // Salin field yang sedang diperbarui tanpa merubah field lainnya
                        const updatedField = {
                          ...hardSkillsFields[index], // Salin seluruh field
                          description: hardSkillsFields[
                            index
                          ].description.filter(
                            (_, i) => i !== descIndex // Hapus deskripsi pada index tertentu
                          ),
                        };

                        // Update array hardSkillsFields tanpa merubah elemen lain
                        const updatedFields = hardSkillsFields.map(
                          (field, idx) => (idx === index ? updatedField : field)
                        );

                        // Panggil upHardSkill untuk memperbarui state
                        upHardSkill(index, updatedField);
                      }}
                    >
                      Hapus Deskripsi
                    </Button>
                  </div>
                ))}

                <Button
                  id={`hardSkills.${index}.description`}
                  type="dashed"
                  htmlType="button"
                  onClick={() => {
                    console.log(hardSkillsFields);
                    upHardSkill(index, {
                      ...getValues("hardSkills")[index],
                      description: [
                        ...getValues("hardSkills")[index].description,
                        { id: "", unit_code: "", unit_title: "" },
                      ],
                    });
                  }}
                  icon={<PlusOutlined />}
                  style={{ marginBottom: "20px" }}
                >
                  Tambah Unit Code dan Title
                </Button>
              </Space>
            </div>
          );
        })}

        <Button
          type="dashed"
          onClick={() =>
            addHardSkill({
              skill_name: "",
              description: [{ unit_code: "", unit_title: "" }],
            })
          }
          block
          icon={<PlusOutlined />}
          style={{ marginBottom: "20px" }}
        >
          Tambah Hard Skill
        </Button>

        <h3 className="text-center font-Poppins text-2xl font-medium p-6">
          Soft Skills
        </h3>

        {softSkillsFields.map((field, index) => {
          return (
            <div key={field.id}>
              <Form.Item label={`Nama Soft Skill ${index + 1}`}>
                <Controller
                  name={`softSkills.${index}.skill_name`}
                  control={control}
                  render={({ field }) => (
                    <Input
                      placeholder="Masukkan nama soft skill"
                      {...field}
                      style={{ width: "100%", height: "50px" }}
                    />
                  )}
                />
                <Button
                  type="text"
                  danger
                  icon={<MinusCircleOutlined />}
                  onClick={() => removeSoftSkill(index)}
                >
                  Hapus
                </Button>
              </Form.Item>

              {/* Tambahkan tombol dan field untuk deskripsi */}
              <Space direction="vertical">
                {field.description.map((descfield, descIndex) => (
                  <div key={descfield.id}>
                    <Form.Item label={`kode unit ${descIndex + 1}`}>
                      <Controller
                        name={`softSkills.${index}.description.${descIndex}.unit_code`}
                        control={control}
                        render={({ field }) => (
                          <Input
                            placeholder="Masukkan unit code"
                            {...field}
                            style={{ width: "100%", height: "50px" }}
                          />
                        )}
                      />
                    </Form.Item>
                    <Form.Item label={`judul unit ${descIndex + 1}`}>
                      <Controller
                        name={`softSkills.${index}.description.${descIndex}.unit_title`}
                        control={control}
                        render={({ field }) => (
                          <Input
                            placeholder="Masukkan unit title"
                            {...field}
                            style={{ width: "100%", height: "50px" }}
                          />
                        )}
                      />
                    </Form.Item>
                    <Button
                      type="text"
                      danger
                      icon={<MinusCircleOutlined />}
                      onClick={() => {
                        // Salin field yang sedang diperbarui tanpa merubah field lainnya
                        const updatedField = {
                          ...softSkillsFields[index], // Salin seluruh field
                          description: softSkillsFields[
                            index
                          ].description.filter(
                            (_, i) => i !== descIndex // Hapus deskripsi pada index tertentu
                          ),
                        };

                        // Update array hardSkillsFields tanpa merubah elemen lain
                        const updatedFields = softSkillsFields.map(
                          (field, idx) => (idx === index ? updatedField : field)
                        );

                        // Panggil upHardSkill untuk memperbarui state
                        upSoftSkill(index, updatedField);
                      }}
                    >
                      Hapus Deskripsi
                    </Button>
                  </div>
                ))}
                <Button
                  id={`softSkills.${index}.description`}
                  type="dashed"
                  htmlType="button"
                  onClick={() => {
                    upSoftSkill(index, {
                      ...getValues("softSkills")[index],
                      description: [
                        ...getValues("softSkills")[index].description,
                        { id: "", unit_code: "", unit_title: "" },
                      ],
                    });
                  }}
                  icon={<PlusOutlined />}
                  style={{ marginBottom: "20px" }}
                >
                  Tambah Unit Code dan Title
                </Button>
              </Space>
            </div>
          );
        })}

        <Button
          type="dashed"
          onClick={() =>
            addSoftSkill({
              skill_name: "",
              description: [{ unit_code: "", unit_title: "" }],
            })
          }
          block
          icon={<PlusOutlined />}
          style={{ marginBottom: "20px" }}
        >
          Tambah Soft Skill
        </Button>

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

export default Tool;
