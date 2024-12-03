import { useEffect, useState } from "react";
import { useForm, Controller, useFieldArray } from "react-hook-form";
import { Form, Input, Button, Space, message, Select } from "antd";
import { PlusOutlined, MinusCircleOutlined, RotateLeftOutlined, BackwardFilled } from "@ant-design/icons";
import MainLayout from "../MainLayout/Layout";
import { useNavigate } from "react-router-dom";
import { Kompetensi } from "../api middleware";

const { Option } = Select;

const Tool = () => {
  const navigate = useNavigate();

  const backHandle = () => {
    navigate("/competence");
  }

  const { control, handleSubmit, reset, watch } = useForm({
    defaultValues: {
      skkni:"",
      divisi:"",
      competenceName: "",
      hardSkills: [
        { skill_name: "", description: [{ unit_code: "", unit_title: "" }] },
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
    remove: removeHardSkill,
  } = useFieldArray({ control, name: "hardSkills" });

  const {
    fields: softSkillsFields,
    append: addSoftSkill,
    remove: removeSoftSkill,
  } = useFieldArray({ control, name: "softSkills" });

  const [competencies, setCompetencies] = useState([]);
  const selectedCompetenceId = watch("selectedCompetenceId");

  // Fetch competencies from the API
  useEffect(() => {
    const fetchCompetencies = async () => {
      try {
        const response = await Kompetensi.get(
          "/"
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

  // Set form values when a competence is selected
  useEffect(() => {
    const selectedCompetence = competencies.find(
      (c) => c._id === selectedCompetenceId
    );
    if (selectedCompetence) {
      reset({
        competenceName: selectedCompetence.nama_kompetensi,
        hardSkills: selectedCompetence.hard_skills || [],
        softSkills: selectedCompetence.soft_skills || [],
        selectedCompetenceId,
      });
    } else {
      reset();
    }
  }, [selectedCompetenceId, competencies, reset]);

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
        await Kompetensi.put(
          `/${data.selectedCompetenceId}`,
          competenceData
        );
        message.success("Kompetensi berhasil diperbarui!");
      } else {
        await Kompetensi.post(
          "/",
          competenceData
        );
        message.success("Kompetensi berhasil ditambahkan!");
      }
      reset();
    } catch (error) {
      console.error("Error saat menyimpan kompetensi:", error);
      message.error("Error saat menyimpan kompetensi!");
    }
    navigate("/competence")
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
                "Input divisi berlebihan atau kurang dari satu! maksimal(1-3 huruf)",
              validate: (value) =>
                value.length <= 6 ||
                "Input divisi berlebihan atau kurang dari satu!",
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
          const {
            fields: descriptionFields,
            append: addDescription,
            remove: removeDescription,
          } = useFieldArray({
            control,
            name: `hardSkills.${index}.description`,
          });

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
                {descriptionFields.map((descField, descIndex) => (
                  <div key={descField.id}>
                    <Form.Item label="Unit Code">
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
                    <Form.Item label="Unit Title">
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
                      onClick={() => removeDescription(descIndex)}
                    >
                      Hapus Deskripsi
                    </Button>
                  </div>
                ))}
              </Space>

              <Button
                type="dashed"
                onClick={() =>
                  addDescription({ unit_code: "", unit_title: "" })
                }
                icon={<PlusOutlined />}
                style={{ marginBottom: "20px" }}
              >
                Tambah Unit Code dan Title
              </Button>
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
          const {
            fields: descriptionFields,
            append: addDescription,
            remove: removeDescription,
          } = useFieldArray({
            control,
            name: `softSkills.${index}.description`,
          });

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
                {descriptionFields.map((descField, descIndex) => (
                  <div key={descField.id}>
                    <Form.Item label="Unit Code">
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
                    <Form.Item label="Unit Title">
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
                      onClick={() => removeDescription(descIndex)}
                    >
                      Hapus Deskripsi
                    </Button>
                  </div>
                ))}
              </Space>

              <Button
                type="dashed"
                onClick={() =>
                  addDescription({ unit_code: "", unit_title: "" })
                }
                icon={<PlusOutlined />}
                style={{ marginBottom: "20px" }}
              >
                Tambah Unit Code dan Title
              </Button>
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
