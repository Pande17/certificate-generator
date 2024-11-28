import MainLayout from "../MainLayout/Layout";
import { Kompetensi } from "../api middleware";
import { message, Table, Col, Row, Button, Input, Modal, Form } from "antd";
import { DeleteOutlined, EditOutlined } from "@ant-design/icons";
import { useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";
import { useForm, Controller } from "react-hook-form";

const SignaturePage = () => {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState([]);
  const [searchText, setSearchText] = useState("");
  const [isEditModalVisible, setIsEditModalVisible] = useState(false);
  const [editData, setEditData] = useState(null);
  const { control, handleSubmit, reset } = useForm();

  const navigate = useNavigate();
  const { confirm } = Modal;

  const filteredData = data.filter((item) =>
    item.config_name?.toLowerCase().includes(searchText.toLowerCase())
  );

  useEffect(() => {
    const fetchSignature = async () => {
      setLoading(true);
      try {
        const respons = await Kompetensi.get(
          `http://127.0.0.1:3000/api/signature`
        );
        const datas = respons.data.data;
        const filterData = datas.filter((item) => !item.deleted_at);
        setData(filterData);
      } catch (error) {
        console.log("error", error);
      } finally {
        setLoading(false);
      }
    };
    fetchSignature();
  }, []);

  const delHandle = async (_id) => {
    try {
      await Kompetensi.delete(`/${_id}`);
      setData((prevData) => prevData.filter((item) => item._id !== _id));
      message.success("Data berhasil dihapus");
    } catch (error) {
      console.error("Error response:", error.response);
      message.error(
        `Gagal menghapus data: ${
          error.response?.data?.message || error.message
        }`
      );
    }
  };

  const delConfirm = (_id, config_name) => {
    confirm({
      title: `Apakah Anda yakin ingin menghapus kompetensi ${config_name}?`,
      content: "Data yang dihapus tidak dapat dikembalikan.",
      okType: "danger",
      okText: "Ya, Hapus",
      cancelText: "Batal",
      onOk() {
        delHandle(_id);
      },
    });
  };

  const createNav = () => {
    navigate("/createParaf");
  };

  const handleEdit = (record) => {
    setEditData(record);
    reset(record); // Isi form dengan data yang diedit
    setIsEditModalVisible(true);
  };

  const onSubmit = async (formData) => {
    try {
      const updatedData = {
        ...editData,
        ...formData,
        signature: formData.ttd,
        stamp: formData.Cap,
        name: formData.atasNama,
        config_name: formData.displayNama,
        role: formData.jabatan,
      };

      await Kompetensi.put(`/${editData._id}`, updatedData);

      setData((prevData) =>
        prevData.map((item) => (item._id === editData._id ? updatedData : item))
      );

      message.success("Data berhasil diperbarui");
      setIsEditModalVisible(false);
    } catch (error) {
      console.error("Error response:", error.response);
      message.error(
        `Gagal memperbarui data: ${
          error.response?.data?.message || error.message
        }`
      );
    }
  };

  const columns = [
    {
      title: "Id",
      align: "center",
      width: 100,
      responsive: ["xs", "sm", "md", "lg"],
      render: (text, record, index) => index + 1,
    },
    {
      responsive: ["xs", "sm", "md", "lg"],
      title: "Signature",
      align: "center",
      dataIndex: "config_name",
      key: "config",
      width: 100,
    },
    {
      width: 100,
      responsive: ["xs", "sm", "md", "lg"],
      title: "Aksi",
      align: "center",
      render: (text, record) => (
        <>
          <Button
            icon={<DeleteOutlined />}
            type="primary"
            style={{ margin: 8 }}
            danger
            onClick={() => delConfirm(record._id, record.config_name)}
          />
          <Button
            icon={<EditOutlined />}
            type="primary"
            onClick={() => handleEdit(record)}
            style={{ margin: 8 }}
          />
        </>
      ),
    },
  ];

  return (
    <MainLayout>
      <div className="flex flex-col items-center  p-5">
        <div>
          <p className="text-xl font-Poppins font-semibold mb-5 text-Text p-3 bg-white rounded-xl">
            List Paraf
          </p>
        </div>

        <Button onClick={createNav} className="m-3">
          Buat Sertifikat
        </Button>

        <Input
          placeholder="Search signature"
          value={searchText}
          onChange={(e) => setSearchText(e.target.value)}
          className="mb-4 p-2 border rounded md:w-1/2"
        />

        <Row
          style={{
            justifyContent: "center",
            width: "90%",
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
                x: "max-content",
                y: filteredData.length > 6 ? 500 : undefined,
              }}
              style={{ width: "100%" }}
            />
          </Col>
        </Row>
      </div>

      {/* Modal Edit */}
      <Modal
        title="Edit Data"
        open={isEditModalVisible}
        onCancel={() => setIsEditModalVisible(false)}
        footer={null}
      >
        <Form layout="vertical" onFinish={handleSubmit(onSubmit)}>
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
                <Input
                  {...field}
                  placeholder="Masukkan Jabatan Penandatangan"
                />
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
    </MainLayout>
  );
};

export default SignaturePage;
