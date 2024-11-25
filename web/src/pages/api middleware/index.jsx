import axios from "axios";

 const api = process.env.REACT_APP_API_URL;


const Signature = axios.create({
  baseURL: `${api}/api/signature`,
});

// Instance untuk kompetensi
const Kompetensi = axios.create({
  baseURL: `${api}/api/competence`,
});

// Instance untuk sertifikat
const Sertifikat = axios.create({
  baseURL: `${api}/api/certificate`,
});

// Instance untuk login
const Login = axios.create({
  baseURL: `${api}/api/login`,
});

// Interceptor untuk menyisipkan token pada setiap request
const applyAuthInterceptor = (instance) => {
  instance.interceptors.request.use(
    (config) => {
      const authToken = localStorage.getItem("authToken"); // Ambil token dari localStorage
      if (authToken) {
        config.headers.Authorization = `Bearer ${authToken}`; // Tambahkan authToken header
      }
      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );
};

// Terapkan interceptor ke semua instance

applyAuthInterceptor(Signature);
applyAuthInterceptor(Kompetensi);
applyAuthInterceptor(Sertifikat);
applyAuthInterceptor(Login);

// Ekspor semua instance untuk digunakan
export { Kompetensi, Sertifikat, Login ,Signature      };
