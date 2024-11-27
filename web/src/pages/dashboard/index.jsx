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

  const calculateTotalSkillScore = (hardSkills, softSkills) => {
    const totalHardSkillsScore = Array.isArray(hardSkills)
      ? hardSkills.reduce((acc, skill) => acc + (skill.skill_score || 0), 0)
      : 0; // Pastikan hardSkills adalah array, jika tidak, set default ke 0

    const totalSoftSkillsScore = Array.isArray(softSkills)
      ? softSkills.reduce((acc, skill) => acc + (skill.skill_score || 0), 0)
      : 0; // Pastikan softSkills adalah array, jika tidak, set default ke 0

    return totalHardSkillsScore + totalSoftSkillsScore;
  };

  const onSubmit = async (formData) => {
    console.log(formData); // Periksa formData yang diterima

    const totalSkillScore = calculateTotalSkillScore(
      formData.hardSkill, // Pastikan ini adalah array
      formData.softSkill // Pastikan ini adalah array
    );

    try {
      const formattedData = {
        savedb: true,
        page_name: "page2a",
        zoom: 1.367,
        data: {
          sertif_name: formData.sertifikat,
          nama_peserta: formData.nama,
          kompeten_bidang: formData.fieldOfStudy,
          kompetensi: dta.find(
            (item) => item._id === formData.selectedCompetenceId
          )?.nama_kompetensi,
          meet_time: formData.meetingTime,
          skkni: formData.skkni,
          valid_date: {
            valid_start: formData.expiredTimeStard?.format("DD MMMM YYYY"),
            valid_end: formData.expiredTimeEnd?.format("DD MMMM YYYY"),
            valid_total: formData.validtime,
          },
          total_meet: formData.totalMeeting,
          kode_referral: {
            referral_id: formData.codeReferralOrder,
            divisi: formData.codeReferralFieldOfStudy,
            bulan_rilis: formData.codeReferralMonth,
            tahun_rilis: formData.codeReferralYear,
          },
          hard_skills: {
            skills: Array.isArray(formData.hardSkill)
              ? formData.hardSkill.map((skill) => ({
                  skill_name: skill.skill_name,
                  skill_jp: skill.jp,
                  description: skill.combined_units.split("\n").map((line) => {
                    const [unit_code, unit_title] = line.split(" - ");
                    return { unit_code, unit_title };
                  }),
                }))
              : [],
            total_skill_jp:
              formData.hardSkill?.reduce(
                (acc, skill) => acc + (skill.jp || 0),
                0
              ) || 0,
            total_skill_score: totalSkillScore, // Replace with actual computation if necessary
          },
          soft_skills: {
            skills: Array.isArray(formData.softSkill)
              ? formData.softSkill.map((skill) => ({
                  skill_name: skill.skill_name,
                  skill_jp: skill.jp,
                  skill_score: skill.skill_score,
                  description: skill.combined_units.split("\n").map((line) => {
                    const [unit_code, unit_title] = line.split(" - ");
                    return { unit_code, unit_title };
                  }),
                }))
              : [],
            total_skill_jp:
              formData.softSkill?.reduce(
                (acc, skill) => acc + (skill.jp || 0),
                0
              ) || 0,
            total_skill_score: totalSkillScore, // Replace with actual computation if necessary
          },
          signature: {},
          total_jp:
            (formData.hardSkill?.reduce(
              (acc, skill) => acc + (skill.jp || 0),
              0
            ) || 0) +
            (formData.softSkill?.reduce(
              (acc, skill) => acc + (skill.jp || 0),
              0
            ) || 0),
        },
      };


      const response = await Sertifikat.put("/", formattedData);

      if (response.status === 200) {
        console.log(dta);
        message.success("Certificate added successfully!");
        reset(); // Clear the form
      }
    } catch (error) {
      console.log(dta);
      console.log("Error adding certificate:", error);
      message.error("Failed to add certificate. Please try again.");
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
  setCurrentRecord(record);
  setIsEditModalVisible(true); // Menampilkan modal edit

  try {
    const response = await Sertifikat.get(`/certificate/${record._id}`); // Ganti dengan endpoint yang benar
    const certificateData = response.data.data;
         console.log("Data yang diambil dari API:", certificateData.data);
    // Setel data form ke state atau form
    reset({
      sertifikat: certificateData.sertif_name,
      nama: certificateData.nama_peserta,
      fieldOfStudy: certificateData.kompeten_bidang,
      selectedCompetenceId: certificateData.kompetensi,
      meetingTime: certificateData.meet_time,
      skkni: certificateData.skkni,
      expiredTimeStard: moment(certificateData.valid_date?.valid_start),
      expiredTimeEnd: moment(certificateData.valid_date?.valid_end),
      validtime: certificateData.valid_date?.valid_total,
      totalMeeting: certificateData.total_meet,
      codeReferralOrder: certificateData.kode_referral?.referral_id,
      codeReferralFieldOfStudy: certificateData.kode_referral?.divisi,
      codeReferralMonth: certificateData.kode_referral?.bulan_rilis,
      codeReferralYear: certificateData.kode_referral?.tahun_rilis,
      hardSkill: certificateData.hard_skills?.skills || [],
      softSkill: certificateData.soft_skills?.skills || [],
    });

    // Jika perlu, Anda bisa menambahkan logika untuk memuat data lainnya (seperti signature, kompetensi, dll.)
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
    // Reset and update hard and soft skills upon competence change
    reset({
      selectedCompetenceId: value,
      hardSkill: [],
      softSkill: [],
    });
    fetchCompetence(value);
  };

  const filteredData = dta.filter((item) =>
    item.sertif_name.toLowerCase().includes(searchText.toLowerCase())
  );

  const downloadPDF = async (_id) => {
    try {
      const response = await Sertifikat.get(`/${_id}/b`, {
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
      title: "ID",
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

        <Row style={{ width: "100%", overflowX: "auto" }}>
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
              style={{ width: "100%" }}
            />
          </Col>
        </Row>
        {/* Modal for Edit */}
        <Modal
          title="Edit Sertifikat"
          open={isEditModalVisible}
          onCancel={() => setIsEditModalVisible(false)}
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
                defaultValue={currentRecord?.sertif_name || ""}
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
                defaultValue={
                  currentRecord?.nama_peserta || " nama peserta tidak ada"
                }
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
                defaultValue={
                  currentRecord?.kompeten_bidang ||
                  "tidak mengisi bidang tidak ada"
                }
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
                defaultValue={
                  currentRecord?.valid_date?.valid_total ||
                  "tidak mengisi Total Tahun tidak ada"
                }
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
                defaultValue={
                  currentRecord?.valid_date?.valid_start
                    ? moment(
                        currentRecord.valid_date.valid_start,
                        "DD MMMM YYYY"
                      )
                    : null
                }
                control={control}
                rules={{ required: "Waktu mulai diperlukan" }}
                render={({ field }) => (
                  <DatePicker
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
                defaultValue={
                  currentRecord?.valid_date?.valid_end
                    ? moment(currentRecord.valid_date.valid_end, "DD MMMM YYYY")
                    : null
                }
                control={control}
                rules={{ required: "Waktu selesai diperlukan" }}
                render={({ field }) => (
                  <DatePicker
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
                defaultValue={
                  currentRecord?.total_meet ||
                  "tidak mengisi Total Pertemuan tidak ada"
                }
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
                defaultValue={
                  currentRecord?.meet_time ||
                  "tidak mengisi Waktu Pertemuan tidak ada"
                }
                control={control}
                rules={{ required: "Waktu pertemuan diperlukan" }}
                render={({ field }) => (
                  <InputNumber
                    {...field}
                    placeholder="Masukkan waktu pertemuan"
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
                defaultValue={currentRecord?.kompeten_bidang}
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
                      defaultValue={currentRecord?.kompeten_bidang}
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
                      defaultValue={currentRecord?.kompeten_bidang}
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
                      defaultValue={currentRecord?.kompeten_bidang}
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
                      defaultValue={currentRecord?.kompeten_bidang}
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
                defaultValue={currentRecord?.kompeten_bidang}
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
