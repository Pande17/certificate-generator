import { useEffect, useState } from "react";
import { Kompetensi } from "../api middleware";
import {
  message,
  Table,
  Col,
  Row,
  Button,
  Modal,
  Form,
  Input,
  Space,
  Select,
} from "antd";
import {
  DeleteOutlined,
  EditOutlined,
  MinusCircleOutlined,
  PlusOutlined,
} from "@ant-design/icons";
import MainLayout from "../MainLayout/Layout";
import { Navigate, useNavigate } from "react-router-dom";
import { useForm, Controller, useFieldArray, } from "react-hook-form";

const competence = () => {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState([]);
  const [searchText, setSearchText] = useState("");
  const [isEditModalVisible, setIsEditModalVisible] = useState(false);
  const [currentRecord, setCurrentRecord] = useState(null);
  const [competencies, setCompetencies] = useState([]);

    const { reset } = useForm();
  const { confirm } = Modal;

  useEffect(() => {
    const fetchingData = async () => {
      setLoading(true);
      try {
        const response = await Kompetensi.get(
          `/`
        );
        const datas = response.data.data;
        const filteredData = datas.filter((item) => !item.deleted_at);
        setData(filteredData);
      } catch (err) {
        console.error("error : ", err);
        message.error("error : ", err);
      } finally {
        setLoading(false);
      }
    };
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
    fetchingData();
  }, []);

  const navigate = useNavigate();

  const backHandle = () => {
    navigate("/competence");
  };

  const { control, handleSubmit } = useForm({
    defaultValues: {
      skkni: "",
      divisi: "",
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

  const handleEdit = (record) => {
    setCurrentRecord(record);
    setIsEditModalVisible(true);
  };


  const filteredData = data.filter((item) =>
    item.nama_kompetensi.toLowerCase().includes(searchText.toLowerCase())
  );
  
  const handleSearch = (e) => {
    setSearchText(e.target.value);
  };


 const delHandle = async (_id) => {
   try {
   await Kompetensi.delete(
       `/${_id}`
     );
     setData((prevData) => prevData.filter((item) => item._id !== _id));
     message.success("Data berhasil dihapus");
   } catch (error) {
     console.error("Error response:", error.response);
     message.error(
       `Gagal menghapus data: ${error.response?.data?.message || error.message}`
     );
   }
 };

[]

  const delConfirm = (_id, nama_kompetensi) => {
    confirm({
      title: `apakah anda yakon ingin menghapus kompetensi ${nama_kompetensi}`,
      content: "data yang di hapus tidak dapat dikembalikan",
      okType: "danger",
      okText: "ya, Hapus",
      cancelText: "Batal",
      onOk() {
        delHandle(_id);
      },
      onCancel() {
        console.log("penghapusan dibatalkan");
      },
    });
  };   

  const onSubmit = async (data) => {
    const competenceData = {
      nama_kompetensi: data.competenceName,
      divisi: data.devisi,
      skkni: data.skkni,
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
  };

  const column = [
    {
      title: "No",
      align: "center",
      width: 100,
      ellipsis: true,
      responsive: ["xs", "sm", "md", "lg"],
      render: (text, record, index) => index + 1,
    },
    {
      title: "Kompetensi",
      align: "center",
      responsive: ["xs", "sm", "md", "lg"],
      key: "nama_kompetensi",
      dataIndex: "nama_kompetensi",
      ellipsis: true,
    },
    {
      width: 300,
      title: "Aksi",
      key: "actions",
      responsive: ["xs", "sm", "md", "lg"],
      align: "center",
      render: (text, record) => {
        return (
          <div>
            <Button
              icon={<DeleteOutlined />}
              style={{ margin: 8 }}
              type="primary"
              danger
              onClick={() => delConfirm(record._id, record.nama_kompetensi)}
            />
            <Button
              icon={<EditOutlined />}
              type="primary"
              style={{ margin: 8 }}
              onClick={() => handleEdit(record)}
            />
          </div>
        );
      },
    },
  ];

  const createNav = () => {
    navigate("/competence/create-competence");
  };
  return (
    <MainLayout>
      <div className="flex flex-col items-center justify-center w-full lg:w-3/4 p-5">
        <p className="text-xl font-Poppins font-semibold mb-5 text-Text p-3 bg-white rounded-xl">
          Daftar Kompetensi
        </p>
        <Button onClick={createNav} className="m-3">
          Buat Kompetensi
        </Button>
        <input
          type="text"
          placeholder="Search"
          value={searchText}
          onChange={handleSearch}
          className="mb-4 p-2 border border-gray-300 rounded w-full md:w-1/2"
        />
        <Row
          style={{ justifyContent: "center", width: "100%", overflowX: "auto" }}
        >
          <Col>
            <Table
              dataSource={filteredData}
              columns={column}
              rowKey="_id"
              pagination={false}
              bordered
              loading={loading}
              scroll={{
                x: "max-content",
                y: filteredData.length > 6 ? 500 : undefined,
              }}
            />
          </Col>
        </Row>
      </div>

      {/* Modal Edit Sertif */}
      <Modal
        title="Edit Kompetensi"
        open={isEditModalVisible}
        onCancel={() => setIsEditModalVisible(false)}
        footer={null}
      >
        <Form
          layout="vertical"
          onFinish={handleSubmit(onSubmit)}
          style={{
            width: "90%",
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
      </Modal>
      {/*  END */}
    </MainLayout>
  );
};

export default competence;
