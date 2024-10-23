import axios from "axios";
import { useEffect, useState } from "react";
import boxicons from 'boxicons';

const Dashboard = () => {
  const[dta, setDta] = useState([]);

  useEffect(()=>{
    const fetchdata = async() =>{
      try{

        const respons = await axios.get(`http://127.0.0.1:3000/api/competence`);
        const isidata = respons.data.data
        setDta(isidata)
      } catch (err){
        console.log("error bg : ",err)
      } 
    }
    fetchdata();
  }, [])
  return (
    <div className="flex justify-center items-center min-h-screen">
      <table className="min-w-full table-auto border-collapse">
        <thead>
          <tr>
            <th className="border px-4 py-2">ID</th>
            <th className="border px-4 py-2">Name</th>
            <th className="border px-4 py-2">Actions</th>
          </tr>
        </thead>
        <tbody>
          {dta.map((data) => (
            <tr key={data.kompetensi_id}>
              <td className="text-center border">
                <p>{data.kompetensi_id}</p>
              </td>
              <td className="border">
                <p>{data.nama_kompetensi}</p>
              </td>
              <td className="text-center">
                <button className="border border-black rounded-md mx-3">
                  <box-icon type="solid" name="trash"></box-icon>
                </button>
                <button className="border border-black rounded-md mx-3">
                  <box-icon name="folder-open" type="solid"></box-icon>
                </button>
                <button className="border border-black rounded-md mx-3">
                  <box-icon name="edit"></box-icon>
                </button>
                <button className="border border-black rounded-md mx-3">
                  <box-icon type="solid" name="download"></box-icon>
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default Dashboard;
