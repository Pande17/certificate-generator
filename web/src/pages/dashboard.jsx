import { Table, Button, Modal, message } from "antd";
import axios from "axios";
import { useEffect, useState } from "react";
import {
  DeleteOutlined,
  EditOutlined,
  DownloadOutlined,
  FolderOpenOutlined,
} from "@ant-design/icons";

const { confirm } = Modal;

const Dashboard = () => {
  const [dta, setDta] = useState([]);
  const [loading, setLoading] = useState(false);

  // Fetch data dari API
  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const response = await axios.get(
          `http://127.0.0.1:3000/api/competence`
        );
        setDta(response.data.data);
      } catch (err) {
        console.error("Error fetching data:", err);
        message.error("Gagal memuat data.");
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  // Fungsi untuk menghapus data dari API dan update tabel
  const deleteCompetence = async (kompetensi_id) => {
    try {
      await axios.delete(
        `http://127.0.0.1:3000/api/competence/${kompetensi_id}`
      );
      // Filter data di state dta agar menghapus item yang telah dihapus di API
      setDta((prevDta) =>
        prevDta.filter((item) => item.kompetensi_id !== kompetensi_id)
      );
      message.success("Kompetensi berhasil dihapus!");
    } catch (error) {
      message.error("Gagal menghapus kompetensi.");
    }
  };

  // Fungsi untuk menampilkan konfirmasi sebelum menghapus
  const showDeleteConfirm = (kompetensi_id) => {
    confirm({
      title: "Apakah Anda yakin ingin menghapus kompetensi ini?",
      content: "Data yang dihapus tidak dapat dikembalikan.",
      okText: "Ya, Hapus",
      okType: "danger",
      cancelText: "Batal",
      onOk() {
        deleteCompetence(kompetensi_id); // Panggil fungsi hapus data
      },
      onCancel() {
        console.log("Penghapusan dibatalkan");
      },
    });
  };

  // Kolom untuk tabel
  const columns = [
    {
      title: "ID",
      key: "index",
      align: "center",
      render: (text, record, index) => index + 1, // index dimulai dari 0, jadi tambahkan 1
    },
    {
      title: "Nama Kompetensi",
      dataIndex: "nama_kompetensi",
      key: "nama_kompetensi",
    },
    {
      title: "Aksi",
      key: "actions",
      align: "center",
      render: (text, record) => (
        <div>
          <Button
            icon={<DeleteOutlined />}
            style={{ marginRight: 8 }}
            type="primary"
            danger
            onClick={() => showDeleteConfirm(record.kompetensi_id)}
          />
          <Button
            icon={<FolderOpenOutlined />}
            style={{ marginRight: 8 }}
            onClick={() =>
              message.info(`Buka folder untuk ID ${record.kompetensi_id}`)
            }
          />
          <Button
            icon={<EditOutlined />}
            style={{ marginRight: 8 }}
            type="primary"
            onClick={() =>
              message.info(`Edit kompetensi dengan ID ${record.kompetensi_id}`)
            }
          />
          <Button
            icon={<DownloadOutlined />}
            type="primary"
            onClick={() =>
              message.info(`Unduh kompetensi dengan ID ${record.kompetensi_id}`)
            }
          />
        </div>
      ),
    },
  ];

  return (
    <div className="flex justify-center items-center min-h-screen">
      <Table
        dataSource={dta}
        columns={columns}
        rowKey="kompetensi_id"
        pagination={false}
        bordered
        loading={loading}
      />
    </div>
  );
};

export default Dashboard;
