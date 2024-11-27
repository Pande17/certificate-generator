import MainLayout from "../MainLayout/Layout";
import { Signature } from "../api middleware";
import { message, Table, Col, Row, Button, Input, Modal, Form } from "antd";
import { DeleteOutlined, EditOutlined } from "@ant-design/icons";
import { useNavigate} from "react-router-dom";
import { useEffect, useState } from "react";

const SignaturePage = () => {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState([]);
  const [searchText, setSearchText] = useState("");
  const [isEditModalVisible, setIsEditModalVisible] = useState(false);
  const [currentSignature, setCurrentSignature] = useState(null);
  const [dataSignature, setDataSignature] = useState(null);
  const [error, setError] = useState(null);
  const [form] = Form.useForm();

  const navigate = useNavigate();
  const { confirm } = Modal;

  const filteredData = data.filter((item) =>
    item.config_name?.toLowerCase().includes(searchText.toLowerCase())
  );

  
  useEffect(() => {
    const fetchSignature = async () => {
      setLoading(true);
      try {
        const respons = await Signature.get(
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

 
  useEffect(() => {
    if (currentSignature) {
      const fetchData = async () => {
        setLoading(true);
        try {
          const response = await Signature.get(
            `http://127.0.0.1:3000/api/signature/${currentSignature._id}`
          );
          const sigData = response.data.data;
          console.log(sigData);
          if (!sigData.deleted_at) {
            setDataSignature(sigData); 
          } else {
            message.warning("Data sertifikat tidak tersedia.");
          }
        } catch (err) {
          console.error("Error fetching data:", err);
          setError("Gagal memuat data sertifikat.");
        } finally {
          setLoading(false);
        }
      };
      fetchData();
    }
  }, [currentSignature]); 

  const delHandle = async (_id) => {
    try {
      await Signature.delete(`http://127.0.0.1:3000/api/signature/${_id}`);
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
      title: `Apakah anda yakin ingin menghapus kompetensi ${config_name}?`,
      content: "Data yang dihapus tidak dapat dikembalikan",
      okType: "danger",
      okText: "Ya, Hapus",
      cancelText: "Batal",
      onOk() {
        delHandle(_id);
      },
      onCancel() {
        console.log("Penghapusan dibatalkan");
      },
    });
  };

  const createNav = () => {
    navigate("/createParaf");
  };


  const handleEdit = async (record) => {
    try {
      setLoading(true);
      console.log("Fetching data for ID:", record._id);
      const response = await Signature.get(
        `http://127.0.0.1:3000/api/signature/${record._id}` 
      );
      const signatureData = response.data.data;

      if (!signatureData.deleted_at) {
        setCurrentSignature(signatureData); 
        form.setFieldsValue(signatureData); 
        setIsEditModalVisible(true); 
      } else {
        message.warning("Data sertifikat tidak tersedia.");
      }
    } catch (error) {
      console.error("Error fetching detailed data:", error);
      message.error("Gagal memuat data.");
    } finally {
      setLoading(false); 
    }
  };


  const handleSubmit = async (values) => {
    setLoading(true);
    try {
      await Signature.put(
        `http://127.0.0.1:3000/api/signature/${currentSignature._id}`,
        values
      );
      message.success("Data berhasil diperbarui");
      setIsEditModalVisible(false); 
      
      setData((prevData) =>
        prevData.map((item) =>
          item._id === currentSignature._id ? { ...item, ...values } : item
        )
      );
    } catch (error) {
      console.error("Error updating signature:", error);
      message.error("Gagal memperbarui data");
    } finally {
      setLoading(false);
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
      <div className="flex flex-col items-center p-5">
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
          }}>
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

        <Modal
          title="Edit Signature"
          open={isEditModalVisible}
          onCancel={() => setIsEditModalVisible(false)}
          footer={null}>
          <Form form={form} layout="vertical" onFinish={handleSubmit}>
            <Form.Item label="Display Nama" name="config_name">
              <Input placeholder="Masukkan nama display" />
            </Form.Item>
            <Form.Item label="Nama Penandatangan" name="name">
              <Input placeholder="Masukkan nama penandatangan" />
            </Form.Item>
            <Form.Item label="Jabatan Penandatangan" name="role">
              <Input placeholder="Masukkan jabatan penandatangan" />
            </Form.Item>
            <Form.Item label="Link Gambar Tanda Tangan" name="signature">
              <Input placeholder="Masukkan link tanda tangan" />
            </Form.Item>
            <Form.Item label="Link Gambar Cap Perusahaan" name="stamp">
              <Input placeholder="Masukkan link cap perusahaan" />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit">
                Simpan
              </Button>
            </Form.Item>
          </Form>
        </Modal>
      </div>
    </MainLayout>
  );
};

export default SignaturePage;
