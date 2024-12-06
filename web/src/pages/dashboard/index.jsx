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
  Spin,
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
  const [selectedSignature, setSelectedSignature] = useState(null);
  const [isSignatureSelected, setIsSignatureSelected] = useState(false);
  const [kompetensiData, setKompetensiData] = useState([]);
  const [skkni, setSkkni] = useState("");
  const [divisi, setDivisi] = useState("");
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [selectedDownload, setSelectedDownload] = useState(null);
  const { control, handleSubmit, reset, setValue, watch } = useForm({
    defaultValues: {
      hardSkill: [],
      softSkill: [],
      selectedCompetenceId: "",
    },
  });

  console.log(watch());

  const navigate = useNavigate();
  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const response = await Sertifikat.get("/");

        if (response.status === 200) {
          const certificates = response.data.data || [];
          const filteredData = certificates
            .filter((item) => !item.deleted_at)
            .sort((a, b) => new Date(b.updated_at) - new Date(a.updated_at));
          setDta(filteredData);
        } else {
          console.error("Error fetching data:", response.status);
          message.error("Gagal memuat data.");
        }
      } catch (err) {
        console.error("Error fetching data:", err);
        message.error("Gagal memuat data.");
      } finally {
        setLoading(false);
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
    const hardSkillScore = Array.isArray(hardSkills)
      ? hardSkills.reduce((acc, skill) => acc + (skill.skillScore || 0), 0)
      : 0;

    const softSkillScore = Array.isArray(softSkills)
      ? softSkills.reduce((acc, skill) => acc + (skill.skillScore || 0), 0)
      : 0;

    return hardSkillScore + softSkillScore;
  };

  const parseDescription = (skills) =>
    Array.isArray(skills)
      ? skills.map((skill) => ({
          ...skill,
          description: skill.combined_units
            ? skill.combined_units.split("\n").map((line) => {
                const [unit_code, unit_title] = line.split(" - ");
                return { unit_code, unit_title };
              })
            : [],
        }))
      : [];

  const onSubmit = async (formData) => {
    console.log("formData:", currentRecord?._id);

    if (!currentRecord?._id) {
      message.error("Certificate ID is missing!");
      return;
    }

    const totalSkillScore = calculateTotalSkillScore(
      formData.hardSkill,
      formData.softSkill
    );

    const selectedSignature = signatureData.find(
      (signature) => signature._id === currentRecord?.selectedSignatureId
    );

    const selectedCompetence = currentRecord?.Kompetensi?.find(
      (item) => item._id === currentRecord?.selectedCompetenceId
    );

    const parsedHardSkills = parseDescription(formData.hardSkill);
    const parsedSoftSkills = parseDescription(formData.softSkill);

    const formattedData = {
      savedb: true,
      page_name: "page2a",
      zoom: 1.367,
      data: {
        sertif_name: formData.sertifikat,
        nama_peserta: formData.nama,
        kompeten_bidang: formData.fieldOfStudy,
        kompetensi: formData.kompetensiDisplay,
        meet_time: formData.meetingTime,
        skkni: formData.skkni,
        validation: formData.validation,
        valid_date: {
          valid_start: formData.expiredTimeStart,
          valid_end: formData.expiredTimeEnd,
          valid_total: formData.validTime,
        },
        total_meet: formData.totalMeeting,
        kode_referral: {
          divisi: formData.codeReferralFieldOfStudy,
        },
        hard_skills: {
          skills: parsedHardSkills,
          total_skill_jp:
            formData.hardSkill?.reduce(
              (acc, skill) => acc + (skill.jp || 0),
              0
            ) || 0,
          total_skill_score: totalSkillScore,
        },
        soft_skills: {
          skills: parsedSoftSkills,
          total_skill_jp:
            formData.softSkill?.reduce(
              (acc, skill) => acc + (skill.jp || 0),
              0
            ) || 0,
          total_skill_score: totalSkillScore,
        },
        signature: {
          config_name: formData.configName || "",
          logo: formData.logo,
          role: formData.role,
          signature: formData.linkGambarPenandatangan,
          name: formData.namaPenandatangan,
          stamp: formData.stamp,
        },
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

    try {
      const response = await Sertifikat.put(
        `/${currentRecord?._id}`,
        formattedData
      );

      if (response.status === 200) {
        message.success("Data updated successfully!");
        reset();
      } else {
        console.error("Server Response:", response?.data || response);
        message.error("Failed to update data. Please try again.");
      }
    } catch (error) {
      console.error("Error updating data:", error?.response?.data || error);
      message.error("Failed to update data. Please try again.");
    }

    navigate(`/dashboard`);
  };

  const { Option } = Select;

  const deleteSertif = async (_id) => {
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
        deleteSertif(_id);
      },
      onCancel() {
        console.log("Penghapusan dibatalkan");
      },
    });
  };

  const handleEdit = async (record) => {
    try {
      const response = await Sertifikat.get(`/${record._id}`);

      console.log("API Response:", response);

      const primaryData = response.data.data.data;
      const secondData = response.data.data;

      const certificateData = {
        ...primaryData,
        ...secondData,
      };
      console.log({ secondData });

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
      setSkkni(skkni);
      setDivisi(divisi);
    } catch (err) {
      console.log("Error fetching competence details:", err);
    }
  };
  const handleCompetenceChange = async (value) => {
    try {
      const competence = kompetensiData.find((item) => item._id === value);
      if (competence) {
        setSkkni(competence.skkni || "");
        setDivisi(competence.divisi || "");

        // Reset form untuk hard skill dan soft skill
        reset((prevValues) => ({
          ...prevValues,
          selectedCompetenceId: value,
          hardSkill: [],
          softSkill: [],
        }));
        // Fetch kompetensi detail
        await fetchCompetence(value);
      }
    } catch (error) {
      console.error("Error handling competence change:", error);
    }
  };

  const fetchSignatureId = async (SignatureId) => {
    try {
      const response = await Signature.get(`/${SignatureId}`);
      return response.data.data;
    } catch (error) {
      console.error("Error fetching signature:", error);
      return null;
    }
  };

  const handleSignatureChange = async (value) => {
    const signature = await fetchSignatureId(value);
    console.log("Signature fetched:", signature);
    if (signature) {
      setSelectedSignature(signature);
      setIsSignatureSelected(true);

      setValue("namaPenandatangan", signature.name || "");
      setValue("role", signature.role || "");
      setValue("linkLogo", signature.logo || "");
      setValue("linkGambarPenandatangan", signature.signature || "");
      setValue("logoPerusahaan", signature.logo || "");
      setValue("stamp", signature.stamp || "");
    } else {
      setIsSignatureSelected(false);
    }
  };

  const filteredData = dta.filter((item) =>
    item.sertif_name.toLowerCase().includes(searchText.toLowerCase())
  );

  const downloadPDF = async (_id, type) => {
    setLoading(true);
    try {
      const response = await Sertifikat.get(`/download/${_id}/${type}`, {
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
    } finally {
      setLoading(false);
    }
  };

  const createNav = () => {
    navigate("/create");
  };

  const handleDownloadClick = (record) => {
    setSelectedDownload(record); // Simpan record yang dipilih ke dalam state
    setIsModalVisible(true); // Tampilkan modal
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
      dataIndex: "sertif_title",
      key: "sertif_title",
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
            onClick={() => handleDownloadClick(record)}
          />
        </div>
      ),
    },
  ];

  console.log(kompetensiData);
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

        <Row
          style={{
            justifyContent: "center",
            width: "100%",
            overflowX: "auto",
          }}
        >
          <Col span={24}>
            <Table
              dataSource={filteredData}
              columns={columns}
              rowKey="_id"
              pagination={false}
              bordered
              loading={loading}
              scroll={{
                x: "min-content",
                y: filteredData.length > 6 ? 400 : undefined,
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
              const matchedCompetence = kompetensiData.find(
                (item) => item.nama_kompetensi === currentRecord.kompetensi
              );

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
                selectedCompetenceId: matchedCompetence?._id || "", // Atur ID hasil pencocokan
                selectedSignatureId:
                  currentRecord?.signature?.config_name || "Tidak mengisi",
                namaPenandatangan: currentRecord?.signature?.name,
                role: currentRecord?.signature?.role,
                stamp: currentRecord?.signature?.stamp,
                linkLogo: currentRecord?.signature?.logo,
                linkGambarPenandatangan: currentRecord?.signature?.signature,
                hardSkill: currentRecord?.data?.hard_skills.skills || [],
                softSkill: currentRecord?.data?.soft_skills.skills || [],
                divisi: currentRecord?.data?.kode_referral.divisi || "",
                skkni: currentRecord?.data?.skkni || "",
                configName:
                  currentRecord?.signature?.config_name || "tidak megisi",
                kompetensiDisplay:currentRecord?.data?.kompetensi || "gak ada",
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

            <Form.Item label="Durasi Pertemuan" required>
              <Controller
                name="meetingTime"
                control={control}
                rules={{ required: "Durasi pertemuan diperlukan" }}
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
              Pilih Kompetensi
            </h1>
            <Controller
              name="selectedCompetenceId"
              control={control}
              render={({ field }) => (
                <Select
                  placeholder="Pilih Kompetensi"
                  {...field}
                  style={{ width: "100%", height: "50px" }}
                  onChange={(value) => {
                    field.onChange(value);
                    handleCompetenceChange(value);
                    setValue("selectedCompetence", value);
                  }}
                >
                  <Option value="" disabled>
                    Pilih Kompetensi
                  </Option>
                  {kompetensiData.map((competence) => (
                    <Option key={competence._id} value={competence._id}>
                      {competence.nama_kompetensi || ""}
                    </Option>
                  ))}
                </Select>
              )}
            />

            <Form.Item label="Kompetensi">
              <Controller
                name="kompetensiDisplay"
                control={control}
                rules={{ required: "kompetensi diperlukan" }}
                render={({ field }) => (
                  <Input
                    {...field}
                    readOnly
                    style={{ width: "100%", height: "50px" }}
                  />
                )}
              />
            </Form.Item>

            <Form.Item label="SKKNI">
              <Controller
                name="skkni"
                control={control}
                rules={{ required: "skkni peserta diperlukan" }}
                render={({ field }) => (
                  <Input
                    {...field}
                    readOnly
                    style={{ width: "100%", height: "50px" }}
                  />
                )}
              />
            </Form.Item>

            <Form.Item label="Divisi">
              <Controller
                name="divisi"
                control={control}
                rules={{ required: "skkni peserta diperlukan" }}
                render={({ field }) => (
                  <Input
                    {...field}
                    readOnly
                    style={{ width: "100%", height: "50px" }}
                  />
                )}
              />
            </Form.Item>

            <div>
              <h2 className="font-Poppins text-2xl font-medium text-center p-6">
                Hard Skills
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
                    name={`hardSkill[${index}].description`}
                    control={control}
                    render={({ field }) => (
                      <Input.TextArea
                        value={field.value
                          ?.map(
                            (desc, idx) =>
                              `${idx + 1}. ${desc.unit_code} - ${
                                desc.unit_title
                              }`
                          )
                          .join("\n")}
                        rows={4}
                        placeholder="Unit Code and Title"
                        readOnly
                        style={{
                          marginBottom: "10px",
                          width: "100%",
                        }}
                        onChange={(e) => {
                          // Split input back into array of objects
                          const updatedDescriptions = e.target.value
                            .split("\n")
                            .map((line) => {
                              const [unit_code, ...unit_title] =
                                line.split(" - ");
                              return {
                                unit_code: unit_code?.trim() || "",
                                unit_title:
                                  unit_title.join(" - ")?.trim() || "",
                              };
                            })
                            .filter(
                              (item) => item.unit_code || item.unit_title
                            );

                          field.onChange(updatedDescriptions);
                        }}
                      />
                    )}
                  />

                  {/* JP Input for each hard skill */}
                  <Controller
                    name={`hardSkill[${index}].skill_jp`}
                    control={control}
                    render={({ field }) => (
                      <InputNumber
                        {...field}
                        placeholder="JP per Skill"
                        style={{
                          width: "100%",
                          height: "50px",
                        }}
                      />
                    )}
                  />
                  <Controller
                    name={`hardSkill[${index}].skill_score`}
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

            <div>
              <h2 className="font-Poppins text-2xl font-medium text-center p-6">
                Soft Skills
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
                    name={`softSkill[${index}].description`}
                    control={control}
                    render={({ field }) => (
                      <Input.TextArea
                        value={field.value
                          ?.map(
                            (desc, idx) =>
                              `${idx + 1}. ${desc.unit_code} - ${
                                desc.unit_title
                              }`
                          )
                          .join("\n")}
                        rows={4}
                        placeholder="Unit Code and Title"
                        readOnly
                        style={{
                          marginBottom: "10px",
                          width: "100%",
                        }}
                        onChange={(e) => {
                          // Split input back into array of objects
                          const updatedDescriptions = e.target.value
                            .split("\n")
                            .map((line) => {
                              const [unit_code, ...unit_title] =
                                line.split(" - ");
                              return {
                                unit_code: unit_code?.trim() || "",
                                unit_title:
                                  unit_title.join(" - ")?.trim() || "",
                              };
                            })
                            .filter(
                              (item) => item.unit_code || item.unit_title
                            );

                          field.onChange(updatedDescriptions);
                        }}
                      />
                    )}
                  />

                  <Controller
                    name={`softSkill[${index}].skill_jp`}
                    control={control}
                    render={({ field }) => (
                      <InputNumber
                        {...field}
                        placeholder="JP per Skill"
                        style={{
                          width: "100%",
                          height: "50px",
                        }}
                      />
                    )}
                  />
                  <Controller
                    name={`softSkill[${index}].skill_score`}
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

            <Controller
              name="selectedSignatureId"
              control={control}
              render={({ field }) => (
                <Select
                  placeholder="Pilih Template Paraf"
                  {...field}
                  onChange={(value) => {
                    const selectedSignature = signatureData.find(
                      (signature) => signature._id === value
                    );
                    field.onChange(value);
                    setValue("selectedSignature", selectedSignature); // Sinkronisasi
                    setValue(
                      "config_name",
                      selectedSignature?.config_name || ""
                    );
                    handleSignatureChange(value);
                  }}
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

            <Form.Item label="Nama Display" required>
              <Controller
                name="configName"
                control={control}
                render={({ field }) => (
                  <Input
                    {...field}
                    readOnly
                    style={{ width: "100%", height: "50px" }}
                  />
                )}
              />
            </Form.Item>

            <Form.Item label="Nama penandatangan" required>
              <Controller
                name="namaPenandatangan"
                control={control}
                render={({ field }) => (
                  <Input
                    {...field}
                    readOnly
                    style={{ width: "100%", height: "50px" }}
                  />
                )}
              />
            </Form.Item>

            <Form.Item label="Jabatan Penandatangan" required>
              <Controller
                name="role"
                control={control}
                render={({ field }) => (
                  <Input
                    {...field}
                    readOnly
                    style={{ width: "100%", height: "50px" }}
                  />
                )}
              />
            </Form.Item>

            <Form.Item label="Stamp Perusahaan" required>
              <Controller
                name="stamp"
                control={control}
                render={({ field }) => (
                  <>
                    <Input
                      {...field}
                      readOnly
                      style={{ width: "100%", height: "50px" }}
                    />
                    {field.value && (
                      <div style={{ marginTop: "10px" }}>
                        <img
                          src={field.value}
                          alt="Logo perusahaan"
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

            <Form.Item label="Link logo" required>
              <Controller
                name="linkLogo"
                control={control}
                render={({ field }) => (
                  <>
                    <Input
                      {...field}
                      readOnly
                      style={{ width: "100%", height: "50px" }}
                    />
                    {field.value && (
                      <div style={{ marginTop: "10px" }}>
                        <img
                          src={field.value}
                          alt="Logo perusahaan"
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

            <Form.Item label="Link gambar penandatangan" required>
              <Controller
                name="linkGambarPenandatangan"
                control={control}
                render={({ field }) => (
                  <>
                    <Input
                      {...field}
                      readOnly
                      style={{ width: "100%", height: "50px" }}
                    />
                    {field.value && (
                      <div style={{ marginTop: "10px" }}>
                        <img
                          src={field.value}
                          alt="Gambar penandatangan"
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
              <Button type="primary" htmlType="submit">
                Submit
              </Button>
            </Form.Item>
          </Form>
        </Modal>

        {/* Modal untuk memilih template */}
        <Modal
          title=""
          open={isModalVisible}
          footer={null}
          loading={loading}
          onCancel={() => setIsModalVisible(false)}
          className="rounded-lg p-6 max-w-lg w-full"
          centered
        >
          <div className="flex flex-col items-center space-y-4">
            <Spin spinning={loading}>
              <p className="text-lg font-semibold text-gray-700">
                Silakan pilih template untuk diunduh:
              </p>
              <div className="flex flex-col sm:flex-row space-y-4 sm:space-x-4 sm:space-y-0 w-full">
                <Button
                  type="primary"
                  className="bg-blue-500 hover:bg-blue-600 text-white font-semibold px-4 py-2 rounded-lg w-full sm:w-auto"
                  onClick={() => downloadPDF(selectedDownload?.data_id, "a")}
                >
                  Download Template V1
                </Button>
                <Button
                  type="primary"
                  className="bg-green-500 hover:bg-green-600 text-white font-semibold px-4 py-2 rounded-lg w-full sm:w-auto"
                  onClick={() => downloadPDF(selectedDownload?.data_id, "b")}
                >
                  Download Template V2
                </Button>
              </div>
            </Spin>
          </div>
        </Modal>
      </div>
    </MainLayout>
  );
};

export default Dashboard;
