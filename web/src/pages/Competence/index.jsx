import { useEffect, useState } from "react"
import axios from "axios";
import{message, Table, Col , Row, Button, Modal} from "antd"
import {
  DeleteOutlined,
  EditOutlined,

} from "@ant-design/icons";
import MainLayout from "../MainLayout/Layout"
import { Navigate, useNavigate } from "react-router-dom";


const competence = () => {
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState([])
    const [searchText, setSearchText] = useState("");

const {confirm} = Modal
const navigate = useNavigate()

    useEffect(() => {
    const fetchingData = async() => {
      setLoading(true)
      try{
        const response = await axios.get(
          `http://127.0.0.1:3000/api/competence`
        );
        const datas = response.data.data
          const filteredData = datas.filter((item) => !item.deleted_at);
        setData(filteredData);
      
      }catch(err){
        console.error('error : ', err)
        message.error('error : ', err)
      }finally {
         setLoading(false)
      }
    }
    fetchingData()
  },[])

  const filteredData = data.filter((item) =>
    item.nama_kompetensi.toLowerCase().includes(searchText.toLowerCase())
  );
  const handleSearch = (e) => {
    setSearchText(e.target.value);
  };

 const delHandle = async (_id) => {
   try {
   await axios.delete(
       `http://127.0.0.1:3000/api/competence/${_id}`
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


  const delConfirm = (_id, nama_kompetensi) => {
    confirm({
      title:`apakah anda yakon ingin menghapus kompetensi ${nama_kompetensi}`,
      content:"data yang di hapus tidak dapat dikembalikan",
      okType:"danger",
      okText:"ya, Hapus",
      cancelText:"Batal",
      onOk() {
        delHandle(_id)
      },
      onCancel() {
        console.log("penghapusan dibatalkan")
      }
    })
  }
  const column = [
    {
      title: "Id",
      align: "center",
      width: 100,
      ellipsis: true,
      responsive: ["xs", "sm", "md", "lg"],
      render: (text, record, index) => index + 1,
    },
    {
      title: "Kompetensi",
      align: "center",
      width: 100,
      responsive: ["xs", "sm", "md", "lg"],
      key: "nama_kompetensi",
      dataIndex: "nama_kompetensi",
      ellipsis: true,
    },
    {
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
              onClick={() => message.info(`Edit in progress`)}
            />
          </div>
        );
      },
    },
  ];


  const createNav = () => {
    navigate("/competence/create-competence");
  }
  return (
    <MainLayout>
      <div className="flex flex-col items-center  p-5">
        <div>
          <p className="text-xl font-Poppins font-semibold mb-5 text-Text p-3 bg-white rounded-xl">
            List Kompetensi
          </p>
        </div>
        <input
          type="text"
          placeholder="Search"
          value={searchText}
          onChange={handleSearch}
          className="mb-4 p-2 border border-gray-300 rounded w-full md:w-1/2"
        />
        <Button onClick={createNav} className="m-3">
          Create Competence
        </Button>
        <Row style={{justifyContent:"center", width: "100%", overflowX: "auto" }}>
          <Col >
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
              style={{ width: "100%" }}
            />
          </Col>
        </Row>
      </div>
    </MainLayout>
  );
}

export default competence ; 
