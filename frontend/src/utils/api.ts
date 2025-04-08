import axios from 'axios';
const api = axios.create({
    baseURL: 'http://localhost:8080/v1', // gRPC-Gateway endpoint
  });

export default api;