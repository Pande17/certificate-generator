import axios from "axios";

const Signature = axios.create({
  baseURL: "http://127.0.0.1:3000/api/signature",
});

// Instance untuk kompetensi
const Kompetensi = axios.create({
  baseURL: "http://127.0.0.1:3000/api/competence",
});

// Instance untuk sertifikat
const Sertifikat = axios.create({
  baseURL: "http://127.0.0.1:3000/api/certificate",
});

// Instance untuk login
const Login = axios.create({
  baseURL: "http://127.0.0.1:3000/api/login",
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
