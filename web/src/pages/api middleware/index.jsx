import axios from "axios";

const Signature = axios.create({
  baseURL: `api/signature`,
});

const Kompetensi = axios.create({
  baseURL: `api/competence`,
});

const Sertifikat = axios.create({
  baseURL: `api/certificate`,
});

const Login = axios.create({
  baseURL: `api/login`,
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
    "X-Requested-With": "XMLHttpRequest",
  },
});

// Interceptor for adding token
const applyAuthInterceptor = (instance) => {
  instance.interceptors.request.use(
    (config) => {
      const authToken = localStorage.getItem("authToken");
      if (authToken) {
        config.headers.Authorization = `Bearer ${authToken}`;
      }
      return config;
    },
    (error) => Promise.reject(error)
  );
};

applyAuthInterceptor(Signature);
applyAuthInterceptor(Kompetensi);
applyAuthInterceptor(Sertifikat);
applyAuthInterceptor(Login);

export { Kompetensi, Sertifikat, Login, Signature };
