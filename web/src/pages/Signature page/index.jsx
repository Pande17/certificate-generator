import MainLayout from "../MainLayout/Layout";
import { Kompetensi } from "../api middleware";
import { message, Table, Col, Row, Button, Input, Modal } from "antd";
import { DeleteOutlined, EditOutlined } from "@ant-design/icons";
import { useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";


const SignaturePage = () => {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState([]);
  const [searchText, setSearchText] = useState("");

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
     await Kompetensi.delete(`http://127.0.0.1:3000/api/signature/${_id}`);
     setData((prevData) => prevData.filter((item) => item._id !== _id));
     message.success("Data berhasil dihapus");
   } catch (error) {
     console.error("Error response:", error.response);
     message.error(
       `Gagal menghapus data: ${error.response?.data?.message || error.message}`
     );
   }
 };

   const delConfirm = (_id, config_name) => {
     confirm({
       title: `apakah anda yakon ingin menghapus kompetensi ${config_name}`,
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


  const handleEdit = (record) => {
    message.info(`Edit triggered for ${record._id}`);
    // Implement edit logic
  };

  const columns = [
    {
      title: "Id",
      align: "center",
      width: 100,
      render: (text, record, index) => index + 1,
    },
    {
      title: "Signature",
      align: "center",
      dataIndex: "config_name",
      key: "config",
    },
    {
      title: "Aksi",
      align: "center",
      render: (text, record) => (
        <>
          <Button
            icon={<DeleteOutlined />}
            type="primary"
            danger
            onClick={() =>delConfirm(record._id , record.config_name)}
          />
          <Button
            icon={<EditOutlined />}
            type="primary"
            onClick={() => handleEdit(record)}
            style={{ marginLeft: 8 }}
          />
        </>
      ),
    },
  ];

  return (
    <MainLayout>
      <div className="flex flex-col items-center  p-5">
        <div>
          <Input
            placeholder="Search signature"
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
            className="mb-4 p-2 border rounded w-full md:w-1/2"
          />
        </div>

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
                x: "max-content",
                y: filteredData.length > 6 ? 500 : undefined,
              }}
              style={{ width: "100%" }}
            />
          </Col>
        </Row>
      </div>
    </MainLayout>
  );
};

export default SignaturePage;
