import {
  Table,
  Button,
  Modal,
  message,
  Row,
  Col,
  Form,
  Input,
  DatePicker,
  InputNumber,
  Select,
} from "antd";
import { useForm, Controller, useFieldArray } from "react-hook-form";
import { Sertifikat, Kompetensi, Signature } from "../api middleware";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  DeleteOutlined,
  EditOutlined,
  DownloadOutlined,
} from "@ant-design/icons";
import MainLayout from "../MainLayout/Layout";

const { confirm } = Modal;

const Dashboard = () => {
  const [dta, setDta] = useState([]);
  const [loading, setLoading] = useState(false);
  const [searchText, setSearchText] = useState("");
  const [isEditModalVisible, setIsEditModalVisible] = useState(false);
  const [currentRecord, setCurrentRecord] = useState(null);
  const [signatureData, setSignatureData] = useState([]);
  const [kompetensiData, setKompetensiData] = useState([]);
  const [skkni, setSkkni] = useState("");
  const [divisi, setDivisi] = useState("");
  const { control, handleSubmit, reset } = useForm({
    defaultValues: {
      hardSkill: [],
      softSkill: [],
      selectedCompetenceId: "",
    },
  });

  const navigate = useNavigate();

  // Fetch data dari API
  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await Sertifikat.get("/");
        const certificates = response.data.data || [];
        const filteredData = certificates.filter((item) => !item.deleted_at);
        setDta(filteredData);
      } catch (err) {
        console.error("Error fetching data:", err);
        message.error("Gagal memuat data.");
      }
    };
    const fetchSignature = async () => {
      try {
        const response = await Signature.get("/");
        setSignatureData(response.data.data);
      } catch (error) {
        console.log("Error fetching signature data:", error);
      }
    };
    const fetchCompetence = async () => {
      try {
        const response = await Kompetensi.get("/");
        setKompetensiData(response.data.data);
      } catch (error) {
        console.log("Error fetching competence data:", error);
      }
    };

    fetchCompetence();
    fetchSignature();
    fetchData();
  }, []);

  const { fields: hardSkillFields, replace: replaceHardSkill } = useFieldArray({
    control,
    name: "hardSkill",
  });

  const { fields: softSkillFields, replace: replaceSoftSkill } = useFieldArray({
    control,
    name: "softSkill",
  });

 const onSubmit = async (formData) => {
   try {
     const sanitizedHardSkills = (formData.hardSkill || []).filter(
       (skill) => skill.skill_name && skill.jp
     );
     const sanitizedSoftSkills = (formData.softSkill || []).filter(
       (skill) => skill.skill_name && skill.skill_score
     );

     const formattedData = {
       ...formData,
       hardSkill: sanitizedHardSkills,
       softSkill: sanitizedSoftSkills,
     };

     const response = await Sertifikat.put(
       `/${currentRecord._id}`,
       formattedData
     );

     if (response.status === 200) {
       message.success("Sertifikat berhasil diperbarui!");
       reset();
       setIsEditModalVisible(false); // Tutup modal
     }
   } catch (error) {
     console.error("Error updating certificate:", error);
     message.error("Gagal memperbarui sertifikat.");
   }
 };

  const { Option } = Select;

  const deleteCompetence = async (_id) => {
    try {
      await Sertifikat.delete(`/${_id}`);
      setDta((prevDta) => prevDta.filter((item) => item._id !== _id));
      message.success("SERTIFIKAT berhasil dihapus!");
    } catch (error) {
      message.error("Gagal menghapus SERTIFIKAT.");
    }
  };

  const handleSearch = (e) => {
    setSearchText(e.target.value);
  };

  const showDeleteConfirm = (_id) => {
    confirm({
      title: "Apakah Anda yakin ingin menghapus SERTIFIKAT ini?",
      content: "Data yang dihapus tidak dapat dikembalikan.",
      okText: "Ya, Hapus",
      okType: "danger",
      cancelText: "Batal",
      onOk() {
        deleteCompetence(_id);
      },
      onCancel() {
        console.log("Penghapusan dibatalkan");
      },
    });
  };

  const handleEdit = async (record) => {
    try {
      const response = await Sertifikat.get(`/${record._id}`); // Ganti dengan endpoint yang benar
      const certificateData = response.data.data.data;
      console.log("Data yang diambil dari API:", certificateData);

      setCurrentRecord(certificateData);
      setIsEditModalVisible(true);
    } catch (error) {
      console.error("Error fetching certificate details:", error);
      message.error("Gagal mengambil data sertifikat.");
    }
  };

  const fetchCompetence = async (competenceId) => {
    const url = `/${competenceId}`;
    try {
      const response = await Kompetensi.get(url);

      const {
        hard_skills = [],
        soft_skills = [],
        skkni = "",
        divisi = "",
      } = response.data.data || {};

      const newHardSkills = hard_skills.map((hardSkill) => ({
        skill_name: hardSkill.skill_name || "",
        combined_units: hardSkill.description
          .map((unit) => `${unit.unit_code} - ${unit.unit_title}`)
          .join("\n"),
      }));

      const newSoftSkills = soft_skills.map((softSkill) => ({
        skill_name: softSkill.skill_name || "",
        combined_units: softSkill.description
          .map((unit) => `${unit.unit_code} - ${unit.unit_title}`)
          .join("\n"),
      }));

      replaceHardSkill(newHardSkills);
      replaceSoftSkill(newSoftSkills);

      // Simpan skkni dan divisi ke state
      setSkkni(skkni);
      setDivisi(divisi);
    } catch (err) {
      console.log(err);
    }
  };

  const handleCompetenceChange = (value) => {
    // Hanya reset field terkait
    reset((prevValues) => ({
      ...prevValues, // Pertahankan nilai sebelumnya
      selectedCompetenceId: value,
      hardSkill: [],
      softSkill: [],
    }));
    fetchCompetence(value);
  };

  const filteredData = dta.filter((item) =>
    item.sertif_name.toLowerCase().includes(searchText.toLowerCase())
  );

  const downloadPDF = async (_id) => {
    try {
      const response = await Sertifikat.get(`/download/${_id}/b`, {
        headers: {
          "Content-Type": "application/pdf",
        },
        responseType: "blob",
      });

      // Membuat link untuk mengunduh file
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", `${_id}.pdf`); // Nama file saat diunduh
      document.body.appendChild(link);
      link.click();
      link.remove(); // Hapus link setelah digunakan
    } catch (error) {
      console.error("Error downloading PDF:", error);
    }
  };

  const createNav = () => {
    navigate("/create");
  };

  const columns = [
    {
      title: "No",
      align: "center",
      width: 100,
      responsive: ["xs", "sm", "md", "lg"],
      ellipsis: true,
      render: (text, record, index) => index + 1,
    },
    {
      title: "Daftar Sertifikat",
      dataIndex: "sertif_name",
      key: "sertif_name",
      responsive: ["xs", "sm", "md", "lg"],
      ellipsis: true,
    },
    {
      title: "Aksi",
      key: "actions",
      align: "center",
      width: 300,
      responsive: ["xs", "sm", "md", "lg"],
      render: (text, record) => (
        <div>
          <Button
            icon={<DeleteOutlined />}
            style={{ margin: 8 }}
            type="primary"
            danger
            onClick={() => showDeleteConfirm(record._id)}
          />
          <Button
            icon={<EditOutlined />}
            style={{ margin: 8 }}
            type="primary"
            onClick={() => handleEdit(record)}
          />
          <Button
            icon={<DownloadOutlined />}
            type="primary"
            style={{ margin: 8 }}
            onClick={() => downloadPDF(record.data_id)}
          />
        </div>
      ),
    },
  ];

  return (
    <MainLayout>
      <div className="flex flex-col items-center justify-center w-full lg:w-3/4 p-5">
        <div>
          <p className="text-xl font-Poppins font-semibold mb-5 text-Text p-3 bg-white rounded-xl">
            List Sertifikat
          </p>
        </div>

        <Button onClick={createNav} className="m-3">
          Buat Sertifikat
        </Button>

        <input
          type="text"
          placeholder="Search"
          value={searchText}
          onChange={handleSearch}
          className="mb-4 p-2 border border-gray-300 rounded w-full md:w-1/2"
        />

        <Row  style={{ justifyContent: "center", width: "100%", overflowX: "auto" }}>
          <Col span={24}>
            <Table
              dataSource={filteredData}
              columns={columns}
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
        {/* Modal for Edit */}
        <Modal
          title="Edit Sertifikat"
          open={isEditModalVisible}
          onCancel={() => setIsEditModalVisible(false)}
          afterOpenChange={(visible) => {
            if (visible && currentRecord) {
              reset({
                sertifikat: currentRecord?.sertif_name || "Tidak mengisi",
                nama: currentRecord?.nama_peserta || "Tidak mengisi",
                fieldOfStudy: currentRecord?.kompeten_bidang || "Tidak mengisi",
                validTime:
                  currentRecord?.valid_date?.valid_total || "Tidak mengisi",
                expiredTimeStart:
                  currentRecord?.valid_date?.valid_start || "Tidak mengisi",
                expiredTimeEnd:
                  currentRecord?.valid_date?.valid_end || "Tidak mengisi",
                totalMeeting: currentRecord?.total_meet || "Tidak mengisi",
                meetingTime: currentRecord?.meet_time || "Tidak mengisi",
                selectedCompetenceId:
                  currentRecord?.kompeten_bidang || "Tidak mengisi", // Atur nilai awal di sini
                selectedSignatureId:
                  currentRecord?.kompeten_bidang || "Tidak mengisi",
                hardSkill: currentRecord?.hardSkills || [],
                softSkill: currentRecord?.softSkills || [],
              });
            }
          }}
          footer={null}
        >
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
            <Form.Item label="Nama Sertifikat" required>
              <Controller
                name="sertifikat"
                control={control}
                rules={{ required: "Nama sertifikat diperlukan" }}
                render={({ field }) => (
                  <Input
                    {...field}
                    placeholder="Masukkan nama sertifikat"
                    style={{ width: "100%", height: "50px" }}
                  />
                )}
              />
            </Form.Item>

            <Form.Item label="Nama" required>
              <Controller
                name="nama"
                control={control}
                rules={{ required: "Nama peserta diperlukan" }}
                render={({ field }) => (
                  <Input
                    {...field}
                    placeholder="Masukkan nama peserta"
                    style={{ width: "100%", height: "50px" }}
                  />
                )}
              />
            </Form.Item>

            <Form.Item label="Bidang Studi" required>
              <Controller
                name="fieldOfStudy"
                control={control}
                rules={{ required: "Bidang studi diperlukan" }}
                render={({ field }) => (
                  <Input
                    {...field}
                    placeholder="Masukkan bidang studi"
                    style={{ width: "100%", height: "50px" }}
                  />
                )}
              />
            </Form.Item>

            <Form.Item label="Total Tahun" required>
              <Controller
                name="validTime"
                control={control}
                rules={{ required: "Waktu validasi diperlukan" }}
                render={({ field }) => (
                  <Input
                    {...field}
                    placeholder="Masukkan jumlah tahun"
                    style={{ width: "100%", height: "50px" }}
                  />
                )}
              />
            </Form.Item>

            <Form.Item label="Waktu Expired (Mulai)" required>
              <Controller
                name="expiredTimeStart"
                control={control}
                rules={{ required: "Waktu mulai diperlukan" }}
                render={({ field }) => (
                  <Input
                    {...field}
                    placeholder="Pilih waktu mulai"
                    style={{ width: "100%", height: "50px" }}
                  />
                )}
              />
            </Form.Item>

            <Form.Item label="Waktu Expired (Selesai)" required>
              <Controller
                name="expiredTimeEnd"
                control={control}
                rules={{ required: "Waktu selesai diperlukan" }}
                render={({ field }) => (
                  <Input
                    {...field}
                    placeholder="Pilih waktu selesai"
                    style={{ width: "100%", height: "50px" }}
                  />
                )}
              />
            </Form.Item>

            <Form.Item label="Total Pertemuan" required>
              <Controller
                name="totalMeeting"
                control={control}
                rules={{ required: "Jumlah pertemuan diperlukan" }}
                render={({ field }) => (
                  <InputNumber
                    {...field}
                    placeholder="Masukkan jumlah pertemuan"
                    style={{ width: "100%", height: "50px" }}
                  />
                )}
              />
            </Form.Item>

            <Form.Item label="Waktu Pertemuan" required>
              <Controller
                name="meetingTime"
                control={control}
                rules={{ required: "Waktu pertemuan diperlukan" }}
                render={({ field }) => (
                  <Input
                    {...field}
                    placeholder="2 Bulan"
                    style={{ width: "100%", height: "50px" }}
                  />
                )}
              />
            </Form.Item>

            <h1 className="text-center font-Poppins text-2xl font-medium p-6">
              Pilih kompetensi
            </h1>
            <Form.Item required>
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
                      handleCompetenceChange(value);
                    }}
                  >
                    <Option value="" disabled>
                      pilih kommpetensi
                    </Option>
                    {kompetensiData.map((competence) => (
                      <Option key={competence._id} value={competence._id}>
                        {competence.nama_kompetensi || ""}
                      </Option>
                    ))}
                  </Select>
                )}
              />
            </Form.Item>

            {skkni && (
              <Form.Item label="SKKNI">
                <Input
                  value={skkni}
                  readOnly
                  style={{ width: "100%", height: "50px" }}
                />
              </Form.Item>
            )}

            {divisi && (
              <Form.Item label="Divisi">
                <Input
                  value={divisi}
                  readOnly
                  style={{ width: "100%", height: "50px" }}
                />
              </Form.Item>
            )}

            {hardSkillFields.length > 0 && (
              <div>
                <h2 className="font-Poppins text-2xl font-medium text-center p-6">
                  Hardskills
                </h2>
                {hardSkillFields.map((skill, index) => (
                  <div key={index} style={{ marginBottom: "20px" }}>
                    <label>{`Hardskill ${index + 1}`}</label>

                    {/* Skill Name Input */}
                    <Controller
                      name={`hardSkill[${index}].skill_name`}
                      control={control}
                      render={({ field }) => (
                        <Input
                          {...field}
                          placeholder="Skill Name"
                          readOnly
                          style={{
                            marginBottom: "10px",
                            width: "100%",
                            height: "50px",
                          }}
                        />
                      )}
                    />

                    {/* Unit Code and Title Input */}
                    <Controller
                      name={`hardSkill[${index}].combined_units`}
                      control={control}
                      render={({ field }) => (
                        <Input.TextArea
                          {...field}
                          rows={4}
                          placeholder="Unit Code and Title"
                          readOnly
                          style={{
                            marginBottom: "10px",
                            width: "100%",
                          }}
                        />
                      )}
                    />

                    {/* JP Input for each hard skill */}
                    <Controller
                      name={`hardSkill[${index}].jp`}
                      control={control}
                      render={({ field }) => (
                        <InputNumber
                          {...field}
                          placeholder="JP per skill"
                          style={{
                            width: "100%",
                            height: "50px",
                          }}
                        />
                      )}
                    />
                    <Controller
                      name={`hardSkill[${index}].skillScore`}
                      control={control}
                      render={({ field }) => (
                        <InputNumber
                          {...field}
                          placeholder="Score"
                          style={{
                            width: "100%",
                            height: "50px",
                          }}
                        />
                      )}
                    />
                  </div>
                ))}
              </div>
            )}

            {softSkillFields.length > 0 && (
              <div>
                <h2 className="font-Poppins text-2xl font-medium text-center p-6">
                  Softskills
                </h2>
                {softSkillFields.map((skill, index) => (
                  <div key={index} style={{ marginBottom: "20px" }}>
                    <label>{`Softskill ${index + 1}`}</label>

                    {/* Skill Name Input */}
                    <Controller
                      name={`softSkill[${index}].skill_name`}
                      control={control}
                      render={({ field }) => (
                        <Input
                          {...field}
                          placeholder="Skill Name"
                          readOnly
                          style={{
                            marginBottom: "10px",
                            width: "100%",
                            height: "50px",
                          }}
                        />
                      )}
                    />

                    {/* Unit Code and Title Input */}
                    <Controller
                      name={`softSkill[${index}].combined_units`}
                      control={control}
                      render={({ field }) => (
                        <Input.TextArea
                          {...field}
                          rows={4}
                          placeholder="Unit Code and Title"
                          readOnly
                          style={{
                            marginBottom: "10px",
                            width: "100%",
                          }}
                        />
                      )}
                    />

                    <Controller
                      name={`softSkill[${index}].jp`}
                      control={control}
                      render={({ field }) => (
                        <InputNumber
                          {...field}
                          placeholder="JP per skill"
                          style={{
                            width: "100%",
                            height: "50px",
                          }}
                        />
                      )}
                    />
                    <Controller
                      name={`softSkill[${index}].skillScore`}
                      control={control}
                      render={({ field }) => (
                        <InputNumber
                          {...field}
                          placeholder="score"
                          style={{
                            width: "100%",
                            height: "50px",
                          }}
                        />
                      )}
                    />
                  </div>
                ))}
              </div>
            )}

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
                        {signature.config_name}
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
        </Modal>
      </div>
    </MainLayout>
  );
};

export default Dashboard;
