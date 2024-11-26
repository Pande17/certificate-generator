import axios from "axios";

// Use fallback if REACT_APP_API_URL is undefined
const api = import.meta.env.REACT_APP_API_URL || "http://localhost:3000";

const Signature = axios.create({
  baseURL: `${api}/api/signature`,
});

const Kompetensi = axios.create({
  baseURL: `${api}/api/competence`,
});

const Sertifikat = axios.create({
  baseURL: `${api}/api/certificate`,
});

const Login = axios.create({
  baseURL: `${api}/api/login`,
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
