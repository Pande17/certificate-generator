  import { Table, Button, Modal, message, Row, Col, Input, Form } from "antd";
  import axios from "axios";
  import { useEffect, useState } from "react";
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

    // Fetch data dari API
    useEffect(() => {
      const fetchData = async () => {
        setLoading(true);
        try {
          const response = await axios.get(
            "http://127.0.0.1:3000/api/certificate"
          );
          const certificates = response.data.data || [];
          const filteredData = certificates.filter((item) => !item.deleted_at);
          setDta(filteredData);
        } catch (err) {
          console.error("Error fetching data:", err);
          message.error("Gagal memuat data.");
        } finally {
          setLoading(false);
        }
      };
      fetchData();
    }, []);

    const deleteCompetence = async (_id) => {
      try {
        await axios.delete(`http://127.0.0.1:3000/api/certificate/${_id}`);
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

    const handleEdit = (record) => {
      setCurrentRecord(record);
      setIsEditModalVisible(true);
    };

    const handleEditSubmit = async (values) => {
      try {
        await axios.put(
          `http://127.0.0.1:3000/api/certificate/${currentRecord._id}`,
          values
        );
        setDta((prevDta) =>
          prevDta.map((item) =>
            item._id === currentRecord._id ? { ...item, ...values } : item
          )
        );
        message.success("SERTIFIKAT berhasil diperbarui!");
        setIsEditModalVisible(false);
        setCurrentRecord(null);
      } catch (error) {
        message.error("Gagal memperbarui SERTIFIKAT.");
      }
    };

    const filteredData = dta.filter((item) =>
      item.sertif_name.toLowerCase().includes(searchText.toLowerCase())
    );

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
              onClick={() =>
                message.info(`Unduh sertifikat dengan ID ${record._id}`)
              }
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
              History list
            </p>
          </div>
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
          open  ={isEditModalVisible}
          onCancel={() => setIsEditModalVisible(false)}
          footer={null}
        >
          <Form
            initialValues={currentRecord}
            onFinish={handleEditSubmit}
            layout="vertical"
          >
            <Form.Item
              name="sertif_name"
              label="Nama Sertifikat"
              rules={[{ required: true, message: "Nama sertifikat wajib!" }]}
            >
              <Input />
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

export default Dashboard;
