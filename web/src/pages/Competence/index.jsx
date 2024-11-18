import { useEffect, useState } from 'react';
import axios from 'axios';
import { message } from 'antd';

const Competence = () => {
	const [loading, setLoading] = useState(false);
	const [data, setData] = useState([]);

	useEffect(() => {
		const fetchingData = async () => {
			setLoading(true);
			try {
				const response = await axios.get(`http://127.0.0.1:3000/api/competence`);
				const fetchedData = response.data.data; // Perbaikan: Simpan data ke state
				setData(fetchedData); // Set data ke state
			} catch (err) {
				console.error('error : ', err);
				message.error('error : ', err.message || err);
			} finally {
				setLoading(false);
			}
		};

		fetchingData(); // Pastikan useEffect memanggil fungsi ini
	}, []); // Tambahkan dependency array agar hanya dijalankan sekali saat komponen dimuat

	return (
		<div>
			<h1>Competence Page</h1>
			{loading ? (
				<p>Loading...</p>
			) : (
				<ul>
					{data.map((item) => (
						<li key={item.id}>{item.name}</li>
					))}
				</ul>
			)}
		</div>
	);
};

// Tambahkan default export di sini
export default Competence;

// import { useEffect, useState } from 'react';
// import axios from 'axios';
// import { message } from 'antd';

// const competence = () => {
// 	const [loading, setLoading] = useState(false);
// 	const [data, setData] = useState([]);

// 	useEffect(() => {
// 		const fetchingData = async () => {
// 			setLoading(true);
// 			try {
// 				const response = await axios.get(`http://127.0.0.1:3000/api/competence`);
// 				const data = response.data.data;
// 			} catch (err) {
// 				console.error('error : ', err);
// 				message.error('error : ', err);
// 			} finally {
// 				setLoading(false);
// 			}
// 		};
// 	});
// };
