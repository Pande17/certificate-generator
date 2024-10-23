import axios from "axios";
import { useEffect, useState } from "react";

const CreatePage = () => {
  const [isiDta, setIsiDta] = useState([]);
  const [pilihan, setPilihan] = useState(""); // Ubah ke string
  const [kompetensi, setKompetensi] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchdata = async () => {
      try {
        const respon = await axios.get("http://127.0.0.1:3000/api/competence");
        const data1 = respon.data.data;

        if (Array.isArray(data1)) {
          setIsiDta(data1);
        } else {
          console.error("Data fetched is not an array");
        }
      } catch (err) {
        console.log("Error during data fetching");
      }
    };
    fetchdata();
  }, []);

  const handleisidata = async (e) => {
    const isiPilihan = e.target.value;
    setPilihan(isiPilihan);

    const pilihanKompetensi = isiDta.find(
      (isiDta) => isiDta.nama_kompetensi === isiPilihan
    );

    if (pilihanKompetensi) {
      setLoading(true);
      try {
        const respons = await axios.get(
          `http://127.0.0.1:3000/api/competence/${pilihanKompetensi.kompetensi_id}`
        );
        setKompetensi(respons.data.data);
      } catch (err) {
        setError(err.message); // Ubah ke .message
      } finally {
        setLoading(false);
      }
    }
  };

  if (loading) {
    return <p>Loading... Please wait.</p>;
  }

  if (error) {
    return <p>Error: {error}</p>;
  }

  return (
    <>
      <h1>Select Competence</h1>
      <select id="dropdown" value={pilihan} onChange={handleisidata}>
        <option value="">Select a Competence</option>
        {Array.isArray(isiDta) &&
          isiDta.map((f) => (
            <option key={f.kompetensi_id} value={f.nama_kompetensi}>
              {f.nama_kompetensi}
            </option>
          ))}
      </select>

      {kompetensi && (
        <div>
          <h2>Competence: {kompetensi.nama_kompetensi}</h2>
          <h3>Hard Skills</h3>
          {kompetensi.hard_skills && kompetensi.hard_skills.length >0 ? (
            <ul>
                {kompetensi.hard_skills.map((skill, i) => (
                    <li key={i}>
                        <h4>{skill.hardskill_name || 'unknown'}</h4>
                        {skill.description.map((desk, d)=> (
                         <li key={d}><strong>{desk.unit_code} : {desk.unit_title}</strong></li>   
                        ))}
                    </li>
                ))};
            </ul>
          ):(
            <p> ??????</p>
          )}

          
        </div>
      )}
    </>
  );
};

export default CreatePage;
